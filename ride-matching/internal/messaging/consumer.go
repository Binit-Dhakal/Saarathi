package messaging

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Binit-Dhakal/Saarathi/pkg/events"
	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/application"
	amqp "github.com/rabbitmq/amqp091-go"
)

func ListenForTripEvents(ch *amqp.Channel, queueName string, svc application.MatchingService) {
	msgs, err := ch.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	for d := range msgs {
		log.Printf("Received a message from RabbitMQ: %s", d.Body)
		go func(d amqp.Delivery) {
			var event events.TripEventCreated
			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Printf("Failed to unmarshal event: %v", err)
				_ = d.Ack(false)
				return
			}
			if err = svc.HandleNewTripEvent(context.Background(), &event); err != nil {
				log.Printf("Failed to handle trip event: %v", err)
				_ = d.Nack(false, false) // retry false for now
				return
			}

			_ = d.Ack(false)
		}(d)
	}
}
