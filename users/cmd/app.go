package main

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type app struct {
	cfg     UserAppConfig
	usersDB *pgxpool.Pool
	rpc     *grpc.Server
	logger  zerolog.Logger
}
