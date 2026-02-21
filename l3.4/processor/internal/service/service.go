package service

import (
	"app/internal/models"
	"bytes"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
	"time"

	"github.com/disintegration/imaging"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type RepositoryInterface interface {
	UploadImages(*models.ProcessedTask) error
}

type Service struct {
	repo RepositoryInterface
	vld  *validator.Validate
}

func New(repo RepositoryInterface) *Service {
	return &Service{repo, validator.New(validator.WithRequiredStructEnabled())}
}

func (s Service) ProcessTask(task *models.KafkaTask) error {
	time.Sleep(5 * time.Second)

	path := uuid.NewString()

	switch task.Ext {
	case "image/jpeg":
		img, err := jpeg.Decode(bytes.NewReader(task.Raw))
		if err != nil {
			return err
		}

		thumbnail := imaging.Resize(img, 150, 150, imaging.Lanczos)

		watermarked := imaging.Clone(img)

		if err = addWatermark(watermarked); err != nil {
			return err
		}

		result := &models.ProcessedTask{ID: task.ID, Path: path, Original: task.Image}

		watermarkedBuf := new(bytes.Buffer)
		if err := jpeg.Encode(watermarkedBuf, watermarked, &jpeg.Options{Quality: 95}); err != nil {
			return err
		}
		result.Watermarked = models.Image{Raw: watermarkedBuf.Bytes(), Ext: "jpeg", Size: watermarkedBuf.Len()}

		thumbnailBuf := new(bytes.Buffer)
		if err := jpeg.Encode(thumbnailBuf, thumbnail, &jpeg.Options{Quality: 85}); err != nil {
			return err
		}
		result.Thumbnail = models.Image{Raw: thumbnailBuf.Bytes(), Ext: "jpeg", Size: thumbnailBuf.Len()}

		return s.repo.UploadImages(result)
	default:
		return models.ErrInvalidImageExt
	}
}

func addWatermark(img draw.Image) error {
	watermarkFile, err := os.Open("./watermark.png")
	if err != nil {
		return err
	}
	defer watermarkFile.Close()

	watermark, _, err := image.Decode(watermarkFile)
	if err != nil {
		return err
	}

	bounds := img.Bounds()
	watermarkBounds := watermark.Bounds()

	x := bounds.Max.X - watermarkBounds.Max.X - 10
	y := bounds.Max.Y - watermarkBounds.Max.Y - 10

	draw.Draw(img, bounds, watermark, image.Point{-x, -y}, draw.Over)

	return nil
}
