package models

import (
	"errors"
	"time"
)

type Comment struct {
	ID        string     `json:"id" db:"id"`
	ParentID  *string    `json:"parent_id,omitempty" db:"parent_id"`
	Content   string     `json:"content" db:"content"`
	Path      string     `json:"path,omitempty" db:"path"`
	Depth     int        `json:"depth" db:"depth"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type CreateCommentRequest struct {
	ParentID *string `json:"parent_id,omitempty" validate:"omitempty,uuid4"`
	Content  string  `json:"content" validate:"required,min=1,max=5000"`
}

type GetCommentsRequest struct {
	ParentID *string `form:"parent_id"`
	Page     int     `form:"page" validate:"min=1"`
	Limit    int     `form:"limit" validate:"min=1,max=100"`
	Query    string  `form:"query"`
	SortBy   string  `form:"sort_by" validate:"oneof=created_at updated_at depth"`
	Order    string  `form:"order" validate:"oneof=asc desc"`
}

type GetCommentsResponse struct {
	Comments   []Comment `json:"comments"`
	Total      int64     `json:"total"`
	Page       int       `json:"page"`
	Limit      int       `json:"limit"`
	TotalPages int       `json:"total_pages"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var (
	ErrInvalidInput     = errors.New("invalid input")
	ErrCommentNotFound  = errors.New("comment not found")
	ErrParentNotFound   = errors.New("parent comment not found")
	ErrMaxDepthExceeded = errors.New("maximum nesting depth exceeded")
	ErrBadPort          = errors.New("bad port")
	ErrBadReleaseMode   = errors.New("bad release mode")
)

const (
	MaxContentLength     = 5000
	MaxNestingDepth      = 20
	DefaultPageSize      = 20
	MaxPageSize          = 100
	MinSearchQueryLength = 2
)
