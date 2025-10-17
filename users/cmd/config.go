package main

import (
	"fmt"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/stackus/dotenv"
)

type (
	RpcConfig struct {
		Host string `default:"0.0.0.0"`
		Port string `default:":8085"`
	}

	PGUsersConfig struct {
		Conn string `required:"true"`
	}

	UserAppConfig struct {
		Environment     string
		LogLevel        string `envconfig:"LOG_LEVEL" default:"DEBUG"`
		PG              PGUsersConfig
		Rpc             RpcConfig
		ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"30s"`
		PrivateKey      string        `envconfig:"JWT_PRIVATE_KEY" required:"true"`
	}
)

func InitConfig() (cfg UserAppConfig, err error) {
	if err = dotenv.Load(dotenv.EnvironmentFiles(os.Getenv("ENVIRONMENT"))); err != nil {
		return
	}

	err = envconfig.Process("", &cfg)

	return
}

func (c RpcConfig) Address() string {
	return fmt.Sprintf("%s%s", c.Host, c.Port)
}
