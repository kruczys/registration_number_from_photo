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

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

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

    imageUrl := "/uploads/" + file.Filename

	c.HTML(http.StatusOK, "result.html", gin.H{
        "result": plateResult,
        "imageUrl": imageUrl,
    })
}

func recognizePlate(filePath string) (map[string]interface{}, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	mw := multipart.NewWriter(body)

	fileWriter, err := mw.CreateFormFile("upload", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err = io.Copy(fileWriter, file); err != nil {
		return nil, fmt.Errorf("failed to copy file content: %w", err)
	}

	if err = mw.WriteField("regions", "pl"); err != nil {
		return nil, fmt.Errorf("failed to write field: %w", err)
	}

	if err = mw.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.platerecognizer.com/v1/plate-reader/", body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Token "+os.Getenv("API_KEY"))
	req.Header.Set("Content-Type", mw.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}
