package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jtonynet/api-gin-rest/config"
)

func ConfigInjectHandler(cfg config.API) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("cfg", cfg)
		c.Next()
	}
}

func ConfigManagerHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfgInterface, exists := c.Get("cfg")

		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Config not found in context",
			})
			c.Abort()
			return
		}

		cfg, ok := cfgInterface.(config.API)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid config type in context",
			})
			c.Abort()
			return
		}

		c.Set("cfg", cfg)

		c.Next()
	}
}
