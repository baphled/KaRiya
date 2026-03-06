// Package facts provides Cobra CLI commands for fact extraction and listing.
package facts

import (
	"context"
	"fmt"
	"io"

	"github.com/baphled/kariya/internal/cli/cmd/cliutil"
	domain "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
)

// FactExtractionService defines the interface for operations needed by ExtractFacts.
type FactExtractionService interface {
	ListEvents(ctx context.Context, filters career.EventListFilters) ([]*domain.Event, error)
	ExtractFactsFromEvent(ctx context.Context, event *domain.Event) ([]domain.Fact, error)
	SaveFact(ctx context.Context, fact *domain.Fact) error
}

// CareerServiceAdapter adapts CareerService to FactExtractionService interface.
type CareerServiceAdapter struct {
	svc *careerservice.Service
}

// NewCareerServiceAdapter creates a new adapter.
//
// Expected:
//   - svc: valid CareerService instance
//
// Returns:
//   - CareerServiceAdapter implementing FactExtractionService
//
// Side effects: None.
func NewCareerServiceAdapter(svc *careerservice.Service) *CareerServiceAdapter {
	return &CareerServiceAdapter{svc: svc}
}

// ListEvents delegates event listing to the career service.
//
// Expected:
//   - ctx: context for the operation
//   - filters: criteria for filtering events
//
// Returns:
//   - Slice of matching events, or error
//
// Side effects:
//   - Reads from database via career service
func (a *CareerServiceAdapter) ListEvents(ctx context.Context, filters career.EventListFilters) ([]*domain.Event, error) {
	return a.svc.ListEvents(ctx, filters)
}

// ExtractFactsFromEvent delegates fact extraction to the career service.
//
// Expected:
//   - ctx: context for the operation
//   - event: event to extract facts from
//
// Returns:
//   - Extracted facts, or error
//
// Side effects:
//   - May read from database via career service
func (a *CareerServiceAdapter) ExtractFactsFromEvent(ctx context.Context, event *domain.Event) ([]domain.Fact, error) {
	return a.svc.ExtractFactsFromEvent(ctx, event)
}

// SaveFact delegates fact persistence to the career service.
//
// Expected:
//   - ctx: context for the operation
//   - fact: fact to save
//
// Returns:
//   - Error if save fails
//
// Side effects:
//   - Writes to database via career service
func (a *CareerServiceAdapter) SaveFact(ctx context.Context, fact *domain.Fact) error {
	return a.svc.SaveFact(ctx, fact)
}

// ExtractFacts extracts facts from all events.
//
// Expected:
//   - svc: Initialized career Service instance
//   - out: io.Writer for output messages
//   - errOut: io.Writer for error output
//   - opts: Optional Bubble Tea program options for testing
//
// Returns:
//   - 0 on success
//   - 1 on error (retrieval, extraction, or save failure)
//
// Side effects:
//   - Reads events from database via svc
//   - Writes facts to database via svc
//   - Writes output to writers
func ExtractFacts(svc *careerservice.Service, out io.Writer, errOut io.Writer, opts ...tea.ProgramOption) int {
	adapter := NewCareerServiceAdapter(svc)
	return ExtractFactsWithService(adapter, out, errOut, opts...)
}

// ExtractFactsWithService extracts facts using the FactExtractionService interface.
//
// Expected:
//   - svc: FactExtractionService implementation
//   - out: io.Writer for output messages
//   - errOut: io.Writer for error output
//   - opts: Optional Bubble Tea program options for testing
//
// Returns:
//   - 0 on success
//   - 1 on error (retrieval, extraction, or save failure)
//
// Side effects:
//   - Reads events from service
//   - Writes facts to service
//   - Writes output to writers
func ExtractFactsWithService(svc FactExtractionService, out io.Writer, _ io.Writer, opts ...tea.ProgramOption) int {
	ctx := context.Background()

	events, err := svc.ListEvents(ctx, career.EventListFilters{Limit: 10000})
	if err != nil {
		cliutil.PrintError(fmt.Sprintf("Error retrieving events: %v", err))
		return 1
	}

	if len(events) == 0 {
		cliutil.PrintInfo("No events found in database.")
		return 0
	}

	factCount, competencyCount, err := extractAndSaveFactsWithService(ctx, svc, events, opts...)
	if err != nil {
		cliutil.PrintError(fmt.Sprintf("Error extracting facts: %v", err))
		return 1
	}

	if factCount == 0 {
		cliutil.PrintInfo("No facts extracted.")
		return 0
	}

	DisplayExtractionResults(out, factCount, len(events), competencyCount)
	cliutil.PrintSuccess("Fact extraction complete!")
	return 0
}

func extractAndSaveFactsWithService(
	ctx context.Context,
	svc FactExtractionService,
	events []*domain.Event,
	opts ...tea.ProgramOption,
) (int, map[string]int, error) {
	factCount := 0
	competencyCount := make(map[string]int)

	err := cliutil.RunWithProgress(fmt.Sprintf("Extracting facts from %d events...", len(events)), len(events), func(update func(int)) error {
		for i := range events {
			processedFacts := processEventFactsWithService(ctx, svc, events[i])
			factCount += processedFacts
			updateCompetencyCountsWithService(ctx, svc, events[i], competencyCount)
			update(i + 1)
		}
		return nil
	}, opts...)

	return factCount, competencyCount, err
}

func processEventFactsWithService(ctx context.Context, svc FactExtractionService, event *domain.Event) int {
	facts, err := svc.ExtractFactsFromEvent(ctx, event)
	if err != nil {
		return 0
	}

	savedCount := 0
	for j := range facts {
		if err := svc.SaveFact(ctx, &facts[j]); err != nil {
			continue
		}
		savedCount++
	}
	return savedCount
}

func updateCompetencyCountsWithService(
	ctx context.Context,
	svc FactExtractionService,
	event *domain.Event,
	competencyCount map[string]int,
) {
	facts, err := svc.ExtractFactsFromEvent(ctx, event)
	if err != nil {
		return
	}

	for i := range facts {
		if len(facts[i].CompetencyCategories) > 0 {
			for _, comp := range facts[i].CompetencyCategories {
				competencyCount[comp]++
			}
		}
	}
}

// DisplayExtractionResults writes formatted extraction results to the output.
//
// Expected:
//   - out: io.Writer to write to
//   - factCount: number of facts extracted
//   - eventCount: number of events processed
//   - competencyCount: map of competency categories and their counts
//
// Returns: None
//
// Side effects:
//   - Writes formatted results to out writer
func DisplayExtractionResults(out io.Writer, factCount, eventCount int, competencyCount map[string]int) {
	results := FormatExtractionResults(factCount, eventCount, competencyCount)
	fmt.Fprint(out, results)
}
