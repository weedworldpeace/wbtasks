package service

import (
	"app/internal/models"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var genres = map[string]bool{
	models.OriginalKey:    true,
	models.WatermarkedKey: true,
	models.ThumbnailKey:   true,
}

type RepositoryInterface interface {
	UploadImage(*models.KafkaTask) (*models.UploadImageResponse, error)
	GetImage(string) (*models.GetImageResponse, error)
	DeleteImage(string) error
}

type Service struct {
	repo RepositoryInterface
	vld  *validator.Validate
}

func New(repo RepositoryInterface) *Service {
	return &Service{repo, validator.New(validator.WithRequiredStructEnabled())}
}

func (s Service) UploadImage(img *models.Image) (*models.UploadImageResponse, error) {
	return s.repo.UploadImage(&models.KafkaTask{
		ID:    uuid.NewString(),
		Image: *img,
	})
}

func (s Service) GetImage(id string) (*models.GetImageResponse, error) {
	if err := isValidId(id); err != nil {
		return nil, errors.Join(models.ErrInvalidId, err)
	}

	return s.repo.GetImage(id)
}

func (s Service) DeleteImage(id string) error {
	if err := isValidId(id); err != nil {
		return errors.Join(models.ErrInvalidId, err)
	}

	return s.repo.DeleteImage(id)
}

func isValidGenre(genre string) bool {
	return genres[genre]
}

func isValidId(id string) error {
	_, err := uuid.Parse(id)
	return err
}
