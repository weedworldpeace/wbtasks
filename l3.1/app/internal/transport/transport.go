package transport

import (
	"app/internal/models"
	"app/pkg/logger"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/rabbitmq"
)

const (
	idParam = "id"
)

type ServiceInterface interface {
	CreateNotification(*models.Notification) (string, error)
	ReadNotification(string) (*models.Notification, error)
	DeleteNotification(string) error
}

type EmailSenderConfig struct {
	Host     string `env:"EMAIL_SENDER_HOST" env-default:"localhost"`
	Port     string `env:"EMAIL_SENDER_PORT" env-default:"1025"`
	Username string `env:"EMAIL_SENDER_USERNAME" env-default:"guest@guest.com"`
	Password string `env:"EMAIL_SENDER_PASSWORD" env-default:"guest"`
	From     string `env:"EMAIL_SENDER_FROM" env-default:"notification@service.com"`
}

type ServerConfig struct {
	Host        string `env:"SERVER_HOST" env-default:"localhost"`
	Port        string `env:"SERVER_PORT" env-default:"8080"`
	ReleaseMode string `env:"RELEASE_MODE" env-default:""`
}

type Server struct {
	httpServer *http.Server
	// cons       *rabbitmq.Consumer
	ctx context.Context
}

func (s *Server) Start() {
	lg := logger.LoggerFromCtx(s.ctx).Lg

	// lg.Info().Msg("trying starting da consumer")

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

	err = s.ctx.Value("rabbitConn").(*rabbitmq.Connection).Close()
	if err != nil {
		lg.Error().Err(err).Send()
	} else {
		lg.Info().Msg("rabbit connection closed gracefully")
	}

	lg.Info().Msg("all services stopped gracefully")
}

func New(service ServiceInterface, rabbitConn *rabbitmq.Connection, rabbitCh *rabbitmq.Channel, serverCfg *ServerConfig, emailSenderCfg *EmailSenderConfig, ctx context.Context) *Server {
	hers := &handlers{ctx, service, *emailSenderCfg}

	cons := rabbitmq.NewConsumer(rabbitCh, &rabbitmq.ConsumerConfig{Queue: "main_queue", AutoAck: true, Exclusive: false, NoLocal: false, NoWait: false, Args: nil})

	mux := ginext.New(serverCfg.ReleaseMode)

	mux.Use(hers.middleware)
	mux.POST("/notify", hers.createNotification)
	mux.GET(fmt.Sprintf("/notify/:%s", idParam), hers.readNotification)
	mux.DELETE(fmt.Sprintf("/notify/:%s", idParam), hers.deleteNotification)

	consumeCh := make(chan []byte)
	go func() {
		if err := cons.Consume(consumeCh); err != nil {
			panic(err)
		}
	}()

	go hers.consume(consumeCh)

	return &Server{&http.Server{Addr: fmt.Sprintf("%s:%s", serverCfg.Host, serverCfg.Port), Handler: mux}, context.WithValue(ctx, "rabbitConn", rabbitConn)}
}
