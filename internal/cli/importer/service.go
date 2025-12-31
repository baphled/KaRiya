package importer

import (
	"context"
	"fmt"

	"github.com/baphled/kariya/internal/domain/career"
	repo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	burst_fact "github.com/baphled/kariya/internal/service/career/burst_fact"
)

// ImportResult represents the result of an import operation
type ImportResult struct {
	TotalRows           int
	SuccessCount        int
	SkippedCount        int
	FailedCount         int
	CreatedEvents       []*career.CareerEvent
	FailedRows          []*ParsedRow
	BurstSuggestions    []burst_fact.BurstSuggestion // Burst suggestions detected from imported events
	ExtractedFactsCount int                           // Number of facts extracted from imported events
	FactsByEventID      map[string][]*career.Fact     // Facts keyed by source event ID
	FactsByCompetency   map[string]int                // Count of facts by competency category
}

// ImportService handles the import workflow
type ImportService struct {
	careerService *careerservice.Service
	parser        *CSVParser
}

// NewImportService creates a new import service
func NewImportService(careerService *careerservice.Service) *ImportService {
	return &ImportService{
		careerService: careerService,
	}
}

// PrepareImport parses CSV and returns parsed rows for review
func (is *ImportService) PrepareImport(ctx context.Context, reader interface{}) ([]*ParsedRow, error) {
	// Get existing events to check for duplicates
	filters := &repo.ListFilters{}
	existingEvents, err := is.careerService.ListEvents(ctx, *filters)
	if err != nil {
		return nil, fmt.Errorf("failed to load existing events: %w", err)
	}

	// Create parser with existing events
	is.parser = NewCSVParserWithMapping(existingEvents)

	// Try to convert reader to io.Reader
	ioReader, ok := reader.(interface{ Read([]byte) (int, error) })
	if !ok {
		return nil, fmt.Errorf("invalid reader type")
	}

	// Parse CSV
	parsedRows, err := is.parser.Parse(ioReader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV: %w", err)
	}

	return parsedRows, nil
}

// ImportRows imports the parsed rows into the database
func (is *ImportService) ImportRows(ctx context.Context, parsedRows []*ParsedRow, selectedRows []int) (*ImportResult, error) {
	result := &ImportResult{
		TotalRows:         len(selectedRows),
		CreatedEvents:     []*career.CareerEvent{},
		FailedRows:        []*ParsedRow{},
		FactsByEventID:    make(map[string][]*career.Fact),
		FactsByCompetency: make(map[string]int),
	}

	// Create a set of selected row numbers for quick lookup
	selectedSet := make(map[int]bool)
	for _, rowNum := range selectedRows {
		selectedSet[rowNum] = true
	}

	// Import each selected row
	for _, rowNum := range selectedRows {
		// Find the parsed row
		var parsedRow *ParsedRow
		for _, pr := range parsedRows {
			if pr.RowNumber == rowNum {
				parsedRow = pr
				break
			}
		}

		if parsedRow == nil {
			result.FailedCount++
			continue
		}

		// Skip invalid or duplicate rows
		if !parsedRow.IsValid || parsedRow.IsDuplicate {
			result.SkippedCount++
			continue
		}

		// Create the event
		event := parsedRow.Event
		err := is.careerService.CaptureEvent(ctx, event, careerservice.ManualEntry)
		if err != nil {
			result.FailedCount++
			parsedRow.ValidationErrors = append(parsedRow.ValidationErrors, fmt.Sprintf("Failed to create event: %v", err))
			result.FailedRows = append(result.FailedRows, parsedRow)
			continue
		}

		result.SuccessCount++
		result.CreatedEvents = append(result.CreatedEvents, event)
	}

	// Detect bursts from new events (Task 2.0: Post-import burst detection)
	if len(result.CreatedEvents) > 0 {
		eventIDs := make([]string, len(result.CreatedEvents))
		for i, event := range result.CreatedEvents {
			eventIDs[i] = event.ID
		}

		suggestions, err := is.careerService.SuggestBursts(ctx, eventIDs)
		if err != nil {
			// Log warning but continue - burst detection is an optional enhancement
			fmt.Printf("Warning: Failed to suggest bursts: %v\n", err)
		} else {
			result.BurstSuggestions = suggestions
			if len(suggestions) > 0 {
				fmt.Printf("Detected %d burst suggestions from %d events\n",
					len(suggestions), len(result.CreatedEvents))
			}
		}
	}

	// Extract facts from new events (Task 3.0: Post-import fact extraction)
	if len(result.CreatedEvents) > 0 {
		for _, event := range result.CreatedEvents {
			// Extract facts from this event
			facts, err := is.careerService.ExtractFactsFromEvent(ctx, event)
			if err != nil {
				// Log warning but continue - fact extraction is an optional enhancement
				fmt.Printf("Warning: Failed to extract facts from event %s: %v\n", event.ID, err)
				continue
			}

			// Persist each extracted fact
			for _, fact := range facts {
				// Set source event ID
				fact.SourceEventID = event.ID

				// Save fact to repository
				if err := is.careerService.SaveFact(ctx, &fact); err != nil {
					fmt.Printf("Warning: Failed to save fact: %v\n", err)
					continue
				}

				// Track the fact
				result.ExtractedFactsCount++
				result.FactsByEventID[event.ID] = append(result.FactsByEventID[event.ID], &fact)

				// Count by competency
				for _, competency := range fact.CompetencyCategories {
					result.FactsByCompetency[competency]++
				}
			}
		}

		// Log summary
		if result.ExtractedFactsCount > 0 {
			fmt.Printf("Extracted %d facts from %d events\n",
				result.ExtractedFactsCount, len(result.CreatedEvents))
			if len(result.FactsByCompetency) > 0 {
				fmt.Printf("Competency breakdown: ")
				for competency, count := range result.FactsByCompetency {
					fmt.Printf("%s: %d, ", competency, count)
				}
				fmt.Printf("\n")
			}
		}
	}

	return result, nil
}

// GetImportSummary returns a summary of the import preparation
func (is *ImportService) GetImportSummary(parsedRows []*ParsedRow) map[string]int {
	summary := map[string]int{
		"total":     len(parsedRows),
		"valid":     0,
		"invalid":   0,
		"duplicate": 0,
	}

	for _, row := range parsedRows {
		if row.IsDuplicate {
			summary["duplicate"]++
		} else if row.IsValid {
			summary["valid"]++
		} else {
			summary["invalid"]++
		}
	}

	return summary
}
