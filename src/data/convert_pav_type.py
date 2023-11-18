import geopandas as gpd


def transfer_type(input_file, output_file):
    gdf = gpd.read_file(input_file)
    gdf["map_type"] = None
    for idx, line in gdf.iterrows():
        # print(line["id"].split("-"))
        prefix = line["id"].split("-")
        # print(type(line[id].split("-")))
        # prefix = str(line[id].split("-")[0])
        if prefix[0] == "J0101":
            gdf["map_type"][idx] = "PAC"
        elif prefix[0] == "J0201":
            gdf["map_type"][idx] = "人行道"
        elif prefix[0] == "J0301":
            gdf["map_type"][idx] = "公園"
        elif prefix[0] == "J0401":
            gdf["map_type"][idx] = "停車場"
        else:
            gdf["map_type"][idx] = "學校"
        # print(line)
    gdf.to_file(output_file, driver="GeoJSON")


transfer_type(
    "/Users/oohyuti/Taipei-City-Dashboard-FE/public/mapData/pav.geojson",
    "/Users/oohyuti/Taipei-City-Dashboard-FE/public/mapData/output_pav.geojson",
)
