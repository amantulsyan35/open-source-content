package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"open-source-content-api/config"
	"open-source-content-api/models"
	"open-source-content-api/services/notion"
)

// NotionHandler handles HTTP requests related to Notion data
type NotionHandler struct {
	service  *notion.Service
	config   *config.Config
}

// NewNotionHandler creates a new NotionHandler instance
func NewNotionHandler(service *notion.Service, cfg *config.Config) *NotionHandler {
	return &NotionHandler{
		service: service,
		config:  cfg,
	}
}

// GetDatabase handles the request to fetch and process Notion database entries
func (h *NotionHandler) GetDatabase(c echo.Context) error {
	startTime := time.Now()

	// Parse pagination parameters
	params := parsePaginationParams(c, h.config.DefaultPageSize)
	
	// Fetch all entries from the Notion database
	entries, err := h.service.FetchAllEntries(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	// Sort entries by creation time (newest first)
	h.service.SortEntriesByCreationTime(entries)
	
	// Apply pagination
	response := applyPagination(entries, params)
	
	executionTime := time.Since(startTime)
	c.Logger().Infof("Total execution time: %v", executionTime)

	return c.JSON(http.StatusOK, response)
}



func (h *NotionHandler) GetIntroduction(c echo.Context) error {
	introduction := models.Introduction{
		Title: "Open Source Content API",
		Description: "For the past three years, I've been tracking the content I consume. It began as a simple behavioral experiment aimed at predicting how my consumption shapes my thinking and problem-solving approaches.\nOver time, it evolved into a curious pursuit and a core thesis on how I operate, Personal Databases.\nThis open API consolidates all the content I consume, making it embeddable and paving the way for innovative applications powered by this dataset—especially in an era dominated by LLMs.",
		Thesis: "https://www.amantulsyan.com/personal-databases",
		GithubLink: "https://github.com/amantulsyan35/open-source-content/tree/main",
	}
	
	return c.JSON(http.StatusOK, introduction)
}

// parsePaginationParams extracts and validates pagination parameters from the request
func parsePaginationParams(c echo.Context, defaultPageSize int) models.PaginationParams {
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
func applyPagination(entries []models.EntryResult, params models.PaginationParams) models.PaginatedResponse {
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