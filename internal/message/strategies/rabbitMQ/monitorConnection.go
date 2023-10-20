package rabbitMQ

import (
	"fmt"
    "time"

	//amqp "github.com/rabbitmq/amqp091-go"
)

func (b *BrokerData) MonitorConnection() {
    ticker := time.NewTicker(1 * time.Second)

    for {
        select {
        case <-ticker.C:
            if !b.IsConnected() {
                ok := b.Reconnect()
                if !ok {
                   fmt.Println("Nao conseguiu reconectar")
                }
            }
        case <-b.done:
            return
        }
    }
}


func (b *BrokerData) Reconnect() bool {
	conn, channel, err := connect(b.cfg)
	if err != nil {
		return false
	}

    fmt.Println("RECONECTOU")
	b.conn = conn
	b.channel = channel
    b.done = make(chan error)

	return true
}
