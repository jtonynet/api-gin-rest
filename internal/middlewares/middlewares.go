package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jtonynet/api-gin-rest/config"

	"github.com/jtonynet/api-gin-rest/internal/message/interfaces"
)

func ConfigInjectHandler(cfg config.API) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("cfg", cfg)
		c.Next()
	}
}

func MessageBrokerInjectHandler(messageBroker interfaces.Broker) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("messageBroker", messageBroker)
		c.Next()
	}
}
