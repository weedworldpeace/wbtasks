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
	CreateComment(ctx context.Context, req models.CreateCommentRequest) (*models.Comment, error)
	GetComments(ctx context.Context, req models.GetCommentsRequest) (*models.GetCommentsResponse, error)
	DeleteComment(ctx context.Context, id string) error
}

type ServerConfig struct {
	Host        string `env:"SERVER_HOST" env-default:"localhost"`
	Port        string `env:"SERVER_PORT" env-default:"8080"`
	ReleaseMode string `env:"RELEASE_MODE" env-default:""`
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
	hers := &handlers{ctx, service}

	mux := ginext.New(serverCfg.ReleaseMode)

	mux.Static("/static", "./web")

	mux.GET("/", func(c *ginext.Context) {
		c.File("././web/index.html")
	})

	api := mux.Group("/comments")

	api.Use(hers.middleware)
	api.POST("", hers.CreateComment)
	api.GET("", hers.GetComments)
	api.DELETE("/:id", hers.DeleteComment)

	return &Server{&http.Server{Addr: fmt.Sprintf("%s:%s", serverCfg.Host, serverCfg.Port), Handler: mux}, ctx}
}
