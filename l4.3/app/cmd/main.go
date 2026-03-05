package main

import (
	"app/internal/config"
	"app/internal/repository"
	"app/internal/service"
	"app/internal/transport"
	"app/pkg/data"
	"app/pkg/logger"
	"app/pkg/sender"
	"app/pkg/wrk/archiver"
	"app/pkg/wrk/notifyer"
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

	repo := repository.New(data)
	service := service.New(repo)

	archiver := archiver.New(cfg.WrkConfig)
	sender := sender.New(cfg.SndConfig)
	notif := notifyer.New(sender)

	server := transport.New(service, &cfg.ServerConfig, archiver, notif, ctx)

	graceCh := make(chan os.Signal, 1)
	signal.Notify(graceCh, syscall.SIGINT, syscall.SIGTERM)

	go server.Start()

	<-graceCh
	server.Stop()
}
