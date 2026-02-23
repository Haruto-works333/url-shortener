package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"github.com/Haruto-works333/url-shortener/internal/handler"
	"github.com/Haruto-works333/url-shortener/internal/repository"
	"github.com/Haruto-works333/url-shortener/internal/service"
)

func main() {
	// Connect to PostgreSQL
	dsn := "host=localhost port=5432 user=app password=password dbname=url_shortener sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("failed to ping database:", err)
	}
	log.Println("connected to database")

	// Initialize layers
	repo := repository.NewURLRepository(db)
	svc := service.NewURLService(repo)
	urlHandler := handler.NewURLHandler(svc)

	clickLogRepo := repository.NewClickLogRepository(db)
	clickLogSvc := service.NewClickLogService(clickLogRepo, repo)
	redirectHandler := handler.NewRedirectHandler(clickLogSvc)

	// Setup router
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	{
		api.POST("/urls", urlHandler.Create)
		api.GET("/urls/:code", urlHandler.GetByShortCode)
	}

	// Redirect: GET /:code -> original URL
	r.GET("/:code", redirectHandler.Redirect)

	r.Run(":8080")
}
