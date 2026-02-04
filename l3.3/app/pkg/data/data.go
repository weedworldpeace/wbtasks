package data

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/wb-go/wbf/dbpg"
)

type Data struct {
	Mu sync.Mutex
	DB *dbpg.DB
}

type DataConfig struct {
	DbHost     string `env:"DB_HOST" env-default:"localhost"`
	DbPort     string `env:"DB_PORT" env-default:"5433"`
	DbUser     string `env:"DB_USER" env-default:"shortener"`
	DbPassword string `env:"DB_PASSWORD" env-default:"shortener_pass"`
	DbName     string `env:"DB_NAME" env-default:"url_shortener"`
}

func New(cfg DataConfig) *Data {
	opts := &dbpg.Options{MaxOpenConns: 10, MaxIdleConns: 5}
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbName)
	fmt.Println(dsn)
	db, err := dbpg.New(dsn, []string{}, opts)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to db: %v", err))
	}

	ctx, canc := context.WithTimeout(context.Background(), 5*time.Second)
	defer canc()
	err = db.Master.PingContext(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed to ping db: %v", err))
	}
	return &Data{
		DB: db,
	}

}
