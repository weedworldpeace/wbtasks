package service

import (
	"app/internal/models"
	"context"
	"strconv"
	"time"

	"github.com/google/uuid"
	"go.uber.org/multierr"
)

type transactionGenre string

const (
	transactionGenreCreate transactionGenre = "create"
	transactionGenreUpdate transactionGenre = "update"
)

type RepoInterface interface {
	CreateTransaction(context.Context, models.Transaction) (*models.Transaction, error)
	GetTransaction(context.Context, string) (*models.Transaction, error)
	ListTransactions(context.Context, time.Time, time.Time) ([]models.Transaction, error)
	UpdateTransaction(context.Context, models.Transaction) (*models.Transaction, error)
	DeleteTransaction(context.Context, string) error
	GetAnalytics(context.Context, time.Time, time.Time) (*models.Analytics, error)
}

type CommentService struct {
	repo RepoInterface
}

func New(repo RepoInterface) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) CreateTransaction(ctx context.Context, trans models.Transaction) (*models.Transaction, error) {
	trans.ID = uuid.NewString()

	if err := valid(trans, transactionGenreCreate); err != nil {
		return nil, err
	}

	return s.repo.CreateTransaction(ctx, trans)
}

func (s *CommentService) GetTransaction(ctx context.Context, id string) (*models.Transaction, error) {
	if err := validUuid(id); err != nil {
		return nil, err
	}

	return s.repo.GetTransaction(ctx, id)
}

func (s *CommentService) ListTransactions(ctx context.Context, from, to, limit, offset string) ([]models.Transaction, error) {
	t1, err := time.Parse(time.RFC3339, from)
	if err != nil {
		return nil, models.ErrInvalidTimestamp
	}

	t2, err := time.Parse(time.RFC3339, to)
	if err != nil {
		return nil, models.ErrInvalidTimestamp
	}

	l, err := strconv.Atoi(limit)
	if err != nil || l < 1 {
		return nil, models.ErrInvalidLimit
	}

	off, err := strconv.Atoi(offset)
	if err != nil || off < 0 {
		return nil, models.ErrInvalidOffset
	}

	res, err := s.repo.ListTransactions(ctx, t1, t2)
	if err != nil {
		return nil, err
	}

	if off >= len(res) {
		return nil, models.ErrListIsEmpty
	}

	if off+l > len(res) {
		return res[off:], nil
	}
	return res[off : off+l], nil
}

func (s *CommentService) UpdateTransaction(ctx context.Context, trans models.Transaction) (*models.Transaction, error) {
	if err := valid(trans, transactionGenreUpdate); err != nil {
		return nil, err
	}

	return s.repo.UpdateTransaction(ctx, trans)
}

func (s *CommentService) DeleteTransaction(ctx context.Context, id string) error {
	if err := validUuid(id); err != nil {
		return err
	}

	return s.repo.DeleteTransaction(ctx, id)
}

func (s *CommentService) GetAnalytics(ctx context.Context, from, to string) (*models.Analytics, error) {
	t1, err := time.Parse(time.RFC3339, from)
	if err != nil {
		return nil, models.ErrInvalidTimestamp
	}

	t2, err := time.Parse(time.RFC3339, to)
	if err != nil {
		return nil, models.ErrInvalidTimestamp
	}

	return s.repo.GetAnalytics(ctx, t1, t2)
}

func valid(t models.Transaction, genre transactionGenre) error {
	var resultErr error

	switch genre {
	case transactionGenreCreate:
		resultErr = multierr.Append(resultErr, validUuid(t.UserID))
		resultErr = multierr.Append(resultErr, validAmount(t.Amount))
		resultErr = multierr.Append(resultErr, validType(t.Type))
		resultErr = multierr.Append(resultErr, validCategory(t.Category))
		resultErr = multierr.Append(resultErr, validDescription(t.Description))
	case transactionGenreUpdate:
		resultErr = multierr.Append(resultErr, validUuid(t.ID))
		resultErr = multierr.Append(resultErr, validUuid(t.UserID))
		resultErr = multierr.Append(resultErr, validCategory(t.Category))
		resultErr = multierr.Append(resultErr, validDescription(t.Description))
	default:
		resultErr = multierr.Append(resultErr, models.ErrInvalidTransactionGenre)
	}

	return resultErr
}

func validUuid(u string) error {
	if _, err := uuid.Parse(u); err != nil {
		return models.ErrInvalidTransactionID
	}
	return nil
}

func validAmount(a float64) error {
	if a <= 0 {
		return models.ErrInvalidTransactionAmount
	}
	return nil
}

func validType(t models.TransactionType) error {
	if t != models.TransactionTypeIncome && t != models.TransactionTypeExpense {
		return models.ErrInvalidTransactionType
	}
	return nil
}

func validCategory(c string) error {
	if c == "" {
		return models.ErrInvalidTransactionCategory
	}
	return nil
}

func validDescription(c string) error {
	if c == "" {
		return models.ErrInvalidTransactionDescription
	}
	return nil
}
