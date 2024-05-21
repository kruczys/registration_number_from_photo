import os
import cv2
import numpy as np
import pandas as pd
import tensorflow as tf
import pytesseract as pt
import plotly.express as px
import matplotlib.pyplot as plt
import tensorflow as tf
import pytesseract as pt
from tensorflow.keras.preprocessing.image import load_img, img_to_array
from tensorflow import keras


model = tf.keras.models.load_model("./object_detection.keras")
print("Model loaded")

def object_detection(path, model):
    image = load_img(path)
    image = np.array(image, dtype=np.uint8)
    resized_image = load_img(path, target_size=(224, 224))

    image_arr_224 = img_to_array(resized_image) / 255.0
    test_arr = image_arr_224.reshape(1, 224, 224, 3)

    coords = model.predict(test_arr)

    h, w, _ = image.shape
    denorm = np.array([w, w, h, h])
    coords = coords * denorm
    coords = coords.astype(np.int32)

    xmin, xmax, ymin, ymax = coords[0]
    pt1 = (xmin, ymin)
    pt2 = (xmax, ymax)

    cv2.rectangle(image, pt1, pt2, (0, 255, 0), 3)
    return image, coords
