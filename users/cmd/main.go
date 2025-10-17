package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Binit-Dhakal/Saarathi/pkg/logger"
	"github.com/Binit-Dhakal/Saarathi/pkg/rest/httpx"
	"github.com/Binit-Dhakal/Saarathi/pkg/rest/jsonutil"
	"github.com/Binit-Dhakal/Saarathi/pkg/setup"
	"github.com/Binit-Dhakal/Saarathi/pkg/waiter"
	"github.com/Binit-Dhakal/Saarathi/users/internal/application"
	"github.com/Binit-Dhakal/Saarathi/users/internal/handlers/grpc"
	"github.com/Binit-Dhakal/Saarathi/users/internal/handlers/rest"
	"github.com/Binit-Dhakal/Saarathi/users/internal/repository/postgres"
	"golang.org/x/sync/errgroup"
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

	app.rpc = setup.SetupRpc()

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
	detailsService := application.NewUserDetailService(userRepo)

	jsonReader := jsonutil.NewReader()
	jsonWriter := jsonutil.NewWriter()
	errorResponder := httpx.NewErrorResponder(jsonWriter, app.logger)

	authHandler := rest.NewUserHandler(authService, tokenService, jsonReader, jsonWriter, errorResponder)

	err = grpc.RegisterServer(detailsService, app.rpc)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/users/riders", authHandler.CreateRiderHandler)
	mux.HandleFunc("POST /api/v1/users/drivers", authHandler.CreateDriverHandler)
	mux.HandleFunc("POST /api/v1/tokens/authentication", authHandler.CreateTokenHandler)
	mux.HandleFunc("POST /api/v1/tokens/refresh", authHandler.RefreshTokenHandler)

	httpServer := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	runner := waiter.NewRunner(app.cfg.ShutdownTimeout)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	group, gCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return runner.RunHTTPServer(gCtx, httpServer)
	})

	group.Go(func() error {
		return runner.RunGRPCServer(gCtx, app.rpc, app.cfg.Rpc.Address())
	})

	app.logger.Info().Msg("Users service is starting up servers...")
	if err := group.Wait(); err != nil {
		app.logger.Error().Err(err).Msg("Service shut down due to an error")
		return err
	}

	app.logger.Info().Msg("All server shut down gracefully")

	return nil
}
