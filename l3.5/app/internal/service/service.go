package service

import (
	"app/internal/models"
	"context"

	"github.com/google/uuid"
)

type RepoInterface interface {
	CreateEvent(context.Context, models.Event) (*models.CreateResponse, error)
	BookEvent(context.Context, models.BookRequest) (*models.BookResponse, error)
	ConfirmEvent(context.Context, models.ConfirmRequest) error
	GetEvent(context.Context, string) (*models.EventResponse, error)
	ListEvents(context.Context) ([]models.EventResponse, error)
}

type CommentService struct {
	repo RepoInterface
}

func New(repo RepoInterface) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) CreateEvent(ctx context.Context, ev models.Event) (*models.CreateResponse, error) {
	ev.EventID = uuid.NewString()

	return s.repo.CreateEvent(ctx, ev)
}

func (s *CommentService) BookEvent(ctx context.Context, book models.BookRequest) (*models.BookResponse, error) {
	if !validUuid(book.EventID) {
		return nil, models.ErrInvalidEventID
	}

	book.BookingID = uuid.NewString()

	return s.repo.BookEvent(ctx, book)
}

func (s *CommentService) ConfirmEvent(ctx context.Context, conf models.ConfirmRequest) error {
	if !validUuid(conf.EventID) {
		return models.ErrInvalidEventID
	}
	if !validUuid(conf.BookingID) {
		return models.ErrInvalidBookingID
	}

	return s.repo.ConfirmEvent(ctx, conf)
}

func (s *CommentService) GetEvent(ctx context.Context, eventId string) (*models.EventResponse, error) {
	if !validUuid(eventId) {
		return nil, models.ErrInvalidEventID
	}

	return s.repo.GetEvent(ctx, eventId)
}

func (s *CommentService) ListEvents(ctx context.Context) ([]models.EventResponse, error) {
	return s.repo.ListEvents(ctx)
}

func validUuid(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
