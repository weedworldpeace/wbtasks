package archiver

import (
	"app/internal/models"
	"app/pkg/logger"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/wb-go/wbf/dbpg"
)

type Wrk struct {
	db       *dbpg.DB
	interval int
}

type WrkConfig struct {
	DbHost     string `env:"DB_HOST" env-default:"localhost"`
	DbPort     int    `env:"DB_PORT" env-default:"5433"`
	DbUser     string `env:"DB_USER" env-default:"user"`
	DbPassword string `env:"DB_PASSWORD" env-default:"12345"`
	DbName     string `env:"DB_NAME" env-default:"events"`
	Interval   int    `env:"WRK_INTERVAL" env-default:"10"`
}

func (w *Wrk) archive(ctx context.Context) {
	lg := logger.LoggerFromCtx(ctx).Lg

	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		lg.Error().Str("worker", "archiver").Err(errors.Join(models.ErrOnDatabase, err)).Send()
		return
	}
	defer tx.Rollback()

	q1 := `DELETE FROM events WHERE date < NOW() RETURNING event_id, user_id, message, date, created_at, updated_at`

	events := make([]models.UserEvent, 0)
	rows, err := tx.QueryContext(ctx, q1)
	if err != nil {
		lg.Error().Str("worker", "archiver").Err(errors.Join(models.ErrOnDatabase, err)).Send()
		return
	}

	for rows.Next() {
		var ev models.UserEvent
		err := rows.Scan(&ev.EventId, &ev.UserId, &ev.Message, &ev.Date, &ev.CreatedAt, &ev.UpdatedAt)
		if err != nil {
			lg.Error().Str("worker", "archiver").Err(errors.Join(models.ErrOnDatabase, err)).Send()
			return
		}
		events = append(events, ev)
	}

	if rows.Err() != nil {
		lg.Error().Str("worker", "archiver").Err(errors.Join(models.ErrOnDatabase, err)).Send()
		return
	}

	q2 := `INSERT INTO archive (user_id, event_id, message, date, created_at, updated_at) values ($1, $2, $3, $4, $5, $6)`
	for _, ev := range events {
		if _, err := tx.ExecContext(ctx, q2, ev.UserId, ev.EventId, ev.Message, ev.Date, ev.CreatedAt, ev.UpdatedAt); err != nil {
			lg.Error().Str("worker", "archiver").Err(errors.Join(models.ErrOnDatabase, err)).Send()
			return
		}
	}

	if err := tx.Commit(); err != nil {
		lg.Error().Str("worker", "archiver").Err(errors.Join(models.ErrOnDatabase, err)).Msg("while commiting")
	} else {
		lg.Info().Str("worker", "archiver").Int("quantity", len(events)).Msg("events archived successfully")
	}
}

func (w *Wrk) Start(ctx context.Context, wg *sync.WaitGroup) {
	t := time.NewTicker(time.Duration(w.interval) * time.Second)
	for {
		select {
		case <-t.C:
			w.archive(ctx)
		case <-ctx.Done():
			wg.Done()
			return
		}
	}
}

func New(cfg WrkConfig) *Wrk {
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
	return &Wrk{
		db:       db,
		interval: cfg.Interval,
	}
}
