package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
    NotionAPIKey     string
    NotionRootPageID string
    ServerPort       string
    DefaultPageSize  int
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