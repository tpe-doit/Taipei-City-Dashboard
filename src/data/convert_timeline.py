import geopandas as gpd
from collections import defaultdict
import json
import pandas as pd

station_df = pd.read_csv(
    "/Users/oohyuti/Taipei-City-Dashboard-FE/public/csv/臺北市水利處雨量站(csv).csv"
)
history_df = pd.read_csv(
    "/Users/oohyuti/Taipei-City-Dashboard-FE/public/csv/歷史日雨量(UTF8).csv"
)

merged_df = pd.merge(station_df, history_df, left_on="NAME", right_on="站名", how="inner")

day_list = []
for i in range(31):
    day_list.append(str(i + 1) + "日雨量")
merged_df["total"] = merged_df[day_list].sum(axis=1)
new_df = merged_df[["TOWN", "STATIONID", "NAME", "X坐標", "Y坐標", "年月", "total"]]
summed_df = new_df.groupby("年月")["total"].sum().reset_index()

# summed_df.to_csv("/Users/oohyuti/Taipei-City-Dashboard-FE/public/csv/sum.csv")


def convert_to_datetime(year_month):
    return pd.to_datetime(str(year_month), format="%Y%m").strftime(
        "%Y-%m-%dT%H:%M:%S+08:00"
    )


summed_df["x"] = summed_df["年月"].apply(convert_to_datetime)
summed_df["y"] = summed_df["total"]
df = summed_df[["y", "x"]]
output_file = "/Users/oohyuti/Taipei-City-Dashboard-FE/public/csv/622.json"
# Print or save the JSON

output_dict = {}
output_dict["data"] = [{"name": "月降雨量"}, {"data": df.to_dict(orient="records")}]
with open(output_file, "w", encoding="utf-8") as file:
    json.dump(output_dict, file, ensure_ascii=False, indent=4)

# print(json_str)
# file_path = (
#     "/Users/oohyuti/Taipei-City-Dashboard-FE/public/mapData/history_flood.geojson"
# )

# gdf = gpd.read_file(file_path)

# monthly_counts = defaultdict(int)
# town_counts = defaultdict(int)
# print(gdf)
# for _, row in gdf.iterrows():
#     date_str = row["properties"]["FDATE"]
#     year_month = date_str[:7]  # Extract year and month
#     monthly_counts[year_month] += 1

#     town_name = row["properties"]["TOWN_NAME"]
#     town_counts[town_name] += 1

# formatted_monthly_data = {
#     "data": [
#         {"y": count, "x": f"{year_month}-01T00:00:00+08:00"}
#         for year_month, count in monthly_counts.items()
#     ]
# }

# # Convert town counts to a simple dictionary (if needed)
# formatted_town_data = dict(town_counts)

# # Output the results
# print("Monthly Counts:")
# print(json.dumps(formatted_monthly_data, indent=4))
# print("\nTown Counts:")
# print(formatted_town_data)
