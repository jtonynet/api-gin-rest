package middlewares

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jtonynet/api-gin-rest/config"

	//"golang.org/x/exp/slog"

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

		log.Println("OLHA O UIUIUIID")
		log.Println(uuid)

		if cfg.FeatureFlags.CacheEnabled {
			if cacheClient.IsConnected() {
				log.Println("TEM ACESSO")
				// Agora, use o cacheClient fornecido como argumento
				msg, err := cacheClient.Get("f1")
				if err != nil {
					log.Printf("Cannot cache get, error: %v", err)
				}
				log.Printf("Cached message: %s", msg)

				msg, err = cacheClient.Get(uuid)
				if err != nil {
					log.Printf("Cannot cache get, error: %v", err)
				}
				log.Printf("Cached message: %s", msg)

			} else {
				log.Println("NOT A NOT!!!")
			}
		} else {
			log.Println("NAO ENTREI NO CACHE")
		}

		c.Next()
	}
}

// func Cache(
// 	key string,
// 	queryFunc func() (interface{}, error),
// 	c *interfaces.CacheClient,
// gin.HandlerFunc {
// 	cachedValue, err := c.Get(key)
// 	if err == nil {
// 		return cachedValue, nil
// 	}

// 	result, err := queryFunc()
// 	if err != nil {
// 		return nil, err
// 	}

// 	expiration := time.Duration(c.cfg.Expiration)
// 	err = c.Set(key, result, expiration)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return result, nil
// }
