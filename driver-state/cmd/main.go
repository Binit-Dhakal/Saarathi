package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Binit-Dhakal/Saarathi/driver-state/internal/application"
	"github.com/Binit-Dhakal/Saarathi/driver-state/internal/handlers/messaging"
	"github.com/Binit-Dhakal/Saarathi/driver-state/internal/handlers/ws"
	"github.com/Binit-Dhakal/Saarathi/driver-state/internal/repository/redis"
	"github.com/Binit-Dhakal/Saarathi/pkg/am"
	"github.com/Binit-Dhakal/Saarathi/pkg/contracts/proto/driverspb"
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
		fmt.Printf("Drivers-State service exitted abnormally:%v\n", err)
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
	var cfg DriverAppConfig
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

	err = driverspb.Registration(reg)
	if err != nil {
		return err
	}

	stream := jetstream.NewStream(cfg.Nats.Stream, app.js, app.logger)

	eventStream := am.NewEventStream(reg, stream)
	domainDispatcher := ddd.NewEventDispatcher[ddd.Event]()

	locationRepo := redis.NewLocationRepo(app.cacheClient)
	wsRepo := redis.NewWSRepo(app.cacheClient)
	offerRepo := redis.NewOfferRepository(app.cacheClient)

	presenceSvc := application.NewPresenceService(wsRepo)
	locationSvc := application.NewLocationService(locationRepo)

	// Solving circular dependency
	offerSvc := application.NewOfferService(domainDispatcher, nil, offerRepo)
	driverStateHandler := ws.NewWebSocketHandler(locationSvc, presenceSvc, offerSvc)
	offerSvc.SetNotifier(driverStateHandler)

	domainHandler := messaging.NewDomainHandlers(eventStream)
	messaging.RegisterDomainEventHandlers(domainDispatcher, domainHandler)

	integrationHandler := messaging.NewIntegrationEventHandlers(offerSvc)
	messaging.RegisterIntegrationHandlers(eventStream, integrationHandler)

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", driverStateHandler.WsHandler)

	fmt.Println("starting server on :8050")
	err = http.ListenAndServe(":8050", mux)
	if err != nil {
		return err
	}

	return nil
}
