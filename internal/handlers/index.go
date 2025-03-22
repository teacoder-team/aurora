package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Get(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello! You've reached Orion Server – a file storage service for teacoder.ru",
	})
}
