package cache

import (
	"github.com/jtonynet/api-gin-rest/config"
	"github.com/jtonynet/api-gin-rest/internal/cache/strategies/redis"
	"github.com/jtonynet/api-gin-rest/internal/interfaces"
)

func NewClient(cfg config.Cache) (interfaces.CacheClient, error) {
	switch cfg.Strategy {
	case "redis":
		return redis.NewClient(cfg)
	default:
		return nil, nil
	}
}
