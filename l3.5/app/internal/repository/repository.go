package repository

import (
	"app/internal/models"
	"app/pkg/data"
	"context"
	"database/sql"
	"errors"
	"time"
)

type Repository struct {
	data *data.Data
}

func New(data *data.Data) *Repository {
	return &Repository{data: data}
}

func (r *Repository) CreateEvent(ctx context.Context, ev models.Event) (*models.CreateResponse, error) {
	q := `INSERT INTO events (event_id, title, date, total_seats, price, time_to_confirm) VALUES ($1, $2, $3, $4, $5, $6)`

	if _, err := r.data.DB.ExecContext(ctx, q, ev.EventID, ev.Title, ev.Date, ev.TotalSeats, ev.Price, ev.TimeToConfirm); err != nil {
		return nil, err
	}

	return &models.CreateResponse{EventID: ev.EventID}, nil
}

func (r *Repository) BookEvent(ctx context.Context, book models.BookRequest) (*models.BookResponse, error) {
	q1 := `SELECT total_seats, time_to_confirm, (SELECT COUNT(*) FROM bookings WHERE event_id = $1) as taken_count FROM events WHERE event_id = $1 FOR UPDATE;`
	q2 := `INSERT INTO bookings (booking_id, event_id, user_name, user_email, booked_at, expires_at) VALUES ($1, $2, $3, $4, $5, $6);`

	tx, err := r.data.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var totalSeats, takenCount, timeToConfirm int

	if err := tx.QueryRowContext(ctx, q1, book.EventID).Scan(&totalSeats, &timeToConfirm, &takenCount); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(models.ErrNonExistEventID, err)
		}
		return nil, err
	}

	if takenCount >= totalSeats {
		return nil, models.ErrSeatsAreTaken
	}

	bookedAt := time.Now()
	expiresAt := bookedAt.Add(time.Duration(timeToConfirm) * time.Second)

	if _, err := tx.ExecContext(ctx, q2, book.BookingID, book.EventID, book.UserName, book.UserEmail, bookedAt, expiresAt); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.BookResponse{BookingID: book.BookingID}, nil
}

func (r *Repository) ConfirmEvent(ctx context.Context, conf models.ConfirmRequest) error {
	q1 := `UPDATE bookings SET status = 'confirmed', confirmed_at = NOW() WHERE booking_id = $1;`

	res, err := r.data.DB.ExecContext(ctx, q1, conf.BookingID)
	if err != nil {
		return err
	}

	if rowsAffected, err := res.RowsAffected(); err == nil && rowsAffected == 0 {
		return models.ErrNonExistBookingID
	}

	return nil
}

func (r *Repository) GetEvent(ctx context.Context, eventId string) (*models.EventResponse, error) {
	q1 := `SELECT event_id, title, date, total_seats, price, created_at, time_to_confirm FROM events WHERE event_id = $1;`
	q2 := `SELECT booking_id, event_id, user_name, user_email, status, booked_at, expires_at, confirmed_at FROM bookings WHERE event_id = $1;`

	var ev models.Event
	if err := r.data.DB.QueryRowContext(ctx, q1, eventId).Scan(&ev.EventID, &ev.Title, &ev.Date, &ev.TotalSeats, &ev.Price, &ev.CreatedAt, &ev.TimeToConfirm); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(models.ErrNonExistEventID, err)
		}
		return nil, err
	}

	rows, err := r.data.DB.QueryContext(ctx, q2, eventId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bookings := make([]models.Booking, 0)
	for rows.Next() {
		var b models.Booking
		if err := rows.Scan(&b.BookingID, &b.EventID, &b.UserName, &b.UserEmail, &b.Status, &b.BookedAt, &b.ExpiresAt, &b.ConfirmedAt); err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &models.EventResponse{
		Event:     ev,
		Bookings:  bookings,
		FreeSeats: ev.TotalSeats - len(bookings),
	}, nil
}

func (r *Repository) ListEvents(ctx context.Context) ([]models.EventResponse, error) {
	q1 := `SELECT event_id, title, date, total_seats, price, created_at, time_to_confirm FROM events;`
	q2 := `SELECT booking_id, event_id, user_name, user_email, status, booked_at, expires_at, confirmed_at FROM bookings WHERE event_id = $1;`

	rows, err := r.data.DB.QueryContext(ctx, q1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.EventResponse

	for rows.Next() {
		var ev models.Event
		if err := rows.Scan(&ev.EventID, &ev.Title, &ev.Date, &ev.TotalSeats, &ev.Price, &ev.CreatedAt, &ev.TimeToConfirm); err != nil {
			return nil, err
		}

		bookRows, err := r.data.DB.QueryContext(ctx, q2, ev.EventID)
		if err != nil {
			return nil, err
		}

		bookings := make([]models.Booking, 0)
		for bookRows.Next() {
			var b models.Booking
			if err := bookRows.Scan(&b.BookingID, &b.EventID, &b.UserName, &b.UserEmail, &b.Status, &b.BookedAt, &b.ExpiresAt, &b.ConfirmedAt); err != nil {
				bookRows.Close()
				return nil, err
			}
			bookings = append(bookings, b)
		}
		bookRows.Close()

		if err := bookRows.Err(); err != nil {
			return nil, err
		}

		events = append(events, models.EventResponse{
			Event:     ev,
			Bookings:  bookings,
			FreeSeats: ev.TotalSeats - len(bookings),
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
