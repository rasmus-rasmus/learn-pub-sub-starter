package pubsub

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	serialisedVal, err := json.Marshal(val)
	if err != nil {
		return err
	}

	ctx := context.Background()
	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        serialisedVal,
	}
	return ch.PublishWithContext(ctx, exchange, key, false, false, msg)
}
