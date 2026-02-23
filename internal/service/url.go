package service

import (
	"context"
	"crypto/rand"
	"math/big"

	"github.com/Haruto-works333/url-shortener/internal/model"
	"github.com/Haruto-works333/url-shortener/internal/repository"
)

const (
	shortCodeLength = 7
	charset         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// URLService contains business logic for URL operations.
type URLService struct {
	repo *repository.URLRepository
}

// NewURLService creates a new URLService.
func NewURLService(repo *repository.URLRepository) *URLService {
	return &URLService{repo: repo}
}

// CreateShortURL generates a short code and stores the URL.
func (s *URLService) CreateShortURL(ctx context.Context, originalURL string) (*model.URL, error) {
	code, err := generateShortCode()
	if err != nil {
		return nil, err
	}
	return s.repo.Create(ctx, code, originalURL)
}

// GetByShortCode retrieves a URL by its short code.
func (s *URLService) GetByShortCode(ctx context.Context, shortCode string) (*model.URL, error) {
	return s.repo.FindByShortCode(ctx, shortCode)
}

// generateShortCode creates a random string of shortCodeLength characters.
func generateShortCode() (string, error) {
	result := make([]byte, shortCodeLength)
	for i := range result {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[n.Int64()]
	}
	return string(result), nil
}