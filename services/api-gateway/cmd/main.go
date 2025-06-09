package main

import (
	"github.com/SkySock/lode/services/api-gateway/internal/app"
)

func main() {
	cfg := app.MustLoad()
	app.Run(cfg)
}
