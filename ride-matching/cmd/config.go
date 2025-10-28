package main

import (
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/stackus/dotenv"
)

type (
	NatsConfig struct {
		URL    string `required:"true"`
		Stream string `default:"saarathi"`
	}

	CacheConfig struct {
		CacheURL string `required:"true"`
	}

	MatchAppConfig struct {
		Environment      string
		LogLevel         string `envconfig:"LOG_LEVEL" default:"DEBUG"`
		Nats             NatsConfig
		Redis            CacheConfig
		UsersGrpcAddress string
		ShutdownTimeout  time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"30s"`
	}
)

func InitConfig() (cfg MatchAppConfig, err error) {
	if err = dotenv.Load(dotenv.EnvironmentFiles(os.Getenv("ENVIRONMENT"))); err != nil {
		return
	}

	err = envconfig.Process("", &cfg)

	return
}
