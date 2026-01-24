// Package burst_management implements the BurstManagement intent for managing career bursts.
package burst_management

import (
	"context"
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/screens"
	burstscreens "github.com/baphled/kariya/internal/cli/screens/burst_management"
	burstmodals "github.com/baphled/kariya/internal/cli/screens/burst_management/modals"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	tea "github.com/charmbracelet/bubbletea"
)

// burstRowFormatter formats a burst for table display.
func burstRowFormatter(burst *career.Burst, index int) []string {
	// Column 1: Name (truncate to 27 chars).
	nameStr := burst.Name
	if len(nameStr) > 27 {
		nameStr = nameStr[:27] + "..."
	}

	// Column 2: Description (truncated preview, max 32 chars).
	descStr := strings.TrimSpace(burst.Description)
	descStr = strings.ReplaceAll(descStr, "\n", " ")
	descStr = strings.ReplaceAll(descStr, "\r", " ")
	if descStr == "" {
		descStr = "-"
	} else if len(descStr) > 32 {
		descStr = descStr[:32] + "..."
	}

	// Column 3: Confirmed Status.
	confirmedStr := "✗ No"
	if burst.Confirmed {
		confirmedStr = "✓ Yes"
	}

	// Column 4: Event Count.
	eventCount := fmt.Sprintf("%d", len(burst.EventIDs))

	// Column 5: Created Date (YYYY-MM-DD).
	createdStr := burst.CreatedAt.Format("2006-01-02")

	return []string{nameStr, descStr, confirmedStr, eventCount, createdStr}
}

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
	// For now, return empty list.
	// TODO: Implement fact loading when burst-fact relationship is defined.
	// This might involve querying facts by burst ID or event IDs.
	_ = ctx
	return []*career.Fact{}
}

// confirmBurst marks the selected burst as confirmed and saves it.
func (i *Intent) confirmBurst() tea.Cmd {
	if i.selectedBurst == nil {
		return nil
	}

	// Mark burst as confirmed.
	i.selectedBurst.Confirmed = true

	// Save to repository if available.
	if i.context.BurstRepository != nil {
		ctx := i.getContext()
		err := i.context.BurstRepository.Update(ctx, i.selectedBurst)
		if err != nil {
			i.confirmError = err
			i.errorModal = feedback.NewErrorModal("Confirmation Failed", err.Error())
			return nil
		}
	}

	// Return to detail modal showing updated burst.
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

	// Update burst with new values.
	i.selectedBurst.Name = msg.Name
	i.selectedBurst.Description = msg.Description

	// Save to repository if available.
	if i.context.BurstRepository != nil {
		ctx := i.getContext()
		err := i.context.BurstRepository.Update(ctx, i.selectedBurst)
		if err != nil {
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
func (i *Intent) showBurstDetailModal(burst *career.Burst) tea.Cmd {
	i.selectedBurst = burst
	i.viewedBursts = append(i.viewedBursts, burst)
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

	// Mark as loading.
	i.suggestionsLoading = true

	return func() tea.Msg {
		ctx := i.getContext()

		// Get all events from the service.
		// List all events to get their IDs (use empty filters to get all).
		events, err := i.context.Service.ListEvents(ctx, careerrepo.ListFilters{})
		if err != nil {
			return BurstSuggestionsLoadedMsg{
				Suggestions: nil,
				Error:       fmt.Errorf("failed to load events: %w", err),
			}
		}

		// Extract event IDs.
		eventIDs := make([]string, len(events))
		for idx, event := range events {
			eventIDs[idx] = event.ID
		}

		// If no events, return empty suggestions.
		if len(eventIDs) == 0 {
			return BurstSuggestionsLoadedMsg{
				Suggestions: nil,
				Error:       fmt.Errorf("no events available for burst detection"),
			}
		}

		// Call the service to detect bursts.
		suggestions, err := i.context.Service.SuggestBursts(ctx, eventIDs)
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

	if msg.Error != nil {
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

	// Show suggestion review modal.
	i.suggestionModal = burstmodals.NewSuggestionReviewModal(msg.Suggestions, i.Theme())
	i.suggestionModal.Show()
	return nil
}

// handleSuggestionReviewComplete handles the SuggestionReviewCompleteMsg.
func (i *Intent) handleSuggestionReviewComplete(msg SuggestionReviewCompleteMsg) tea.Cmd {
	if msg.Cancelled || len(msg.AcceptedSuggestions) == 0 {
		// No suggestions accepted or user cancelled - return to list.
		i.transitionToScreen(burstscreens.NewBurstListScreen(i.filteredBursts))
		return nil
	}

	// Create bursts from accepted suggestions.
	for _, suggestion := range msg.AcceptedSuggestions {
		burst := &career.Burst{
			ID:          "", // Will be generated by repository
			Name:        suggestion.Name,
			Description: suggestion.Description,
			EventIDs:    suggestion.EventIDs,
			Confirmed:   false,
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

		// Add to filtered bursts list.
		i.filteredBursts = append(i.filteredBursts, burst)
		i.context.Bursts = append(i.context.Bursts, burst)
	}

	// Refresh list screen with new bursts.
	i.transitionToScreen(burstscreens.NewBurstListScreen(i.filteredBursts))
	return nil
}
