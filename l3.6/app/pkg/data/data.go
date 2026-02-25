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
	DbPort     int    `env:"DB_PORT" env-default:"5433"`
	DbUser     string `env:"DB_USER" env-default:"user"`
	DbPassword string `env:"DB_PASSWORD" env-default:"12345"`
	DbName     string `env:"DB_NAME" env-default:"transactions"`
}

func New(cfg DataConfig) *Data {
	opts := &dbpg.Options{MaxOpenConns: 10, MaxIdleConns: 5}
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbName)
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
