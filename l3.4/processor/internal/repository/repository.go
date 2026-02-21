package repository

import (
	"app/internal/models"
	"app/pkg/data"
	"app/pkg/storage"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/wb-go/wbf/retry"
)

type Repository struct {
	data *data.Data
	str  *storage.StorageClient
}

func New(data *data.Data, str *storage.StorageClient) *Repository {
	return &Repository{data: data, str: str}
}

func (r *Repository) UploadImages(task *models.ProcessedTask) error {
	if err := r.str.UploadFile(fmt.Sprintf("%s/%s", models.OriginalKey, task.Path), task.Original.Raw, task.Original.Ext); err != nil {
		return errors.Join(models.ErrImageUploadFailed, err)
	}
	if err := r.str.UploadFile(fmt.Sprintf("%s/%s", models.WatermarkedKey, task.Path), task.Watermarked.Raw, task.Watermarked.Ext); err != nil {
		return errors.Join(models.ErrImageUploadFailed, err)
	}
	if err := r.str.UploadFile(fmt.Sprintf("%s/%s", models.ThumbnailKey, task.Path), task.Thumbnail.Raw, task.Thumbnail.Ext); err != nil {
		return errors.Join(models.ErrImageUploadFailed, err)
	}

	ctx, canc := context.WithTimeout(context.Background(), 5*time.Second)
	defer canc()

	_, err := r.data.DB.ExecWithRetry(ctx, retry.Strategy{Attempts: 3, Delay: time.Second, Backoff: 1}, `UPDATE tasks SET status = 'completed', path = $1 WHERE id = $2`, task.Path, task.ID)
	if err != nil {
		return err
	}
	return nil
}
