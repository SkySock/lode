package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SkySock/lode/libs/utils/http/middleware"
	"github.com/gorilla/mux"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func reverseProxyWithLogger(l *slog.Logger) func(target string) http.Handler {
	return func(target string) http.Handler {
		targetURL, err := url.Parse(target)
		if err != nil {
			l.Error("Failed to parse target URL", "target", target, "error", err)
		}
		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		return proxy
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

func Run(cfg *Config) {
	log := setupLogger(cfg.Env)
	log.Info("starting application", slog.String("env", cfg.Env))

	newReverseProxy := reverseProxyWithLogger(log)

	r := mux.NewRouter()

	r.Use(middleware.Logging(log))
	r.PathPrefix("/api/v1/auth").Handler(newReverseProxy("http://user-service:8080"))

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
