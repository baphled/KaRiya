// Package burst_management implements the BurstManagement intent for managing career bursts.
package burst_management

import (
	"github.com/baphled/kariya/internal/cli/screens"
	burstscreens "github.com/baphled/kariya/internal/cli/screens/burst_management"
	tea "github.com/charmbracelet/bubbletea"
)

// Init bootstraps the intent by loading burst data and transitioning to the list screen.
//
// Returns:
//   - Always nil; the initial screen transition is performed synchronously.
//
// Side effects:
//   - Loads bursts from the repository into filteredBursts.
//   - Sets the initial selectedBurst to the first burst if available.
//   - Transitions the intent to StateList with a new BurstListScreen.
func (i *Intent) Init() tea.Cmd {
	// Load bursts from repository if available.
	if err := i.context.LoadBursts(); err != nil {
		// Error is acceptable - context.Bursts may already be populated (e.g., in tests).
		_ = err
	}

	// Initialize filtered bursts with the provided bursts.
	i.filteredBursts = i.context.Bursts
	if len(i.filteredBursts) > 0 {
		i.selectedBurst = i.filteredBursts[0]
	}

	// Set initial state and transition to list screen.
	i.state = StateList
	i.transitionToScreen(burstscreens.NewBurstListScreen(i.filteredBursts))
	return nil
}

// Update is the central message dispatcher for the intent, routing messages to modals, shortcuts, or the active screen.
//
// Expected:
//   - msg must be a valid tea.Msg (keyboard input, async completion, or custom message).
//
// Returns:
//   - A tea.Cmd for any async follow-up work, or nil when the message is fully handled.
//
// Side effects:
//   - May mutate intent state, open/close modals, or transition screens depending on the message type.
func (i *Intent) Update(msg tea.Msg) tea.Cmd {
	if !i.active {
		return nil
	}

	// Handle custom messages (from modals, etc).
	switch msg := msg.(type) {
	case EditBurstMsg:
		return i.handleEditBurstMsg(msg)
	case BurstEventsLoadedMsg:
		return i.handleBurstEventsLoaded(msg)
	case BurstFactsLoadedMsg:
		return i.handleBurstFactsLoaded(msg)
	case BurstSuggestionsLoadedMsg:
		return i.handleBurstSuggestionsLoaded(msg)
	case SuggestionReviewCompleteMsg:
		// This case handles direct message sends (e.g., from tests).
		// In normal flow, handleModalUpdates calls the handler directly when modal closes.
		return i.handleSuggestionReviewComplete(msg)
	case FactExtractionCompleteMsg:
		return i.handleFactExtractionComplete(msg)
	}

	// Handle modals first (highest priority).
	if cmd := i.handleModalUpdates(msg); cmd != nil || i.HasActiveModal() {
		return cmd
	}

	// Handle key shortcuts when no modal is active.
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if cmd := i.handleKeyShortcuts(keyMsg); cmd != nil {
			return cmd
		}
	}

	// Delegate to active screen.
	if i.activeScreen != nil {
		cmd, result := i.activeScreen.Update(msg)
		if result != nil {
			return i.handleScreenResult(result)
		}
		return cmd
	}

	return nil
}

// View composes the visual output by rendering the active screen with breadcrumbs, help text, and any modal overlays.
//
// Returns:
//   - A fully rendered string representing the current intent UI, including any visible modal overlay.
//
// Side effects:
//   - Rebuilds the modal registry on each call to ensure overlay state is current.
func (i *Intent) View() string {
	if !i.active {
		return "BurstManagement intent is not active"
	}

	// Use screen-based rendering (activeScreen is always set after Init).
	if i.activeScreen != nil {
		return i.renderWithScreen(i.activeScreen)
	}

	// Fallback rendering when no screen is active (e.g., during modal-only states).
	// This ensures modals still render correctly.
	view := i.CreateViewWithBreadcrumbs("Main Menu", "Burst Management", i.getStateName())
	view.WithContent("Loading...")
	view.WithHelp(i.getContextHelp())
	baseView := view.Render()

	i.rebuildModalRegistry()
	return i.modalRegistry.RenderOverlay(baseView)
}

// renderWithScreen renders the current screen with modal overlays.
func (i *Intent) renderWithScreen(screen screens.Screen) string {
	view := i.CreateViewWithBreadcrumbs("Main Menu", "Burst Management", i.getStateName())
	view.WithContent(screen.RenderContent())
	view.WithHelp(i.getContextHelp())
	baseView := view.Render()

	// Rebuild registry and render any visible modal as overlay.
	i.rebuildModalRegistry()
	return i.modalRegistry.RenderOverlay(baseView)
}
