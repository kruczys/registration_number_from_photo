package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
    "image"
    "image/color"
    "image/jpeg"
    

    "github.com/fogleman/gg"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type PlateResult struct {
    Results []struct {
        Box struct {
            Xmin int `json:"xmin"`
            Ymin int `json:"ymin"`
            Xmax int `json:"xmax"`
            Ymax int `json:"ymax"`
        } `json:"box"`
        Plate string `json:"plate"`
    } `json:"results"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20

	router.Static("/static", "./static")
    router.Static("/uploads", "./uploads")
	router.LoadHTMLGlob("./static/*.html")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.GET("/favicon.ico", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	router.POST("/upload", handleUpload)
	router.Run(":8080")
}

func handleUpload(c *gin.Context) {

	file, err := c.FormFile("image")
	if err != nil {
		log.Println("Error retrieving the file:", err)
		c.String(http.StatusBadRequest, "Bad Request: unable to get form file")
		return
	}

	filePath := filepath.Join("uploads", file.Filename)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		log.Println("Error saving the file:", err)
		c.String(http.StatusInternalServerError, "Could not save file: %v", err)
		return
	}


	plateResult, err := recognizePlate(filePath)
	if err != nil {
		log.Println("Error saving the file:", err)
		c.String(http.StatusInternalServerError, "Error recognizing plate: %v", err)
		return
	}

    boxedImagePath := drawBoundingBox(filePath, plateResult)

    originalImageUrl := "/uploads/" + filepath.Base(filePath)
    boxedImageUrl := "/uploads/" + filepath.Base(boxedImagePath)

	c.HTML(http.StatusOK, "result.html", gin.H{
        "result": plateResult,
        "originalImageUrl": originalImageUrl,
        "boxedImageUrl": boxedImageUrl,
    })
}

func recognizePlate(filePath string) (PlateResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return PlateResult{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	mw := multipart.NewWriter(body)

	fileWriter, err := mw.CreateFormFile("upload", filepath.Base(filePath))
	if err != nil {
		return PlateResult{}, fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err = io.Copy(fileWriter, file); err != nil {
		return PlateResult{}, fmt.Errorf("failed to copy file content: %w", err)
	}

	if err = mw.WriteField("regions", "pl"); err != nil {
		return PlateResult{}, fmt.Errorf("failed to write field: %w", err)
	}

	if err = mw.Close(); err != nil {
		return PlateResult{}, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.platerecognizer.com/v1/plate-reader/", body)
	if err != nil {
		return PlateResult{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Token "+os.Getenv("API_KEY"))
	req.Header.Set("Content-Type", mw.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return PlateResult{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	var result PlateResult
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return PlateResult{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

func drawBoundingBox(imagePath string, plateResult PlateResult) string {
    imgFile, err := os.Open(imagePath)
    if err != nil {
        log.Fatal(err)
    }
    defer imgFile.Close()

    img, _, err := image.Decode(imgFile)
    if err != nil {
        log.Fatal(err)
    }

    dc := gg.NewContextForImage(img)
    dc.SetColor(color.RGBA{255, 0, 0, 255}) 
    dc.SetLineWidth(2)

    for _, result := range plateResult.Results {
        dc.DrawRectangle(float64(result.Box.Xmin), float64(result.Box.Ymin),
            float64(result.Box.Xmax-result.Box.Xmin), float64(result.Box.Ymax-result.Box.Ymin))
        dc.Stroke()
    }

    boxedImagePath := filepath.Join("uploads", "boxed_"+filepath.Base(imagePath))
    outFile, err := os.Create(boxedImagePath)
    if err != nil {
        log.Fatal(err)
    }
    defer outFile.Close()

    jpeg.Encode(outFile, dc.Image(), nil)

    return boxedImagePath
}
