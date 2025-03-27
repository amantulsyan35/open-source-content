package models

import (
	"time"
)

// EntryResult represents a formatted result from the Notion database
type EntryResult struct {
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	CreatedAt time.Time `json:"createdAt"`
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

type Introduction struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Thesis      string `json:"thesis"`
	GithubLink  string `json:"github_link"`
}

type YoutubeResponse struct {
	Data []string `json:"data"`
}
type WebResponse struct {
	Data []string `json:"data"`
}