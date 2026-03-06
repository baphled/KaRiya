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

// ExtractFacts extracts facts from all events.
//
// Expected:
//   - svc: Initialized career Service instance
//   - out: io.Writer for output messages
//   - _: io.Writer for error output (unused)
//   - opts: Optional Bubble Tea program options for testing
//
// Returns:
//   - 0 on success
//   - 1 on error (retrieval, extraction, or save failure)
//
// Side effects:
//   - Reads events from database via svc
//   - Writes facts to database via svc
//   - Writes output to out writer
func ExtractFacts(svc *careerservice.Service, out io.Writer, _ io.Writer, opts ...tea.ProgramOption) int {
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

	factCount, competencyCount, err := extractAndSaveFacts(ctx, svc, events, opts...)
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

func extractAndSaveFacts(
	ctx context.Context,
	svc *careerservice.Service,
	events []*domain.Event,
	opts ...tea.ProgramOption,
) (int, map[string]int, error) {
	factCount := 0
	competencyCount := make(map[string]int)

	err := cliutil.RunWithProgress(fmt.Sprintf("Extracting facts from %d events...", len(events)), len(events), func(update func(int)) error {
		for i := range events {
			processedFacts := processEventFacts(ctx, svc, events[i])
			factCount += processedFacts
			updateCompetencyCounts(ctx, svc, events[i], competencyCount)
			update(i + 1)
		}
		return nil
	}, opts...)

	return factCount, competencyCount, err
}

func processEventFacts(ctx context.Context, svc *careerservice.Service, event *domain.Event) int {
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

func updateCompetencyCounts(ctx context.Context, svc *careerservice.Service, event *domain.Event, competencyCount map[string]int) {
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

// DisplayExtractionResults writes formatted fact extraction results to the output.
//
// Expected:
//   - out: io.Writer to write output to
//   - factCount: Number of facts extracted
//   - eventCount: Number of events processed
//   - competencyCount: Map of competency names to counts
//
// Side effects:
//   - Writes formatted output to out writer
func DisplayExtractionResults(out io.Writer, factCount, eventCount int, competencyCount map[string]int) {
	output := FormatExtractionResults(factCount, eventCount, competencyCount)
	fmt.Fprint(out, output)
}
