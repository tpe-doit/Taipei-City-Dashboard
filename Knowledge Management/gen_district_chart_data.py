from get_district import is_point_in_district, is_point_inside_polygon
import json
from typing import List

districts = [
    "北投區",
    "士林區",
    "內湖區",
    "南港區",
    "松山區",
    "信義區",
    "中山區",
    "大同區",
    "中正區",
    "萬華區",
    "大安區",
    "文山區",
]

example_chart_config: {
    "color": ["#C3E0E5", "#41729F"],
    "types": ["BarPercentChart", "ColumnChart"],
    "unit": "UNIT",
    "categories": [
        "北投區",
        "士林區",
        "內湖區",
        "南港區",
        "松山區",
        "信義區",
        "中山區",
        "大同區",
        "中正區",
        "萬華區",
        "大安區",
        "文山區",
    ],
}


def calculate_center(points):
    sum_x = sum(point[0] for point in points)
    sum_y = sum(point[1] for point in points)
    num_points = len(points)
    center_x = sum_x / num_points
    center_y = sum_y / num_points
    return center_x, center_y


def multiPolygon_data(coordinates):
    data = [0] * len(districts)
    for i, district in enumerate(districts):
        for cord in coordinates:
            area = cord[0]
            c_y, c_x = calculate_center(area)
            if is_point_inside_polygon(c_x, c_y, district):
                data[i] += 1
    return data


def gen_district_data_from_geojson(file_path):
    ret_data = []
    record = {}
    with open(file_path) as f:
        geojson = json.load(f)
        for index, feature in enumerate(geojson["features"]):
            # name is the category of the BarPercentChart
            name = feature["properties"]["map_type"]
            if name not in record:
                record[name] = {}
                record[name]["name"] = name
                record[name]["data"] = [0] * len(districts)

            for i, district in enumerate(districts):
                try:
                    for cord in feature["geometry"]["coordinates"]:
                        area = cord[0]
                        if not isinstance(area[0], List):
                            area = cord
                        c_y, c_x = calculate_center(area)
                        if is_point_inside_polygon(c_x, c_y, district):
                            record[name]["data"][i] += int(feature["properties"]["面積"])
                except:
                    continue
    for k, v in record.items():
        ret_data.append(v)
    return ret_data


geojson_path = "./public/mapData/all_pav.geojson"

ret = {}
ret["data"] = gen_district_data_from_geojson(geojson_path)

with open("./public/chartData/new_file.json", "w+") as f:
    json.dump(ret, f, ensure_ascii=False)
