package main

import (
	"log"
	"storage/config"
	"storage/routes"
	"storage/utils"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	cfg, err := utils.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}

	if err := config.LoadConfigTags(); err != nil {
		log.Fatalf("❌ Failed to load config tags: %v", err)
	}

	config.ConnectDatabase(cfg)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Next()
	})

	router.GET("/", routes.IndexHandler)
	router.POST("/upload", routes.UploadHandler)

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.S3Region),
		Endpoint:    aws.String(cfg.S3Endpoint),
		Credentials: credentials.NewStaticCredentials(cfg.S3AccessKeyId, cfg.S3SecretAccessKey, ""),
	})
	if err != nil {
		log.Fatalf("❌ Failed to initialize AWS session: %v", err)
	}

	s3Svc := s3.New(sess)
	_, err = s3Svc.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		log.Fatalf("❌ Failed to list S3 buckets: %v", err)
	} else {
		log.Printf("✅ Successfully connected to AWS S3")
	}

	log.Printf("🚀 Server is running at: %s\n", cfg.ApplicationURL)

	if err := router.Run(":" + strconv.Itoa(cfg.ApplicationPort)); err != nil {
		log.Fatalf("❌ Error starting server: %v\n", err)
	}
}
