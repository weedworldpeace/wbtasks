package models

import "errors"

const (
	StatusPending   = "pending"
	StatusCompleted = "completed"
	StatusFailed    = "failed"
	OriginalKey     = "original"
	WatermarkedKey  = "watermarked"
	ThumbnailKey    = "thumbnail"
)

type Image struct {
	Raw  []byte `json:"raw"`
	Ext  string `json:"extension"`
	Size int    `json:"size"`
}

type KafkaTask struct {
	ID string `json:"id"`
	Image
}

type PostgresTask struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Path   string `json:"path"`
}

type ProcessedTask struct {
	ID          string `json:"id"`
	Path        string `json:"path"`
	Original    Image  `json:"original"`
	Watermarked Image  `json:"watermarked"`
	Thumbnail   Image  `json:"thumbnail"`
}

var (
	ErrInvalidReleaseMode = errors.New("invalid release mode")
	ErrInvalidPort        = errors.New("invalid port")
	ErrInvalidWorkerCount = errors.New("invalid worker count")
	ErrInvalidLogLevel    = errors.New("invalid log level")
	ErrInvalidImageExt    = errors.New("invalid image extension")
	ErrImageUploadFailed  = errors.New("failed to upload image")
)
