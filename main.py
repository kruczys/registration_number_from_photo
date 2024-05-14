import os
import cv2
import numpy as np
import pandas as pd
import tensorflow as tf
import pytesseract as pt
import plotly.express as px
import matplotlib.pyplot as plt
import xml.etree.ElementTree as xet

from glob import glob
from skimage import io
from shutil import copy
from tensorflow.keras.models import Model
from tensorflow.keras.callbacks import TensorBoard
from sklearn.model_selection import train_test_split
from tensorflow.keras.applications import InceptionResNetV2
from tensorflow.keras.layers import Dense, Dropout, Flatten, Input
from tensorflow.keras.preprocessing.image import load_img, img_to_array

path = glob("./images/*.xml")
labels_dict = dict(filepath=[], xmin=[], xmax=[], ymin=[], ymax=[])
for filename in path:

    info = xet.parse(filename)
    root = info.getroot()
    member_object = root.find("object")
    labels_info = member_object.find("bndbox")
    xmin = int(labels_info.find("xmin").text)
    xmax = int(labels_info.find("xmax").text)
    ymin = int(labels_info.find("ymin").text)
    ymax = int(labels_info.find("ymax").text)

    labels_dict["filepath"].append(filename)
    labels_dict["xmin"].append(xmin)
    labels_dict["xmax"].append(xmax)
    labels_dict["ymin"].append(ymin)
    labels_dict["ymax"].append(ymax)

df = pd.DataFrame(labels_dict)
df.to_csv("labels.csv", index=False)

filename = df["filepath"][0]


def getFilename(filename):
    filename_image = xet.parse(filename).getroot().find("filename").text
    filepath_image = os.path.join("./images", filename_image)
    return filepath_image


image_path = list(df["filepath"].apply(getFilename))

labels = df.iloc[:, 1:].values
data = []
output = []
for ind in range(len(image_path)):
    image = image_path[ind]
    img_arr = cv2.imread(image)
    h, w, d = img_arr.shape
    load_image = load_img(image, target_size=(224, 224))
    load_image_arr = img_to_array(load_image)
    norm_load_image_arr = load_image_arr / 255.0
    xmin, xmax, ymin, ymax = labels[ind]
    nxmin, nxmax = xmin / w, xmax / w
    nymin, nymax = ymin / h, ymax / h
    label_norm = (nxmin, nxmax, nymin, nymax)
    data.append(norm_load_image_arr)
    output.append(label_norm)

X = np.array(data, dtype=np.float32)
y = np.array(output, dtype=np.float32)
x_train, x_test, y_train, y_test = train_test_split(
    X, y, train_size=0.8, random_state=0
)

inception_resnet = InceptionResNetV2(
    weights="imagenet", include_top=False, input_tensor=Input(shape=(224, 224, 3))
)

headmodel = inception_resnet.output
headmodel = Flatten()(headmodel)
headmodel = Dense(500, activation="relu")(headmodel)
headmodel = Dense(250, activation="relu")(headmodel)
headmodel = Dense(4, activation="sigmoid")(headmodel)


# ---------- model
model = Model(inputs=inception_resnet.input, outputs=headmodel)
model.compile(loss="mse", optimizer=tf.keras.optimizers.Adam(learning_rate=1e-4))
tfb = TensorBoard("object_detection")
history = model.fit(
    x=x_train,
    y=y_train,
    batch_size=10,
    epochs=180,
    validation_data=(x_test, y_test),
    callbacks=[tfb],
)
model.save("./object_detection.h5")
