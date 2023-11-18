import json
import pandas as pd

json_file_path = 'clean-data/merged_solar_energy_data.json'

# Load the JSON data
with open(json_file_path, 'r') as file:
    json_data = json.load(file)

# Helper function to create a GeoJSON feature


def create_feature(record, color):
    # Extract latitude and longitude to create the geometry
    if pd.notnull(record['經度座標']) and pd.notnull(record['緯度座標']):
        # use the coordinates to create a rectangle
        HALF_SIDE = 0.001
        rectangle = [
                    [record['經度座標'] - HALF_SIDE, record['緯度座標'] - HALF_SIDE],
                    [record['經度座標'] + HALF_SIDE, record['緯度座標'] - HALF_SIDE],
                    [record['經度座標'] + HALF_SIDE, record['緯度座標'] + HALF_SIDE],
                    [record['經度座標'] - HALF_SIDE, record['緯度座標'] + HALF_SIDE],
                    [record['經度座標'] - HALF_SIDE, record['緯度座標'] - HALF_SIDE]
        ]
        geometry = {
            "type": "Polygon",
            "coordinates": [rectangle]
        }
    else:
        geometry = None

    if pd.notnull(record['裝置容量-瓩']):
        record['density'] = record['裝置容量-瓩']
        record['height'] = record['裝置容量-瓩']

    record['color'] = color
    record['type'] = "柱狀圖"

    # Use the rest of the record as properties, and remove the coordinates
    properties = record.copy()
    properties.pop('經度座標', None)
    properties.pop('緯度座標', None)

    feature = {
        "type": "Feature",
        "geometry": geometry,
        "properties": properties
    }

    return feature


def value_to_hsl(value, min_value, max_value, base_color_hue=120):
    # Scale the value to a range of [0, 1]
    scaled_value = (value - min_value) / (max_value - min_value)
    # Calculate the hue shift based on the base color hue
    hue = base_color_hue - (scaled_value * (base_color_hue - 0))
    # The base color #80e3d4 corresponds to an HSL value of (180, 100%, 66.7%)
    # We will keep the saturation and lightness of the base color
    saturation = 58.8  # Adjusted for the base color
    lightness = 67.8  # Adjusted for the base color
    # Return the HSL color as a string
    return f"hsl({hue}, {saturation}%, {lightness}%)"


min_value = min([record['裝置容量-瓩'] for record in json_data])
max_value = max([record['裝置容量-瓩'] for record in json_data])

colors = [value_to_hsl(value, min_value, max_value)
          for value in [record['裝置容量-瓩'] for record in json_data]]

datas = zip(json_data, colors)

# Create GeoJSON features
features = [create_feature(record, color) for record, color in datas]

# Define the GeoJSON structure
geojson_data = {
    "type": "FeatureCollection",
    "features": features
}

# Save the GeoJSON data to a file
geojson_file_path = 'clean-data/merged_solar_energy_data.geojson'
with open(geojson_file_path, 'w') as file:
    json.dump(geojson_data, file, ensure_ascii=False)

geojson_file_path
