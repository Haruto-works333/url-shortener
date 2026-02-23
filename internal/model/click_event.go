package model

import "time"

// ClickEvent represents a row in BigQuery click_events table
type ClickEvent struct {
	ShortCode   string    `bigquery:"short_code"`
	OriginalURL string    `bigquery:"original_url"`
	IPAddress   string    `bigquery:"ip_address"`
	UserAgent   string    `bigquery:"user_agent"`
	Referer     string    `bigquery:"referer"`
	ClickedAt   time.Time `bigquery:"clicked_at"`
}
