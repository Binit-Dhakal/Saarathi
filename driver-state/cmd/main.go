package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Binit-Dhakal/Saarathi/driver-state/internal/application"
	"github.com/Binit-Dhakal/Saarathi/driver-state/internal/handlers/messaging"
	"github.com/Binit-Dhakal/Saarathi/driver-state/internal/handlers/ws"
	"github.com/Binit-Dhakal/Saarathi/driver-state/internal/repository/redis"
	"github.com/Binit-Dhakal/Saarathi/pkg/events"
	log "github.com/Binit-Dhakal/Saarathi/pkg/logger"
	"github.com/Binit-Dhakal/Saarathi/pkg/messagebus"
	"github.com/Binit-Dhakal/Saarathi/pkg/setup"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	logger := log.NewStandardLogger()

	redisClient, err := setup.SetupRedis()
	if err != nil {
		logger.Error("failed to connect to redis", err)
		os.Exit(1)
	}
	defer redisClient.Close()

	rConn, rCh, err := setup.SetupRabbitMQ()
	if err != nil {
		logger.Error("RabbitMQ not setup", err)
		os.Exit(1)
	}
	defer rConn.Close()
	defer rCh.Close()

	err = buildRabbitMQEntity(rCh)
	if err != nil {
		logger.Error("RabbitMQ setup", err)
		os.Exit(1)
	}

	bus := messagebus.NewRabbitMQBus(rCh)

	locationRepo := redis.NewLocationRepo(redisClient)
	wsRepo := redis.NewWSRepo(redisClient)

	presenceSvc := application.NewPresenceService(wsRepo)
	locationSvc := application.NewLocationService(locationRepo)
	offerSvc := application.NewOfferService(bus)
	driverStateHandler := ws.NewWebSocketHandler(locationSvc, presenceSvc, offerSvc)

	offerHandler := messaging.NewTripOfferHandler(driverStateHandler)

	var wg sync.WaitGroup

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	instanceID, _ := os.Hostname()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var event events.TripOfferRequest
		bus.Subscribe(
			context.Background(),
			fmt.Sprintf("trip-offer-request-%s", instanceID),
			event.EventName(),
			offerHandler.HandleOfferRequest,
		)
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", driverStateHandler.WsHandler)

	logger.Info("Driver state service started on :8084")
	err = http.ListenAndServe(":8084", mux)
	if err != nil {
		logger.Error("ListenAndServe: ", err)
	}
}

func buildRabbitMQEntity(ch *amqp.Channel) error {
	instanceID, err := os.Hostname()
	if err != nil {
		return err
	}

	configs := []setup.QueueConfig{
		{
			Name:       fmt.Sprintf("trip-offer-request-%s", instanceID),
			Exchange:   messagebus.TripOfferExchange,
			RoutingKey: messagebus.DriverRoutingKey(events.EventOfferRequest, instanceID),
			Type:       "topic",
			Durable:    true,
		},
	}

	err = setup.SetupQueues(ch, configs)
	if err != nil {
		return err
	}

	return nil
}
