// Package facts provides Cobra CLI commands for fact extraction and listing.
package facts

import (
	"context"
	"fmt"
	"io"
	"sort"

	"github.com/baphled/kariya/internal/cli/cmd/cliutil"
	domain "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
)

// ExtractFacts extracts facts from all events.
func ExtractFacts(svc *careerservice.Service, out io.Writer, _ io.Writer) int {
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

	factCount, competencyCount, err := extractAndSaveFacts(ctx, svc, events)
	if err != nil {
		cliutil.PrintError(fmt.Sprintf("Error extracting facts: %v", err))
		return 1
	}

	if factCount == 0 {
		cliutil.PrintInfo("No facts extracted.")
		return 0
	}

	displayExtractionResults(out, factCount, len(events), competencyCount)
	cliutil.PrintSuccess("Fact extraction complete!")
	return 0
}

func extractAndSaveFacts(ctx context.Context, svc *careerservice.Service, events []*domain.Event) (int, map[string]int, error) {
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
	})

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

func displayExtractionResults(out io.Writer, factCount, eventCount int, competencyCount map[string]int) {
	fmt.Fprintf(out, "\n=== Fact Extraction Results ===\n")
	fmt.Fprintf(out, "Extracted %d facts from %d events\n\n", factCount, eventCount)

	if len(competencyCount) == 0 {
		return
	}

	fmt.Fprintf(out, "Competency breakdown:\n")
	competencies := make([]string, 0, len(competencyCount))
	for c := range competencyCount {
		competencies = append(competencies, c)
	}
	sort.Strings(competencies)

	for _, c := range competencies {
		fmt.Fprintf(out, "  - %s: %d facts\n", c, competencyCount[c])
	}
	fmt.Fprintf(out, "\n")
}
