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
	v1 "github.com/SkySock/lode/services/user-service/internal/handler/http/v1"
	"github.com/SkySock/lode/services/user-service/internal/repository"
	"github.com/SkySock/lode/services/user-service/internal/usecase"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func Run(cfg *config.Config) {
	log := setupLogger(cfg.Env)

	log.Info("starting application", slog.String("env", cfg.Env))

	pool := db.InitPool(context.Background(), cfg.DB, log)
	defer db.ClosePool(pool)

	accountRepository := repository.NewAccountRepository(pool)

	userUsecase := usecase.NewUserUsecase(accountRepository)

	controllers := controllers{
		SignUp: v1.NewSignUp(log, userUsecase),
		SignIn: v1.NewSignIn(log, userUsecase),
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
