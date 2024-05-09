import cv2
import imutils
import pytesseract

pytesseract.pytesseract.tesseract_cmd = "tesseract"

image = cv2.imread("test.jpg")
image = imutils.resize(image, width=500)
gray_image = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)
gray_image = cv2.bilateralFilter(gray_image, 11, 17, 17)
edge = cv2.Canny(gray_image, 30, 200)
cnts, new = cv2.findContours(edge.copy(), cv2.RETR_LIST, cv2.CHAIN_APPROX_SIMPLE)
image1 = image.copy()
cv2.drawContours(image1, cnts, -1, (0, 255, 0), 3)
cv2.imshow("original image", image1)
cv2.waitKey(0)
