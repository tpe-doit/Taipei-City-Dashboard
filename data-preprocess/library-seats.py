from collections import defaultdict
import json
import requests

file_url = "https://seat.tpml.edu.tw/sm/service/getAllArea"
data = requests.get(file_url).json()

# Initialize a dictionary to hold the aggregated data
free_counts = defaultdict(int)
total_counts = defaultdict(int)

keys = ["總館", "稻香分館", "廣慈分館", "西湖分館"]
for key in keys:
    free_counts[key] = 0
    total_counts[key] = 0

# Aggregate the free counts by branch name
for entry in data:
    branch_name = entry["branchName"]
    free_count = entry["freeCount"]
    total_count = entry["totalCount"]
    free_counts[branch_name] += free_count
    total_counts[branch_name] += total_count

# Transform the aggregated data into the desired format
transformed_data = {
    "data": [
        {
            "name": "已使用",
            "data": [
                total_counts[branch_name] - free_counts[branch_name]
                for branch_name in free_counts.keys()
            ],
        },
        {
            "name": "未使用",
            "data": [free for free in free_counts.values()],
        },
    ]
}

print(total_counts.keys())

with open("public/chartData/2000.json", "w", encoding="utf-8") as file:
    json.dump(transformed_data, file, ensure_ascii=False)
