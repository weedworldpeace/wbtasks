package data

import "calendar/internal/models"

type Data struct {
	Users map[string][]models.Event
}

func New() *Data {
	return &Data{make(map[string][]models.Event)}
}
