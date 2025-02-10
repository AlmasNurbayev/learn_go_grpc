package app

import (
	"log/slog"
	grpcapp "sso/internal/app/grpc"
	"sso/internal/services/auth"
	"sso/internal/storage/postgres"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func NewApp(log *slog.Logger, port int, storagePath string, tokenTTL time.Duration) *App {
	// init storage

	postgresStorage, err := postgres.NewStorage(storagePath, log)
	if err != nil {
		panic(err)
	}
	// init auth service
	authService := auth.NewService(log, postgresStorage, tokenTTL)
	if err != nil {
		panic(err)
	}
	//
	grpcApp := grpcapp.NewApp(log, port, authService)

	return &App{
		GRPCServer: grpcApp}
}
