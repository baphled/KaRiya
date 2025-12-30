package models

import (
	"github.com/baphled/kariya/internal/domain/career"
)

// ViewEventMsg is sent when the user wants to view event details
type ViewEventMsg struct {
	Event *career.CareerEvent
}

// EventActionMenuMsg is sent when an event is selected to show its action menu
type EventActionMenuMsg struct {
	Event *career.CareerEvent
}
