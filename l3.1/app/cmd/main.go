package main

import (
	"app/internal/config"
	"app/internal/repository"
	"app/internal/service"
	"app/internal/transport"
	"app/pkg/data/inmem"
	"app/pkg/data/rabbit"
	"app/pkg/logger"
	"context"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	lg := logger.New()
	ctx := context.WithValue(context.Background(), logger.LoggerKey, lg)

	cfg := config.New()
	rabbitConn, rabbitCh := rabbit.New(cfg.RabbitConfig)
	data := inmem.New()
	repo := repository.New(data, rabbitCh)
	service := service.New(repo)
	server := transport.New(service, rabbitConn, rabbitCh, &cfg.ServerConfig, &cfg.EmailSenderConfig, ctx)

	graceCh := make(chan os.Signal, 1)
	signal.Notify(graceCh, syscall.SIGINT, syscall.SIGTERM)

	go server.Start()

	<-graceCh
	server.Stop()
}
