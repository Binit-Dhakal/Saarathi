package messagebus

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQBus struct {
	ch *amqp.Channel
}

func NewRabbitMQBus(ch *amqp.Channel) *RabbitMQBus {
	return &RabbitMQBus{
		ch: ch,
	}
}

func (r *RabbitMQBus) Publish(ctx context.Context, exchange string, topic string, message any) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return r.ch.Publish(
		exchange,
		topic,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (r *RabbitMQBus) Consume(ctx context.Context, queue string) (<-chan amqp.Delivery, error) {
	return r.ch.Consume(
		queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
}
