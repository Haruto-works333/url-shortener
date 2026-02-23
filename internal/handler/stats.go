package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Haruto-works333/url-shortener/internal/repository"
)

// StatsHandler handles HTTP requests for analytics
type StatsHandler struct {
	statsRepo *repository.StatsRepository
}

// NewStatsHandler creates a new StatsHandler
func NewStatsHandler(statsRepo *repository.StatsRepository) *StatsHandler {
	return &StatsHandler{statsRepo: statsRepo}
}

// GetStats handles GET /api/urls/:code/stats
func (h *StatsHandler) GetStats(c *gin.Context) {
	code := c.Param("code")
	from, to := parseTimeRange(c)

	stats, err := h.statsRepo.GetClickStats(c.Request.Context(), code, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetReferers handles GET /api/urls/:code/referers
func (h *StatsHandler) GetReferers(c *gin.Context) {
	code := c.Param("code")
	from, to := parseTimeRange(c)

	referers, err := h.statsRepo.GetRefererStats(c.Request.Context(), code, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get referer stats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"short_code": code,
		"referers":   referers,
	})
}

// parseTimeRange extracts from/to query parameters with defaults
func parseTimeRange(c *gin.Context) (time.Time, time.Time) {
	now := time.Now()
	from := now.AddDate(0, 0, -30) // default: last 30 days
	to := now.AddDate(0, 0, 1)     // default: tomorrow (to include today)

	if v := c.Query("from"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			from = t
		}
	}
	if v := c.Query("to"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			to = t.AddDate(0, 0, 1) // make "to" inclusive
		}
	}

	return from, to
}
