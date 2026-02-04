package service

import (
	"app/internal/models"
	"context"
	"strings"
)

type RepoInterface interface {
	Create(ctx context.Context, parentID *string, content string) (*models.Comment, error)
	GetPaginated(ctx context.Context, params models.GetCommentsRequest) (*models.GetCommentsResponse, error)
	Delete(ctx context.Context, id string) error
}

type CommentService struct {
	repo RepoInterface
}

func New(repo RepoInterface) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) CreateComment(ctx context.Context, req models.CreateCommentRequest) (*models.Comment, error) {
	req.Content = strings.TrimSpace(req.Content)
	if req.Content == "" {
		return nil, models.ErrInvalidInput
	}

	if len(req.Content) > models.MaxContentLength {
		return nil, models.ErrInvalidInput
	}

	if req.ParentID != nil && *req.ParentID == "" {
		req.ParentID = nil
	}

	comment, err := s.repo.Create(ctx, req.ParentID, req.Content)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *CommentService) GetComments(ctx context.Context, req models.GetCommentsRequest) (*models.GetCommentsResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	result, err := s.repo.GetPaginated(ctx, req)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *CommentService) DeleteComment(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
