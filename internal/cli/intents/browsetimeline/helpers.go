package browsetimeline

import (
	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/intents"
	burstModals "github.com/baphled/kariya/internal/cli/screens/burst_management/modals"
	skillModals "github.com/baphled/kariya/internal/cli/screens/skills/modals"
	"github.com/baphled/kariya/internal/cli/screens/timeline"
	"github.com/baphled/kariya/internal/cli/screens/timeline/modals"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/skillinference"
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

// RefreshData re-applies current filters and rebuilds the timeline screen
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (i *Intent) RefreshData() tea.Cmd {
	i.applyFilters()
	i.transitionToScreen(timeline.NewTimelineEventListScreen(i.filteredEvents))
	return nil
}

// HasVisibleSkillsModal checks whether the skills detail modal is currently
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (i *Intent) HasVisibleSkillsModal() bool {
	return i.viewSkillsModal != nil && i.viewSkillsModal.IsVisible()
}

// HasVisibleErrorModal checks whether an error modal is currently displayed
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (i *Intent) HasVisibleErrorModal() bool {
	return i.errorModal != nil
}

// HasVisibleQuickAddModal checks whether the quick-add event modal is currently
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (i *Intent) HasVisibleQuickAddModal() bool {
	return i.quickAddModal != nil && i.quickAddModal.IsVisible()
}

// HasVisibleEditModal checks whether the event edit modal is currently
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (i *Intent) HasVisibleEditModal() bool {
	return i.editModal != nil && i.editModal.IsVisible()
}

// ShowErrorModal presents an error overlay to the user, taking highest
//
// Expected:
//   - Must be a valid string.
//
// Side effects:
//   - None.
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

	// View modals (priority order: picker > suggestion > skills > detail).
	if i.skillPickerModal != nil {
		i.modalRegistry.Register(intents.NewViewModalAdapter(
			i.skillPickerModal.IsVisible,
			i.skillPickerModal.View,
			i.skillPickerModal.Update,
		))
	}

	if i.skillSuggestionModal != nil {
		i.modalRegistry.Register(intents.NewViewModalAdapter(
			i.skillSuggestionModal.IsVisible,
			i.skillSuggestionModal.View,
			i.skillSuggestionModal.Update,
		))
	}

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

func (i *Intent) openSkillPickerModal() tea.Cmd {
	if i.selectedEvent == nil {
		return nil
	}
	ctx := i.getContext()
	allSkills, err := i.context.CLIEventService.ListAllSkills(ctx)
	if err != nil {
		i.ShowErrorModal("Error Loading Skills", err.Error())
		return nil
	}
	eventSkills, err := i.context.CLIEventService.GetSkillsForEvent(ctx, i.selectedEvent.ID)
	if err != nil {
		i.ShowErrorModal("Error Loading Skills", err.Error())
		return nil
	}
	eventSkillIDs := make(map[string]bool)
	for _, s := range eventSkills {
		eventSkillIDs[s.ID] = true
	}
	availableSkills := make([]*career.Skill, 0)
	for _, s := range allSkills {
		if !eventSkillIDs[s.ID] {
			availableSkills = append(availableSkills, s)
		}
	}
	width, height := i.getTerminalDimensions()
	i.skillPickerModal = modals.NewSkillPickerModal(availableSkills, i.Theme())
	i.skillPickerModal.SetDimensions(width, height)
	i.skillPickerModal.Show()
	return nil
}

func (i *Intent) linkSkillToCurrentEvent(skill *career.Skill) tea.Cmd {
	return i.performSkillLinkOperation(skill, true)
}

func (i *Intent) unlinkSkillFromCurrentEvent(skill *career.Skill) tea.Cmd {
	return i.performSkillLinkOperation(skill, false)
}

func (i *Intent) performSkillLinkOperation(skill *career.Skill, link bool) tea.Cmd {
	if i.selectedEvent == nil || skill == nil {
		return nil
	}
	ctx := i.getContext()
	var err error
	if link {
		err = i.context.CLIEventService.LinkSkillToEvent(ctx, i.selectedEvent.ID, skill.ID)
	} else {
		err = i.context.CLIEventService.UnlinkSkillFromEvent(ctx, i.selectedEvent.ID, skill.ID)
	}
	if err != nil {
		action := "Linking"
		if !link {
			action = "Unlinking"
		}
		i.ShowErrorModal("Error "+action+" Skill", err.Error())
		return nil
	}
	return func() tea.Msg {
		if link {
			return SkillLinkedMsg{EventID: i.selectedEvent.ID, SkillID: skill.ID}
		}
		return SkillUnlinkedMsg{EventID: i.selectedEvent.ID, SkillID: skill.ID}
	}
}

func (i *Intent) refreshSkillsModal() tea.Cmd {
	if i.selectedEvent == nil || i.viewSkillsModal == nil {
		return nil
	}
	ctx := i.getContext()
	skills, err := i.context.CLIEventService.GetSkillsForEvent(ctx, i.selectedEvent.ID)
	if err != nil {
		i.ShowErrorModal("Error Loading Skills", err.Error())
		return nil
	}
	i.viewSkillsModal.SetSkills(skills)
	return nil
}

func (i *Intent) openSkillAddModal() tea.Cmd {
	width, height := i.getTerminalDimensions()
	i.skillAddModal = skillModals.NewAddEditModal(nil, width, height)
	return i.skillAddModal.Init()
}

func (i *Intent) createAndLinkSkill(skillData *skillModals.SkillEditData) tea.Cmd {
	if i.selectedEvent == nil || skillData == nil {
		return nil
	}
	ctx := i.getContext()
	newSkill := skillData.ToSkill("")
	if err := i.context.CLISkillService.Create(ctx, newSkill); err != nil {
		i.ShowErrorModal("Error Creating Skill", err.Error())
		return nil
	}
	if err := i.context.CLIEventService.LinkSkillToEvent(ctx, i.selectedEvent.ID, newSkill.ID); err != nil {
		i.ShowErrorModal("Error Linking Skill", err.Error())
		return nil
	}
	return func() tea.Msg { return SkillCreatedMsg{Skill: newSkill} }
}

func (i *Intent) inferSkillsFromEvent() tea.Cmd {
	if i.selectedEvent == nil || i.context.SkillInferenceService == nil {
		return nil
	}
	ctx, service, event := i.getContext(), i.context.SkillInferenceService, i.selectedEvent
	return func() tea.Msg {
		result, err := service.InferSkillsFromEvents(ctx, []*career.Event{event})
		if err != nil {
			return SkillSuggestionsErrorMsg{Err: err}
		}
		return SkillSuggestionsLoadedMsg{Suggestions: result.Suggestions, ExistingSkillNames: result.ExistingSkillNames}
	}
}

func (i *Intent) handleSkillSuggestionsLoaded(msg SkillSuggestionsLoadedMsg) tea.Cmd {
	if len(msg.Suggestions) == 0 {
		i.ShowErrorModal("No Skills Detected", "No skills detected from event.")
		return nil
	}
	newSuggestions := filterNewSkillSuggestions(msg.Suggestions, msg.ExistingSkillNames)
	if len(newSuggestions) == 0 {
		i.ShowErrorModal("All Skills Tracked", "All detected skills already in profile.")
		return nil
	}
	width, height := i.getTerminalDimensions()
	i.skillSuggestionModal = burstModals.NewSkillSuggestionModal(newSuggestions, i.Theme())
	i.skillSuggestionModal.SetDimensions(width, height)
	i.skillSuggestionModal.Show()
	return nil
}

func filterNewSkillSuggestions(suggestions []skillinference.SkillSuggestion, existingNames []string) []skillinference.SkillSuggestion {
	existingMap := make(map[string]bool)
	for _, name := range existingNames {
		existingMap[name] = true
	}
	result := make([]skillinference.SkillSuggestion, 0)
	for _, s := range suggestions {
		if !existingMap[s.Name] {
			result = append(result, s)
		}
	}
	return result
}

func (i *Intent) saveSkillFromSuggestion(suggestion skillinference.SkillSuggestion) {
	if i.context.SkillInferenceService == nil {
		return
	}
	ctx := i.getContext()
	skills, err := i.context.SkillInferenceService.CreateSkillsFromSuggestions(ctx, []skillinference.SkillSuggestion{suggestion})
	if err != nil {
		i.ShowErrorModal("Skill Creation Failed", err.Error())
		return
	}
	if i.selectedEvent != nil && len(skills) > 0 {
		_ = i.context.CLIEventService.LinkSkillToEvent(ctx, i.selectedEvent.ID, skills[0].ID)
	}
}
