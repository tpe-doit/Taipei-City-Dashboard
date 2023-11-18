import json


def merge_geojson(file1, file2, output_file):
    # Read the contents of the first GeoJSON file
    with open(file1, 'r', encoding='utf8') as f1:
        data1 = json.load(f1)

    # Read the contents of the second GeoJSON file
    with open(file2, 'r', encoding='utf8') as f2:
        data2 = json.load(f2)

    # Combine the features from both files into a single feature collection
    combined_features = data1['features'] + data2['features']
    combined_data = {
        'type': 'FeatureCollection',
        'features': combined_features
    }

    # Write the combined feature collection to a new GeoJSON file
    with open(output_file, 'w', encoding='utf8') as output:
        json.dump(combined_data, output, ensure_ascii=False)


# Usage example
file1 = 'clean-data/merged_solar_energy_data.geojson'
file2 = 'clean-data/solar_energy_contours.geojson'
output_file = 'clean-data/solar_energy.geojson'

merge_geojson(file1, file2, output_file)
