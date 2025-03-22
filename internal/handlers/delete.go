package handlers

import (
	"net/http"
	"orion/internal/repositories"
	"orion/pkg/utils"

	"github.com/gin-gonic/gin"
)

func SoftDelete(c *gin.Context) {
	cfg, err := utils.LoadConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error loading config"})
		return
	}

	secretKey := c.GetHeader("X-Upload-Secret")
	if secretKey != cfg.UploadSecretKey {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid secret key"})
		return
	}

	tag := c.Param("tag")
	id := c.Param("id")

	fileRepo := repositories.NewFileRepository()

	fileRecord, err := fileRepo.GetFileByID(id, tag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch file from database"})
		return
	}

	if fileRecord == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	err = fileRepo.SoftDeleteFile(id, tag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark file as deleted"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File marked as deleted successfully"})
}
