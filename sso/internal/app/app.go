package app

import (
	"context"
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

func NewApp(ctx context.Context, log *slog.Logger, port int, storagePath string, tokenTTL time.Duration) *App {
	// init storage

	postgresStorage, err := postgres.NewStorage(ctx, storagePath, log)
	if err != nil {
		panic(err)
	}
	// init auth service
	authService := auth.NewService(log, postgresStorage, tokenTTL)

	//
	grpcApp := grpcapp.NewApp(log, port, authService)

	return &App{
		GRPCServer:       grpcApp,
		PostrgresStorage: postgresStorage,
	}
}
