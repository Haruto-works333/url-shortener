package model

// URLStats represents click statistics for a short URL
type URLStats struct {
	ShortCode   string       `json:"short_code"`
	TotalClicks int64        `json:"total_clicks"`
	DailyClicks []DailyClick `json:"daily_clicks"`
}

// DailyClick represents click count for a single day
type DailyClick struct {
	Date   string `json:"date"`
	Clicks int64  `json:"clicks"`
}

// RefererStat represents click count grouped by referer
type RefererStat struct {
	Referer string `json:"referer"`
	Clicks  int64  `json:"clicks"`
}
