import geopandas as gpd
from collections import defaultdict
import json
import os
import pandas as pd
from shapely.geometry import Polygon


merge_df = pd.read_csv(
    "/Users/oohyuti/Taipei-City-Dashboard-FE/public/csv/merge.csv", index_col=0
)

june_df = merge_df[merge_df["年月"] == 202306]
summed_df = june_df.groupby("TOWN")["total"].sum().reset_index()
base_path = "/Users/oohyuti/Taipei-City-Dashboard-FE/public/csv/台北市"

# print(summed_df.head())

gdf_list = []


for index, row in summed_df.iterrows():
    json_file = os.path.join(base_path, f"{row['TOWN']}.json")
    with open(json_file, "r") as file:
        data = json.load(file)
    for pair in data:
        tmp = pair[0]
        pair[0] = pair[1]
        pair[1] = tmp
        print(pair)
    polygon = Polygon(data)
    gdf = gpd.GeoDataFrame(
        {"district": row["TOWN"], "geometry": [polygon], "rain": row["total"]}
    )
    gdf_list.append(gdf)

# 合并所有 GeoDataFrame

combined_gdf = gpd.GeoDataFrame(pd.concat(gdf_list, ignore_index=True))
output_geojson_file = "/Users/oohyuti/Taipei-City-Dashboard-FE/public/mapData/district_history_rain.geojson"
combined_gdf.to_file(output_geojson_file, driver="GeoJSON")
# 查看合并后的 GeoDataFrame
# print(combined_gdf)

# gdf = gpd.GeoDataFrame({"district":row["TOWN"], "rain": }
#     [row["TOWN"], row["total"]], columns=["district", "rain"], geometry=[polygon]
# )
# print(gdf)


# final_gdf = gpd.GeoDataFrame(pd.concat(gdf_list, ignore_index=True))

# final_gdf.to_file(
#     "/Users/oohyuti/Taipei-City-Dashboard-FE/public/mapData/district_history_rain.geojson",
#     driver="GeoJSON",
# )
