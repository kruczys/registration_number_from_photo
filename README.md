
# License Plate Recognizer

This project is a web application for recognizing license plates from uploaded images. The application leverages the Deep License Plate Recognition API from [parkpow](https://github.com/parkpow/deep-license-plate-recognition) to identify and highlight license plates in images.

## Features

- Upload an image of a vehicle.
- Display the original image and the image with a bounding box around the recognized license plate.
- Show the recognized license plate number.

## Motivation

Initially, I attempted to build a deep neural network model for license plate recognition from scratch. As seen in train_model.py and main.py files. However, the results were not satisfactory due to various challenges in training and optimizing the model. My main issue was using pytesseract to recognize letters from cropped out images picked by neural network. Consequently, I switched to using the well-established Deep License Plate Recognition API, which provided much better accuracy and reliability.

## Technologies Used

- Go (Golang)
- Gin Web Framework
- HTMX for AJAX requests
- Deep License Plate Recognition API
- HTML/CSS for the frontend
- Docker

## Setup

### Prerequisites

- Go 1.22.3 or later
- [Git](https://git-scm.com/)
- An API key from [parkpow](https://github.com/parkpow/deep-license-plate-recognition)

### Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/kruczys/registration_number_from_photo.git
    cd license-plate-recognizer/web_app
    ```

2. Set up environment variables:

    Create a `.env` file in the web_app directory of the project and add your API key:

    ```sh
    API_KEY=your_api_key_here
    ```

3. Install dependencies:

    ```sh
    go mod tidy
    ```

4. Run the application:

    ```sh
    go run main.go
    ```

### Using Docker (Optional)

1. Build the Docker image:

    ```sh
    git clone https://github.com/kruczys/registration_number_from_photo.git
    cd license-plate-recognizer/web_app
    docker build -t license-plate-recognizer .
    ```

2. Run the Docker container:

    ```sh
    docker run -p 8080:8080 --env-file .env license-plate-recognizer
    ```

## Usage

1. Open your web browser and navigate to `http://localhost:8080`.
2. Upload an image of a vehicle. You can use photos provided in test_photos directory located in root directory of project.
3. The application will display the original image alongside the image with a bounding box around the recognized license plate, as well as the recognized license plate number.

## Acknowledgements

- [parkpow](https://github.com/parkpow/deep-license-plate-recognition) for the Deep License Plate Recognition API which works lovely.
- [Gin Web Framework](https://github.com/gin-gonic/gin) for making web development in Go a breeze <3.
- [HTMX](https://htmx.org/) for providing an easy way to handle AJAX requests in HTML.

## Contact

If you have any questions or feedback, feel free to reach out to me at konradkreczko@gmail.com.

