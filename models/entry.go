package models

import (
	"time"
)

// EntryResult represents a formatted result from the Notion database
type EntryResult struct {
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	CreatedTime time.Time `json:"createdTime"`
}

// PaginatedResponse wraps the results with pagination metadata
type PaginatedResponse struct {
	Entries    []EntryResult `json:"entries"`
	NextCursor string        `json:"nextCursor,omitempty"`
	HasMore    bool          `json:"hasMore"`
}

// PaginationParams holds query parameters for pagination
type PaginationParams struct {
	Cursor   string
	PageSize int
	Offset   int
}