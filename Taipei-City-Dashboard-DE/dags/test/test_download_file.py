"""Offline regression tests for downloads, without requiring an Airflow runtime."""

import ast
from pathlib import Path
import tempfile
import traceback
import unittest
from unittest.mock import Mock, patch

import requests


SOURCE_PATH = Path(__file__).parents[1] / "utils" / "extract_stage.py"
source_tree = ast.parse(SOURCE_PATH.read_text(encoding="utf-8"))
download_node = next(
    node for node in source_tree.body
    if isinstance(node, ast.FunctionDef) and node.name == "download_file"
)
namespace = {"requests": requests, "DATA_PATH": "", "PROXIES": {}}
exec(
    compile(ast.Module(body=[download_node], type_ignores=[]), str(SOURCE_PATH), "exec"),
    namespace,
)
download_file = namespace["download_file"]


class DownloadFileTests(unittest.TestCase):
    secret = "FAKE_DOWNLOAD_CREDENTIAL_ONLY"
    url = f"https://example.invalid/data?Authorization={secret}&format=JSON"

    def setUp(self):
        self.tempdir = tempfile.TemporaryDirectory()
        self.addCleanup(self.tempdir.cleanup)

    def make_response(self, status=200):
        response = requests.Response()
        response.status_code = status
        response.url = self.url
        response.reason = "Download response"
        response.request = requests.Request("GET", self.url).prepare()
        response.raw = Mock()
        response.iter_content = Mock(return_value=iter([b"weather", b" data"]))
        return response

    def assert_safe_failure(self, callback, expected_type):
        try:
            callback()
        except expected_type as error:
            self.assertIs(type(error), expected_type)
            self.assertNotIn(self.secret, str(error))
            self.assertNotIn(self.secret, repr(error))
            formatted = "".join(traceback.format_exception(
                type(error), error, error.__traceback__,
            ))
            self.assertNotIn(self.secret, formatted)
            self.assertIsNone(error.request)
            self.assertIsNone(error.response)
            self.assertIsNone(error.__context__)
            self.assertIsNone(error.__cause__)
            return str(error)
        self.fail(f"Expected {expected_type.__name__}")

    def download(self):
        return download_file("weather.json", self.url, file_folder=self.tempdir.name)

    def test_success_preserves_bytes_and_request_options(self):
        response = self.make_response()
        with patch.object(requests, "get", return_value=response) as get:
            result = self.download()
        self.assertEqual(Path(result).read_bytes(), b"weather data")
        get.assert_called_once_with(
            self.url, stream=True, proxies=None, verify=True, timeout=60,
        )

    def test_http_error_retains_type_and_status_without_response_url(self):
        with patch.object(requests, "get", return_value=self.make_response(401)):
            message = self.assert_safe_failure(self.download, requests.exceptions.HTTPError)
        self.assertIn("401", message)

    def test_connection_failures_do_not_retain_nested_credential_context(self):
        for error_type in [
            requests.exceptions.Timeout,
            requests.exceptions.ConnectTimeout,
            requests.exceptions.ReadTimeout,
            requests.exceptions.ConnectionError,
            requests.exceptions.ProxyError,
            requests.exceptions.SSLError,
            requests.exceptions.TooManyRedirects,
            requests.exceptions.InvalidURL,
            requests.exceptions.InvalidSchema,
            requests.exceptions.MissingSchema,
        ]:
            with self.subTest(error_type=error_type.__name__):
                original = error_type(self.url, request=self.make_response().request)
                original.__cause__ = OSError(self.url)
                original.__context__ = OSError(self.url)
                with patch.object(requests, "get", side_effect=original):
                    self.assert_safe_failure(self.download, error_type)

    def test_stream_failures_do_not_expose_request_or_response(self):
        for error_type in [
            requests.exceptions.ChunkedEncodingError,
            requests.exceptions.ContentDecodingError,
            requests.exceptions.StreamConsumedError,
        ]:
            with self.subTest(error_type=error_type.__name__):
                response = self.make_response()
                def fail_during_stream(chunk_size):
                    yield b"partial chunk"
                    raise error_type(
                        self.url, request=response.request, response=response,
                    )
                response.iter_content = fail_during_stream
                with patch.object(requests, "get", return_value=response):
                    self.assert_safe_failure(self.download, error_type)

    def test_file_errors_preserve_their_original_type(self):
        with patch.object(requests, "get", return_value=self.make_response()):
            with self.assertRaises(FileNotFoundError):
                download_file(
                    "weather.json", self.url,
                    file_folder=str(Path(self.tempdir.name) / "missing"),
                )


if __name__ == "__main__":
    unittest.main()
