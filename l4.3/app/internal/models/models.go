package models

import (
	"context"
	"errors"
	"time"
)

type Event struct {
	EventId   string    `json:"event_id"`
	Message   string    `json:"message" binding:"required"`
	Date      time.Time `json:"date" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserEvent struct {
	UserId string `json:"user_id"`
	Event
}

type EventTask struct {
	Event Event
	Email string
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type ToLog struct {
	Level   string
	Message string
	Error   error
	Ctx     context.Context
}

var (
	ErrInvalidPort        = errors.New("invalid port")
	ErrInvalidReleaseMode = errors.New("invalid release mode")
	ErrInvalidRequestBody = errors.New("invalid request body")
	ErrInvalidTimePeriod  = errors.New("invalid time period")
	ErrInvalidUserId      = errors.New("invalid user id")
	ErrInvalidEventId     = errors.New("invalid event id")
	ErrInvalidDate        = errors.New("invalid date")
	ErrInvalidEmail       = errors.New("invalid email")
	ErrOnDatabase         = errors.New("on database")
	ErrNonExistEvent      = errors.New("non exist event")
	ErrUnsupportedLevel   = errors.New("unsupported level")
	ErrUnexpected         = errors.New("unexpected")
)
