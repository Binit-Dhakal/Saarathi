package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Binit-Dhakal/Saarathi/pkg/events"
	"github.com/Binit-Dhakal/Saarathi/pkg/messagebus"
	"github.com/Binit-Dhakal/Saarathi/pkg/setup"
	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/application"
	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/handlers/messaging"
	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/repository/postgres"
	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/repository/redis"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	client, err := setup.SetupRedis()
	if err != nil {
		log.Fatal("Couldn't setup redis:", err)
	}

	usersDB, err := setup.SetupPostgresDB()
	if err != nil {
		log.Fatal("Couldn't setup users db:", err)
	}

	conn, ch, err := setup.SetupRabbitMQ()
	if err != nil {
		log.Println("RabbitMQ error: ", err)
		os.Exit(1)
	}
	defer conn.Close()
	defer ch.Close()

	buildRabbitMQEntity(ch)

	rideRepo := redis.NewRideMatchingRepository(client)
	redisMetaRepo := redis.NewCacheDriverMetaRepo(client)
	pgMetaRepo := postgres.NewPGMetaRepo(usersDB)
	availabilityRepo := redis.NewDriverAvailableRepo(client)
	presenceRepo := redis.NewPresenceRepo(client)

	bus := messagebus.NewRabbitMQBus(ch)

	matchingSvc := application.NewMatchingService(bus, rideRepo)
	driverInfoSvc := application.NewDriverInfoService(redisMetaRepo, pgMetaRepo, availabilityRepo)
	presenceSvc := application.NewPresenceService(presenceRepo)

	handler := messaging.NewTripEventHandler(matchingSvc, driverInfoSvc, presenceSvc, bus)

	listenForTripEvents(ch, "ride-matching.trip-create", handler)
}

func buildRabbitMQEntity(ch *amqp.Channel) {
	instanceID, err := os.Hostname()
	if err != nil {
		log.Fatal("Hostname not set")
	}

	configs := []setup.QueueConfig{
		{Name: "ride-matching.trip-create", Exchange: messagebus.TripEventsExchange, RoutingKey: events.EventTripCreated, Type: "topic", Durable: true},
		{Name: fmt.Sprintf("ride-matching.instance.%s", instanceID), Exchange: messagebus.TripOfferExchange, RoutingKey: events.EventOfferResponse, Type: "topic", Durable: true},
	}

	err = setup.SetupQueues(ch, configs)
	if err != nil {
		log.Fatal(err)
	}
}

func listenForTripEvents(ch *amqp.Channel, queueName string, handler *messaging.TripEventHandler) {
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
			if err = handler.HandleTripEvent(d.Body); err != nil {
				log.Printf("Failed to handle trip event: %v", err)
				_ = d.Nack(false, false) // retry false for now
				return
			}
			_ = d.Ack(false)
		}(d)
	}
}
