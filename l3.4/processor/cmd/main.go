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
	str := storage.New(cfg.StorageConfig)
	repo := repository.New(data, str)
	service := service.New(repo)
	brk := broker.New(cfg.BrokerConfig)
	cons := transport.New(service, &cfg.ConsumerConfig, ctx, brk)

	graceCh := make(chan os.Signal, 1)
	signal.Notify(graceCh, syscall.SIGINT, syscall.SIGTERM)

	cons.Start()

	<-graceCh
	cons.Stop()
}
