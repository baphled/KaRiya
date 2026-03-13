package bursts

import (
	"fmt"
	"strings"

	domain "github.com/baphled/kariya/internal/domain/career"
	burst_fact "github.com/baphled/kariya/internal/service/career/burstfact"
	"github.com/baphled/kariya/internal/ui/uikit/primitives"
	"github.com/baphled/kariya/internal/ui/uikit/theme"
)

// FormatBurstDetectionResults generates formatted output for burst detection results.
//
// Expected:
//   - suggestions: Slice of burst suggestions to format
//   - eventCount: Total number of events analyzed
//
// Returns:
//   - Formatted string with title, summary, and detailed burst information
//   - Each burst includes: name (or auto-numbered), event count, description, confidence score
//
// Side effects:
//   - None (pure function)
func FormatBurstDetectionResults(suggestions []burst_fact.BurstSuggestion, eventCount int) string {
	var b strings.Builder

	th := theme.Default()
	b.WriteString(primitives.Title("Burst Detection Results", th).Render())
	b.WriteString("\n")
	fmt.Fprintf(&b, "Detected %d bursts from %d events:\n\n", len(suggestions), eventCount)

	for i, burst := range suggestions {
		burstName := burst.Name
		if burstName == "" {
			burstName = fmt.Sprintf("Burst %d", i+1)
		}
		fmt.Fprintf(&b, "%d. %s\n", i+1, burstName)
		fmt.Fprintf(&b, "   Events: %d\n", len(burst.EventIDs))
		if burst.Description != "" {
			fmt.Fprintf(&b, "   Description: %s\n", burst.Description)
		}
		fmt.Fprintf(&b, "   Confidence: %.1f%%\n", burst.ConfidenceScore*100)
		b.WriteString("\n")
	}

	return b.String()
}

// FormatSaveResults generates formatted output for save operation results.
//
// Expected:
//   - savedBursts: Slice of successfully saved bursts
//   - totalSuggestions: Total number of burst suggestions attempted
//
// Returns:
//   - message: Formatted result message
//   - isSuccess: true if any bursts were saved, false otherwise
//
// Side effects:
//   - None (pure function)
func FormatSaveResults(savedBursts []*domain.Burst, totalSuggestions int) (message string, isSuccess bool) {
	savedCount := len(savedBursts)
	if savedCount > 0 {
		return fmt.Sprintf("Burst detection complete! Saved %d of %d bursts to database.", savedCount, totalSuggestions), true
	}
	return "Burst detection complete, but no bursts were saved (repository may not be configured).", false
}

// FormatBurstList generates formatted output for listing existing bursts.
//
// Expected:
//   - bursts: Slice of domain.Burst pointers to format
//
// Returns:
//   - Formatted string with title, count, and detailed burst information
//   - Each burst includes: name (or auto-numbered), ID, event count, description, competency focus, creation time
//
// Side effects:
//   - None (pure function)
func FormatBurstList(bursts []*domain.Burst) string {
	var b strings.Builder

	th := theme.Default()
	b.WriteString(primitives.Title("Bursts", th).Render())
	b.WriteString("\n")
	fmt.Fprintf(&b, "Total bursts: %d\n\n", len(bursts))

	for i, burst := range bursts {
		burstName := burst.Name
		if burstName == "" {
			burstName = fmt.Sprintf("Burst %d", i+1)
		}
		fmt.Fprintf(&b, "%d. %s\n", i+1, burstName)
		fmt.Fprintf(&b, "   ID: %s\n", burst.ID)
		fmt.Fprintf(&b, "   Events: %d\n", len(burst.EventIDs))
		if burst.Description != "" {
			fmt.Fprintf(&b, "   Description: %s\n", burst.Description)
		}
		if burst.Description != "" {
			fmt.Fprintf(&b, "   Competency Focus: %s\n", burst.Description)
		}
		fmt.Fprintf(&b, "   Created: %s\n", burst.CreatedAt.Format("2006-01-02 15:04:05"))
		b.WriteString("\n")
	}

	return b.String()
}
