// Package burst_management implements the BurstManagement intent for managing career bursts.
package burst_management

import (
	"context"
	"fmt"

	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/screens"
	burstscreens "github.com/baphled/kariya/internal/cli/screens/burst_management"
	burstmodals "github.com/baphled/kariya/internal/cli/screens/burst_management/modals"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/service/career/burst_fact"
	tea "github.com/charmbracelet/bubbletea"
)

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
func (i *Intent) transitionToScreen(screen screens.Screen) {
	i.activeScreen = screen

	// Set terminal info if available.
	termInfo := i.GetTerminalInfo()
	if termInfo != nil && termInfo.Width > 0 && termInfo.Height > 0 {
		screen.SetTerminalInfo(termInfo.Width, termInfo.Height)
	}

	// Set theme if available.
	if theme := i.Theme(); theme != nil {
		screen.SetTheme(theme)
	}

	// Set logo if available.
	if logo := i.GetLogo(); logo != nil {
		screen.SetLogo(logo, i.GetLogoSpacing())
	}
}

// clearSuggestionState clears state related to suggestion review.
// This is called when exiting suggestion review to prevent race conditions
// with background fact extraction operations.
func (i *Intent) clearSuggestionState() {
	i.suggestionModal = nil
	i.suggestionsLoading = false
	// Clear selected burst since we're returning from suggestion review.
	// This prevents handleFactExtractionComplete from showing detail modal.
	i.selectedBurst = nil
}

// hasVisibleModal returns true if any modal is currently visible.
// This uses IsVisible() checks rather than nil checks to handle
// cases where modals are created but not yet shown or already hidden.
func (i *Intent) hasVisibleModal() bool {
	if i.detailModal != nil && i.detailModal.IsVisible() {
		return true
	}
	if i.eventsModal != nil && i.eventsModal.IsVisible() {
		return true
	}
	if i.factsModal != nil && i.factsModal.IsVisible() {
		return true
	}
	if i.editModal != nil && i.editModal.IsVisible() {
		return true
	}
	if i.deleteModal != nil && i.deleteModal.IsVisible() {
		return true
	}
	if i.confirmModal != nil && i.confirmModal.IsVisible() {
		return true
	}
	if i.errorModal != nil {
		return true
	}
	if i.loadingModal != nil {
		return true
	}

	// Check loading flags (these indicate async operations in progress).
	if i.loadingEvents || i.loadingFacts {
		return true
	}

	return false
}

// deleteBurst deletes a burst and refreshes the list.
func (i *Intent) deleteBurst(burst *career.Burst) tea.Cmd {
	if burst == nil {
		return nil
	}

	ctx := i.getContext()
	if i.context.BurstRepository != nil {
		err := i.context.BurstRepository.Delete(ctx, burst.ID)
		if err != nil {
			i.deleteError = err
			i.errorModal = feedback.NewErrorModal("Delete Failed", err.Error())
			i.deleteModal = nil
			i.selectedBurst = nil
			i.state = StateList
			return nil
		}
	}

	// Remove burst from lists.
	deletedID := burst.ID
	newBursts := make([]*career.Burst, 0, len(i.context.Bursts)-1)
	for _, b := range i.context.Bursts {
		if b.ID != deletedID {
			newBursts = append(newBursts, b)
		}
	}
	i.context.Bursts = newBursts

	// Remove from filtered list.
	newFiltered := make([]*career.Burst, 0, len(i.filteredBursts)-1)
	for _, b := range i.filteredBursts {
		if b.ID != deletedID {
			newFiltered = append(newFiltered, b)
		}
	}
	i.filteredBursts = newFiltered

	// Reset state and refresh list screen.
	i.deleteModal = nil
	i.selectedBurst = nil
	i.state = StateList
	i.transitionToScreen(burstscreens.NewBurstListScreen(i.filteredBursts))

	return nil
}

// RefreshData reloads/refreshes the filtered data.
func (i *Intent) RefreshData() tea.Cmd {
	// Reload bursts from context if available.
	if err := i.context.LoadBursts(); err != nil {
		// Error loading - keep existing data.
		_ = err
	}

	// Update filtered bursts.
	i.filteredBursts = i.context.Bursts

	// Refresh the list screen.
	i.transitionToScreen(burstscreens.NewBurstListScreen(i.filteredBursts))
	return nil
}

// loadBurstEvents loads career events for the given burst.
func (i *Intent) loadBurstEvents(burst *career.Burst) []*career.CareerEvent {
	if burst == nil || i.context.Service == nil {
		return []*career.CareerEvent{}
	}

	ctx := i.getContext()
	events := make([]*career.CareerEvent, 0, len(burst.EventIDs))

	for _, eventID := range burst.EventIDs {
		event, err := i.context.Service.GetEventByID(ctx, eventID)
		if err == nil && event != nil {
			events = append(events, event)
		}
	}

	return events
}

// loadBurstFacts loads facts for the given burst.
func (i *Intent) loadBurstFacts(burst *career.Burst) []*career.Fact {
	if burst == nil || i.context.Service == nil {
		return []*career.Fact{}
	}

	ctx := i.getContext()
	facts, err := i.context.Service.GetFactsBySourceBurstID(ctx, burst.ID)
	if err != nil {
		return []*career.Fact{}
	}

	return facts
}

// confirmBurst marks the selected burst as confirmed and saves it.
// Note: Fact extraction happens when accepting suggestions, not here.
func (i *Intent) confirmBurst() tea.Cmd {
	if i.selectedBurst == nil {
		return nil
	}

	// Save to repository first (if available) BEFORE modifying in-memory state.
	if i.context.BurstRepository != nil {
		// Create a copy with Confirmed=true for the update.
		ctx := i.getContext()
		i.selectedBurst.Confirmed = true
		err := i.context.BurstRepository.Update(ctx, i.selectedBurst)
		if err != nil {
			// Rollback in-memory change on failure.
			i.selectedBurst.Confirmed = false
			i.confirmError = err
			i.errorModal = feedback.NewErrorModal("Confirmation Failed", err.Error())
			return nil
		}
	} else {
		// No repository - just update in memory.
		i.selectedBurst.Confirmed = true
	}

	// Stay on detail view - just update the confirmed status visually.
	i.state = StateDetail
	return i.showBurstDetailModal(i.selectedBurst)
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

// HasVisibleErrorModal returns true if the error modal is visible.
func (i *Intent) HasVisibleErrorModal() bool {
	return i.errorModal != nil
}

// HasVisibleDeleteModal returns true if the delete confirmation modal is visible.
func (i *Intent) HasVisibleDeleteModal() bool {
	return i.deleteModal != nil && i.deleteModal.IsVisible()
}

// HasVisibleConfirmModal returns true if the confirm burst modal is visible.
func (i *Intent) HasVisibleConfirmModal() bool {
	return i.confirmModal != nil && i.confirmModal.IsVisible()
}

// HasVisibleEditModal returns true if the edit burst modal is visible.
func (i *Intent) HasVisibleEditModal() bool {
	return i.editModal != nil && i.editModal.IsVisible()
}

// ShowErrorModal creates and shows an error modal with the given title and message.
func (i *Intent) ShowErrorModal(title, message string) {
	i.errorModal = feedback.NewErrorModal(title, message)
}

// GetTerminalDimensions returns current terminal dimensions (exported for testing).
func (i *Intent) GetTerminalDimensions() (width, height int) {
	return i.getTerminalDimensions()
}

// GetStateName returns human-readable state name (exported for testing).
func (i *Intent) GetStateName() string {
	return i.getStateName()
}

// GetContextHelp returns context help string (exported for testing).
func (i *Intent) GetContextHelp() string {
	return i.getContextHelp()
}

// showBurstDetailModal creates and shows the burst detail modal.
// Note: viewedBursts tracking is handled by HandleNavigate, not here.
func (i *Intent) showBurstDetailModal(burst *career.Burst) tea.Cmd {
	i.selectedBurst = burst
	width, height := i.getTerminalDimensions()
	theme := i.Theme()

	modal := burstmodals.NewBurstDetailModal(burst, theme)
	modal.SetDimensions(width, height)
	modal.Show()
	i.detailModal = modal

	return nil
}

// showBurstEventsModal loads and shows events for the current burst.
func (i *Intent) showBurstEventsModal() tea.Cmd {
	if i.selectedBurst == nil {
		return nil
	}

	// Check for nil service before starting async load.
	if i.context.Service == nil {
		// No service available - return empty events immediately.
		return func() tea.Msg {
			return BurstEventsLoadedMsg{Events: []*career.CareerEvent{}}
		}
	}

	i.loadingEvents = true
	return func() tea.Msg {
		events := make([]*career.CareerEvent, 0, len(i.selectedBurst.EventIDs))
		for _, eventID := range i.selectedBurst.EventIDs {
			event, err := i.context.Service.GetEventByID(i.getContext(), eventID)
			if err != nil {
				continue
			}
			events = append(events, event)
		}
		return BurstEventsLoadedMsg{Events: events}
	}
}

// showBurstFactsModal loads and shows facts for the current burst.
func (i *Intent) showBurstFactsModal() tea.Cmd {
	if i.selectedBurst == nil {
		return nil
	}

	// Check for nil service before starting async load.
	if i.context.Service == nil {
		// No service available - return empty facts immediately.
		return func() tea.Msg {
			return BurstFactsLoadedMsg{Facts: []*career.Fact{}}
		}
	}

	i.loadingFacts = true
	return func() tea.Msg {
		facts, err := i.context.Service.GetFactsBySourceBurstID(
			i.getContext(),
			i.selectedBurst.ID,
		)
		if err != nil {
			return BurstFactsLoadedMsg{Error: err}
		}
		return BurstFactsLoadedMsg{Facts: facts}
	}
}

// openEditModal creates and shows the edit modal for a burst.
func (i *Intent) openEditModal(burst *career.Burst) tea.Cmd {
	width, height := i.getTerminalDimensions()
	i.editModal = burstmodals.NewEditBurstModal(burst, width, height)
	return i.editModal.Init()
}

// openDeleteModal creates and shows the delete confirmation modal.
func (i *Intent) openDeleteModal(burst *career.Burst) tea.Cmd {
	burstName := burst.Name
	if len(burstName) > 50 {
		burstName = burstName[:47] + "..."
	}
	i.deleteModal = feedback.NewConfirmModal(
		"Delete Burst",
		fmt.Sprintf("Are you sure you want to delete '%s'?", burstName),
	).WithVariant(feedback.ConfirmDestructive)
	i.selectedBurst = burst
	return i.deleteModal.Init()
}

// updateDetailModalRegistry updates the modal registry with detail modals.
func (i *Intent) updateDetailModalRegistry() {
	// Detail modals (detail, events, facts).
	if i.detailModal != nil && i.detailModal.IsVisible() {
		i.modalRegistry.Register(intents.NewViewModalAdapter(
			i.detailModal.IsVisible,
			i.detailModal.View,
			i.detailModal.Update,
		))
	}

	if i.eventsModal != nil && i.eventsModal.IsVisible() {
		i.modalRegistry.Register(intents.NewViewModalAdapter(
			i.eventsModal.IsVisible,
			i.eventsModal.View,
			i.eventsModal.Update,
		))
	}

	if i.factsModal != nil && i.factsModal.IsVisible() {
		i.modalRegistry.Register(intents.NewViewModalAdapter(
			i.factsModal.IsVisible,
			i.factsModal.View,
			i.factsModal.Update,
		))
	}
}

// RebuildModalRegistry creates a fresh modal registry with all current modals (exported for testing).
func (i *Intent) RebuildModalRegistry() {
	i.rebuildModalRegistry()
}

// rebuildModalRegistry creates a fresh modal registry with all current modals.
// Call this whenever a modal is created or destroyed to keep the registry current.
func (i *Intent) rebuildModalRegistry() {
	if i.modalRegistry == nil {
		i.modalRegistry = intents.NewModalRegistry()
	}
	i.modalRegistry.Clear()

	// Register modals in priority order (highest priority first).
	// Error modal has highest priority.
	if i.errorModal != nil {
		width, height := i.getTerminalDimensions()
		i.modalRegistry.Register(intents.NewErrorModalAdapter(i.errorModal, width, height, i.Theme()))
	}

	// Loading modal (for StateExtractingFacts and StateSuggesting).
	// When loading is active, it should be the only modal visible.
	if i.loadingModal != nil {
		width, height := i.getTerminalDimensions()
		i.modalRegistry.Register(intents.NewErrorModalAdapter(i.loadingModal, width, height, i.Theme()))
		// Don't register other modals when loading - loading takes full precedence.
		return
	}

	// Detail modals (detail, events, facts).
	i.updateDetailModalRegistry()

	// Edit burst modal.
	if i.editModal != nil && i.editModal.IsVisible() {
		i.modalRegistry.Register(intents.NewFormModalAdapter(
			i.editModal.IsVisible,
			i.editModal.View,
			i.editModal.Update,
		))
	}

	// Delete confirmation modal.
	if i.deleteModal != nil {
		i.modalRegistry.Register(intents.NewConfirmModalAdapter(i.deleteModal))
	}

	// Confirm burst modal.
	if i.confirmModal != nil {
		i.modalRegistry.Register(intents.NewConfirmModalAdapter(i.confirmModal))
	}

	// Suggestion review modal.
	if i.suggestionModal != nil && i.suggestionModal.IsVisible() {
		i.modalRegistry.Register(intents.NewViewModalAdapter(
			i.suggestionModal.IsVisible,
			i.suggestionModal.View,
			i.suggestionModal.Update,
		))
	}
}

// handleBurstEventsLoaded handles the BurstEventsLoadedMsg.
func (i *Intent) handleBurstEventsLoaded(msg BurstEventsLoadedMsg) tea.Cmd {
	i.loadingEvents = false

	if msg.Error != nil {
		i.ShowErrorModal("Load Events Failed", msg.Error.Error())
		return nil
	}

	if i.selectedBurst == nil {
		return nil
	}

	// Show events modal.
	width, height := i.getTerminalDimensions()
	theme := i.Theme()
	i.eventsModal = burstmodals.NewBurstEventsModal(
		i.selectedBurst.ID,
		i.selectedBurst.Name,
		msg.Events,
		theme,
	)
	i.eventsModal.SetDimensions(width, height)
	i.eventsModal.Show()

	return nil
}

// handleBurstFactsLoaded handles the BurstFactsLoadedMsg.
func (i *Intent) handleBurstFactsLoaded(msg BurstFactsLoadedMsg) tea.Cmd {
	i.loadingFacts = false

	if msg.Error != nil {
		i.ShowErrorModal("Load Facts Failed", msg.Error.Error())
		return nil
	}

	if i.selectedBurst == nil {
		return nil
	}

	// Show facts modal.
	width, height := i.getTerminalDimensions()
	theme := i.Theme()
	i.factsModal = burstmodals.NewBurstFactsModal(
		i.selectedBurst.ID,
		i.selectedBurst.Name,
		msg.Facts,
		theme,
	)
	i.factsModal.SetDimensions(width, height)
	i.factsModal.Show()

	return nil
}

// startBurstDetection loads all events and triggers AI burst detection.
func (i *Intent) startBurstDetection() tea.Cmd {
	if i.context.Service == nil {
		i.ShowErrorModal("Detection Failed", "Service not available")
		i.state = StateList
		return nil
	}

	// Cancel any previous async operation.
	if i.cancelFunc != nil {
		i.cancelFunc()
	}

	// Create cancellable context for this operation.
	ctx, cancel := context.WithCancel(context.Background())
	i.cancelFunc = cancel

	// Mark as loading and create loading modal.
	i.suggestionsLoading = true
	i.loadingModal = feedback.NewLoadingModal("Detecting burst patterns...", true).WithTheme(i.Theme())

	// Capture existing bursts before async operation.
	existingBursts := i.context.Bursts

	// Capture service reference to avoid race conditions.
	service := i.context.Service

	return func() tea.Msg {
		// Check if cancelled before starting.
		if ctx.Err() != nil {
			return BurstSuggestionsLoadedMsg{
				Suggestions: nil,
				Error:       ctx.Err(),
			}
		}

		// Get all events from the service.
		// List all events to get their IDs (use Limit=-1 for no limit).
		events, err := service.ListEvents(ctx, careerrepo.ListFilters{Limit: -1})
		if err != nil {
			return BurstSuggestionsLoadedMsg{
				Suggestions: nil,
				Error:       fmt.Errorf("failed to load events: %w", err),
			}
		}

		// Check if cancelled after loading events.
		if ctx.Err() != nil {
			return BurstSuggestionsLoadedMsg{
				Suggestions: nil,
				Error:       ctx.Err(),
			}
		}

		// Build a set of event IDs that are already in confirmed bursts.
		usedEventIDs := make(map[string]bool)
		for _, burst := range existingBursts {
			for _, eventID := range burst.EventIDs {
				usedEventIDs[eventID] = true
			}
		}

		// Extract event IDs, excluding those already in bursts.
		var eventIDs []string
		for _, event := range events {
			if !usedEventIDs[event.ID] {
				eventIDs = append(eventIDs, event.ID)
			}
		}

		// If no events available for detection, return appropriate message.
		if len(eventIDs) == 0 {
			return BurstSuggestionsLoadedMsg{
				Suggestions: nil,
				Error:       fmt.Errorf("no unassigned events available for burst detection"),
			}
		}

		// Call the service to detect bursts from unassigned events only.
		suggestions, err := service.SuggestBursts(ctx, eventIDs)
		if err != nil {
			return BurstSuggestionsLoadedMsg{
				Suggestions: nil,
				Error:       err,
			}
		}

		return BurstSuggestionsLoadedMsg{
			Suggestions: suggestions,
			Error:       nil,
		}
	}
}

// handleBurstSuggestionsLoaded handles the BurstSuggestionsLoadedMsg.
func (i *Intent) handleBurstSuggestionsLoaded(msg BurstSuggestionsLoadedMsg) tea.Cmd {
	i.suggestionsLoading = false
	i.loadingModal = nil

	if msg.Error != nil {
		// Silently ignore cancelled operations - user already knows they cancelled.
		if msg.Error == context.Canceled {
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
		i.ShowErrorModal("No Suggestions Found", "No suggestions were generated from your events. Try adding more events or adjusting detection settings.")
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
	// Guard: Only process if we're in suggestion review state.
	// This prevents double-processing when the modal close triggers direct handling
	// and a subsequent message is sent.
	if i.state != StateSuggestionReview && i.state != StateSuggesting {
		return nil
	}

	if len(msg.AcceptedSuggestions) == 0 {
		// No suggestions accepted - return to list.
		// Note: We don't check msg.Cancelled here because user may have accepted
		// some suggestions and then pressed Esc to close. Those should still be saved.
		i.state = StateList
		i.transitionToScreen(burstscreens.NewBurstListScreen(i.filteredBursts))
		return nil
	}

	// Track successfully created bursts for fact extraction.
	createdBursts := make([]*career.Burst, 0, len(msg.AcceptedSuggestions))

	// Create bursts from accepted suggestions.
	for _, suggestion := range msg.AcceptedSuggestions {
		burst := &career.Burst{
			ID:          "", // Will be generated by repository
			Name:        suggestion.Name,
			Description: suggestion.Description,
			EventIDs:    suggestion.EventIDs,
			Confirmed:   true, // Auto-confirm accepted suggestions
		}

		// Save to repository if available.
		if i.context.BurstRepository != nil {
			ctx := i.getContext()
			err := i.context.BurstRepository.Create(ctx, burst)
			if err != nil {
				i.ShowErrorModal("Create Failed", fmt.Sprintf("Failed to create burst: %v", err))
				continue
			}
		}

		// Add to filtered bursts list and track for fact extraction.
		i.filteredBursts = append(i.filteredBursts, burst)
		i.context.Bursts = append(i.context.Bursts, burst)
		createdBursts = append(createdBursts, burst)
	}

	// Refresh list screen with new bursts.
	i.transitionToScreen(burstscreens.NewBurstListScreen(i.filteredBursts))

	// Trigger fact extraction only for successfully created bursts.
	var cmds []tea.Cmd
	for _, burst := range createdBursts {
		cmds = append(cmds, i.extractFactsForBurst(burst))
	}

	if len(cmds) > 0 {
		i.extractingFacts = true
		i.state = StateExtractingFacts
		i.loadingModal = feedback.NewLoadingModal("Extracting facts from accepted bursts...", true).WithTheme(i.Theme())
		return tea.Batch(cmds...)
	}

	// No bursts created (all failed) - stay on list.
	i.state = StateList
	return nil
}

// saveAndExtractBurst saves a burst immediately and triggers fact extraction.
// This is called when user presses 'a' to provide instant feedback and ensure
// fact extraction happens even if the app crashes before modal closes.
func (i *Intent) saveAndExtractBurst(suggestion burst_fact.BurstSuggestion) tea.Cmd {
	_, cmd := i.saveAndExtractBurstWithResult(suggestion)
	return cmd
}

// saveAndExtractBurstWithResult saves a burst and returns both the burst and the extraction command.
// This is used when we need to display the burst detail modal after saving.
func (i *Intent) saveAndExtractBurstWithResult(suggestion burst_fact.BurstSuggestion) (*career.Burst, tea.Cmd) {
	burst := &career.Burst{
		ID:          "", // Will be generated by repository
		Name:        suggestion.Name,
		Description: suggestion.Description,
		EventIDs:    suggestion.EventIDs,
		Confirmed:   true, // Auto-confirm accepted suggestions
	}

	// Save to repository if available.
	if i.context.BurstRepository != nil {
		ctx := i.getContext()
		err := i.context.BurstRepository.Create(ctx, burst)
		if err != nil {
			i.ShowErrorModal("Create Failed", fmt.Sprintf("Failed to create burst: %v", err))
			return nil, nil
		}
	}

	// Add to filtered bursts list immediately.
	i.filteredBursts = append(i.filteredBursts, burst)
	i.context.Bursts = append(i.context.Bursts, burst)

	// Trigger fact extraction immediately (runs async in background).
	return burst, i.extractFactsForBurst(burst)
}

func (i *Intent) extractFactsForBurst(burst *career.Burst) tea.Cmd {
	if burst == nil {
		return func() tea.Msg {
			return FactExtractionCompleteMsg{Error: fmt.Errorf("no burst provided")}
		}
	}

	// Cancel any previous async operation.
	if i.cancelFunc != nil {
		i.cancelFunc()
	}

	// Create cancellable context for this operation.
	ctx, cancel := context.WithCancel(context.Background())
	i.cancelFunc = cancel

	// Capture service reference to avoid race conditions.
	service := i.context.Service

	return func() tea.Msg {
		// Check if cancelled before starting.
		if ctx.Err() != nil {
			return FactExtractionCompleteMsg{Error: ctx.Err()}
		}

		if service == nil {
			return FactExtractionCompleteMsg{Error: fmt.Errorf("service not available")}
		}

		facts, err := service.ExtractFactsFromBurst(ctx, burst)
		if err != nil {
			return FactExtractionCompleteMsg{Error: err}
		}

		// Check if cancelled after extraction.
		if ctx.Err() != nil {
			return FactExtractionCompleteMsg{Error: ctx.Err()}
		}

		savedFacts := make([]*career.Fact, 0, len(facts))
		for idx := range facts {
			fact := &facts[idx]
			fact.SourceBurstID = burst.ID

			// Check if cancelled during save loop.
			if ctx.Err() != nil {
				return FactExtractionCompleteMsg{Error: ctx.Err()}
			}

			if err := service.SaveFact(ctx, fact); err != nil {
				continue
			}
			savedFacts = append(savedFacts, fact)
		}

		return FactExtractionCompleteMsg{Facts: savedFacts}
	}
}
