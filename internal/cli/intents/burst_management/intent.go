// Package burst_management implements the BurstManagement intent for managing career bursts.
package burst_management

import (
	"fmt"

	"github.com/baphled/kariya/internal/cli/screens"
	burstscreens "github.com/baphled/kariya/internal/cli/screens/burst_management"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
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

	// Delegate to active screen (if available).
	if i.activeScreen != nil {
		cmd, result := i.activeScreen.Update(msg)
		if result != nil {
			return tea.Batch(cmd, i.handleScreenResult(result))
		}
		return cmd
	}

	// Fallback to state-based handling (for states without screens yet)
	switch i.state {
	case StateList:
		return i.updateListView(msg)
	case StateDetail:
		return i.updateDetailView(msg)
	case StateDetailEvents:
		return i.updateDetailEventsView(msg)
	case StateDetailFacts:
		return i.updateDetailFactsView(msg)
	case StateEdit:
		return i.updateEditView(msg)
	case StateDeleteConfirm:
		return i.updateDeleteConfirmView(msg)
	case StateConfirm:
		return i.updateConfirmView(msg)
	case StateExtractingFacts:
		return i.updateExtractingFactsView(msg)
	case StateSuggesting:
		return i.updateSuggestingView(msg)
	case StateSuggestionReview:
		return i.updateSuggestionReviewView(msg)
	}

	return nil
}

// noopCmd is a sentinel command to indicate a message was consumed.
// This prevents the message from propagating to the screen after a modal closes.
func noopCmd() tea.Msg { return nil }

// handleModalUpdates handles updates for all modals in priority order.
// Returns a command (possibly noopCmd) if a modal consumed the message.
func (i *Intent) handleModalUpdates(msg tea.Msg) tea.Cmd {
	// Error modal has special handling (highest priority).
	if i.errorModal != nil {
		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.Type == tea.KeyEsc {
			i.errorModal = nil
			return noopCmd
		}
		return noopCmd
	}

	// Delete confirmation modal.
	if i.deleteModal != nil && i.deleteModal.IsVisible() {
		cmd, confirmed := i.deleteModal.Update(msg)
		if !i.deleteModal.IsVisible() {
			// Modal was closed.
			if confirmed && i.selectedBurst != nil {
				// User confirmed deletion - delete the burst.
				return tea.Batch(cmd, i.deleteBurst(i.selectedBurst))
			}
			// User cancelled or modal closed without confirmation.
			i.deleteModal = nil
			i.selectedBurst = nil
			i.state = StateList
			return cmd
		}
		return cmd
	}

	// Confirm burst modal.
	if i.confirmModal != nil && i.confirmModal.IsVisible() {
		cmd, confirmed := i.confirmModal.Update(msg)
		if !i.confirmModal.IsVisible() {
			// Modal was closed.
			if confirmed && i.selectedBurst != nil {
				// User confirmed - mark burst as confirmed.
				return tea.Batch(cmd, i.confirmBurst())
			}
			// User cancelled - show detail modal again.
			i.confirmModal = nil
			if i.selectedBurst != nil {
				return i.showBurstDetailModal(i.selectedBurst)
			}
			return noopCmd
		}
		return cmd
	}

	// Detail modals (detail, events, facts) - handle keyboard shortcuts.
	if i.detailModal != nil && i.detailModal.IsVisible() {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "v":
				// View events - show events modal.
				i.detailModal.Hide()
				return i.showBurstEventsModal()
			case "f":
				// View facts - show facts modal.
				i.detailModal.Hide()
				return i.showBurstFactsModal()
			case "e":
				// Edit burst.
				i.detailModal.Hide()
				return i.openEditModal(i.selectedBurst)
			case "d":
				// Delete burst.
				i.detailModal.Hide()
				return i.openDeleteModal(i.selectedBurst)
			case "c":
				// Confirm burst.
				i.detailModal.Hide()
				i.confirmModal = feedback.NewConfirmModal(
					"Confirm Burst",
					"Mark this burst as confirmed and extract facts?",
				).WithVariant(feedback.ConfirmDefault)
				return i.confirmModal.Init()
			case "esc", "enter":
				// Close detail modal.
				i.detailModal = nil
				return noopCmd
			}
		}
		// Update the modal.
		_, cmd := i.detailModal.Update(msg)
		return cmd
	}

	// Events modal.
	if i.eventsModal != nil && i.eventsModal.IsVisible() {
		if keyMsg, ok := msg.(tea.KeyMsg); ok && (keyMsg.String() == "esc" || keyMsg.String() == "enter") {
			// Close events modal and show detail modal again.
			i.eventsModal = nil
			if i.selectedBurst != nil {
				return i.showBurstDetailModal(i.selectedBurst)
			}
			return noopCmd
		}
		// Update the modal.
		_, cmd := i.eventsModal.Update(msg)
		return cmd
	}

	// Facts modal.
	if i.factsModal != nil && i.factsModal.IsVisible() {
		if keyMsg, ok := msg.(tea.KeyMsg); ok && (keyMsg.String() == "esc" || keyMsg.String() == "enter") {
			// Close facts modal and show detail modal again.
			i.factsModal = nil
			if i.selectedBurst != nil {
				return i.showBurstDetailModal(i.selectedBurst)
			}
			return noopCmd
		}
		// Update the modal.
		_, cmd := i.factsModal.Update(msg)
		return cmd
	}

	// Edit burst modal.
	if i.editModal != nil && i.editModal.IsVisible() {
		cmd, completed, formData := i.editModal.Update(msg)
		if !i.editModal.IsVisible() {
			// Modal was closed.
			if completed && formData != nil && i.selectedBurst != nil {
				// User completed form - send EditBurstMsg.
				editMsg := EditBurstMsg{
					BurstID:     i.selectedBurst.ID,
					Name:        formData.Name,
					Description: formData.Description,
				}
				// Clear modal and return msg to trigger edit handling.
				i.editModal = nil
				return tea.Batch(cmd, func() tea.Msg { return editMsg })
			}
			// User cancelled - return to list view.
			i.editModal = nil
			return noopCmd
		}
		return cmd
	}

	return nil
}

// handleKeyShortcuts handles keyboard shortcuts when no modal is active.
func (i *Intent) handleKeyShortcuts(keyMsg tea.KeyMsg) tea.Cmd {
	// TODO: Add global shortcuts (e.g., 's' for suggest, 'a' for add)
	// when implemented
	return nil
}

// updateListView handles messages while viewing the burst list.
func (i *Intent) updateListView(msg tea.Msg) tea.Cmd {
	// Minimal implementation - will be expanded later.
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return tea.Quit
		case "esc":
			i.SetCancelled()
			return nil
		}
	}
	return nil
}

// updateDetailView handles messages while viewing burst details.
func (i *Intent) updateDetailView(msg tea.Msg) tea.Cmd {
	return nil
}

// updateDetailEventsView handles messages while viewing events in a burst.
func (i *Intent) updateDetailEventsView(msg tea.Msg) tea.Cmd {
	return nil
}

// updateDetailFactsView handles messages while viewing facts from a burst.
func (i *Intent) updateDetailFactsView(msg tea.Msg) tea.Cmd {
	return nil
}

// updateEditView handles the edit state.
func (i *Intent) updateEditView(msg tea.Msg) tea.Cmd {
	return nil
}

// updateDeleteConfirmView handles the delete confirmation state.
func (i *Intent) updateDeleteConfirmView(msg tea.Msg) tea.Cmd {
	return nil
}

// updateConfirmView handles the confirmation state.
// The confirm modal is handled in handleModalUpdates, so this is a no-op.
func (i *Intent) updateConfirmView(msg tea.Msg) tea.Cmd {
	return nil
}

// updateExtractingFactsView handles the fact extraction state.
func (i *Intent) updateExtractingFactsView(msg tea.Msg) tea.Cmd {
	return nil
}

// updateSuggestingView handles the burst suggestion loading state.
func (i *Intent) updateSuggestingView(msg tea.Msg) tea.Cmd {
	return nil
}

// updateSuggestionReviewView handles the suggestion review state.
func (i *Intent) updateSuggestionReviewView(msg tea.Msg) tea.Cmd {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "esc":
			// Cancel suggestion review and return to list.
			i.state = StateList
			i.suggestions = nil
			i.currentSuggestionIdx = 0
			// Restore list screen.
			i.transitionToScreen(burstscreens.NewBurstListScreen(i.filteredBursts))
			return nil

		case "n", "right":
			// Next suggestion.
			if i.currentSuggestionIdx < len(i.suggestions)-1 {
				i.currentSuggestionIdx++
			}
			return nil

		case "p", "left":
			// Previous suggestion.
			if i.currentSuggestionIdx > 0 {
				i.currentSuggestionIdx--
			}
			return nil

		case "a":
			// Accept current suggestion - create burst.
			if i.currentSuggestionIdx < len(i.suggestions) {
				return i.acceptSuggestion(i.suggestions[i.currentSuggestionIdx])
			}
			return nil

		case "r":
			// Reject current suggestion - move to next or return to list.
			if i.currentSuggestionIdx < len(i.suggestions)-1 {
				// Remove current suggestion and stay at same index.
				i.suggestions = append(i.suggestions[:i.currentSuggestionIdx], i.suggestions[i.currentSuggestionIdx+1:]...)
				// Adjust index if we're past the end.
				if i.currentSuggestionIdx >= len(i.suggestions) {
					i.currentSuggestionIdx = len(i.suggestions) - 1
				}
			} else {
				// Last suggestion - remove and return to list.
				i.suggestions = i.suggestions[:i.currentSuggestionIdx]
			}

			// If no suggestions left, return to list.
			if len(i.suggestions) == 0 {
				i.state = StateList
				i.currentSuggestionIdx = 0
				// Restore list screen.
				i.transitionToScreen(burstscreens.NewBurstListScreen(i.filteredBursts))
			}
			return nil
		}
	}
	return nil
}

// View renders the intent's current state.
func (i *Intent) View() string {
	if !i.active {
		return "BurstManagement intent is not active"
	}

	// Use screen-based rendering when screen is available.
	if i.activeScreen != nil {
		return i.renderWithScreen(i.activeScreen)
	}

	// Fallback to state-based content (for states without screens yet).
	content := i.getStateContent()

	// Render with breadcrumbs and help.
	view := i.CreateViewWithBreadcrumbs("Main Menu", "Burst Management", i.getStateName())
	view.WithContent(content)
	view.WithHelp(i.getContextHelp())
	return view.Render()
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

// getStateContent returns the content for the current state.
func (i *Intent) getStateContent() string {
	switch i.state {
	case StateList:
		return i.viewList()
	case StateDetail:
		return i.viewDetail()
	case StateDetailEvents:
		return i.viewDetailEvents()
	case StateDetailFacts:
		return i.viewDetailFacts()
	case StateEdit:
		return i.viewEdit()
	case StateDeleteConfirm:
		return i.viewDeleteConfirm()
	case StateConfirm:
		return i.viewConfirm()
	case StateExtractingFacts:
		return i.viewExtractingFacts()
	case StateSuggesting:
		return i.viewSuggesting()
	case StateSuggestionReview:
		return i.viewSuggestionReview()
	}
	return "Unknown state"
}

// viewList renders the burst list view.
func (i *Intent) viewList() string {
	if len(i.filteredBursts) == 0 {
		return "No bursts found."
	}
	return "Burst list view"
}

// viewDetail renders the burst detail view.
func (i *Intent) viewDetail() string {
	if i.selectedBurst == nil {
		return "No burst selected."
	}
	return "Burst detail view"
}

// viewDetailEvents renders the events view for a burst.
func (i *Intent) viewDetailEvents() string {
	return "Burst events view"
}

// viewDetailFacts renders the facts view for a burst.
func (i *Intent) viewDetailFacts() string {
	return "Burst facts view"
}

// viewEdit renders the edit view for a burst.
func (i *Intent) viewEdit() string {
	return "Edit burst view"
}

// viewDeleteConfirm renders the delete confirmation view.
func (i *Intent) viewDeleteConfirm() string {
	return "Delete confirmation view"
}

// viewConfirm renders the confirmation view.
func (i *Intent) viewConfirm() string {
	return "Confirm burst view"
}

// viewExtractingFacts renders the fact extraction progress view.
func (i *Intent) viewExtractingFacts() string {
	return "Extracting facts view"
}

// viewSuggesting renders the burst suggestion loading view.
func (i *Intent) viewSuggesting() string {
	theme := i.Theme()
	if theme == nil {
		return "Detecting burst patterns..."
	}

	// Show loading spinner/message.
	content := primitives.Title("🔍 Detecting Burst Patterns", theme).Render() + "\n\n"
	content += primitives.Body("Analyzing your events to find related patterns...\n", theme).Render()
	content += primitives.Body("This may take a moment.\n", theme).Render()

	return content
}

// viewSuggestionReview renders the suggestion review view.
func (i *Intent) viewSuggestionReview() string {
	if len(i.suggestions) == 0 {
		return "No suggestions available"
	}

	// Show current suggestion.
	currentIdx := i.currentSuggestionIdx
	if currentIdx >= len(i.suggestions) {
		currentIdx = 0
	}

	suggestion := i.suggestions[currentIdx]

	theme := i.Theme()

	// Build content with or without theme.
	var content string
	if theme != nil {
		content = primitives.Title(fmt.Sprintf("Burst Suggestion %d of %d", currentIdx+1, len(i.suggestions)), theme).Render() + "\n\n"
		content += primitives.Body(fmt.Sprintf("Name: %s\n", suggestion.Name), theme).Render()
		content += primitives.Body(fmt.Sprintf("Description: %s\n", suggestion.Description), theme).Render()
		content += primitives.Body(fmt.Sprintf("Events: %d\n", len(suggestion.EventIDs)), theme).Render()
		content += primitives.Body(fmt.Sprintf("Confidence: %.1f%%\n", suggestion.ConfidenceScore*100), theme).Render()
		content += "\n"
		content += primitives.Body("Press 'a' to accept, 'r' to reject, 'n' for next, 'p' for previous, 'esc' to cancel", theme).Render()
	} else {
		// Fallback without theme (e.g., in tests).
		content = fmt.Sprintf("Burst Suggestion %d of %d\n\n", currentIdx+1, len(i.suggestions))
		content += fmt.Sprintf("Name: %s\n", suggestion.Name)
		content += fmt.Sprintf("Description: %s\n", suggestion.Description)
		content += fmt.Sprintf("Events: %d\n", len(suggestion.EventIDs))
		content += fmt.Sprintf("Confidence: %.1f%%\n", suggestion.ConfidenceScore*100)
		content += "\n"
		content += "Press 'a' to accept, 'r' to reject, 'n' for next, 'p' for previous, 'esc' to cancel"
	}

	return content
}
