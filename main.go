package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

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
	
	// Setup routes
	api.SetupRoutes(e, cfg)
	
	// Start server
	log.Printf("Starting server on port %s", cfg.ServerPort)
	e.Logger.Fatal(e.Start(":" + cfg.ServerPort))
}