package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/jtonynet/api-gin-rest/config"

	"github.com/jtonynet/api-gin-rest/internal/interfaces"
)

func ConfigInject(cfg config.API) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("cfg", cfg)
		c.Next()
	}
}

func MessageBrokerInject(messageBroker interfaces.MessageBroker) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("messageBroker", messageBroker)
		c.Next()
	}
}

func CacheClientInject(cacheClient interfaces.CacheClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("cacheClient", cacheClient)
		c.Next()
	}
}
