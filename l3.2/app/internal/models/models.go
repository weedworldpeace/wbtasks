package models

import (
	"errors"
)

var (
	ErrNonExistURL       = errors.New("non exist url")
	ErrNonShorten        = errors.New("non exist shorten")
	ErrAlreadyExistURL   = errors.New("already exist url")
	ErrOccupiedShortCode = errors.New("occupied short code")
	ErrBadPort           = errors.New("bad port value")
	ErrBadReleaseMode    = errors.New("bad release mode value")
	ErrBadURL            = errors.New("bad url value")
)

type ShortenRequest struct {
	URL string `json:"url" validate:"required,url"`
}

type ShortenResponse struct {
	OriginalURL string `json:"original_url"`
	ShortCode   string `json:"short_code"`
}

type ClickData struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type AnalyticsResponse struct {
	ShortCode   string      `json:"short_code"`
	TotalClicks int         `json:"total_clicks"`
	ClicksByDay []ClickData `json:"clicks_by_day"`
	UserAgents  []string    `json:"user_agents"`
}
