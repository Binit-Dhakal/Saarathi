package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	log "github.com/Binit-Dhakal/Saarathi/pkg/logger"
	"github.com/Binit-Dhakal/Saarathi/pkg/messagebus"
	"github.com/Binit-Dhakal/Saarathi/pkg/rest/httpx"
	"github.com/Binit-Dhakal/Saarathi/pkg/rest/jsonutil"
	"github.com/Binit-Dhakal/Saarathi/pkg/setup"
	"github.com/Binit-Dhakal/Saarathi/trips/internals/application"
	"github.com/Binit-Dhakal/Saarathi/trips/internals/handlers/rest"
	"github.com/Binit-Dhakal/Saarathi/trips/internals/repository/postgres"
	"github.com/Binit-Dhakal/Saarathi/trips/internals/repository/redis"
)

func main() {
	logger := log.NewStandardLogger()
	dbpool, err := setup.SetupPostgresDB()
	if err != nil {
		logger.Error("failed to connect to the database", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	redisClient, err := setup.SetupRedis()
	if err != nil {
		logger.Error("failed to connect to redis", err)
		os.Exit(1)
	}
	defer redisClient.Close()

	rConn, rCh, err := setup.SetupRabbitMQ()
	if err != nil {
		logger.Error("failed to connect to rabbitMQ", err)
		os.Exit(1)
	}
	defer rConn.Close()
	defer rCh.Close()

	bus := messagebus.NewRabbitMQBus(rCh)

	redisRepo := redis.NewRedisFareRepository(redisClient)
	tripRepo := postgres.NewTripRepository(dbpool)

	rideService := application.NewRideService(redisRepo, tripRepo, bus)
	routeService := application.NewRouteService()

	jsonWriter := jsonutil.NewWriter()
	jsonReader := jsonutil.NewReader()
	errorResponder := httpx.NewErrorResponder(jsonWriter, logger)

	tripHandler := rest.NewTripHandler(rideService, routeService, jsonReader, jsonWriter, errorResponder)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/fare/preview", tripHandler.PreviewFare)
	mux.HandleFunc("/api/v1/fare/confirm", tripHandler.ConfirmFare)

	server := &http.Server{
		Addr:         ":8082",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	fmt.Println("starting server on :8082")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Fprint(os.Stderr, "server failed", err)
		os.Exit(1)
	}
}
