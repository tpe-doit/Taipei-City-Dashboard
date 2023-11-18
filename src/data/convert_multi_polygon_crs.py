import geopandas as gpd
from pyproj import Proj, Transformer
from shapely.ops import transform

toushui_list = ["PAC", "parking", "road", "school", "sidewalk"]


def transform_geojson(input_file, output_file, original_crs, target_crs):
    gdf = gpd.read_file(input_file)

    if gdf.crs is None:
        gdf.set_crs(original_crs, inplace=True)

    transformer = Transformer.from_crs(original_crs, target_crs, always_xy=True)
    gdf["geometry"] = gdf["geometry"].apply(
        lambda geom: transform(transformer.transform, geom)
    )

    gdf.set_crs(target_crs, allow_override=True)
    gdf.to_file(output_file, driver="GeoJSON")


# for file in toushui_list:
transform_geojson(
    "/Users/oohyuti/Taipei-City-Dashboard-FE/public/mapData/K1001-3.geojson",
    "/Users/oohyuti/Taipei-City-Dashboard-FE/public/mapData/toushui/output_K1001-3.geojson",
    "epsg:3826",
    "epsg:4326",
)
