package pubsub

import (
	"encoding/json"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeJSON[T any](conn *amqp.Connection, exchange, queueName, key string, simpleQueueType int, handler func(T) routing.AckType) error {
	ch, queue, err := DeclareAndBind(conn, exchange, queueName, key, simpleQueueType)
	if err != nil {
		return err
	}

	deliveryChan, err := ch.Consume(queue.Name, "", false, false, false, false, amqp.Table{})
	if err != nil {
		return err
	}

	go func() {

		for delivery := range deliveryChan {
			var val T
			err = json.Unmarshal(delivery.Body, &val)
			if err != nil {
				log.Fatalf("Fatal error in deserializing delivery. Error: %v", err)
			}
			acktype := handler(val)
			switch acktype {
			case routing.ACK:
				delivery.Ack(false)
				// fmt.Printf(("Acked"))
			case routing.NACKREQUEUE:
				delivery.Nack(false, true)
				// fmt.Printf("Nacked and requeued")
			case routing.NACKDISCARD:
				delivery.Nack(false, false)
				// fmt.Printf("Nacked and discarded")
			default:
				log.Fatalf("Invalid AckType: %v", acktype)
			}

		}
	}()

	return nil
}
