package main

import (
	"app/internal/config"
	"app/internal/repository"
	"app/internal/service"
	"app/internal/transport"
	"app/pkg/broker"
	"app/pkg/data"
	"app/pkg/logger"
	"app/pkg/storage"
	"context"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	lg := logger.New()
	ctx := context.WithValue(context.Background(), logger.LoggerKey, lg)

	cfg := config.New()
	data := data.New(cfg.DataConfig)
	brk := broker.New(cfg.BrokerConfig)
	str := storage.New(cfg.StorageConfig)
	repo := repository.New(data, brk, str)
	service := service.New(repo)
	server := transport.New(service, &cfg.ServerConfig, ctx)

	graceCh := make(chan os.Signal, 1)
	signal.Notify(graceCh, syscall.SIGINT, syscall.SIGTERM)

	go server.Start()

	<-graceCh
	server.Stop()
}
