package storage

import (
	"app/internal/models"
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type StorageClient struct {
	client     *minio.Client
	bucketName string
	endpoint   string
	useSSL     bool
}

type StorageConfig struct {
	Endpoint   string `env:"MINIO_ENDPOINT" env-default:"localhost:9000"`
	AccessKey  string `env:"MINIO_ACCESS_KEY" env-default:"minioadmin"`
	SecretKey  string `env:"MINIO_SECRET_KEY" env-default:"minioadmin"`
	BucketName string `env:"MINIO_BUCKET_NAME" env-default:"images"`
	UseSSL     bool   `env:"MINIO_USE_SSL" env-default:"false"`
	Region     string `env:"MINIO_REGION" env-default:"us-east-1"`
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
		client:     client,
		bucketName: cfg.BucketName,
		endpoint:   cfg.Endpoint,
		useSSL:     cfg.UseSSL,
	}
}

func (m *StorageClient) UploadFile(path string, data []byte, contentType string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	reader := bytes.NewReader(data)

	_, err := m.client.PutObject(ctx, m.bucketName, path, reader, int64(len(data)),
		minio.PutObjectOptions{
			ContentType: contentType,
			UserMetadata: map[string]string{
				"Uploaded": time.Now().Format(time.RFC3339),
			},
		})

	if err != nil {
		return errors.Join(models.ErrImageUploadFailed, err)
	}

	return nil
}
