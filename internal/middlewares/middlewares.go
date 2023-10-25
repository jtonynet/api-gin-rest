package middlewares

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
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

func CachedRequest() gin.HandlerFunc {
	//func CachedRequest(c *gin.Context) {
	return func(c *gin.Context) {
		cfg := c.MustGet("cfg").(config.API)
		var cacheClient interfaces.CacheClient
		var param, paramValue string

		if cfg.FeatureFlags.CacheEnabled {

			cacheClient = c.MustGet("cacheClient").(interfaces.CacheClient)

			if cacheClient.IsConnected() {
				param = extractParamFromPath(c.FullPath())
				paramValue = c.Params.ByName(param)
				// TODO: capturar queryString completa para cachear TODAS requests inclusive
				// as parametrizadas. CachedRequest deve ser agnostica quanto a rotas e params

				cachedData, err := cacheClient.Get(paramValue)
				if err == nil {
					var returnData map[string]interface{}
					if err := json.Unmarshal([]byte(cachedData), &returnData); err != nil {
						slog.Error("middlewares:CachedRequest:json.Unmarshal error: ", err)
						c.JSON(http.StatusInternalServerError, gin.H{
							"error": "Erro ao recuperar dados",
						})
						c.Abort()
					}

					statusCodeReturn := http.StatusOK
					cachedHttpStatus := gjson.Get(cachedData, "HTTPStatusCode")
					if cachedHttpStatus.Exists() {
						statusCodeReturn = int(cachedHttpStatus.Int())
					}

					currentTime := time.Now()
					timeFormatted := currentTime.Format("15:04:05.000000")
					fmt.Println("1 - RETORNANDO DO MIDDLEWARE (HH:MM:SS.mmmuuu):", timeFormatted)

					c.JSON(statusCodeReturn, returnData)
					c.Abort()
				}
			}
		}

		c.Next()

		if cfg.FeatureFlags.CacheEnabled {
			queryResult := c.GetString("queryResult")

			resultJSON, err := json.Marshal(queryResult)
			if err != nil {
				slog.Error("cannot unmarshal resultJSON error: %v", err)
			}

			err = cacheClient.Set(paramValue, string(resultJSON), cacheClient.GetDefaultExpiration())
			if err != nil {
				slog.Error("cannot set key: %s error: %v", paramValue, err)
			}
		}
	}
}

func extractParamFromPath(path string) string {
	pathSplited := strings.Split(path, ":")
	if len(pathSplited) > 1 {
		return pathSplited[1]
	}
	return path
}
