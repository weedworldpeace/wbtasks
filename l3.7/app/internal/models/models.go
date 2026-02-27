package models

import (
	"errors"
	"time"
)

var Roles = map[string]bool{"admin": true, "manager": true, "viewer": true}

type Item struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Quantity    int       `json:"quantity" db:"quantity"`
	Price       float64   `json:"price" db:"price"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type ItemHistory struct {
	ID        string    `json:"id" db:"id"`
	ItemID    string    `json:"item_id" db:"item_id"`
	Action    string    `json:"action" db:"action"`
	UserID    string    `json:"user_id" db:"user_id"`
	UserRole  string    `json:"user_role" db:"user_role"`
	OldData   *Item     `json:"old_data,omitempty" db:"old_data"`
	NewData   *Item     `json:"new_data,omitempty" db:"new_data"`
	ChangedAt time.Time `json:"changed_at" db:"changed_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var (
	ErrInvalidRequestBody     = errors.New("invalid request body")
	ErrInvalidItemID          = errors.New("invalid item ID")
	ErrInvalidItemName        = errors.New("invalid name")
	ErrInvalidItemDescription = errors.New("invalid description")
	ErrInvalidItemQuantity    = errors.New("invalid quantity")
	ErrInvalidItemPrice       = errors.New("invalid price")
	ErrInvalidRole            = errors.New("invalid role")
	ErrOnDatabase             = errors.New("error on database")
	ErrItemNotFound           = errors.New("item not found")
	ErrInvalidTimestamp       = errors.New("invalid timestamp")
	ErrInvalidLimit           = errors.New("invalid limit")
	ErrInvalidOffset          = errors.New("invalid offset")
	ErrListIsEmpty            = errors.New("list is empty")
	ErrInvalidPort            = errors.New("invalid port")
	ErrInvalidReleaseMode     = errors.New("invalid release mode")
	ErrInvalidSigningMethod   = errors.New("invalid signing method")
	ErrInvalidGenre           = errors.New("invalid genre")
	ErrUnmarshalRaw           = errors.New("unmarshal raw json from bd")
	ErrUnexpected             = errors.New("unexpected")
	ErrNoPermission           = errors.New("no permission")
)
