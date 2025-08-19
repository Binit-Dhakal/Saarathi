package setup

import (
	"context"
	"fmt"
	"time"

	"github.com/Binit-Dhakal/Saarathi/pkg/env"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupPostgresDB() (*pgxpool.Pool, error) {
	dbHost := env.GetEnvWithDefault("DB_HOST", "localhost")
	dbPort := env.GetEnvWithDefault("DB_PORT", "5432")
	dbUser := env.GetEnvWithDefault("DB_USER", "postgres")
	dbPass := env.GetEnvWithDefault("DB_PASSWORD", "postgres")
	dbName := env.GetEnvWithDefault("DB_NAME", "saarathi")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName)

	// Create a connection pool
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("error parsing database config: %w", err)
	}
	// Set pool configuration
	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	// Create the connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	// Verify connection
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("error verifying database connection: %w", err)
	}

	return pool, nil
}
