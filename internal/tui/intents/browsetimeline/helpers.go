package browsetimeline

import (
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/domain/skills"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	"github.com/baphled/kariya/internal/tui/forms"
	"github.com/baphled/kariya/internal/tui/intents"
	burstviews "github.com/baphled/kariya/internal/tui/views/burst"
	eventviews "github.com/baphled/kariya/internal/tui/views/event"
	"github.com/baphled/kariya/internal/tui/views/shared"
	skillviews "github.com/baphled/kariya/internal/tui/views/skill"
	"github.com/baphled/kariya/internal/ui/behaviors"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/uikit/feedback"
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

// openViewSkillsModal returns a tea.Cmd that loads skills for the selected event asynchronously.
func (i *Intent) openViewSkillsModal() tea.Cmd {
	if i.selectedEvent == nil {
		return nil
	}
	ctx := i.getContext()
	eventID := i.selectedEvent.ID
	svc := i.context.CLIEventService
	return func() tea.Msg {
		if svc == nil {
			return SkillsForModalLoadedMsg{EventID: eventID, Skills: []*career.Skill{}}
		}
		eventSkills, err := svc.GetSkillsForEvent(ctx, eventID)
		if err != nil {
			return SkillsForModalLoadedMsg{EventID: eventID, Skills: nil, Error: err}
		}
		return SkillsForModalLoadedMsg{EventID: eventID, Skills: eventSkills}
	}
}

func (i *Intent) findEventByID(eventID string) *career.Event {
	if eventID == "" {
		return nil
	}

	for idx := range i.context.Events {
		evt := i.context.Events[idx]
		if evt != nil && evt.ID == eventID {
			return evt
		}
	}

	return nil
}

func (i *Intent) findSkillByID(skillID string, skillList []*career.Skill) *career.Skill {
	if skillID == "" {
		return nil
	}

	for idx := range skillList {
		skill := skillList[idx]
		if skill != nil && skill.ID == skillID {
			return skill
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

// openSearchModal creates and initializes the search modal via shared SearchView.
func (i *Intent) openSearchModal() tea.Cmd {
	width, height := i.getTerminalDimensions()
	searchView := shared.NewSearchView("Search events...")
	searchView.SetTerminalInfo(width, height)
	i.searchAdapter = intents.NewFormViewAdapter(searchView)
	i.searchAdapter.Show()
	return searchView.Init()
}

// openSortModal creates and initializes the sort modal via shared SortView.
func (i *Intent) openSortModal() tea.Cmd {
	width, height := i.getTerminalDimensions()
	sortFields := []shared.SortFieldOption{
		{Key: "date", Label: "Date"},
		{Key: "text", Label: "Text"},
	}
	sortView := shared.NewSortView(sortFields, i.filters.SortBy, i.filters.SortOrder)
	sortView.SetTerminalInfo(width, height)
	i.sortAdapter = intents.NewFormViewAdapter(sortView)
	i.sortAdapter.Show()
	return sortView.Init()
}

// openFilterModal creates and initializes the filter modal via event Filter view.
func (i *Intent) openFilterModal() tea.Cmd {
	width, height := i.getTerminalDimensions()
	fields := extractFilterFields(i.context.Events)
	sortFields := []shared.SortFieldOption{
		{Key: "date", Label: "Date"},
		{Key: "text", Label: "Text"},
	}
	defaults := &eventviews.FilterFormData{
		Selections: map[string][]string{
			"companies":  i.filters.Companies,
			"categories": i.filters.Categories,
			"projects":   i.filters.Projects,
		},
		SortBy:    i.filters.SortBy,
		SortOrder: i.filters.SortOrder,
	}
	filterView := eventviews.NewFilter(fields, sortFields, defaults)
	filterView.SetTerminalInfo(width, height)
	i.filterAdapter = intents.NewFormViewAdapter(filterView)
	i.filterAdapter.Show()
	return filterView.Init()
}

// extractFilterFields extracts unique companies, categories, and projects
// from events as FilterFieldConfig slices for the filter view.
func extractFilterFields(events []*career.Event) []eventviews.FilterFieldConfig {
	companyMap := make(map[string]bool)
	categoryMap := make(map[string]bool)
	projectMap := make(map[string]bool)

	for _, evt := range events {
		if evt.Company != "" {
			companyMap[evt.Company] = true
		}
		for _, cat := range evt.Categories {
			categoryMap[cat] = true
		}
		if evt.Project != "" {
			projectMap[evt.Project] = true
		}
	}

	var fields []eventviews.FilterFieldConfig

	if field := buildFilterField("companies", "Filter by Company", companyMap); field != nil {
		fields = append(fields, *field)
	}

	if field := buildFilterField("categories", "Filter by Category", categoryMap); field != nil {
		fields = append(fields, *field)
	}

	if field := buildFilterField("projects", "Filter by Project", projectMap); field != nil {
		fields = append(fields, *field)
	}

	return fields
}

// buildFilterField builds a filter field configuration from map values.
//
// Expected:
//   - values contains unique filter values.
//
// Returns:
//   - A FilterFieldConfig pointer, or nil when no values exist.
//
// Side effects:
//   - None.
func buildFilterField(key, title string, values map[string]bool) *eventviews.FilterFieldConfig {
	if len(values) == 0 {
		return nil
	}
	opts := make([]forms.SelectOption, 0, len(values))
	for value := range values {
		opts = append(opts, forms.SelectOption{Key: value, Value: value})
	}
	return &eventviews.FilterFieldConfig{Key: key, Title: title, Options: opts}
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

// RefreshData re-applies current filters and rebuilds the timeline screen
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (i *Intent) RefreshData() tea.Cmd {
	i.applyFilters()
	i.transitionToView(eventviews.NewListView(display.EventsFromDomain(i.filteredEvents)))
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
	return i.viewSkills != nil && i.viewSkills.IsVisible()
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
	return i.quickAdd != nil && i.quickAdd.IsVisible()
}

// HasVisibleEditModal checks whether the event edit modal is currently
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (i *Intent) HasVisibleEditModal() bool {
	return i.edit != nil && i.edit.IsVisible()
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

// rebuildModalRegistry creates a fresh modal registry with all current views.
// Call this whenever a modal is created or destroyed to keep the registry current.
func (i *Intent) rebuildModalRegistry() {
	if i.modalRegistry == nil {
		i.modalRegistry = intents.NewModalRegistry()
	}
	i.modalRegistry.Clear()

	// Register views in priority order (highest priority first).
	// Error modal has highest priority.
	if i.errorModal != nil {
		width, height := i.getTerminalDimensions()
		i.modalRegistry.Register(intents.NewErrorModalAdapter(i.errorModal, width, height, i.Theme()))
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

	if i.quickAdd != nil {
		i.modalRegistry.Register(intents.NewFormModalAdapter(
			i.quickAdd.IsVisible,
			i.quickAdd.View,
			i.quickAdd.Update,
		))
	}

	if i.edit != nil {
		i.modalRegistry.Register(intents.NewFormModalAdapter(
			i.edit.IsVisible,
			i.edit.View,
			i.edit.Update,
		))
	}

	// Confirm modal (delete).
	if i.deleteModal != nil {
		i.modalRegistry.Register(intents.NewConfirmModalAdapter(i.deleteModal))
	}

	// View overlays (priority order: picker > suggestion > skills > detail).
	if i.skillPicker != nil {
		i.modalRegistry.Register(intents.NewViewModalAdapter(
			i.skillPicker.IsVisible,
			i.skillPicker.View,
			i.skillPicker.Update,
		))
	}

	if i.skillSuggestionModal != nil {
		i.modalRegistry.Register(intents.NewViewModalAdapter(
			i.skillSuggestionModal.IsVisible,
			i.skillSuggestionModal.View,
			i.skillSuggestionModal.Update,
		))
	}

	if i.viewSkills != nil {
		i.modalRegistry.Register(intents.NewViewModalAdapter(
			i.viewSkills.IsVisible,
			i.viewSkills.View,
			i.viewSkills.Update,
		))
	}

	if i.viewDetail != nil {
		i.modalRegistry.Register(intents.NewViewModalAdapter(
			i.viewDetail.IsVisible,
			i.viewDetail.View,
			i.viewDetail.Update,
		))
	}
}

func (i *Intent) openSkillPickerModal() tea.Cmd {
	if i.selectedEvent == nil {
		return nil
	}
	ctx := i.getContext()
	eventID := i.selectedEvent.ID
	svc := i.context.CLIEventService
	return func() tea.Msg {
		allSkills, err := svc.ListAllSkills(ctx)
		if err != nil {
			return SkillPickerDataLoadedMsg{Error: err}
		}
		eventSkills, err := svc.GetSkillsForEvent(ctx, eventID)
		if err != nil {
			return SkillPickerDataLoadedMsg{Error: err}
		}
		return SkillPickerDataLoadedMsg{
			AllSkills:   allSkills,
			EventSkills: eventSkills,
			Error:       nil,
		}
	}
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
	eventID := i.selectedEvent.ID
	skillID := skill.ID
	svc := i.context.CLIEventService
	return func() tea.Msg {
		var err error
		if link {
			err = svc.LinkSkillToEvent(ctx, eventID, skillID)
			return SkillLinkedMsg{EventID: eventID, SkillID: skillID, Error: err}
		}
		err = svc.UnlinkSkillFromEvent(ctx, eventID, skillID)
		return SkillUnlinkedMsg{EventID: eventID, SkillID: skillID, Error: err}
	}
}

func (i *Intent) refreshSkillsModal() tea.Cmd {
	if i.selectedEvent == nil || i.viewSkills == nil {
		return nil
	}
	ctx := i.getContext()
	eventID := i.selectedEvent.ID
	svc := i.context.CLIEventService
	return func() tea.Msg {
		eventSkills, err := svc.GetSkillsForEvent(ctx, eventID)
		return SkillsRefreshedMsg{Skills: eventSkills, Error: err}
	}
}

func (i *Intent) openSkillAddModal() tea.Cmd {
	width, height := i.getTerminalDimensions()
	i.skillAddModal = skillviews.NewAddEdit(nil, width, height)
	return i.skillAddModal.Init()
}

func (i *Intent) createAndLinkSkill(skillData *skillviews.EditData) tea.Cmd {
	if i.selectedEvent == nil || skillData == nil {
		return nil
	}
	input := skills.SkillInput{
		Name:      skillData.Name,
		Category:  skillData.Category,
		Level:     skillData.Level,
		YearsUsed: skillData.YearsUsed,
	}
	newSkill, err := skills.NewSkillFromInput(input, "")
	if err != nil {
		i.ShowErrorModal("Invalid Skill Data", err.Error())
		return nil
	}
	ctx := i.getContext()
	eventID := i.selectedEvent.ID
	svc := i.context.CLIEventService
	skillSvc := i.context.CLISkillCreator
	return func() tea.Msg {
		if err := skillSvc.Create(ctx, newSkill); err != nil {
			return SkillCreatedMsg{Skill: nil, Error: err}
		}
		if err := svc.LinkSkillToEvent(ctx, eventID, newSkill.ID); err != nil {
			return SkillCreatedMsg{Skill: newSkill, Error: err}
		}
		return SkillCreatedMsg{Skill: newSkill}
	}
}

func (i *Intent) inferSkillsFromEvent() tea.Cmd {
	if i.selectedEvent == nil || i.context.SkillInferenceService == nil {
		return nil
	}
	ctx, service, selectedEvent := i.getContext(), i.context.SkillInferenceService, i.selectedEvent
	return func() tea.Msg {
		result, err := service.InferSkillsFromEvents(ctx, []*career.Event{selectedEvent})
		if err != nil {
			return SkillSuggestionsErrorMsg{Error: err}
		}
		return SkillSuggestionsLoadedMsg{Suggestions: result.Suggestions, ExistingSkillNames: result.ExistingSkillNames}
	}
}

func (i *Intent) handleSkillSuggestionsLoaded(msg SkillSuggestionsLoadedMsg) {
	if len(msg.Suggestions) == 0 {
		i.ShowErrorModal("No Skills Detected", "No skills detected from event.")
		return
	}
	newSuggestions := filterNewSkillSuggestions(msg.Suggestions, msg.ExistingSkillNames)
	if len(newSuggestions) == 0 {
		i.ShowErrorModal("All Skills Tracked", "All detected skills already in profile.")
		return
	}
	width, height := i.getTerminalDimensions()
	i.skillSuggestions = append([]skillinference.SkillSuggestion(nil), newSuggestions...)
	i.skillSuggestionModal = burstviews.NewSkillSuggestion(display.SkillSuggestionsFromDomain(newSuggestions), i.Theme())
	i.skillSuggestionModal.SetDimensions(width, height)
	i.skillSuggestionModal.Show()
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
	createdSkills, err := i.context.SkillInferenceService.CreateSkillsFromSuggestions(ctx, []skillinference.SkillSuggestion{suggestion})
	if err != nil {
		i.ShowErrorModal("Skill Creation Failed", err.Error())
		return
	}
	if i.selectedEvent != nil && len(createdSkills) > 0 {
		if err := i.context.CLIEventService.LinkSkillToEvent(ctx, i.selectedEvent.ID, createdSkills[0].ID); err != nil {
			i.ShowErrorModal("Error Linking Skill", err.Error())
		}
	}
}

// handleSkillsForModalLoaded applies skills data to the view skills modal.
//
// Expected:
//   - msg contains skills for the selected event.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - Updates skill modal state or shows error modal.
func (i *Intent) handleSkillsForModalLoaded(msg SkillsForModalLoadedMsg) {
	if msg.Error != nil {
		i.ShowErrorModal(errorLoadingSkillsTitle, msg.Error.Error())
		return
	}

	i.selectedEventSkills = msg.Skills
	width, height := i.getTerminalDimensions()
	theme := i.Theme()
	i.viewSkills = eventviews.NewSkillsDetail(msg.EventID, display.SkillsFromDomain(msg.Skills), theme)
	i.viewSkills.SetDimensions(width, height)
	i.viewSkills.Show()
}

// handleSkillPickerDataLoaded prepares the skill picker modal with available skills.
//
// Expected:
//   - msg contains skills for the current event and all available skills.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - Updates skill picker modal state or shows error modal.
func (i *Intent) handleSkillPickerDataLoaded(msg SkillPickerDataLoadedMsg) {
	if msg.Error != nil {
		i.ShowErrorModal(errorLoadingSkillsTitle, msg.Error.Error())
		return
	}

	i.availableSkills = msg.AllSkills
	i.selectedEventSkills = msg.EventSkills
	eventSkillIDs := make(map[string]bool)
	for idx := range msg.EventSkills {
		s := msg.EventSkills[idx]
		eventSkillIDs[s.ID] = true
	}
	availableSkills := make([]display.Skill, 0, len(msg.AllSkills))
	for idx := range msg.AllSkills {
		s := msg.AllSkills[idx]
		if !eventSkillIDs[s.ID] {
			availableSkills = append(availableSkills, display.SkillFromDomain(s))
		}
	}
	width, height := i.getTerminalDimensions()
	i.skillPicker = eventviews.NewSkillPicker(availableSkills, i.Theme())
	i.skillPicker.SetDimensions(width, height)
	i.skillPicker.Show()
}

// handleSkillsRefreshed refreshes the skills shown in the skills modal.
//
// Expected:
//   - msg contains updated skill data for the current event.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - Updates modal state or shows error modal.
func (i *Intent) handleSkillsRefreshed(msg SkillsRefreshedMsg) {
	if msg.Error != nil {
		i.ShowErrorModal(errorLoadingSkillsTitle, msg.Error.Error())
		return
	}

	i.selectedEventSkills = msg.Skills
	if i.viewSkills != nil {
		i.viewSkills.SetSkills(display.SkillsFromDomain(msg.Skills))
	}
}

// buildModalHandlers returns the ordered slice of modal handlers for the intent.
func (i *Intent) buildModalHandlers() []modalHandler {
	return []modalHandler{
		{
			isActive: func() bool { return i.searchAdapter != nil && i.searchAdapter.IsVisible() },
			update:   i.updateSearchModal,
			isClosed: func() bool { return i.searchAdapter == nil || !i.searchAdapter.IsVisible() },
		},
		{
			isActive: func() bool { return i.filterAdapter != nil && i.filterAdapter.IsVisible() },
			update:   i.updateFilterModal,
			isClosed: func() bool { return i.filterAdapter == nil || !i.filterAdapter.IsVisible() },
		},
		{
			isActive: func() bool { return i.sortAdapter != nil && i.sortAdapter.IsVisible() },
			update:   i.updateSortModal,
			isClosed: func() bool { return i.sortAdapter == nil || !i.sortAdapter.IsVisible() },
		},
		{
			isActive: func() bool { return i.quickAdd != nil && i.quickAdd.IsVisible() },
			update:   i.updateQuickAddModal,
			isClosed: func() bool { return i.quickAdd == nil || !i.quickAdd.IsVisible() },
		},
		{
			isActive: func() bool { return i.edit != nil && i.edit.IsVisible() },
			update:   i.updateEditModal,
			isClosed: func() bool { return i.edit == nil || !i.edit.IsVisible() },
		},
		{
			isActive: func() bool { return i.deleteModal != nil && i.deleteModal.IsVisible() },
			update:   i.updateDeleteModal,
			isClosed: func() bool { return i.deleteModal == nil || !i.deleteModal.IsVisible() },
		},
		{
			isActive: func() bool { return i.skillPicker != nil && i.skillPicker.IsVisible() },
			update:   i.updateSkillPickerModal,
			isClosed: func() bool { return i.skillPicker == nil || !i.skillPicker.IsVisible() },
		},
		{
			isActive: func() bool { return i.skillAddModal != nil && i.skillAddModal.IsVisible() },
			update:   i.updateSkillAddModal,
			isClosed: func() bool { return i.skillAddModal == nil || !i.skillAddModal.IsVisible() },
		},
		{
			isActive: func() bool {
				return i.skillSuggestionModal != nil && i.skillSuggestionModal.IsVisible()
			},
			update:   i.updateSkillSuggestionModal,
			isClosed: func() bool { return i.skillSuggestionModal == nil || !i.skillSuggestionModal.IsVisible() },
		},
		{
			isActive: func() bool { return i.viewSkills != nil && i.viewSkills.IsVisible() },
			update:   i.updateViewSkillsModal,
			isClosed: func() bool { return i.viewSkills == nil || !i.viewSkills.IsVisible() },
		},
		{
			isActive: func() bool { return i.viewDetail != nil && i.viewDetail.IsVisible() },
			update:   i.updateViewDetailModal,
			isClosed: func() bool { return i.viewDetail == nil || !i.viewDetail.IsVisible() },
		},
	}
}
