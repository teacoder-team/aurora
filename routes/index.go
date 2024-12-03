package routes

import (
	"net/http"
	"storage/config"
	"time"

	"github.com/gin-gonic/gin"
)

func IndexHandler(c *gin.Context) {
	cfgTags, err := config.GetConfigTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load config tags"})
		return
	}

	currentTime := time.Now().Format("20060102150405")

	response := map[string]interface{}{
		"timestamp":    currentTime,
		"tags":         cfgTags.Tags,
		"jpeg_quality": cfgTags.JpegQuality,
	}

	c.JSON(http.StatusOK, response)
}
