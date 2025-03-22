package config

import (
	"orion/pkg/utils"

	"github.com/gin-gonic/gin"
)

func Cors(cfg *utils.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		for _, allowedOrigin := range cfg.AllowedOrigins {
			if allowedOrigin == origin {
				c.Header("Access-Control-Allow-Origin", origin)
				break
			}
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")

		c.Next()
	}
}
