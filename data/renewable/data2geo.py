import json
import pandas as pd

json_file_path = 'clean-data/merged_solar_energy_data.json'

# Load the JSON data
with open(json_file_path, 'r') as file:
    json_data = json.load(file)

# Helper function to create a GeoJSON feature
def create_feature(record):
    # Extract latitude and longitude to create the geometry
    if pd.notnull(record['經度座標']) and pd.notnull(record['緯度座標']):
        geometry = {
            "type": "Point",
            "coordinates": [record['經度座標'], record['緯度座標']]
        }
    else:
        geometry = None

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

# Create GeoJSON features
features = [create_feature(record) for record in json_data]

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
