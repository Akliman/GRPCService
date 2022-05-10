package main

import (
	"GRPCService/api/grpc/server"
	"GRPCService/core"
	"GRPCService/external/memcache"
	"GRPCService/logger"
	"context"
	"os"
	"os/signal"
)

const (
	memcacheAdr = "localhost:11211" //Адрес memcached сервера
	GRPCPort    = "8080"            //Порт на котором поднимается GRPC сервер
)

func main() {

	logger.CreateLogger()

	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Kill, os.Interrupt)

	go func() {
		<-c
		logger.LogrusLogger.Info("Stopping service")
		cancel()
		memcache.Close()
	}()

	err := core.CreateMemcahed(ctx, memcacheAdr)
	if err != nil {
		logger.LogrusLogger.Fatal("error memcached connection. err: ", err)
	}

	err = server.StartGRPCServer(ctx, GRPCPort)
	if err != nil {
		logger.LogrusLogger.Fatal("error starting server. err: ", err)
	}
}
