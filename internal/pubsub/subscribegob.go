package pubsub

import (
	"bytes"
	"encoding/gob"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeGob[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType int,
	handler func(T) routing.AckType,
) error {
	ch, queue, err := DeclareAndBind(conn, exchange, queueName, key, simpleQueueType)
	if err != nil {
		return err
	}

	ch.Qos(10, 0, false)

	deliveryChan, err := ch.Consume(queue.Name, "", false, false, false, false, amqp.Table{})
	if err != nil {
		return err
	}

	go func() {
		for delivery := range deliveryChan {
			buffer := bytes.NewBuffer(delivery.Body)
			decoder := gob.NewDecoder(buffer)
			var body T
			err = decoder.Decode(&body)
			if err != nil {
				log.Fatalf("Fatal error in deserializing delivery. Error %v", err)
			}
			acktype := handler(body)
			switch acktype {
			case routing.ACK:
				delivery.Ack(false)
			case routing.NACKREQUEUE:
				delivery.Nack(false, true)
			case routing.NACKDISCARD:
				delivery.Nack(false, false)
			default:
				log.Fatalf("Invalid AckType: %v", acktype)
			}
		}
	}()

	return nil
}
