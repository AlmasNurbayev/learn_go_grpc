package grpcapp

import (
	"fmt"
	"log/slog"
	"net"
	authgrpc "sso/internal/grpc/auth"
	"sso/internal/grpc/middleware"
	"sso/internal/services/auth"
	"time"

	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
	timeout    time.Duration
}

func NewApp(
	log *slog.Logger,
	port int,
	authService *auth.AuthService,
	timeout time.Duration,

) *App {
	gRPCServer := grpc.NewServer(grpc.UnaryInterceptor(middleware.ContextMiddleware(timeout, log)))
	authgrpc.Register(gRPCServer, authService, log)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
		timeout:    timeout,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(slog.String("op", op), slog.Int("port", a.port))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		log.Error("failed to start grpc server", op, err)
	}

	log.Info("starting grpc is running", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		log.Error("failed to start grpc server", op, err)
	}
	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).Warn("stopping grpc server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()

}
