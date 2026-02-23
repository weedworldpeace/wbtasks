package service

import (
	"context"
	"time"
	"worker/pkg/logger"
)

type RepoInterface interface {
	CleanupBookings() (int, error)
}

type ServiceConfig struct {
	ServiceInterval int `env:"SERVICE_INTERVAL" env-default:"60"` // seconds
}

type Service struct {
	repo   RepoInterface
	cfg    ServiceConfig
	ctx    context.Context
	stopCh chan struct{}
}

func New(repo RepoInterface, cfg ServiceConfig, ctx context.Context) *Service {
	return &Service{repo: repo, cfg: cfg, ctx: ctx, stopCh: make(chan struct{})}
}

func (s *Service) Start() {
	lg := logger.LoggerFromCtx(s.ctx).Lg

	t := time.NewTicker(time.Duration(s.cfg.ServiceInterval) * time.Second)
	for {
		select {
		case <-t.C:
			if n, err := s.repo.CleanupBookings(); err != nil {
				lg.Error().Err(err).Msg("failed to cleanup bookings")
			} else {
				lg.Info().Int("count", n).Msg("bookings cleaned up")
			}
		case <-s.stopCh:
			t.Stop()
			return
		}
	}
}

func (s *Service) Stop() {
	close(s.stopCh)
}
