package browsetimeline

import (
	"fmt"

	"github.com/baphled/kariya/internal/domain/capture"
	"github.com/baphled/kariya/internal/domain/career"
	eventviews "github.com/baphled/kariya/internal/tui/views/event"
	"github.com/baphled/kariya/internal/tui/views/shared"
	"github.com/baphled/kariya/internal/ui/behaviors"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/uikit/feedback"
	tea "github.com/charmbracelet/bubbletea"
)

const editFailedTitle = "Edit Failed"

// openQuickAddModal creates and shows the quick add modal.
func (i *Intent) openQuickAddModal() tea.Cmd {
	width, height := i.getTerminalDimensions()
	i.quickAdd = eventviews.NewQuickAdd(width, height)
	return i.quickAdd.Init()
}

// openEditModalForEvent creates and shows the edit modal for an event.
func (i *Intent) openEditModalForEvent(evt *career.Event) tea.Cmd {
	width, height := i.getTerminalDimensions()
	i.edit = eventviews.NewEdit(display.EventFromDomain(evt), width, height)
	return i.edit.Init()
}

// openDeleteModalForEvent creates and shows the delete confirmation modal.
func (i *Intent) openDeleteModalForEvent(evt *career.Event) tea.Cmd {
	eventText := evt.Text
	if len(eventText) > 50 {
		eventText = eventText[:47] + "..."
	}
	i.deleteModal = feedback.NewConfirmModal(
		"Delete Event",
		fmt.Sprintf("Are you sure you want to delete '%s'?", eventText),
	).WithVariant(feedback.ConfirmDestructive)
	i.selectedEvent = evt
	return i.deleteModal.Init()
}

// showEventDetailModal creates and shows the event detail modal.
func (i *Intent) showEventDetailModal(evt *career.Event) tea.Cmd {
	i.selectedEvent = evt
	i.viewedEvents = append(i.viewedEvents, evt)
	width, height := i.getTerminalDimensions()
	theme := i.Theme()
	i.viewDetail = eventviews.NewDetail(display.EventFromDomain(evt), theme)
	i.viewDetail.SetDimensions(width, height)
	i.viewDetail.Show()
	return nil
}

// updateSearchModal handles search adapter updates via FormViewAdapter.
func (i *Intent) updateSearchModal(msg tea.Msg) tea.Cmd {
	result := i.searchAdapter.HandleUpdate(msg)
	if result.Applied {
		if searchData, ok := result.Data.(shared.SearchFormData); ok {
			i.filters.SearchText = searchData.SearchText
			if searchData.SearchText != "" {
				i.filterStack.Push(behaviors.FilterLayerSearch)
			}
			i.RefreshData()
		}
	}
	if result.Closed && !result.Applied {
		i.searchAdapter = nil
	}
	return result.Cmd
}

// updateFilterModal handles filter adapter updates via FormViewAdapter.
func (i *Intent) updateFilterModal(msg tea.Msg) tea.Cmd {
	result := i.filterAdapter.HandleUpdate(msg)
	if result.Applied {
		i.applyFilterFormData(result.Data)
	}
	if result.Closed && !result.Applied {
		i.filterAdapter = nil
	}
	return result.Cmd
}

// updateSortModal handles sort adapter updates via FormViewAdapter.
func (i *Intent) updateSortModal(msg tea.Msg) tea.Cmd {
	result := i.sortAdapter.HandleUpdate(msg)
	if result.Applied {
		if sortData, ok := result.Data.(shared.SortFormData); ok {
			i.filters.SortBy = sortData.SortBy
			i.filters.SortOrder = sortData.SortOrder
			if sortData.SortBy != "date" || sortData.SortOrder != "desc" {
				i.filterStack.Push(behaviors.FilterLayerSort)
			}
			i.RefreshData()
		}
	}
	if result.Closed && !result.Applied {
		i.sortAdapter = nil
	}
	return result.Cmd
}

// updateQuickAddModal handles quick add modal updates.
func (i *Intent) updateQuickAddModal(msg tea.Msg) tea.Cmd {
	cmd, completed, eventData := i.quickAdd.Update(msg)
	if !i.quickAdd.IsVisible() {
		if completed && eventData != nil {
			ctx := i.getContext()
			newEvent, err := capture.NewEventFromInput(capture.EventInput{
				Text: eventData.Text,
				Date: eventData.Date,
			})
			if err != nil {
				i.ShowErrorModal("Quick Add Failed", err.Error())
				return cmd
			}
			if err := i.context.CLIEventService.CaptureEvent(ctx, newEvent.Text, newEvent.Date, "manual"); err != nil {
				i.ShowErrorModal("Quick Add Failed", err.Error())
				return cmd
			}
			refreshedEvents, listErr := i.context.CLIEventService.ListEvents(ctx, nil)
			if listErr != nil {
				i.ShowErrorModal("Refresh Failed", listErr.Error())
				return cmd
			}
			i.context.Events = refreshedEvents
			i.RefreshData()
		}
		i.quickAdd = nil
	}
	return cmd
}

// updateEditModal handles edit modal updates.
func (i *Intent) updateEditModal(msg tea.Msg) tea.Cmd {
	cmd, completed, eventData := i.edit.Update(msg)
	if !i.edit.IsVisible() {
		if completed && eventData != nil {
			i.applyEditModalChanges(eventData)
		}
		i.edit = nil
	}
	return cmd
}

// applyFilterFormData applies filter form selections to the intent.
//
// Expected:
//   - data contains filter form data of type eventviews.FilterFormData.
//
// Returns:
//   - None.
//
// Side effects:
//   - Updates filter state and refreshes data when applied.
func (i *Intent) applyFilterFormData(data interface{}) {
	filterData, ok := data.(eventviews.FilterFormData)
	if !ok {
		return
	}
	i.applyFilterSelections(filterData)
	i.filters.SortBy = filterData.SortBy
	i.filters.SortOrder = filterData.SortOrder
	i.pushActiveFilterLayers()
	i.RefreshData()
}

// applyFilterSelections sets company/category/project selections from form data.
//
// Expected:
//   - data contains selection values keyed by filter category.
//
// Returns:
//   - None.
//
// Side effects:
//   - Updates filter fields on the intent.
func (i *Intent) applyFilterSelections(data eventviews.FilterFormData) {
	if companies, exists := data.Selections["companies"]; exists {
		i.filters.Companies = companies
	}
	if categories, exists := data.Selections["categories"]; exists {
		i.filters.Categories = categories
	}
	if projects, exists := data.Selections["projects"]; exists {
		i.filters.Projects = projects
	}
}

// pushActiveFilterLayers pushes non-empty filter layers onto the stack.
//
// Expected:
//   - Filters are already populated on the intent.
//
// Returns:
//   - None.
//
// Side effects:
//   - Updates the filter stack state.
func (i *Intent) pushActiveFilterLayers() {
	if len(i.filters.Companies) > 0 {
		i.filterStack.Push(behaviors.FilterLayerCompany)
	}
	if len(i.filters.Categories) > 0 {
		i.filterStack.Push(behaviors.FilterLayerCategory)
	}
	if len(i.filters.Projects) > 0 {
		i.filterStack.Push(behaviors.FilterLayerProject)
	}
}

// applyEditModalChanges updates the selected event using edit modal data.
//
// Expected:
//   - eventData contains edited event fields.
//   - cmd is the command returned by the modal update.
//
// Returns:
//   - None.
//
// Side effects:
//   - Updates event data, persists changes, and refreshes list state.
func (i *Intent) applyEditModalChanges(eventData *eventviews.EditData) {
	ctx := i.getContext()
	originalEvent := i.findEventByID(i.edit.GetOriginalEvent().ID)
	if originalEvent == nil {
		i.ShowErrorModal(editFailedTitle, "Original event could not be resolved.")
		return
	}
	updatedEvent, err := capture.UpdateEventFromInput(capture.EditEventInput{
		EventID:    originalEvent.ID,
		Text:       eventData.Text,
		Date:       eventData.Date,
		Company:    eventData.Company,
		Project:    eventData.Project,
		Tags:       eventData.Tags,
		Categories: eventData.Categories,
		CreatedAt:  originalEvent.CreatedAt,
		UpdatedAt:  originalEvent.UpdatedAt,
	})
	if err != nil {
		i.ShowErrorModal(editFailedTitle, err.Error())
		return
	}
	if err := i.context.CLIEventService.UpdateEventMetadata(ctx, updatedEvent); err != nil {
		i.ShowErrorModal(editFailedTitle, err.Error())
		return
	}
	for idx, evt := range i.context.Events {
		if evt.ID == updatedEvent.ID {
			i.context.Events[idx] = updatedEvent
			break
		}
	}
	i.RefreshData()
}

// handleSkillMessages routes skill-related messages for the intent.
//
// Expected:
//   - msg is a Bubble Tea message.
//
// Returns:
//   - A tea.Cmd and a boolean indicating whether the message was handled.
//
// Side effects:
//   - May update skill modal state or show error modals.
func (i *Intent) handleSkillMessages(msg tea.Msg) (tea.Cmd, bool) {
	switch msg := msg.(type) {
	case SkillLinkedMsg:
		return i.handleSkillLinkResult("Error Linking Skill", msg.Error)
	case SkillUnlinkedMsg:
		return i.handleSkillLinkResult("Error Unlinking Skill", msg.Error)
	case SkillCreatedMsg:
		return i.handleSkillLinkResult("Error Creating Skill", msg.Error)
	case SkillSuggestionsLoadedMsg:
		i.handleSkillSuggestionsLoaded(msg)
		return nil, true
	case SkillSuggestionsErrorMsg:
		i.ShowErrorModal("Skill Inference Failed", msg.Error.Error())
		return nil, true
	case SkillsForModalLoadedMsg:
		i.handleSkillsForModalLoaded(msg)
		return nil, true
	case SkillPickerDataLoadedMsg:
		i.handleSkillPickerDataLoaded(msg)
		return nil, true
	case SkillsRefreshedMsg:
		i.handleSkillsRefreshed(msg)
		return nil, true
	default:
		return nil, false
	}
}

// handleSkillLinkResult handles skill link/create responses by showing errors
// or refreshing the skills modal.
//
// Expected:
//   - errorTitle describes the user-facing error title.
//   - err is the message error (may be nil).
//
// Returns:
//   - A tea.Cmd and a boolean indicating the message was handled.
//
// Side effects:
//   - May show error modal or refresh skills modal data.
func (i *Intent) handleSkillLinkResult(errorTitle string, err error) (tea.Cmd, bool) {
	if err != nil {
		i.ShowErrorModal(errorTitle, err.Error())
		return nil, true
	}
	return i.refreshSkillsModal(), true
}

// routeModalAndKeys processes modal updates and keyboard shortcuts in order.
//
// Expected:
//   - msg is a Bubble Tea message.
//
// Returns:
//   - A tea.Cmd and a boolean indicating whether routing consumed the message.
//
// Side effects:
//   - May update modal state, open modals, or trigger filters.
func (i *Intent) routeModalAndKeys(msg tea.Msg) (tea.Cmd, bool) {
	if cmd := i.handleModalUpdates(msg); cmd != nil || i.hasActiveModal() {
		return cmd, true
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if cmd := i.handleKeyShortcuts(keyMsg); cmd != nil {
			return cmd, true
		}
	}

	return nil, false
}

// handleCancelViewResult processes cancel results from the active view.
//
// Expected:
//   - cmd is the command returned by the view update.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - May hide modals or mark the intent as cancelled.
func (i *Intent) handleCancelViewResult(cmd tea.Cmd) tea.Cmd {
	if i.hasActiveModal() {
		if i.viewDetail != nil && i.viewDetail.IsVisible() {
			i.viewDetail.Hide()
			i.viewDetail = nil
			return nil
		}
		return cmd
	}

	i.setCancelled()
	return cmd
}

// noopCmd is a sentinel command to indicate a message was consumed.
// This prevents the message from propagating to the screen after a modal closes.
func noopCmd() tea.Msg { return nil }

// modalHandler defines the interface for handling a modal's update cycle.
type modalHandler struct {
	isActive func() bool
	update   func(tea.Msg) tea.Cmd
	isClosed func() bool
}

// handleModalUpdates handles updates for all modals in priority order.
// Returns a command (possibly noopCmd) if a modal consumed the message.
func (i *Intent) handleModalUpdates(msg tea.Msg) tea.Cmd {
	if cmd := i.handleErrorModalUpdate(msg); cmd != nil {
		return cmd
	}

	handlers := i.buildModalHandlers()

	return i.processModalHandlers(msg, handlers)
}

// handleErrorModalUpdate handles error modal updates with highest priority.
//
// Expected:
//   - msg is a Bubble Tea message.
//
// Returns:
//   - A tea.Cmd value when the error modal is active; nil otherwise.
//
// Side effects:
//   - May close the error modal on escape.
func (i *Intent) handleErrorModalUpdate(msg tea.Msg) tea.Cmd {
	if i.errorModal == nil {
		return nil
	}
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.Type == tea.KeyEsc {
		i.errorModal = nil
		return noopCmd
	}
	return noopCmd
}

// processModalHandlers updates the first active modal in the list.
//
// Expected:
//   - msg is a Bubble Tea message.
//   - handlers are ordered by priority.
//
// Returns:
//   - A tea.Cmd for the active modal, or nil if none are active.
//
// Side effects:
//   - May return noopCmd when a modal closes without a command.
func (i *Intent) processModalHandlers(msg tea.Msg, handlers []modalHandler) tea.Cmd {
	for _, h := range handlers {
		if h.isActive() {
			cmd := h.update(msg)
			if h.isClosed() && cmd == nil {
				return noopCmd
			}
			return cmd
		}
	}

	return nil
}

// hasActiveModal returns true if any modal is currently visible.
func (i *Intent) hasActiveModal() bool {
	i.rebuildModalRegistry()
	return i.modalRegistry.HasVisibleModal()
}

// updateDeleteModal handles delete modal updates.
func (i *Intent) updateDeleteModal(msg tea.Msg) tea.Cmd {
	cmd, confirmed := i.deleteModal.Update(msg)
	if !i.deleteModal.IsVisible() {
		if confirmed && i.selectedEvent != nil {
			ctx := i.getContext()
			if err := i.context.CLIEventService.DeleteEvent(ctx, i.selectedEvent.ID); err != nil {
				i.ShowErrorModal("Delete Failed", err.Error())
				i.deleteError = err
				i.deleteModal = nil
				return cmd
			}
			deletedID := i.selectedEvent.ID
			newEvents := make([]*career.Event, 0, len(i.context.Events)-1)
			for _, evt := range i.context.Events {
				if evt.ID != deletedID {
					newEvents = append(newEvents, evt)
				}
			}
			i.context.Events = newEvents
			i.RefreshData()
		}
		i.selectedEvent = nil
		i.deleteModal = nil
	}
	return cmd
}

// updateViewSkillsModal handles view skills modal updates.
func (i *Intent) updateViewSkillsModal(msg tea.Msg) tea.Cmd {
	_, cmd := i.viewSkills.Update(msg)

	if !i.viewSkills.IsVisible() {
		i.viewSkills = nil
		return cmd
	}

	action := i.viewSkills.GetAction()
	if action != eventviews.ActionNone {
		i.viewSkills.ClearAction()

		switch action {
		case eventviews.ActionAddExisting:
			return i.openSkillPickerModal()

		case eventviews.ActionAddNew:
			return i.openSkillAddModal()

		case eventviews.ActionInfer:
			return i.inferSkillsFromEvent()

		case eventviews.ActionRemove:
			selectedSkill := i.viewSkills.GetSelectedSkill()
			return i.unlinkSkillFromCurrentEvent(i.findSkillByID(selectedSkill.ID, i.selectedEventSkills))
		}
	}

	return cmd
}

// updateSkillPickerModal handles skill picker modal updates.
func (i *Intent) updateSkillPickerModal(msg tea.Msg) tea.Cmd {
	_, cmd := i.skillPicker.Update(msg)

	if !i.skillPicker.IsVisible() {
		if i.skillPicker.HasSelection() {
			selectedSkill := i.skillPicker.GetSelectedSkill()
			i.skillPicker = nil
			return tea.Batch(cmd, i.linkSkillToCurrentEvent(i.findSkillByID(selectedSkill.ID, i.availableSkills)))
		}
		i.skillPicker = nil
	}

	return cmd
}

// updateSkillAddModal handles skill add modal updates.
func (i *Intent) updateSkillAddModal(msg tea.Msg) tea.Cmd {
	cmd, completed, skillData := i.skillAddModal.Update(msg)
	if !i.skillAddModal.IsVisible() {
		if completed && skillData != nil {
			i.skillAddModal = nil
			return i.createAndLinkSkill(skillData)
		}
		i.skillAddModal = nil
	}
	return cmd
}

// updateSkillSuggestionModal handles skill suggestion modal updates.
func (i *Intent) updateSkillSuggestionModal(msg tea.Msg) tea.Cmd {
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "a" {
		if selected := i.skillSuggestionModal.GetCurrentSkill(); selected != nil {
			domainSuggestion := i.findSkillSuggestionByName(selected.Name)
			if domainSuggestion != nil {
				i.saveSkillFromSuggestion(*domainSuggestion)
			}
		}
	}
	_, cmd := i.skillSuggestionModal.Update(msg)
	if !i.skillSuggestionModal.IsVisible() {
		accepted := i.skillSuggestionModal.GetAcceptedSkills()
		if len(accepted) > 0 {
			i.ShowErrorModal("Skills Created", fmt.Sprintf("Successfully created %d skill(s)", len(accepted)))
		}
		i.skillSuggestionModal = nil
		if i.viewDetail != nil && i.selectedEvent != nil {
			i.viewDetail.SetEvent(display.EventFromDomain(i.selectedEvent))
		}
		return tea.Batch(cmd, i.refreshSkillsModal())
	}
	return cmd
}

// updateViewDetailModal handles view detail modal updates.
func (i *Intent) updateViewDetailModal(msg tea.Msg) tea.Cmd {
	if cmd, handled := i.handleViewDetailKeyMsg(msg); handled {
		return cmd
	}
	_, cmd := i.viewDetail.Update(msg)
	if !i.viewDetail.IsVisible() {
		i.viewDetail = nil
	}
	return cmd
}

// handleViewDetailKeyMsg routes keyboard actions for the detail modal.
//
// Expected:
//   - msg is a Bubble Tea message.
//
// Returns:
//   - A tea.Cmd and a bool indicating whether the key was handled.
//
// Side effects:
//   - May hide the detail modal or open edit/delete/skills modals.
func (i *Intent) handleViewDetailKeyMsg(msg tea.Msg) (tea.Cmd, bool) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return nil, false
	}
	if cmd, handled := i.handleViewDetailActionKey(keyMsg); handled {
		return cmd, true
	}
	return i.handleViewDetailEscapeKey(keyMsg)
}

// handleViewDetailActionKey handles action keys for the detail modal.
//
// Expected:
//   - keyMsg is a Bubble Tea key message.
//
// Returns:
//   - A tea.Cmd and a bool indicating whether the key was handled.
//
// Side effects:
//   - May open related modals or update modal state.
func (i *Intent) handleViewDetailActionKey(keyMsg tea.KeyMsg) (tea.Cmd, bool) {
	switch keyMsg.String() {
	case "s":
		return i.openViewSkillsModal(), true
	case "e":
		if i.selectedEvent == nil {
			return nil, true
		}
		i.viewDetail.Hide()
		i.viewDetail = nil
		return i.openEditModalForEvent(i.selectedEvent), true
	case "d":
		if i.selectedEvent == nil {
			return nil, true
		}
		return i.openDeleteModalFromDetail(), true
	default:
		return nil, false
	}
}

// handleViewDetailEscapeKey handles escape key behavior for the detail modal.
//
// Expected:
//   - keyMsg is a Bubble Tea key message.
//
// Returns:
//   - A tea.Cmd and a bool indicating whether the key was handled.
//
// Side effects:
//   - May hide the detail modal and return a noop command.
func (i *Intent) handleViewDetailEscapeKey(keyMsg tea.KeyMsg) (tea.Cmd, bool) {
	if keyMsg.Type != tea.KeyEsc && keyMsg.String() != "esc" {
		return nil, false
	}
	i.viewDetail.Hide()
	if !i.viewDetail.IsVisible() {
		i.viewDetail = nil
		return noopCmd, true
	}
	return nil, true
}

// openDeleteModalFromDetail opens the delete modal for the selected event.
//
// Expected:
//   - selectedEvent is set on the intent.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - Hides the detail modal and shows the delete confirmation modal.
func (i *Intent) openDeleteModalFromDetail() tea.Cmd {
	eventText := i.selectedEvent.Text
	if len(eventText) > 50 {
		eventText = eventText[:47] + "..."
	}
	i.viewDetail.Hide()
	i.viewDetail = nil
	i.deleteModal = feedback.NewConfirmModal(
		"Delete Event",
		fmt.Sprintf("Are you sure you want to delete '%s'?", eventText),
	).WithVariant(feedback.ConfirmDestructive)
	return i.deleteModal.Init()
}

// handleKeyShortcuts handles keyboard shortcuts when no modal is active.
func (i *Intent) handleKeyShortcuts(keyMsg tea.KeyMsg) tea.Cmd {
	switch keyMsg.String() {
	case "/":
		return i.openSearchModal()
	case "f":
		return i.openFilterModal()
	case "s":
		return i.openSortModal()
	case "x":
		if i.HasActiveFilters() {
			i.ClearFilters()
			i.RefreshData()
			return nil
		}
	}
	return nil
}
