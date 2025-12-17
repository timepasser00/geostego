package main

import (
	"bytes"
	"fmt"
	"image/png"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/timepasser00/geostego/pkg/stego"
)

func main() {

	r := gin.Default()

	r.MaxMultipartMemory = 8 << 20 // 8 MiB

	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "GeoStego (Gin) is operational")
	})

	r.POST("/encode", handleEncode)
	r.POST("/decode", handleDecode)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}

func handleEncode(c *gin.Context) {
	message := c.PostForm("message")
	if message == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message field is required"})
		return
	}

	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image file is required"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not open uploaded file"})
		return
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image format. Only PNG is supported."})
		return
	}

	stegoImg, err := stego.Encode(img, []byte(message))
	if err != nil {
		if err == stego.ErrMessageTooLarge {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "Message is too long for this image size"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Encoding failed"})
		return
	}

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, stegoImg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate output image"})
		return
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Expose-Headers", "Content-Disposition");
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileHeader.Filename))
	c.Data(http.StatusOK, "image/png", buf.Bytes())
}

func handleDecode(c *gin.Context) {

	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image file is required"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not open uploaded file"})
		return
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image format. Only PNG is supported."})
		return
	}

	extractedBytes, err := stego.Decode(img)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No hidden message found or decryption failed"})
		return
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": string(extractedBytes),
	})
}
