package rabbitMQ

import (
    "fmt"
    "log"

    "github.com/jtonynet/api-gin-rest/config"
    amqp "github.com/rabbitmq/amqp091-go"
)

/*
Fortemente baseado nos exemplos da lib streadway
https://github.com/streadway/amqp/tree/master/_examples
*/

type BrokerData struct {
    conn    *amqp.Connection
    channel *amqp.Channel
    cfg		config.MessageBroker
    done    chan error

    userConsumerHandler func(string) error
}

var Broker *BrokerData
var strConn string

func InitBroker(cfg config.MessageBroker) (*BrokerData, error) {

    conn, channel, err := connect(cfg)
    if err != nil {
        return nil, err
    }

    // Reliable publisher confirms require confirm.select support from the connection.
    if cfg.ReliableMessagesEnable {
        log.Printf("enabling publishing confirms.")
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

func (b *BrokerData) IsConnected() bool {
    if b.conn == nil || b.channel == nil {
        log.Println("conn and channel nil values")
        return false
    }

    if b.conn.IsClosed() {
        log.Println("conn closed")
        return false
    }

    if b.channel.IsClosed() {
        log.Println("channel closed")
        return false
    }

    return true
}

func connect(cfg config.MessageBroker) (*amqp.Connection, *amqp.Channel, error) {
    strConn = fmt.Sprintf("amqp://%s:%s@%s:%s/",
        cfg.User,
        cfg.Pass,
        cfg.Host,
        cfg.Port)

    conn, err := amqp.Dial(strConn)
    if err != nil {
        return nil, nil, err
    }

    channel, err := conn.Channel()
    if err != nil {
        conn.Close()
        return nil, nil, err
    }

    exchangeDeclare(channel, cfg.Exchange, cfg.ExchangeType)
    exchangeDeclare(channel, cfg.ExchangeDL, cfg.ExchangeTypeDL)

    queueDeclare(channel, cfg.Queue)
    queueDeclare(channel, cfg.QueueDL)

    queueBind(channel, cfg.Queue, cfg.RoutingKey, cfg.Exchange)
    queueBind(channel, cfg.QueueDL, cfg.RoutingKeyDL, cfg.ExchangeDL)

    return conn, channel, nil
}

func exchangeDeclare(channel *amqp.Channel, exchange string, exchangeType string) error {
    log.Printf("got Channel, declaring %q Exchange (%q)", exchange, exchangeType)
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

func (b *BrokerData) reconnect() error {
    conn, channel, err := connect(b.cfg)
    if err != nil {
        return err
    }

    b.conn = conn
    b.channel = channel
    b.done = make(chan error)

    return nil
}
