import csv
import json

# Open the CSV file and read its contents
data = {
  "data": [
    {
      "name": "",
      "data": [
        { "x": "北投區", "y": 0 },
        { "x": "士林區", "y": 0 },
        { "x": "內湖區", "y": 0 },
        { "x": "南港區", "y": 0 },
        { "x": "松山區", "y": 0 },
        { "x": "信義區", "y": 0 },
        { "x": "中山區", "y": 0 },
        { "x": "大同區", "y": 0 },
        { "x": "中正區", "y": 0 },
        { "x": "萬華區", "y": 0 },
        { "x": "大安區", "y": 0 },
        { "x": "文山區", "y": 0 }
      ]
    }
  ]
}
with open('/Users/owenowenisme/Desktop/Taipei-City-Dashboard-FE/mydata/Wc.csv', 'r',encoding = 'utf-8-sig' )as csv_file:
	csv_reader = csv.DictReader(csv_file)
    # Convert each row of the CSV file to a dictionary
	for row in csv_reader:
		if row['dist'] == '北投區':
			data['data'][0]['data'][0]['y']+=1
		elif row['dist'] == '士林區':
			data['data'][0]['data'][1]['y']+=1
		elif row['dist'] == '內湖區':
			data['data'][0]['data'][2]['y']+=1
		elif row['dist'] == '南港區':
			data['data'][0]['data'][3]['y']+=1
		elif row['dist'] == '松山區':
			data['data'][0]['data'][4]['y']+=1
		elif row['dist'] == '信義區':
			data['data'][0]['data'][5]['y']+=1
		elif row['dist'] == '中山區':
			data['data'][0]['data'][6]['y']+=1
		elif row['dist'] == '大同區':
			data['data'][0]['data'][7]['y']+=1
		elif row['dist'] == '中正區':
			data['data'][0]['data'][8]['y']+=1
		elif row['dist'] == '萬華區':
			data['data'][0]['data'][9]['y']+=1
		elif row['dist'] == '大安區':
			data['data'][0]['data'][10]['y']+=1
		elif row['dist'] == '文山區':
			data['data'][0]['data'][11]['y']+=1
		else:
			print(row)
			


# Write the dictionary to a JSON file
with open('/Users/owenowenisme/Desktop/Taipei-City-Dashboard-FE/mydata/Wc_num', 'w' ,encoding='utf-8') as json_file:
	print(data)
	json.dump(data, json_file,ensure_ascii=False)
