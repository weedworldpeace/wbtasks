package repository

import (
	"app/internal/models"
	"app/pkg/data"
	"context"
	"errors"
	"time"
)

type Repository struct {
	data *data.Data
}

func (r *Repository) CreateEvent(ctx context.Context, userEvent models.UserEvent) (*models.Event, error) {
	q := `INSERT INTO events (user_id, event_id, message, date) VALUES ($1, $2, $3, $4) RETURNING created_at, updated_at`
	if err := r.data.DB.QueryRowContext(ctx, q, userEvent.UserId, userEvent.EventId, userEvent.Message, userEvent.Date).Scan(&userEvent.CreatedAt, &userEvent.UpdatedAt); err != nil {
		return nil, errors.Join(models.ErrOnDatabase, err)
	}

	return &userEvent.Event, nil
}

func (r *Repository) UpdateEvent(ctx context.Context, userEvent models.UserEvent) (*models.Event, error) {
	q := `UPDATE events SET message = $1, updated_at = NOW() WHERE user_id = $2 AND event_id = $3 RETURNING date, created_at, updated_at`

	if err := r.data.DB.QueryRowContext(ctx, q, userEvent.Message, userEvent.UserId, userEvent.EventId).Scan(&userEvent.Date, &userEvent.CreatedAt, &userEvent.UpdatedAt); err != nil {
		return nil, errors.Join(models.ErrOnDatabase, err)
	}

	return &userEvent.Event, nil
}

func (r *Repository) DeleteEvent(ctx context.Context, userEvent models.UserEvent) error {
	q := `DELETE FROM events WHERE user_id = $1 AND event_id = $2`

	res, err := r.data.DB.ExecContext(ctx, q, userEvent.UserId, userEvent.EventId)
	if err != nil {
		return errors.Join(models.ErrOnDatabase, err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.Join(models.ErrOnDatabase, err)
	}
	if rowsAffected == 0 {
		return models.ErrNonExistEvent
	}

	return nil
}

func (r *Repository) ReadEvents(ctx context.Context, userId string, dateFrom, dateTo time.Time) ([]models.Event, error) {
	q := `SELECT event_id, message, date, created_at, updated_at FROM events WHERE user_id = $1 AND date >= $2 AND date < $3`

	rows, err := r.data.DB.QueryContext(ctx, q, userId, dateFrom, dateTo)
	if err != nil {
		return nil, errors.Join(models.ErrOnDatabase, err)
	}

	events := make([]models.Event, 0)

	for rows.Next() {
		var ev models.Event

		err = rows.Scan(&ev.EventId, &ev.Message, &ev.Date, &ev.CreatedAt, &ev.UpdatedAt)
		if err != nil {
			return nil, errors.Join(models.ErrOnDatabase, err)
		}
		events = append(events, ev)
	}

	if rows.Err() != nil {
		return nil, errors.Join(models.ErrOnDatabase, err)
	}

	return events, nil
}

func New(data *data.Data) *Repository {
	return &Repository{data}
}
