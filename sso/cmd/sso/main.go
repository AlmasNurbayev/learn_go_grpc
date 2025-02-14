package main

import (
	"fmt"
	"os"
	"os/signal"
	"sso/internal/app"
	"sso/internal/config"
	"sso/internal/lib/logger"
	"sso/internal/utils"
	"syscall"
)

func main() {
	s := fmt.Sprintf("%s", "фывфыв")
	fmt.Println(s)

	cfg := config.MustLoad()
	Log := logger.InitLogger(cfg.Env)
	p, err := utils.PrintAsJSON(cfg)
	if err != nil {
		panic(err)
	}
	Log.Info("load config: ")
	Log.Info(string(*p))

	application := app.NewApp(Log, cfg.GRPC.Port, cfg.DSN, cfg.TokenTTL, cfg.GRPC.Timeout)

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
