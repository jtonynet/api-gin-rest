package middlewares

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jtonynet/api-gin-rest/config"
	"github.com/tidwall/gjson"

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

// Tratamento adequedo de Middlewares para obter "Separation of Concerns"
// https://gin-gonic.com/docs/examples/custom-middleware/
func CachedRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := c.MustGet("cfg").(config.API)
		var cacheClient interfaces.CacheClient
		var cacheKey string

		if cfg.FeatureFlags.CacheEnabled {
			cacheClient = c.MustGet("cacheClient").(interfaces.CacheClient)

			if cacheClient.IsConnected() {
				paramKeyName, paramKeyValue, err := cacheClient.GetNameAndKeyFromPath(c.FullPath())
				if err != nil {
					slog.Error("cacheClient.GetNameAndKeyFromPath error: ", err)
					c.Abort()
				}
				cacheKey = fmt.Sprintf("%s:%s", paramKeyName, c.Params.ByName(paramKeyValue))

				cachedData, err := cacheClient.Get(cacheKey)
				slog.Info("cachedData ", cacheKey, cachedData)
				if err == nil {
					var returnData map[string]interface{}
					if err := json.Unmarshal([]byte(cachedData), &returnData); err != nil {
						slog.Error("middlewares:CachedRequest:json.Unmarshal error: ", err)
						c.Abort()
					}

					statusCodeReturn := http.StatusOK
					cachedHttpStatus := gjson.Get(cachedData, "HTTPStatusCode")
					if cachedHttpStatus.Exists() {
						statusCodeReturn = int(cachedHttpStatus.Int())
					}

					currentTime := time.Now()
					timeFormatted := currentTime.Format("15:04:05.000000")
					fmt.Println("MIDDLEWARE CachedRequest (HH:MM:SS.mmmuuu):", timeFormatted)

					c.JSON(statusCodeReturn, returnData)
					c.Abort()
				}
			}
		}

		c.Next()

		if cfg.FeatureFlags.CacheEnabled && cacheClient.IsConnected() {
			queryResult, queryResultExists := c.Get("queryResult")
			if queryResultExists {
				queryResultJSON, err := json.Marshal(queryResult)
				if err != nil {
					slog.Error("middlewares:CachedRequest:json.Marshal error: ", err)
					c.Abort()
				}

				err = cacheClient.Set(cacheKey, string(queryResultJSON), cacheClient.GetDefaultExpiration())
				if err != nil {
					slog.Error("cannot set key: %s error: %v", cacheKey, err)
				}
			}
		}
	}
}
