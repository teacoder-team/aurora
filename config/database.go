package config

import (
	"fmt"
	"orion/internal/models"
	"orion/pkg/logger"
	"orion/pkg/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase(cfg *utils.Config) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
		cfg.PostgresHost,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDatabase,
		cfg.PostgresPort,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("❌ Failed to establish database connection: %v", err)
	}

	logger.Info("✅ Database connection established successfully")

	err = DB.AutoMigrate(&models.File{})
	if err != nil {
		logger.Error("❌ Error during migration: %v", err)
	}

	logger.Info("✅ Database migrated successfully")
}
