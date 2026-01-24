// Package burst_management implements the BurstManagement intent for managing career bursts.
package burst_management

import (
	"context"
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/domain/career"
)

// burstRowFormatter formats a burst for table display.
func burstRowFormatter(burst *career.Burst, index int) []string {
	// Column 1: Name (truncate to 27 chars).
	nameStr := burst.Name
	if len(nameStr) > 27 {
		nameStr = nameStr[:27] + "..."
	}

	// Column 2: Description (truncated preview, max 32 chars).
	descStr := strings.TrimSpace(burst.Description)
	descStr = strings.ReplaceAll(descStr, "\n", " ")
	descStr = strings.ReplaceAll(descStr, "\r", " ")
	if descStr == "" {
		descStr = "-"
	} else if len(descStr) > 32 {
		descStr = descStr[:32] + "..."
	}

	// Column 3: Confirmed Status.
	confirmedStr := "✗ No"
	if burst.Confirmed {
		confirmedStr = "✓ Yes"
	}

	// Column 4: Event Count.
	eventCount := fmt.Sprintf("%d", len(burst.EventIDs))

	// Column 5: Created Date (YYYY-MM-DD).
	createdStr := burst.CreatedAt.Format("2006-01-02")

	return []string{nameStr, descStr, confirmedStr, eventCount, createdStr}
}

// getTerminalDimensions returns current terminal dimensions with fallback defaults.
func (i *Intent) getTerminalDimensions() (width, height int) {
	width, height = behaviors.DefaultModalDimensions()
	if termInfo := i.GetTerminalInfo(); termInfo != nil && termInfo.Width > 0 && termInfo.Height > 0 {
		width, height = termInfo.Width, termInfo.Height
	}
	return
}

// getContext returns a context for service calls.
func (i *Intent) getContext() context.Context {
	return context.Background()
}

// getStateName returns a human-readable name for the current state.
func (i *Intent) getStateName() string {
	switch i.state {
	case StateList:
		return "Burst List"
	case StateDetail:
		return "Burst Details"
	case StateDetailEvents:
		return "Burst Events"
	case StateDetailFacts:
		return "Burst Facts"
	case StateEdit:
		return "Edit Burst"
	case StateDeleteConfirm:
		return "Delete Confirmation"
	case StateConfirm:
		return "Confirm Burst"
	case StateExtractingFacts:
		return "Extracting Facts"
	case StateSuggesting:
		return "Suggesting Bursts"
	case StateSuggestionReview:
		return "Review Suggestion"
	default:
		return "Unknown"
	}
}

// getContextHelp returns themed keyboard shortcuts for the current state.
func (i *Intent) getContextHelp() string {
	theme := i.Theme()

	switch i.state {
	case StateList:
		badges := []*primitives.Badge{
			primitives.NavigateBadge(theme),
			primitives.HelpKeyBadge("Enter", "View Details", theme),
			primitives.AddBadge(theme),
			primitives.EditBadge(theme),
			primitives.DeleteBadge(theme),
			primitives.HelpKeyBadge("s", "Suggest", theme),
			primitives.BackBadge(theme),
		}

		return intents.CombineThemedFooters(
			intents.ThemedCustomFooter(theme, badges...),
			intents.ThemedGlobalBadges(theme),
		)

	case StateDetail:
		badges := []*primitives.Badge{
			primitives.HelpKeyBadge("v", "View Events", theme),
			primitives.HelpKeyBadge("f", "View Facts", theme),
			primitives.EditBadge(theme),
			primitives.DeleteBadge(theme),
			primitives.HelpKeyBadge("c", "Confirm", theme),
			primitives.BackBadge(theme),
		}

		return intents.CombineThemedFooters(
			intents.ThemedCustomFooter(theme, badges...),
			intents.ThemedGlobalBadges(theme),
		)

	case StateDetailEvents, StateDetailFacts:
		badges := []*primitives.Badge{
			primitives.BackBadge(theme),
		}

		return intents.CombineThemedFooters(
			intents.ThemedCustomFooter(theme, badges...),
			intents.ThemedGlobalBadges(theme),
		)

	default:
		return intents.ThemedGlobalBadges(theme)
	}
}

// transitionToScreen sets the active screen and updates state.
func (i *Intent) transitionToScreen(screen interface{}) {
	// TODO: Implement screen transition when screens are created
	// For now, store the screen in activeScreen field
	// i.activeScreen = screen

	// Set terminal info if available
	termInfo := i.GetTerminalInfo()
	if termInfo != nil && termInfo.Width > 0 && termInfo.Height > 0 {
		// screen.SetTerminalInfo(termInfo.Width, termInfo.Height)
	}

	// Set theme if available
	if theme := i.Theme(); theme != nil {
		// screen.SetTheme(theme)
	}

	// Set logo if available
	if logo := i.GetLogo(); logo != nil {
		// screen.SetLogo(logo, i.GetLogoSpacing())
	}
}
