package inmem

import (
	"app/internal/models"
	"sync"
)

type Data struct {
	Notifications map[string]*models.Notification
	Mu            sync.Mutex
}

func New() *Data {
	return &Data{Notifications: make(map[string]*models.Notification)}
}
