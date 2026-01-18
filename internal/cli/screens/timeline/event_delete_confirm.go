package timeline

import (
	"fmt"

	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/domain/career"
)

// State constant for state matrix tracking (REQUIRED)
const EventDeleteConfirmState = "event_delete_confirm"

// EventDeleteConfirmScreen provides a confirmation dialog for deleting a career event.
//
// This screen wraps BaseConfirmScreen with event-specific context:
// - Shows event text (truncated) in confirmation message
// - Customizes button text ("Delete" / "Cancel")
// - Preserves event data for the caller
//
// Usage:
//
//	screen := timeline.NewEventDeleteConfirmScreen(event)
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
// - BaseConfirmScreen provides the confirmation UI
// - docs/TUI_DEVELOPER_GUIDE.md (Screen patterns)
// - docs/TUI_STANDARDS.md (Keyboard shortcuts)
type EventDeleteConfirmScreen struct {
	*base.BaseConfirmScreen
	event *career.CareerEvent
}

// NewEventDeleteConfirmScreen creates a new event deletion confirmation screen.
//
// The screen:
// - Displays event text (truncated to 60 chars) in the confirmation message
// - Uses "Delete" / "Cancel" button text
// - Defaults to "No" (Cancel) for safety
// - Supports standard navigation: ←→/hl to toggle, y/n for direct, Enter to confirm
//
// Parameters:
//   - event: The event to delete
func NewEventDeleteConfirmScreen(event *career.CareerEvent) *EventDeleteConfirmScreen {
	breadcrumbs := []string{"Main Menu", "Timeline", "Delete Confirmation"}
	title := "Delete Event"

	// Truncate event text for display
	eventText := event.Text
	if len(eventText) > 60 {
		eventText = eventText[:60] + "..."
	}

	message := fmt.Sprintf(
		"Are you sure you want to delete this event?\n\n\"%s\"\n\nThis action cannot be undone.",
		eventText,
	)

	confirmScreen := base.NewBaseConfirmScreen(breadcrumbs, title, message)
	confirmScreen.SetYesText("Delete")
	confirmScreen.SetNoText("Cancel")

	return &EventDeleteConfirmScreen{
		BaseConfirmScreen: confirmScreen,
		event:             event,
	}
}

// GetEvent returns the event being considered for deletion.
func (s *EventDeleteConfirmScreen) GetEvent() *career.CareerEvent {
	return s.event
}
