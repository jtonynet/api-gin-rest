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
					slog.Error("middlewares:CachedGetRequest:cacheClient.getCacheKeyFromPath error: ", err)
					c.Abort()
				}

				cachedData, err := cacheClient.Get(cacheKey)
				if err == nil {
					var returnData interface{}
					if err := json.Unmarshal([]byte(cachedData), &returnData); err != nil {
						slog.Error("middlewares:CachedGetRequest:json.Unmarshal error: ", err)
						c.Abort()
					}

					statusCodeReturn := http.StatusOK
					cachedHttpStatus := gjson.Get(cachedData, "HTTPStatusCode")
					if cachedHttpStatus.Exists() {
						statusCodeReturn = int(cachedHttpStatus.Int())
					}

					c.JSON(statusCodeReturn, returnData)
					c.Abort()
				}
			}
		}

		c.Next()

		if cfg.FeatureFlags.CacheEnabled && cacheClient.IsConnected() {
			setCacheResult(c, cacheKey)
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

				cacheKey = fmt.Sprintf("%s:%s", segmentCacheKey, UUID)

				var msgJson map[string]interface{}
				msg := fmt.Sprintf(`{"HTTPStatusCode":202, "SegmentKey":"%s", "uuid":"%s", "Message":" in processing"}`, segmentCacheKey, UUID)
				err = json.Unmarshal([]byte(msg), &msgJson)
				if err != nil {
					slog.Error("middlewares:CachedPostRequest:json.Unmarshal error: ", err)
					c.Abort()
				}

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
			setCacheResult(c, cacheKey)
		}
	}
}

func CachedDeleteRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		cfg := c.MustGet("cfg").(config.API)

		if cfg.FeatureFlags.CacheEnabled {
			cacheClient := c.MustGet("cacheClient").(interfaces.CacheClient)

			if cacheClient.IsConnected() {
				cacheKey, err := getCacheKeyFromPath(c)
				if err != nil {
					slog.Error("middlewares:CachedEditAndDeleteRequest:cacheClient.getCacheKeyFromPath error: ", err)
					c.Abort()
				}
				cacheClient.Delete(cacheKey)
			}
		}
	}
}

func getCacheKeyFromPath(c *gin.Context) (string, error) {
	path := c.FullPath()
	var key string
	queryStringCacheKey := ""

	queryParameters := c.Request.URL.Query()
	for queryKey, queryValues := range queryParameters {
		for _, queryValue := range queryValues {
			queryStringCacheKey = queryStringCacheKey + fmt.Sprintf(":%s=%s", queryKey, queryValue)
		}
	}

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
		key = fmt.Sprintf("%s:%s%s", paramName, c.Params.ByName(paramKey), queryStringCacheKey)
		return key, nil
	}

	pathSplited := strings.Split(path, "/")
	if len(pathSplited) > 0 {
		key = fmt.Sprintf("%s%s", pathSplited[len(pathSplited)-1], queryStringCacheKey)
		return key, nil
	}

	return "", errors.New("improperly formatted path")
}

func setCacheResult(c *gin.Context, cacheKey string) {
	result, resultExists := c.Get("result")
	if resultExists {
		cacheClient := c.MustGet("cacheClient").(interfaces.CacheClient)

		resultJSON, err := json.Marshal(result)
		if err != nil {
			slog.Error("middlewares:setCacheResult:json.Marshal error: ", err)
			c.Abort()
		}

		err = cacheClient.Set(cacheKey, string(resultJSON), cacheClient.GetDefaultExpiration())
		if err != nil {
			slog.Error("cannot set key: %s error: %v", cacheKey, err)
		}
	}
}
