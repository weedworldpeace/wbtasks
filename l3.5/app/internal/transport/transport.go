package transport

import (
	"app/internal/models"
	"app/pkg/logger"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/wb-go/wbf/ginext"
)

const (
	shorten = "shorten"
)

type ServiceInterface interface {
	CreateEvent(context.Context, models.Event) (*models.CreateResponse, error)
	BookEvent(context.Context, models.BookRequest) (*models.BookResponse, error)
	ConfirmEvent(context.Context, models.ConfirmRequest) error
	GetEvent(context.Context, string) (*models.EventResponse, error)
	ListEvents(context.Context) ([]models.EventResponse, error)
}

type ServerConfig struct {
	Host        string `env:"SERVER_HOST" env-default:"localhost"`
	Port        int    `env:"SERVER_PORT" env-default:"8080"`
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
	mux.Use(cors.Default())
	mux.Use(hers.middleware)

	mux.Static("/admin/static", "./web/admin")
	mux.GET("/admin", func(c *ginext.Context) {
		c.File("./web/admin/index.html")
	})
	mux.Static("/user/static", "./web/user")
	mux.GET("/user", func(c *ginext.Context) {
		c.File("./web/user/index.html")
	})

	mux.GET("/", func(c *ginext.Context) {
		c.Redirect(http.StatusFound, "/user")
	})

	events := mux.Group("/api/v1/events")

	events.POST("/", hers.CreateEvent)
	events.POST("/:id/book", hers.BookEvent)
	events.POST("/:id/confirm", hers.ConfirmEvent)
	events.GET("/:id", hers.GetEvent)
	events.GET("/list", hers.ListEvents)

	return &Server{&http.Server{Addr: fmt.Sprintf("%s:%d", serverCfg.Host, serverCfg.Port), Handler: mux}, ctx}
}
