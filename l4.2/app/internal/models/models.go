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

type ChStruct struct {
	Data []string
	Err  error
}

var (
	ErrBadFlagD               = errors.New("the delimiter must be a single character")
	ErrNoFlagF                = errors.New("you must specify fields")
	ErrInvalidFieldValue      = errors.New("invalid field value")
	ErrInvalidDecreasingRange = errors.New("invalid decreasing range")
	ErrInvalidFieldRange      = errors.New("invalid field range")
	ErrFieldsNumberedFrom     = errors.New("fields are numbered from 1")
	ErrInvalidRangeEndpoint   = errors.New("invalid range with no endpoint")
	ErrInvalidPort            = errors.New("invalid port")
	ErrUnexpected             = errors.New("unexpected")
	ErrNoQuorum               = errors.New("no quorum")
	ErrInvalidWrkCount        = errors.New("invalid wrk count: should be between 1 and 20")
)

const (
	ExitErrOccurred = 1
)
