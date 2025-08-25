package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Binit-Dhakal/Saarathi/pkg/events"
	"github.com/Binit-Dhakal/Saarathi/pkg/messagebus"
	"github.com/Binit-Dhakal/Saarathi/pkg/setup"
	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/application"
	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/messaging"
	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/repository/postgres"
	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/repository/redis"
)

// test
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

	instanceID, err := os.Hostname()
	if err != nil {
		log.Fatal("Hostname not set")
	}

	// for event from "trips" service to "ride-matching" service
	trip_create_queue_name := "ride-matching.trip-create"
	queue, err := setup.DeclareQueue(ch, trip_create_queue_name)
	if err != nil {
		log.Fatal("Failed to declare queue: ", err)
	}

	err = setup.BindQueue(ch, queue.Name, events.TripCreatedEvent, messagebus.TripEventsExchange)
	if err != nil {
		log.Fatal("Failed to bind queue: ", err)
	}

	// For event from "ride-matching service" to "driver state service"
	err = setup.DeclareExchange(ch, messagebus.TripOfferExchange, "topic")
	if err != nil {
		log.Fatal("Failed to declare exchange: ", err)
	}

	instance_queue_name := fmt.Sprintf("ride-matching.instance.%s", instanceID)
	queue, err = setup.DeclareQueue(ch, instance_queue_name)
	if err != nil {
		log.Fatal("Failed to declare queue: ", err)
	}

	routing_key := messagebus.RideMatchingRoutingKey(instanceID)
	err = setup.BindQueue(ch, queue.Name, routing_key, messagebus.TripOfferExchange)
	if err != nil {
		log.Fatal("Failed to bind queue to the exchange: ", err)
	}

	rideRepo := redis.NewRideMatchingRepository(client)
	redisMetaRepo := redis.NewCacheDriverMetaRepo(client)
	pgMetaRepo := postgres.NewPGMetaRepo(usersDB)
	bus := messagebus.NewRabbitMQBus(ch)
	matchingSvc := application.NewMatchingService(bus, rideRepo, redisMetaRepo, pgMetaRepo)

	messaging.ListenForTripEvents(ch, trip_create_queue_name, matchingSvc)
}
