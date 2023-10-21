package message

import (
    "github.com/jtonynet/api-gin-rest/config"
    "github.com/jtonynet/api-gin-rest/internal/message/interfaces"

    "github.com/jtonynet/api-gin-rest/internal/message/strategies/rabbitMQ"
)

func InitBroker(cfg config.MessageBroker) (interfaces.Broker, error) {
    switch cfg.Strategy {
    case "rabbitmq":
        return rabbitMQ.InitBroker(cfg)
    default:
        return nil, nil
    }
}
