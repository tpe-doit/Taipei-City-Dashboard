from get_district import is_point_in_district, is_point_inside_polygon
import json

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
    "color": ["#9c7a3e", "#b7e28e", "#7febd0", "#8ce8ff"],
    "types": ["BarPercentChart", "ColumnChart"],
    "unit": "UNIT",
    "categories": districts,
}

example_data = {
    "data": [
        {
            "name": ">40年",
            "data": [
                15742,
                20453,
                6060,
                6410,
                10855,
                12653,
                15154,
                12984,
                11563,
                17647,
                20243,
                12710,
            ],
        },
        {
            "name": "20-40年",
            "data": [
                6312,
                7805,
                12366,
                2857,
                4194,
                6702,
                6043,
                2500,
                3174,
                3731,
                7165,
                8214,
            ],
        },
        {
            "name": "5-20年",
            "data": [
                1963,
                2335,
                3357,
                1351,
                684,
                1022,
                1947,
                1025,
                1171,
                716,
                1469,
                2039,
            ],
        },
        {"name": "<5年", "data": [48, 91, 71, 12, 6, 5, 49, 25, 33, 14, 38, 73]},
    ]
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
    ret_data = {}
    with open(file_path) as f:
        geojson = json.load(f)
        for feature in geojson["features"]:
            # name is the category of the BarPercentChart
            name = feature["properties"]["Name"]
            ret_data["name"] = name
            if feature["geometry"]["type"] == "MultiPolygon":
                ret_data["data"] = multiPolygon_data(feature["geometry"]["coordinates"])
            else:
                print("Not Implemented")
    return ret_data


geojson_path = "./public/mapData/tp_flood.geojson"
ret = gen_district_data_from_geojson(geojson_path)
with open("./public/chartData/new_file.json", "w+") as f:
    json.dump(ret, f)
