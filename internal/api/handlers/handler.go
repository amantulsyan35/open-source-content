package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"open-source-content-api/config"
	"open-source-content-api/models"
	"open-source-content-api/services/notion"
	"open-source-content-api/utils"
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
	params := utils.ParsePaginationParams(c, h.config.DefaultPageSize)
	
	// Fetch all entries from the Notion database
	entries, err := h.service.FetchAllEntries(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	utils.SortEntriesByCreationTime(entries)
	
	// Apply pagination
	response := utils.ApplyPagination(entries, params)
	
	executionTime := time.Since(startTime)
	c.Logger().Infof("Total execution time: %v", executionTime)

	return c.JSON(http.StatusOK, response)
}

func (h *NotionHandler) GetYoutubeVideos(c echo.Context) error {
		// Fetch all entries from the Notion database
		entries, err := h.service.FetchAllEntries(c.Request().Context())
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		utils.SortEntriesByCreationTime(entries)

		youtubeEntries := utils.FilterYoutubeEntries(entries)

		response := models.YoutubeResponse{
			Data: youtubeEntries,
		}

		return c.JSON(http.StatusOK, response)
}

func (h *NotionHandler) GetWebLinks(c echo.Context) error {
		// Fetch all entries from the Notion database
		entries, err := h.service.FetchAllEntries(c.Request().Context())
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		utils.SortEntriesByCreationTime(entries)

		webEntries := utils.FilterWebEntries(entries)

		response := models.WebResponse{
			Data: webEntries,
		}

		return c.JSON(http.StatusOK, response)
}


func (h *NotionHandler) GetIntroduction(c echo.Context) error {
	log.Println("GetIntroduction handler called")
	introduction := models.Introduction{
		Title: "Open Source Content API",
		Description: "For the past three years, I've been tracking the content I consume. It began as a simple behavioral experiment aimed at predicting how my consumption shapes my thinking and problem-solving approaches.\nOver time, it evolved into a curious pursuit and a core thesis on how I operate, Personal Databases.\nThis open API consolidates all the content I consume, making it embeddable and paving the way for innovative applications powered by this dataset—especially in an era dominated by LLMs.",
		Thesis: "https://www.amantulsyan.com/personal-databases",
		GithubLink: "https://github.com/amantulsyan35/open-source-content/tree/main",
	}

	log.Printf("Returning response with status %d", http.StatusOK)
	
	return c.JSON(http.StatusOK, introduction)
}

