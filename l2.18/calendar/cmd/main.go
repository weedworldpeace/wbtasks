package main

import (
	"calendar/internal/config"
	"calendar/internal/repository"
	"calendar/internal/service"
	"calendar/internal/transport"
	"calendar/pkg/data"
	"calendar/pkg/logger"
	"context"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	lg := logger.New()
	ctx := context.WithValue(context.Background(), logger.LoggerKey, lg)

	cfg := config.New()

	data := data.New()
	repo := repository.New(data)
	service := service.New(repo)
	server := transport.New(service, &cfg.ServerConfig, ctx)

	graceCh := make(chan os.Signal, 1)
	signal.Notify(graceCh, syscall.SIGINT, syscall.SIGTERM)

	go server.Start()

	<-graceCh
	server.Stop()
}
