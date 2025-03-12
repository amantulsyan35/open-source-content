package api

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"open-source-content-api/config"
	"open-source-content-api/internal/api/handlers"
	"open-source-content-api/services/notion"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(e *echo.Echo, cfg *config.Config) {
	// Initialize the Notion service
	notionService := notion.NewService(cfg.NotionAPIKey, cfg.NotionRootPageID)
	
	// Initialize handlers
	notionHandler := handlers.NewNotionHandler(notionService, cfg)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "Test route works!")
	})
		

	// API v1 group
	v1 := e.Group("/v1")
	
	// Routes
	v1.GET("", notionHandler.GetDatabase)
}