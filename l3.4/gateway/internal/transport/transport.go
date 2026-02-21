package transport

import (
	"app/internal/models"
	"app/pkg/logger"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/wb-go/wbf/ginext"
)

const (
	shorten = "shorten"
)

type ServiceInterface interface {
	UploadImage(*models.Image) (*models.UploadImageResponse, error)
	DeleteImage(string) error
	GetImage(string) (*models.GetImageResponse, error)
}

type ServerConfig struct {
	Host        string `env:"SERVER_HOST" env-default:"localhost"`
	Port        string `env:"SERVER_PORT" env-default:"8080"`
	ReleaseMode string `env:"RELEASE_MODE" env-default:""`
	LogLevel    int    `env:"LOG_LEVEL" env-default:"1"`
}

type Server struct {
	httpServer *http.Server
	ctx        context.Context
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

}

func New(service ServiceInterface, serverCfg *ServerConfig, ctx context.Context) *Server {
	// ctx = context.WithValue(ctx, logger.LoggerKey, logger.LoggerFromCtx(ctx).LoggerLevel(serverCfg.LogLevel))

	hers := &handlers{ctx, service}

	mux := ginext.New(serverCfg.ReleaseMode)

	mux.Static("/static", "./web")

	mux.GET("/", func(c *ginext.Context) {
		c.File("./web/index.html")
	})

	mux.Use(hers.middleware)
	mux.POST("/api/v1/upload", hers.UploadImage)
	mux.GET("/api/v1/image/:id", hers.GetImage)
	mux.DELETE("/api/v1/image/:id", hers.DeleteImage)

	return &Server{&http.Server{Addr: fmt.Sprintf("%s:%s", serverCfg.Host, serverCfg.Port), Handler: mux}, ctx}
}
