import matplotlib.pyplot as plt
import matplotlib.colors as mcolors
import pandas as pd
import math
import json
from shapely.geometry import Point
import geopandas as gpd

building_df = pd.read_csv(
    "/Users/oohyuti/Taipei-City-Dashboard-FE/public/csv/危樓經緯度.csv"
)
output_file = (
    "/Users/oohyuti/Taipei-City-Dashboard-FE/public/csv/dangerous_building_area.geojson"
)

area_list = list(building_df["面積"])

normalized_numbers = mcolors.Normalize(min(area_list), max(area_list))(area_list)

# colormap = plt.cm.viridis
green_red = mcolors.LinearSegmentedColormap.from_list("", ["green", "red"])

colors = [green_red(norm_num) for norm_num in normalized_numbers]

hex_colors = [mcolors.to_hex(color) for color in colors]

# print(hex_colors)
radius_list = []
for area in building_df["面積"]:
    radius_list.append(math.sqrt(area / math.pi))

building_df["color"] = hex_colors
building_df["radius"] = radius_list

building_df.to_csv(
    "/Users/oohyuti/Taipei-City-Dashboard-FE/public/csv/危樓經緯度with顏色半徑.csv"
)


def convert_to_point(coord_str):
    coord_list = json.loads(coord_str)
    return Point(coord_list)


building_df["latlng"] = building_df["latlng"].apply(convert_to_point)

# Create a GeoDataFrame
gdf = gpd.GeoDataFrame(building_df, geometry="latlng")

# Export to GeoJSON
# geojson = gdf.to_json()
gdf.to_file(output_file, driver="GeoJSON")

# The GeoJSON can now be used for plotting
# print(geojson)
# geojson = {"type": "FeatureCollection", "features": []}

# # Populate GeoJSON with data from DataFrame
# for index, row in building_df.iterrows():
#     coordinates = json.loads(row["Coordinates"])
#     feature = {
#         "type": "Feature",
#         "properties": {"color": row["color"], "radius": row["radius"]},
#         "geometry": {"type": "Point", "coordinates": coordinates},
#     }
#     geojson["features"].append(feature)

# print(json.dumps(geojson, indent=2))
# gdf.to_file(output_file, driver="GeoJSON")
