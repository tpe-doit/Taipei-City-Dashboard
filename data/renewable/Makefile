all: convert-encoding merge-solar-energy convert2geo heatmap copy-2-mapdata

convert-encoding:
	python big5_to_utf8.py

merge-solar-energy:
	python merge_solar_energy_data.py

convert2geo:
	python data2geo.py

heatmap:
	python terrain.py

copy-2-mapdata:
	cp clean-data/merged_solar_energy_data.geojson ../../public/mapData/solar_energy.geojson
	cp clean-data/solar_energy_contours.geojson ../../public/mapData/solar_energy_contour.geojson