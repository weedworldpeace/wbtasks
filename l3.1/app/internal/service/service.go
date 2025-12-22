package service

import (
	"app/internal/models"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

type RepositoryInterface interface {
	CreateNotification(*models.Notification) error
	ReadNotification(string) (*models.Notification, error)
	DeleteNotification(string) error
}

type Service struct {
	repo RepositoryInterface
}

func New(repo RepositoryInterface) *Service {
	return &Service{repo}
}

func (s Service) CreateNotification(notif *models.Notification) (string, error) { // not finished
	notif.Email = strings.TrimSpace(notif.Email)
	if err := validEmail(notif.Email); err != nil {
		return "", err
	}

	notif.CreationDate = time.Now()
	notif.Id = uuid.NewString()

	if err := s.repo.CreateNotification(notif); err != nil {
		return "", err
	} else {
		return notif.Id, nil
	}
}

func (s Service) ReadNotification(id string) (*models.Notification, error) { // not finished
	return s.repo.ReadNotification(id)
}

func (s Service) DeleteNotification(id string) error { // not finished
	return s.repo.DeleteNotification(id)
}

func validEmail(email string) error {
	if email == "" {
		return models.ErrBadEmail
	}

	if !emailRegex.MatchString(email) {
		return models.ErrBadEmail
	}

	return nil
}
