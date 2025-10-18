package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Binit-Dhakal/Saarathi/offers/internal/application"
	"github.com/Binit-Dhakal/Saarathi/offers/internal/handlers/messaging"
	"github.com/Binit-Dhakal/Saarathi/offers/internal/logging"
	"github.com/Binit-Dhakal/Saarathi/offers/internal/repository/postgres"
	"github.com/Binit-Dhakal/Saarathi/offers/internal/repository/redis"
	"github.com/Binit-Dhakal/Saarathi/pkg/am"
	"github.com/Binit-Dhakal/Saarathi/pkg/contracts/proto/driverspb"
	"github.com/Binit-Dhakal/Saarathi/pkg/contracts/proto/offerspb"
	"github.com/Binit-Dhakal/Saarathi/pkg/contracts/proto/rmspb"
	"github.com/Binit-Dhakal/Saarathi/pkg/contracts/proto/tripspb"
	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
	"github.com/Binit-Dhakal/Saarathi/pkg/jetstream"
	"github.com/Binit-Dhakal/Saarathi/pkg/logger"
	"github.com/Binit-Dhakal/Saarathi/pkg/registry"
	"github.com/Binit-Dhakal/Saarathi/pkg/setup"
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
	app.DB, err = setup.SetupPostgresDB(app.cfg.PG.Conn)
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

	app.logger = logger.New(logger.LogConfig{
		Environment: app.cfg.Environment,
		LogLevel:    logger.Level(app.cfg.LogLevel),
	})

	return nil
}

func run() (err error) {
	var cfg OfferSvcConfig
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
	defer app.DB.Close()
	defer app.cacheClient.Close()
	defer app.nc.Close()

	reg := registry.NewRegistry()

	err = tripspb.Registration(reg)
	if err != nil {
		return err
	}
	err = offerspb.Registration(reg)
	if err != nil {
		return err
	}

	err = rmspb.Registration(reg)
	if err != nil {
		return err
	}

	err = driverspb.Registration(reg)
	if err != nil {
		return err
	}

	stream := jetstream.NewStream(cfg.Nats.Stream, app.js, app.logger)
	domainDispatcher := ddd.NewEventDispatcher[ddd.Event]()

	evtStream := am.NewEventStream(reg, stream)

	tripReadRepo := postgres.NewTripReadModelRepo(app.DB)
	candidatesRepo := redis.NewTripCandidatesRepo(app.cacheClient)
	driverRepo := redis.NewDriverAvailabilityRepo(app.cacheClient)

	offerSvc := application.NewService(candidatesRepo, driverRepo, tripReadRepo, domainDispatcher)

	integrationHandlers := logging.LogEventHandlerAccess(
		messaging.NewIntegrationEventHandlers(offerSvc),
		"IntegrationEvents",
		app.logger,
	)
	domainHandlers := logging.LogEventHandlerAccess(
		messaging.NewDomainEventHandlers(evtStream),
		"DomainEvents",
		app.logger,
	)

	messaging.RegisterDomainEventHandlers(domainDispatcher, domainHandlers)
	err = messaging.RegisterIntegrationHandlers(evtStream, integrationHandlers)
	if err != nil {
		return err
	}

	// Wait for shutdown signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	<-sigs
	fmt.Println("Shutdown signal received")

	evtStream.Unsubscribe()

	fmt.Println("Graceful shutdown")

	return nil
}
