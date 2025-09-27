package main

import (
	"log/slog"
	"os"

	"github.com/phenirain/sso/internal"
	"github.com/phenirain/sso/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Could not load config", "err", err)
		os.Exit(1)
	}
	if err := internal.Run(cfg); err != nil {
		slog.Error("Failed to run server", "err", err)
		os.Exit(1)
	}
}
