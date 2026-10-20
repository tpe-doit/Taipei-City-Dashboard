from airflow import DAG
from operators.common_pipeline import CommonDag

def _ensure_ready_table(engine, table_name, col_map):
    """建表(若不存在)。冪等,每次 DAG run 都會跑。"""
    from sqlalchemy.sql import text as sa_text
    from utils.generate_sql_to_create_DB_table import generate_sql_to_create_db_table

    sql = generate_sql_to_create_db_table(table_name, col_map)
    sql = sql.replace(
        f"CREATE TRIGGER {table_name}_mtime",
        f"DROP TRIGGER IF EXISTS {table_name}_mtime ON public.{table_name};\n"
        f"    CREATE TRIGGER {table_name}_mtime",
    )
    with engine.connect() as conn:
        conn.execute(sa_text(sql).execution_options(autocommit=True))

def _pharmacies(**kwargs):
    # === Imports(全部寫在函式內)===
    import pandas as pd
    from sqlalchemy import create_engine
    from airflow.models import Variable
    from utils.extract_stage import get_data_taipei_api
    from utils.load_stage import (
        save_geodataframe_to_postgresql,
        update_lasttime_in_data_to_dataset_info,
    )
    from utils.transform_address import get_addr_xy_parallel
    from utils.transform_geometry import add_point_wkbgeometry_column_to_df

    # === Config ===
    ready_data_db_uri = kwargs.get("ready_data_db_uri")
    dag_infos = kwargs.get("dag_infos")
    dag_id = dag_infos.get("dag_id")
    load_behavior = dag_infos.get("load_behavior")
    default_table = dag_infos.get("ready_data_default_table")
    history_table = dag_infos.get("ready_data_history_table")

    COL_MAP = {
        "data_time": "timestamp with time zone DEFAULT CURRENT_TIMESTAMP",
        "seq": "integer",
        "pharmacy_name": "text",
        "zipcode": "text",
        "address": "text",
        "phone": "text",
        "lng": "double precision",
        "lat": "double precision",
        "wkb_geometry": "geometry(Point,4326)",
    }
    SELECT_COLUMNS = list(COL_MAP.keys())

    # === Extract ===
    raw_data = pd.DataFrame(get_data_taipei_api("42cfc382-f2b8-4c3a-87ad-37249634f78e"))
    # === Transform ===
    data = raw_data.rename(columns={
        "_id": "seq",
        "機構名稱": "pharmacy_name",
        "地址": "address",
        "電話": "phone",
        "經度": "lng",
        "緯度": "lat"
    })
    if "phone" not in data.columns:
        data["phone"] = None
    if "lng" not in data.columns:
        data["lng"] = None
    if "lat" not in data.columns:
        data["lat"] = None
    data["zipcode"] = None
    data["data_time"] = pd.to_datetime("now")

    # === Geocode（地址→經緯度;若原始資料缺座標則打 TPGOS 補點）===
    data["lng"] = pd.to_numeric(data["lng"], errors="coerce")
    data["lat"] = pd.to_numeric(data["lat"], errors="coerce")
    if Variable.get("TPGOS_GET_ADDR_XY", default_var=None):
        missing_mask = data["lng"].isna() | data["lat"].isna()
        uniq = data.loc[missing_mask, "address"].dropna().drop_duplicates().tolist()
        if uniq:
            lng, lat = get_addr_xy_parallel(uniq, sleep_time=0.5)
            addr_xy = pd.DataFrame({"address": uniq, "geo_lng": lng, "geo_lat": lat})
            data = data.merge(addr_xy, on="address", how="left")
            data["lng"] = data["lng"].fillna(pd.to_numeric(data["geo_lng"], errors="coerce"))
            data["lat"] = data["lat"].fillna(pd.to_numeric(data["geo_lat"], errors="coerce"))
            data = data.drop(columns=["geo_lng", "geo_lat"], errors="ignore")
    # 經緯度(WGS84)→ Point wkb_geometry
    data = add_point_wkbgeometry_column_to_df(
        data, x=data["lng"], y=data["lat"], from_crs=4326, to_crs=4326, is_add_xy_columns=True
    )
    # keep-all-rows:對不到座標的列保留,wkb_geometry 設 null(非空點)
    data.loc[data["lng"].isna() | data["lat"].isna(), "wkb_geometry"] = None
    data = data[SELECT_COLUMNS]

    # === Load ===
    engine = create_engine(ready_data_db_uri)
    _ensure_ready_table(engine, default_table, COL_MAP)
    save_geodataframe_to_postgresql(
        engine,
        gdata=data,
        load_behavior=load_behavior,
        geometry_type="Point",
        default_table=default_table,
        history_table=history_table,
    )
    update_lasttime_in_data_to_dataset_info(
        engine, dag_id, data["data_time"].max()
    )

dag = CommonDag(proj_folder="proj_city_dashboard", dag_folder="pharmacies")
dag.create_dag(etl_func=_pharmacies)
