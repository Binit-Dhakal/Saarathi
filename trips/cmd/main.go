package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Binit-Dhakal/Saarathi/pkg/am"
	"github.com/Binit-Dhakal/Saarathi/pkg/jetstream"
	"github.com/Binit-Dhakal/Saarathi/pkg/logger"
	log "github.com/Binit-Dhakal/Saarathi/pkg/logger"
	"github.com/Binit-Dhakal/Saarathi/pkg/rest/httpx"
	"github.com/Binit-Dhakal/Saarathi/pkg/rest/jsonutil"
	"github.com/Binit-Dhakal/Saarathi/pkg/setup"
	"github.com/Binit-Dhakal/Saarathi/trips/internals/application"
	"github.com/Binit-Dhakal/Saarathi/trips/internals/handlers/rest"
	"github.com/Binit-Dhakal/Saarathi/trips/internals/repository/postgres"
	"github.com/Binit-Dhakal/Saarathi/trips/internals/repository/redis"
	"github.com/nats-io/nats.go"
)

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Trips service exitted abnormally: %v\n", err)
		os.Exit(1)
	}
}

func infraSetup(app *app) (err error) {
	app.tripsDB, err = setup.SetupPostgresDB(app.cfg.PG.Conn)
	if err != nil {
		return err
	}

	app.cacheClient, err = setup.SetupRedis(app.cfg.Redis.CacheURL)
	if err != nil {
		return err
	}

	app.nc, err = nats.Connect(app.cfg.Nats.URL)
	if err != nil {
		return err
	}

	app.js, err = setup.SetupJetStream(app.cfg.Nats.Stream, app.nc)
	if err != nil {
		return err
	}

	app.logger = logger.New(log.LogConfig{
		Environment: app.cfg.Environment,
		LogLevel:    logger.Level(app.cfg.LogLevel),
	})

	return nil
}

func run() (err error) {
	var cfg TripAppConfig
	cfg, err = InitConfig()
	if err != nil {
		return err
	}

	app := &app{
		cfg: cfg,
	}

	err = infraSetup(app)
	if err != nil {
		return err
	}

	defer app.tripsDB.Close()
	defer app.cacheClient.Close()
	defer app.nc.Close()

	stream := jetstream.NewStream(cfg.Nats.Stream, app.js, app.logger)
	eventStream := am.NewEventPublisher(stream)

	redisRepo := redis.NewRedisFareRepository(app.cacheClient)
	tripRepo := postgres.NewTripRepository(app.tripsDB)

	rideService := application.NewRideService(redisRepo, tripRepo, eventStream)
	routeService := application.NewRouteService()

	jsonWriter := jsonutil.NewWriter()
	jsonReader := jsonutil.NewReader()
	errorResponder := httpx.NewErrorResponder(jsonWriter, app.logger)

	tripHandler := rest.NewTripHandler(rideService, routeService, jsonReader, jsonWriter, errorResponder)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/fare/preview", tripHandler.PreviewFare)
	mux.HandleFunc("/api/v1/fare/confirm", tripHandler.ConfirmFare)
	mux.HandleFunc("/api/v1/trip/updates", tripHandler.TripUpdate)

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

	return
}
