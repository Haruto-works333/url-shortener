package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Haruto-works333/url-shortener/internal/model"
	"github.com/Haruto-works333/url-shortener/internal/service"
)

// URLHandler handles HTTP requests for URL operations.
type URLHandler struct {
	service *service.URLService
}

// NewURLHandler creates a new URLHandler.
func NewURLHandler(service *service.URLService) *URLHandler {
	return &URLHandler{service: service}
}

// Create handles POST /api/urls
func (h *URLHandler) Create(c *gin.Context) {
	var req model.CreateURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	url, err := h.service.CreateShortURL(c.Request.Context(), req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create short url"})
		return
	}

	c.JSON(http.StatusCreated, model.CreateURLResponse{
		ShortCode:   url.ShortCode,
		OriginalURL: url.OriginalURL,
		ShortURL:    fmt.Sprintf("http://localhost:8080/%s", url.ShortCode),
	})
}

// GetByShortCode handles GET /api/urls/:code
func (h *URLHandler) GetByShortCode(c *gin.Context) {
	code := c.Param("code")

	url, err := h.service.GetByShortCode(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "url not found"})
		return
	}

	c.JSON(http.StatusOK, url)
}