package message

import (
	"fmt"
	"github.com/streadway/amqp"
)

var ConsumerChannel = make(chan string)

func (b *BrokerData) Consume() error {
	fmt.Printf("Queue bound to Exchange, starting Consume (consumer tag %q)", b.ConsumerTag)
	deliveries, err := b.channel.Consume(
		b.Queue,            // name
		b.ConsumerTag,      // consumerTag,
		false,      		// noAck
		false,      		// exclusive
		false,      		// noLocal
		false,      		// noWait
		nil,        		// arguments
	)
	if err != nil {
		fmt.Printf("Queue Consume: %s", err)
		return err
	}

	go handle(deliveries, b.done)
	return nil
}

func (b *BrokerData) Shutdown() error {
	// will close() the deliveries channel
	if err := b.channel.Cancel(b.ConsumerTag, true); err != nil {
		return fmt.Errorf("Consumer cancel failed: %s", err)
	}

	if err := b.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer fmt.Printf("AMQP shutdown OK")

	// wait for handle() to exit
	return <-b.done
}

func handle(deliveries <-chan amqp.Delivery, done chan error) {
	for d := range deliveries {
		fmt.Printf(
			"got %dB delivery: [%v] %q",
			len(d.Body),
			d.DeliveryTag,
			d.Body,
		)

		fmt.Println(string(d.Body))
		ConsumerChannel <- string(d.Body)

		d.Ack(false)
	}
	fmt.Printf("handle: deliveries channel closed")
	done <- nil
}
