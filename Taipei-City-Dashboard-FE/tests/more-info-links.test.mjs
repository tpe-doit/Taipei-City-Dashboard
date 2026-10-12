import assert from "node:assert/strict";
import { readFileSync } from "node:fs";
import test from "node:test";
import vm from "node:vm";

// Exercise the component's actual link functions without requiring a DOM runner.
const source = readFileSync(
	new URL("../src/components/dialogs/MoreInfo.vue", import.meta.url),
	"utf8",
);
const script = source.match(/<script setup>([\s\S]*?)<\/script>/)[1];
const { getSafeLink, getLinkTag } = vm.runInNewContext(
	`${script.slice(script.indexOf("function getSafeLink"))}; ({ getSafeLink, getLinkTag });`,
	{ URL },
);

test("keeps HTTP/S links and existing labels for the actual data hosts", () => {
	const cases = [
		["https://data.taipei/dataset/detail", "資料集 - 2 (data.taipei)"],
		["https://data.ntpc.gov.tw/datasets", "資料集 - 2 (data.ntpc)"],
		["https://tuic.gov.taipei/projects", "大數據中心專案網頁"],
		["https://github.com/taipei-doit/Taipei-City-Dashboard", "GitHub 程式庫"],
		["http://DATA.TAIPEI:80/dataset", "資料集 - 2 (data.taipei)"],
	];
	for (const [link, label] of cases) {
		assert.equal(getLinkTag(link, 1), label);
		assert.equal(getSafeLink(link).href, new URL(link).href);
	}
});

test("does not label substring lookalikes as a trusted data host", () => {
	for (const host of ["data.taipei", "data.ntpc.gov.tw", "tuic.gov.taipei", "github.com"]) {
		for (const link of [
			`https://evil.example/${host}`,
			`https://evil.example/?source=${host}`,
			`https://evil.example/#${host}`,
			`https://${host}.evil.example/`,
			`https://${host}@evil.example/`,
		]) {
			assert.equal(getLinkTag(link, 0), "資料集 - 1 (其他)");
		}
	}
});

test("leaves unsafe schemes and malformed values without a navigable href", () => {
	for (const link of [
		"javascript:alert('github.com')",
		"JaVaScRiPt:alert(1)",
		"data:text/html,<script>alert(1)</script>",
		"file:///etc/passwd",
		"//github.com/taipei-doit",
		"/github.com",
		"not a URL",
		"",
		null,
		42,
	]) {
		assert.equal(getSafeLink(link), null);
		assert.equal(getLinkTag(link, 0), "資料集 - 1 (其他)");
	}
	assert.match(source, /:href="getSafeLink\(link\)\?\.href"/);
});
