package repository

import (
	"context"
	"database/sql"

	"github.com/Haruto-works333/url-shortener/internal/model"
)

// URLRepository handles database operations for URLs.
type URLRepository struct {
	db *sql.DB
}

// NewURLRepository creates a new URLRepository.
func NewURLRepository(db *sql.DB) *URLRepository {
	return &URLRepository{db: db}
}

// Create inserts a new URL record and returns it.
func (r *URLRepository) Create(ctx context.Context, shortCode, originalURL string) (*model.URL, error) {
	query := `
		INSERT INTO urls (short_code, original_url)
		VALUES ($1, $2)
		RETURNING id, short_code, original_url, created_at, updated_at`

	var u model.URL
	err := r.db.QueryRowContext(ctx, query, shortCode, originalURL).
		Scan(&u.ID, &u.ShortCode, &u.OriginalURL, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// FindByShortCode retrieves a URL by its short code.
func (r *URLRepository) FindByShortCode(ctx context.Context, shortCode string) (*model.URL, error) {
	query := `SELECT id, short_code, original_url, created_at, updated_at FROM urls WHERE short_code = $1`

	var u model.URL
	err := r.db.QueryRowContext(ctx, query, shortCode).
		Scan(&u.ID, &u.ShortCode, &u.OriginalURL, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
