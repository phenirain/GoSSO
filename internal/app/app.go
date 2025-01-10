package app

import (
	grpcapp "github.com/EtoNeAnanasbI95/auth-grpc-demo/internal/app/grpc"
	"github.com/EtoNeAnanasbI95/auth-grpc-demo/internal/config"
	"github.com/EtoNeAnanasbI95/auth-grpc-demo/internal/repository/user"
	"github.com/EtoNeAnanasbI95/auth-grpc-demo/internal/services/auth"
	"github.com/EtoNeAnanasbI95/auth-grpc-demo/internal/storage"
	"log/slog"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, config *config.Config) *App {
	const op = "App.New"
	log.Info(op, "msg", "attempting to init app")
	log.Info(op, "msg", "initializing db connection")
	db := storage.MustInitDb(config.ConnectionString)
	log.Info(op, "msg", "initializing app")
	usersRepository := user.New(db)
	log.Info(op, "msg", "initializing auth service")
	authService := auth.New(time.Second*60, log, usersRepository, []byte(config.Secret))

	grpcApp := grpcapp.New(log, config.GRPC.Port, authService)
	return &App{
		GRPCSrv: grpcApp,
	}
}
