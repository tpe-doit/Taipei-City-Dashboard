import geopandas as gpd
from pyproj import Transformer
from shapely.geometry import MultiLineString
from shapely.ops import transform


def transform_multilinestring(mls, transformer):
    transformed_lines = [transform(transformer.transform, line) for line in mls]
    return MultiLineString(transformed_lines)


def transform_geometry(geom, transformer):
    if geom is None:
        return None
    if isinstance(geom, MultiLineString):
        return transform_multilinestring(geom, transformer)
    else:
        return transform(transformer.transform, geom)


def transform_geojson(input_file, output_file, original_crs, target_crs):
    gdf = gpd.read_file(input_file)

    if gdf.crs is None:
        gdf = gdf.set_crs(original_crs)

    transformer = Transformer.from_crs(original_crs, target_crs, always_xy=True)

    gdf["geometry"] = gdf["geometry"].apply(
        lambda geom: transform_geometry(geom, transformer)
    )

    gdf = gdf.set_crs(target_crs, allow_override=True)

    gdf.to_file(output_file, driver="GeoJSON")


transform_geojson(
    "/Users/oohyuti/Taipei-City-Dashboard-FE/public/mapData/K1001-3.geojson",
    "/Users/oohyuti/Taipei-City-Dashboard-FE/public/mapData/toushui/output_K1001-3.geojson",
    "epsg:3826",
    "epsg:4326",
)
