import csv
import json

# Open the CSV file and read its contents
data = {}
data['data'] = []
with open('/Users/owenowenisme/Desktop/Taipei-City-Dashboard-FE/mydata/Wc.csv', 'r',encoding = 'utf-8-sig' )as csv_file:
	csv_reader = csv.DictReader(csv_file)
    # Convert each row of the CSV file to a dictionary
	for row in csv_reader:
		# Print the dictionary
		row.update
		if(int(row['a'])>= 1):
			row['lvl']= 4
		elif(int(row['b'])>= 1):
			row['lvl']= 3
		elif(int(row['c'])>= 1):
			row['lvl']= 2
		elif(int(row['d'])>= 1):
			row['lvl']= 1
		row.pop('a')
		row.pop('b')
		row.pop('c')
		row.pop('d')
		data['data'].append(row)
			


# Write the dictionary to a JSON file
with open('/Users/owenowenisme/Desktop/Taipei-City-Dashboard-FE/mydata/Wc.json', 'w' ,encoding='utf-8') as json_file:
	print(data)
	json.dump(data, json_file,ensure_ascii=False)
