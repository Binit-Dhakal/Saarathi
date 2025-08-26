package messagebus

import (
	"context"
	"encoding/json"

	"github.com/Binit-Dhakal/Saarathi/pkg/events"
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

func (r *RabbitMQBus) Publish(ctx context.Context, exchange string, routingKey string, event events.Event) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return r.ch.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (r *RabbitMQBus) Subscribe(ctx context.Context, queue string, eventName string, handler EventHandler) error {
	msgs, err := r.ch.Consume(
		queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case d, ok := <-msgs:
				if !ok {
					return
				}

				evt, err := events.DecodeEvent(eventName, d.Body)
				if err != nil {
					d.Nack(false, false)
					continue
				}

				if err := handler(ctx, evt); err != nil {
					d.Nack(false, false)
				} else {
					d.Ack(false)
				}

			}
		}
	}()

	return nil
}
