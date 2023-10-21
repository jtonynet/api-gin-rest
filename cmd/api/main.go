package main

import (
    "log"
    "time"

    "github.com/jtonynet/api-gin-rest/cmd/common"
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
        log.Fatal("cannot initialize Database: ", err)
    }

    messageBroker, err := common.InitMessageBroker(cfg.MessageBroker, RetryMaxElapsedTime)
    if err != nil {
        log.Fatal("cannot initialize MessageBroker: ", err)
    }
    //Monitora quedas de conexão e tenta reconectar
    go messageBroker.AutoReconnect() 

    routes.HandleRequests(cfg.API, messageBroker)
}
