package repository

import (
	"app/internal/models"
	"app/pkg/data/inmem"
	"encoding/json"

	"github.com/rabbitmq/amqp091-go"
	"github.com/wb-go/wbf/rabbitmq"
)

// type PublisherConfig struct {
// 	ExchangeName string
// 	RoutingKey   string
// 	ContentType  string
// }

type Repository struct {
	data *inmem.Data
	pub  *rabbitmq.Publisher
}

func New(data *inmem.Data, rabbitCh *rabbitmq.Channel) *Repository {
	return &Repository{data, rabbitmq.NewPublisher(rabbitCh, "main_exchange")}
}

func (r *Repository) CreateNotification(notif *models.Notification) error {
	marshalled, err := json.Marshal(notif)
	if err != nil {
		return err
	}

	r.data.Mu.Lock()
	r.data.Notifications[notif.Id] = notif
	r.data.Mu.Unlock()

	headers := amqp091.Table{"x-delay": (notif.SendingDate.Sub(notif.CreationDate)).Milliseconds()}
	err = r.pub.Publish(marshalled, "main_routing_key", "json", rabbitmq.PublishingOptions{Headers: headers})
	if err != nil {
		r.data.Mu.Lock()
		delete(r.data.Notifications, notif.Id)
		r.data.Mu.Unlock()
		return err
	}

	return nil
}

func (r *Repository) ReadNotification(id string) (*models.Notification, error) {
	r.data.Mu.Lock()
	defer r.data.Mu.Unlock()

	if v, b := r.data.Notifications[id]; !b {
		return nil, models.ErrNonExistId
	} else {
		return v, nil
	}
}

func (r *Repository) DeleteNotification(id string) error {
	r.data.Mu.Lock()
	defer r.data.Mu.Unlock()

	if _, b := r.data.Notifications[id]; !b {
		return models.ErrNonExistId
	} else {
		delete(r.data.Notifications, id)
		return nil
	}
}
