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

type ServiceInterface interface {
	CreateItem(context.Context, models.Item) (*models.Item, error)
	ListItems(context.Context, string, string) ([]models.Item, error)
	GetItem(context.Context, string) (*models.Item, error)
	UpdateItem(context.Context, models.Item) (*models.Item, error)
	DeleteItem(context.Context, string) error
	ListHistory(context.Context, string, string) ([]models.ItemHistory, error)
	GetToken(context.Context, string) (string, error)
	// GetAnalytics(context.Context, string, string) (*models.Analytics, error)
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
	// config := cors.Config{
	// 	AllowOrigins:     []string{"http://localhost:8080", "http://127.0.0.1:8080"},
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// 	MaxAge:           12 * time.Hour,
	// }

	// mux.Use(cors.New(config))
	mux.Use(cors.Default())

	mux.Static("/auth/static", "./web/auth")
	mux.GET("/auth", func(c *ginext.Context) {
		c.File("./web/auth/index.html")
	})
	mux.Static("/home/static", "./web/home")
	mux.GET("/home", func(c *ginext.Context) {
		c.File("./web/home/index.html")
	})

	mux.GET("/", func(c *ginext.Context) {
		c.File("./web/home/index.html")
	})

	api := mux.Group("/api/v1")
	api.Use(hers.middleware)

	api.GET("/auth", hers.getToken)

	history := api.Group("/history")
	history.Use(hers.jwtChecker)
	history.GET("/", hers.listHistory)

	items := api.Group("/items")
	items.Use(hers.jwtChecker)
	items.POST("/", hers.createItem)
	items.GET("/", hers.listItems)
	items.GET("/:id", hers.getItem)
	items.PUT("/:id", hers.updateItem)
	items.DELETE("/:id", hers.deleteItem)

	return &Server{&http.Server{Addr: fmt.Sprintf("%s:%d", serverCfg.Host, serverCfg.Port), Handler: mux}, ctx}
}
