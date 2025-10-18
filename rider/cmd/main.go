package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Binit-Dhakal/Saarathi/pkg/am"
	"github.com/Binit-Dhakal/Saarathi/pkg/contracts/proto/tripspb"
	"github.com/Binit-Dhakal/Saarathi/pkg/jetstream"
	"github.com/Binit-Dhakal/Saarathi/pkg/logger"
	"github.com/Binit-Dhakal/Saarathi/pkg/registry"
	"github.com/Binit-Dhakal/Saarathi/pkg/setup"
	"github.com/Binit-Dhakal/Saarathi/rider/internal/application"
	"github.com/Binit-Dhakal/Saarathi/rider/internal/handlers/messaging"
	"github.com/Binit-Dhakal/Saarathi/rider/internal/handlers/rest"
	"github.com/Binit-Dhakal/Saarathi/rider/internal/repository/redis"
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

	app.logger = logger.New(logger.LogConfig{
		Environment: app.cfg.Environment,
		LogLevel:    logger.Level(app.cfg.LogLevel),
	})

	return nil
}

func run() (err error) {
	var cfg RiderAppConfig
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
	defer app.cacheClient.Close()
	defer app.nc.Close()

	reg := registry.NewRegistry()

	err = tripspb.Registration(reg)
	if err != nil {
		return err
	}

	stream := jetstream.NewStream(cfg.Nats.Stream, app.js, app.logger)
	evtStream := am.NewEventStream(reg, stream)

	repo := redis.NewTripPayloadRepository(app.cacheClient)
	updateSvc := application.NewRiderUpdateService(repo)

	integrationHandler := messaging.NewIntegrationEventHandlers(updateSvc)
	err = messaging.RegisterIntegrationHandlers(evtStream, integrationHandler)
	if err != nil {
		return err
	}

	tripHandler := rest.NewTripUpdateHandler()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/trip/updates", tripHandler.TripUpdate)

	server := &http.Server{
		Addr:         ":8010",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	fmt.Println("starting server on :8010")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Fprint(os.Stderr, "server failed", err)
		os.Exit(1)
	}

	return

}
