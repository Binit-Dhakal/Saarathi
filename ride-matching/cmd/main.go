package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
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
	ctx, cancel := context.WithCancel(context.Background())

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

	var wg sync.WaitGroup

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	wg.Add(1)
	go func() {
		defer wg.Done()
		bus.Subscribe(
			ctx,
			"ride-matching-trip-create",
			events.EventTripCreated,
			handler.HandleTripEvent,
		)
	}()

	<-sigs
	fmt.Println("Shutdown signal received")
	cancel()

	wg.Wait()
	fmt.Println("Graceful shutdown")
}

func buildRabbitMQEntity(ch *amqp.Channel) {
	configs := []setup.QueueConfig{
		{
			Name:       "ride-matching-trip-create",
			Exchange:   messagebus.TripEventsExchange,
			RoutingKey: events.EventTripCreated,
			Type:       "topic",
			Durable:    true,
		},
		{
			Name:       "ride-matching-offer-response",
			Exchange:   messagebus.TripOfferExchange,
			RoutingKey: events.EventOfferResponse,
			Type:       "topic",
			Durable:    true,
		},
	}

	err := setup.SetupQueues(ch, configs)
	if err != nil {
		log.Fatal(err)
	}
}
