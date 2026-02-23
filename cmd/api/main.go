package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/bigquery"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"github.com/Haruto-works333/url-shortener/internal/handler"
	"github.com/Haruto-works333/url-shortener/internal/repository"
	"github.com/Haruto-works333/url-shortener/internal/service"
)

func main() {
	// Connect to PostgreSQL
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost port=5432 user=app password=password dbname=url_shortener sslmode=disable"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("failed to ping database:", err)
	}
	log.Println("connected to database")

	// Initialize BigQuery client (optional: skip if GCP_PROJECT is not set)
	var bqRepo *repository.ClickEventBQRepository
	var statsRepo *repository.StatsRepository
	if project := os.Getenv("GCP_PROJECT"); project != "" {
		bqClient, err := bigquery.NewClient(context.Background(), project)
		if err != nil {
			log.Printf("WARNING: BigQuery client init failed: %v (continuing without BQ)", err)
		} else {
			defer bqClient.Close()
			bqRepo = repository.NewClickEventBQRepository(bqClient)
			statsRepo = repository.NewStatsRepository(bqClient)
			log.Println("connected to BigQuery")
		}
	} else {
		log.Println("GCP_PROJECT not set, BigQuery disabled")
	}

	// Initialize layers
	repo := repository.NewURLRepository(db)
	svc := service.NewURLService(repo)
	urlHandler := handler.NewURLHandler(svc)

	clickLogRepo := repository.NewClickLogRepository(db)
	clickLogSvc := service.NewClickLogService(clickLogRepo, repo, bqRepo)
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

	// Analytics endpoints (only if BigQuery is available)
	if statsRepo != nil {
		statsHandler := handler.NewStatsHandler(statsRepo)
		api.GET("/urls/:code/stats", statsHandler.GetStats)
		api.GET("/urls/:code/referers", statsHandler.GetReferers)
	}

	// Redirect: GET /:code -> original URL
	r.GET("/:code", redirectHandler.Redirect)

	// Use PORT env var (Cloud Run sets this automatically)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
