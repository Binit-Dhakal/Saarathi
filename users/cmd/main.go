package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Binit-Dhakal/Saarathi/pkg/logger"
	"github.com/Binit-Dhakal/Saarathi/pkg/rest/httpx"
	"github.com/Binit-Dhakal/Saarathi/pkg/rest/jsonutil"
	"github.com/Binit-Dhakal/Saarathi/pkg/setup"
	"github.com/Binit-Dhakal/Saarathi/users/internal/application"
	"github.com/Binit-Dhakal/Saarathi/users/internal/handlers/rest"
	"github.com/Binit-Dhakal/Saarathi/users/internal/repository/postgres"
)

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Trips service exitted abnormally: %v\n", err)
		os.Exit(1)
	}
}

func infraSetup(app *app) (err error) {
	app.usersDB, err = setup.SetupPostgresDB(app.cfg.PG.Conn)
	if err != nil {
		return err
	}

	app.logger = logger.New(logger.LogConfig{
		Environment: app.cfg.Environment,
		LogLevel:    logger.Level(app.cfg.LogLevel),
	})

	return nil
}

func run() (err error) {
	var cfg UserAppConfig
	cfg, err = InitConfig()
	if err != nil {
		return err
	}

	app := &app{
		cfg: cfg,
	}

	err = infraSetup(app)
	if err != nil {
		return
	}
	defer app.usersDB.Close()

	userRepo := postgres.NewUserRepo(app.usersDB)
	tokenRepo := postgres.NewTokenRepo(app.usersDB)

	authService := application.NewAuthService(app.usersDB, userRepo)
	tokenService := application.NewJWTService(app.cfg.PrivateKey, tokenRepo)

	jsonReader := jsonutil.NewReader()
	jsonWriter := jsonutil.NewWriter()
	errorResponder := httpx.NewErrorResponder(jsonWriter, app.logger)

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

	fmt.Println("starting serve on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return
}
