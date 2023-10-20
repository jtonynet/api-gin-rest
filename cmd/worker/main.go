package main

import (
	"log"
	"time"

	"github.com/jtonynet/api-gin-rest/cmd/common"
	"github.com/jtonynet/api-gin-rest/cmd/worker/handlers"

	"github.com/jtonynet/api-gin-rest/internal/message/interfaces"
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
	// go messageBroker.MonitorConnection()
	
	err = messageBroker.RunConsumer(handlers.InsertAluno)
	if err != nil {
		log.Fatal("cannot consume messages from Broker: ", err)
	}

	
    go func (mb interfaces.Broker,  userHandler func(string) error) {
		for {
			// Verificar a conexão do messageBroker continuamente
			if !mb.IsConnected() {
				log.Println("Ta DESconnectada")

				log.Println("Connection is down. Attempting to reestablish...")
				ok := mb.Reconnect()
				if ok {
					log.Println("Reconnectei")
				} else {
					log.Println("NAO Reconnectei")
				}

				err = mb.RunConsumer(userHandler)
				if err == nil {
				    log.Println("COMSUMINDO")
				} else {
				    log.Println("NAO TO TINTINDO NADA!")
				}
			} else {
				log.Println("Ta connectada")
			}
			time.Sleep(1 * time.Second)
		}
	}(messageBroker, handlers.InsertAluno)

	select {}
}
