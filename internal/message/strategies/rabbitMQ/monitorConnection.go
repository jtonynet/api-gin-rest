package rabbitMQ

import (
	"fmt"
    "time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (b *BrokerData) MonitorConnection() {
    ticker := time.NewTicker(1 * time.Second)

    for {
        select {
        case <-ticker.C:
            if !b.IsConnected() {
                ok := b.reconnect()
                if !ok {
                   fmt.Println("Nao conseguiu reconectar")
                } else {
                    fmt.Println("RECONECTOU")
                }

            }
        case <-b.done:
            return
        }
    }
}

func (b *BrokerData) IsConnected() bool {
    if b.conn == nil || b.channel == nil {
        return false
    }

    if b.conn.IsClosed() {
        return false
    }

    if b.channel.IsClosed() {
		return false
    }

    return true
}

func (b *BrokerData) reconnect() bool {
	conn, err := amqp.Dial(strConn)
	if err != nil {
		fmt.Println(strConn)
		return false
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return false
	}
	
	b.conn = conn
	b.channel = channel

	return true
}
