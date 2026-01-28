// Package burst_management implements the BurstManagement intent for managing career bursts.
package burst_management

import (
	"context"
	"errors"
	"fmt"

	"github.com/baphled/kariya/internal/cli/screens"
	burstscreens "github.com/baphled/kariya/internal/cli/screens/burst_management"
	burstmodals "github.com/baphled/kariya/internal/cli/screens/burst_management/modals"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	tea "github.com/charmbracelet/bubbletea"
)

// Init is called when the intent is activated.
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

// Update processes a message in the intent.
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
			return tea.Batch(cmd, i.handleScreenResult(result))
		}
		return cmd
	}

	return nil
}

// noopCmd is a sentinel command to indicate a message was consumed.
// This prevents the message from propagating to the screen after a modal closes.
func noopCmd() tea.Msg { return nil }

// handleModalUpdates handles updates for all modals in priority order.
// Returns a command (possibly noopCmd) if a modal consumed the message.
func (i *Intent) handleModalUpdates(msg tea.Msg) tea.Cmd {
	// Process modals in priority order - first match handles the message.
	if cmd := i.handleErrorModalUpdate(msg); cmd != nil {
		return cmd
	}
	if cmd := i.handleLoadingModalUpdate(msg); cmd != nil {
		return cmd
	}
	if cmd := i.handleDeleteModalUpdate(msg); cmd != nil {
		return cmd
	}
	if cmd := i.handleConfirmModalUpdate(msg); cmd != nil {
		return cmd
	}
	if cmd := i.handleSuggestionModalUpdate(msg); cmd != nil {
		return cmd
	}
	if cmd := i.handleDetailModalUpdate(msg); cmd != nil {
		return cmd
	}
	if cmd := i.handleEventsModalUpdate(msg); cmd != nil {
		return cmd
	}
	if cmd := i.handleFactsModalUpdate(msg); cmd != nil {
		return cmd
	}
	if cmd := i.handleEditModalUpdate(msg); cmd != nil {
		return cmd
	}
	return nil
}

// handleErrorModalUpdate handles error modal updates (highest priority).
func (i *Intent) handleErrorModalUpdate(msg tea.Msg) tea.Cmd {
	if i.errorModal == nil {
		return nil
	}
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.Type == tea.KeyEsc {
		i.errorModal = nil
	}
	return noopCmd
}

// handleLoadingModalUpdate handles loading modal updates (cancellable with Esc).
func (i *Intent) handleLoadingModalUpdate(msg tea.Msg) tea.Cmd {
	if i.loadingModal == nil {
		return nil
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc {
			i.cancelAsyncOperation()
			return noopCmd
		}
		return i.loadingModal.Init()
	case feedback.ModalSpinnerTickMsg:
		return i.loadingModal.Update(msg)
	default:
		return i.loadingModal.Init()
	}
}

// cancelAsyncOperation cancels the current async operation and resets state.
func (i *Intent) cancelAsyncOperation() {
	if i.cancelFunc != nil {
		i.cancelFunc()
		i.cancelFunc = nil
	}
	i.loadingModal = nil
	i.suggestionsLoading = false
	i.extractingFacts = false
	i.state = StateList
}

// handleDeleteModalUpdate handles delete confirmation modal updates.
func (i *Intent) handleDeleteModalUpdate(msg tea.Msg) tea.Cmd {
	if i.deleteModal == nil || !i.deleteModal.IsVisible() {
		return nil
	}
	cmd, confirmed := i.deleteModal.Update(msg)
	if !i.deleteModal.IsVisible() {
		if confirmed && i.selectedBurst != nil {
			i.deleteBurst(i.selectedBurst)
			return cmd
		}
		i.deleteModal = nil
		i.selectedBurst = nil
		i.state = StateList
		return noopCmd
	}
	return cmd
}

// handleConfirmModalUpdate handles burst confirmation modal updates.
func (i *Intent) handleConfirmModalUpdate(msg tea.Msg) tea.Cmd {
	if i.confirmModal == nil || !i.confirmModal.IsVisible() {
		return nil
	}
	cmd, confirmed := i.confirmModal.Update(msg)
	if !i.confirmModal.IsVisible() {
		if confirmed && i.selectedBurst != nil {
			return tea.Batch(cmd, i.confirmBurst())
		}
		i.confirmModal = nil
		if i.selectedBurst != nil {
			return i.showBurstDetailModal(i.selectedBurst)
		}
		return noopCmd
	}
	return cmd
}

// handleSuggestionModalUpdate handles suggestion review modal updates.
func (i *Intent) handleSuggestionModalUpdate(msg tea.Msg) tea.Cmd {
	if i.suggestionModal == nil || !i.suggestionModal.IsVisible() {
		return nil
	}

	// Handle accept key specially.
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "a" {
		return i.handleSuggestionAccept(msg)
	}

	// Handle other keys normally.
	_, cmd := i.suggestionModal.Update(msg)
	if !i.suggestionModal.IsVisible() {
		return i.handleSuggestionModalClosed(cmd)
	}
	return cmd
}

// handleSuggestionAccept handles accepting a suggestion via 'a' key.
func (i *Intent) handleSuggestionAccept(msg tea.Msg) tea.Cmd {
	currentSuggestion := i.suggestionModal.GetCurrentSuggestion()
	if currentSuggestion == nil {
		return noopCmd
	}

	cmd := i.saveAndExtractBurst(*currentSuggestion)
	i.suggestionModal.Update(msg)

	if !i.suggestionModal.IsVisible() {
		i.clearSuggestionState()
		i.state = StateList
		i.transitionToScreen(burstscreens.NewBurstListScreen(i.filteredBursts))
	}
	return cmd
}

// handleSuggestionModalClosed handles suggestion modal closure.
func (i *Intent) handleSuggestionModalClosed(cmd tea.Cmd) tea.Cmd {
	accepted := i.suggestionModal.GetAcceptedSuggestions()
	cancelled := i.suggestionModal.GetAction() == burstmodals.SuggestionActionCancel
	i.clearSuggestionState()

	if len(accepted) > 0 && !cancelled {
		completeMsg := SuggestionReviewCompleteMsg{
			AcceptedSuggestions: accepted,
			Cancelled:           false,
		}
		return tea.Batch(cmd, func() tea.Msg { return completeMsg })
	}

	i.state = StateList
	i.transitionToScreen(burstscreens.NewBurstListScreen(i.filteredBursts))
	return noopCmd
}

// handleDetailModalUpdate handles detail modal updates with keyboard shortcuts.
func (i *Intent) handleDetailModalUpdate(msg tea.Msg) tea.Cmd {
	if i.detailModal == nil || !i.detailModal.IsVisible() {
		return nil
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if cmd := i.handleDetailModalKeypress(keyMsg); cmd != nil {
			return cmd
		}
	}

	_, cmd := i.detailModal.Update(msg)
	return cmd
}

// handleDetailModalKeypress handles keyboard shortcuts in detail modal.
func (i *Intent) handleDetailModalKeypress(keyMsg tea.KeyMsg) tea.Cmd {
	switch keyMsg.String() {
	case "v":
		i.detailModal.Hide()
		return i.showBurstEventsModal()
	case "f":
		i.detailModal.Hide()
		return i.showBurstFactsModal()
	case "e":
		i.detailModal.Hide()
		return i.openEditModal(i.selectedBurst)
	case "d":
		i.detailModal.Hide()
		return i.openDeleteModal(i.selectedBurst)
	case "c":
		i.detailModal.Hide()
		return i.showConfirmBurstModal()
	}

	if keyMsg.Type == tea.KeyEsc || keyMsg.Type == tea.KeyEnter {
		i.detailModal = nil
		return noopCmd
	}
	return nil
}

// showConfirmBurstModal shows the confirmation modal for the selected burst.
func (i *Intent) showConfirmBurstModal() tea.Cmd {
	existingFacts := i.loadBurstFacts(i.selectedBurst)
	if len(existingFacts) > 0 {
		i.confirmModal = feedback.NewConfirmModal(
			"Re-extract Facts",
			fmt.Sprintf("This burst already has %d facts. Re-extract and add more?", len(existingFacts)),
		).WithVariant(feedback.ConfirmDefault)
	} else {
		i.confirmModal = feedback.NewConfirmModal(
			"Confirm Burst",
			"Mark this burst as confirmed and extract facts?",
		).WithVariant(feedback.ConfirmDefault)
	}
	return i.confirmModal.Init()
}

// handleEventsModalUpdate handles events modal updates.
func (i *Intent) handleEventsModalUpdate(msg tea.Msg) tea.Cmd {
	if i.eventsModal == nil || !i.eventsModal.IsVisible() {
		return nil
	}
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.Type == tea.KeyEsc || keyMsg.Type == tea.KeyEnter {
			i.eventsModal = nil
			if i.selectedBurst != nil {
				return i.showBurstDetailModal(i.selectedBurst)
			}
			return noopCmd
		}
	}
	_, cmd := i.eventsModal.Update(msg)
	return cmd
}

// handleFactsModalUpdate handles facts modal updates.
func (i *Intent) handleFactsModalUpdate(msg tea.Msg) tea.Cmd {
	if i.factsModal == nil || !i.factsModal.IsVisible() {
		return nil
	}
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.Type == tea.KeyEsc || keyMsg.Type == tea.KeyEnter {
			i.factsModal = nil
			if i.selectedBurst != nil {
				return i.showBurstDetailModal(i.selectedBurst)
			}
			return noopCmd
		}
	}
	_, cmd := i.factsModal.Update(msg)
	return cmd
}

// handleEditModalUpdate handles edit modal updates.
func (i *Intent) handleEditModalUpdate(msg tea.Msg) tea.Cmd {
	if i.editModal == nil || !i.editModal.IsVisible() {
		return nil
	}
	cmd, completed, formData := i.editModal.Update(msg)
	if !i.editModal.IsVisible() {
		if completed && formData != nil && i.selectedBurst != nil {
			editMsg := EditBurstMsg{
				BurstID:     i.selectedBurst.ID,
				Name:        formData.Name,
				Description: formData.Description,
			}
			i.editModal = nil
			return tea.Batch(cmd, func() tea.Msg { return editMsg })
		}
		i.editModal = nil
		i.state = StateList
		i.transitionToScreen(burstscreens.NewBurstListScreen(i.filteredBursts))
		return noopCmd
	}
	return cmd
}

// handleKeyShortcuts handles keyboard shortcuts when no modal is active.
func (i *Intent) handleKeyShortcuts(keyMsg tea.KeyMsg) tea.Cmd {
	if keyMsg.String() == "s" {
		// Trigger burst suggestion detection.
		i.state = StateSuggesting
		return i.startBurstDetection()
	}
	return nil
}

// handleFactExtractionComplete handles the FactExtractionCompleteMsg.
func (i *Intent) handleFactExtractionComplete(msg FactExtractionCompleteMsg) tea.Cmd {
	i.extractingFacts = false
	i.loadingModal = nil

	if msg.Error != nil {
		// Silently ignore cancelled operations - user already knows they cancelled.
		if errors.Is(msg.Error, context.Canceled) {
			i.state = StateList
			return nil
		}
		i.errorModal = feedback.NewErrorModal("Extraction Failed", msg.Error.Error())
		// Only transition to list if suggestion modal is not visible.
		// User may still be reviewing remaining suggestions.
		if i.suggestionModal == nil || !i.suggestionModal.IsVisible() {
			i.state = StateList
		}
		return nil
	}

	i.extractedFactsCount = len(msg.Facts)

	// If the suggestion modal is still visible (user reviewing remaining suggestions),
	// do NOT change state or transition screens. Let user continue reviewing.
	if i.suggestionModal != nil && i.suggestionModal.IsVisible() {
		return nil
	}

	i.state = StateList

	// Only show detail modal if a burst was selected (e.g., from confirm flow).
	// When coming from suggestion acceptance flow, selectedBurst is nil,
	// so we stay on the list view showing the newly created bursts.
	if i.selectedBurst != nil {
		return i.showBurstDetailModal(i.selectedBurst)
	}

	// Refresh list screen to ensure visual update after extraction completes.
	i.transitionToScreen(burstscreens.NewBurstListScreen(i.filteredBursts))
	return nil
}

// View renders the intent's current state.
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
