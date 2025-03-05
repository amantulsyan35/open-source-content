package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"

	"open-source-content-api/config"
	"open-source-content-api/internal/api"
)



func main() {
	// Load configuration
	cfg := config.Load()
	
	// Create a new Echo instance
	e := echo.New()
	
	// Add middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.RateLimiterWithConfig(
		middleware.RateLimiterConfig{
			Skipper: middleware.DefaultSkipper,
			Store: middleware.NewRateLimiterMemoryStoreWithConfig(
				middleware.RateLimiterMemoryStoreConfig{
					Rate:      rate.Limit(cfg.RateLimiterRate),
					Burst:     cfg.RateLimiterBurst,
					ExpiresIn: cfg.RateLimiterExpiry,
				},
			),
			IdentifierExtractor: func(ctx echo.Context) (string, error) {
				id := ctx.RealIP()
				return id, nil
			},
			ErrorHandler: func(context echo.Context, err error) error {
				return context.JSON(http.StatusForbidden, nil)
			},
			DenyHandler: func(context echo.Context, identifier string, err error) error {
				return context.JSON(http.StatusTooManyRequests, nil)
			},
		},
	))
	

	
	// Setup routes
	api.SetupRoutes(e, cfg)
	
	// Start server
	log.Printf("Starting server on port %s", cfg.ServerPort)
	e.Logger.Fatal(e.Start(":" + cfg.ServerPort))
}