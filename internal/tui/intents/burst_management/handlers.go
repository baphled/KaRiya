package burst_management

import (
	"context"
	"errors"
	"fmt"
	"strings"

	domcapture "github.com/baphled/kariya/internal/domain/capture"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/burstfact"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	burstviews "github.com/baphled/kariya/internal/tui/views/burst"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/uikit/feedback"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
	tea "github.com/charmbracelet/bubbletea"
)

// noopCmd is a command that does nothing but prevents message propagation.
func noopCmd() tea.Msg { return nil }

// =============================================================================
// View Result Handlers
// =============================================================================

// handleViewResult processes a view result and dispatches to the appropriate handler.
func (i *Intent) handleViewResult(result widgets.ViewResult) tea.Cmd {
	switch result.Type() {
	case widgets.ResultCancel:
		return i.HandleCancel(result)
	case widgets.ResultNavigate:
		return i.HandleNavigate(result)
	case widgets.ResultSubmit:
		return i.HandleSubmit(result)
	case widgets.ResultError:
		return i.HandleError(result)
	default:
		return nil
	}
}

// HandleCancel processes view cancellation by navigating back through the state hierarchy.
//
// Expected:
//   - result must be a valid widgets.ViewResult of type ResultCancel.
//
// Returns:
//   - A tea.Cmd if further processing is needed, or nil when the cancellation is handled internally.
//
// Side effects:
//   - Transitions the intent state toward StateList or marks the intent as cancelled.
//   - May replace the active view with a new burst List view.
func (i *Intent) HandleCancel(result widgets.ViewResult) tea.Cmd {
	_ = result
	switch i.state {
	case StateList:
		if i.hasVisibleModal() {
			return nil
		}
		i.SetCancelled()
		return nil

	case StateDeleteConfirm:
		i.deleteModal = nil
		fallthrough
	default:
		i.state = StateList
		i.transitionToView(burstviews.NewList(display.BurstsFromDomain(i.filteredBursts)))
		return nil
	}
}

// HandleNavigate dispatches navigation results to the appropriate action or detail view.
//
// Expected:
//   - result must be a valid widgets.ViewResult of type ResultNavigate.
//   - result.Data() contains either a burst.Nav for actions, a map[string]interface{} for suggest,
//     or a *career.Burst for detail viewing.
//
// Returns:
//   - A tea.Cmd for async operations triggered by the navigation, or nil when no further action is needed.
//
// Side effects:
//   - May update selectedBurst, change intent state, or open modals depending on the navigation data.
func (i *Intent) HandleNavigate(result widgets.ViewResult) tea.Cmd {
	data := result.Data()

	if nav, ok := data.(burstviews.Nav); ok {
		if nav.Action == burstviews.ActionView {
			domainBurst := i.findBurstByID(nav.Burst.ID)
			if domainBurst == nil {
				return nil
			}
			i.selectedBurst = domainBurst
			i.AddViewedBurst(domainBurst)
			return i.showDetail(domainBurst)
		}
		return i.handleActionData(map[string]interface{}{
			"action": string(nav.Action),
			"burst":  nav.Burst,
		})
	}

	if actionData, ok := data.(map[string]interface{}); ok {
		return i.handleActionData(actionData)
	}

	if burst, ok := data.(display.Burst); ok {
		domainBurst := i.findBurstByID(burst.ID)
		if domainBurst == nil {
			return nil
		}
		i.selectedBurst = domainBurst
		i.AddViewedBurst(domainBurst)
		return i.showDetail(domainBurst)
	}

	if burst, ok := data.(*career.Burst); ok {
		i.selectedBurst = burst
		i.AddViewedBurst(burst)
		return i.showDetail(burst)
	}

	return nil
}

// HandleSubmit processes form submission results from views.
// Currently a no-op as burst management does not use view-level form submissions.
//
// Expected:
//   - result must be a valid widgets.ViewResult of type ResultSubmit.
//
// Returns:
//   - Always nil.
//
// Side effects:
//   - None.
func (i *Intent) HandleSubmit(result widgets.ViewResult) tea.Cmd {
	_ = result
	return nil
}

// HandleError captures view-level errors and presents them to the user via an error modal.
//
// Expected:
//   - result must be a valid widgets.ViewResult of type ResultError with a populated error.
//
// Returns:
//   - Always nil; the error modal is displayed on the next render cycle.
//
// Side effects:
//   - Stores the error in deleteError and creates a visible error modal overlay.
func (i *Intent) HandleError(result widgets.ViewResult) tea.Cmd {
	errResult, ok := result.(*widgets.ErrorViewResult)
	if !ok || errResult.Err == nil {
		return nil
	}
	i.deleteError = errResult.Err
	i.feedbackModal = feedback.NewErrorModal("Operation Failed", errResult.Err.Error())
	return nil
}

// =============================================================================
// Action Data Handler
// =============================================================================

// resolveBurstFromAction resolves a domain burst from action data.
// Returns the resolved burst or nil if not found.
func (i *Intent) resolveBurstFromAction(actionData map[string]interface{}) *career.Burst {
	displayBurst, hasDisplayBurst := actionData["burst"].(display.Burst)
	if !hasDisplayBurst || displayBurst.ID == "" {
		return nil
	}
	domainBurst := i.findBurstByID(displayBurst.ID)
	if domainBurst == nil && i.selectedBurst != nil && i.selectedBurst.ID == displayBurst.ID {
		domainBurst = i.selectedBurst
	}
	return domainBurst
}

// openDeleteConfirmation opens a destructive confirm modal for burst deletion.
func (i *Intent) openDeleteConfirmation(burst *career.Burst) tea.Cmd {
	i.selectedBurst = burst
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

// openBurstConfirmation opens a confirm modal for burst confirmation.
func (i *Intent) openBurstConfirmation(burst *career.Burst) tea.Cmd {
	i.selectedBurst = burst
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

// handleActionData processes action data from navigation results.
func (i *Intent) handleActionData(actionData map[string]interface{}) tea.Cmd {
	action, _ := actionData["action"].(string)
	burst := i.resolveBurstFromAction(actionData)
	if burst != nil {
		actionData["burst"] = burst
	}

	switch action {
	case "view":
		if burst, ok := actionData["burst"].(*career.Burst); ok {
			i.selectedBurst = burst
			i.AddViewedBurst(burst)
			return i.showDetail(burst)
		}
		return nil

	case "edit":
		if burst, ok := actionData["burst"].(*career.Burst); ok {
			i.selectedBurst = burst
			i.state = StateEdit
			return i.openEditModal(burst)
		}
		return nil

	case "delete":
		if burst, ok := actionData["burst"].(*career.Burst); ok {
			return i.openDeleteConfirmation(burst)
		}
		return nil

	case "suggest":
		i.state = StateSuggesting
		return i.startBurstDetection()

	case "view_events":
		if burst, ok := actionData["burst"].(*career.Burst); ok {
			i.selectedBurst = burst
			i.state = StateDetailEvents
			return i.showEvents()
		}
		return nil

	case "view_facts":
		if burst, ok := actionData["burst"].(*career.Burst); ok {
			i.selectedBurst = burst
			i.state = StateDetailFacts
			return i.showFacts()
		}
		return nil

	case "confirm":
		if burst, ok := actionData["burst"].(*career.Burst); ok {
			return i.openBurstConfirmation(burst)
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
	if cmd := i.handleFeedbackModalUpdate(msg); cmd != nil {
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
	if cmd := i.handleSuggestionEventsModalUpdate(msg); cmd != nil {
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
	if cmd := i.handleSkillsModalUpdate(msg); cmd != nil {
		return cmd
	}
	if cmd := i.handleEditModalUpdate(msg); cmd != nil {
		return cmd
	}
	return nil
}

// handleFeedbackModalUpdate handles feedback modal updates (highest priority).
func (i *Intent) handleFeedbackModalUpdate(msg tea.Msg) tea.Cmd {
	if i.feedbackModal == nil {
		return nil
	}
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.Type == tea.KeyEsc {
		i.feedbackModal = nil
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
			return tea.Batch(cmd, i.deleteBurst(i.selectedBurst))
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
			i.showDetail(i.selectedBurst)
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

	domainSuggestion := i.findBurstSuggestionByName(currentSuggestion.Name)
	if domainSuggestion == nil {
		return noopCmd
	}

	cmd := i.saveAndExtractBurst(*domainSuggestion)
	i.suggestionModal.Update(msg)

	if !i.suggestionModal.IsVisible() {
		i.clearSuggestionState()
		i.state = StateList
		i.transitionToView(burstviews.NewList(display.BurstsFromDomain(i.filteredBursts)))
	}
	return cmd
}

// handleSuggestionModalClosed handles suggestion modal closure.
func (i *Intent) handleSuggestionModalClosed(cmd tea.Cmd) tea.Cmd {
	acceptedDisplay := i.suggestionModal.GetAcceptedSuggestions()
	cancelled := i.suggestionModal.GetAction() == burstviews.SuggestionActionCancel
	accepted := make([]burstfact.BurstSuggestion, 0, len(acceptedDisplay))
	for idx := range acceptedDisplay {
		domainSuggestion := i.findBurstSuggestionByName(acceptedDisplay[idx].Name)
		if domainSuggestion != nil {
			accepted = append(accepted, *domainSuggestion)
		}
	}
	i.clearSuggestionState()

	if len(accepted) > 0 && !cancelled {
		completeMsg := SuggestionReviewCompleteMsg{
			AcceptedSuggestions: accepted,
			Cancelled:           false,
		}
		return tea.Batch(cmd, func() tea.Msg { return completeMsg })
	}

	i.state = StateList
	i.transitionToView(burstviews.NewList(display.BurstsFromDomain(i.filteredBursts)))
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
		return i.showEvents()
	case "f":
		i.detailModal.Hide()
		return i.showFacts()
	case "e":
		i.detailModal.Hide()
		return i.openEditModal(i.selectedBurst)
	case "d":
		i.detailModal.Hide()
		return i.openDeleteModal(i.selectedBurst)
	case "c":
		i.detailModal.Hide()
		return i.showConfirmBurstModal()
	case "s":
		i.detailModal.Hide()
		return i.showSkills()
	case "i":
		if i.selectedBurst != nil && !i.selectedBurst.Confirmed {
			i.ShowWarningModal("Burst Not Confirmed", "Please confirm this burst before inferring skills.")
			return noopCmd
		}
		i.detailModal.Hide()
		i.inferredFromDetail = true
		return i.startSkillInference()
	}

	if keyMsg.Type == tea.KeyEsc || keyMsg.Type == tea.KeyEnter {
		i.detailModal = nil
		return noopCmd
	}
	return nil
}

// showConfirmBurstModal starts async loading of facts to determine confirm modal text.
func (i *Intent) showConfirmBurstModal() tea.Cmd {
	if i.selectedBurst == nil {
		return nil
	}
	ctx := i.getContext()
	burstID := i.selectedBurst.ID
	service := i.context.Service
	return func() tea.Msg {
		if service == nil {
			return ConfirmBurstFactsLoadedMsg{Facts: []*career.Fact{}}
		}
		facts, err := service.GetFactsBySourceBurstID(ctx, burstID)
		if err != nil {
			return ConfirmBurstFactsLoadedMsg{Facts: []*career.Fact{}, Error: err}
		}
		return ConfirmBurstFactsLoadedMsg{Facts: facts}
	}
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
				return i.showDetail(i.selectedBurst)
			}
			return noopCmd
		}
	}
	_, cmd := i.eventsModal.Update(msg)
	return cmd
}

// handleSkillsModalUpdate handles skills modal updates.
func (i *Intent) handleSkillsModalUpdate(msg tea.Msg) tea.Cmd {
	if i.skillsModal == nil || !i.skillsModal.IsVisible() {
		return nil
	}
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.Type == tea.KeyEsc || keyMsg.Type == tea.KeyEnter {
			i.skillsModal = nil
			if i.selectedBurst != nil {
				return i.showDetail(i.selectedBurst)
			}
			return noopCmd
		}
	}
	_, cmd := i.skillsModal.Update(msg)
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
				return i.showDetail(i.selectedBurst)
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
		i.transitionToView(burstviews.NewList(display.BurstsFromDomain(i.filteredBursts)))
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
		i.feedbackModal = feedback.NewErrorModal("Extraction Failed", msg.Error.Error())
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

	targetBurst := i.selectedBurst
	if targetBurst == nil {
		targetBurst = msg.Burst
	}
	if targetBurst != nil && targetBurst.Confirmed && len(msg.Facts) > 0 {
		i.selectedBurst = targetBurst
		return i.startSkillInference()
	}

	i.state = StateList

	// Only show detail modal if a burst was selected (e.g., from confirm flow).
	// When coming from suggestion acceptance flow, selectedBurst is nil,
	// so we stay on the list view showing the newly created bursts.
	if i.selectedBurst != nil {
		return i.showDetail(i.selectedBurst)
	}

	// Refresh list screen to ensure visual update after extraction completes.
	i.transitionToView(burstviews.NewList(display.BurstsFromDomain(i.filteredBursts)))
	return nil
}

// handleEditBurstMsg handles the EditBurstMsg sent by the edit modal.
// Delegates validation and field application to the pure domain function
// capture.ApplyBurstEdit, keeping only repository persistence and UI
// transitions in the intent.
func (i *Intent) handleEditBurstMsg(msg EditBurstMsg) tea.Cmd {
	if i.selectedBurst == nil || i.selectedBurst.ID != msg.BurstID {
		i.feedbackModal = feedback.NewErrorModal("Edit Failed", "Burst not found")
		i.state = StateDetail
		return nil
	}

	originalName := i.selectedBurst.Name
	originalDescription := i.selectedBurst.Description

	input := domcapture.BurstEditInputFromFields(msg.Name, msg.Description)
	_, err := domcapture.ApplyBurstEdit(i.selectedBurst, input)
	if err != nil {
		i.state = StateEdit
		return nil
	}

	i.editModal = nil

	if i.context.BurstRepository == nil {
		return i.showDetail(i.selectedBurst)
	}

	ctx := i.getContext()
	burst := i.selectedBurst
	repo := i.context.BurstRepository
	return func() tea.Msg {
		repoErr := repo.Update(ctx, burst)
		return BurstEditCompleteMsg{
			Burst:               burst,
			Error:               repoErr,
			OriginalName:        originalName,
			OriginalDescription: originalDescription,
		}
	}
}

// handleBurstEditComplete handles the BurstEditCompleteMsg sent when a burst edit is complete.
func (i *Intent) handleBurstEditComplete(msg BurstEditCompleteMsg) tea.Cmd {
	if msg.Error != nil {
		if msg.Burst != nil && (msg.OriginalName != "" || msg.OriginalDescription != "") {
			msg.Burst.Name = msg.OriginalName
			msg.Burst.Description = msg.OriginalDescription
		}
		i.editError = msg.Error
		i.feedbackModal = feedback.NewErrorModal("Update Failed", msg.Error.Error())
		return nil
	}

	return i.showDetail(msg.Burst)
}

// handleConfirmBurstFactsLoaded handles the ConfirmBurstFactsLoadedMsg after async facts load.
func (i *Intent) handleConfirmBurstFactsLoaded(msg ConfirmBurstFactsLoadedMsg) tea.Cmd {
	if len(msg.Facts) > 0 {
		i.confirmModal = feedback.NewConfirmModal(
			"Re-extract Facts",
			fmt.Sprintf("This burst already has %d facts. Re-extract and add more?", len(msg.Facts)),
		).WithVariant(feedback.ConfirmDefault)
	} else {
		i.confirmModal = feedback.NewConfirmModal(
			"Confirm Burst",
			"Mark this burst as confirmed and extract facts?",
		).WithVariant(feedback.ConfirmDefault)
	}
	i.confirmModal.Show()
	return i.confirmModal.Init()
}

// handleBurstSaved handles the BurstSavedMsg after async burst creation.
func (i *Intent) handleBurstSaved(msg BurstSavedMsg) tea.Cmd {
	if msg.Error != nil {
		i.ShowErrorModal("Create Failed", msg.Error.Error())
		return nil
	}
	i.filteredBursts = append(i.filteredBursts, msg.Burst)
	i.context.Bursts = append(i.context.Bursts, msg.Burst)
	return i.extractFactsForBurst(msg.Burst)
}

// handleBurstConfirmed handles the BurstConfirmedMsg after async confirmation.
func (i *Intent) handleBurstConfirmed(msg BurstConfirmedMsg) tea.Cmd {
	i.confirmModal = nil
	if msg.Error != nil {
		i.confirmError = msg.Error
		i.feedbackModal = feedback.NewErrorModal("Confirmation Failed", msg.Error.Error())
		return nil
	}
	return i.extractFactsForBurst(msg.Burst)
}

// handleBurstDeleted handles the BurstDeletedMsg after async deletion.
func (i *Intent) handleBurstDeleted(msg BurstDeletedMsg) tea.Cmd {
	if msg.Error != nil {
		i.deleteError = msg.Error
		i.feedbackModal = feedback.NewErrorModal("Delete Failed", msg.Error.Error())
		i.deleteModal = nil
		i.selectedBurst = nil
		i.state = StateList
		return nil
	}
	i.context.Bursts = i.removeBurstFromSlice(i.context.Bursts, msg.BurstID)
	i.filteredBursts = i.removeBurstFromSlice(i.filteredBursts, msg.BurstID)
	i.deleteModal = nil
	i.selectedBurst = nil
	i.state = StateList
	i.transitionToView(burstviews.NewList(display.BurstsFromDomain(i.filteredBursts)))
	return nil
}

// handleBurstEventsLoaded handles the BurstEventsLoadedMsg.
func (i *Intent) handleBurstEventsLoaded(msg BurstEventsLoadedMsg) tea.Cmd {
	i.loadingEvents = false
	if !i.handleLoadError(msg.Error, "Load Events Failed") {
		return nil
	}
	i.eventsModal = burstviews.NewEvents(
		i.selectedBurst.ID, i.selectedBurst.Name, display.EventsFromDomain(msg.Events), i.Theme())
	i.showModalWithDimensions(i.eventsModal)
	return nil
}

// handleBurstSkillsLoaded handles the BurstSkillsLoadedMsg.
func (i *Intent) handleBurstSkillsLoaded(msg BurstSkillsLoadedMsg) tea.Cmd {
	if !i.handleLoadError(msg.Error, "Load Skills Failed") {
		return nil
	}
	i.skillsModal = burstviews.NewSkills(
		i.selectedBurst.ID, i.selectedBurst.Name, display.SkillsFromDomain(msg.Skills), i.Theme())
	i.showModalWithDimensions(i.skillsModal)
	return nil
}

// handleBurstFactsLoaded handles the BurstFactsLoadedMsg.
func (i *Intent) handleBurstFactsLoaded(msg BurstFactsLoadedMsg) tea.Cmd {
	i.loadingFacts = false
	if !i.handleLoadError(msg.Error, "Load Facts Failed") {
		return nil
	}
	displayFacts := make([]display.Fact, len(msg.Facts))
	for idx := range msg.Facts {
		displayFacts[idx] = display.FactFromDomain(msg.Facts[idx])
	}
	i.factsModal = burstviews.NewFacts(
		i.selectedBurst.ID, i.selectedBurst.Name, displayFacts, i.Theme())
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
		i.ShowWarningModal("No Suggestions Found",
			"No suggestions were generated from your events. "+
				"Try adding more events or adjusting detection settings.")
		i.state = StateList
		return nil
	}

	// Show suggestion review modal and update state.
	// Clear selectedBurst since we're entering suggestion review mode, not viewing a specific burst.
	// This prevents handleFactExtractionComplete from showing the detail modal for an old burst.
	i.selectedBurst = nil
	i.burstSuggestions = append([]burstfact.BurstSuggestion(nil), msg.Suggestions...)
	i.suggestionModal = burstviews.NewSuggestionReview(display.BurstSuggestionsFromDomain(msg.Suggestions), i.Theme())
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
		i.transitionToView(burstviews.NewList(display.BurstsFromDomain(i.filteredBursts)))
		return nil
	}

	createdBursts := i.createBurstsFromSuggestions(msg.AcceptedSuggestions)
	i.transitionToView(burstviews.NewList(display.BurstsFromDomain(i.filteredBursts)))

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
		i.ShowWarningModal("No Skills Detected",
			"No skills were detected from the burst events. "+
				"The events may not contain enough technical details.")
		i.state = StateList
		return nil
	}

	newSuggestions := filterNewSuggestions(msg.Suggestions, msg.ExistingSkillNames)
	if len(newSuggestions) == 0 {
		i.ShowSuccessModal("All Skills Already Tracked",
			"Detected skills already in your profile: "+strings.Join(msg.ExistingSkillNames, ", "))
		i.state = StateList
		return nil
	}

	i.skillSuggestions = append([]skillinference.SkillSuggestion(nil), newSuggestions...)
	i.skillSuggestionModal = burstviews.NewSkillSuggestion(display.SkillSuggestionsFromDomain(newSuggestions), i.Theme())
	width, height := i.getTerminalDimensions()
	i.skillSuggestionModal.SetDimensions(width, height)
	i.skillSuggestionModal.Show()
	i.state = StateSkillSuggestionReview
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

	successMsg := fmt.Sprintf("Successfully created %d skill(s)", len(msg.Skills))
	i.ShowSuccessModal("Skills Created", successMsg)
	i.state = StateList
	return nil
}

// handleSkillSuggestionModalUpdate handles skill suggestion modal updates.
func (i *Intent) handleSkillSuggestionModalUpdate(msg tea.Msg) tea.Cmd {
	if i.skillSuggestionModal == nil || !i.skillSuggestionModal.IsVisible() {
		return nil
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "a" {
		return i.handleSkillSuggestionAccept(msg)
	}

	_, cmd := i.skillSuggestionModal.Update(msg)

	if !i.skillSuggestionModal.IsVisible() {
		action := i.skillSuggestionModal.GetAction()

		switch action {
		case burstviews.SuggestionActionReject:
			if !i.skillSuggestionModal.IsVisible() {
				i.skillSuggestionModal = nil
				i.state = StateList
			}
			return cmd

		case burstviews.SuggestionActionViewEvents:
			return i.openSuggestionEventsModal()

		case burstviews.SuggestionActionCancel:
			i.skillSuggestionModal = nil
			if i.inferredFromDetail && i.selectedBurst != nil {
				i.inferredFromDetail = false
				i.state = StateList
				i.showDetail(i.selectedBurst)
				return noopCmd
			}
			i.inferredFromDetail = false
			i.state = StateList
			return noopCmd
		}
	}

	return cmd
}

// handleSkillSuggestionAccept creates async command for accepted skill, mirroring burst acceptance.
func (i *Intent) handleSkillSuggestionAccept(msg tea.Msg) tea.Cmd {
	selected := i.skillSuggestionModal.GetCurrentSkill()
	if selected == nil {
		return noopCmd
	}

	domainSuggestion := i.findSkillSuggestionByName(selected.Name)
	if domainSuggestion == nil {
		return noopCmd
	}

	cmd := i.saveSkillFromSuggestion(*domainSuggestion)
	i.skillSuggestionModal.Update(msg)

	if !i.skillSuggestionModal.IsVisible() {
		acceptedDisplay := i.skillSuggestionModal.GetAcceptedSkills()
		accepted := make([]skillinference.SkillSuggestion, 0, len(acceptedDisplay))
		for idx := range acceptedDisplay {
			domainSuggestion := i.findSkillSuggestionByName(acceptedDisplay[idx].Name)
			if domainSuggestion != nil {
				accepted = append(accepted, *domainSuggestion)
			}
		}
		successMsg := fmt.Sprintf("Successfully created %d skill(s)", len(accepted))
		i.ShowSuccessModal("Skills Created", successMsg)
		i.skillSuggestionModal = nil
		i.state = StateList
	}
	return cmd
}

// handleSuggestionEventsModalUpdate handles updates when the suggestion events modal is visible.
func (i *Intent) handleSuggestionEventsModalUpdate(msg tea.Msg) tea.Cmd {
	if i.suggestionEventsModal == nil || !i.suggestionEventsModal.IsVisible() {
		return nil
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.Type == tea.KeyEsc {
			i.suggestionEventsModal = nil
			i.skillSuggestionModal.Show()
			return noopCmd
		}
	}

	_, cmd := i.suggestionEventsModal.Update(msg)
	return cmd
}
