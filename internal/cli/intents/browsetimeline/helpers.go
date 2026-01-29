package browsetimeline

import (
	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/screens/timeline"
	"github.com/baphled/kariya/internal/cli/screens/timeline/modals"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/domain/career"
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

// showSkillsForCurrentEvent loads and displays skills for the currently selected event.
func (i *Intent) showSkillsForCurrentEvent() tea.Cmd {
	if i.selectedEvent == nil {
		return nil
	}

	var skills []*career.Skill
	if i.context.CLIEventService != nil {
		ctx := i.getContext()
		var err error
		skills, err = i.context.CLIEventService.GetSkillsForEvent(ctx, i.selectedEvent.ID)
		if err != nil {
			skills = []*career.Skill{}
		}
	}

	width, height := i.getTerminalDimensions()
	theme := i.Theme()
	i.viewSkillsModal = modals.NewSkillsDetailModal(i.selectedEvent.ID, skills, theme)
	i.viewSkillsModal.SetDimensions(width, height)
	i.viewSkillsModal.Show()

	return nil
}

// openSearchModal creates and initializes the search modal.
func (i *Intent) openSearchModal() tea.Cmd {
	width, height := i.getTerminalDimensions()
	i.searchModal = modals.NewSearchModal(i.filters.SearchText, width, height)
	return i.searchModal.Init()
}

// openSortModal creates and initializes the sort modal.
func (i *Intent) openSortModal() tea.Cmd {
	width, height := i.getTerminalDimensions()
	current := &modals.SortConfig{
		SortBy:    i.filters.SortBy,
		SortOrder: i.filters.SortOrder,
	}
	i.sortModal = modals.NewSortModal(i.context.Events, current, width, height)
	return i.sortModal.Init()
}

// openFilterModal creates and initializes the filter modal.
func (i *Intent) openFilterModal() tea.Cmd {
	width, height := i.getTerminalDimensions()
	currentFilters := &modals.TimelineFilters{
		SearchText: i.filters.SearchText,
		Tags:       i.filters.Tags,
		Companies:  i.filters.Companies,
		Categories: i.filters.Categories,
		Projects:   i.filters.Projects,
		SortBy:     i.filters.SortBy,
		SortOrder:  i.filters.SortOrder,
	}
	i.filterModal = modals.NewFilterModal(i.context.Events, currentFilters, width, height)
	return i.filterModal.Init()
}

// getStateName returns a human-readable name for the current state.
func (i *Intent) getStateName() string {
	switch i.state {
	case StateTimeline:
		return "Timeline"
	case StateDeleteConfirm:
		return "Delete Confirmation"
	default:
		return "Unknown"
	}
}

// getContextHelp returns themed keyboard shortcuts for the current state.
func (i *Intent) getContextHelp() string {
	theme := i.Theme()

	switch i.state {
	case StateTimeline:
		badges := []*primitives.Badge{
			primitives.NavigateBadge(theme),
			primitives.HelpKeyBadge("Enter", "View Details", theme),
			primitives.AddBadge(theme),
			primitives.EditBadge(theme),
			primitives.DeleteBadge(theme),
			primitives.SearchBadge(theme),
			primitives.FilterBadge(theme),
			primitives.HelpKeyBadge("s", "Sort", theme),
		}

		if i.HasActiveFilters() {
			badges = append(badges, primitives.HelpKeyBadge("x", "Clear filters", theme))
		}

		badges = append(badges, primitives.BackBadge(theme))

		return intents.CombineThemedFooters(
			intents.ThemedCustomFooter(theme, badges...),
			intents.ThemedGlobalBadges(theme),
		)

	case StateDeleteConfirm:
		return ""

	default:
		return intents.ThemedGlobalBadges(theme)
	}
}

// removeEventFromList removes an event from both the full list and filtered list.
func (i *Intent) removeEventFromList(eventID string) {
	for idx, evt := range i.context.Events {
		if evt.ID == eventID {
			i.context.Events = append(i.context.Events[:idx], i.context.Events[idx+1:]...)
			break
		}
	}

	for idx, evt := range i.filteredEvents {
		if evt.ID == eventID {
			i.filteredEvents = append(i.filteredEvents[:idx], i.filteredEvents[idx+1:]...)
			break
		}
	}

	i.transitionToScreen(timeline.NewTimelineEventListScreen(i.filteredEvents))
}

// RefreshData reloads/refreshes the filtered data.
func (i *Intent) RefreshData() tea.Cmd {
	i.applyFilters()
	i.transitionToScreen(timeline.NewTimelineEventListScreen(i.filteredEvents))
	return nil
}

// HasVisibleSkillsModal returns true if the skills modal is currently visible.
func (i *Intent) HasVisibleSkillsModal() bool {
	return i.viewSkillsModal != nil && i.viewSkillsModal.IsVisible()
}

// HasVisibleErrorModal returns true if the error modal is currently visible.
func (i *Intent) HasVisibleErrorModal() bool {
	return i.errorModal != nil
}

// HasVisibleQuickAddModal returns true if the quick add modal is currently visible.
func (i *Intent) HasVisibleQuickAddModal() bool {
	return i.quickAddModal != nil && i.quickAddModal.IsVisible()
}

// HasVisibleEditModal returns true if the edit modal is currently visible.
func (i *Intent) HasVisibleEditModal() bool {
	return i.editModal != nil && i.editModal.IsVisible()
}

// ShowErrorModal displays an error modal with the given title and message.
func (i *Intent) ShowErrorModal(title, message string) {
	i.errorModal = feedback.NewErrorModal(title, message)
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

	// Form modals (search, filter, sort, quickAdd, edit).
	if i.searchModal != nil {
		i.modalRegistry.Register(intents.NewFormModalAdapter(
			i.searchModal.IsVisible,
			i.searchModal.View,
			i.searchModal.Update,
		))
	}

	if i.filterModal != nil {
		i.modalRegistry.Register(intents.NewFormModalAdapter(
			i.filterModal.IsVisible,
			i.filterModal.View,
			i.filterModal.Update,
		))
	}

	if i.sortModal != nil {
		i.modalRegistry.Register(intents.NewFormModalAdapter(
			i.sortModal.IsVisible,
			i.sortModal.View,
			i.sortModal.Update,
		))
	}

	if i.quickAddModal != nil {
		i.modalRegistry.Register(intents.NewFormModalAdapter(
			i.quickAddModal.IsVisible,
			i.quickAddModal.View,
			i.quickAddModal.Update,
		))
	}

	if i.editModal != nil {
		i.modalRegistry.Register(intents.NewFormModalAdapter(
			i.editModal.IsVisible,
			i.editModal.View,
			i.editModal.Update,
		))
	}

	// Confirm modal (delete).
	if i.deleteModal != nil {
		i.modalRegistry.Register(intents.NewConfirmModalAdapter(i.deleteModal))
	}

	// View modals.
	if i.viewSkillsModal != nil {
		i.modalRegistry.Register(intents.NewViewModalAdapter(
			i.viewSkillsModal.IsVisible,
			i.viewSkillsModal.View,
			i.viewSkillsModal.Update,
		))
	}

	if i.viewDetailModal != nil {
		i.modalRegistry.Register(intents.NewViewModalAdapter(
			i.viewDetailModal.IsVisible,
			i.viewDetailModal.View,
			i.viewDetailModal.Update,
		))
	}
}
