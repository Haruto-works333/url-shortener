// internal/handler/redirect.go
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Haruto-works333/url-shortener/internal/service"
)

// RedirectHandler handles redirect requests
type RedirectHandler struct {
	clickLogService *service.ClickLogService
}

// NewRedirectHandler creates a new RedirectHandler
func NewRedirectHandler(clickLogService *service.ClickLogService) *RedirectHandler {
	return &RedirectHandler{clickLogService: clickLogService}
}

// Redirect looks up the short code, logs the click, and redirects to the original URL
func (h *RedirectHandler) Redirect(c *gin.Context) {
	shortCode := c.Param("code")

	originalURL, err := h.clickLogService.RedirectAndLog(
		c.Request.Context(),
		shortCode,
		c.ClientIP(),
		c.Request.UserAgent(),
		c.Request.Referer(),
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "short code not found"})
		return
	}

	c.Redirect(http.StatusFound, originalURL)
}