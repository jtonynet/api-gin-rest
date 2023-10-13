package message

import (
	"fmt"

	"github.com/jtonynet/api-gin-rest/config"
	"github.com/streadway/amqp"
)

type BrokerData struct {
	conn    *amqp.Connection
	channel *amqp.Channel

	ExchangeAluno   string
	RoutingKeyAluno string
	QueueAluno      string
}

var (
	Broker          *BrokerData
	DefaultExchange = "amq.default"
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

	if cfg.ExchangeAluno == DefaultExchange {
		_, err = channel.QueueDeclare(
			cfg.QueueAluno,
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}

	Broker = &BrokerData{
		conn:    conn,
		channel: channel,

		ExchangeAluno:   cfg.ExchangeAluno,
		RoutingKeyAluno: cfg.RoutingKeyAluno,
		QueueAluno:      cfg.QueueAluno,
	}

	return nil
}

func (b *BrokerData) Publish(exchange, routingKey, message string) error {

	exchangeTopub := func() string {
		if exchange == DefaultExchange {
			return ""
		}
		return exchange
	}()

	err := b.channel.Publish(
		exchangeTopub,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	fmt.Println(message)
	return err
}

func (b *BrokerData) Consume(queue string, handler func([]byte)) error {
	messages, err := b.channel.Consume(
		queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for message := range messages {
			handler(message.Body)
		}
	}()
	return nil
}

func (b *BrokerData) CheckReadiness() error {
	if b.conn == nil || b.channel == nil {
		return fmt.Errorf("MessageBroker is not ready")
	}
	return nil
}
