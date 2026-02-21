package repository

import (
	"app/internal/models"
	"app/pkg/broker"
	"app/pkg/data"
	"app/pkg/storage"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/wb-go/wbf/retry"
)

type Repository struct {
	data *data.Data
	brk  *broker.Broker
	str  *storage.StorageClient
}

func New(data *data.Data, broker *broker.Broker, str *storage.StorageClient) *Repository {
	return &Repository{data: data, brk: broker, str: str}
}

func (r *Repository) UploadImage(task *models.KafkaTask) (*models.UploadImageResponse, error) {
	ctx, canc := context.WithTimeout(context.Background(), 20*time.Second)
	defer canc()

	tx, err := r.data.DB.BeginTxWithRetry(ctx, retry.Strategy{Attempts: 3, Delay: 2 * time.Second, Backoff: 1.5}, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if _, err = tx.ExecContext(ctx, `INSERT INTO tasks (id, status, path) VALUES ($1, 'pending', '')`, task.ID); err != nil {
		return nil, err
	}

	data, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}

	err = r.brk.Producer.SendWithRetry(ctx, retry.Strategy{Attempts: 3, Delay: 2 * time.Second, Backoff: 1.5}, []byte("tasks"), data)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &models.UploadImageResponse{ID: task.ID}, nil
}

func (r *Repository) GetImage(id string) (*models.GetImageResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var status, path string
	if err := r.data.DB.QueryRowContext(ctx, `SELECT status, path FROM tasks WHERE id = $1`, id).Scan(&status, &path); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrImageNotFound
		}
		return nil, err
	}

	switch status {
	case models.StatusPending:
		return nil, models.ErrImagePending
	case models.StatusFailed:
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if _, err := r.data.DB.Master.ExecContext(ctx, `DELETE FROM tasks WHERE id = $1`, id); err != nil {
			return nil, errors.Join(models.ErrImageFailed, err)
		}

		return nil, models.ErrImageFailed
	case models.StatusCompleted:
		return &models.GetImageResponse{
			OriginalUrl:    fmt.Sprintf("http://%s/%s/%s/%s", r.str.ClientEndpoint, r.str.BucketName, models.OriginalKey, path),
			WatermarkedUrl: fmt.Sprintf("http://%s/%s/%s/%s", r.str.ClientEndpoint, r.str.BucketName, models.WatermarkedKey, path),
			ThumbnailUrl:   fmt.Sprintf("http://%s/%s/%s/%s", r.str.ClientEndpoint, r.str.BucketName, models.ThumbnailKey, path),
		}, nil
	default:
		return nil, models.ErrInvalidStatus
	}
}

func (r *Repository) DeleteImage(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var status, path string
	if err := r.data.DB.QueryRowContext(ctx, `SELECT status, path FROM tasks WHERE id = $1`, id).Scan(&status, &path); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ErrImageNotFound
		}
		return err
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tx, err := r.data.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM tasks WHERE id = $1`, id); err != nil {
		return err
	}

	for _, v := range []string{models.OriginalKey, models.WatermarkedKey, models.ThumbnailKey} {
		if err := r.DeleteImageByType(path, v); err != nil {
			return err
		}
	}

	err = retry.DoContext(ctx, retry.Strategy{Attempts: 3, Delay: 1, Backoff: 1}, tx.Commit)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteImageByType(path string, genre string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := retry.DoContext(ctx, retry.Strategy{Attempts: 3, Delay: 1, Backoff: 1}, func() error {
		return r.str.DeleteFile(fmt.Sprintf("%s/%s", genre, path))
	})
	if err != nil {
		return errors.Join(models.ErrImageDeleteFailed, err)
	}
	return nil
}
