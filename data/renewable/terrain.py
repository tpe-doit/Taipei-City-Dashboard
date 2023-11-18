import numpy as np
from scipy.stats import gaussian_kde
import matplotlib.pyplot as plt
from shapely.geometry import mapping, Polygon
import geopandas as gpd
import json

# Load the JSON data
file_path = 'clean-data/merged_solar_energy_data.json'
with open(file_path, 'r', encoding='utf-8') as file:
    data = json.load(file)

# Extract coordinates and capacity
coords = np.array([(item['經度座標'], item['緯度座標'], item['裝置容量-瓩'])
                  for item in data])

# Perform KDE on the coordinates weighted by capacity
x, y, z = coords[:, 0], coords[:, 1], coords[:, 2]
kde = gaussian_kde([x, y], weights=z, bw_method='scott')

# Create grid to evaluate kde upon
xmin, xmax = min(x), max(x)
ymin, ymax = min(y), max(y)
xx, yy = np.mgrid[xmin:xmax:100j, ymin:ymax:100j]
grid_coords = np.vstack([xx.ravel(), yy.ravel()])
zz = kde(grid_coords)

# Reshape the output to fit the grid shape
heatmap = zz.reshape(xx.shape)

# Since the user wants more contour lines, we need to increase the number of levels in the contour plot.
# Redefine the heatmap plot with more contour levels
plt.figure(figsize=(8, 6))
# Increased number of levels
heatmap_plot = plt.contourf(xx, yy, heatmap, cmap='hot', levels=100)
# Increased number of levels
contour_plot = plt.contour(xx, yy, heatmap, colors='black', levels=100)
plt.colorbar(heatmap_plot)

# Plot contour values
plt.clabel(contour_plot, inline=1, fontsize=8)
plt.title('Heatmap with More Contour Lines')
plt.xlabel('Longitude')
plt.ylabel('Latitude')

# Save the new contour plot to a file
contour_plot_file_more = 'clean-data/solar_energy_heatmap_contours.png'
# plt.savefig(contour_plot_file_more)

# Extract new contour paths for GeoJSON conversion
contour_paths_more = []
for contour in contour_plot.collections:
    for path in contour.get_paths():
        v = path.vertices
        coords = list(map(lambda xy: (xy[0], xy[1]), v))
        # Close the polygon if it's not closed
        if coords[0] != coords[-1]:
            coords.append(coords[0])
        contour_paths_more.append(coords)

contour_values = []
for i, contour in enumerate(contour_plot.collections):
    for path in contour.get_paths():
        contour_values.append(contour_plot.levels[i])

# Now we'll associate each contour path with its corresponding value
contour_data = zip(contour_paths_more, contour_values)


def value_to_hsl(value, min_value, max_value):
    # Scale the value to a range of [0, 1]
    scaled_value = (value - min_value) / (max_value - min_value)
    # Convert the scaled value to a hue in the HSL color space
    # Hue value will change from 120 (green, low) to 0 (red, high)
    hue = 120 - (scaled_value * 120)
    # Return the HSL color as a string
    return f"hsl({hue}, 100%, 50%)"


min_value = min(contour_values)
max_value = max(contour_values)

colors = [value_to_hsl(value, min_value, max_value)
          for value in contour_values]


# Create GeoDataFrame for new GeoJSON
gdf_more = gpd.GeoDataFrame({
    'geometry': [Polygon(coords) for coords in contour_paths_more],
    # 'value': [value for coords, value in contour_data],
    'height': [value for coords, value in contour_data],
    'color': colors
})

# Convert to new GeoJSON
geojson_file_path_more = 'clean-data/solar_energy_contours.geojson'
gdf_more.to_file(geojson_file_path_more, driver='GeoJSON')

contour_plot_file_more, geojson_file_path_more
