import os

# list all csv under the same directory

file_paths = [
    os.path.join("./", file_name)
    for file_name in os.listdir("./") if file_name.endswith(".csv") and not file_name.endswith("_utf8.csv")
]

print(file_paths)

# Define a function to convert a file from Big5 to UTF-8 encoding
def convert_file_to_utf8(file_path):
    # Create a new file name by replacing the original extension with "_utf8.csv"
    new_file_path = os.path.join("utf8/", file_path)
    print(new_file_path)

    try:
        if not os.path.exists("./utf8"):
            os.mkdir("./utf8")
        # Read the content of the file with Big5 encoding
        with open(file_path, 'r', encoding='big5', errors='ignore') as file:
            content = file.read()

        # Write the content to a new file with UTF-8 encoding
        with open(new_file_path, 'w', encoding='utf8') as new_file:
            new_file.write(content)

        return new_file_path, "Success"
    except Exception as e:
        return file_path, f"Failed: {e}"
    
# Convert all files in the list
for file_path in file_paths:
	print(convert_file_to_utf8(file_path))