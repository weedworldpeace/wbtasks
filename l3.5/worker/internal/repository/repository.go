package repository

import (
	"context"
	"worker/pkg/data"
)

type Repository struct {
	data *data.Data
}

func New(data *data.Data) *Repository {
	return &Repository{data: data}
}

func (r *Repository) CleanupBookings() (int, error) {
	q := `DELETE FROM bookings WHERE expires_at < NOW() AND status = 'pending';`
	res, err := r.data.DB.ExecContext(context.Background(), q)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(rowsAffected), nil
}
