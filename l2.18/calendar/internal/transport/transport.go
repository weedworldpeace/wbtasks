package transport

import (
	"calendar/internal/models"
	"calendar/pkg/logger"
	"context"
	"fmt"
	"net/http"
	"time"
)

type ServiceInterface interface {
	CreateEvent(*models.UserEvent) error
	UpdateEvent(*models.UserEvent) error
	DeleteEvent(*models.UserEvent) error
	ReadEvents(string, string, string) ([]models.Event, error)
}

type ServerConfig struct {
	Host string `env:"SERVER_HOST" env-default:"localhost"`
	Port string `env:"SERVER_PORT" env-default:"8080"`
}

type Server struct {
	httpServer *http.Server
	ctx        context.Context
}

func (s *Server) Start() {
	lg := logger.LoggerFromCtx(s.ctx).Lg

	lg.Info(fmt.Sprintf("trying starting da server on %s...", s.httpServer.Addr))

	err := s.httpServer.ListenAndServe()
	if err != http.ErrServerClosed {
		panic(err)
	}
}

func (s *Server) Stop() {
	lg := logger.LoggerFromCtx(s.ctx).Lg

	lg.Info("trying stopping da server...")

	ctx, canc := context.WithTimeout(context.Background(), time.Second*15)
	defer canc()
	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		lg.Error(err.Error())
	} else {
		lg.Info("server stopped gracefully")
	}
}

func New(service ServiceInterface, cfg *ServerConfig, ctx context.Context) *Server {
	hers := &handlers{ctx, service}
	mux := http.NewServeMux()

	mux.HandleFunc("/create_event", hers.middleware(http.HandlerFunc(hers.createEvent)))
	mux.HandleFunc("/update_event", hers.middleware(http.HandlerFunc(hers.updateEvent)))
	mux.HandleFunc("/delete_event", hers.middleware(http.HandlerFunc(hers.deleteEvent)))
	mux.HandleFunc("/events_for_day", hers.middleware(http.HandlerFunc(hers.eventsForDay)))
	mux.HandleFunc("/events_for_week", hers.middleware(http.HandlerFunc(hers.eventsForWeek)))
	mux.HandleFunc("/events_for_month", hers.middleware(http.HandlerFunc(hers.eventsForMonth)))

	return &Server{&http.Server{Addr: fmt.Sprintf("%s:%s", cfg.Host, cfg.Port), Handler: mux}, ctx}
}
