package bursts

import (
	"context"
	"fmt"
	"io"

	"github.com/baphled/kariya/internal/cli/cmd/cliutil"
	domain "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	burst_fact "github.com/baphled/kariya/internal/service/career/burstfact"
	tea "github.com/charmbracelet/bubbletea"
)

// DetectBursts analyzes events and suggests burst groupings.
//
// Expected:
//   - svc: Career service instance providing event and burst repository access
//   - out: Writer for formatted output
//   - _: Error writer (unused)
//   - opts: Bubble Tea program options for spinner display
//
// Returns:
//   - Exit code: 0 on success, 1 on error
//
// Side effects:
//   - Queries event repository via context
//   - Calls burst detection service
//   - Saves burst suggestions to repository
//   - Writes formatted output and status messages
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

	DisplayBurstSuggestions(out, suggestions, len(events))
	return SaveBurstSuggestions(ctx, svc, suggestions, opts...)
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
		if err != nil {
			return fmt.Errorf("burst detection failed: %w", err)
		}
		return nil
	}, opts...)
	return suggestions, err
}

// DisplayBurstSuggestions is a thin wrapper around FormatBurstDetectionResults.
//
// Expected:
//   - out: Writer to output formatted results
//   - suggestions: Slice of burst suggestions to display
//   - eventCount: Total number of events analyzed
//
// Returns:
//   - None
//
// Side effects:
//   - Writes formatted output to provided writer
func DisplayBurstSuggestions(out io.Writer, suggestions []burst_fact.BurstSuggestion, eventCount int) {
	output := FormatBurstDetectionResults(suggestions, eventCount)
	fmt.Fprint(out, output)
}

// SaveBurstSuggestions saves burst suggestions and displays results.
//
// Expected:
//   - ctx: Context for database operations
//   - svc: Career service instance with burst repository
//   - suggestions: Slice of burst suggestions to save
//   - opts: Bubble Tea program options for spinner display
//
// Returns:
//   - Exit code: 0 if all bursts saved successfully, 1 on error
//
// Side effects:
//   - Saves bursts to repository via service
//   - Displays spinner during save operation
//   - Writes success/error messages to stderr via cliutil
func SaveBurstSuggestions(
	ctx context.Context,
	svc *careerservice.Service,
	suggestions []burst_fact.BurstSuggestion,
	opts ...tea.ProgramOption,
) int {
	var savedBursts []*domain.Burst
	err := cliutil.RunWithSpinner("Saving burst suggestions...", func() error {
		var err error
		savedBursts, err = svc.SaveBurstSuggestions(ctx, suggestions)
		if err != nil {
			return fmt.Errorf("burst detection failed: %w", err)
		}
		return nil
	}, opts...)

	if err != nil {
		cliutil.PrintError(fmt.Sprintf("Error saving burst suggestions: %v", err))
		return 1
	}

	message, isSuccess := FormatSaveResults(savedBursts, len(suggestions))
	if isSuccess {
		cliutil.PrintSuccess(message)
	} else {
		cliutil.PrintInfo(message)
	}

	return 0
}
