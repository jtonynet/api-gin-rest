package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/jtonynet/api-gin-rest/config"
	"github.com/jtonynet/api-gin-rest/internal/database"
	"github.com/jtonynet/api-gin-rest/internal/message"
	"github.com/jtonynet/api-gin-rest/models"
)

// @title api-gin-rest
// @version 0.0.3
// @description Estudo API Rest em Golang com Gin
// @contact.name API GIN Support
// @termsOfService https://github.com/jtonynet/api-gin-rest
// @contact.url https://github.com/jtonynet/api-gin-rest
// @contact.email learningingenuity@gmail.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8083
func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	retryCfg := backoff.NewExponentialBackOff()
	retryCfg.MaxElapsedTime = time.Duration(cfg.API.RetryMaxElapsedTimeInMs) * time.Millisecond

	var dbErr, msgBrokerErr error

	err = backoff.RetryNotify(func() error {
		dbErr = database.Init(cfg.Database)
		return dbErr
	}, retryCfg, func(err error, t time.Duration) {
		log.Printf("Retrying connect to Database after error: %v", err)
	})

	if err != nil {
		log.Fatal("cannot initialize Database: ", dbErr)
	}

	err = backoff.RetryNotify(func() error {
		msgBrokerErr = message.InitBroker(cfg.MessageBroker)
		return msgBrokerErr
	}, retryCfg, func(err error, t time.Duration) {
		log.Printf("Retrying connect to MessageBroker after error: %v", err)
	})

	if err != nil {
		log.Fatal("cannot initialize MessageBroker: ", msgBrokerErr)
	}

	err = message.Broker.Consume()
	if err != nil {
		log.Fatal("cannot consume messages from Broker: ", msgBrokerErr)
	}

    go func() {
        for msg := range message.ConsumerChannel {
			var aluno models.Aluno
			err := json.Unmarshal([]byte(msg), &aluno)
			if err != nil {
				fmt.Println("Erro na análise JSON:", err)
				return
			}

			err = database.DB.Create(&aluno).Error
			if err != nil {
				fmt.Println("REQUEUE!")
			}
        }
    }()

	select {}
}
