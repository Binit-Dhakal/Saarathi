package main

import (
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/stackus/dotenv"
)

type (
	PGTripsConfig struct {
		Conn string `required:"true"`
	}

	NatsConfig struct {
		URL    string `required:"true"`
		Stream string `default:"saarathi"`
	}

	CacheConfig struct {
		CacheURL string `required:"true"`
	}

	TripAppConfig struct {
		Environment      string
		LogLevel         string `envconfig:"LOG_LEVEL" default:"DEBUG"`
		PG               PGTripsConfig
		Nats             NatsConfig
		UsersGrpcAddress string
		Redis            CacheConfig
		ShutdownTimeout  time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"30s"`
	}
)

func InitConfig() (cfg TripAppConfig, err error) {
	if err = dotenv.Load(dotenv.EnvironmentFiles(os.Getenv("ENVIRONMENT"))); err != nil {
		return
	}

	err = envconfig.Process("", &cfg)

	return
}
