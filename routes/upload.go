package routes

import (
	"bytes"
	"image"
	"log"
	"net/http"
	"storage/config"
	"storage/models"
	"storage/utils"
	"time"

	_ "image/jpeg" // Для поддержки JPEG
	_ "image/png"  // Для поддержки PNG

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

func UploadHandler(c *gin.Context) {
	// Получаем файл из запроса
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file from request"})
		return
	}

	// Открываем файл
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer src.Close()

	// Читаем содержимое файла
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	// Генерация уникального ID для файла
	fileID, err := utils.GenerateID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate file ID"})
		return
	}

	// Генерация уникального ID для metadata
	metadataID, err := utils.GenerateID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate metadata ID"})
		return
	}

	// Инициализация S3-сессии
	cfg, _ := utils.LoadConfig()
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.S3Region),
		Endpoint:    aws.String(cfg.S3Endpoint),
		Credentials: credentials.NewStaticCredentials(cfg.S3AccessKeyId, cfg.S3SecretAccessKey, ""),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize S3 session"})
		return
	}

	s3Svc := s3.New(sess)

	// Определяем папку в S3 в зависимости от тега
	tag := c.PostForm("tag")
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

	// Загружаем файл в S3
	_, err = s3Svc.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(cfg.S3BucketName),
		Key:         aws.String(s3Folder + fileID), // Используем сгенерированный fileID и путь для тега
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String(file.Header.Get("Content-Type")),
	})
	if err != nil {
		log.Printf("S3 upload error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file to S3"})
		return
	}

	// Получаем размеры изображения, если файл — это изображение
	var width, height int
	if file.Header.Get("Content-Type") == "image/jpeg" || file.Header.Get("Content-Type") == "image/png" {
		img, _, err := image.Decode(bytes.NewReader(buf.Bytes()))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode image"})
			return
		}
		width = img.Bounds().Dx()
		height = img.Bounds().Dy()
	}

	// Создаем запись в таблице metadata
	metadata := models.Metadata{
		ID:     metadataID,
		Type:   file.Header.Get("Content-Type"),
		Width:  width,
		Height: height,
	}

	if err := config.DB.Create(&metadata).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create metadata"})
		return
	}

	// Создаем запись о файле в таблице files
	fileRecord := models.File{
		ID:          fileID,
		Tag:         tag,
		Filename:    file.Filename,
		MetadataID:  metadata.ID, // Используем только что созданный metadataID
		ContentType: file.Header.Get("Content-Type"),
		Size:        int(file.Size),
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}

	if err := config.DB.Create(&fileRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file metadata in database"})
		return
	}

	// Возвращаем успешный ответ
	c.JSON(http.StatusOK, gin.H{
		"message":  "File uploaded successfully",
		"file_id":  fileID,
		"filename": file.Filename,
	})
}
