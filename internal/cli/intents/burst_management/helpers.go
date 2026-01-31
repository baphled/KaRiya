// Package burst_management implements the BurstManagement intent for managing career bursts.
package burst_management

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/screens"
	burstscreens "github.com/baphled/kariya/internal/cli/screens/burst_management"
	burstmodals "github.com/baphled/kariya/internal/cli/screens/burst_management/modals"
	skillmodals "github.com/baphled/kariya/internal/cli/screens/skills/modals"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/service/career/burstfact"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	tea "github.com/charmbracelet/bubbletea"
)

// filterNewSuggestions removes suggestions whose names appear in existingNames.
func filterNewSuggestions(
	suggestions []skillinference.SkillSuggestion,
	existingNames []string,
) []skillinference.SkillSuggestion {
	if len(existingNames) == 0 {
		return suggestions
	}

	existingSet := make(map[string]bool, len(existingNames))
	for _, name := range existingNames {
		existingSet[strings.ToLower(name)] = true
	}

	filtered := make([]skillinference.SkillSuggestion, 0, len(suggestions))
	for _, s := range suggestions {
		if !existingSet[strings.ToLower(s.Name)] {
			filtered = append(filtered, s)
		}
	}
	return filtered
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
	case StateInferringSkills:
		return "Inferring Skills"
	case StateSkillSuggestionReview:
		return "Review Skills"
	default:
		return "Unknown"
	}
}

// getContextHelp returns themed keyboard shortcuts for the current state.
func (i *Intent) getContextHelp() string {
	theme := i.Theme()

	// Only show badges for keys that are actually handled by this intent.
	// Global keys (q, m) are not handled, so ThemedGlobalBadges is not used.
	switch i.state {
	case StateList:
		badges := []*primitives.Badge{
			primitives.NavigateBadge(theme),
			primitives.PageBadge(theme),
			primitives.ViewBadge(theme),
			primitives.AddBadge(theme),
			primitives.EditBadge(theme),
			primitives.DeleteBadge(theme),
			primitives.SuggestBadge(theme),
			primitives.BackBadge(theme),
		}
		return intents.ThemedCustomFooter(theme, badges...)

	case StateDetail:
		badges := []*primitives.Badge{
			primitives.ViewEventsBadge(theme),
			primitives.ViewFactsBadge(theme),
			primitives.EditBadge(theme),
			primitives.DeleteBadge(theme),
			primitives.ConfirmActionBadge(theme),
			primitives.BackBadge(theme),
		}
		return intents.ThemedCustomFooter(theme, badges...)

	case StateDetailEvents, StateDetailFacts:
		badges := []*primitives.Badge{
			primitives.BackBadge(theme),
		}
		return intents.ThemedCustomFooter(theme, badges...)

	default:
		return ""
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
	return i.hasVisibleContentModal() ||
		i.hasVisibleActionModal() ||
		i.hasVisibleFeedbackModal() ||
		i.isLoadingAsync()
}

// hasVisibleContentModal checks if content display modals are visible.
func (i *Intent) hasVisibleContentModal() bool {
	return (i.detailModal != nil && i.detailModal.IsVisible()) ||
		(i.eventsModal != nil && i.eventsModal.IsVisible()) ||
		(i.factsModal != nil && i.factsModal.IsVisible()) ||
		(i.skillsModal != nil && i.skillsModal.IsVisible()) ||
		(i.suggestionEventsModal != nil && i.suggestionEventsModal.IsVisible())
}

// hasVisibleActionModal checks if action modals (edit, delete, confirm) are visible.
func (i *Intent) hasVisibleActionModal() bool {
	return (i.editModal != nil && i.editModal.IsVisible()) ||
		(i.deleteModal != nil && i.deleteModal.IsVisible()) ||
		(i.confirmModal != nil && i.confirmModal.IsVisible())
}

// hasVisibleFeedbackModal checks if feedback modals (error, loading) are active.
func (i *Intent) hasVisibleFeedbackModal() bool {
	return i.feedbackModal != nil || i.loadingModal != nil
}

// isLoadingAsync returns true if any async loading operation is in progress.
func (i *Intent) isLoadingAsync() bool {
	return i.loadingEvents || i.loadingFacts
}

// deleteBurst deletes a burst and refreshes the list.
func (i *Intent) deleteBurst(burst *career.Burst) {
	if burst == nil {
		return
	}

	ctx := i.getContext()
	if i.context.BurstRepository != nil {
		if err := i.context.BurstRepository.Delete(ctx, burst.ID); err != nil {
			i.deleteError = err
			i.feedbackModal = feedback.NewErrorModal("Delete Failed", err.Error())
			i.deleteModal = nil
			i.selectedBurst = nil
			i.state = StateList
			return
		}
	}

	// Remove burst from context and filtered lists.
	i.context.Bursts = i.removeBurstFromSlice(i.context.Bursts, burst.ID)
	i.filteredBursts = i.removeBurstFromSlice(i.filteredBursts, burst.ID)

	// Reset state and refresh list screen.
	i.deleteModal = nil
	i.selectedBurst = nil
	i.state = StateList
	i.transitionToScreen(burstscreens.NewBurstListScreen(i.filteredBursts))
}

// removeBurstFromSlice removes a burst with the given ID from the slice.
func (i *Intent) removeBurstFromSlice(bursts []*career.Burst, id string) []*career.Burst {
	result := make([]*career.Burst, 0, len(bursts))
	for _, b := range bursts {
		if b.ID != id {
			result = append(result, b)
		}
	}
	return result
}

// RefreshData reloads burst data from the repository and refreshes the list screen.
//
// Returns:
//   - Always nil; screen transitions are handled internally.
//
// Side effects:
//   - Reloads bursts from the repository via IntentContext.
//   - Replaces filteredBursts with the reloaded data.
//   - Transitions the active screen to a new BurstListScreen.
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
func (i *Intent) loadBurstEvents(burst *career.Burst) []*career.Event {
	if burst == nil || i.context.Service == nil {
		return []*career.Event{}
	}

	ctx := i.getContext()
	events := make([]*career.Event, 0, len(burst.EventIDs))

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

// confirmBurst marks the selected burst as confirmed and triggers fact extraction.
func (i *Intent) confirmBurst() tea.Cmd {
	if i.selectedBurst == nil {
		return nil
	}

	// Delegate to the service which handles Confirmed, ConfirmedAt, UpdatedAt,
	// and rolls back correctly on failure.
	if i.context.Service != nil {
		ctx := i.getContext()
		if err := i.context.Service.ConfirmBurst(ctx, i.selectedBurst); err != nil {
			i.confirmError = err
			i.feedbackModal = feedback.NewErrorModal("Confirmation Failed", err.Error())
			return nil
		}
	} else {
		now := time.Now()
		i.selectedBurst.Confirmed = true
		i.selectedBurst.ConfirmedAt = &now
		i.selectedBurst.UpdatedAt = now
	}

	// Trigger fact extraction for the confirmed burst.
	return i.extractFactsForBurst(i.selectedBurst)
}

// HasVisibleFeedbackModal returns true if the feedback modal is visible.
func (i *Intent) HasVisibleFeedbackModal() bool {
	return i.feedbackModal != nil
}

// HasVisibleErrorModal checks whether an error modal is currently displayed.
//
// Deprecated: Use HasVisibleFeedbackModal instead.
//
// Returns:
//   - True if the error modal reference is non-nil.
//
// Side effects:
//   - None.
func (i *Intent) HasVisibleErrorModal() bool {
	return i.HasVisibleFeedbackModal()
}

// GetFeedbackModal returns the current feedback modal for testing.
func (i *Intent) GetFeedbackModal() *feedback.Modal {
	return i.feedbackModal
}

// GetLoadingModal returns the current loading modal for testing.
func (i *Intent) GetLoadingModal() *feedback.Modal {
	return i.loadingModal
}

// HasVisibleDeleteModal checks whether the delete confirmation modal is currently displayed.
//
// Returns:
//   - True if the delete modal exists and reports itself as visible.
//
// Side effects:
//   - None.
func (i *Intent) HasVisibleDeleteModal() bool {
	return i.deleteModal != nil && i.deleteModal.IsVisible()
}

// HasVisibleConfirmModal checks whether the burst confirmation modal is currently displayed.
//
// Returns:
//   - True if the confirm modal exists and reports itself as visible.
//
// Side effects:
//   - None.
func (i *Intent) HasVisibleConfirmModal() bool {
	return i.confirmModal != nil && i.confirmModal.IsVisible()
}

// HasVisibleEditModal checks whether the edit burst modal is currently displayed.
//
// Returns:
//   - True if the edit modal exists and reports itself as visible.
//
// Side effects:
//   - None.
func (i *Intent) HasVisibleEditModal() bool {
	return i.editModal != nil && i.editModal.IsVisible()
}

// ShowErrorModal creates and activates an error modal overlay for the intent.
//
// Expected:
//   - title must be a non-empty human-readable heading.
//   - message must describe the error condition.
//
// Side effects:
//   - Replaces any existing error modal on the intent.
func (i *Intent) ShowErrorModal(title, message string) {
	i.feedbackModal = feedback.NewErrorModal(title, message)
}

// ShowWarningModal creates and shows a warning modal with the given title and message.
func (i *Intent) ShowWarningModal(title, message string) {
	i.feedbackModal = feedback.NewWarningModal(title, message)
}

// ShowSuccessModal creates and shows a success modal with the given message.
func (i *Intent) ShowSuccessModal(title, message string) {
	modal := feedback.NewSuccessModal(message)
	modal.Title = title
	i.feedbackModal = modal
}

// GetTerminalDimensions exposes terminal dimensions for testing purposes.
//
// Returns:
//   - width and height in columns and rows, with fallback defaults when terminal info is unavailable.
//
// Side effects:
//   - None.
func (i *Intent) GetTerminalDimensions() (width, height int) {
	return i.getTerminalDimensions()
}

// GetStateName exposes the human-readable name of the current intent state for testing purposes.
//
// Returns:
//   - A display-friendly label for the current state, or "Unknown" for unrecognized states.
//
// Side effects:
//   - None.
func (i *Intent) GetStateName() string {
	return i.getStateName()
}

// GetContextHelp exposes the themed keyboard shortcut help string for testing purposes.
//
// Returns:
//   - A rendered help string with keyboard badges for the current state, or empty for states without help.
//
// Side effects:
//   - None.
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
			return BurstEventsLoadedMsg{Events: []*career.Event{}}
		}
	}

	i.loadingEvents = true
	return func() tea.Msg {
		events := make([]*career.Event, 0, len(i.selectedBurst.EventIDs))
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

// showBurstSkillsModal loads and shows skills for the current burst.
func (i *Intent) showBurstSkillsModal() tea.Cmd {
	if i.selectedBurst == nil {
		return nil
	}

	if i.context.SkillRepository == nil {
		return func() tea.Msg {
			return BurstSkillsLoadedMsg{Skills: []*career.Skill{}}
		}
	}

	return func() tea.Msg {
		ctx := i.getContext()
		skillMap := make(map[string]*career.Skill)
		for _, eventID := range i.selectedBurst.EventIDs {
			skills, err := i.context.SkillRepository.GetSkillsForEvent(ctx, eventID)
			if err != nil {
				continue
			}
			for _, skill := range skills {
				skillMap[skill.ID] = skill
			}
		}

		dedupedSkills := make([]*career.Skill, 0, len(skillMap))
		for _, skill := range skillMap {
			dedupedSkills = append(dedupedSkills, skill)
		}
		return BurstSkillsLoadedMsg{Skills: dedupedSkills}
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

	if i.skillsModal != nil && i.skillsModal.IsVisible() {
		i.modalRegistry.Register(intents.NewViewModalAdapter(
			i.skillsModal.IsVisible,
			i.skillsModal.View,
			i.skillsModal.Update,
		))
	}
}

// RebuildModalRegistry recreates the modal registry from the current modal state, exported for testing.
//
// Side effects:
//   - Clears and repopulates the modal registry based on all active modals.
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
	// Feedback modal has highest priority.
	if i.feedbackModal != nil {
		width, height := i.getTerminalDimensions()
		i.modalRegistry.Register(intents.NewErrorModalAdapter(i.feedbackModal, width, height, i.Theme()))
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

	// Suggestion events modal (drill-down from skill suggestions).
	if i.suggestionEventsModal != nil && i.suggestionEventsModal.IsVisible() {
		i.modalRegistry.Register(intents.NewViewModalAdapter(
			i.suggestionEventsModal.IsVisible,
			i.suggestionEventsModal.View,
			i.suggestionEventsModal.Update,
		))
	}

	// Skill suggestion review modal.
	if i.skillSuggestionModal != nil && i.skillSuggestionModal.IsVisible() {
		i.modalRegistry.Register(intents.NewViewModalAdapter(
			i.skillSuggestionModal.IsVisible,
			i.skillSuggestionModal.View,
			i.skillSuggestionModal.Update,
		))
	}
}

// modalWithDimensions is an interface for modals that support dimensions.
type modalWithDimensions interface {
	SetDimensions(width, height int)
	Show()
}

// showModalWithDimensions sets dimensions and shows a modal.
func (i *Intent) showModalWithDimensions(modal modalWithDimensions) {
	width, height := i.getTerminalDimensions()
	modal.SetDimensions(width, height)
	modal.Show()
}

// openSuggestionEventsModal resolves event IDs from the currently selected skill suggestion
// and opens an events modal to display them.
func (i *Intent) openSuggestionEventsModal() tea.Cmd {
	if i.skillSuggestionModal == nil {
		return noopCmd
	}

	selected := i.skillSuggestionModal.GetCurrentSkill()
	if selected == nil {
		return noopCmd
	}

	events := i.resolveEventsByIDs(selected.EventIDs)

	width, height := i.getTerminalDimensions()
	i.suggestionEventsModal = skillmodals.NewEventsModal(
		"suggestion",
		selected.Name,
		events,
		i.Theme(),
	)
	i.suggestionEventsModal.SetDimensions(width, height)
	i.suggestionEventsModal.Show()

	return noopCmd
}

// resolveEventsByIDs resolves event IDs to Event objects using the burst service.
func (i *Intent) resolveEventsByIDs(eventIDs []string) []*career.Event {
	if i.context.Service == nil || len(eventIDs) == 0 {
		return []*career.Event{}
	}

	ctx := i.getContext()
	events := make([]*career.Event, 0, len(eventIDs))
	for _, id := range eventIDs {
		event, err := i.context.Service.GetEventByID(ctx, id)
		if err == nil && event != nil {
			events = append(events, event)
		}
	}
	return events
}

// startBurstDetection loads all events and triggers AI burst detection.
func (i *Intent) startBurstDetection() tea.Cmd {
	if i.context.Service == nil {
		i.ShowErrorModal("Detection Failed", "Service not available")
		i.state = StateList
		return nil
	}

	i.cancelPreviousOperation()
	ctx, cancel := context.WithCancel(context.Background())
	i.cancelFunc = cancel

	i.suggestionsLoading = true
	i.loadingModal = feedback.NewLoadingModal("Detecting burst patterns...", true).WithTheme(i.Theme())

	asyncCmd := i.createBurstDetectionCmd(ctx, i.context.Service, i.context.Bursts)
	return tea.Batch(asyncCmd, i.loadingModal.Init())
}

// startSkillInference triggers skill inference for the selected burst.
func (i *Intent) startSkillInference() tea.Cmd {
	if i.selectedBurst == nil {
		i.ShowErrorModal("Inference Failed", "No burst selected")
		i.state = StateList
		return nil
	}

	if i.context.SkillInferenceService == nil {
		i.ShowErrorModal("Inference Failed", "Skill inference service not available")
		i.state = StateList
		return nil
	}

	i.cancelPreviousOperation()

	i.inferringSkills = true
	i.state = StateInferringSkills
	i.loadingModal = feedback.NewLoadingModal("Detecting skills from burst events...", true).WithTheme(i.Theme())

	asyncCmd := i.inferSkillsFromBurst(i.selectedBurst)
	return tea.Batch(asyncCmd, i.loadingModal.Init())
}

// cancelPreviousOperation cancels any in-progress async operation.
func (i *Intent) cancelPreviousOperation() {
	if i.cancelFunc != nil {
		i.cancelFunc()
	}
}

// createBurstDetectionCmd creates the async command for burst detection.
func (i *Intent) createBurstDetectionCmd(
	ctx context.Context,
	service BurstService,
	existingBursts []*career.Burst,
) func() tea.Msg {
	return func() tea.Msg {
		if ctx.Err() != nil {
			return BurstSuggestionsLoadedMsg{Error: ctx.Err()}
		}

		eventIDs, err := i.getUnassignedEventIDs(ctx, service, existingBursts)
		if err != nil {
			return BurstSuggestionsLoadedMsg{Error: err}
		}

		if len(eventIDs) == 0 {
			return BurstSuggestionsLoadedMsg{
				Error: fmt.Errorf("no unassigned events available for burst detection"),
			}
		}

		suggestions, err := service.SuggestBursts(ctx, eventIDs)
		return BurstSuggestionsLoadedMsg{Suggestions: suggestions, Error: err}
	}
}

// getUnassignedEventIDs returns event IDs not already assigned to bursts.
func (i *Intent) getUnassignedEventIDs(
	ctx context.Context,
	service BurstService,
	existingBursts []*career.Burst,
) ([]string, error) {
	events, err := service.ListEvents(ctx, careerrepo.EventListFilters{Limit: -1})
	if err != nil {
		return nil, fmt.Errorf("failed to load events: %w", err)
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	usedEventIDs := i.buildUsedEventIDSet(existingBursts)

	var eventIDs []string
	for _, event := range events {
		if !usedEventIDs[event.ID] {
			eventIDs = append(eventIDs, event.ID)
		}
	}
	return eventIDs, nil
}

// buildUsedEventIDSet builds a set of event IDs already in bursts.
func (i *Intent) buildUsedEventIDSet(bursts []*career.Burst) map[string]bool {
	usedEventIDs := make(map[string]bool)
	for _, burst := range bursts {
		for _, eventID := range burst.EventIDs {
			usedEventIDs[eventID] = true
		}
	}
	return usedEventIDs
}

// isInSuggestionState returns true if currently in suggestion review state.
func (i *Intent) isInSuggestionState() bool {
	return i.state == StateSuggestionReview || i.state == StateSuggesting
}

// createBurstsFromSuggestions creates bursts from accepted suggestions.
func (i *Intent) createBurstsFromSuggestions(suggestions []burstfact.BurstSuggestion) []*career.Burst {
	createdBursts := make([]*career.Burst, 0, len(suggestions))
	for _, suggestion := range suggestions {
		burst := i.createBurstFromSuggestion(suggestion)
		if burst != nil {
			createdBursts = append(createdBursts, burst)
		}
	}
	return createdBursts
}

// createBurstFromSuggestion creates and saves a single burst from a suggestion.
func (i *Intent) createBurstFromSuggestion(suggestion burstfact.BurstSuggestion) *career.Burst {
	burst := &career.Burst{
		Name:        suggestion.Name,
		Description: suggestion.Description,
		EventIDs:    suggestion.EventIDs,
		Confirmed:   true,
	}

	if i.context.BurstRepository != nil {
		if err := i.context.BurstRepository.Create(i.getContext(), burst); err != nil {
			i.ShowErrorModal("Create Failed", fmt.Sprintf("Failed to create burst: %v", err))
			return nil
		}
	}

	i.filteredBursts = append(i.filteredBursts, burst)
	i.context.Bursts = append(i.context.Bursts, burst)
	return burst
}

// startFactExtractionForBursts starts fact extraction for multiple bursts.
func (i *Intent) startFactExtractionForBursts(bursts []*career.Burst) tea.Cmd {
	if len(bursts) == 0 {
		i.state = StateList
		return nil
	}

	var cmds []tea.Cmd
	for _, burst := range bursts {
		cmds = append(cmds, i.extractFactsForBurst(burst))
	}

	i.extractingFacts = true
	i.state = StateExtractingFacts
	i.loadingModal = feedback.NewLoadingModal("Extracting facts from accepted bursts...", true).WithTheme(i.Theme())
	cmds = append(cmds, i.loadingModal.Init())
	return tea.Batch(cmds...)
}

// saveAndExtractBurst saves a burst immediately and triggers fact extraction.
// This is called when user presses 'a' to provide instant feedback and ensure
// fact extraction happens even if the app crashes before modal closes.
func (i *Intent) saveAndExtractBurst(suggestion burstfact.BurstSuggestion) tea.Cmd {
	_, cmd := i.saveAndExtractBurstWithResult(suggestion)
	return cmd
}

// saveAndExtractBurstWithResult saves a burst and returns both the burst and the extraction command.
// This is used when we need to display the burst detail modal after saving.
func (i *Intent) saveAndExtractBurstWithResult(suggestion burstfact.BurstSuggestion) (*career.Burst, tea.Cmd) {
	burst := &career.Burst{
		ID:          "",
		Name:        suggestion.Name,
		Description: suggestion.Description,
		EventIDs:    suggestion.EventIDs,
		Confirmed:   true,
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
			return FactExtractionCompleteMsg{Burst: burst, Error: ctx.Err()}
		}

		if service == nil {
			return FactExtractionCompleteMsg{Burst: burst, Error: fmt.Errorf("service not available")}
		}

		facts, err := service.ExtractFactsFromBurst(ctx, burst)
		if err != nil {
			return FactExtractionCompleteMsg{Burst: burst, Error: err}
		}

		// Check if cancelled after extraction.
		if ctx.Err() != nil {
			return FactExtractionCompleteMsg{Burst: burst, Error: ctx.Err()}
		}

		savedFacts := make([]*career.Fact, 0, len(facts))
		for idx := range facts {
			fact := &facts[idx]
			fact.SourceBurstID = burst.ID

			// Check if cancelled during save loop.
			if ctx.Err() != nil {
				return FactExtractionCompleteMsg{Burst: burst, Error: ctx.Err()}
			}

			if err := service.SaveFact(ctx, fact); err != nil {
				continue
			}
			savedFacts = append(savedFacts, fact)
		}

		return FactExtractionCompleteMsg{Facts: savedFacts, Burst: burst}
	}
}

// inferSkillsFromBurst runs skill inference on burst events and returns suggestions.
func (i *Intent) inferSkillsFromBurst(burst *career.Burst) tea.Cmd {
	if burst == nil {
		return func() tea.Msg {
			return SkillSuggestionsErrorMsg{Err: fmt.Errorf("no burst provided")}
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
	service := i.context.SkillInferenceService

	return func() tea.Msg {
		// Check if cancelled before starting.
		if ctx.Err() != nil {
			return SkillSuggestionsErrorMsg{Err: ctx.Err()}
		}

		if service == nil {
			return SkillSuggestionsErrorMsg{Err: fmt.Errorf("skill inference service not available")}
		}

		// Load events for this burst.
		events := i.loadBurstEvents(burst)
		if len(events) == 0 {
			return SkillSuggestionsErrorMsg{Err: fmt.Errorf("no events found for burst")}
		}

		// Check if cancelled after loading events.
		if ctx.Err() != nil {
			return SkillSuggestionsErrorMsg{Err: ctx.Err()}
		}

		result, err := service.InferSkillsFromEvents(ctx, events)
		if err != nil {
			return SkillSuggestionsErrorMsg{Err: fmt.Errorf("skill detection failed: %w", err)}
		}

		return SkillSuggestionsLoadedMsg{
			Suggestions:        result.Suggestions,
			ExistingSkillNames: result.ExistingSkillNames,
		}
	}
}

// saveSkillFromSuggestion persists a single accepted skill suggestion synchronously.
// This mirrors saveAndExtractBurst: save immediately on each 'a' press.
func (i *Intent) saveSkillFromSuggestion(suggestion skillinference.SkillSuggestion) {
	service := i.context.SkillInferenceService
	if service == nil {
		i.ShowErrorModal("Skill Creation Failed", "skill inference service not available")
		return
	}

	ctx := i.getContext()
	_, err := service.CreateSkillsFromSuggestions(ctx, []skillinference.SkillSuggestion{suggestion})
	if err != nil {
		i.ShowErrorModal("Skill Creation Failed", fmt.Sprintf("Failed to create skill: %v", err))
	}
}
