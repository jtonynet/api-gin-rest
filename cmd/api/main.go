package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jtonynet/api-gin-rest/cmd/common"
	"github.com/jtonynet/api-gin-rest/routes"

	"github.com/jtonynet/api-gin-rest/internal/interfaces"

	"github.com/jtonynet/api-gin-rest/internal/cache"
)

func main() {
	cfg, err := common.InitConfig()
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	RetryMaxElapsedTime := time.Duration(cfg.API.RetryMaxElapsedTimeInMs) * time.Millisecond

	err = common.InitDatabase(cfg.Database, RetryMaxElapsedTime)
	if err != nil {
		log.Fatal("cannot initialize Database: ", err)
	}

	var cacheClient interfaces.CacheClient
	cacheClient, err = cache.NewClient(cfg.Cache)
	if err != nil {
		fmt.Println("NAO SE CONECTOU AO CACHE!")
	}

	var messageBroker interfaces.Broker
	if cfg.API.FeatureFlags.PostAlunoAsMessageEnabled {
		messageBroker, err = common.NewMessageBroker(cfg.MessageBroker, cacheClient, RetryMaxElapsedTime)
		if err != nil {
			log.Fatal("cannot initialize MessageBroker: ", err)
		}
	}

	routes.HandleRequests(
		cfg.API,
		messageBroker,
		cacheClient,
	)
}
