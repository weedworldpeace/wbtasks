package repository

import (
	"app/internal/models"
	"app/pkg/data"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Repository struct {
	data *data.Data
}

func New(data *data.Data) *Repository {
	return &Repository{data: data}
}

func (r *Repository) Create(ctx context.Context, parentID *string, content string) (*models.Comment, error) {
	tx, err := r.data.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	id := uuid.New().String()
	var path string
	var depth int

	if parentID == nil || *parentID == "" {
		path = id
		depth = 0
	} else {
		var parentPath string
		var parentDepth int
		err = tx.QueryRowContext(ctx, `
			SELECT path, depth FROM comments 
			WHERE id = $1 AND deleted_at IS NULL
		`, *parentID).Scan(&parentPath, &parentDepth)
		if err != nil {
			return nil, models.ErrParentNotFound
		}

		if parentDepth >= models.MaxNestingDepth {
			return nil, models.ErrMaxDepthExceeded
		}

		path = parentPath + "/" + id
		depth = parentDepth + 1
	}

	now := time.Now()
	comment := &models.Comment{
		ID:        id,
		ParentID:  parentID,
		Content:   strings.TrimSpace(content),
		Path:      path,
		Depth:     depth,
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO comments (id, parent_id, content, path, depth, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, comment.ID, comment.ParentID, comment.Content, comment.Path, comment.Depth,
		comment.CreatedAt, comment.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert comment: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return comment, nil
}

func (r *Repository) GetPaginated(ctx context.Context, params models.GetCommentsRequest) (*models.GetCommentsResponse, error) {
	query := `
        SELECT id, parent_id, content, path, depth, created_at, updated_at
        FROM comments 
        WHERE deleted_at IS NULL
    `

	var args []interface{}
	argIdx := 1

	if params.ParentID != nil && *params.ParentID != "" {
		var parentPath string
		err := r.data.DB.QueryRowContext(ctx,
			`SELECT path FROM comments WHERE id = $1 AND deleted_at IS NULL`,
			*params.ParentID).Scan(&parentPath)

		if err != nil {
			return nil, models.ErrParentNotFound
		}

		query += fmt.Sprintf(` AND (path = $%d OR path LIKE $%d || '/%%')`, argIdx, argIdx)
		args = append(args, parentPath)
		argIdx++
	} else {
		query += ` AND parent_id IS NULL`
	}

	if params.Query != "" {
		query += fmt.Sprintf(` AND content ILIKE $%d`, argIdx)
		args = append(args, "%"+params.Query+"%")
		argIdx++
	}

	countQuery := strings.Replace(query,
		"SELECT id, parent_id, content, path, depth, created_at, updated_at",
		"SELECT COUNT(*)", 1)

	var total int64
	err := r.data.DB.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count comments: %w", err)
	}

	sortBy := "created_at"
	if params.SortBy != "" && (params.SortBy == "created_at" || params.SortBy == "updated_at" || params.SortBy == "depth") {
		sortBy = params.SortBy
	}

	order := "DESC"
	if params.Order == "asc" {
		order = "ASC"
	}

	limit := params.Limit
	if limit <= 0 {
		limit = models.DefaultPageSize
	}
	if limit > models.MaxPageSize {
		limit = models.MaxPageSize
	}

	page := params.Page
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	query += fmt.Sprintf(` ORDER BY %s %s LIMIT $%d OFFSET $%d`, sortBy, order, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.data.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query comments: %w", err)
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(
			&comment.ID,
			&comment.ParentID,
			&comment.Content,
			&comment.Path,
			&comment.Depth,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	totalPages := 0
	if total > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}

	return &models.GetCommentsResponse{
		Comments:   comments,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	var path string
	err := r.data.DB.QueryRowContext(ctx, `
		SELECT path FROM comments 
		WHERE id = $1 AND deleted_at IS NULL
	`, id).Scan(&path)
	if err != nil {
		return models.ErrCommentNotFound
	}

	_, err = r.data.DB.ExecContext(ctx, `
		UPDATE comments 
		SET deleted_at = $1 
		WHERE path LIKE $2 || '%'
	`, time.Now(), path)
	if err != nil {
		return fmt.Errorf("failed to delete comments: %w", err)
	}

	return nil
}
