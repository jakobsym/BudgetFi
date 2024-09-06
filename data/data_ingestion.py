import pandas as pd

# Specify the path to your Excel file
file_path = 'filetypes/excel/CPIForecast.xlsx'  # Replace with your actual file path

# Read the Excel file
# For .xlsx files, the engine 'openpyxl' is used by default
data = pd.read_excel(file_path, sheet_name=None)  # sheet_name=None reads all sheets

# Print the names of the sheets
print("Sheets available in the Excel file:", data.keys())

# To read a specific sheet, use:
# df = pd.read_excel(file_path, sheet_name='Sheet1')  # Replace 'Sheet1' with your desired sheet name

# Print the first few rows of the data from the first sheet
first_sheet = list(data.keys())[0]
df = data[first_sheet]
print(df.head())

# Filter rows where a specific column has a certain value
filtered_data = df[df['ColumnName'] == 'SomeValue']

# Display the filtered data
print(filtered_data)