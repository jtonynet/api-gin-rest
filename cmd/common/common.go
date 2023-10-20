package common

import (
	"log"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/jtonynet/api-gin-rest/config"
	"github.com/jtonynet/api-gin-rest/internal/database"
	"github.com/jtonynet/api-gin-rest/internal/message"
	"github.com/jtonynet/api-gin-rest/internal/message/interfaces"
)

func InitConfig() (*config.Config, error) {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func InitDatabase(cfg config.Database, RetryMaxElapsedTime time.Duration) error {
	retryCfg := backoff.NewExponentialBackOff()
	retryCfg.MaxElapsedTime = RetryMaxElapsedTime

	var dbErr error
	err := backoff.RetryNotify(func() error {
		dbErr = database.Init(cfg)
		return dbErr
	}, retryCfg, func(err error, t time.Duration) {
		log.Printf("Retrying connect to Database after error: %v", err)
	})

	return err
}

func InitMessageBroker(cfg config.MessageBroker, RetryMaxElapsedTime time.Duration) (interfaces.Broker, error) {
	retryCfg := backoff.NewExponentialBackOff()
	retryCfg.MaxElapsedTime = RetryMaxElapsedTime

	var msgBrokerErr error
	var messageBroker interfaces.Broker
	
	err := backoff.RetryNotify(func() error {
		messageBroker, msgBrokerErr = message.InitBroker(cfg)
		return msgBrokerErr
	}, retryCfg, func(err error, t time.Duration) {
		log.Printf("Retrying connect to MessageBroker after error: %v", err)
	})

	return messageBroker, err
}
