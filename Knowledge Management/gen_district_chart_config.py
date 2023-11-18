from get_district import is_point_in_district
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


geojson_path = "./public/mapData/tp_flood.geojson"


def multiPolygon_config(coordinates):
    data = { "name": "", "data": [0] * len(coordinates)}
    for cord in coordinates:
        area = cord[0]
        for boundary in area:
            x = boundary[0]
            y = boundary[1]
            for district in districts:
                if is_point_in_district(x, y, district):
                    
                    


def gen_district_data_from_geojson(file_path):
    with open(file_path) as f:
        geojson = json.load(f)
        for feature in geojson["features"]:
            # name is the category of the BarPercentChart
            name = feature["properties"]["Name"]
            if feature["geometry"]["type"] == "MultiPolygon":
                return multiPolygon_config(feature["geometry"]["coordinates"], name)

            else:
                print("Not Implemented")
                # try:
                #     area = sum([float(num) for num in gdf["面積"]])
                # except:
                #     area = sum(gdf.geometry.area)


gen_district_data_from_geojson(geojson_path)
