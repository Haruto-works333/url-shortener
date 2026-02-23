package service

import (
	"context"
	"time"

	"github.com/Haruto-works333/url-shortener/internal/model"
	"github.com/Haruto-works333/url-shortener/internal/repository"
)

// ClickLogService handles business logic for click logging
type ClickLogService struct {
	clickLogRepo *repository.ClickLogRepository
	urlRepo      *repository.URLRepository
	bqRepo       *repository.ClickEventBQRepository
}

// NewClickLogService creates a new ClickLogService
func NewClickLogService(clickLogRepo *repository.ClickLogRepository, urlRepo *repository.URLRepository, bqRepo *repository.ClickEventBQRepository) *ClickLogService {
	return &ClickLogService{
		clickLogRepo: clickLogRepo,
		urlRepo:      urlRepo,
		bqRepo:       bqRepo,
	}
}

// RedirectAndLog finds the original URL and records the click event
func (s *ClickLogService) RedirectAndLog(ctx context.Context, shortCode, ipAddress, userAgent, referer string) (string, error) {
	// 1. Find the original URL
	url, err := s.urlRepo.FindByShortCode(ctx, shortCode)
	if err != nil {
		return "", err
	}

	// 2. Record the click log to Cloud SQL (source of truth)
	clickLog := &model.ClickLog{
		URLID:     url.ID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Referer:   referer,
	}
	if err := s.clickLogRepo.Create(ctx, clickLog); err != nil {
		return "", err
	}

	// 3. Send to BigQuery asynchronously (best-effort)
	if s.bqRepo != nil {
		event := &model.ClickEvent{
			ShortCode:   shortCode,
			OriginalURL: url.OriginalURL,
			IPAddress:   ipAddress,
			UserAgent:   userAgent,
			Referer:     referer,
			ClickedAt:   time.Now(),
		}
		s.bqRepo.InsertAsync(ctx, event)
	}

	return url.OriginalURL, nil
}
