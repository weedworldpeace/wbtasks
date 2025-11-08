package service

import (
	"calendar/internal/models"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidUserId     = errors.New("invalid user id")
	ErrInvalidEventId    = errors.New("invalid event id")
	ErrInvalidDate       = errors.New("invalid date")
	ErrInvalidTimePeriod = errors.New("invalid time period")
)

type RepositoryInterface interface {
	CreateEvent(*models.UserEvent) error
	UpdateEvent(*models.UserEvent) error
	DeleteEvent(*models.UserEvent) error
	ReadEvents(string, int64, int64) ([]models.Event, error)
}

func validUserId(userId string) error {
	if uuid.Validate(userId) != nil {
		return ErrInvalidUserId
	}
	return nil
}

func validDate(date string) error {
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		return errors.Join(ErrInvalidDate, err)
	}
	return nil
}

func validEventId(eventId string) error {
	if uuid.Validate(eventId) != nil {
		return ErrInvalidEventId
	}
	return nil
}

type Service struct {
	repo RepositoryInterface
}

func (s *Service) CreateEvent(userEvent *models.UserEvent) error {
	if err := validUserId(userEvent.UserId); err != nil {
		return err
	}
	if err := validDate(userEvent.Date); err != nil {
		return err
	}
	userEvent.EventId = uuid.NewString()

	return s.repo.CreateEvent(userEvent)
}

func (s *Service) UpdateEvent(userEvent *models.UserEvent) error {
	if err := validUserId(userEvent.UserId); err != nil {
		return err
	}
	if err := validEventId(userEvent.EventId); err != nil {
		return err
	}

	return s.repo.UpdateEvent(userEvent)
}

func (s *Service) DeleteEvent(userEvent *models.UserEvent) error {
	if err := validUserId(userEvent.UserId); err != nil {
		return err
	}
	if err := validEventId(userEvent.EventId); err != nil {
		return err
	}

	return s.repo.DeleteEvent(userEvent)
}

func (s *Service) ReadEvents(userId, rawDate, genre string) ([]models.Event, error) {
	if err := validUserId(userId); err != nil {
		return []models.Event{}, err
	}
	if err := validDate(rawDate); err != nil {
		return []models.Event{}, err
	}

	dateFrom, _ := time.Parse("2006-01-02", rawDate)

	switch genre {
	case "day":
		dateTo := dateFrom.AddDate(0, 0, 1)
		return s.repo.ReadEvents(userId, dateFrom.Unix(), dateTo.Unix())
	case "week":
		dateTo := dateFrom.AddDate(0, 0, 7)
		return s.repo.ReadEvents(userId, dateFrom.Unix(), dateTo.Unix())
	case "month":
		dateTo := dateFrom.AddDate(0, 1, 0)
		return s.repo.ReadEvents(userId, dateFrom.Unix(), dateTo.Unix())
	default:
		return []models.Event{}, ErrInvalidTimePeriod
	}
}

func New(repo RepositoryInterface) *Service {
	return &Service{repo}
}
