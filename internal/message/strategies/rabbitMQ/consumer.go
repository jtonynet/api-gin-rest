package rabbitMQ

import (
	"fmt"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/tidwall/gjson"
)

/*
Fortemente baseado no exemplo da lib streadway consumer.go
https://github.com/streadway/amqp/blob/master/_examples/simple-consumer/consumer.go
*/

func (b *Broker) RunConsumer(consumerHandler func(string) (string, error)) error {
	slog.Info("Queue bound to Exchange, starting Consume (consumer tag %q)", b.cfg.ConsumerTag)
	deliveries, err := b.channel.Consume(
		b.cfg.Queue,       // name
		b.cfg.ConsumerTag, // consumerTag,
		false,             // noAck
		false,             // exclusive
		false,             // noLocal
		false,             // noWait
		nil,               // arguments
	)
	if err != nil {
		slog.Info("Queue Consume: %s", err)
		return err
	}

	go b.handle(consumerHandler, deliveries, b.done)

	return nil
}

func (b *Broker) Shutdown() error {
	// will close() the deliveries channel
	if err := b.channel.Cancel(b.cfg.ConsumerTag, true); err != nil {
		return fmt.Errorf("consumer cancel failed: %s", err)
	}

	if err := b.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer slog.Info("AMQP shutdown OK")

	// wait for handle() to exit
	return <-b.done
}

// ESSE METODO DEVERIA SER DECORADO PELO CACHE
func (b *Broker) handle(consumerHandler func(string) (string, error), deliveries <-chan amqp.Delivery, done chan error) {
	var attempt int32
	requeue := true

	b.consumerHandler = consumerHandler

	for d := range deliveries {

		msgValue, err := consumerHandler(string(d.Body))
		msgUUID := gjson.Get(msgValue, "uuid")
		msgKey := fmt.Sprintf("aluno:%s", msgUUID)

		fmt.Println("//------------------------------------------")
		fmt.Println(msgKey)
		fmt.Println("//------------------------------------------")

		//err := json.Unmarshal([]byte(msg), &msgValue)
		//msgValueJson := json.Marshal(msgValue)

		if err != nil {

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
					slog.Error("cmd:worker:main:messageBroker:RunConsumer:handle:b.moveToDeadQueue error %v", err)
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

		if b.cacheClient != nil {
			b.cacheClient.Delete(msgKey)
			err = b.cacheClient.Set(msgKey, string(msgValue), b.cacheClient.GetDefaultExpiration())
			if err != nil {
				slog.Error("cmd:worker:main:messageBroker:RunConsumer:handle:b.cacheClient.Set error set%v", err)
			}
		}
	}

	slog.Info("handle: deliveries channel closed")
	done <- nil
}

func (b *Broker) moveToDeadQueue(message string) error {
	initialAttempt := int32(0)
	return b.publish(message, initialAttempt, b.cfg.ExchangeDL, b.cfg.RoutingKeyDL)
}

func (b *Broker) reconnectAndConsume() error {
	err := b.reconnect()
	if err != nil {
		return err
	}

	go b.RunConsumer(b.consumerHandler)
	return nil
}
