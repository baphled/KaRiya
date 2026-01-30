package burst_management

import (
	"context"
	"errors"
	"fmt"

	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/screens"
	burstscreens "github.com/baphled/kariya/internal/cli/screens/burst_management"
	burstmodals "github.com/baphled/kariya/internal/cli/screens/burst_management/modals"
	"github.com/baphled/kariya/internal/cli/screens/facts"
	"github.com/baphled/kariya/internal/cli/screens/timeline"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

// Ensure Intent implements ScreenResultHandler interface.
var _ behaviors.ScreenResultHandler = (*Intent)(nil)

// noopCmd is a command that does nothing but prevents message propagation.
func noopCmd() tea.Msg { return nil }

// =============================================================================
// Screen Result Handlers (ScreenResultHandler interface)
// =============================================================================

// handleScreenResult processes a screen result and determines next action.
func (i *Intent) handleScreenResult(result interface{}) tea.Cmd {
	if result == nil {
		return nil
	}

	screenResult, ok := result.(screens.ScreenResult)
	if !ok {
		return nil
	}

	return behaviors.NewScreenResultDispatcher(i).Dispatch(screenResult)
}

// HandleCancel processes screen cancellation by navigating back through the state hierarchy.
//
// Expected:
//   - result must satisfy the ScreenResultHandler interface contract.
//
// Returns:
//   - A tea.Cmd if further processing is needed, or nil when the cancellation is handled internally.
//
// Side effects:
//   - Transitions the intent state toward StateList or marks the intent as cancelled.
//   - May replace the active screen with a new BurstListScreen.
//
//nolint:revive // result parameter required by ScreenResultHandler interface
func (i *Intent) HandleCancel(result *screens.CancelResult) tea.Cmd {
	_ = result
	switch i.state {
	case StateList:
		// Check if any modal is visible - if so, Esc should close the modal, not cancel the intent.
		// Use hasVisibleModal() to properly check visibility rather than just nil checks.
		// This prevents race conditions where async operations might leave stale modal references.
		if i.hasVisibleModal() {
			// Modal is active or loading - ignore cancel from screen.
			return nil
		}
		// No modals active - cancel from list returns to main menu.
		i.SetCancelled()
		return nil

	case StateDetail, StateDetailEvents, StateDetailFacts:
		// Return to list from detail views.
		i.state = StateList
		i.transitionToScreen(burstscreens.NewBurstListScreen(i.filteredBursts))
		return nil

	case StateDeleteConfirm:
		// Cancel delete returns to list.
		i.state = StateList
		i.deleteModal = nil
		i.transitionToScreen(burstscreens.NewBurstListScreen(i.filteredBursts))
		return nil

	default:
		// Default: return to list.
		i.state = StateList
		i.transitionToScreen(burstscreens.NewBurstListScreen(i.filteredBursts))
		return nil
	}
}

// HandleNavigate dispatches navigation results to the appropriate action or detail view.
//
// Expected:
//   - result must be non-nil with ResultData containing either a map[string]interface{} for actions or a *career.Burst for detail viewing.
//
// Returns:
//   - A tea.Cmd for async operations triggered by the navigation, or nil when no further action is needed.
//
// Side effects:
//   - May update selectedBurst, change intent state, or open modals depending on the navigation data.
func (i *Intent) HandleNavigate(result *screens.NavigateResult) tea.Cmd {
	// Handle action data (add, edit, delete, suggest).
	if actionData, ok := result.ResultData.(map[string]interface{}); ok {
		return i.handleActionData(actionData)
	}

	// Handle burst selection (view details).
	if burst, ok := result.ResultData.(*career.Burst); ok {
		i.selectedBurst = burst
		i.AddViewedBurst(burst)
		// Show detail modal instead of transitioning to detail screen.
		return i.showBurstDetailModal(burst)
	}

	return nil
}

// HandleSubmit processes form submission results from screens.
// Currently a no-op as burst management does not use screen-level form submissions.
//
// Expected:
//   - result must satisfy the ScreenResultHandler interface contract.
//
// Returns:
//   - Always nil.
//
// Side effects:
//   - None.
//
//nolint:revive // result parameter required by ScreenResultHandler interface
func (i *Intent) HandleSubmit(result *screens.SubmitResult) tea.Cmd {
	_ = result
	return nil
}

// HandleError captures screen-level errors and presents them to the user via an error modal.
//
// Expected:
//   - result must be non-nil with a populated Err field.
//
// Returns:
//   - Always nil; the error modal is displayed on the next render cycle.
//
// Side effects:
//   - Stores the error in deleteError and creates a visible error modal overlay.
func (i *Intent) HandleError(result *screens.ErrorResult) tea.Cmd {
	// Store error and show error modal.
	i.deleteError = result.Err
	i.errorModal = feedback.NewErrorModal("Operation Failed", result.Err.Error())
	return nil
}

// =============================================================================
// Action Data Handler
// =============================================================================

// handleActionData processes action data from navigation results.
//
//nolint:funlen // Action routing requires handling multiple action types in one function
func (i *Intent) handleActionData(actionData map[string]interface{}) tea.Cmd {
	//nolint:errcheck // Type assertion ok value intentionally ignored - empty string is acceptable default
	action, _ := actionData["action"].(string)
	switch action {
	case "edit":
		// Extract burst from action data.
		if burst, ok := actionData["burst"].(*career.Burst); ok {
			i.selectedBurst = burst
			i.state = StateEdit
			// Open edit modal using helper.
			return i.openEditModal(burst)
		}
		return nil

	case "delete":
		// Extract burst from action data.
		if burst, ok := actionData["burst"].(*career.Burst); ok {
			i.selectedBurst = burst
			// Create delete confirmation modal.
			burstName := burst.Name
			if len(burstName) > 50 {
				burstName = burstName[:47] + "..."
			}
			i.deleteModal = feedback.NewConfirmModal(
				"Delete Burst",
				fmt.Sprintf("Are you sure you want to delete '%s'?", burstName),
			).WithVariant(feedback.ConfirmDestructive)
			i.state = StateDeleteConfirm
			return i.deleteModal.Init()
		}
		return nil

	case "suggest":
		// Trigger burst suggestion AI detection.
		i.state = StateSuggesting
		// Load all events and trigger burst detection.
		return i.startBurstDetection()

	case "view_events":
		// View events in burst.
		if burst, ok := actionData["burst"].(*career.Burst); ok {
			i.selectedBurst = burst
			i.state = StateDetailEvents
			// Load events for this burst.
			i.burstEvents = i.loadBurstEvents(burst)
			// Use timeline.EventListScreen to display burst events.
			i.transitionToScreen(timeline.NewTimelineEventListScreen(i.burstEvents))
		}
		return nil

	case "view_facts":
		// View facts extracted from burst.
		if burst, ok := actionData["burst"].(*career.Burst); ok {
			i.selectedBurst = burst
			i.state = StateDetailFacts
			// Load facts for this burst.
			i.burstFacts = i.loadBurstFacts(burst)
			// Use facts.FactListScreen to display burst facts.
			i.transitionToScreen(facts.NewFactListScreen(i.burstFacts))
		}
		return nil

	case "confirm":
		// Confirm burst (mark as confirmed).
		if burst, ok := actionData["burst"].(*career.Burst); ok {
			i.selectedBurst = burst
			// Show confirmation modal.
			burstName := burst.Name
			if len(burstName) > 50 {
				burstName = burstName[:47] + "..."
			}
			i.confirmModal = feedback.NewConfirmModal(
				"Confirm Burst",
				fmt.Sprintf("Mark '%s' as confirmed? This indicates the burst is validated and complete.", burstName),
			)
			i.state = StateConfirm
			return i.confirmModal.Init()
		}
		return nil

	default:
		return nil
	}
}

// =============================================================================
// Modal Update Handlers
// =============================================================================

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
	if cmd := i.handleSkillSuggestionModalUpdate(msg); cmd != nil {
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
// Only forwards spinner tick messages to the modal; other messages are consumed
// without restarting the spinner to avoid multiple concurrent tick loops.
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
		// Consume other keys without restarting spinner ticks.
		return noopCmd
	case feedback.ModalSpinnerTickMsg:
		// Forward tick to loading modal to advance spinner.
		return i.loadingModal.Update(msg)
	default:
		// Ignore unrelated messages while loading modal is active.
		return noopCmd
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
	i.inferringSkills = false
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
			// Show detail modal and return noopCmd to prevent message propagation.
			// The message (e.g., Esc) was already consumed by the confirm modal.
			i.showBurstDetailModal(i.selectedBurst)
			return noopCmd
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
	case "i":
		i.detailModal.Hide()
		return i.startSkillInference()
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

// =============================================================================
// Keyboard Shortcut Handler
// =============================================================================

// handleKeyShortcuts handles keyboard shortcuts when no modal is active.
func (i *Intent) handleKeyShortcuts(keyMsg tea.KeyMsg) tea.Cmd {
	if keyMsg.String() == "s" {
		// Trigger burst suggestion detection.
		i.state = StateSuggesting
		return i.startBurstDetection()
	}
	return nil
}

// =============================================================================
// Message Handlers
// =============================================================================

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

	// Auto-trigger skill inference if we have a confirmed burst with facts
	if i.selectedBurst != nil && i.selectedBurst.Confirmed && len(msg.Facts) > 0 {
		// Automatically infer skills from this burst
		return i.startSkillInference()
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

// handleEditBurstMsg handles the EditBurstMsg sent by the edit modal.
func (i *Intent) handleEditBurstMsg(msg EditBurstMsg) tea.Cmd {
	if i.selectedBurst == nil || i.selectedBurst.ID != msg.BurstID {
		// Burst mismatch or nil - show error.
		i.errorModal = feedback.NewErrorModal("Edit Failed", "Burst not found")
		i.state = StateDetail
		return nil
	}

	// Validate the name is not empty.
	if msg.Name == "" {
		// Name is required - keep modal open and show error.
		// Re-open the edit modal with an error indication.
		i.state = StateEdit
		return nil
	}

	// Clear the edit modal (it's already been closed by handleModalUpdates or test).
	i.editModal = nil

	// Save original values in case we need to rollback.
	originalName := i.selectedBurst.Name
	originalDescription := i.selectedBurst.Description

	// Update burst with new values.
	i.selectedBurst.Name = msg.Name
	i.selectedBurst.Description = msg.Description

	// Save to repository if available.
	if i.context.BurstRepository != nil {
		ctx := i.getContext()
		err := i.context.BurstRepository.Update(ctx, i.selectedBurst)
		if err != nil {
			// Rollback in-memory changes on failure.
			i.selectedBurst.Name = originalName
			i.selectedBurst.Description = originalDescription
			i.editError = err
			i.errorModal = feedback.NewErrorModal("Update Failed", err.Error())
			return nil
		}
	}

	// Return to detail modal showing updated burst.
	return i.showBurstDetailModal(i.selectedBurst)
}

// handleBurstEventsLoaded handles the BurstEventsLoadedMsg.
func (i *Intent) handleBurstEventsLoaded(msg BurstEventsLoadedMsg) tea.Cmd {
	i.loadingEvents = false
	if !i.handleLoadError(msg.Error, "Load Events Failed") {
		return nil
	}
	i.eventsModal = burstmodals.NewBurstEventsModal(
		i.selectedBurst.ID, i.selectedBurst.Name, msg.Events, i.Theme())
	i.showModalWithDimensions(i.eventsModal)
	return nil
}

// handleBurstFactsLoaded handles the BurstFactsLoadedMsg.
func (i *Intent) handleBurstFactsLoaded(msg BurstFactsLoadedMsg) tea.Cmd {
	i.loadingFacts = false
	if !i.handleLoadError(msg.Error, "Load Facts Failed") {
		return nil
	}
	i.factsModal = burstmodals.NewBurstFactsModal(
		i.selectedBurst.ID, i.selectedBurst.Name, msg.Facts, i.Theme())
	i.showModalWithDimensions(i.factsModal)
	return nil
}

// handleLoadError handles load errors and returns false if processing should stop.
func (i *Intent) handleLoadError(err error, title string) bool {
	if err != nil {
		i.ShowErrorModal(title, err.Error())
		return false
	}
	return i.selectedBurst != nil
}

// handleBurstSuggestionsLoaded handles the BurstSuggestionsLoadedMsg.
func (i *Intent) handleBurstSuggestionsLoaded(msg BurstSuggestionsLoadedMsg) tea.Cmd {
	i.suggestionsLoading = false
	i.loadingModal = nil

	if msg.Error != nil {
		// Silently ignore cancelled operations - user already knows they cancelled.
		if errors.Is(msg.Error, context.Canceled) {
			i.state = StateList
			return nil
		}
		i.suggestionsError = msg.Error
		i.ShowErrorModal("Burst Detection Failed", msg.Error.Error())
		i.state = StateList
		return nil
	}

	if len(msg.Suggestions) == 0 {
		// No suggestions found - show message and return to list.
		i.ShowErrorModal("No Suggestions Found",
			"No suggestions were generated from your events. "+
				"Try adding more events or adjusting detection settings.")
		i.state = StateList
		return nil
	}

	// Show suggestion review modal and update state.
	// Clear selectedBurst since we're entering suggestion review mode, not viewing a specific burst.
	// This prevents handleFactExtractionComplete from showing the detail modal for an old burst.
	i.selectedBurst = nil
	i.suggestionModal = burstmodals.NewSuggestionReviewModal(msg.Suggestions, i.Theme())
	width, height := i.getTerminalDimensions()
	i.suggestionModal.SetDimensions(width, height)
	i.suggestionModal.Show()
	i.state = StateSuggestionReview
	return nil
}

// handleSuggestionReviewComplete handles the SuggestionReviewCompleteMsg.
func (i *Intent) handleSuggestionReviewComplete(msg SuggestionReviewCompleteMsg) tea.Cmd {
	if !i.isInSuggestionState() {
		return nil
	}

	if len(msg.AcceptedSuggestions) == 0 {
		i.state = StateList
		i.transitionToScreen(burstscreens.NewBurstListScreen(i.filteredBursts))
		return nil
	}

	createdBursts := i.createBurstsFromSuggestions(msg.AcceptedSuggestions)
	i.transitionToScreen(burstscreens.NewBurstListScreen(i.filteredBursts))

	return i.startFactExtractionForBursts(createdBursts)
}

// =============================================================================
// Skill Inference Message Handlers
// =============================================================================

// handleSkillSuggestionsLoaded handles the SkillSuggestionsLoadedMsg.
func (i *Intent) handleSkillSuggestionsLoaded(msg SkillSuggestionsLoadedMsg) tea.Cmd {
	i.inferringSkills = false
	i.loadingModal = nil

	if msg.Error != nil {
		// Silently ignore cancelled operations.
		if errors.Is(msg.Error, context.Canceled) {
			i.state = StateList
			return nil
		}
		i.skillInferenceError = msg.Error
		i.ShowErrorModal("Skill Inference Failed", msg.Error.Error())
		i.state = StateList
		return nil
	}

	if len(msg.Suggestions) == 0 {
		// No skills detected - show message and return to list.
		i.ShowErrorModal("No Skills Detected",
			"No skills were detected from the burst events. "+
				"The events may not contain enough technical details.")
		i.state = StateList
		return nil
	}

	// Show skill suggestion modal.
	i.skillSuggestionModal = burstmodals.NewSkillSuggestionModal(msg.Suggestions, i.Theme())
	width, height := i.getTerminalDimensions()
	i.skillSuggestionModal.SetDimensions(width, height)
	i.skillSuggestionModal.Show()
	i.state = StateSkillSuggestionReview
	return nil
}

// handleSkillSuggestionsError handles the SkillSuggestionsErrorMsg.
func (i *Intent) handleSkillSuggestionsError(msg SkillSuggestionsErrorMsg) tea.Cmd {
	i.inferringSkills = false
	i.loadingModal = nil

	// Silently ignore cancelled operations.
	if errors.Is(msg.Err, context.Canceled) {
		i.state = StateList
		return nil
	}

	i.skillInferenceError = msg.Err
	i.ShowErrorModal("Skill Inference Failed", msg.Err.Error())
	i.state = StateList
	return nil
}

// handleSkillsCreated handles the SkillsCreatedMsg.
func (i *Intent) handleSkillsCreated(msg SkillsCreatedMsg) tea.Cmd {
	i.loadingModal = nil

	if msg.Error != nil {
		// Silently ignore cancelled operations.
		if errors.Is(msg.Error, context.Canceled) {
			i.state = StateList
			return nil
		}
		i.ShowErrorModal("Skill Creation Failed", msg.Error.Error())
		i.state = StateList
		return nil
	}

	// Skills created successfully - show success message and return to list.
	successMsg := fmt.Sprintf("Successfully created %d skill(s)", len(msg.Skills))
	i.ShowErrorModal("Skills Created", successMsg)
	i.state = StateList
	return nil
}

// handleSkillSuggestionModalUpdate handles skill suggestion modal updates.
func (i *Intent) handleSkillSuggestionModalUpdate(msg tea.Msg) tea.Cmd {
	if i.skillSuggestionModal == nil || !i.skillSuggestionModal.IsVisible() {
		return nil
	}

	_, cmd := i.skillSuggestionModal.Update(msg)

	// Check if modal was closed (action != "").
	if !i.skillSuggestionModal.IsVisible() {
		action := i.skillSuggestionModal.GetAction()

		switch action {
		case burstmodals.SuggestionActionAccept, burstmodals.SuggestionActionReject:
			// Individual accept/reject - keep modal open until all processed.
			// Modal automatically closes when no suggestions remain.
			if !i.skillSuggestionModal.IsVisible() {
				// All suggestions processed - create skills from accepted.
				accepted := i.skillSuggestionModal.GetAcceptedSkills()

				if len(accepted) > 0 {
					// Create skills from accepted suggestions.
					i.loadingModal = feedback.NewLoadingModal(
						fmt.Sprintf("Creating %d skill(s)...", len(accepted)),
						true,
					).WithTheme(i.Theme())
					i.skillSuggestionModal = nil
					return tea.Batch(cmd, i.createSkillsFromSuggestions(accepted))
				}

				// No accepted suggestions - just return to list.
				i.skillSuggestionModal = nil
				i.state = StateList
			}
			return cmd

		case burstmodals.SuggestionActionCancel:
			// User cancelled - return to list.
			i.skillSuggestionModal = nil
			i.state = StateList
			return cmd
		}
	}

	return cmd
}
