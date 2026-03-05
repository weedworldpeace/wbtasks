package service

import (
	"app/internal/models"
	"context"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type RepositoryInterface interface {
	CreateEvent(context.Context, models.UserEvent) (*models.Event, error)
	UpdateEvent(context.Context, models.UserEvent) (*models.Event, error)
	DeleteEvent(context.Context, models.UserEvent) error
	ReadEvents(context.Context, string, time.Time, time.Time) ([]models.Event, error)
}

func validUserId(userId string) error {
	if uuid.Validate(userId) != nil {
		return models.ErrInvalidUserId
	}
	return nil
}

func validDate(date string) error {
	_, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return errors.Join(models.ErrInvalidDate, err)
	}
	return nil
}

func validEventId(eventId string) error {
	if uuid.Validate(eventId) != nil {
		return models.ErrInvalidEventId
	}
	return nil
}

type Service struct {
	repo RepositoryInterface
	vld  validator.Validate
}

func (s *Service) CreateEvent(ctx context.Context, userEvent models.UserEvent, email string) (*models.Event, error) {
	if err := validUserId(userEvent.UserId); err != nil {
		return nil, err
	}
	if err := s.vld.Var(email, "required,email"); err != nil {
		return nil, errors.Join(models.ErrInvalidEmail, err)
	}

	userEvent.EventId = uuid.NewString()

	return s.repo.CreateEvent(ctx, userEvent)
}

func (s *Service) UpdateEvent(ctx context.Context, userEvent models.UserEvent) (*models.Event, error) {
	if err := validUserId(userEvent.UserId); err != nil {
		return nil, err
	}
	if err := validEventId(userEvent.EventId); err != nil {
		return nil, err
	}

	return s.repo.UpdateEvent(ctx, userEvent)
}

func (s *Service) DeleteEvent(ctx context.Context, userEvent models.UserEvent) error {
	if err := validUserId(userEvent.UserId); err != nil {
		return err
	}
	if err := validEventId(userEvent.EventId); err != nil {
		return err
	}

	return s.repo.DeleteEvent(ctx, userEvent)
}

func (s *Service) ReadEvents(ctx context.Context, userId, rawDate, genre string) ([]models.Event, error) {
	if err := validUserId(userId); err != nil {
		return []models.Event{}, err
	}
	if err := validDate(rawDate); err != nil {
		return []models.Event{}, err
	}

	dateFrom, _ := time.Parse(time.RFC3339, rawDate)

	switch genre {
	case "day":
		dateTo := dateFrom.AddDate(0, 0, 1)
		return s.repo.ReadEvents(ctx, userId, dateFrom, dateTo)
	case "week":
		dateTo := dateFrom.AddDate(0, 0, 7)
		return s.repo.ReadEvents(ctx, userId, dateFrom, dateTo)
	case "month":
		dateTo := dateFrom.AddDate(0, 1, 0)
		return s.repo.ReadEvents(ctx, userId, dateFrom, dateTo)
	default:
		return []models.Event{}, models.ErrInvalidTimePeriod
	}
}

func New(repo RepositoryInterface) *Service {
	return &Service{repo, *validator.New()}
}
