package app

import (
	"log/slog"
	grpcapp "sso/internal/app/grpc"
	"sso/internal/services/auth"
	"sso/internal/storage/postgres"
	"time"
)

type App struct {
	GRPCServer       *grpcapp.App
	PostrgresStorage *postgres.Storage
}

func NewApp(log *slog.Logger, port int, DSN string, tokenTTL time.Duration,
	timeout time.Duration) *App {
	// init storage

	postgresStorage, err := postgres.NewStorage(DSN, log, timeout)
	if err != nil {
		panic(err)
	}
	// init auth service
	authService := auth.NewService(log, postgresStorage, tokenTTL)

	//
	grpcApp := grpcapp.NewApp(log, port, authService, timeout)

	return &App{
		GRPCServer:       grpcApp,
		PostrgresStorage: postgresStorage,
	}
}
