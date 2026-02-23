package models

import (
	"errors"
	"time"
)

const (
	StatusPending   = "pending"
	StatusConfirmed = "confirmed"
	StatusCancelled = "cancelled"
)

type Event struct {
	EventID       string    `json:"event_id"`
	Title         string    `json:"title" binding:"required"`
	Date          time.Time `json:"date" binding:"required"`
	TotalSeats    int       `json:"total_seats" binding:"required,min=1"`
	Price         float64   `json:"price" binding:"required,min=0"`
	CreatedAt     time.Time `json:"created_at"`
	TimeToConfirm int       `json:"time_to_confirm" binding:"required,min=1"` // in seconds
}

type User struct {
	UserName  string `json:"user_name" binding:"required"`
	UserEmail string `json:"user_email" binding:"required,email"`
}

type Booking struct {
	BookingID string `json:"booking_id"`
	EventID   string `json:"event_id,omitempty"`
	User
	Status      string     `json:"status"`
	BookedAt    time.Time  `json:"booked_at"`
	ExpiresAt   time.Time  `json:"expires_at"`
	ConfirmedAt *time.Time `json:"confirmed_at,omitempty"`
}

type BookRequest struct {
	BookingID string `json:"booking_id"`
	EventID   string `json:"event_id"`
	User
}

type ConfirmRequest struct {
	BookingID string `json:"booking_id" binding:"required"`
	EventID   string `json:"event_id"`
}

type EventResponse struct {
	Event     Event     `json:"event"`
	Bookings  []Booking `json:"bookings"`
	FreeSeats int       `json:"free_seats"`
}

type CreateResponse struct {
	EventID string `json:"event_id"`
}

type BookResponse struct {
	BookingID string `json:"booking_id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var (
	ErrInvalidPort        = errors.New("invalid port")
	ErrInvalidReleaseMode = errors.New("invalid release mode")
	ErrInvalidRequestBody = errors.New("invalid request body")
	ErrInvalidEventID     = errors.New("invalid event ID")
	ErrInvalidBookingID   = errors.New("invalid booking ID")
	ErrSeatsAreTaken      = errors.New("all seats are already taken")
	ErrNonExistEventID    = errors.New("event with given ID does not exist")
	ErrNonExistBookingID  = errors.New("booking with given ID does not exist")
)
