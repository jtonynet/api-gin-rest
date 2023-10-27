package rabbitMQ

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/cenkalti/backoff"
)

func (b *Broker) autoReconnect() {
	tickInterval := 100 * time.Millisecond
	ticker := time.NewTicker(tickInterval)
	defer ticker.Stop()

	for range ticker.C {
		if !b.IsConnected() {
			var err error
			RetryMaxElapsedTime := time.Duration(b.cfg.AutoReconnectRetryMaxElapsedInMs) * time.Millisecond

			retryCfg := backoff.NewExponentialBackOff()
			retryCfg.MaxElapsedTime = RetryMaxElapsedTime

			if b.consumerHandler != nil {
				var msgReconnectAndConsumeErr error
				err = backoff.RetryNotify(func() error {
					msgReconnectAndConsumeErr = b.reconnectAndConsume()
					return msgReconnectAndConsumeErr
				}, retryCfg, func(err error, t time.Duration) {
					slog.Error("Attempting to resume consumer: %v", err)
				})

			} else {
				var msgReconnectErr error
				err = backoff.RetryNotify(func() error {
					msgReconnectErr = b.reconnect()
					return msgReconnectErr
				}, retryCfg, func(err error, t time.Duration) {
					slog.Error("Attempting to resume publish: %v", err)
				})
			}

			if err != nil {
				panic(fmt.Sprintf("Connection lost, failed to reestablish the connection: %v", err))
			}
		}
	}
}
