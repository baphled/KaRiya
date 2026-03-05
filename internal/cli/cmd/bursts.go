// Package cmd provides CLI command handlers for the KaRiya application.
package cmd

import (
	"context"
	"fmt"
	"io"

	"github.com/baphled/kariya/internal/cli/cmd/cliutil"
	"github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
)

// DetectBursts analyzes events and suggests burst groupings.
// Detects bursts from all events and saves them to the repository.
// It uses ProgressRunner to display a spinner during detection.
func DetectBursts(
	svc *careerservice.Service,
	out io.Writer,
	errOut io.Writer,
	runner cliutil.ProgressRunner,
	opts ...tea.ProgramOption,
) int {
	ctx := context.Background()

	var suggestionsCount int
	var savedCount int

	err := runner.RunWithSpinner("Detecting bursts from events", func() error {
		var err error
		suggestionsCount, savedCount, err = svc.DetectAndSaveBursts(ctx)
		return err
	}, opts...)

	if err != nil {
		fmt.Fprintf(errOut, "Error detecting bursts: %v\n", err)
		return 1
	}

	if suggestionsCount == 0 {
		fmt.Fprintf(out, "No bursts detected or no events found.\n")
		return 0
	}

	fmt.Fprintf(out, "\n=== Burst Detection Results ===\n")
	fmt.Fprintf(out, "Detected %d potential bursts.\n", suggestionsCount)

	if savedCount > 0 {
		fmt.Fprintf(out, "✓ Burst detection complete! Saved %d of %d bursts to database.\n", savedCount, suggestionsCount)
	} else {
		fmt.Fprintf(out, "⚠ Burst detection complete, but no bursts were saved (repository may not be configured).\n")
	}

	return 0
}

// ListBursts displays all existing bursts.
func ListBursts(svc *careerservice.Service, out io.Writer, errOut io.Writer) int {
	ctx := context.Background()

	burstRepo := svc.GetBurstRepository()
	if burstRepo == nil {
		fmt.Fprintf(errOut, "Error: Burst repository not configured\n")
		return 1
	}

	bursts, err := burstRepo.List(ctx, career.BurstListFilters{Limit: 10000})
	if err != nil {
		fmt.Fprintf(errOut, "Error retrieving bursts: %v\n", err)
		return 1
	}

	if len(bursts) == 0 {
		fmt.Fprintf(out, "No bursts found.\n")
		return 0
	}

	fmt.Fprintf(out, "\n=== Bursts ===\n")
	fmt.Fprintf(out, "Total bursts: %d\n\n", len(bursts))

	for i, burst := range bursts {
		burstName := burst.Name
		if burstName == "" {
			burstName = fmt.Sprintf("Burst %d", i+1)
		}
		fmt.Fprintf(out, "%d. %s\n", i+1, burstName)
		fmt.Fprintf(out, "   ID: %s\n", burst.ID)
		fmt.Fprintf(out, "   Events: %d\n", len(burst.EventIDs))
		if burst.Description != "" {
			fmt.Fprintf(out, "   Description: %s\n", burst.Description)
		}
		fmt.Fprintf(out, "   Created: %s\n", burst.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Fprintf(out, "\n")
	}

	return 0
}
