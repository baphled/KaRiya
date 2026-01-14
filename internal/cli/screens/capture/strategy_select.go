package capture

import (
	"fmt"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/screens/base"
)

// StrategySelectScreen allows users to choose between Quick and Manual capture strategies.
//
// This screen presents two capture strategy options:
// - Quick: Minimal fields (event text only), date defaults to today
// - Manual: Full form with optional fields (date, company, project, tags)
//
// The screen uses BaseSelectScreen[T] for list navigation and selection.
//
// Keyboard Shortcuts:
// - ↑/↓/j/k: Navigate between strategies
// - Enter: Select strategy and proceed to form
// - Esc: Cancel and return to main menu
//
// Related:
// - internal/cli/intents/capture_event.go (CaptureStrategy constants)
// - internal/cli/screens/base/select_screen.go (BaseSelectScreen[T])
// - tasks/tasks-42-tui-architecture-refactor.md (Phase 1: CaptureEvent Migration)
type StrategySelectScreen struct {
	*base.BaseSelectScreen[intents.CaptureStrategy]
}

// strategyItem represents a strategy with its label and description for rendering.
type strategyItem struct {
	strategy    intents.CaptureStrategy
	label       string
	description string
}

// NewStrategySelectScreen creates a new StrategySelectScreen.
//
// Parameters:
//   - breadcrumbs: Breadcrumb trail for header (e.g., ["Main Menu", "Capture Event"])
//
// Returns a StrategySelectScreen with Quick strategy selected by default (index 0).
func NewStrategySelectScreen(breadcrumbs []string) *StrategySelectScreen {
	// Define strategy options with labels and descriptions
	strategies := []intents.CaptureStrategy{
		intents.StrategyQuick,
		intents.StrategyManual,
	}

	// Create item renderer that shows strategy label and description
	renderer := func(strategy intents.CaptureStrategy) string {
		var label, description string
		switch strategy {
		case intents.StrategyQuick:
			label = "Quick"
			description = "Capture with minimal fields (event text only)"
		case intents.StrategyManual:
			label = "Manual"
			description = "Full form with optional fields (date, company, project, tags)"
		default:
			label = string(strategy)
			description = "Unknown strategy"
		}
		return fmt.Sprintf("%s - %s", label, description)
	}

	// Create base select screen
	baseScreen := base.NewBaseSelectScreen(
		strategies,
		renderer,
		breadcrumbs,
		"Select Capture Strategy",
	)

	return &StrategySelectScreen{
		BaseSelectScreen: baseScreen,
	}
}

// WithInitialSelection sets the initial selection index.
//
// This is useful for restoring state when navigating back.
// If the index is out of bounds, it will be clamped to valid range.
//
// Returns the screen for method chaining.
func (s *StrategySelectScreen) WithInitialSelection(index int) *StrategySelectScreen {
	s.BaseSelectScreen.WithInitialSelection(index)
	return s
}
