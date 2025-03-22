package config

import (
	"orion/pkg/logger"
	"orion/pkg/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func InitS3Session() (*s3.S3, error) {
	cfg, err := utils.LoadConfig()
	if err != nil {
		logger.Error("❌ Failed to load config: %v", err)
		return nil, err
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.S3Region),
		Endpoint:    aws.String(cfg.S3Endpoint),
		Credentials: credentials.NewStaticCredentials(cfg.S3AccessKeyId, cfg.S3SecretAccessKey, ""),
	})
	if err != nil {
		logger.Error("❌ Failed to initialize S3 session: %v", err)
		return nil, err
	}

	return s3.New(sess), nil
}
