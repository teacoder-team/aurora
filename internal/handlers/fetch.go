package handlers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"orion/config"
	"orion/internal/repositories"
	"orion/pkg/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

func Fetch(c *gin.Context) {
	tag := c.Param("tag")
	id := c.Param("id")

	fileRepo := repositories.NewFileRepository()

	fileRecord, err := fileRepo.GetFileByID(id, tag)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	if fileRecord == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	cfg, _ := utils.LoadConfig()

	s3Svc, err := config.InitS3Session()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize S3 session"})
		return
	}

	var folder string
	switch tag {
	case "courses":
		folder = "courses/"
	case "avatars":
		folder = "avatars/"
	case "attachments":
		folder = "attachments/"
	default:
		folder = "misc/"
	}

	objectKey := fmt.Sprintf("%s%s", folder, fileRecord.ID)

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

	fileBytes := buf.Bytes()
	contentType := *result.ContentType

	c.Header("Content-Type", contentType)
	c.Header("Content-Length", fmt.Sprintf("%d", len(fileBytes)))
	c.Data(http.StatusOK, contentType, fileBytes)
}
