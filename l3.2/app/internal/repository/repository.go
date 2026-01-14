package repository

import (
	"app/internal/models"
	"app/pkg/data"
	"context"
	"database/sql"
	"net/http"
)

type Repository struct {
	data *data.Data
}

func New(data *data.Data) *Repository {
	return &Repository{data: data}
}

func (r *Repository) CreateLink(orig string, shorten string) error {
	q := "INSERT INTO urls (original_url, short_code, created_at, is_active) VALUES ($1, $2, NOW(), true) RETURNING id, original_url, short_code, created_at;"
	_, err := r.data.DB.ExecContext(context.Background(), q, orig, shorten)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) CheckOriginalLink(url string) bool {
	q := "SELECT short_code FROM urls WHERE original_url = $1;"

	err := r.data.DB.QueryRowContext(context.Background(), q, url).Scan()
	if err == sql.ErrNoRows {
		return false
	} else {
		return true
	}
}

func (r *Repository) CheckShorten(url string) bool {
	q := "SELECT id FROM urls WHERE short_code = $1;"

	err := r.data.DB.QueryRowContext(context.Background(), q, url).Scan()
	if err == sql.ErrNoRows {
		return false
	} else {
		return true
	}
}

func (r *Repository) Redirect(req *http.Request) (string, error) {
	q1 := "SELECT id, original_url FROM urls WHERE short_code = $1 AND is_active = true;"
	q2 := "INSERT INTO clicks (url_id, clicked_at, user_agent) VALUES ($1, NOW(), $2);"

	var originalURL, id string
	err := r.data.DB.QueryRowContext(context.Background(), q1, req.URL.String()[3:]).Scan(&id, &originalURL)
	if err != nil {
		return "", models.ErrNonExistURL
	} else {
		// go func() {
		// 	_, err = r.data.DB.ExecContext(context.Background(), q2, id, req.UserAgent())
		// }()
		_, err = r.data.DB.ExecContext(context.Background(), q2, id, req.UserAgent())
		if err != nil {
			return "", err
		}
		return originalURL, nil
	}
}

func (r *Repository) GetAnalytics(shortCode string) (*models.AnalyticsResponse, error) {
	var res models.AnalyticsResponse
	var urlID int

	err := r.data.DB.QueryRowContext(context.Background(), "SELECT id FROM urls WHERE short_code = $1", shortCode).Scan(&urlID)
	if err != nil {
		return nil, models.ErrNonExistURL
	}
	res.ShortCode = shortCode

	err = r.data.DB.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM clicks WHERE url_id = $1", urlID).Scan(&res.TotalClicks)
	if err != nil {
		return nil, err
	}

	rows, err := r.data.DB.QueryContext(context.Background(), `SELECT DATE(clicked_at) as date, COUNT(*) as count FROM clicks WHERE url_id = $1 GROUP BY DATE(clicked_at) ORDER BY date DESC`, urlID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cd models.ClickData
		err = rows.Scan(&cd.Date, &cd.Count)
		if err != nil {
			return nil, err
		}
		res.ClicksByDay = append(res.ClicksByDay, cd)
	}

	rows, err = r.data.DB.QueryContext(context.Background(), `SELECT user_agent FROM clicks WHERE url_id = $1 AND user_agent IS NOT NULL GROUP BY user_agent ORDER BY MAX(clicked_at) DESC LIMIT 20`, urlID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ua string
		err = rows.Scan(&ua)
		if err != nil {
			return nil, err
		}
		res.UserAgents = append(res.UserAgents, ua)
	}

	return &res, nil
}
