package setup

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupPostgresDB(conn string) (*pgxpool.Pool, error) {
	// Create a connection pool
	config, err := pgxpool.ParseConfig(conn)
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
