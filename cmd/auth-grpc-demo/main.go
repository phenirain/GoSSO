package main

import (
	"github.com/EtoNeAnanasbI95/auth-grpc-demo/internal/app"
	"github.com/EtoNeAnanasbI95/auth-grpc-demo/internal/config"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	//os.Setenv("GRPC_AUTH_CONFIG_PATH", "./config/local.yaml")
	cfg := config.MustLoadConfig()

	log := setupLogger(cfg.Env)

	log.Info("Starting GRPC api",
		slog.String("env", cfg.Env),
		slog.Int("port", cfg.GRPC.Port))

	application := app.New(log, cfg)
	go application.GRPCSrv.MustRun()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	log.Info("Shutting down server...")
	application.GRPCSrv.Stop()

}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envLocal:
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return logger
}
