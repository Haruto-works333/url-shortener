// internal/repository/click_log.go
package repository

import (
	"context"
	"database/sql"

	"github.com/Haruto-works333/url-shortener/internal/model"
)

// ClickLogRepository handles DB operations for click_logs
type ClickLogRepository struct {
	db *sql.DB
}

// NewClickLogRepository creates a new ClickLogRepository
func NewClickLogRepository(db *sql.DB) *ClickLogRepository {
	return &ClickLogRepository{db: db}
}

// Create inserts a new click log record
func (r *ClickLogRepository) Create(ctx context.Context, log *model.ClickLog) error {
	query := `
		INSERT INTO click_logs (url_id, ip_address, user_agent, referer)
		VALUES ($1, $2, $3, $4)
		RETURNING id, clicked_at`
	return r.db.QueryRowContext(ctx, query,
		log.URLID,
		log.IPAddress,
		log.UserAgent,
		log.Referer,
	).Scan(&log.ID, &log.ClickedAt)
}