package service

import (
	"app/internal/models"
	"context"
	"slices"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/multierr"
)

var secretKey = []byte("your-secret-key")

type Genre string

const (
	GenreCreate  Genre = "create"
	GenreRead    Genre = "read"
	GenreUpdate  Genre = "update"
	GenreDelete  Genre = "delete"
	GenreHistory Genre = "history"
)

type RepoInterface interface {
	CreateItem(context.Context, models.Item, string, string) (*models.Item, error)
	GetItem(context.Context, string) (*models.Item, error)
	ListItems(context.Context) ([]models.Item, error)
	UpdateItem(context.Context, models.Item, string, string) (*models.Item, error)
	DeleteItem(context.Context, string, string, string) error
	ListHistory(context.Context) ([]models.ItemHistory, error)
	GetUuidByRole(context.Context, string) (string, error)
	GetRoleByUuid(context.Context, string) (string, error)
}

type CommentService struct {
	repo RepoInterface
}

func New(repo RepoInterface) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) CreateItem(ctx context.Context, item models.Item) (*models.Item, error) {
	userId, userRole, err := s.CheckRole(ctx, GenreCreate)
	if err != nil {
		return nil, err
	}

	item.ID = uuid.NewString()

	if err := valid(item, GenreCreate); err != nil {
		return nil, err
	}

	return s.repo.CreateItem(ctx, item, userId, userRole)
}

func (s *CommentService) GetItem(ctx context.Context, id string) (*models.Item, error) {
	_, _, err := s.CheckRole(ctx, GenreRead)
	if err != nil {
		return nil, err
	}

	if err := validUuid(id); err != nil {
		return nil, err
	}

	return s.repo.GetItem(ctx, id)
}

func (s *CommentService) ListItems(ctx context.Context, limit, offset string) ([]models.Item, error) {
	_, _, err := s.CheckRole(ctx, GenreRead)
	if err != nil {
		return nil, err
	}

	l, err := strconv.Atoi(limit)
	if err != nil || l < 1 {
		return nil, models.ErrInvalidLimit
	}

	off, err := strconv.Atoi(offset)
	if err != nil || off < 0 {
		return nil, models.ErrInvalidOffset
	}

	res, err := s.repo.ListItems(ctx)
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

func (s *CommentService) UpdateItem(ctx context.Context, item models.Item) (*models.Item, error) {
	userId, userRole, err := s.CheckRole(ctx, GenreUpdate)
	if err != nil {
		return nil, err
	}

	if err := valid(item, GenreUpdate); err != nil {
		return nil, err
	}

	return s.repo.UpdateItem(ctx, item, userId, userRole)
}

func (s *CommentService) DeleteItem(ctx context.Context, id string) error {
	userId, userRole, err := s.CheckRole(ctx, GenreDelete)
	if err != nil {
		return err
	}

	if err := validUuid(id); err != nil {
		return err
	}

	return s.repo.DeleteItem(ctx, id, userId, userRole)
}

func (s *CommentService) ListHistory(ctx context.Context, limit, offset string) ([]models.ItemHistory, error) {
	if _, _, err := s.CheckRole(ctx, GenreHistory); err != nil {
		return nil, err
	}
	l, err := strconv.Atoi(limit)
	if err != nil || l < 1 {
		return nil, models.ErrInvalidLimit
	}

	off, err := strconv.Atoi(offset)
	if err != nil || off < 0 {
		return nil, models.ErrInvalidOffset
	}

	res, err := s.repo.ListHistory(ctx)
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

func (s *CommentService) GetToken(ctx context.Context, role string) (string, error) {
	if !validRole(role) {
		return "", models.ErrInvalidRole
	}

	id, err := s.repo.GetUuidByRole(ctx, role)
	if err != nil {
		return "", err
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user_id": id})

	tokString, err := tok.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokString, nil
}

func (s *CommentService) CheckRole(ctx context.Context, genre Genre) (string, string, error) {
	id := ctx.Value("user_id")

	idStr, b := id.(string)
	if !b {
		return "", "", models.ErrUnexpected
	}

	if err := validUuid(idStr); err != nil {
		return "", "", err
	}

	role, err := s.repo.GetRoleByUuid(ctx, idStr)
	if err != nil {
		return "", "", err
	}

	if !validRole(role) {
		return "", "", err
	}

	switch genre {
	case GenreCreate, GenreUpdate, GenreDelete:
		if slices.Contains([]string{"admin", "manager"}, role) {
			return idStr, role, nil
		}
	case GenreRead, GenreHistory:
		if slices.Contains([]string{"admin", "manager", "viewer"}, role) {
			return idStr, role, nil
		}
	default:
		return "", "", models.ErrUnexpected
	}

	return "", "", models.ErrNoPermission
}

func valid(item models.Item, genre Genre) error {
	var resultErr error

	switch genre {
	case GenreCreate:
		resultErr = multierr.Append(resultErr, validName(item.Name))
		resultErr = multierr.Append(resultErr, validDescription(item.Description))
		resultErr = multierr.Append(resultErr, validQuantity(item.Quantity))
		resultErr = multierr.Append(resultErr, validPrice(item.Price))
	case GenreUpdate:
		resultErr = multierr.Append(resultErr, validUuid(item.ID))
		resultErr = multierr.Append(resultErr, validName(item.Name))
		resultErr = multierr.Append(resultErr, validDescription(item.Description))
		resultErr = multierr.Append(resultErr, validQuantity(item.Quantity))
		resultErr = multierr.Append(resultErr, validPrice(item.Price))
	default:
		resultErr = multierr.Append(resultErr, models.ErrInvalidGenre)
	}

	return resultErr
}

func validUuid(u string) error {
	if _, err := uuid.Parse(u); err != nil {
		return models.ErrInvalidItemID
	}
	return nil
}

func validDescription(c string) error {
	if c == "" {
		return models.ErrInvalidItemDescription
	}
	return nil
}

func validName(c string) error {
	if c == "" {
		return models.ErrInvalidItemName
	}
	return nil
}

func validQuantity(q int) error {
	if q < 0 {
		return models.ErrInvalidItemQuantity
	}
	return nil
}

func validPrice(p float64) error {
	if p < 0 {
		return models.ErrInvalidItemPrice
	}
	return nil
}

func validRole(role string) bool {
	_, b := models.Roles[role]
	return b
}
