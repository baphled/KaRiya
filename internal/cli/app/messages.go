package app

import (
	"github.com/baphled/kariya/internal/domain/career"
)

// FormSubmittedMsg is sent when a form is successfully submitted
type FormSubmittedMsg struct {
	Event *career.CareerEvent
	Err   error
}

// NavigateMsg is sent to navigate to a specific screen
type NavigateMsg struct {
	Screen Screen
}

// SuccessNavigateMsg is sent from success screen for navigation
type SuccessNavigateMsg struct {
	Action string // "capture_another", "view_list", "exit"
}

// EditEventMsg is sent to edit a specific event
type EditEventMsg struct {
	Event *career.CareerEvent
}

// DeleteEventMsg is sent to delete a specific event
type DeleteEventMsg struct {
	EventID string
}

// ConfirmDeleteMsg is sent when delete is confirmed
type ConfirmDeleteMsg struct {
	EventID string
}

// CancelDeleteMsg is sent when delete is cancelled
type CancelDeleteMsg struct{}

// EventDeletedMsg is sent when an event is successfully deleted
type EventDeletedMsg struct {
	EventID string
	Err     error
}

// EventUpdatedMsg is sent when an event is successfully updated
type EventUpdatedMsg struct {
	Event *career.CareerEvent
	Err   error
}

// BackMsg is sent to go back to the previous screen
type BackMsg struct{}

// QuitMsg is sent to quit the application
type QuitMsg struct{}

// ViewEventMsg is sent to view a specific event's details
type ViewEventMsg struct {
	Event *career.CareerEvent
}

// EventActionMenuMsg is sent when an event is selected in the list to show its action menu
type EventActionMenuMsg struct {
	Event *career.CareerEvent
}

