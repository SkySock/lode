package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SkySock/lode/services/user-service/internal/config"
	"github.com/SkySock/lode/services/user-service/internal/db"
	v1Auth "github.com/SkySock/lode/services/user-service/internal/handler/http/v1/auth"
	"github.com/SkySock/lode/services/user-service/internal/repository"
	"github.com/SkySock/lode/services/user-service/internal/usecase"
	"github.com/valkey-io/valkey-go"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func Run(cfg *config.Config) {
	log := setupLogger(cfg.Env)

	log.Info("starting application", slog.String("env", cfg.Env))

	ctx := context.Background()
	pool := db.InitPool(ctx, cfg.DB, log)
	defer db.ClosePool(pool)

	client, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{cfg.Valkey.Addr}})
	if err != nil {
		panic(err)
	}
	defer client.Close()

	sessionRepo := repository.NewSessionRepository(client)
	accountRepo := repository.NewAccountRepository()
	profileRepo := repository.NewProfileRepository()

	authUsecase := usecase.NewAuthUsecase(pool, accountRepo, sessionRepo, profileRepo, cfg.Auth)

	controllers := controllers{
		SignUp:  v1Auth.NewSignUp(log, authUsecase),
		SignIn:  v1Auth.NewSignIn(log, authUsecase),
		SignOut: v1Auth.NewSignOut(log, authUsecase),
	}

	r := newRouter(log, controllers)

	s := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port),
		Handler:      r,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	serverErrors := make(chan error, 1)

	go func() {
		log.Info("Starting HTTP server", slog.String("addr", s.Addr))
		err := s.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	select {
	case err := <-serverErrors:
		log.Error("Server error", slog.String("error", err.Error()))
		os.Exit(1)

	case sig := <-stop:
		log.Info("Received signal:", slog.String("signal", sig.String()))
		log.Info("Shutting down gracefully...")

		tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := s.Shutdown(tc); err != nil {
			log.Error("Graceful shutdown failed", slog.String("error", err.Error()))
		} else {
			log.Info("Server stopped.")
		}
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
