package service

import (
	"calendar/internal/models"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	errInvalidUserId     = errors.New("invalid user id")
	errInvalidEventId    = errors.New("invalid event id")
	errInvalidDate       = errors.New("invalid date")
	errInvalidTimePeriod = errors.New("invalid time period")
)

type RepositoryInterface interface {
	CreateEvent(*models.UserEvent) error
	UpdateEvent(*models.UserEvent) error
	DeleteEvent(*models.UserEvent) error
	ReadEvents(string, int64, int64) []models.Event
}

func validUserId(userId string) error {
	// if uuid.Validate(userId) != nil { // to backkkk
	// 	return errInvalidUserId
	// }
	return nil
}

func validDate(date string) (int64, error) {
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return 0, fmt.Errorf("%v: %v", errInvalidDate, err)
	}
	return parsedDate.Unix(), nil
}

func validEventId(eventId string) error {
	if uuid.Validate(eventId) != nil {
		return errInvalidEventId
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
	unixTime, err := validDate(userEvent.RawTime)
	if err != nil {
		return err
	}
	userEvent.UnixTime = unixTime
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

func (s *Service) ReadEvents(userId, date, genre string) ([]models.Event, error) {
	if err := validUserId(userId); err != nil {
		return []models.Event{}, err
	}
	unixFromTime, err := validDate(date)
	if err != nil {
		return []models.Event{}, err
	}

	switch genre {
	case "day":
		unixToTime := time.Unix(unixFromTime, 0).AddDate(0, 0, 1).Unix()
		return s.repo.ReadEvents(userId, unixFromTime, unixToTime), nil
	case "week":
		unixToTime := time.Unix(unixFromTime, 0).AddDate(0, 0, 7).Unix()
		return s.repo.ReadEvents(userId, unixFromTime, unixToTime), nil
	case "month":
		unixToTime := time.Unix(unixFromTime, 0).AddDate(0, 1, 0).Unix()
		return s.repo.ReadEvents(userId, unixFromTime, unixToTime), nil
	default:
		return []models.Event{}, errInvalidTimePeriod
	}
}

func New(repo RepositoryInterface) *Service {
	return &Service{repo}
}
