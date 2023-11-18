all: convert-encoding merge-solar-energy convert2geo heatmap merge-jeojson copy-2-mapdata

convert-encoding:
	python big5_to_utf8.py

merge-solar-energy:
	python merge_solar_energy_data.py

convert2geo:
	python data2geo.py

heatmap:
	python terrain.py

merge-jeojson:
	python main.py

copy-2-mapdata:
	cp clean-data/solar_energy.geojson ../../public/mapData/solar_energy.geojson