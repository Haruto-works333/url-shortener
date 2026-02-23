package model

import "time"

// URL represents a shortened URL record in the database.
type URL struct {
	ID          int64     `json:"id"`
	ShortCode   string    `json:"short_code"`
	OriginalURL string    `json:"original_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateURLRequest is the payload for creating a new short URL.
type CreateURLRequest struct {
	URL string `json:"url" binding:"required,url"`
}

// CreateURLResponse is returned after successfully creating a short URL.
type CreateURLResponse struct {
	ShortCode   string `json:"short_code"`
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}
