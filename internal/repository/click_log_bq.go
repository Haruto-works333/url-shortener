package repository

import (
	"context"
	"log"

	"cloud.google.com/go/bigquery"
	"github.com/Haruto-works333/url-shortener/internal/model"
)

// ClickEventBQRepository handles BigQuery writes for click events
type ClickEventBQRepository struct {
	inserter *bigquery.Inserter
}

// NewClickEventBQRepository creates a new BigQuery click event repository
func NewClickEventBQRepository(client *bigquery.Client) *ClickEventBQRepository {
	table := client.Dataset("analytics").Table("click_events")
	return &ClickEventBQRepository{
		inserter: table.Inserter(),
	}
}

// InsertAsync sends a click event to BigQuery in a separate goroutine.
// Errors are logged but do not affect the caller.
func (r *ClickEventBQRepository) InsertAsync(ctx context.Context, event *model.ClickEvent) {
	go func() {
		if err := r.inserter.Put(context.Background(), event); err != nil {
			log.Printf("ERROR: BigQuery insert failed: %v", err)
		}
	}()
}
