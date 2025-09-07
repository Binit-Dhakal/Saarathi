package main

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type app struct {
	cfg     UserAppConfig
	usersDB *pgxpool.Pool
	logger  zerolog.Logger
}
