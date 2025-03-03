package routes

import (
	"bytes"
	"fmt"
	"image"
	"log"
	"net/http"
	"storage/config"
	"storage/models"
	"storage/utils"

	_ "image/jpeg"
	_ "image/png"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

func findFileByID(id, tag string) (*models.File, error) {
	var file models.File
	result := config.DB.Where("id = ? AND tag = ? AND deleted IS NULL", id, tag).First(&file)
	if result.Error != nil {
		return nil, result.Error
	}
	return &file, nil
}

func Fetch(c *gin.Context) {
	tag := c.Param("tag")
	id := c.Param("id")

	fileRecord, err := findFileByID(id, tag)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found in database"})
		return
	}

	cfg, _ := utils.LoadConfig()
	s3Svc, err := config.InitS3Session()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize S3 session"})
		return
	}

	var s3Folder string
	switch tag {
	case "courses":
		s3Folder = "courses/"
	case "avatars":
		s3Folder = "avatars/"
	case "attachments":
		s3Folder = "attachments/"
	default:
		s3Folder = "misc/"
	}

	objectKey := fmt.Sprintf("%s%s", s3Folder, id)
	var fileBytes []byte

	result, err := s3Svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(cfg.S3BucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		log.Printf("S3 fetch error: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found in S3"})
		return
	}
	defer result.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(result.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	fileBytes = buf.Bytes()

	var width, height int
	if fileRecord.ContentType == "image/jpeg" || fileRecord.ContentType == "image/png" {
		img, _, err := image.Decode(bytes.NewReader(fileBytes))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode image"})
			return
		}
		width = img.Bounds().Dx()
		height = img.Bounds().Dy()

		log.Printf("Image dimensions: %dx%d", width, height)
	}

	c.Header("Content-Type", fileRecord.ContentType)
	c.Header("Content-Length", fmt.Sprintf("%d", len(fileBytes)))
	c.Data(http.StatusOK, fileRecord.ContentType, fileBytes)
}
