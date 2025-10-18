package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Binit-Dhakal/Saarathi/pkg/am"
	"github.com/Binit-Dhakal/Saarathi/pkg/contracts/proto/driverspb"
	"github.com/Binit-Dhakal/Saarathi/pkg/contracts/proto/tripspb"
	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
	"github.com/Binit-Dhakal/Saarathi/pkg/jetstream"
	"github.com/Binit-Dhakal/Saarathi/pkg/logger"
	"github.com/Binit-Dhakal/Saarathi/pkg/registry"
	"github.com/Binit-Dhakal/Saarathi/pkg/setup"
	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/application"
	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/handlers/messaging"
	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/infrastructure"
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

	if err := tripspb.Registration(reg); err != nil {
		return err
	}

	if err := driverspb.Registration(reg); err != nil {
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

	domainHandler := messaging.NewDomainEventHandlers(eventStream)
	messaging.RegisterDomainEventHandlers(domainDispatcher, domainHandler)

	integrationHandler := messaging.NewIntegrationEventHandlers(matchingSvc)
	messaging.RegisterIntegrationEventHandlers(eventStream, integrationHandler)

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
