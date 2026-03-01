package models

import (
	"errors"
)

type Entity struct {
	Data []string  `json:"data"`
	Args Arguments `json:"args"`
}

type Arguments struct {
	F []int  `json:"f"`
	D string `json:"d"`
	S bool   `json:"s"`
}

type Response struct {
	Data []string `json:"data"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var (
	ErrInvalidDataSize        = errors.New("invalid data size")
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
