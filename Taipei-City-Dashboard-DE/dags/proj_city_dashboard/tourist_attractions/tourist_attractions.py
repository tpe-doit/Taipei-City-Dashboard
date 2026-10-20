from airflow import DAG
from operators.common_pipeline import CommonDag


def _transfer(**kwargs):
    import pandas as pd
    import requests
    from sqlalchemy import create_engine
    from utils.extract_stage import get_tdx_data
    from utils.load_stage import (
        save_geodataframe_to_postgresql,
        update_lasttime_in_data_to_dataset_info,
    )
    from utils.transform_geometry import add_point_wkbgeometry_column_to_df

    # Config
    ready_data_db_uri = kwargs.get("ready_data_db_uri")
    dag_infos = kwargs.get("dag_infos")
    dag_id = dag_infos.get("dag_id")
    load_behavior = dag_infos.get("load_behavior")
    default_table = dag_infos.get("ready_data_default_table")
    history_table = dag_infos.get("ready_data_history_table")
    # Load
    GEOMETRY_TYPE = "Point"
    TPE_URL = "https://tdx.transportdata.tw/api/basic/v2/Tourism/ScenicSpot/Taipei?%24format=JSON"

    df = None
    try:
        res = get_tdx_data(TPE_URL, output_format='dataframe')
        df = res.copy()
        df["longitude"] = df["Position"].apply(lambda pos: pos.get("PositionLon") if isinstance(pos, dict) else None)
        df["latitude"] = df["Position"].apply(lambda pos: pos.get("PositionLat") if isinstance(pos, dict) else None)
        zipcode_to_area = {
            "100": "中正區",
            "103": "大同區",
            "104": "中山區",
            "105": "松山區",
            "106": "大安區",
            "108": "萬華區",
            "110": "信義區",
            "111": "士林區",
            "112": "北投區",
            "114": "內湖區",
            "115": "南港區",
            "116": "文山區"
        }
        df['distric'] = df['ZipCode'].astype(str).map(zipcode_to_area)
        df = df.rename(columns={
            "ScenicSpotName": "name",
            "DescriptionDetail": "introduction",
            "Phone": "tel",
            "Class1": "type",
            "City": "city"
        })
    except Exception as e:
        print(f"TDX API failed with {e}, falling back to travel.taipei open API")
        headers = {
            "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
            "Accept": "application/json",
            "Referer": "https://www.travel.taipei/",
        }
        all_data = []
        page = 1
        while True:
            url = f"https://www.travel.taipei/open-api/zh-tw/Attractions/All?page={page}"
            res = requests.get(url, headers=headers, timeout=30)
            if res.status_code != 200:
                break
            items = res.json().get("data", [])
            if not items:
                break
            all_data.extend(items)
            page += 1

        if not all_data:
            print("No data retrieved from travel.taipei")
            return

        df = pd.DataFrame(all_data)
        df["longitude"] = pd.to_numeric(df["elong"], errors="coerce")
        df["latitude"] = pd.to_numeric(df["nlat"], errors="coerce")
        df["type"] = df["category"].apply(
            lambda c: c[0]["name"] if isinstance(c, list) and len(c) > 0 and isinstance(c[0], dict) else None
        )
        df["city"] = "臺北市"

    gdata = add_point_wkbgeometry_column_to_df(
            df, x=df["longitude"], y=df["latitude"], from_crs=4326
        )

    df = gdata[["name", "type", "introduction", "city", "distric", "tel", "longitude", "latitude", "wkb_geometry"]]
    df["data_time"] = pd.to_datetime("now").strftime("%Y-%m-%d %H:%M:%S")

    engine = create_engine(ready_data_db_uri)
    save_geodataframe_to_postgresql(
        engine,
        gdata=df,
        load_behavior=load_behavior,
        default_table=default_table,
        history_table=history_table,
        geometry_type=GEOMETRY_TYPE,
    )
    update_lasttime_in_data_to_dataset_info(
            engine, dag_id, df["data_time"].max()
        )

dag = CommonDag(proj_folder="proj_city_dashboard", dag_folder="tourist_attractions")
dag.create_dag(etl_func=_transfer)
