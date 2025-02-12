package main

import (
	"context"
	"os"
	"os/signal"
	"sso/internal/app"
	"sso/internal/config"
	"sso/internal/lib/logger"
	"sso/internal/utils"
	"syscall"
	"time"
)

func main() {
	cfg := config.MustLoad()
	Log := logger.InitLogger(cfg.Env)
	p, err := utils.PrintAsJSON(cfg)
	if err != nil {
		panic(err)
	}
	Log.Info("load config: ")
	Log.Info(string(*p))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	application := app.NewApp(ctx, Log, cfg.GRPC.Port, cfg.DSN, cfg.TokenTTL)

	go func() {
		application.GRPCServer.MustRun()
	}()
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	signalString := <-done
	Log.Info("received signal " + signalString.String())

	application.GRPCServer.Stop()
	application.PostrgresStorage.Close()

	Log.Info("server stopped")

}
