package rabbitMQ

import (
	"fmt"

	"github.com/jtonynet/api-gin-rest/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type BrokerData struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	cfg		config.MessageBroker
	done    chan error
}

var Broker *BrokerData

func InitBroker(cfg config.MessageBroker) (*BrokerData, error) {
	strConn := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.User,
		cfg.Pass,
		cfg.Host,
		cfg.Port)

	conn, err := amqp.Dial(strConn)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	exchangeDeclare(channel, cfg.Exchange, cfg.ExchangeType)
	exchangeDeclare(channel, cfg.ExchangeDL, cfg.ExchangeTypeDL)

	queueDeclare(channel, cfg.Queue)
	queueDeclare(channel, cfg.QueueDL)

	queueBind(channel, cfg.Queue, cfg.RoutingKey, cfg.Exchange)
	queueBind(channel, cfg.QueueDL, cfg.RoutingKeyDL, cfg.ExchangeDL)

	// Reliable publisher confirms require confirm.select support from the connection.
	if cfg.ReliableMessagesEnable {
		fmt.Printf("enabling publishing confirms.")
		if err := channel.Confirm(false); err != nil {
			return nil, fmt.Errorf("Channel could not be put into confirm mode: %s", err)
		}

		confirms := channel.NotifyPublish(make(chan amqp.Confirmation, 1))

		defer confirmOne(confirms)
	}

	Broker = &BrokerData{
		conn:    conn,
		channel: channel,
		cfg: cfg,
		done:    make(chan error),
	}

	return Broker, nil
}

func exchangeDeclare(channel *amqp.Channel, exchange string, exchangeType string) error {
	fmt.Printf("got Channel, declaring %q Exchange (%q)", exchange, exchangeType)
	if err := channel.ExchangeDeclare(
		exchange,		// name
		exchangeType,	// type
		true,			// durable
		false,			// auto-deleted
		false,			// internal
		false,			// noWait
		nil,			// arguments
	); err != nil {
		return err
	}

	return nil
}

func queueDeclare(channel *amqp.Channel, queue string) error {
	_, err := channel.QueueDeclare(
		queue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}

func queueBind(channel *amqp.Channel, queue string, routingKey string, exchange string) error {
	if err := channel.QueueBind(
		queue,
		routingKey,
		exchange,
		false,
		nil,
	); err != nil {
		return err
	}

	return nil
}

func (b *BrokerData) CheckReadiness() error {
	if b.conn != nil || b.channel != nil {
		err := exchangeDeclare(b.channel, b.cfg.Exchange, b.cfg.ExchangeType)
		if err != nil {
			return fmt.Errorf("MessageBroker is not ready")
		}
	} 

	return nil
}
