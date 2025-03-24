package handlers

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"orion/config"
	"orion/internal/models"
	"orion/internal/repositories"
	"orion/pkg/logger"
	"orion/pkg/utils"
	"strings"

	_ "image/jpeg"
	_ "image/png"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

func Upload(c *gin.Context) {
	cfg, err := utils.LoadConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load config"})
		return
	}

	secretKey := c.GetHeader("X-Upload-Secret")
	if secretKey != cfg.UploadSecretKey {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid secret key"})
		return
	}

	file, err := getFileFromRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	buf, err := readFile(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tag, err := validateTag(c, cfg)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileID, err := utils.GenerateID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate file ID"})
		return
	}

	s3Svc, err := config.InitS3Session()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize S3 session"})
		return
	}

	err = uploadToS3(s3Svc, cfg, tag, file, fileID, buf)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fileRepo := repositories.NewFileRepository()

	fileRecord := models.File{
		ID:          fileID,
		Tag:         tag,
		Filename:    file.Filename,
		ContentType: file.Header.Get("Content-Type"),
		Size:        int(file.Size),
		Deleted:     nil,
	}

	if err := fileRepo.CreateFile(&fileRecord); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file metadata in database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "File uploaded successfully",
		"file_id":  fileID,
		"filename": file.Filename,
	})
}

func getFileFromRequest(c *gin.Context) (*multipart.FileHeader, error) {
	file, err := c.FormFile("file")
	if err != nil {
		logger.Error("Failed to get file from request", err)
		return nil, fmt.Errorf("failed to get file from request: %w", err)
	}
	return file, nil
}

func readFile(file *multipart.FileHeader) (*bytes.Buffer, error) {
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(src)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return buf, nil
}

func validateTag(c *gin.Context, cfg *utils.Config) (string, error) {
	tag := c.PostForm("tag")
	allowedTags := map[string]struct{}{}
	for _, allowedTag := range strings.Split(cfg.AllowedTags, ",") {
		allowedTags[strings.TrimSpace(allowedTag)] = struct{}{}
	}

	if _, exists := allowedTags[tag]; !exists {
		return "Invalid tag", nil
	}

	return tag, nil
}

func uploadToS3(s3Svc *s3.S3, cfg *utils.Config, tag string, file *multipart.FileHeader, fileID string, buf *bytes.Buffer) error {
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

	_, err := s3Svc.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(cfg.S3BucketName),
		Key:         aws.String(folder + fileID),
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String(file.Header.Get("Content-Type")),
	})
	return err
}
