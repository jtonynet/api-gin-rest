package rabbitMQ

import (
    "fmt"
    "log"

    amqp "github.com/rabbitmq/amqp091-go"
)

/*
Fortemente baseado no exemplo da lib streadway consumer.go
https://github.com/streadway/amqp/blob/master/_examples/simple-consumer/consumer.go
*/

func (b *BrokerData) RunConsumer(userConsumerHandler func(string) error) error {
    log.Printf("Queue bound to Exchange, starting Consume (consumer tag %q)", b.cfg.ConsumerTag)
    deliveries, err := b.channel.Consume(
        b.cfg.Queue,		// name
        b.cfg.ConsumerTag,	// consumerTag,
        false,				// noAck
        false,				// exclusive
        false,				// noLocal
        false,				// noWait
        nil,				// arguments
    )
    if err != nil {
        log.Printf("Queue Consume: %s", err)
        return err
    }

    go b.handle(userConsumerHandler, deliveries, b.done)

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

    defer log.Printf("AMQP shutdown OK")

    // wait for handle() to exit
    return <-b.done
}

func (b *BrokerData) handle(userConsumerHandler func(string) error, deliveries <-chan amqp.Delivery, done chan error) {
    var attempt int32
    requeue := true

    b.userConsumerHandler = userConsumerHandler

    for d := range deliveries {
        
        if err := userConsumerHandler(string(d.Body)); err != nil {

            // TODO:
            // Ajustar para o descrito em https://www.rabbitmq.com/dlx.html de maneira automatizada. ADOTAR!
            // No mundo ideal fariamos um Nack incrementando Headers["X-Attempt"] Mas infelizmente a lib do
            // rabbitMQ não me permitiu essa abordagem e fiz manualmente. Estou adicionando e requeue na mão
            // incrementando X-Attempt. Pesquisar possiveis melhorias para essa gestão
            // https://www.rabbitmq.com/dlx.html
            // https://www.rabbitmq.com/quorum-queues.html#poison-message-handling
            // https://www.inanzzz.com/index.php/post/1p7m/creating-a-rabbitmq-dlx-dead-letter-exchange-example-with-golang
            
            if d.Headers["X-Attempt"] == nil {
                attempt = 1
            } else if attemptTemp, ok := d.Headers["X-Attempt"].(int32); ok && attemptTemp < b.cfg.MaxAttempts {
                attempt = d.Headers["X-Attempt"].(int32) + 1
            } else {
                if err := b.moveToDeadQueue(string(d.Body)); err != nil {
                    log.Println("Erro ao mover para a dead message queue:", err)
                } else {
                    d.Ack(false)
                    requeue = false
                }
            }

            if requeue {
                d.Ack(false)
                b.publish(string(d.Body), attempt, b.cfg.Exchange, b.cfg.RoutingKey)
            }

            requeue = true
        } else {
            d.Ack(false)
        }
        
    }

    log.Printf("handle: deliveries channel closed")
    done <- nil
}

func (b *BrokerData) moveToDeadQueue(message string) error {
    initialAttempt := int32(0)
    return b.publish(message, initialAttempt, b.cfg.ExchangeDL, b.cfg.RoutingKeyDL)
}

func (b *BrokerData) reconnectAndConsume() error {
    err := b.reconnect()
    if err != nil {
        return err
    }

    go b.RunConsumer(b.userConsumerHandler)
    return nil
}
