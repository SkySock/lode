package main

import (
	"github.com/SkySock/lode/services/sso/internal/app"
	"github.com/SkySock/lode/services/sso/internal/config"
)

func main() {
	cfg := config.MustLoad()
	app.Run(cfg)
}
