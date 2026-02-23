package repository

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"

	"github.com/Haruto-works333/url-shortener/internal/model"
)

// StatsRepository handles BigQuery reads for analytics
type StatsRepository struct {
	client *bigquery.Client
}

// NewStatsRepository creates a new StatsRepository
func NewStatsRepository(client *bigquery.Client) *StatsRepository {
	return &StatsRepository{client: client}
}

// GetClickStats returns total and daily click counts for a short code
func (r *StatsRepository) GetClickStats(ctx context.Context, shortCode string, from, to time.Time) (*model.URLStats, error) {
	query := r.client.Query(`
		SELECT
			FORMAT_TIMESTAMP('%Y-%m-%d', clicked_at) AS date,
			COUNT(*) AS clicks
		FROM ` + "`haru-url-shortener.analytics.click_events`" + `
		WHERE short_code = @short_code
		  AND clicked_at >= @from
		  AND clicked_at < @to
		GROUP BY date
		ORDER BY date
	`)
	query.Parameters = []bigquery.QueryParameter{
		{Name: "short_code", Value: shortCode},
		{Name: "from", Value: from},
		{Name: "to", Value: to},
	}

	it, err := query.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("bigquery query failed: %w", err)
	}

	stats := &model.URLStats{
		ShortCode:   shortCode,
		DailyClicks: []model.DailyClick{},
	}

	var total int64
	for {
		var row struct {
			Date   string `bigquery:"date"`
			Clicks int64  `bigquery:"clicks"`
		}
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("bigquery row read failed: %w", err)
		}
		stats.DailyClicks = append(stats.DailyClicks, model.DailyClick{
			Date:   row.Date,
			Clicks: row.Clicks,
		})
		total += row.Clicks
	}
	stats.TotalClicks = total

	return stats, nil
}

// GetRefererStats returns click counts grouped by referer
func (r *StatsRepository) GetRefererStats(ctx context.Context, shortCode string, from, to time.Time) ([]model.RefererStat, error) {
	query := r.client.Query(`
		SELECT
			IFNULL(NULLIF(referer, ''), 'direct') AS referer,
			COUNT(*) AS clicks
		FROM ` + "`haru-url-shortener.analytics.click_events`" + `
		WHERE short_code = @short_code
		  AND clicked_at >= @from
		  AND clicked_at < @to
		GROUP BY referer
		ORDER BY clicks DESC
	`)
	query.Parameters = []bigquery.QueryParameter{
		{Name: "short_code", Value: shortCode},
		{Name: "from", Value: from},
		{Name: "to", Value: to},
	}

	it, err := query.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("bigquery query failed: %w", err)
	}

	var results []model.RefererStat
	for {
		var row struct {
			Referer string `bigquery:"referer"`
			Clicks  int64  `bigquery:"clicks"`
		}
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("bigquery row read failed: %w", err)
		}
		results = append(results, model.RefererStat{
			Referer: row.Referer,
			Clicks:  row.Clicks,
		})
	}

	return results, nil
}
