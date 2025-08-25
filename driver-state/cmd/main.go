package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Binit-Dhakal/Saarathi/driver-state/internal/application"
	"github.com/Binit-Dhakal/Saarathi/driver-state/internal/handlers/ws"
	"github.com/Binit-Dhakal/Saarathi/driver-state/internal/messaging"
	"github.com/Binit-Dhakal/Saarathi/driver-state/internal/repository/redis"
	log "github.com/Binit-Dhakal/Saarathi/pkg/logger"
	"github.com/Binit-Dhakal/Saarathi/pkg/messagebus"
	"github.com/Binit-Dhakal/Saarathi/pkg/setup"
)

func main() {
	logger := log.NewStandardLogger()

	instanceID, err := os.Hostname()
	if err != nil {
		logger.Error("Hostname not set", nil)
		os.Exit(1)
	}

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

	err = setup.DeclareExchange(rCh, messagebus.TripOfferExchange, "topic")
	if err != nil {
		logger.Error("Error in declaring exchange", err)
		os.Exit(1)
	}

	queueName := fmt.Sprintf("driver-state.ride-matching.%s.queue", instanceID)
	queue, err := setup.DeclareQueue(rCh, queueName)
	if err != nil {
		logger.Error("Error in declaring queue", err)
		os.Exit(1)
	}

	routingKey := messagebus.DriverRoutingKey(instanceID)
	err = setup.BindQueue(rCh, queue.Name, routingKey, messagebus.TripOfferExchange)
	if err != nil {
		logger.Error("Error in binding queue", err)
		os.Exit(1)
	}

	locationRepo := redis.NewLocationRepo(redisClient)
	wsRepo := redis.NewWSRepo(redisClient)

	presenceSvc := application.NewPresenceService(wsRepo)
	locationSvc := application.NewLocationService(locationRepo)
	driverStateHandler := ws.NewWebSocketHandler(locationSvc, presenceSvc)

	offerSvc := application.NewOfferService(driverStateHandler)
	go messaging.ListenForOfferEvents(rCh, queueName, offerSvc)

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", driverStateHandler.WsHandler)

	logger.Info("Driver state service started on :8084")
	err = http.ListenAndServe(":8084", mux)
	if err != nil {
		logger.Error("ListenAndServe: ", err)
	}
}
