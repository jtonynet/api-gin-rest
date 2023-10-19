package message

import (
    "github.com/jtonynet/api-gin-rest/config"
	"github.com/jtonynet/api-gin-rest/internal/message/strategies/rabbitMQ"
	"github.com/jtonynet/api-gin-rest/internal/message/interfaces"
)

func InitBroker(cfg config.MessageBroker) (error, interfaces.Broker) {
    switch cfg.Strategy {
    case "rabbitmq":
        return rabbitMQ.InitBroker(cfg)
    default:
        return nil, nil
    }
}