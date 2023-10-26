package middlewares

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func CachedGetRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := c.MustGet("cfg").(config.API)
		var cacheClient interfaces.CacheClient
		var cacheKey string

		if cfg.FeatureFlags.CacheEnabled {
			cacheClient = c.MustGet("cacheClient").(interfaces.CacheClient)

			if cacheClient.IsConnected() {
				var err error
				cacheKey, err = getCacheKeyFromPath(c)
				if err != nil {
					slog.Error("cacheClient.getCacheKeyFromPath error: ", err)
					c.Abort()
				}

				cachedData, err := cacheClient.Get(cacheKey)
				slog.Info("cachedData ", cacheKey, cachedData)
				if err == nil {
					var returnData map[string]interface{}
					if err := json.Unmarshal([]byte(cachedData), &returnData); err != nil {
						slog.Error("middlewares:CachedGetRequest:json.Unmarshal error: ", err)
						c.Abort()
					}

					statusCodeReturn := http.StatusOK
					cachedHttpStatus := gjson.Get(cachedData, "HTTPStatusCode")
					if cachedHttpStatus.Exists() {
						statusCodeReturn = int(cachedHttpStatus.Int())
					}

					currentTime := time.Now()
					timeFormatted := currentTime.Format("15:04:05.000000")
					fmt.Println("MIDDLEWARE CachedGetRequest (HH:MM:SS.mmmuuu):", timeFormatted)

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
					slog.Error("middlewares:CachedGetRequest:json.Marshal error: ", err)
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

func CachedPostRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := c.MustGet("cfg").(config.API)
		var cacheClient interfaces.CacheClient
		var segmentCacheKey, cacheKey string

		if cfg.FeatureFlags.CacheEnabled {
			cacheClient = c.MustGet("cacheClient").(interfaces.CacheClient)

			if cacheClient.IsConnected() {
				var err error
				segmentCacheKey, err = getCacheKeyFromPath(c)
				if err != nil {
					slog.Error("middlewares:CachedPostRequest:cacheClient.getCacheKeyFromPath error: ", err)
					c.Abort()
				}

				UUID := uuid.New().String()
				c.Set("UUID", UUID)

				var msgJson map[string]interface{}
				msg := fmt.Sprintf(`{"HTTPStatusCode":202, "SegmentKey":"%s", "uuid":"%s", "Message":" in processing"}`, segmentCacheKey, UUID)
				err = json.Unmarshal([]byte(msg), &msgJson)
				if err != nil {
					slog.Error("middlewares:CachedPostRequest:json.Unmarshal error: ", err)
					c.Abort()
				}

				cacheKey = fmt.Sprintf("%s:%s", segmentCacheKey, UUID)

				// Request Post data no expires INFO cache. Only Request Get data cache expires
				// Request Get data has the same key as the Request Post, it overwrites.
				expiration := time.Duration(0 * time.Millisecond)

				err = cacheClient.Set(cacheKey, msg, expiration)
				if err != nil {
					slog.Error("middlewares:CachedPostRequest:cacheClient.Set error: ", err)
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
					slog.Info("middlewares:CachedRequest:json.Marshal error: ", err)
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

// func PublishPostRequestAsync() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		cfg := c.MustGet("cfg").(config.API)
// 		var messageBroker interfaces.Broker

// 		c.Next()

// 		if cfg.FeatureFlags.PostAlunoAsMessageEnabled {
// 			messageBroker := c.MustGet("messageBroker").(interfaces.Broker)
// 		}
// 	}
// }

func getCacheKeyFromPath(c *gin.Context) (string, error) {
	path := c.FullPath()

	hasColon := strings.Contains(path, ":")
	if hasColon {
		var paramName, paramKey string = "", ""

		pathSplited := strings.Split(path, "/")
		if len(pathSplited) > 2 {
			paramName = pathSplited[len(pathSplited)-2]
		}

		paramSplited := strings.Split(path, ":")
		if len(paramSplited) > 1 {
			paramKey = paramSplited[1]
		}

		if paramName == "" || paramKey == "" {
			return "", errors.New("improperly formatted path")
		}
		key := fmt.Sprintf("%s:%s", paramName, c.Params.ByName(paramKey))
		return key, nil
	}

	pathSplited := strings.Split(path, "/")
	if len(pathSplited) > 0 {
		return pathSplited[len(pathSplited)-1], nil
	}

	return "", errors.New("improperly formatted path")
}
