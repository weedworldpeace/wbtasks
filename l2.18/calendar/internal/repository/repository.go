package repository

import (
	"calendar/internal/models"
	"calendar/pkg/data"
	"errors"
	"slices"
	"time"
)

var (
	ErrNonExistEventId = errors.New("non exist event id")
	ErrNonExistUserId  = errors.New("non exist user id")
)

type Repository struct {
	data *data.Data
}

func (r *Repository) checkUserId(userId string) bool {
	r.data.Mu.RLock()
	defer r.data.Mu.RUnlock()
	_, b := r.data.Users[userId]
	return b
}

func (r *Repository) CreateEvent(userEvent *models.UserEvent) error {
	r.data.Mu.Lock()
	defer r.data.Mu.Unlock()
	unixDate, _ := time.Parse("2006-01-02", userEvent.Date)
	for i, v := range r.data.Users[userEvent.UserId] {
		curUnixDate, _ := time.Parse("2006-01-02", v.Date)
		if curUnixDate.Unix() > unixDate.Unix() {
			r.data.Users[userEvent.UserId] = slices.Insert(r.data.Users[userEvent.UserId], i, userEvent.Event)
			return nil
		}
	}
	r.data.Users[userEvent.UserId] = append(r.data.Users[userEvent.UserId], userEvent.Event)
	return nil
}

func (r *Repository) UpdateEvent(userEvent *models.UserEvent) error {
	if !r.checkUserId(userEvent.UserId) {
		return ErrNonExistUserId
	}

	r.data.Mu.Lock()
	defer r.data.Mu.Unlock()

	for i, v := range r.data.Users[userEvent.UserId] {
		if v.EventId == userEvent.EventId {
			r.data.Users[userEvent.UserId][i].Message = userEvent.Message
			return nil
		}
	}
	return ErrNonExistEventId
}

func (r *Repository) DeleteEvent(userEvent *models.UserEvent) error {
	if !r.checkUserId(userEvent.UserId) {
		return ErrNonExistUserId
	}

	r.data.Mu.Lock()
	defer r.data.Mu.Unlock()

	for i, v := range r.data.Users[userEvent.UserId] {
		if v.EventId == userEvent.EventId {
			r.data.Users[userEvent.UserId] = slices.Delete(r.data.Users[userEvent.UserId], i, i+1)
			return nil
		}
	}
	return ErrNonExistEventId
}

func (r *Repository) ReadEvents(userId string, dateFrom, dateTo int64) ([]models.Event, error) {
	result := make([]models.Event, 0)

	if !r.checkUserId(userId) {
		return result, ErrNonExistUserId
	}

	r.data.Mu.RLock()
	defer r.data.Mu.RUnlock()

	for _, v := range r.data.Users[userId] {
		curUnixDate, _ := time.Parse("2006-01-02", v.Date)
		if curUnixDate.Unix() >= dateFrom {
			if curUnixDate.Unix() < dateTo {
				result = append(result, v)
			} else {
				return result, nil
			}
		}
	}
	return result, nil
}

func New(data *data.Data) *Repository {
	return &Repository{data}
}
