package data

import (
	"calendar/internal/models"
	"sync"
)

type Data struct {
	Users map[string][]models.Event
	Mu    sync.RWMutex
}

func New() *Data {
	return &Data{Users: make(map[string][]models.Event)}
}
