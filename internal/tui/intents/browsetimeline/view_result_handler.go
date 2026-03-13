package browsetimeline

import (
	"fmt"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/tui/views/event"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
	tea "github.com/charmbracelet/bubbletea"
)

// handleViewResult processes a view result with typed action dispatch (ADR pattern).
// Extracts Nav payload and switches on action type.
func (i *Intent) handleViewResult(res widgets.ViewResult) tea.Cmd {
	if res.Type() != widgets.ResultNavigate {
		return nil
	}

	if skillNav, ok := res.Data().(event.SkillNav); ok {
		return i.handleSkillNavAction(skillNav)
	}

	payload, ok := res.Data().(event.Nav)
	if !ok {
		return func() tea.Msg {
			return fmt.Errorf("invalid event action payload type: %T", res.Data())
		}
	}

	switch payload.Action {
	case event.ActionEdit:
		return i.handleEventAction(payload.Event.ID, i.openEditModalForEvent)
	case event.ActionDelete:
		return i.handleEventAction(payload.Event.ID, i.openDeleteModalForEvent)
	case event.ActionView:
		return i.handleEventAction(payload.Event.ID, i.showEventDetailModal)
	case event.ActionAdd:
		return i.openQuickAddModal()
	case event.ActionFilter:
		return i.openFilterModal()
	case event.ActionSearch:
		return nil
	case event.ActionSort:
		return nil
	case event.ActionClear:
		return nil
	default:
		return func() tea.Msg {
			return fmt.Errorf("unknown event action: %s", payload.Action)
		}
	}
}

// handleSkillNavAction routes skill navigation actions.
//
// Expected:
//   - nav contains the skill action and event context.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - Updates selected event and opens skill-related modals.
func (i *Intent) handleSkillNavAction(nav event.SkillNav) tea.Cmd {
	if nav.Event.ID != "" {
		i.selectedEvent = i.findEventByID(nav.Event.ID)
	}
	switch nav.Action {
	case event.ActionShowEventSkills:
		return i.openViewSkillsModal()
	case event.ActionInferSkills:
		return i.inferSkillsFromEvent()
	default:
		return nil
	}
}

// handleEventAction resolves an event and executes the given action.
//
// Expected:
//   - eventID is the target event identifier.
//   - action is a handler for the resolved event.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - May open modals or update intent state.
func (i *Intent) handleEventAction(eventID string, action func(*career.Event) tea.Cmd) tea.Cmd {
	if eventID == "" {
		return nil
	}
	evt := i.findEventByID(eventID)
	if evt == nil {
		return nil
	}
	return action(evt)
}
