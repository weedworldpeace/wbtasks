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

func (r *Repository) CreateTransaction(ctx context.Context, trans models.Transaction) (*models.Transaction, error) {
	q := `INSERT INTO transactions (id, user_id, amount, type, category, description) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`

	if err := r.data.DB.QueryRowContext(ctx, q, trans.ID, trans.UserID, trans.Amount, trans.Type, trans.Category, trans.Description).Scan(&trans.ID, &trans.CreatedAt, &trans.UpdatedAt); err != nil {
		return nil, errors.Join(models.ErrOnDatabase, err)
	}

	return &trans, nil
}

func (r *Repository) GetTransaction(ctx context.Context, id string) (*models.Transaction, error) {
	q := `SELECT id, user_id, amount, type, category, description, created_at, updated_at FROM transactions WHERE id = $1`

	var trans models.Transaction
	if err := r.data.DB.QueryRowContext(ctx, q, id).Scan(&trans.ID, &trans.UserID, &trans.Amount, &trans.Type, &trans.Category, &trans.Description, &trans.CreatedAt, &trans.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(models.ErrTransactionNotFound, err)
		}
		return nil, errors.Join(models.ErrOnDatabase, err)
	}

	return &trans, nil
}

func (r *Repository) ListTransactions(ctx context.Context, from, to time.Time) ([]models.Transaction, error) {
	q := `SELECT id, user_id, amount, type, category, description, created_at, updated_at FROM transactions WHERE created_at BETWEEN $1 AND $2 ORDER BY created_at DESC`

	var transactions []models.Transaction
	rows, err := r.data.DB.QueryContext(ctx, q, from, to)
	if err != nil {
		return nil, errors.Join(models.ErrOnDatabase, err)
	}
	defer rows.Close()

	for rows.Next() {
		var trans models.Transaction
		err := rows.Scan(&trans.ID, &trans.UserID, &trans.Amount, &trans.Type, &trans.Category, &trans.Description, &trans.CreatedAt, &trans.UpdatedAt)
		if err != nil {
			return nil, errors.Join(models.ErrOnDatabase, err)
		}
		transactions = append(transactions, trans)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Join(models.ErrOnDatabase, err)
	}

	return transactions, nil
}

func (r *Repository) UpdateTransaction(ctx context.Context, trans models.Transaction) (*models.Transaction, error) {
	q := `UPDATE transactions SET category = $1, description = $2, updated_at = NOW() WHERE id = $3 RETURNING amount, type, created_at, updated_at`

	if err := r.data.DB.QueryRowContext(ctx, q, trans.Category, trans.Description, trans.ID).Scan(&trans.Amount, &trans.Type, &trans.CreatedAt, &trans.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(models.ErrTransactionNotFound, err)
		}
		return nil, errors.Join(models.ErrOnDatabase, err)
	}

	return &trans, nil
}

func (r *Repository) DeleteTransaction(ctx context.Context, id string) error {
	q := `DELETE FROM transactions WHERE id = $1`

	res, err := r.data.DB.ExecContext(ctx, q, id)
	if err != nil {
		return errors.Join(models.ErrOnDatabase, err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.Join(models.ErrOnDatabase, err)
	}
	if rowsAffected == 0 {
		return models.ErrTransactionNotFound
	}

	return nil
}

func (r *Repository) GetAnalytics(ctx context.Context, from, to time.Time) (*models.Analytics, error) {
	q := `SELECT COUNT(*) AS total,
		COALESCE(SUM(amount), 0) AS total_sum,
		COALESCE(SUM(CASE WHEN type = 'income' THEN amount END), 0) AS income_sum,
		COALESCE(SUM(CASE WHEN type = 'expense' THEN amount END), 0) AS expense_sum,
		COALESCE(AVG(amount), 0) AS total_avg,
		COALESCE(AVG(CASE WHEN type = 'income' THEN amount END), 0) AS income_avg,
		COALESCE(AVG(CASE WHEN type = 'expense' THEN amount END), 0) AS expense_avg,
		COALESCE(PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY amount), 0) AS total_median,
		COALESCE(PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY CASE WHEN type = 'income' THEN amount END), 0) AS income_median,
		COALESCE(PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY CASE WHEN type = 'expense' THEN amount END), 0) AS expense_median,
		COALESCE(PERCENTILE_CONT(0.9) WITHIN GROUP (ORDER BY amount), 0) AS total_percentile_90,
		COALESCE(PERCENTILE_CONT(0.9) WITHIN GROUP (ORDER BY CASE WHEN type = 'income' THEN amount END), 0) AS income_percentile_90,
		COALESCE(PERCENTILE_CONT(0.9) WITHIN GROUP (ORDER BY CASE WHEN type = 'expense' THEN amount END), 0) AS expense_percentile_90
		FROM transactions
		WHERE created_at >= $1 AND created_at <= $2`

	var analytics models.Analytics
	if err := r.data.DB.QueryRowContext(ctx, q, from, to).Scan(&analytics.Total, &analytics.Sum.Amount, &analytics.Sum.Income, &analytics.Sum.Expense, &analytics.Average.Amount, &analytics.Average.Income, &analytics.Average.Expense, &analytics.Median.Amount, &analytics.Median.Income, &analytics.Median.Expense, &analytics.Percentile90.Amount, &analytics.Percentile90.Income, &analytics.Percentile90.Expense); err != nil {
		return nil, errors.Join(models.ErrOnDatabase, err)
	}

	return &analytics, nil
}
