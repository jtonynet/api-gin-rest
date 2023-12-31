package message

import (
	"github.com/jtonynet/api-gin-rest/config"

	"github.com/jtonynet/api-gin-rest/internal/interfaces"
	"github.com/jtonynet/api-gin-rest/internal/message/strategies/rabbitMQ"
)

func NewBroker(cfg config.MessageBroker) (interfaces.MessageBroker, error) {
	switch cfg.Strategy {
	case "rabbitmq":
		return rabbitMQ.NewBroker(cfg)
	default:
		return nil, nil
	}
}
