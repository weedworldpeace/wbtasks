package storage

import (
	"app/internal/models"
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type StorageClient struct {
	Client         *minio.Client
	BucketName     string
	Endpoint       string
	ClientEndpoint string
	UseSSL         bool
}

type StorageConfig struct {
	Endpoint       string `env:"MINIO_ENDPOINT" env-default:"localhost:9000"`
	ClientEndpoint string `env:"MINIO_CLIENT_ENDPOINT" env-default:"localhost:9000"`
	AccessKey      string `env:"MINIO_ACCESS_KEY" env-default:"minioadmin"`
	SecretKey      string `env:"MINIO_SECRET_KEY" env-default:"minioadmin"`
	BucketName     string `env:"MINIO_BUCKET_NAME" env-default:"images"`
	UseSSL         bool   `env:"MINIO_USE_SSL" env-default:"false"`
	Region         string `env:"MINIO_REGION" env-default:"us-east-1"`
}

func New(cfg StorageConfig) *StorageClient {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
		Region: cfg.Region,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to create MinIO client: %s", err.Error()))
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.BucketName)
	if err != nil {
		panic(fmt.Sprintf("failed to check if bucket exists: %s", err.Error()))
	}

	if !exists {
		err = client.MakeBucket(ctx, cfg.BucketName, minio.MakeBucketOptions{Region: cfg.Region})
		if err != nil && err.Error() != minio.BucketAlreadyExists {
			panic(fmt.Sprintf("failed to create bucket: %s", err.Error()))
		}

		policy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::` + cfg.BucketName + `/*"]}]}`
		err = client.SetBucketPolicy(ctx, cfg.BucketName, policy)
		if err != nil {
			panic(fmt.Sprintf("failed to set bucket policy: %s", err.Error()))
		}
	}

	return &StorageClient{
		Client:         client,
		BucketName:     cfg.BucketName,
		Endpoint:       cfg.Endpoint,
		UseSSL:         cfg.UseSSL,
		ClientEndpoint: cfg.ClientEndpoint,
	}
}

func (m *StorageClient) DownloadFile(path string) ([]byte, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	obj, err := m.Client.GetObject(ctx, m.BucketName, path, minio.GetObjectOptions{})
	if err != nil {
		return nil, "", fmt.Errorf("failed to get object %s: %w", path, err)
	}
	defer obj.Close()

	data, err := io.ReadAll(obj)
	if err != nil {
		return nil, "", errors.Join(fmt.Errorf("failed to read object data %s: %w", path, err))
	}

	stat, err := obj.Stat()
	if err != nil {
		return nil, "", errors.Join(fmt.Errorf("failed to stat object %s: %w", path, err))
	}

	return data, stat.ContentType, nil
}

func (m *StorageClient) DeleteFile(path string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.Client.RemoveObject(ctx, m.BucketName, path, minio.RemoveObjectOptions{})
	if err != nil {
		return errors.Join(models.ErrImageDeleteFailed, fmt.Errorf("%s: %w", path, err))
	}
	return nil
}
