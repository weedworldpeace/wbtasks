package service

import (
	"app/internal/models"
	"math/rand/v2"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type RepositoryInterface interface {
	CheckOriginalLink(string) bool
	CheckShorten(string) bool
	CreateLink(string, string) error
	Redirect(*http.Request) (string, error)
	GetAnalytics(string) (*models.AnalyticsResponse, error)
}

type Service struct {
	repo RepositoryInterface
	vld  *validator.Validate
}

func New(repo RepositoryInterface) *Service {
	return &Service{repo, validator.New(validator.WithRequiredStructEnabled())}
}

func (s Service) CreateLink(data *models.ShortenRequest) (string, error) {
	if err := s.vld.Struct(data); err != nil {
		return "", models.ErrBadURL
	}

	if s.repo.CheckOriginalLink(data.URL) {
		return "", models.ErrAlreadyExistURL
	}

	shorten := newShorten()
	for s.repo.CheckShorten(shorten) {
		shorten = newShorten()
	}

	err := s.repo.CreateLink(data.URL, shorten)
	if err != nil {
		return "", err
	}

	return shorten, nil

}

func (s Service) Redirect(req *http.Request) (string, error) {
	return s.repo.Redirect(req)
}

func (s Service) GetAnalytics(url string) (*models.AnalyticsResponse, error) {
	return s.repo.GetAnalytics(url)
}

func newShorten() string {
	var shorten []byte
	charset1 := "abcdefghijklmnopqrstuvwxyz"
	for range 3 {
		shorten = append(shorten, charset1[rand.IntN(len(charset1))])
	}
	charset2 := "0123456789"
	for range 3 {
		shorten = append(shorten, charset2[rand.IntN(len(charset2))])
	}
	return string(shorten)
}
