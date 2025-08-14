package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	log "github.com/Binit-Dhakal/Saarathi/pkg/logger"
	"github.com/Binit-Dhakal/Saarathi/pkg/rest/httpx"
	"github.com/Binit-Dhakal/Saarathi/pkg/rest/jsonutil"
	"github.com/Binit-Dhakal/Saarathi/users/internal/application"
	"github.com/Binit-Dhakal/Saarathi/users/internal/handlers/rest"
	"github.com/Binit-Dhakal/Saarathi/users/internal/repository/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	logger := log.NewStandardLogger()

	dbpool, err := setupDatabase()
	if err != nil {
		logger.Error("failed to connect to the database", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	userRepo := postgres.NewUserRepo(dbpool)
	tokenRepo := postgres.NewTokenRepo(dbpool)

	jwtSecretKey := getEnvWithDefault("JWT_PRIVATE_KEY", "")
	// key, err := getPrivateKey(jwtSecretKey)
	// if err != nil {
	// 	panic(err)
	// }

	authService := application.NewAuthService(dbpool, userRepo)
	tokenService := application.NewJWTService(jwtSecretKey, tokenRepo)

	jsonReader := jsonutil.NewReader()
	jsonWriter := jsonutil.NewWriter()
	errorResponder := httpx.NewErrorResponder(jsonWriter, logger)

	authHandler := rest.NewUserHandler(authService, tokenService, jsonReader, jsonWriter, errorResponder)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/users/riders", authHandler.CreateRiderHandler)
	mux.HandleFunc("POST /api/v1/users/drivers", authHandler.CreateDriverHandler)
	mux.HandleFunc("POST /api/v1/tokens/authentication", authHandler.CreateTokenHandler)
	mux.HandleFunc("POST /api/v1/tokens/refresh", authHandler.RefreshTokenHandler)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	logger.Info("starting serve on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("server failed", err)
		os.Exit(1)
	}
}

func setupDatabase() (*pgxpool.Pool, error) {
	dbHost := getEnvWithDefault("DB_HOST", "localhost")
	dbPort := getEnvWithDefault("DB_PORT", "5432")
	dbUser := getEnvWithDefault("DB_USER", "postgres")
	dbPass := getEnvWithDefault("DB_PASSWORD", "postgres")
	dbName := getEnvWithDefault("DB_NAME", "saarathi")

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

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
