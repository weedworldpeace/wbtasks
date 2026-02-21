package models

import (
	"errors"
)

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

type GetImageResponse struct {
	OriginalUrl    string `json:"original_url,omitempty"`
	WatermarkedUrl string `json:"watermarked_url,omitempty"`
	ThumbnailUrl   string `json:"thumbnail_url,omitempty"`
}

type UploadImageResponse struct {
	ID string `json:"id"`
}

var (
	ErrInvalidReleaseMode = errors.New("invalid release mode")
	ErrInvalidPort        = errors.New("invalid port")
	ErrInvalidId          = errors.New("invalid id")
	ErrImageNotFound      = errors.New("image not found")
	ErrImagePending       = errors.New("image is still pending")
	ErrImageFailed        = errors.New("image processing failed")
	// ErrInvalidGenre        = errors.New("invalid genre")
	ErrInvalidImage        = errors.New("invalid image")
	ErrInvalidImageType    = errors.New("invalid image type")
	ErrInvalidStatus       = errors.New("invalid status")
	ErrImageUploadFailed   = errors.New("failed to upload image")
	ErrImageDeleteFailed   = errors.New("failed to delete image")
	ErrImageDownloadFailed = errors.New("failed to download image")
)
