import geopandas as gpd
from shapely.geometry import Point, Polygon
import json
import os
import warnings

warnings.simplefilter(action="ignore", category=FutureWarning)


taipei_districts = {}
dir = "./Knowledge Management/towns/台北市"
for filename in os.listdir(dir):
    if filename.endswith(".json"):
        file_path = os.path.join(dir, filename)
        with open(file_path) as f:
            coordinates = json.load(f)
            taipei_districts[filename[:-5]] = coordinates

taipei_districts["test"] = [(0, 0), (2, 0), (2, 2), (0, 2)]


# Function to check if points are within the district
def is_point_in_district(x, y, district):
    district_polygon = Polygon(taipei_districts[district])
    district_gdf = gpd.GeoDataFrame([1], geometry=[district_polygon])
    print(district_gdf)
    point = Point(x, y)
    return district_gdf.contains(point).bool()


def is_point_inside_polygon(x, y, district):
    num = len(taipei_districts[district])
    inside = False

    p1x, p1y = taipei_districts[district][0]
    for i in range(num + 1):
        p2x, p2y = taipei_districts[district][i % num]
        if y > min(p1y, p2y):
            if y <= max(p1y, p2y):
                if x <= max(p1x, p2x):
                    if p1y != p2y:
                        xints = (y - p1y) * (p2x - p1x) / (p2y - p1y) + p1x
                    if p1x == p2x or x <= xints:
                        inside = not inside
        p1x, p1y = p2x, p2y

    return inside


# # Example usage
# print(is_point_in_district(25.1942, 121.578, "士林區"))
