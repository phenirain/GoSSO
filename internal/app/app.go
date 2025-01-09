package app

import (
	grpcapp "github.com/EtoNeAnanasbI95/auth-grpc-demo/internal/app/grpc"
	"github.com/EtoNeAnanasbI95/auth-grpc-demo/internal/config"
	"github.com/EtoNeAnanasbI95/auth-grpc-demo/internal/services/auth"
	"log/slog"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, config *config.Config) *App {
	// TODO: init db

	// TODO: init service layer

	// TODO: put user repo into the argument, instead of nil
	authService := auth.New(time.Second*60, log, nil, []byte(config.Secret))

	grpcApp := grpcapp.New(log, config.GRPC.Port, authService)
	return &App{
		GRPCSrv: grpcApp,
	}
}
