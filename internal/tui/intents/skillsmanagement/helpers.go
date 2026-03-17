package skillsmanagement

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	domain "github.com/baphled/kariya/internal/domain/career"
	career "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/service/career/skillinference"

	"github.com/baphled/kariya/internal/tui/forms"
	"github.com/baphled/kariya/internal/tui/intents"
	eventviews "github.com/baphled/kariya/internal/tui/views/event"
	"github.com/baphled/kariya/internal/tui/views/shared"
	skillviews "github.com/baphled/kariya/internal/tui/views/skill"
	"github.com/baphled/kariya/internal/ui/behaviors"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/uikit/feedback"
	tea "github.com/charmbracelet/bubbletea"
)

// skillRowFormatterWithCounts creates a row formatter that includes event counts.
func skillRowFormatterWithCounts(eventCounts map[string]int) behaviors.RowFormatter[*domain.Skill] {
	return func(skill *domain.Skill, _ int) []string {
		name := skill.Name
		if len(name) > 22 {
			name = name[:22] + "..."
		}

		category := skill.Category
		if category == "" {
			category = "-"
		}

		level := skill.Level
		if level == "" {
			level = "-"
		}

		years := "-"
		if skill.YearsUsed != nil {
			years = strconv.Itoa(*skill.YearsUsed)
		}

		eventCount := "-"
		if eventCounts != nil {
			if count, ok := eventCounts[skill.ID]; ok {
				eventCount = strconv.Itoa(count)
			}
		}

		return []string{name, category, level, years, eventCount}
	}
}

// syncTableSelection syncs the TableBehavior selection with the intent's data.
func (i *Intent) syncTableSelection() {
	i.selectedIndex = i.tableBehavior.GetSelectedIndex()
	if selected := i.tableBehavior.GetSelectedItem(); selected != nil {
		i.selectedSkill = *selected
	} else {
		i.selectedSkill = nil
	}
}

func (i *Intent) findSkillByID(skillID string) *domain.Skill {
	if skillID == "" {
		return nil
	}

	for _, skill := range i.skills {
		if skill != nil && skill.ID == skillID {
			return skill
		}
	}

	return nil
}

func (i *Intent) findEventByID(eventID string) *domain.Event {
	if eventID == "" {
		return nil
	}

	for _, event := range i.skillEvents {
		if event != nil && event.ID == eventID {
			return event
		}
	}

	return nil
}

func (i *Intent) findSkillSuggestionByName(name string) *skillinference.SkillSuggestion {
	if name == "" {
		return nil
	}

	for idx := range i.skillSuggestions {
		suggestion := i.skillSuggestions[idx]
		if suggestion.Name == name {
			return &i.skillSuggestions[idx]
		}
	}

	return nil
}

// getBreadcrumbs returns breadcrumbs for the current state.
func (i *Intent) getBreadcrumbs() []string {
	breadcrumbs := []string{"Skills"}

	switch i.state {
	case StateDetail, StateDetailEvents, StateDetailEventDetail:
		if i.selectedSkill != nil {
			breadcrumbs = append(breadcrumbs, i.selectedSkill.Name)
		}
		if i.state == StateDetailEvents {
			breadcrumbs = append(breadcrumbs, "Events")
		}
		if i.state == StateDetailEventDetail {
			breadcrumbs = append(breadcrumbs, "Events", "Detail")
		}
	case StateAdd:
		breadcrumbs = append(breadcrumbs, "Add")
	case StateEdit:
		breadcrumbs = append(breadcrumbs, "Edit")
	case StateDelete:
		breadcrumbs = append(breadcrumbs, "Delete")
	}

	return breadcrumbs
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
	if i.feedbackModal != nil {
		width, height := i.getTerminalDimensions()
		i.modalRegistry.Register(intents.NewErrorModalAdapter(i.feedbackModal, width, height, i.Theme()))
	}

	// Loading modal (for StateInferringSkills).
	// When loading is active, it should be the only modal visible.
	if i.loadingModal != nil {
		width, height := i.getTerminalDimensions()
		i.modalRegistry.Register(intents.NewErrorModalAdapter(i.loadingModal, width, height, i.Theme()))
		// Don't register other modals when loading - loading takes full precedence.
		return
	}

	// Form view adapters (search, filter, sort) — implement ManagedModal directly.
	if i.searchAdapter != nil {
		i.modalRegistry.Register(i.searchAdapter)
	}

	if i.filterAdapter != nil {
		i.modalRegistry.Register(i.filterAdapter)
	}

	if i.sortAdapter != nil {
		i.modalRegistry.Register(i.sortAdapter)
	}

	if i.addEditModal != nil {
		i.modalRegistry.Register(intents.NewFormModalAdapter(
			i.addEditModal.IsVisible,
			i.addEditModal.View,
			i.addEditModal.Update,
		))
	}

	// Confirm modal (delete).
	if i.deleteModal != nil {
		i.modalRegistry.Register(intents.NewConfirmModalAdapter(i.deleteModal))
	}

	// View modals.
	if i.viewDetailModal != nil {
		i.modalRegistry.Register(intents.NewViewModalAdapter(
			i.viewDetailModal.IsVisible,
			i.viewDetailModal.View,
			i.viewDetailModal.Update,
		))
	}

	if i.skillEventsModal != nil {
		i.modalRegistry.Register(intents.NewViewModalAdapter(
			i.skillEventsModal.IsVisible,
			i.skillEventsModal.View,
			i.skillEventsModal.Update,
		))
	}

	if i.eventDetailModal != nil {
		i.modalRegistry.Register(intents.NewViewModalAdapter(
			i.eventDetailModal.IsVisible,
			i.eventDetailModal.View,
			i.eventDetailModal.Update,
		))
	}

	if i.suggestionEventsModal != nil && i.suggestionEventsModal.IsVisible() {
		i.modalRegistry.Register(intents.NewViewModalAdapter(
			i.suggestionEventsModal.IsVisible,
			i.suggestionEventsModal.View,
			i.suggestionEventsModal.Update,
		))
	}

	if i.skillSuggestionModal != nil && i.skillSuggestionModal.IsVisible() {
		i.modalRegistry.Register(intents.NewViewModalAdapter(
			i.skillSuggestionModal.IsVisible,
			i.skillSuggestionModal.View,
			i.skillSuggestionModal.Update,
		))
	}
}

// getTerminalDimensions returns current terminal dimensions with fallback defaults.
func (i *Intent) getTerminalDimensions() (width, height int) {
	width, height = behaviors.DefaultModalDimensions()
	if termInfo := i.GetTerminalInfo(); termInfo != nil && termInfo.Width > 0 && termInfo.Height > 0 {
		width, height = termInfo.Width, termInfo.Height
	}
	return
}

// extractSkillFilterFields extracts unique categories and levels from skills
// as FilterFieldConfig slices for the filter view.
func extractSkillFilterFields(skills []*domain.Skill) []skillviews.FilterFieldConfig {
	categoryMap := make(map[string]bool)
	levelMap := make(map[string]bool)

	for _, s := range skills {
		if s.Category != "" {
			categoryMap[s.Category] = true
		}
		if s.Level != "" {
			levelMap[s.Level] = true
		}
	}

	var fields []skillviews.FilterFieldConfig

	if len(categoryMap) > 0 {
		opts := make([]forms.SelectOption, 0, len(categoryMap))
		for cat := range categoryMap {
			opts = append(opts, forms.SelectOption{Key: cat, Value: cat})
		}
		fields = append(fields, skillviews.FilterFieldConfig{Key: "categories", Title: "Filter by Category", Options: opts})
	}

	if len(levelMap) > 0 {
		opts := make([]forms.SelectOption, 0, len(levelMap))
		for lvl := range levelMap {
			opts = append(opts, forms.SelectOption{Key: lvl, Value: lvl})
		}
		fields = append(fields, skillviews.FilterFieldConfig{Key: "levels", Title: "Filter by Level", Options: opts})
	}

	return fields
}

// openFilterModal opens the filter modal with current filters pre-populated.
func (i *Intent) openFilterModal() tea.Cmd {
	width, height := i.getTerminalDimensions()
	fields := extractSkillFilterFields(i.skills)

	var defaults *skillviews.FilterFormData
	if i.context.Filters != nil {
		defaults = &skillviews.FilterFormData{
			Selections: map[string][]string{},
		}
		if i.context.Filters.Category != "" {
			defaults.Selections["categories"] = []string{i.context.Filters.Category}
		}
		if i.context.Filters.Level != "" {
			defaults.Selections["levels"] = []string{i.context.Filters.Level}
		}
		if i.context.Filters.MinEvents > 0 {
			defaults.MinYearsStr = strconv.Itoa(i.context.Filters.MinEvents)
		}
	}

	filterView := skillviews.NewFilter(fields, defaults)
	filterView.SetTerminalInfo(width, height)
	i.filterAdapter = intents.NewFormViewAdapter(filterView)
	i.filterAdapter.Show()
	return filterView.Init()
}

// openSortModal opens the sort modal with current sort config pre-populated.
func (i *Intent) openSortModal() tea.Cmd {
	width, height := i.getTerminalDimensions()
	sortFields := []shared.SortFieldOption{
		{Key: "name", Label: "Name"},
		{Key: "category", Label: "Category"},
		{Key: "level", Label: "Level"},
		{Key: "years", Label: "Years of Experience"},
		{Key: "events", Label: "Events Count"},
	}

	defaultSortBy := ""
	defaultSortOrder := ""
	if i.context.Filters != nil {
		defaultSortBy = i.context.Filters.SortBy
		defaultSortOrder = i.context.Filters.SortOrder
	}

	sortView := shared.NewSortView(sortFields, defaultSortBy, defaultSortOrder)
	sortView.SetTerminalInfo(width, height)
	i.sortAdapter = intents.NewFormViewAdapter(sortView)
	i.sortAdapter.Show()
	return sortView.Init()
}

// openSearchModal opens the search modal with current search text pre-populated.
func (i *Intent) openSearchModal() tea.Cmd {
	width, height := i.getTerminalDimensions()
	searchView := shared.NewSearchView("Search skills...")
	searchView.SetTerminalInfo(width, height)
	i.searchAdapter = intents.NewFormViewAdapter(searchView)
	i.searchAdapter.Show()
	return searchView.Init()
}

// openViewDetailModal opens the view detail modal for the selected skill.
func (i *Intent) openViewDetailModal() tea.Cmd {
	if i.selectedSkill == nil {
		return nil
	}

	skill := i.selectedSkill

	// Get event count for this skill.
	eventCount := 0
	if i.eventCounts != nil {
		eventCount = i.eventCounts[skill.ID]
	}

	width, height := i.getTerminalDimensions()

	i.viewDetailModal = skillviews.NewDetail(display.SkillFromDomain(skill), i.Theme(), eventCount, nil)
	i.viewDetailModal.SetDimensions(width, height)
	i.viewDetailModal.Show()

	return nil
}

// openAddEditModal opens the add/edit modal for a skill.
func (i *Intent) openAddEditModal(skill *domain.Skill) tea.Cmd {
	width, height := i.getTerminalDimensions()
	formData := forms.GetSkillFormData(skill)
	originalSkillID := ""
	if skill != nil {
		originalSkillID = skill.ID
	}

	i.addEditModal = skillviews.NewAddEdit(formData, width, height, originalSkillID)
	return i.addEditModal.Init()
}

// openDeleteModal opens the delete confirmation modal for a skill.
func (i *Intent) openDeleteModal(skill *domain.Skill) tea.Cmd {
	if skill == nil {
		return nil
	}

	i.selectedSkill = skill
	skillName := skill.Name
	if len(skillName) > 50 {
		skillName = skillName[:47] + "..."
	}

	i.deleteModal = feedback.NewConfirmModal(
		"Delete Skill",
		fmt.Sprintf("Are you sure you want to delete '%s'?", skillName),
	).WithVariant(feedback.ConfirmDestructive)
	return i.deleteModal.Init()
}

// openSkillEventsModal opens the skill events modal for the selected skill.
func (i *Intent) openSkillEventsModal(events []*domain.Event) tea.Cmd {
	if i.selectedSkill == nil {
		return nil
	}

	i.skillEvents = append([]*domain.Event(nil), events...)

	width, height := i.getTerminalDimensions()

	i.skillEventsModal = skillviews.NewEvents(
		i.selectedSkill.ID,
		i.selectedSkill.Name,
		display.EventsFromDomain(events),
		i.Theme(),
	)
	i.skillEventsModal.SetDimensions(width, height)
	i.skillEventsModal.Show()

	return nil
}

// openEventDetailModal opens the event detail modal for a selected event.
func (i *Intent) openEventDetailModal(event *domain.Event) tea.Cmd {
	if event == nil {
		return nil
	}

	width, height := i.getTerminalDimensions()

	i.eventDetailModal = eventviews.NewDetail(display.EventFromDomain(event), i.Theme()).WithShowSkillsOption(false)
	i.eventDetailModal.SetDimensions(width, height)
	i.eventDetailModal.Show()

	return nil
}

// loadEventsForSkillModal loads events and opens the skill events modal.
func (i *Intent) loadEventsForSkillModal() tea.Cmd {
	return func() tea.Msg {
		events, err := i.context.SkillRepository.GetEventsUsingSkill(context.Background(), i.selectedSkill.ID)
		return SkillEventsForModalLoadedMsg{
			Events: events,
			Error:  err,
		}
	}
}

// startSkillInference triggers skill inference from all events.
func (i *Intent) startSkillInference() tea.Cmd {
	if i.context.SkillInferenceService == nil {
		i.feedbackModal = feedback.NewErrorModal("Inference Failed", "Skill inference service not available")
		return nil
	}

	if i.context.EventRepository == nil {
		i.feedbackModal = feedback.NewErrorModal("Inference Failed", "Event repository not available")
		return nil
	}

	i.state = StateInferringSkills
	i.loadingModal = feedback.NewLoadingModal("Analyzing all events for skills...", true).WithTheme(i.Theme())

	return tea.Batch(
		i.loadingModal.Init(),
		i.inferSkillsFromAllEvents(),
	)
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

	events := i.resolveEventsFromIDs(selected.EventIDs)

	width, height := i.getTerminalDimensions()
	i.suggestionEventsModal = skillviews.NewEvents(
		"suggestion",
		selected.Name,
		display.EventsFromDomain(events),
		i.Theme(),
	)
	i.suggestionEventsModal.SetDimensions(width, height)
	i.suggestionEventsModal.Show()

	return noopCmd
}

// resolveEventsFromIDs resolves event IDs to Event objects using the EventRepository.
func (i *Intent) resolveEventsFromIDs(eventIDs []string) []*domain.Event {
	if i.context.EventRepository == nil || len(eventIDs) == 0 {
		return []*domain.Event{}
	}

	events := make([]*domain.Event, 0, len(eventIDs))
	for _, id := range eventIDs {
		event, err := i.context.EventRepository.GetByID(context.Background(), id)
		if err == nil && event != nil {
			events = append(events, event)
		}
	}
	return events
}

// filterNewSuggestions removes suggestions whose names appear in existingNames.
// Returns only suggestions for skills not already tracked.
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

// inferSkillsFromAllEvents creates async command for skill inference from all events.
func (i *Intent) inferSkillsFromAllEvents() tea.Cmd {
	return func() tea.Msg {
		// Get all events from repository
		events, err := i.context.EventRepository.List(context.Background(), career.EventListFilters{})
		if err != nil {
			return SkillSuggestionsLoadedMsg{Error: err}
		}

		if len(events) == 0 {
			return SkillSuggestionsLoadedMsg{
				Error: errors.New("no events available for skill analysis"),
			}
		}

		result, err := i.context.SkillInferenceService.InferSkillsFromEvents(context.Background(), events)
		if err != nil {
			return SkillSuggestionsLoadedMsg{Error: err}
		}

		return SkillSuggestionsLoadedMsg{
			Suggestions:        result.Suggestions,
			ExistingSkillNames: result.ExistingSkillNames,
		}
	}
}
