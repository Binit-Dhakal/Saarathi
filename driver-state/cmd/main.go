package main

import (
	"net/http"
	"os"

	"github.com/Binit-Dhakal/Saarathi/driver-state/internal/application"
	"github.com/Binit-Dhakal/Saarathi/driver-state/internal/handlers/ws"
	"github.com/Binit-Dhakal/Saarathi/driver-state/internal/repository/redis"
	log "github.com/Binit-Dhakal/Saarathi/pkg/logger"
	"github.com/Binit-Dhakal/Saarathi/pkg/setup"
)

func main() {
	logger := log.NewStandardLogger()

	redisClient, err := setup.SetupRedis()
	if err != nil {
		logger.Error("failed to connect to redis", err)
		os.Exit(1)
	}
	defer redisClient.Close()

	rabbitConn, rabbitCh, err := setup.SetupRabbitMQ()
	if err != nil {
		logger.Error("RabbitMQ not setup", err)
		os.Exit(1)
	}
	defer rabbitConn.Close()
	defer rabbitCh.Close()

	redisRepo := redis.NewLocationRepo(redisClient)

	locationSvc := application.NewLocationService(redisRepo)
	driverStateHandler := ws.NewWebSocketHandler(locationSvc)

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", driverStateHandler.WsHandler)

	logger.Info("Driver state service started on 8084")
	err = http.ListenAndServe(":8084", mux)
	if err != nil {
		logger.Error("ListenAndServe: ", err)
	}
}
