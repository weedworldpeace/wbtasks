package repository

import (
	"calendar/internal/models"
	"calendar/pkg/data"
	"errors"
	"slices"
)

var (
	errNonExistEventId = errors.New("non exist event id")
)

type Repository struct {
	data *data.Data
}

func (r *Repository) CreateEvent(userEvent *models.UserEvent) error {
	r.data.Users[userEvent.UserId] = append(r.data.Users[userEvent.UserId], userEvent.Event)
	return nil
}

func (r *Repository) UpdateEvent(userEvent *models.UserEvent) error {
	for i, v := range r.data.Users[userEvent.UserId] {
		if v.EventId == userEvent.EventId {
			r.data.Users[userEvent.UserId][i].Message = userEvent.Message
			return nil
		}
	}
	return errNonExistEventId
}

func (r *Repository) DeleteEvent(userEvent *models.UserEvent) error {
	for i, v := range r.data.Users[userEvent.UserId] {
		if v.EventId == userEvent.EventId {
			r.data.Users[userEvent.UserId] = slices.Delete(r.data.Users[userEvent.UserId], i, i+1)
			return nil
		}
	}
	return errNonExistEventId
}

func (r *Repository) ReadEvents(userId string, from, to int64) []models.Event {
	result := make([]models.Event, 0)
	for _, v := range r.data.Users[userId] {
		if v.UnixTime >= from && v.UnixTime < to {
			result = append(result, v)
		}
	}
	slices.SortFunc(result, func(a, b models.Event) int {
		if a.UnixTime < b.UnixTime {
			return -1
		} else if a.UnixTime > b.UnixTime {
			return 1
		} else {
			return 0
		}
	})
	return result
}

func New(data *data.Data) *Repository {
	return &Repository{data}
}
