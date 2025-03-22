package main

import (
	"orion/config"
	"orion/internal/handlers"
	"orion/pkg/logger"
	"orion/pkg/utils"
	"strconv"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	logger.InitLogger()

	cfg, err := utils.LoadConfig()
	if err != nil {
		logger.Error("❌ Failed to load config", err)
		return
	}

	config.ConnectDatabase(cfg)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Use(config.Cors(cfg))

	router.GET("/", handlers.Get)
	router.POST("/upload", handlers.Upload)
	router.GET("/:tag/:id", handlers.Fetch)
	router.DELETE("/:tag/:id", handlers.SoftDelete)

	svc, err := config.InitS3Session()
	if err != nil {
		logger.Error("❌ Failed to initialize AWS session", err)
		return
	}

	_, err = svc.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		logger.Error("❌ Failed to list S3 buckets", err)
		return
	} else {
		logger.Info("✅ Successfully connected to AWS S3")
	}

	address := ":" + strconv.Itoa(cfg.ApplicationPort)

	logger.Info("🚀 Server is running at " + cfg.ApplicationURL)

	if err := router.Run(address); err != nil {
		logger.Error("Error starting server", err)
	}
}
