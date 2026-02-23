// internal/model/click_log.go
package model

import "time"

// ClickLog represents a single redirect event
type ClickLog struct {
	ID        int64     `json:"id"`
	URLID     int64     `json:"url_id"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	Referer   string    `json:"referer"`
	ClickedAt time.Time `json:"clicked_at"`
}