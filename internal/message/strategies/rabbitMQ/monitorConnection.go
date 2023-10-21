package rabbitMQ

import (
    "log"
    "time"

	"github.com/cenkalti/backoff"
)

func (b *BrokerData) MonitorConnection() {

    for {
        if !b.IsConnected() {
            var err error
            RetryMaxElapsedTime := 3 * time.Minute

            retryCfg := backoff.NewExponentialBackOff()
            retryCfg.MaxElapsedTime = RetryMaxElapsedTime

            if b.userHandler != nil { 
                var msgConsumerAgainErr error
                err = backoff.RetryNotify(func() error {
                    msgConsumerAgainErr = b.reconnectAndConsume()
                    return msgConsumerAgainErr
                }, retryCfg, func(err error, t time.Duration) {
                    log.Printf("Tentando voltar a consumir: %v", err)
                })
            } else {
                var msgReconnectErr error
                err = backoff.RetryNotify(func() error {
                    msgReconnectErr = b.reconnect()
                    return msgReconnectErr
                }, retryCfg, func(err error, t time.Duration) {
                    log.Printf("Tentando voltar a publicar: %v", err)
                })
            }

            if err != nil {
                log.Printf("Nao deu para voltar a consumir error: %v", err)
            }
        }
        time.Sleep(1 * time.Second)
    }
}

