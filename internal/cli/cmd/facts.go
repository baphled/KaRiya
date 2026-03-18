package cmd

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/baphled/kariya/internal/cli/cmd/cliutil"
	"github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
)

// ExtractFacts executes the fact extraction logic for all events.
// It uses ProgressRunner to display a progress bar.
func ExtractFacts(
	svc *careerservice.Service,
	out io.Writer,
	errOut io.Writer,
	runner cliutil.ProgressRunner,
	opts ...tea.ProgramOption,
) int {
	ctx := context.Background()

	var factCount int
	var competencyCount map[string]int

	err := runner.RunWithProgress("Extracting facts from events", 0, func(update func(current int)) error {
		var err error
		factCount, competencyCount, err = svc.ExtractFactsFromAllEvents(ctx, func(current, _ int) {
			update(current)
		})
		if err != nil {
			return fmt.Errorf("fact extraction failed in ExtractFactsFromAllEvents: %w", err)
		}
		return nil
	}, opts...)

	if err != nil {
		fmt.Fprintf(errOut, "Error during fact extraction: %v\n", err)
		return 1
	}

	if factCount == 0 {
		fmt.Fprintf(out, "No facts extracted or no events found.\n")
		return 0
	}

	fmt.Fprintf(out, "\n=== Fact Extraction Results ===\n")
	fmt.Fprintf(out, "Extracted %d facts\n\n", factCount)

	if len(competencyCount) > 0 {
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

	fmt.Fprintf(out, "✓ Fact extraction complete!\n")
	return 0
}

// ListFacts displays all existing facts.
// ListFacts lists all facts from the fact repository and prints them to the provided writer.
func ListFacts(svc *careerservice.Service, out io.Writer, errOut io.Writer) int {
	ctx := context.Background()

	factRepo := svc.GetFactRepository()
	if factRepo == nil {
		fmt.Fprintf(errOut, "Error: Fact repository not configured\n")
		return 1
	}

	facts, err := factRepo.List(ctx, career.FactListFilters{Limit: 10000})
	if err != nil {
		fmt.Fprintf(errOut, "Error retrieving facts: %v\n", err)
		return 1
	}

	if len(facts) == 0 {
		fmt.Fprintf(out, "No facts found in database.\n")
		return 0
	}

	fmt.Fprintf(out, "\n=== Existing Facts ===\n")
	fmt.Fprintf(out, "Total facts: %d\n\n", len(facts))

	for i, fact := range facts {
		fmt.Fprintf(out, "%d. %s\n", i+1, fact.Text)
		fmt.Fprintf(out, "   ID: %s\n", fact.ID)
		fmt.Fprintf(out, "   Source Event: %s\n", fact.SourceEventID)
		if len(fact.CompetencyCategories) > 0 {
			fmt.Fprintf(out, "   Competencies: %s\n", strings.Join(fact.CompetencyCategories, ", "))
		}
		if fact.RoleFit != "" {
			fmt.Fprintf(out, "   Role Fit: %s\n", fact.RoleFit)
		}
		if len(fact.AudienceRelevance) > 0 {
			fmt.Fprintf(out, "   Target Audiences: %s\n", strings.Join(fact.AudienceRelevance, ", "))
		}
		fmt.Fprintf(out, "   Created: %s\n", fact.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Fprintf(out, "\n")
	}

	return 0
}
