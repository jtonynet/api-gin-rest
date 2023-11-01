package main

import (
	"log/slog"
	"time"

	"github.com/jtonynet/api-gin-rest/cmd/common"
	"github.com/jtonynet/api-gin-rest/cmd/worker/messageHandlers"
	"github.com/jtonynet/api-gin-rest/internal/cache"
	"github.com/jtonynet/api-gin-rest/internal/interfaces"
)

func main() {
	cfg, err := common.InitConfig()
	if err != nil {
		slog.Error("cannot load config error: %v", err)
	}

	RetryMaxElapsedTime := time.Duration(cfg.API.RetryMaxElapsedTimeInMs) * time.Millisecond

	err = common.InitDatabase(cfg.Database, RetryMaxElapsedTime)
	if err != nil {
		slog.Error("cannot initialize Database, error: %v", err)
	}

	var cacheClient interfaces.CacheClient
	if cfg.API.FeatureFlags.CacheEnabled {
		cacheClient, err = cache.NewClient(cfg.Cache)
		if err != nil {
			slog.Error("cannot initialize CacheClient, error: %v", err)
		}
	}

	if cfg.API.FeatureFlags.PostAlunoAsMessageEnabled {
		messageBroker, err := common.NewMessageBroker(
			cfg.MessageBroker,
			RetryMaxElapsedTime,
		)

		if err != nil {
			slog.Error("cannot initialize MessageBroker, error: %v", err)
		}

		insertAluno := messageHandlers.NewInsertAluno()
		inserAlunoCached := messageHandlers.NewCached(insertAluno, cacheClient, cfg.MessageBroker.Queue)

		err = messageBroker.RunConsumer(inserAlunoCached.Execute)
		if err != nil {
			slog.Error("cannot consume messages from Broker, error: %v", err)
		}

		select {}
	}
}
