package capture

import (
	"fmt"

	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/cli/types"
)

// StrategySelectScreen allows users to choose between Quick and Manual capture strategies.
//
// This screen presents two capture strategy options:
// - Quick: Minimal fields (event text only), date defaults to today
// - Manual: Full form with optional fields (date, company, project, tags)
//
// The screen uses SelectScreen[T] for list navigation and selection.
//
// Keyboard Shortcuts:
// - ↑/↓/j/k: Navigate between strategies
// - Enter: Select strategy and proceed to form
// - Esc: Cancel and return to main menu
//
// Related:
// - internal/cli/intents/captureevent/constants.go (CaptureStrategy constants)
// - internal/cli/screens/base/select_screen.go (SelectScreen[T])
// - tasks/tasks-42-tui-architecture-refactor.md (Phase 1: CaptureEvent Migration).
type StrategySelectScreen struct {
	*base.SelectScreen[types.CaptureStrategy]
}

// NewStrategySelectScreen creates a new StrategySelectScreen.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized StrategySelectScreen ready for use.
//
// Side effects:
//   - None.
func NewStrategySelectScreen(breadcrumbs []string) *StrategySelectScreen {
	// Define strategy options with labels and descriptions
	strategies := []types.CaptureStrategy{
		types.StrategyQuick,
		types.StrategyManual,
	}

	// Create item renderer that shows strategy label and description
	renderer := func(strategy types.CaptureStrategy) string {
		var label, description string
		switch strategy {
		case types.StrategyQuick:
			label = "Quick"
			description = "Capture with minimal fields (event text only)"
		case types.StrategyManual:
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
		SelectScreen: baseScreen,
	}
}

// WithInitialSelection sets the initial selection index.
//
// Expected:
//   - int must be valid.
//
// Returns:
//   - A fully initialized StrategySelectScreen ready for use.
//
// Side effects:
//   - None.
func (s *StrategySelectScreen) WithInitialSelection(index int) *StrategySelectScreen {
	s.SelectScreen.WithInitialSelection(index)
	return s
}
