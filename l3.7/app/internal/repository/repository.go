package repository

import (
	"app/internal/models"
	"app/pkg/data"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
)

type Repository struct {
	data *data.Data
}

func New(data *data.Data) *Repository {
	return &Repository{data: data}
}

func (r *Repository) CreateItem(ctx context.Context, item models.Item, userId, userRole string) (*models.Item, error) {
	_, err := r.data.DB.ExecContext(ctx, `SELECT set_config('app.current_user_id', $1, false), set_config('app.current_user_role', $2, false)`, userId, userRole)
	if err != nil {
		return nil, errors.Join(models.ErrOnDatabase, err)
	}

	q := `INSERT INTO items (id, name, description, quantity, price) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`

	if err := r.data.DB.QueryRowContext(ctx, q, item.ID, item.Name, item.Description, item.Quantity, item.Price).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt); err != nil {
		return nil, errors.Join(models.ErrOnDatabase, err)
	}

	return &item, nil
}

func (r *Repository) GetItem(ctx context.Context, id string) (*models.Item, error) {
	q := `SELECT id, name, description, quantity, price, created_at, updated_at FROM items WHERE id = $1`

	var item models.Item
	if err := r.data.DB.QueryRowContext(ctx, q, id).Scan(&item.ID, &item.Name, &item.Description, &item.Quantity, &item.Price, &item.CreatedAt, &item.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(models.ErrItemNotFound, err)
		}
		return nil, errors.Join(models.ErrOnDatabase, err)
	}

	return &item, nil
}

func (r *Repository) ListItems(ctx context.Context) ([]models.Item, error) {
	q := `SELECT id, name, description, quantity, price, created_at, updated_at FROM items`

	var items []models.Item
	rows, err := r.data.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, errors.Join(models.ErrOnDatabase, err)
	}
	defer rows.Close()

	for rows.Next() {
		var item models.Item
		err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.Quantity, &item.Price, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			return nil, errors.Join(models.ErrOnDatabase, err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Join(models.ErrOnDatabase, err)
	}

	return items, nil
}

func (r *Repository) UpdateItem(ctx context.Context, item models.Item, userId, userRole string) (*models.Item, error) {
	_, err := r.data.DB.ExecContext(ctx, `SELECT set_config('app.current_user_id', $1, false), set_config('app.current_user_role', $2, false)`, userId, userRole)
	if err != nil {
		return nil, errors.Join(models.ErrOnDatabase, err)
	}

	q := `UPDATE items SET name = $1, description = $2, quantity = $3, price = $4 WHERE id = $5 RETURNING name, description, quantity, price, created_at, updated_at`

	if err := r.data.DB.QueryRowContext(ctx, q, item.Name, item.Description, item.Quantity, item.Price, item.ID).Scan(&item.Name, &item.Description, &item.Quantity, &item.Price, &item.CreatedAt, &item.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(models.ErrItemNotFound, err)
		}
		return nil, errors.Join(models.ErrOnDatabase, err)
	}

	return &item, nil
}

func (r *Repository) DeleteItem(ctx context.Context, id string, userId, userRole string) error {
	_, err := r.data.DB.ExecContext(ctx, `SELECT set_config('app.current_user_id', $1, false), set_config('app.current_user_role', $2, false)`, userId, userRole)
	if err != nil {
		return errors.Join(models.ErrOnDatabase, err)
	}

	q := `DELETE FROM items WHERE id = $1`

	res, err := r.data.DB.ExecContext(ctx, q, id)
	if err != nil {
		return errors.Join(models.ErrOnDatabase, err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.Join(models.ErrOnDatabase, err)
	}
	if rowsAffected == 0 {
		return models.ErrItemNotFound
	}

	return nil
}

func (r *Repository) ListHistory(ctx context.Context) ([]models.ItemHistory, error) {
	q := `SELECT id, item_id, action, user_id, user_role, old_data, new_data, changed_at FROM items_history ORDER BY changed_at DESC`

	var items []models.ItemHistory
	rows, err := r.data.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, errors.Join(models.ErrOnDatabase, err)
	}
	defer rows.Close()

	for rows.Next() {
		var itemHistory models.ItemHistory
		var oldItem, newItem []byte
		err := rows.Scan(&itemHistory.ID, &itemHistory.ItemID, &itemHistory.Action, &itemHistory.UserID, &itemHistory.UserRole, &oldItem, &newItem, &itemHistory.ChangedAt)
		if err != nil {
			return nil, errors.Join(models.ErrOnDatabase, err)
		}

		itemHistory.OldData = &models.Item{}
		itemHistory.NewData = &models.Item{}
		if err = json.Unmarshal(oldItem, itemHistory.OldData); err != nil {
			if err.Error() == "unexpected end of JSON input" {
				itemHistory.OldData = nil
			} else {
				return nil, errors.Join(models.ErrUnmarshalRaw, err)
			}
		}
		if err = json.Unmarshal(newItem, itemHistory.NewData); err != nil {
			if err.Error() == "unexpected end of JSON input" {
				itemHistory.NewData = nil
			} else {
				return nil, errors.Join(models.ErrUnmarshalRaw, err)
			}
		}

		items = append(items, itemHistory)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Join(models.ErrOnDatabase, err)
	}

	return items, nil
}

func (r *Repository) GetUuidByRole(ctx context.Context, role string) (string, error) {
	q := `SELECT id FROM users WHERE role = $1`

	var id string
	if err := r.data.DB.QueryRowContext(ctx, q, role).Scan(&id); err != nil {
		return "", errors.Join(models.ErrOnDatabase, err)
	}

	return id, nil
}

func (r *Repository) GetRoleByUuid(ctx context.Context, id string) (string, error) {
	q := `SELECT role FROM users WHERE id = $1`

	var role string
	if err := r.data.DB.QueryRowContext(ctx, q, id).Scan(&role); err != nil {
		return "", errors.Join(models.ErrOnDatabase, err)
	}

	return role, nil
}
