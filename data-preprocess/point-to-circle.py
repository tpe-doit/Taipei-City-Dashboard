import geopy.distance
import json

file_path = "public/mapData/dev_bd.geojson"
new_file_path = "public/mapData/dev_bd_circle.geojson"

with open(file_path, "r") as file:
    geojson_data = json.load(file)


def create_circle_coordinates_lon_lat(center, radius, num_points=36):
    """
    Create a circle approximation using polygons around a center point (longitude, latitude).
    :param center: tuple of (longitude, latitude) for the center of the circle.
    :param radius: radius of the circle in meters.
    :param num_points: number of points to use in the approximation, more points = smoother circle.
    :return: list of coordinates forming the circle.
    """
    lon, lat = center
    circle_coords = []
    for i in range(num_points):
        angle = 360 / num_points * i
        point = geopy.distance.distance(meters=radius).destination(
            (lat, lon), bearing=angle
        )
        circle_coords.append((point.longitude, point.latitude))
    circle_coords.append(
        circle_coords[0]
    )  # Complete the circle by adding the first point at the end
    return circle_coords


# Create circle geometries for each point, considering the corrected coordinate order (longitude, latitude)
for feature in geojson_data["features"]:
    if feature["geometry"]["type"] == "Point":
        center_point = tuple(feature["geometry"]["coordinates"])
        circle_coords = create_circle_coordinates_lon_lat(
            center_point, 200
        )  # 200 meters radius
        feature["geometry"] = {"type": "Polygon", "coordinates": [circle_coords]}


with open(new_file_path, "w") as new_file:
    json.dump(geojson_data, new_file)

new_file_path
