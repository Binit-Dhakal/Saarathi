package main

import (
	"net/http"
	"os"
	"time"

	"github.com/Binit-Dhakal/Saarathi/pkg/env"
	log "github.com/Binit-Dhakal/Saarathi/pkg/logger"
	"github.com/Binit-Dhakal/Saarathi/pkg/rest/httpx"
	"github.com/Binit-Dhakal/Saarathi/pkg/rest/jsonutil"
	"github.com/Binit-Dhakal/Saarathi/pkg/setup"
	"github.com/Binit-Dhakal/Saarathi/users/internal/application"
	"github.com/Binit-Dhakal/Saarathi/users/internal/handlers/rest"
	"github.com/Binit-Dhakal/Saarathi/users/internal/repository/postgres"
)

func main() {
	logger := log.NewStandardLogger()

	dbpool, err := setup.SetupPostgresDB()
	if err != nil {
		logger.Error("failed to connect to the database", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	userRepo := postgres.NewUserRepo(dbpool)
	tokenRepo := postgres.NewTokenRepo(dbpool)

	jwtSecretKey, err := env.GetEnv("JWT_PRIVATE_KEY")
	if err != nil {
		panic(err)
	}

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
