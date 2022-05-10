package main

import (
	"GRPCService/api/grpc/server"
	"GRPCService/core"
	"GRPCService/logger"
	"context"
	"os"
	"os/signal"
)

const (
	GRPCPort = "8080" //Порт на котором поднимается GRPC сервер
)

func main() {
	logger.CreateLogger()

	err := core.CreateLocal()
	if err != nil {
		logger.LogrusLogger.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Kill, os.Interrupt)

	go func() {
		<-c
		logger.LogrusLogger.Info("Stopping service")
		cancel()
	}()

	err = server.StartGRPCServer(ctx, GRPCPort)
	if err != nil {
		logger.LogrusLogger.Fatal(err)
	}
}
