package main

import (
	"app/internal/config"
	"app/internal/service"
	"app/pkg/logger"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	lg := logger.New()
	ctx, canc := context.WithTimeout(context.WithValue(context.Background(), logger.LoggerKey, lg), time.Second*30)

	cfg := config.New()
	doneCh := make(chan struct{})
	service := service.New(cfg.ServiceConfig, ctx, doneCh)
	go service.Cut()

	graceCh := make(chan os.Signal, 1)
	signal.Notify(graceCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-graceCh:
		canc()
		time.Sleep(2 * time.Second)
	case <-doneCh:
		canc()
	}
}
