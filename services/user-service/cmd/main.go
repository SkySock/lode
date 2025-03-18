package main

import (
	"github.com/SkySock/lode/services/user-service/internal/app"
	"github.com/SkySock/lode/services/user-service/internal/config"
	"github.com/SkySock/lode/services/user-service/internal/db"
)

func main() {
	cfg := config.MustLoad()

	if cfg.DB.Migrate {
		if err := db.RunMigrations(cfg.DB.URL); err != nil {
			panic(err)
		}
	}

	app.Run(cfg)
}
