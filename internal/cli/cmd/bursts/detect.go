package bursts

import (
	"context"
	"fmt"
	"io"

	"github.com/baphled/kariya/internal/cli/cmd/cliutil"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/cli/uikit/theme"
	domain "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	burst_fact "github.com/baphled/kariya/internal/service/career/burstfact"
	tea "github.com/charmbracelet/bubbletea"
)

// DetectBursts analyzes events and suggests burst groupings.
// Detects bursts from all events and saves them to the repository.
func DetectBursts(svc *careerservice.Service, out io.Writer, _ io.Writer, opts ...tea.ProgramOption) int {
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

	eventIDs := extractEventIDs(events)
	suggestions, err := detectBurstsFromEvents(ctx, svc, eventIDs, len(events), opts...)
	if err != nil {
		cliutil.PrintError(fmt.Sprintf("Error detecting bursts: %v", err))
		return 1
	}

	if len(suggestions) == 0 {
		cliutil.PrintInfo("No bursts detected.")
		return 0
	}

	displayBurstSuggestions(out, suggestions, len(events))
	return saveBurstSuggestions(ctx, svc, suggestions, opts...)
}

func extractEventIDs(events []*domain.Event) []string {
	eventIDs := make([]string, len(events))
	for i, event := range events {
		eventIDs[i] = event.ID
	}
	return eventIDs
}

func detectBurstsFromEvents(
	ctx context.Context,
	svc *careerservice.Service,
	eventIDs []string,
	eventCount int,
	opts ...tea.ProgramOption,
) ([]burst_fact.BurstSuggestion, error) {
	var suggestions []burst_fact.BurstSuggestion
	err := cliutil.RunWithSpinner(fmt.Sprintf("Detecting bursts from %d events...", eventCount), func() error {
		var err error
		suggestions, err = svc.SuggestBursts(ctx, eventIDs)
		return err
	}, opts...)
	return suggestions, err
}

func displayBurstSuggestions(out io.Writer, suggestions []burst_fact.BurstSuggestion, eventCount int) {
	th := theme.Default()
	fmt.Fprintln(out, primitives.Title("Burst Detection Results", th).Render())
	fmt.Fprintf(out, "Detected %d bursts from %d events:\n\n", len(suggestions), eventCount)

	for i, burst := range suggestions {
		burstName := burst.Name
		if burstName == "" {
			burstName = fmt.Sprintf("Burst %d", i+1)
		}
		fmt.Fprintf(out, "%d. %s\n", i+1, burstName)
		fmt.Fprintf(out, "   Events: %d\n", len(burst.EventIDs))
		if burst.Description != "" {
			fmt.Fprintf(out, "   Description: %s\n", burst.Description)
		}
		fmt.Fprintf(out, "   Confidence: %.1f%%\n", burst.ConfidenceScore*100)
		fmt.Fprintf(out, "\n")
	}
}

func saveBurstSuggestions(
	ctx context.Context,
	svc *careerservice.Service,
	suggestions []burst_fact.BurstSuggestion,
	opts ...tea.ProgramOption,
) int {
	var savedBursts []*domain.Burst
	err := cliutil.RunWithSpinner("Saving burst suggestions...", func() error {
		var err error
		savedBursts, err = svc.SaveBurstSuggestions(ctx, suggestions)
		return err
	}, opts...)

	if err != nil {
		cliutil.PrintError(fmt.Sprintf("Error saving burst suggestions: %v", err))
		return 1
	}

	savedCount := len(savedBursts)
	if savedCount > 0 {
		cliutil.PrintSuccess(fmt.Sprintf("Burst detection complete! Saved %d of %d bursts to database.", savedCount, len(suggestions)))
	} else {
		cliutil.PrintInfo("Burst detection complete, but no bursts were saved (repository may not be configured).")
	}

	return 0
}
