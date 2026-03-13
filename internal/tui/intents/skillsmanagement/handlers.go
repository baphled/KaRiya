package skillsmanagement

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	domain "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/domain/skills"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	"github.com/baphled/kariya/internal/tui/intents"
	burstviews "github.com/baphled/kariya/internal/tui/views/burst"
	"github.com/baphled/kariya/internal/tui/views/shared"
	skillview "github.com/baphled/kariya/internal/tui/views/skill"
	"github.com/baphled/kariya/internal/ui/behaviors"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/uikit/feedback"
	tea "github.com/charmbracelet/bubbletea"
)

// noopCmd is a command that does nothing but prevents message propagation.
func noopCmd() tea.Msg { return nil }

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
func (i *Intent) handleLoadingModalUpdate(msg tea.Msg) tea.Cmd {
	if i.loadingModal == nil {
		return nil
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc {
			i.loadingModal = nil
			i.state = StateList
			return noopCmd
		}
		return noopCmd
	case feedback.ModalSpinnerTickMsg:
		return i.loadingModal.Update(msg)
	default:
		return noopCmd
	}
}

// handleSkillsLoaded handles the SkillsLoadedMsg.
func (i *Intent) handleSkillsLoaded(msg SkillsLoadedMsg) {
	if msg.Error != nil {
		i.result = &intents.IntentResult[*Result]{
			Status: intents.Failed,
			Error: &intents.IntentError{
				Code:    "LOAD_FAILED",
				Message: "Failed to load skills",
				Cause:   msg.Error,
			},
		}
		return
	}

	i.skills = msg.Skills
	i.selectedIndex = 0

	// Update TableBehavior with new skills data.
	i.tableBehavior.SetItems(i.skills)
	i.syncTableSelection()

	// Load event counts.
	if counts, err := i.context.GetEventCounts(); err == nil {
		i.eventCounts = counts
		// Recreate table with event counts formatter.
		columns := []behaviors.ColumnDef{
			{Title: "Name", Width: 25},
			{Title: "Category", Width: 15},
			{Title: "Level", Width: 12},
			{Title: "Years", Width: 8},
			{Title: "Events", Width: 8},
		}
		i.tableBehavior = behaviors.NewTableBehavior(
			nil, columns, skillRowFormatterWithCounts(i.eventCounts),
		).PageSize(15).PaginationPrefix("Skills").EmptyMessage("No skills found. Press 'a' to add a new skill.")

		i.tableBehavior.SetItems(i.skills)
		if theme := i.Theme(); theme != nil {
			i.tableBehavior.SetTheme(theme)
		}
	}

	// Create a new list view with updated data.
	view := skillview.NewListView(display.SkillsFromDomain(i.skills), i.eventCounts)
	i.transitionToView(view)
}

// handleSkillCreated handles the SkillCreatedMsg.
func (i *Intent) handleSkillCreated(msg SkillCreatedMsg) tea.Cmd {
	if msg.Error != nil {
		i.result = &intents.IntentResult[*Result]{
			Status: intents.Failed,
			Error: &intents.IntentError{
				Code:    "CREATE_FAILED",
				Message: "Failed to create skill",
				Cause:   msg.Error,
			},
		}
		return nil
	}

	i.state = StateList
	return i.Init()
}

// handleSkillUpdated handles the SkillUpdatedMsg.
func (i *Intent) handleSkillUpdated(msg SkillUpdatedMsg) tea.Cmd {
	if msg.Error != nil {
		i.result = &intents.IntentResult[*Result]{
			Status: intents.Failed,
			Error: &intents.IntentError{
				Code:    "UPDATE_FAILED",
				Message: "Failed to update skill",
				Cause:   msg.Error,
			},
		}
		return nil
	}

	i.state = StateList
	return i.Init()
}

// handleSkillDeleted handles the SkillDeletedMsg.
func (i *Intent) handleSkillDeleted(msg SkillDeletedMsg) tea.Cmd {
	if msg.Error != nil {
		i.result = &intents.IntentResult[*Result]{
			Status: intents.Failed,
			Error: &intents.IntentError{
				Code:    "DELETE_FAILED",
				Message: "Failed to delete skill",
				Cause:   msg.Error,
			},
		}
		return nil
	}

	i.state = StateList
	return i.Init()
}

// handleSkillEventsLoaded handles the SkillEventsLoadedMsg.
func (i *Intent) handleSkillEventsLoaded(msg SkillEventsLoadedMsg) {
	if msg.Error != nil {
		return
	}

	i.skillEvents = msg.Events
	i.eventsLoaded = true
}

// handleSkillEventsForModalLoaded handles the SkillEventsForModalLoadedMsg (modal flow).
func (i *Intent) handleSkillEventsForModalLoaded(msg SkillEventsForModalLoadedMsg) tea.Cmd {
	if msg.Error != nil {
		return nil
	}

	// Open the skill events modal with the loaded events.
	return i.openSkillEventsModal(msg.Events)
}

// handleKeyShortcuts handles keyboard shortcuts when no modal is active.
// These are intent-level shortcuts that open modals.
func (i *Intent) handleKeyShortcuts(keyMsg tea.KeyMsg) tea.Cmd {
	// Handle back key at root state.
	if intents.HandleGlobalKeys(keyMsg) == intents.KeyBack {
		i.SetCancelled()
		return nil
	}

	switch keyMsg.String() {
	case "f":
		return i.openFilterModal()
	case "s":
		return i.openSortModal()
	case "/":
		return i.openSearchModal()
	case "i":
		return i.startSkillInference()
	case "x":
		if i.HasActiveFilters() {
			i.ClearFilters()
			return i.RefreshData()
		}
	}
	return nil
}

// handleNavigateData processes navigation data from screen results.
func (i *Intent) handleNavigateData(actionData map[string]interface{}) tea.Cmd {
	action, _ := actionData["action"].(string)
	switch action {
	case "view":
		if skill, ok := actionData["skill"].(display.Skill); ok {
			i.selectedSkill = i.findSkillByID(skill.ID)
			return i.openViewDetailModal()
		}
		return nil

	case "add":
		return i.openAddEditModal(nil)

	case "edit":
		if skill, ok := actionData["skill"].(display.Skill); ok {
			i.selectedSkill = i.findSkillByID(skill.ID)
			return i.openAddEditModal(i.selectedSkill)
		}
		return nil

	case "delete":
		if skill, ok := actionData["skill"].(display.Skill); ok {
			selectedSkill := i.findSkillByID(skill.ID)
			return i.openDeleteModal(selectedSkill)
		}
		return nil

	default:
		return nil
	}
}

// Modal update handlers.
// These handle updates when modals are visible.

// parseYears converts a string to int, returning 0 for empty or invalid input.
func parseYears(s string) int {
	if s == "" {
		return 0
	}
	v, err := strconv.Atoi(s)
	if err != nil || v < 0 {
		return 0
	}
	return v
}

// handleFilterModalUpdate handles updates when filter adapter is visible.
func (i *Intent) handleFilterModalUpdate(msg tea.Msg) tea.Cmd {
	result := i.filterAdapter.HandleUpdate(msg)
	if result.Applied {
		if filterData, ok := result.Data.(skillview.FilterFormData); ok {
			return i.applyFilterSelections(filterData)
		}
	}
	if result.Closed && !result.Applied {
		i.filterAdapter = nil
	}
	return result.Cmd
}

// applyFilterSelections applies filter form data to the intent's filters.
func (i *Intent) applyFilterSelections(filterData skillview.FilterFormData) tea.Cmd {
	if i.context.Filters == nil {
		i.context.Filters = &Filters{}
	}
	if categories, exists := filterData.Selections["categories"]; exists && len(categories) > 0 {
		i.context.Filters.Category = categories[0]
	} else {
		i.context.Filters.Category = ""
	}
	if levels, exists := filterData.Selections["levels"]; exists && len(levels) > 0 {
		i.context.Filters.Level = levels[0]
	} else {
		i.context.Filters.Level = ""
	}
	i.context.Filters.MinEvents = parseYears(filterData.MinYearsStr)
	return i.reloadSkills()
}

// handleSortModalUpdate handles updates when sort adapter is visible.
func (i *Intent) handleSortModalUpdate(msg tea.Msg) tea.Cmd {
	result := i.sortAdapter.HandleUpdate(msg)
	if result.Applied {
		if sortData, ok := result.Data.(shared.SortFormData); ok {
			if i.context.Filters == nil {
				i.context.Filters = &Filters{}
			}
			i.context.Filters.SortBy = sortData.SortBy
			i.context.Filters.SortOrder = sortData.SortOrder
			return i.reloadSkills()
		}
	}
	if result.Closed && !result.Applied {
		i.sortAdapter = nil
	}
	return result.Cmd
}

// handleSearchModalUpdate handles updates when search adapter is visible.
func (i *Intent) handleSearchModalUpdate(msg tea.Msg) tea.Cmd {
	result := i.searchAdapter.HandleUpdate(msg)
	if result.Applied {
		if searchData, ok := result.Data.(shared.SearchFormData); ok {
			if i.context.Filters == nil {
				i.context.Filters = &Filters{}
			}
			i.context.Filters.SearchText = searchData.SearchText
			return i.reloadSkills()
		}
	}
	if result.Closed && !result.Applied {
		i.searchAdapter = nil
	}
	return result.Cmd
}

// handleViewDetailModalUpdate handles updates when view detail modal is visible.
func (i *Intent) handleViewDetailModalUpdate(msg tea.Msg) tea.Cmd {
	_, cmd := i.viewDetailModal.Update(msg)

	if !i.viewDetailModal.IsVisible() {
		// Modal was closed - check for action.
		action := i.viewDetailModal.GetAction()
		switch action {
		case "events":
			// Load events for this skill and show modal.
			i.viewDetailModal = nil
			return i.loadEventsForSkillModal()
		case "edit":
			// Open add/edit modal for this skill.
			return i.openAddEditModal(i.selectedSkill)
		case "delete":
			// Open delete confirmation modal.
			return i.openDeleteModal(i.selectedSkill)
		}
		// Simple close - clear the modal.
		i.viewDetailModal = nil
	}

	return cmd
}

// handleAddEditModalUpdate handles updates when add/edit modal is visible.
func (i *Intent) handleAddEditModalUpdate(msg tea.Msg) tea.Cmd {
	cmd, completed, skillData := i.addEditModal.Update(msg)

	if !i.addEditModal.IsVisible() {
		if completed && skillData != nil {
			input := skills.SkillInput{
				Name:      skillData.Name,
				Category:  skillData.Category,
				Level:     skillData.Level,
				YearsUsed: skillData.YearsUsed,
			}
			skillID := i.addEditModal.GetOriginalSkillID()
			originalSkill := i.findSkillByID(skillID)
			skill, err := skills.NewSkillFromInput(input, skillID)
			if err != nil {
				i.addEditModal = nil
				i.feedbackModal = feedback.NewErrorModal("Invalid Skill Data", err.Error())
				i.state = StateList
				return nil
			}
			i.addEditModal = nil
			if originalSkill != nil {
				return i.updateSkill(skill)
			}
			return i.createSkill(skill)
		}
		i.addEditModal = nil
	}

	return cmd
}

// handleDeleteModalUpdate handles updates when delete confirmation modal is visible.
func (i *Intent) handleDeleteModalUpdate(msg tea.Msg) tea.Cmd {
	cmd, confirmed := i.deleteModal.Update(msg)

	if !i.deleteModal.IsVisible() {
		if confirmed && i.selectedSkill != nil {
			// User confirmed deletion.
			skillID := i.selectedSkill.ID
			i.deleteModal = nil
			return func() tea.Msg {
				err := i.context.SkillRepository.Delete(context.Background(), skillID)
				return SkillDeletedMsg{
					SkillID: skillID,
					Error:   err,
				}
			}
		}
		// User cancelled.
		i.deleteModal = nil
	}

	return cmd
}

// handleSkillEventsModalUpdate handles updates when skill events modal is visible.
func (i *Intent) handleSkillEventsModalUpdate(msg tea.Msg) tea.Cmd {
	_, cmd := i.skillEventsModal.Update(msg)

	if !i.skillEventsModal.IsVisible() {
		// Check if user selected an event.
		if i.skillEventsModal.HasSelection() {
			selectedEvent := i.skillEventsModal.GetSelectedEvent()
			i.skillEventsModal.ClearSelection()
			if selectedEvent != nil {
				return i.openEventDetailModal(i.findEventByID(selectedEvent.ID))
			}
			return nil
		}
		// Simple close - clear the modal.
		i.skillEventsModal = nil
	}

	return cmd
}

// handleEventDetailModalUpdate handles updates when event detail modal is visible.
func (i *Intent) handleEventDetailModalUpdate(msg tea.Msg) tea.Cmd {
	_, cmd := i.eventDetailModal.Update(msg)

	if !i.eventDetailModal.IsVisible() {
		// Event detail modal was closed.
		i.eventDetailModal = nil
		// Re-show the skill events modal if it exists.
		if i.skillEventsModal != nil {
			i.skillEventsModal.Show()
		}
	}

	return cmd
}

// createSkill creates a new skill.
func (i *Intent) createSkill(skill *domain.Skill) tea.Cmd {
	return func() tea.Msg {
		err := i.context.SkillRepository.Create(context.Background(), skill)
		return SkillCreatedMsg{
			Skill: skill,
			Error: err,
		}
	}
}

// updateSkill updates an existing skill.
func (i *Intent) updateSkill(skill *domain.Skill) tea.Cmd {
	return func() tea.Msg {
		err := i.context.SkillRepository.Update(context.Background(), skill)
		return SkillUpdatedMsg{
			Skill: skill,
			Error: err,
		}
	}
}

// handleSkillSuggestionsLoaded handles the SkillSuggestionsLoadedMsg.
func (i *Intent) handleSkillSuggestionsLoaded(msg SkillSuggestionsLoadedMsg) {
	i.loadingModal = nil

	if msg.Error != nil {
		// Silently ignore cancelled operations.
		if errors.Is(msg.Error, context.Canceled) {
			i.state = StateList
			return
		}
		i.feedbackModal = feedback.NewErrorModal("Skill Inference Failed", msg.Error.Error())
		i.state = StateList
		return
	}

	if len(msg.Suggestions) == 0 {
		i.feedbackModal = feedback.NewWarningModal("No Skills Found", "No skills were detected from the events")
		i.state = StateList
		return
	}

	newSuggestions := filterNewSuggestions(msg.Suggestions, msg.ExistingSkillNames)
	if len(newSuggestions) == 0 {
		successModal := feedback.NewSuccessModal(
			"Detected skills already in your profile: " +
				strings.Join(msg.ExistingSkillNames, ", "))
		successModal.Title = "All Skills Already Tracked"
		i.feedbackModal = successModal
		i.state = StateList
		return
	}

	i.skillSuggestions = append([]skillinference.SkillSuggestion(nil), newSuggestions...)
	i.skillSuggestionModal = burstviews.NewSkillSuggestion(display.SkillSuggestionsFromDomain(newSuggestions), i.Theme())
	width, height := i.getTerminalDimensions()
	i.skillSuggestionModal.SetDimensions(width, height)
	i.skillSuggestionModal.Show()
	i.state = StateSkillSuggestionReview
}

// handleSkillSuggestionModalUpdate handles skill suggestion modal updates.
func (i *Intent) handleSkillSuggestionModalUpdate(msg tea.Msg) tea.Cmd {
	if i.skillSuggestionModal == nil || !i.skillSuggestionModal.IsVisible() {
		return nil
	}

	_, cmd := i.skillSuggestionModal.Update(msg)

	if !i.skillSuggestionModal.IsVisible() {
		action := i.skillSuggestionModal.GetAction()

		switch action {
		case burstviews.SuggestionActionAccept, burstviews.SuggestionActionReject:
			return i.processAcceptedSkillSuggestions(cmd)

		case burstviews.SuggestionActionViewEvents:
			return i.openSuggestionEventsModal()

		case burstviews.SuggestionActionCancel:
			i.skillSuggestionModal = nil
			i.state = StateList
			return noopCmd
		}
	}

	return noopCmd
}

// processAcceptedSkillSuggestions collects accepted skill suggestions and creates them.
func (i *Intent) processAcceptedSkillSuggestions(cmd tea.Cmd) tea.Cmd {
	acceptedDisplay := i.skillSuggestionModal.GetAcceptedSkills()
	accepted := make([]skillinference.SkillSuggestion, 0, len(acceptedDisplay))
	for idx := range acceptedDisplay {
		domainSuggestion := i.findSkillSuggestionByName(acceptedDisplay[idx].Name)
		if domainSuggestion != nil {
			accepted = append(accepted, *domainSuggestion)
		}
	}

	if len(accepted) > 0 {
		i.loadingModal = feedback.NewLoadingModal(
			fmt.Sprintf("Creating %d skill(s)...", len(accepted)),
			true,
		).WithTheme(i.Theme())
		i.skillSuggestionModal = nil
		return tea.Batch(cmd, i.loadingModal.Init(), i.createSkillsFromSuggestions(accepted))
	}

	i.skillSuggestionModal = nil
	i.state = StateList
	return noopCmd
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

// createSkillsFromSuggestions persists accepted skill suggestions as confirmed skills.
func (i *Intent) createSkillsFromSuggestions(suggestions []skillinference.SkillSuggestion) tea.Cmd {
	if len(suggestions) == 0 {
		return func() tea.Msg {
			return SkillsCreatedMsg{Skills: nil, Error: nil}
		}
	}

	service := i.context.SkillInferenceService

	return func() tea.Msg {
		if service == nil {
			return SkillsCreatedMsg{Error: errors.New("skill inference service not available")}
		}

		createdSkills, err := service.CreateSkillsFromSuggestions(context.Background(), suggestions)
		if err != nil {
			return SkillsCreatedMsg{Error: fmt.Errorf("failed to create skills: %w", err)}
		}

		return SkillsCreatedMsg{
			Skills: createdSkills,
			Error:  nil,
		}
	}
}

// handleSkillsCreatedFromInference handles skills created from accepted suggestions.
func (i *Intent) handleSkillsCreatedFromInference(msg SkillsCreatedMsg) tea.Cmd {
	i.loadingModal = nil

	if msg.Error != nil {
		i.feedbackModal = feedback.NewErrorModal("Skill Creation Failed", msg.Error.Error())
		i.state = StateList
		return nil
	}

	i.feedbackModal = feedback.NewSuccessModal(fmt.Sprintf("Successfully created %d skill(s)", len(msg.Skills)))
	i.state = StateList
	return i.RefreshData()
}
