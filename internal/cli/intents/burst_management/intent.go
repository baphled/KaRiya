// Package burst_management implements the BurstManagement intent for managing career bursts.
package burst_management

import (
	"fmt"

	"github.com/baphled/kariya/internal/cli/screens"
	burstscreens "github.com/baphled/kariya/internal/cli/screens/burst_management"
	burstmodals "github.com/baphled/kariya/internal/cli/screens/burst_management/modals"
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
	case SuggestionReviewCompleteMsg:
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

	// Loading modal - cancellable with Esc.
	if i.loadingModal != nil {
		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.Type == tea.KeyEsc {
			// Cancel the loading operation.
			i.loadingModal = nil
			i.suggestionsLoading = false
			i.extractingFacts = false
			i.state = StateList
			return noopCmd
		}
		// Loading modal is visible - consume all other messages.
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

	// Suggestion review modal.
	if i.suggestionModal != nil && i.suggestionModal.IsVisible() {
		_, cmd := i.suggestionModal.Update(msg)
		if !i.suggestionModal.IsVisible() {
			// Modal was closed - send completion message.
			accepted := i.suggestionModal.GetAcceptedSuggestions()
			cancelled := i.suggestionModal.GetAction() == burstmodals.SuggestionActionCancel

			completeMsg := SuggestionReviewCompleteMsg{
				AcceptedSuggestions: accepted,
				Cancelled:           cancelled,
			}

			// Clear modal.
			i.suggestionModal = nil

			// Send completion message to be handled.
			return tea.Batch(cmd, func() tea.Msg { return completeMsg })
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
	switch keyMsg.String() {
	case "s":
		// Trigger burst suggestion detection.
		i.state = StateSuggesting
		return i.startBurstDetection()
	}
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

// handleFactExtractionComplete handles the FactExtractionCompleteMsg.
func (i *Intent) handleFactExtractionComplete(msg FactExtractionCompleteMsg) tea.Cmd {
	i.extractingFacts = false
	i.loadingModal = nil

	if msg.Error != nil {
		i.errorModal = feedback.NewErrorModal("Extraction Failed", msg.Error.Error())
		i.state = StateList
		return nil
	}

	i.extractedFactsCount = len(msg.Facts)
	i.state = StateList
	return i.showBurstDetailModal(i.selectedBurst)
}

// updateExtractingFactsView handles the fact extraction state.
func (i *Intent) updateExtractingFactsView(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case FactExtractionCompleteMsg:
		return i.handleFactExtractionComplete(msg)

	case tea.KeyMsg:
		if msg.String() == "esc" {
			i.extractingFacts = false
			i.state = StateList
			return nil
		}
	}
	return nil
}

// updateSuggestingView handles the burst suggestion loading state.
func (i *Intent) updateSuggestingView(msg tea.Msg) tea.Cmd {
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
	baseView := view.Render()

	// Apply modal overlay if any modals are visible (loading, error, etc.).
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
	theme := i.Theme()
	if theme == nil {
		return "Extracting facts..."
	}

	content := primitives.Title("Extracting and Saving Facts", theme).Render() + "\n\n"
	content += primitives.Body("Analyzing burst events to extract facts...\n", theme).Render()
	content += primitives.Body("Facts will be saved to the database automatically.\n", theme).Render()

	return content
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
