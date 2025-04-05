package utils

import (
	"open-source-content-api/models"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

// SortEntriesByCreationTime sorts entries by creation time (newest first)
func SortEntriesByCreationTime(entries []models.EntryResult) {
	for i := 0; i < len(entries); i++ {
		for j := i + 1; j < len(entries); j++ {
			if entries[i].CreatedAt.Before(entries[j].CreatedAt) {
				entries[i], entries[j] = entries[j], entries[i]
			}
		}
	}
}

func FilterYoutubeEntries(entries []models.EntryResult) []models.EntryResult {
	var youtubeEntries []models.EntryResult

	for _, entry := range entries {
		// Check for both www.youtube.com and youtu.be formats
		if strings.Contains(entry.URL, "youtube.com") || strings.Contains(entry.URL, "youtu.be") {
			youtubeEntries = append(youtubeEntries, entry)
		}
	}

	return youtubeEntries

}

func FilterWebEntries(entries []models.EntryResult) []models.EntryResult {
	var webLinks []models.EntryResult

	for _, entry := range entries {
		// Check for both www.youtube.com and youtu.be formats
		if !strings.Contains(entry.URL, "youtube.com") && !strings.Contains(entry.URL, "youtu.be") {
			webLinks = append(webLinks, entry)
		}
	}

	return webLinks
}

// parsePaginationParams extracts and validates pagination parameters from the request
func ParsePaginationParams(c echo.Context, defaultPageSize int) models.PaginationParams {
	params := models.PaginationParams{
		PageSize: defaultPageSize,
		Offset:   0,
	}
	
	// Parse page size if provided
	if sizeParm := c.QueryParam("pageSize"); sizeParm != "" {
		if size, err := strconv.Atoi(sizeParm); err == nil && size > 0 {
			params.PageSize = size
			// Cap at reasonable maximum
			if params.PageSize > 100 {
				params.PageSize = 100
			}
		}
	}
	
	// Parse cursor (which is our offset)
	params.Cursor = c.QueryParam("cursor")
	if params.Cursor != "" {
		if offset, err := strconv.Atoi(params.Cursor); err == nil && offset >= 0 {
			params.Offset = offset
		}
	}
	
	return params
}

// applyPagination applies pagination parameters to the full list of entries
func ApplyPagination(entries []models.EntryResult, params models.PaginationParams) models.PaginatedResponse {
	totalEntries := len(entries)
	hasMore := params.Offset + params.PageSize < totalEntries
	
	response := models.PaginatedResponse{
		HasMore:    hasMore,
		NextCursor: "",
		Entries:    []models.EntryResult{},
	}
	
	// Add the entries for this page
	endIndex := min(params.Offset+params.PageSize, totalEntries)
	if params.Offset < totalEntries {
		response.Entries = entries[params.Offset:endIndex]
	}
	
	// Set the next cursor if there are more entries
	if hasMore {
		response.NextCursor = strconv.Itoa(params.Offset + params.PageSize)
	}
	
	return response
}

// min returns the smaller of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}