package main

import (
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type app struct {
	cfg         MatchAppConfig
	cacheClient *redis.Client
	nc          *nats.Conn
	js          nats.JetStreamContext
	logger      zerolog.Logger
}
