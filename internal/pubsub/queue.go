package pubsub

import (
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func DeclareAndBind(conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType int, // 1 for transient 0 for durable
) (*amqp.Channel, amqp.Queue, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}
	table := amqp.Table{}
	table["x-dead-letter-exchange"] = routing.ExchangePerilDeadLetter
	queue, err := ch.QueueDeclare(queueName,
		simpleQueueType == 0,
		simpleQueueType == 1,
		simpleQueueType == 1,
		false,
		table)
	if err != nil {
		return nil, queue, err
	}
	err = ch.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	return ch, queue, nil
}
