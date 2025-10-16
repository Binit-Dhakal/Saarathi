package main

import (
	"fmt"
	"os"

	"github.com/Binit-Dhakal/Saarathi/offers/internal/handlers/messaging"
	"github.com/Binit-Dhakal/Saarathi/pkg/am"
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
	app.cacheClient, err = setup.SetupRedis(app.cfg.Redis.CacheURL)
	if err != nil {
		return err
	}

	app.nc, err = nats.Connect(app.cfg.Nats.URL)
	if err != nil {
		return err
	}

	app.js, err = setup.SetupJetStream(app.nc)
	if err != nil {
		return err
	}

	err = setup.SetupStreams(app.js, app.cfg.Nats.TripStream, app.cfg.Nats.SagaStream)
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
	defer app.cacheClient.Close()
	defer app.nc.Close()

	reg := registry.NewRegistry()

	err = tripspb.Registration(reg)
	if err != nil {
		return err
	}

	sagaStream := jetstream.NewStream(cfg.Nats.SagaStream, app.js, app.logger)
	domainDispatcher := ddd.NewEventDispatcher[ddd.Event]()

	sagaEvtStream := am.NewEventStream(reg, sagaStream)

	integrationHandlers := messaging.NewIntegrationEventHandlers(domainDispatcher)
	err = messaging.RegisterIntegrationHandlers(sagaEvtStream, integrationHandlers)
	if err != nil {
		return err
	}

	return nil
}
