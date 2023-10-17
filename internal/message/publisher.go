package message

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (b *BrokerData) Publish(body string) error {
	initialAttempt := int32(0)
    return b.publish(body, initialAttempt, b.cfg.Exchange, b.cfg.RoutingKey)
}

func (b *BrokerData) publish(body string, attempt int32, exchange string, routingKey string) error {
	headers := amqp.Table{
		"X-Attempt": attempt,
	}

	if err := b.channel.Publish(
		exchange,		// publish to an exchange
		routingKey,		// routing to 0 or more queues
		false,			// mandatory
		false,			// immediate
		amqp.Publishing{
			Headers:         headers,
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            []byte(body),
			DeliveryMode:    amqp.Transient,	// 1=non-persistent, 2=persistent
			Priority:        0,					// 0-9
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		return fmt.Errorf("Exchange Publish: %s", err)
	}

	return nil
}

// One would typically keep a channel of publishings, a sequence number, and a
// set of unacknowledged sequence numbers and loop until the publishing channel
// is closed.
func confirmOne(confirms <-chan amqp.Confirmation) {
	fmt.Printf("waiting for confirmation of one publishing")

	if confirmed := <-confirms; confirmed.Ack {
		fmt.Printf("confirmed delivery with delivery tag: %d", confirmed.DeliveryTag)
	} else {
		fmt.Printf("failed delivery of delivery tag: %d", confirmed.DeliveryTag)
	}
}
