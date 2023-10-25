package middlewares

import (
	"log/slog"

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

func MessageBrokerInject(messageBroker interfaces.Broker) gin.HandlerFunc {
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

func CachedRequest(cacheClient interfaces.CacheClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := c.MustGet("cfg").(config.API)
		uuid := c.Params.ByName("uuid")

		if cfg.FeatureFlags.CacheEnabled {
			if cacheClient.IsConnected() {
				msg, err := cacheClient.Get("f1")
				if err != nil {
					slog.Error("Cannot cache get, error: %v", err)
				}
				slog.Info("Cached message: %s", msg)

				msg, err = cacheClient.Get(uuid)
				if err != nil {
					slog.Error("Cannot cache get, error: %v", err)
				}
				slog.Info("Cached message: %s", msg)

			} else {
				slog.Info("CacheClient is disconnected")
			}
		} else {
			slog.Info("CacheClient is disabled")
		}

		c.Next()
	}
}
