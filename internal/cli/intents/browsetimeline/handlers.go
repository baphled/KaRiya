package browsetimeline

import (
	"fmt"

	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/timeline"
	"github.com/baphled/kariya/internal/cli/screens/timeline/modals"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

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

// HandleCancel handles screen cancellation (back/escape).
func (i *Intent) HandleCancel(result *screens.CancelResult) tea.Cmd {
	switch i.state {
	case StateTimeline:
		i.setCancelled()
		return nil

	case StateDeleteConfirm:
		i.state = StateTimeline
		i.transitionToScreen(timeline.NewTimelineEventListScreen(i.filteredEvents))
		return nil

	default:
		i.setCancelled()
		return nil
	}
}

// HandleNavigate handles screen navigation results.
func (i *Intent) HandleNavigate(result *screens.NavigateResult) tea.Cmd {
	if actionData, ok := result.ResultData.(map[string]interface{}); ok {
		return i.handleActionData(actionData)
	}

	if confirmed, ok := result.ResultData.(bool); ok {
		return i.handleDeleteConfirmation(confirmed)
	}

	if event, ok := result.ResultData.(*career.Event); ok {
		return i.showEventDetailModal(event)
	}

	return nil
}

// handleActionData processes action data from navigation results.
func (i *Intent) handleActionData(actionData map[string]interface{}) tea.Cmd {
	action, ok := actionData["action"].(string)
	if !ok {
		return nil
	}
	switch action {
	case "add":
		return i.openQuickAddModal()

	case "edit":
		if event, ok := actionData["event"].(*career.Event); ok {
			return i.openEditModalForEvent(event)
		}
		return nil

	case "delete":
		if event, ok := actionData["event"].(*career.Event); ok {
			return i.openDeleteModalForEvent(event)
		}
		return nil

	case "filter":
		return i.openFilterModal()

	default:
		return nil
	}
}

// openQuickAddModal creates and shows the quick add modal.
func (i *Intent) openQuickAddModal() tea.Cmd {
	width, height := i.getTerminalDimensions()
	i.quickAddModal = modals.NewQuickAddModal(width, height)
	return i.quickAddModal.Init()
}

// openEditModalForEvent creates and shows the edit modal for an event.
func (i *Intent) openEditModalForEvent(event *career.Event) tea.Cmd {
	width, height := i.getTerminalDimensions()
	i.editModal = modals.NewEditModal(event, width, height)
	return i.editModal.Init()
}

// openDeleteModalForEvent creates and shows the delete confirmation modal.
func (i *Intent) openDeleteModalForEvent(event *career.Event) tea.Cmd {
	eventText := event.Text
	if len(eventText) > 50 {
		eventText = eventText[:47] + "..."
	}
	i.deleteModal = feedback.NewConfirmModal(
		"Delete Event",
		fmt.Sprintf("Are you sure you want to delete '%s'?", eventText),
	).WithVariant(feedback.ConfirmDestructive)
	i.selectedEvent = event
	return i.deleteModal.Init()
}

// showEventDetailModal creates and shows the event detail modal.
func (i *Intent) showEventDetailModal(event *career.Event) tea.Cmd {
	i.selectedEvent = event
	i.viewedEvents = append(i.viewedEvents, event)
	width, height := i.getTerminalDimensions()
	theme := i.Theme()
	i.viewDetailModal = modals.NewEventDetailModal(event, theme)
	i.viewDetailModal.SetDimensions(width, height)
	i.viewDetailModal.Show()
	return nil
}

// handleDeleteConfirmation handles the result of the delete confirmation dialog.
func (i *Intent) handleDeleteConfirmation(confirmed bool) tea.Cmd {
	if !confirmed {
		i.state = StateTimeline
		i.transitionToScreen(timeline.NewTimelineEventListScreen(i.filteredEvents))
		return nil
	}

	if i.selectedEvent == nil {
		i.state = StateTimeline
		i.transitionToScreen(timeline.NewTimelineEventListScreen(i.filteredEvents))
		return nil
	}

	if i.context.CLIEventService != nil {
		eventID := i.selectedEvent.ID
		err := i.context.CLIEventService.DeleteEvent(i.getContext(), eventID)
		if err != nil {
			i.deleteError = err
			i.state = StateTimeline
			i.transitionToScreen(timeline.NewTimelineEventListScreen(i.filteredEvents))
			return nil
		}

		i.removeEventFromList(eventID)
	}

	i.state = StateTimeline
	i.transitionToScreen(timeline.NewTimelineEventListScreen(i.filteredEvents))
	return nil
}

// HandleSubmit handles form submission results.
func (i *Intent) HandleSubmit(result *screens.SubmitResult) tea.Cmd {
	return nil
}

// HandleError handles error results from screens.
func (i *Intent) HandleError(result *screens.ErrorResult) tea.Cmd {
	i.deleteError = result.Err
	return nil
}
