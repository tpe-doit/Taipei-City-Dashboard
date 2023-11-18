from collections import defaultdict
import json
import requests

file_url = "https://seat.tpml.edu.tw/sm/service/getAllArea"
data = requests.get(file_url).json()

# Initialize a dictionary to hold the aggregated data
aggregated_data = defaultdict(int)

# Aggregate the free counts by branch name
for entry in data:
    branch_name = entry["branchName"]
    free_count = entry["freeCount"]
    aggregated_data[branch_name] += free_count

# Transform the aggregated data into the desired format
transformed_data = {
    "data": [
        {
            "name": "",
            "data": [
                {"x": branch, "y": free} for branch, free in aggregated_data.items()
            ],
        }
    ]
}

with open("public/chartData/2000.json", "w", encoding="utf-8") as file:
    json.dump(transformed_data, file, ensure_ascii=False)
