package utils

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	ApplicationPort int
	ApplicationURL  string
	AllowedOrigins  []string

	PostgresUser     string
	PostgresPassword string
	PostgresHost     string
	PostgresPort     int
	PostgresDatabase string

	S3Region          string
	S3Endpoint        string
	S3BucketName      string
	S3AccessKeyId     string
	S3SecretAccessKey string

	UseS3       bool
	AllowedTags string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	applicationPort, err := strconv.Atoi(getEnv("APPLICATION_PORT", "4000"))
	if err != nil {
		return nil, err
	}

	postgresPort, err := strconv.Atoi(getEnv("POSTGRES_PORT", "5433"))
	if err != nil {
		return nil, err
	}

	allowedOrigins := getEnv("ALLOWED_ORIGINS", "")
	allowedOriginsList := strings.Split(allowedOrigins, ",")

	config := &Config{
		ApplicationPort: applicationPort,
		ApplicationURL:  os.ExpandEnv(getEnv("APPLICATION_URL", "")),
		AllowedOrigins:  allowedOriginsList,

		PostgresUser:     getEnv("POSTGRES_USER", ""),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", ""),
		PostgresHost:     getEnv("POSTGRES_HOST", ""),
		PostgresPort:     postgresPort,
		PostgresDatabase: getEnv("POSTGRES_DATABASE", ""),

		S3Region:          getEnv("S3_REGION", ""),
		S3Endpoint:        getEnv("S3_ENDPOINT", ""),
		S3BucketName:      getEnv("S3_BUCKET_NAME", ""),
		S3AccessKeyId:     getEnv("S3_ACCESS_KEY_ID", ""),
		S3SecretAccessKey: getEnv("S3_SECRET_ACCESS_KEY", ""),

		UseS3:       getEnv("S3_REGION", "") != "" && getEnv("S3_ENDPOINT", "") != "",
		AllowedTags: getEnv("ALLOWED_TAGS", ""),
	}

	if config.UseS3 {
		if config.S3AccessKeyId == "" || config.S3SecretAccessKey == "" {
			log.Fatalf("❌ Missing S3 credentials: S3_ACCESS_KEY_ID and S3_SECRET_ACCESS_KEY must be set.")
			return nil, err
		}
	} else {
		log.Fatalf("❌ S3 is not configured. Make sure S3_REGION and S3_ENDPOINT are set if you intend to use S3.")
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
