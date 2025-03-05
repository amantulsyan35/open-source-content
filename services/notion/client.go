package notion

import (
	"context"
	"fmt"
	"sync"

	"github.com/jomei/notionapi"

	"open-source-content-api/models"
)

// Service encapsulates Notion API interactions
type Service struct {
	client  *notionapi.Client
	rootID  string
}

// NewService creates a new Notion service
func NewService(apiKey, rootPageID string) *Service {
	return &Service{
		client: notionapi.NewClient(notionapi.Token(apiKey)),
		rootID: rootPageID,
	}
}

// FetchAllEntries retrieves all entries from all Notion databases
func (s *Service) FetchAllEntries(ctx context.Context) ([]models.EntryResult, error) {
	// Step 1: Get all child pages
	childPages, err := s.getAllChildPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get child pages: %w", err)
	}
	
	// Step 2: Get all databases from all child pages
	databaseIDs, err := s.getAllDatabaseIDs(ctx, childPages)
	if err != nil {
		return nil, fmt.Errorf("failed to get database IDs: %w", err)
	}
	
	// Step 3: Get all entries from all databases
	entries, err := s.getAllDatabaseEntries(ctx, databaseIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get database entries: %w", err)
	}
	
	return entries, nil
}

// SortEntriesByCreationTime sorts entries by creation time (newest first)
func (s *Service) SortEntriesByCreationTime(entries []models.EntryResult) {
	for i := 0; i < len(entries); i++ {
		for j := i + 1; j < len(entries); j++ {
			if entries[i].CreatedTime.Before(entries[j].CreatedTime) {
				entries[i], entries[j] = entries[j], entries[i]
			}
		}
	}
}

// getAllChildPages retrieves all child pages from the root page
func (s *Service) getAllChildPages(ctx context.Context) ([]string, error) {
	var allChildPages []string
	var startCursorStr string
	hasMorePages := true
	hasStartCursor := false
	
	for hasMorePages {
		pageBlocksRequest := &notionapi.Pagination{
			PageSize: 100, // Get a large batch to reduce API calls
		}
		
		if hasStartCursor {
			pageBlocksRequest.StartCursor = notionapi.Cursor(startCursorStr)
		}
		
		pageBlocks, err := s.client.Block.GetChildren(ctx, notionapi.BlockID(s.rootID), pageBlocksRequest)
		if err != nil {
			return nil, err
		}
		
		// Process the child pages
		for _, block := range pageBlocks.Results {
			if block.GetType() == "child_page" {
				if childPageBlock, ok := block.(*notionapi.ChildPageBlock); ok {
					allChildPages = append(allChildPages, string(childPageBlock.ID))
				}
			}
		}
		
		// Check if there are more pages to fetch
		if pageBlocks.HasMore {
			startCursorStr = string(pageBlocks.NextCursor)
			hasStartCursor = true
		} else {
			hasMorePages = false
		}
	}
	
	return allChildPages, nil
}

// getAllDatabaseIDs retrieves all database IDs from the child pages
func (s *Service) getAllDatabaseIDs(ctx context.Context, childPageIDs []string) ([]notionapi.DatabaseID, error) {
	var allDatabaseIDs []notionapi.DatabaseID
	var databaseIDsMutex sync.Mutex
	var wg sync.WaitGroup
	
	for _, pageID := range childPageIDs {
		wg.Add(1)
		go func(blockID string) {
			defer wg.Done()
			
			blocks, err := s.client.Block.GetChildren(ctx, notionapi.BlockID(blockID), nil)
			if err != nil {
				// Log error but continue with other pages
				fmt.Printf("Error fetching child blocks: %v\n", err)
				return
			}
			
			var dbIDs []notionapi.DatabaseID
			for _, block := range blocks.Results {
				if block.GetType() == "child_database" {
					if dbBlock, ok := block.(*notionapi.ChildDatabaseBlock); ok {
						dbIDs = append(dbIDs, notionapi.DatabaseID(dbBlock.ID))
					}
				}
			}
			
			// Safely add to the global list
			if len(dbIDs) > 0 {
				databaseIDsMutex.Lock()
				allDatabaseIDs = append(allDatabaseIDs, dbIDs...)
				databaseIDsMutex.Unlock()
			}
		}(pageID)
	}
	
	// Wait for all database IDs to be collected
	wg.Wait()
	
	return allDatabaseIDs, nil
}

// getAllDatabaseEntries retrieves all entries from the databases
func (s *Service) getAllDatabaseEntries(ctx context.Context, databaseIDs []notionapi.DatabaseID) ([]models.EntryResult, error) {
	var allEntries []models.EntryResult
	var allEntriesMutex sync.Mutex
	var wg sync.WaitGroup
	
	for _, dbID := range databaseIDs {
		wg.Add(1)
		go func(databaseID notionapi.DatabaseID) {
			defer wg.Done()
			
			// Query the entire database - we'll handle pagination in memory
			query := &notionapi.DatabaseQueryRequest{
				PageSize: 100, // Get a large batch to reduce API calls
			}
			
			resp, err := s.client.Database.Query(ctx, databaseID, query)
			if err != nil {
				// Log error but continue with other databases
				fmt.Printf("Error querying database: %v\n", err)
				return
			}
			
			var entries []models.EntryResult
			for _, page := range resp.Results {
				entry := models.EntryResult{}
				
				// Extract title from Name property
				if nameProperty, ok := page.Properties["Name"]; ok {
					if titleProperty, ok := nameProperty.(*notionapi.TitleProperty); ok && len(titleProperty.Title) > 0 {
						entry.Title = titleProperty.Title[0].PlainText
					}
				}
				
				// Extract URL from Link property
				if linkProperty, ok := page.Properties["Link"]; ok {
					if urlProperty, ok := linkProperty.(*notionapi.URLProperty); ok {
						entry.URL = urlProperty.URL
					}
				}
				
				entry.CreatedTime = page.CreatedTime
				
				// Only add valid entries
				// if entry.Title != "" && entry.URL != "" {
					entries = append(entries, entry)
				// }
			}
			
			// Safely add to the global list
			if len(entries) > 0 {
				allEntriesMutex.Lock()
				allEntries = append(allEntries, entries...)
				allEntriesMutex.Unlock()
			}
		}(dbID)
	}
	
	// Wait for all entries to be collected
	wg.Wait()
	
	return allEntries, nil
}