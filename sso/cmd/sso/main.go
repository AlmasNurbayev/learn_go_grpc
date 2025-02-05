package main

import (
	"fmt"
	"os"
	"os/signal"
	"sso/internal/config"
	"strconv"
	"syscall"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println("sso")
	fmt.Println(cfg)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		fmt.Println("go")
	}()
	fmt.Println("starting server on " + strconv.Itoa(cfg.GRPC.Port))
	<-done
	fmt.Println("stopping server on " + strconv.Itoa(cfg.GRPC.Port))
}
