import os
import pandas as pd
import glob
# merge the csv row if the files have "台北市再生能源設置資料-太陽光電"

def extract_district(address):
    if pd.isna(address) or not isinstance(address, str):
        return None
    # The district is typically the second element when splitting by '市' or '縣', followed by '區'
    # Split the string and return the element that ends with '區'
    parts = address.split('市', 1)[-1].split('區', 1)
    district = parts[0] + '區' if parts[0] else None
    return district

# Get a list of all CSV files in the directory
file_list = glob.glob("utf8/*.csv")

# Create an empty DataFrame to store the merged data
merged_data = pd.DataFrame()

# Iterate over each file
for file in file_list:
	# Check if the file name contains the specified string
	if "台北市再生能源設置資料-太陽光電" in file:
		# Read the CSV file into a DataFrame
		data = pd.read_csv(file)
		# Append the data to the merged_data DataFrame
		merged_data = pd.concat([merged_data, data], ignore_index=True)

# if any row miss "行政區" column, fill it with extracting the 行政區 from "地址" column
# usually the "行政區" is the the words between "市" and "區" in "地址" column

for i in range(len(merged_data)):
	if pd.isnull(merged_data.loc[i, "行政區"]):
		# print(type(merged_data.loc[i, "地址"].split("市")[1]))
		# tmp = merged_data.loc[i, "地址"].split("市")[1]
		# merged_data.loc[i, "行政區"] = tmp.split("區")[0]
		merged_data.loc[i, "行政區"] = extract_district(merged_data.loc[i, "地址"])

if not os.path.exists("clean-data"):
	os.mkdir("clean-data")
# Save the merged data to a new CSV file
merged_data.to_csv("clean-data/merged_solar_energy_data.csv", index=False)

json_data = merged_data.to_json(orient="records", force_ascii=False)
with open("clean-data/merged_solar_energy_data.json", "w", encoding="utf8") as file:
	file.write(json_data)
