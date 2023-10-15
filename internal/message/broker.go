package message

import (
	"fmt"

	"github.com/jtonynet/api-gin-rest/config"
	"github.com/streadway/amqp"
)

type BrokerData struct {
	conn    *amqp.Connection
	channel *amqp.Channel

	Exchange      string
	ExchangeType  string
	RoutingKey    string
	Queue         string
	ConsumerTag   string

	done    chan error
}

var (
	Broker          *BrokerData
)

func InitBroker(cfg config.MessageBroker) error {
	strConn := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.User,
		cfg.Pass,
		cfg.Host,
		cfg.Port)

	conn, err := amqp.Dial(strConn)
	if err != nil {
		return err
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return err
	}

	fmt.Printf("got Channel, declaring %q Exchange (%q)", cfg.Exchange, cfg.ExchangeType)
	if err := channel.ExchangeDeclare(
		cfg.Exchange,     // name
		cfg.ExchangeType, // type
		true,             // durable
		false,            // auto-deleted
		false,            // internal
		false,            // noWait
		nil,              // arguments
	); err != nil {
		return err
	}

	_, err = channel.QueueDeclare(
		cfg.Queue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

    if err := channel.QueueBind(
        cfg.Queue,
        cfg.RoutingKey,
        cfg.Exchange,
        false,
        nil,
    ); err != nil {
        return err
    }

	// Reliable publisher confirms require confirm.select support from the connection.
	if cfg.ReliableMessagesEnable {
		fmt.Printf("enabling publishing confirms.")
		if err := channel.Confirm(false); err != nil {
			return fmt.Errorf("Channel could not be put into confirm mode: %s", err)
		}

		confirms := channel.NotifyPublish(make(chan amqp.Confirmation, 1))

		defer confirmOne(confirms)
	}

	Broker = &BrokerData{
		conn:    conn,
		channel: channel,

		Exchange:     cfg.Exchange,
		ExchangeType: cfg.ExchangeType,
		RoutingKey:   cfg.RoutingKey,
		Queue:        cfg.Queue,
		ConsumerTag:  cfg.ConsumerTag,

		done:    make(chan error),
	}

	return nil
}

func (b *BrokerData) CheckReadiness() error {
	if b.conn == nil || b.channel == nil {
		return fmt.Errorf("MessageBroker is not ready")
	}
	return nil
}
