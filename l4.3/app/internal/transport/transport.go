package transport

import (
	"app/internal/models"
	"app/pkg/logger"
	"app/pkg/wrk/archiver"
	"app/pkg/wrk/notifyer"
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/wb-go/wbf/ginext"
)

type ServiceInterface interface {
	CreateEvent(context.Context, models.UserEvent, string) (*models.Event, error)
	UpdateEvent(context.Context, models.UserEvent) (*models.Event, error)
	DeleteEvent(context.Context, models.UserEvent) error
	ReadEvents(context.Context, string, string, string) ([]models.Event, error)
}

type ServerConfig struct {
	Host        string `env:"SERVER_HOST" env-default:"localhost"`
	Port        int    `env:"SERVER_PORT" env-default:"8080"`
	ReleaseMode string `env:"RELEASE_MODE" env-default:""`
}

type Server struct {
	httpServer *http.Server
	ctx        context.Context
	canc       context.CancelFunc
	wg         *sync.WaitGroup
	logCh      chan models.ToLog
}

func (s *Server) Start() {
	lg := logger.LoggerFromCtx(s.ctx).Lg

	lg.Info().Msg(fmt.Sprintf("trying starting da server on %s...", s.httpServer.Addr))

	if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		panic(err)
	}
}

func (s *Server) Stop() {
	lg := logger.LoggerFromCtx(s.ctx).Lg

	lg.Info().Msg("trying stopping da server...")

	ctx, canc := context.WithTimeout(context.Background(), time.Second*15)
	defer canc()
	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		lg.Error().Err(err).Send()
	} else {
		lg.Info().Msg("server stopped gracefully")
	}

	s.canc()
	close(s.logCh)
	lg.Info().Msg("waiting logs and workers(archive + notif)...")
	s.wg.Wait()
}

func New(service ServiceInterface, serverCfg *ServerConfig, wrk *archiver.Wrk, notif *notifyer.Notifyer, ctx context.Context) *Server {
	logCh := make(chan models.ToLog, logger.BufSize)
	notifCh := make(chan models.EventTask, notifyer.BufSize)
	hers := &handlers{ctx, service, logCh, notifCh}

	mux := ginext.New(serverCfg.ReleaseMode)
	mux.Use(cors.Default())

	mux.Use(hers.middleware)

	mux.POST("/create_event", hers.createEvent)
	mux.POST("/update_event", hers.updateEvent)
	mux.POST("/delete_event", hers.deleteEvent)
	mux.GET("/events_for_day", hers.eventsForDay)
	mux.GET("/events_for_week", hers.eventsForWeek)
	mux.GET("/events_for_month", hers.eventsForMonth)

	ctx, canc := context.WithCancel(ctx)

	var wg sync.WaitGroup
	wg.Add(3)

	go logger.LoggerFromCtx(ctx).Start(logCh, &wg)
	go wrk.Start(ctx, &wg)
	go notif.Start(ctx, &wg, notifCh)

	return &Server{&http.Server{Addr: fmt.Sprintf("%s:%d", serverCfg.Host, serverCfg.Port), Handler: mux}, ctx, canc, &wg, logCh}
}
