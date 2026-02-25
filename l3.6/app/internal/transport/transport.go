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
	CreateTransaction(context.Context, models.Transaction) (*models.Transaction, error)
	ListTransactions(context.Context, string, string, string, string) ([]models.Transaction, error)
	GetTransaction(context.Context, string) (*models.Transaction, error)
	UpdateTransaction(context.Context, models.Transaction) (*models.Transaction, error)
	DeleteTransaction(context.Context, string) error
	GetAnalytics(context.Context, string, string) (*models.Analytics, error)
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

	mux.Static("/static", "./web")
	mux.GET("/", func(c *ginext.Context) {
		c.File("./web/index.html")
	})

	transactions := mux.Group("/api/v1/transactions")
	transactions.Use(hers.middleware)

	transactions.POST("/", hers.CreateTransaction)
	transactions.GET("/", hers.ListTransactions)
	transactions.GET("/:id", hers.GetTransaction)
	transactions.PUT("/:id", hers.UpdateTransaction)
	transactions.DELETE("/:id", hers.DeleteTransaction)

	analytics := mux.Group("/api/v1/analytics")

	analytics.GET("/", hers.GetAnalytics)

	return &Server{&http.Server{Addr: fmt.Sprintf("%s:%d", serverCfg.Host, serverCfg.Port), Handler: mux}, ctx}
}
