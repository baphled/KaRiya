// Package burst_management implements the BurstManagement intent for managing career bursts.
package burst_management

import (
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
		// Intercept 'a' key to save burst and extract facts immediately.
		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "a" {
			// Get current suggestion BEFORE modal removes it from the list.
			currentSuggestion := i.suggestionModal.GetCurrentSuggestion()
			if currentSuggestion != nil {
				// Save burst and trigger fact extraction immediately.
				cmd := i.saveAndExtractBurst(*currentSuggestion)

				// Now let modal update its internal state.
				i.suggestionModal.Update(msg)

				// If modal closed (no more suggestions), return to list view.
				if !i.suggestionModal.IsVisible() {
					i.suggestionModal = nil
					i.state = StateList
					i.transitionToScreen(burstscreens.NewBurstListScreen(i.filteredBursts))
				}

				return cmd
			}
		}

		// Handle other keys normally.
		_, cmd := i.suggestionModal.Update(msg)
		if !i.suggestionModal.IsVisible() {
			// Modal was closed (Esc or rejected all).
			accepted := i.suggestionModal.GetAcceptedSuggestions()
			cancelled := i.suggestionModal.GetAction() == burstmodals.SuggestionActionCancel

			// Clear modal.
			i.suggestionModal = nil

			// If user accepted some suggestions before cancelling, create them.
			if len(accepted) > 0 && !cancelled {
				completeMsg := SuggestionReviewCompleteMsg{
					AcceptedSuggestions: accepted,
					Cancelled:           false,
				}
				return tea.Batch(cmd, func() tea.Msg { return completeMsg })
			}

			// No suggestions or cancelled - return to list.
			// Return noopCmd to prevent the esc key from propagating to the list screen.
			i.state = StateList
			i.transitionToScreen(burstscreens.NewBurstListScreen(i.filteredBursts))
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

// handleFactExtractionComplete handles the FactExtractionCompleteMsg.
func (i *Intent) handleFactExtractionComplete(msg FactExtractionCompleteMsg) tea.Cmd {
	i.extractingFacts = false
	i.loadingModal = nil

	if msg.Error != nil {
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
