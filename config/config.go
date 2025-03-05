package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
    NotionAPIKey     string
    NotionRootPageID string
    ServerPort       string
    DefaultPageSize  int

    // Rate Limiter Configuration
	RateLimiterRate     float64      // Requests per second
	RateLimiterBurst    int          // Maximum burst size
	RateLimiterExpiry   time.Duration // How long to keep rate limiter state
}



// Load returns a new Config struct populated from environment variables
func Load() *Config {
    // Load .env file if it exists
    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found or error loading it. Using environment variables.")
    }

    return &Config{
        NotionAPIKey:     getEnv("NOTION_API_KEY", ""),
        NotionRootPageID: getEnv("NOTION_ROOT_PAGE_ID", ""),
        ServerPort:       getEnv("SERVER_PORT", "1323"),
        DefaultPageSize:  20,

 	    // Rate Limiter configuration aligned with Notion API limits
		// Notion allows 3 req/sec average with some bursts allowed
		RateLimiterRate:     3.0,  // 3 requests per second to match Notion's limit
		RateLimiterBurst:    9,    // Allow small bursts (3 seconds worth of requests)
		RateLimiterExpiry:   1 * time.Minute, // Reset rate limiter state after 1 minute
    }
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}