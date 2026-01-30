package skillsmanagement

import (
	"context"
	"errors"
	"fmt"

	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/screens"
	burstmodals "github.com/baphled/kariya/internal/cli/screens/burst_management/modals"
	"github.com/baphled/kariya/internal/cli/screens/skills"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	domain "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/skillinference"
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
func (i *Intent) handleSkillsLoaded(msg SkillsLoadedMsg) tea.Cmd {
	if msg.Error != nil {
		i.result = &intents.IntentResult[*Result]{
			Status: intents.Failed,
			Error: &intents.IntentError{
				Code:    "LOAD_FAILED",
				Message: "Failed to load skills",
				Cause:   msg.Error,
			},
		}
		return nil
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

	// Create a new list screen with updated data.
	i.listScreen = skills.NewSkillsListScreen(i.skills)
	if i.eventCounts != nil {
		i.listScreen.SetEventCounts(i.eventCounts)
	}
	i.transitionToScreen(i.listScreen)

	return nil
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
func (i *Intent) handleSkillEventsLoaded(msg SkillEventsLoadedMsg) tea.Cmd {
	if msg.Error != nil {
		return nil
	}

	i.skillEvents = msg.Events
	i.eventsLoaded = true

	return nil
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

// handleScreenResult processes a screen result and determines next action.
func (i *Intent) handleScreenResult(result interface{}) tea.Cmd {
	if result == nil {
		return nil
	}

	screenResult, ok := result.(screens.ScreenResult)
	if !ok {
		return nil
	}

	return behaviors.NewScreenResultDispatcher(i).Dispatch(screenResult)
}

// handleNavigateData processes navigation data from screen results.
func (i *Intent) handleNavigateData(actionData map[string]interface{}) tea.Cmd {
	action, _ := actionData["action"].(string)
	switch action {
	case "view":
		if skill, ok := actionData["skill"].(*domain.Skill); ok {
			i.selectedSkill = skill
			return i.openViewDetailModal()
		}
		return nil

	case "add":
		return i.openAddEditModal(nil)

	case "edit":
		if skill, ok := actionData["skill"].(*domain.Skill); ok {
			i.selectedSkill = skill
			return i.openAddEditModal(skill)
		}
		return nil

	case "delete":
		if skill, ok := actionData["skill"].(*domain.Skill); ok {
			return i.openDeleteModal(skill)
		}
		return nil

	default:
		return nil
	}
}

// Modal update handlers.
// These handle updates when modals are visible.

// handleFilterModalUpdate handles updates when filter modal is visible.
// Takes tea.Msg (not tea.KeyMsg) to allow huh forms to work correctly.
func (i *Intent) handleFilterModalUpdate(msg tea.Msg) tea.Cmd {
	cmd, applied, filterData := i.filterModal.Update(msg)

	if applied && filterData != nil {
		// User confirmed filters - apply them.
		newFilters := i.filterModal.ToFilters()

		// Update internal filter state.
		if i.context.Filters == nil {
			i.context.Filters = &Filters{}
		}
		// Map Filters to internal Filters format.
		if len(newFilters.Categories) > 0 {
			i.context.Filters.Category = newFilters.Categories[0]
		} else {
			i.context.Filters.Category = ""
		}
		if len(newFilters.Levels) > 0 {
			i.context.Filters.Level = newFilters.Levels[0]
		} else {
			i.context.Filters.Level = ""
		}
		i.context.Filters.MinEvents = newFilters.MinYears

		// Reload skills with new filters.
		return i.reloadSkills()
	}

	// Modal was closed without completion (Esc) or still being edited.
	return cmd
}

// handleSortModalUpdate handles updates when sort modal is visible.
// Takes tea.Msg (not tea.KeyMsg) to allow huh forms to work correctly.
func (i *Intent) handleSortModalUpdate(msg tea.Msg) tea.Cmd {
	cmd, applied, sortData := i.sortModal.Update(msg)

	if applied && sortData != nil {
		// User confirmed sort - apply it.
		sortConfig := i.sortModal.ToSortConfig()

		// Update internal filter state.
		if i.context.Filters == nil {
			i.context.Filters = &Filters{}
		}
		i.context.Filters.SortBy = sortConfig.SortBy
		i.context.Filters.SortOrder = sortConfig.SortOrder

		// Reload skills with new sort.
		return i.reloadSkills()
	}

	// Modal was closed without completion (Esc) or still being edited.
	return cmd
}

// handleSearchModalUpdate handles updates when search modal is visible.
// Takes tea.Msg (not tea.KeyMsg) to allow huh forms to work correctly.
func (i *Intent) handleSearchModalUpdate(msg tea.Msg) tea.Cmd {
	cmd, applied, searchData := i.searchModal.Update(msg)

	if applied && searchData != nil {
		// User confirmed search - apply it.
		if i.context.Filters == nil {
			i.context.Filters = &Filters{}
		}
		i.context.Filters.SearchText = searchData.SearchText

		// Reload skills with new search.
		return i.reloadSkills()
	}

	// Modal was closed without completion (Esc) or still being edited.
	return cmd
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
			// User completed form - save skill.
			originalSkill := i.addEditModal.GetOriginalSkill()
			if originalSkill != nil {
				// Editing existing skill.
				skill := skillData.ToSkill(originalSkill.ID)
				i.addEditModal = nil
				return i.updateSkill(skill)
			}
			// Creating new skill.
			skill := skillData.ToSkill("")
			i.addEditModal = nil
			return i.createSkill(skill)
		}
		// User cancelled - close modal.
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
				err := i.context.SkillRepository.Delete(i.context.Ctx, skillID)
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
			// Open event detail modal.
			return i.openEventDetailModal(selectedEvent)
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
		err := i.context.SkillRepository.Create(i.context.Ctx, skill)
		return SkillCreatedMsg{
			Skill: skill,
			Error: err,
		}
	}
}

// updateSkill updates an existing skill.
func (i *Intent) updateSkill(skill *domain.Skill) tea.Cmd {
	return func() tea.Msg {
		err := i.context.SkillRepository.Update(i.context.Ctx, skill)
		return SkillUpdatedMsg{
			Skill: skill,
			Error: err,
		}
	}
}

// handleSkillSuggestionsLoaded handles the SkillSuggestionsLoadedMsg.
func (i *Intent) handleSkillSuggestionsLoaded(msg SkillSuggestionsLoadedMsg) tea.Cmd {
	i.loadingModal = nil

	if msg.Error != nil {
		// Silently ignore cancelled operations.
		if errors.Is(msg.Error, context.Canceled) {
			i.state = StateList
			return nil
		}
		i.feedbackModal = feedback.NewErrorModal("Skill Inference Failed", msg.Error.Error())
		i.state = StateList
		return nil
	}

	if len(msg.Suggestions) == 0 {
		i.feedbackModal = feedback.NewWarningModal("No Skills Found", "No skills were detected from the events")
		i.state = StateList
		return nil
	}

	i.skillSuggestionModal = burstmodals.NewSkillSuggestionModal(msg.Suggestions, i.Theme())
	width, height := i.getTerminalDimensions()
	i.skillSuggestionModal.SetDimensions(width, height)
	i.skillSuggestionModal.Show()
	i.state = StateSkillSuggestionReview
	return nil
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
		case burstmodals.SuggestionActionAccept, burstmodals.SuggestionActionReject:
			if !i.skillSuggestionModal.IsVisible() {
				accepted := i.skillSuggestionModal.GetAcceptedSkills()

				if len(accepted) > 0 {
					i.loadingModal = feedback.NewLoadingModal(
						fmt.Sprintf("Creating %d skill(s)...", len(accepted)),
						true,
					).WithTheme(i.Theme())
					i.skillSuggestionModal = nil
					return tea.Batch(cmd, i.createSkillsFromSuggestions(accepted))
				}

				i.skillSuggestionModal = nil
				i.state = StateList
			}
			return noopCmd

		case burstmodals.SuggestionActionViewEvents:
			return i.openSuggestionEventsModal()

		case burstmodals.SuggestionActionCancel:
			i.skillSuggestionModal = nil
			i.state = StateList
			return noopCmd
		}
	}

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
			return SkillsCreatedMsg{Error: fmt.Errorf("skill inference service not available")}
		}

		skills, err := service.CreateSkillsFromSuggestions(i.context.Ctx, suggestions)
		if err != nil {
			return SkillsCreatedMsg{Error: fmt.Errorf("failed to create skills: %w", err)}
		}

		return SkillsCreatedMsg{
			Skills: skills,
			Error:  nil,
		}
	}
}

// handleSkillsCreatedFromInference handles skills created from accepted suggestions.
func (i *Intent) handleSkillsCreatedFromInference(msg SkillsCreatedMsg) tea.Cmd {
	if msg.Error != nil {
		i.feedbackModal = feedback.NewErrorModal("Skill Creation Failed", msg.Error.Error())
		i.state = StateList
		return nil
	}

	// Skills created successfully - refresh the list to show them
	i.state = StateList
	return i.RefreshData()
}
