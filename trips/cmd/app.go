package main

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type app struct {
	cfg         TripAppConfig
	tripsDB     *pgxpool.Pool
	cacheClient *redis.Client
	nc          *nats.Conn
	js          nats.JetStreamContext
	logger      zerolog.Logger
}
