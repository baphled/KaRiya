package skillsmanagement

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/screens/skills/modals"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	domain "github.com/baphled/kariya/internal/domain/career"
	career "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/service/career/skillinference"
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
				eventCount = fmt.Sprintf("%d", count)
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

// getContextHelp returns the help text for the current state.
func (i *Intent) getContextHelp() string {
	theme := i.Theme()

	switch i.state {
	case StateList:
		return intents.CombineThemedFooters(
			intents.ThemedCustomFooter(theme,
				primitives.NavigateBadge(theme),
				primitives.SelectBadge(theme),
				primitives.InferSkillsBadge(theme),
				primitives.SearchBadge(theme),
				primitives.BackBadge(theme),
			),
			intents.ThemedGlobalBadges(theme),
		)
	case StateDetail:
		return intents.CombineThemedFooters(
			intents.ThemedDetailViewFooter(theme),
			intents.ThemedGlobalBadges(theme),
		)
	case StateInferringSkills, StateSkillSuggestionReview:
		return intents.ThemedCustomFooter(theme,
			primitives.CancelBadge(theme),
		)
	case StateDelete:
		return intents.ThemedGlobalBadges(theme)
	default:
		return intents.ThemedGlobalBadges(theme)
	}
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

	// Form modals (search, filter, sort, add/edit).
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

// openFilterModal opens the filter modal with current filters pre-populated.
func (i *Intent) openFilterModal() tea.Cmd {
	width, height := i.getTerminalDimensions()

	// Build current filters for pre-population.
	var currentFilters *modals.Filters
	if i.context.Filters != nil {
		currentFilters = &modals.Filters{
			Categories: []string{},
			Levels:     []string{},
			MinYears:   i.context.Filters.MinEvents,
			MaxYears:   0,
		}
		if i.context.Filters.Category != "" {
			currentFilters.Categories = []string{i.context.Filters.Category}
		}
		if i.context.Filters.Level != "" {
			currentFilters.Levels = []string{i.context.Filters.Level}
		}
	}

	// Create filter modal.
	i.filterModal = modals.NewFilterModal(
		i.skills,
		currentFilters,
		width,
		height,
	)

	// Call Init() for immediate rendering.
	return i.filterModal.Init()
}

// openSortModal opens the sort modal with current sort config pre-populated.
func (i *Intent) openSortModal() tea.Cmd {
	width, height := i.getTerminalDimensions()

	// Build current sort config for pre-population.
	var currentSort *modals.SortConfig
	if i.context.Filters != nil {
		currentSort = &modals.SortConfig{
			SortBy:    i.context.Filters.SortBy,
			SortOrder: i.context.Filters.SortOrder,
		}
	}

	// Create sort modal.
	i.sortModal = modals.NewSortModal(
		i.skills,
		currentSort,
		width,
		height,
	)

	// Call Init() for immediate rendering.
	return i.sortModal.Init()
}

// openSearchModal opens the search modal with current search text pre-populated.
func (i *Intent) openSearchModal() tea.Cmd {
	width, height := i.getTerminalDimensions()

	// Get current search text.
	searchText := ""
	if i.context.Filters != nil {
		searchText = i.context.Filters.SearchText
	}

	// Create search modal.
	i.searchModal = modals.NewSearchModal(
		searchText,
		width,
		height,
	)

	// Call Init() for immediate rendering.
	return i.searchModal.Init()
}

// openViewDetailModal opens the view detail modal for the selected skill.
// BUG-016: Uses i.selectedSkill (set by handleNavigateData) instead of
// indexing i.skills[i.selectedIndex], which was stale and pointed to
// the wrong skill.
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

	i.viewDetailModal = modals.NewDetailModal(skill, i.Theme(), eventCount, nil)
	i.viewDetailModal.SetDimensions(width, height)
	i.viewDetailModal.Show()

	return nil
}

// openAddEditModal opens the add/edit modal for a skill.
func (i *Intent) openAddEditModal(skill *domain.Skill) tea.Cmd {
	width, height := i.getTerminalDimensions()

	i.addEditModal = modals.NewAddEditModal(skill, width, height)
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

	width, height := i.getTerminalDimensions()

	i.skillEventsModal = modals.NewEventsModal(
		i.selectedSkill.ID,
		i.selectedSkill.Name,
		events,
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

	i.eventDetailModal = components.NewViewEventDetailModal(event, i.Theme()).
		WithShowSkillsOption(false)
	i.eventDetailModal.SetDimensions(width, height)
	i.eventDetailModal.Show()

	return nil
}

// loadEventsForSkillModal loads events and opens the skill events modal.
func (i *Intent) loadEventsForSkillModal() tea.Cmd {
	return func() tea.Msg {
		events, err := i.context.SkillRepository.GetEventsUsingSkill(i.context.Ctx, i.selectedSkill.ID)
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
	i.suggestionEventsModal = modals.NewEventsModal(
		"suggestion",
		selected.Name,
		events,
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
		event, err := i.context.EventRepository.GetByID(i.context.Ctx, id)
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
		events, err := i.context.EventRepository.List(i.context.Ctx, career.EventListFilters{})
		if err != nil {
			return SkillSuggestionsLoadedMsg{Error: err}
		}

		if len(events) == 0 {
			return SkillSuggestionsLoadedMsg{
				Error: fmt.Errorf("no events available for skill analysis"),
			}
		}

		result, err := i.context.SkillInferenceService.InferSkillsFromEvents(i.context.Ctx, events)
		if err != nil {
			return SkillSuggestionsLoadedMsg{Error: err}
		}

		return SkillSuggestionsLoadedMsg{
			Suggestions:        result.Suggestions,
			ExistingSkillNames: result.ExistingSkillNames,
		}
	}
}
