package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"worker/internal/config"
	"worker/internal/repository"
	"worker/internal/service"
	"worker/pkg/data"
	"worker/pkg/logger"
)

func main() {
	lg := logger.New()
	ctx := context.WithValue(context.Background(), logger.LoggerKey, lg)

	cfg := config.New()
	data := data.New(cfg.DataConfig)
	repo := repository.New(data)
	srv := service.New(repo, cfg.ServiceConfig, ctx)

	graceCh := make(chan os.Signal, 1)
	signal.Notify(graceCh, syscall.SIGINT, syscall.SIGTERM)

	go srv.Start()

	<-graceCh
	srv.Stop()
}
