package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ApplicationPort int
	ApplicationURL  string

	PostgresUser     string
	PostgresPassword string
	PostgresHost     string
	PostgresPort     int
	PostgresDatabase string
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

	config := &Config{
		ApplicationPort: applicationPort,
		ApplicationURL:  os.ExpandEnv(getEnv("APPLICATION_URL", "http://localhost:4000")),

		PostgresUser:     getEnv("POSTGRES_USER", "root"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", ""),
		PostgresHost:     getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:     postgresPort,
		PostgresDatabase: getEnv("POSTGRES_DATABASE", ""),
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
