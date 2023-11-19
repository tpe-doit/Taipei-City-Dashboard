import geopandas as gpd
import os

os.environ["SHAPE_RESTORE_SHX"] = "YES"


def convert_shp_to_geojson(input_shp, output_geojson):
    gdf = gpd.read_file(input_shp, SHAPE_RESTORE_SHX="yes")
    gdf.to_file(output_geojson, driver="GeoJSON")


convert_shp_to_geojson(
    "/Users/oohyuti/Taipei-City-Dashboard-FE/public/mapData/water/臺北市歷史積淹水範圍(1121023更新).shp",
    "/Users/oohyuti/Taipei-City-Dashboard-FE/public/mapData/history_flood.geojson",
)
