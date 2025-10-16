package main

import (
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/stackus/dotenv"
)

type (
	PGConfig struct {
		Conn string `required:"true"`
	}

	NatsConfig struct {
		URL        string `required:"true"`
		TripStream string `default:"saarathi"`
		SagaStream string `default:"trip-creation-saga"`
	}

	CacheConfig struct {
		CacheURL string `required:"true"`
	}

	OfferSvcConfig struct {
		Environment     string
		LogLevel        string `envconfig:"LOG_LEVEL" default:"DEBUG"`
		PG              PGConfig
		Nats            NatsConfig
		Redis           CacheConfig
		ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"30s"`
	}
)

func InitConfig() (cfg OfferSvcConfig, err error) {
	if err = dotenv.Load(dotenv.EnvironmentFiles(os.Getenv("ENVIRONMENT"))); err != nil {
		return
	}

	err = envconfig.Process("", &cfg)

	return
}
