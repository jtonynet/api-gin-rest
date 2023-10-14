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

	go func() {
		for {
			err := message.Broker.Consume(message.Broker.QueueAluno, func(body []byte) {
				//log.Printf("Mensagem recebida: %s", string(body))

				msg := string(body)

				var aluno models.Aluno
				err := json.Unmarshal([]byte(msg), &aluno)
				if err != nil {
					fmt.Println("Erro na análise JSON:", err)
					return
				}

				// fmt.Printf("Aluno: %+v\n", aluno)
				err = database.DB.Create(&aluno).Error
				if err != nil {
					fmt.Println("REQUEUE!")
				}

			})
			if err != nil {
				log.Printf("Erro ao consumir mensagens: %v", err)
			}

			duracao := 10 * time.Millisecond
			time.Sleep(duracao)
		}
	}()

	select {}
}
