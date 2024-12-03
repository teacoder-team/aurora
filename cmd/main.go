package main

import (
	"log"
	"storage/config"
	"storage/routes"
	"storage/utils"
	"strconv"

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
	
	router.Use(cors.CORSHandler(cfg))

	router.GET("/", routes.IndexHandler)
	router.POST("/upload", routes.UploadHandler)
	router.GET("/:tag/:id", routes.GetFileHandler)

	s3Svc, err := config.InitS3Session()
	if err != nil {
		log.Fatalf("❌ Failed to initialize AWS session: %v", err)
	}

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
