package timeline

import (
	"fmt"

	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/domain/career"
)

// EventDeleteConfirmState identifies the event deletion confirmation view
// in the state matrix. On this screen the user sees the event summary and
// a yes/no prompt to confirm permanent deletion. The event text is
// truncated to 60 characters inside a confirmation dialog with Delete and
// Cancel buttons. A warning states the action cannot be undone. The Cancel
// button is focused by default for safety. Left/right arrows or h/l
// toggle selection, y/n submit directly, and Escape cancels.
const EventDeleteConfirmState = "event_delete_confirm"

// EventDeleteConfirmScreen provides a confirmation dialog for deleting a career event.
//
// This screen wraps ConfirmScreen with event-specific context:
// - Shows event text (truncated) in confirmation message
// - Customizes button text ("Delete" / "Cancel")
// - Preserves event data for the caller
//
// Usage:
//
//	screen := browsetimeline.NewEventDeleteConfirmScreen(event)
//	cmd, result := screen.Update(msg)
//	if result != nil && result.Type() == screens.ResultNavigate {
//	    confirmed := result.Data().(bool)
//	    if confirmed {
//	        event := screen.GetEvent()
//	        // Proceed with deletion
//	    } else {
//	        // User cancelled - go back
//	    }
//	}
//
// Related:
// - ConfirmScreen provides the confirmation UI
// - docs/TUI_DEVELOPER_GUIDE.md (Screen patterns)
// - docs/TUI_STANDARDS.md (Keyboard shortcuts).
type EventDeleteConfirmScreen struct {
	*base.ConfirmScreen
	event *career.Event
}

// NewEventDeleteConfirmScreen creates a new event deletion confirmation screen.
//
// Expected:
//   - event must be valid.
//
// Returns:
//   - A fully initialized EventDeleteConfirmScreen ready for use.
//
// Side effects:
//   - None.
func NewEventDeleteConfirmScreen(event *career.Event) *EventDeleteConfirmScreen {
	breadcrumbs := []string{"Main Menu", "Timeline", "Delete Confirmation"}
	title := "Delete Event"

	eventText := event.Text
	if len(eventText) > 60 {
		eventText = eventText[:60] + "..."
	}

	message := fmt.Sprintf(
		"Are you sure you want to delete this event?\n\n%q\n\nThis action cannot be undone.",
		eventText,
	)

	confirmScreen := base.NewBaseConfirmScreen(breadcrumbs, title, message)
	confirmScreen.SetYesText("Delete")
	confirmScreen.SetNoText("Cancel")

	return &EventDeleteConfirmScreen{
		ConfirmScreen: confirmScreen,
		event:         event,
	}
}

// GetEvent returns the event being considered for deletion.
//
// Returns:
//   - A fully initialized career.Event ready for use.
//
// Side effects:
//   - None.
func (s *EventDeleteConfirmScreen) GetEvent() *career.Event {
	return s.event
}
