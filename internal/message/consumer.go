package message

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (b *BrokerData) Consume(customHandler func(string) error) error {
	fmt.Printf("Queue bound to Exchange, starting Consume (consumer tag %q)", b.cfg.ConsumerTag)
	deliveries, err := b.channel.Consume(
		b.cfg.Queue,            // name
		b.cfg.ConsumerTag,      // consumerTag,
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

	go b.handle(customHandler, deliveries, b.done)
	return nil
}

func (b *BrokerData) Shutdown() error {
	// will close() the deliveries channel
	if err := b.channel.Cancel(b.cfg.ConsumerTag, true); err != nil {
		return fmt.Errorf("Consumer cancel failed: %s", err)
	}

	if err := b.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer fmt.Printf("AMQP shutdown OK")

	// wait for handle() to exit
	return <-b.done
}

func (b *BrokerData) handle(customHandler func(string) error, deliveries <-chan amqp.Delivery, done chan error) {
	var attempt int32
	requeue := true

	for d := range deliveries {
		fmt.Printf(
			"got %dB delivery: [%v] %q",
			len(d.Body),
			d.DeliveryTag,
			d.Body,
		)
		
		if err := customHandler(string(d.Body)); err != nil { //--> CUSTOM HANDLER
			if d.Headers["X-Attempt"] == nil {
				attempt = 1
			} else if attemptTemp, ok := d.Headers["X-Attempt"].(int32); ok && attemptTemp < b.cfg.MaxAttempts {				
				attempt = d.Headers["X-Attempt"].(int32) + 1
			} else {
                // Limite de tentativas atingido, mova para a dead message queue.
                if err := b.moveToDeadQueue(string(d.Body)); err != nil {
                    fmt.Println("Erro ao mover para a dead message queue:", err)
                } else {
					d.Ack(false)
					requeue = false
				}
            }

			if requeue {
				// No mundo ideal fariamos um Nack incrementando Headers["X-Attempt"]
				// Mas infelizmente a lib do rabbitMQ não me permitiu essa abordagem e fiz manualmente
				// err = d.Nack(false, true)
				// if err != nil {
				// 	fmt.Println("Erro ao reenfileirar a mensagem:", err)
				// }

				d.Ack(false)
				b.publish(string(d.Body), b.cfg.Exchange, b.cfg.RoutingKey, attempt)
		    }

			requeue = true
        } else {
			d.Ack(false)
		}
	}

	fmt.Printf("handle: deliveries channel closed")
	done <- nil
}

func (b *BrokerData) moveToDeadQueue(message string) error {
	initialAttempt := int32(0)
	return b.publish(message, b.cfg.ExchangeDL, b.cfg.RoutingKeyDL, initialAttempt)
}
