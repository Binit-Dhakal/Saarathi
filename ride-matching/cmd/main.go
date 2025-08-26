package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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
	fmt.Println("Subscribing to trip created event")

	var event events.TripEventCreated
	bus.Subscribe(
		context.Background(),
		event.EventName(),
		event.EventName(),
		handler.HandleTripEvent,
	)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	<-sigs
	fmt.Println("Shutting down")
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
