package routes

import (
	"bytes"
	"image"
	"log"
	"net/http"
	"storage/config"
	"storage/models"
	"storage/utils"
	"strings"
	"time"

	_ "image/jpeg"
	_ "image/png"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

func Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file from request"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer src.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	cfg, _ := utils.LoadConfig()

	tag := c.PostForm("tag")

	allowedTags := map[string]struct{}{}
	for _, allowedTag := range strings.Split(cfg.AllowedTags, ",") {
		allowedTags[strings.TrimSpace(allowedTag)] = struct{}{}
	}

	if _, exists := allowedTags[tag]; !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag. Allowed tags are: " + cfg.AllowedTags})
		return
	}

	fileID, err := utils.GenerateID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate file ID"})
		return
	}

	metadataID, err := utils.GenerateID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate metadata ID"})
		return
	}

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

	_, err = s3Svc.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(cfg.S3BucketName),
		Key:         aws.String(s3Folder + fileID),
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String(file.Header.Get("Content-Type")),
	})
	if err != nil {
		log.Printf("S3 upload error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file to S3"})
		return
	}

	fileType := getFileType(file.Filename, file.Header.Get("Content-Type"))

	var width, height int
	if fileType == "Image" {
		img, _, err := image.Decode(bytes.NewReader(buf.Bytes()))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode image"})
			return
		}
		width = img.Bounds().Dx()
		height = img.Bounds().Dy()
	}

	metadata := models.Metadata{
		ID:     metadataID,
		Type:   fileType,
		Width:  width,
		Height: height,
	}

	if err := config.DB.Create(&metadata).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create metadata"})
		return
	}

	fileRecord := models.File{
		ID:          fileID,
		Tag:         tag,
		Filename:    file.Filename,
		MetadataID:  metadata.ID,
		ContentType: file.Header.Get("Content-Type"),
		Size:        int(file.Size),
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}

	if err := config.DB.Create(&fileRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file metadata in database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "File uploaded successfully",
		"file_id":  fileID,
		"filename": file.Filename,
	})
}

func getFileType(filename, contentType string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch {
	case ext == ".jpg" || ext == ".jpeg" || ext == ".webp" || ext == ".png" || ext == ".gif" || ext == ".bmp":
		return "Image"
	case ext == ".mp4" || ext == ".avi" || ext == ".mkv" || ext == ".mov":
		return "Video"
	case ext == ".mp3" || ext == ".wav" || ext == ".ogg" || ext == ".flac":
		return "Audio"
	case ext == ".zip" || ext == ".rar" || ext == ".tar" || ext == ".gz":
		return "Zip"
	default:
		return "Misc"
	}
}
