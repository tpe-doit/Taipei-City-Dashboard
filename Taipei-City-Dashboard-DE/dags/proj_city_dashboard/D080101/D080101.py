import logging

from airflow import DAG
from requests.exceptions import RequestException, Timeout

from operators.common_pipeline import CommonDag


def _D080101(**kwargs):
    from proj_city_dashboard.D080101.cdc_visit_case_etl import cdc_visit_case_etl
    
    URL = "https://od.cdc.gov.tw/eic/NHI_AcuteUpperRespiratoryInfections.csv"
    HTTP_FALLBACK_URL = "http://od.cdc.gov.tw/eic/NHI_AcuteUpperRespiratoryInfections.csv"
    DIEASE_NAME = '急性上呼吸道感染'
    
    try:
        cdc_visit_case_etl(URL, DIEASE_NAME, **kwargs)
    except (Timeout, RequestException, Exception) as e:
        logging.warning(f"CDC HTTPS connection failed ({e}); using HTTP CSV fallback")
        cdc_visit_case_etl(HTTP_FALLBACK_URL, DIEASE_NAME, **kwargs)


dag = CommonDag(proj_folder="proj_city_dashboard", dag_folder="D080101")
dag.create_dag(etl_func=_D080101)
