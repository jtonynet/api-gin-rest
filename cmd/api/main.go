package main

import (
	"log"
	"log/slog"
	"time"

	"github.com/jtonynet/api-gin-rest/cmd/common"

	"github.com/jtonynet/api-gin-rest/internal/interfaces"
	"github.com/jtonynet/api-gin-rest/routes"
)

func main() {
	cfg, err := common.InitConfig()
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	RetryMaxElapsedTime := time.Duration(cfg.API.RetryMaxElapsedTimeInMs) * time.Millisecond

	err = common.InitDatabase(cfg.Database, RetryMaxElapsedTime)
	if err != nil {
		slog.Error("cannot initialize Database, error: %v", err)
	}

	var cacheClient interfaces.CacheClient
	cacheClient, err = common.NewCacheClient(cfg.Cache)
	if err != nil {
		slog.Error("cannot initialize cacheClient, error: %v", err)
	}

	var messageBroker interfaces.MessageBroker
	if cfg.API.FeatureFlags.PostAlunoAsMessageEnabled {
		messageBroker, err = common.NewMessageBroker(cfg.MessageBroker, RetryMaxElapsedTime)
		if err != nil {
			slog.Error("cannot initialize MessageBroker, error: %v", err)
		}
	}

	routes.HandleRequests(
		cfg.API,
		messageBroker,
		cacheClient,
	)
}
