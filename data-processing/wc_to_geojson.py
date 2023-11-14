import json
import geojson

# Read JSON data from a file
json_file_path = './mydata/Wc.json'  # Replace with the actual input file path
with open(json_file_path, 'r', encoding='utf-8') as file:
    data = json.load(file)

# Convert to GeoJSON and add ID
features = []
for idx, item in enumerate(data['data'], start=1):
    coordinates = [float(item["lng"]) , float(item["lat"])]
    properties = {key: item[key] for key in item if key not in ["lng", "lat"]}
    properties["id"] = idx
    point = geojson.Point(coordinates)
    feature = geojson.Feature(geometry=point, properties=properties)
    features.append(feature)

feature_collection = geojson.FeatureCollection(features)

# Convert to GeoJSON string
geojson_string = geojson.dumps(feature_collection, indent=2)

# Write GeoJSON to a new file
geojson_output_file_path = './mydata/public_wc.geojson'  # Replace with the desired output file path
with open(geojson_output_file_path, 'w', encoding='utf-8') as output_file:
    output_file.write(geojson_string)

print(f"GeoJSON data written to: {geojson_output_file_path}")
