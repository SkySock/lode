package db

import (
	"context"
	"log/slog"
	"os"

	"github.com/SkySock/lode/services/user-service/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitPool(ctx context.Context, cfg config.DBConfig, log *slog.Logger) *pgxpool.Pool {
	connString := cfg.URL

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		panic(err)
	}
	config.MaxConns = 10
	config.MinConns = 1

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Error("Failed to create connection pool", slog.String("error", err.Error()))
		os.Exit(1)
	}

	err = pool.Ping(ctx)
	if err != nil {
		log.Error("Failed to ping database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	return pool
}

func ClosePool(pool *pgxpool.Pool) {
	if pool != nil {
		pool.Close()
	}
}
