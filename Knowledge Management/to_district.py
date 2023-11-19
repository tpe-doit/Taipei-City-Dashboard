# Provided data
data = [
    {"name": "停車場", "data": [3346, 0, 2520, 235, 0, 0, 777, 0, 735, 0, 0, 652]},
    {"name": "PAC", "data": [0, 0, 0, 0, 0, 0, 17691, 0, 4438, 5438, 27201, 0]},
    {"name": "人行道", "data": [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]},
    {
        "name": "學校",
        "data": [20360, 409, 2608, 0, 400, 618, 1527, 3283, 0, 2799, 4375, 0],
    },
    {
        "name": "公園",
        "data": [17047, 12707, 5721, 4383, 3577, 768, 0, 1689, 5643, 1552, 894, 2200],
    },
]

# Taipei districts
districts = [
    "北投區",
    "士林區",
    "內湖區",
    "南港區",
    "松山區",
    "信義區",
    "中山區",
    "大同區",
    "中正區",
    "萬華區",
    "大安區",
    "文山區",
]

# Iterate and sum all numbers by index
sums_by_district = [0] * 12  # 12 districts, initialize all sums to 0
for item in data:
    for i, value in enumerate(item["data"]):
        sums_by_district[i] += value

# Map each index to the district of Taipei
district_data = [
    {"x": district, "y": sum_val}
    for district, sum_val in zip(districts, sums_by_district)
]

print(district_data)
