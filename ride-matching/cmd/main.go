package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Binit-Dhakal/Saarathi/pkg/am"
	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
	"github.com/Binit-Dhakal/Saarathi/pkg/jetstream"
	"github.com/Binit-Dhakal/Saarathi/pkg/logger"
	"github.com/Binit-Dhakal/Saarathi/pkg/proto/offerspb"
	"github.com/Binit-Dhakal/Saarathi/pkg/proto/rmspb"
	"github.com/Binit-Dhakal/Saarathi/pkg/registry"
	"github.com/Binit-Dhakal/Saarathi/pkg/setup"
	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/application"
	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/handlers/messaging"
	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/infrastructure"
	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/logging"
	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/repository/postgres"
	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/repository/redis"
	"github.com/nats-io/nats.go"
)

func main() {
	err := run()
	if err != nil {
		fmt.Printf("RMS service exitted abnormally: %v\n", err)
		os.Exit(1)
	}
}

func infraSetup(app *app) (err error) {
	app.usersDB, err = setup.SetupPostgresDB(app.cfg.PG.Conn)
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
	var cfg MatchAppConfig
	cfg, err = InitConfig()
	if err != nil {
		return err
	}

	app := &app{
		cfg: cfg,
	}

	err = infraSetup(app)
	if err != nil {
		return
	}
	defer app.usersDB.Close()
	defer app.cacheClient.Close()
	defer app.nc.Close()

	reg := registry.NewRegistry()

	if err := offerspb.Registration(reg); err != nil {
		return err
	}

	if err := rmspb.Registration(reg); err != nil {
		return err
	}

	domainDispatcher := ddd.NewEventDispatcher[ddd.Event]()

	stream := jetstream.NewStream(cfg.Nats.Stream, app.js, app.logger)
	eventStream := am.NewEventStream(reg, stream)

	_, cancel := context.WithCancel(context.Background())

	rideRepo := redis.NewRideMatchingRepository(app.cacheClient)
	redisMetaRepo := redis.NewCacheDriverMetaRepo(app.cacheClient)
	pgMetaRepo := postgres.NewPGMetaRepo(app.usersDB)
	availabilityRepo := redis.NewDriverAvailableRepo(app.cacheClient)

	driverInfoSvc := application.NewDriverInfoService(redisMetaRepo, pgMetaRepo, availabilityRepo)

	driverInfoAdapter := infrastructure.NewDriverInfoAdapter(driverInfoSvc)

	matchingSvc := application.NewMatchingService(domainDispatcher, rideRepo, driverInfoAdapter, driverInfoAdapter)

	domainHandler := logging.LogEventHandlerAccess(
		messaging.NewDomainEventHandlers(eventStream),
		"DomainEvents",
		app.logger,
	)
	integrationHandler := logging.LogEventHandlerAccess(
		messaging.NewIntegrationEventHandlers(matchingSvc),
		"IntegrationEvents",
		app.logger,
	)

	messaging.RegisterDomainEventHandlers(domainDispatcher, domainHandler)
	err = messaging.RegisterIntegrationEventHandlers(eventStream, integrationHandler)
	if err != nil {
		return err
	}

	// Wait for shutdown signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	<-sigs
	fmt.Println("Shutdown signal received")
	cancel()

	eventStream.Unsubscribe()
	fmt.Println("Graceful shutdown")

	return nil
}
