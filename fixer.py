import os
import xml.etree.ElementTree as ET

# Define the directory where the XML files are located
directory = "./images/"

# Loop through each XML file in the directory
for filename in os.listdir(directory):
    if filename.endswith(".xml"):
        # Construct the full file path
        filepath = os.path.join(directory, filename)

        # Parse the XML file
        tree = ET.parse(filepath)
        root = tree.getroot()

        # Find the 'path' element
        path_element = root.find("./path")

        # Replace the 'path' element text with the new path, keeping the filename the same
        new_path = f"/home/konrad/Desktop/coding/machine_learning/projekt1/registration_number_from_photo/images/{root.find('./filename').text}"
        path_element.text = new_path

        # Write the changes back to the XML file
        tree.write(filepath)
