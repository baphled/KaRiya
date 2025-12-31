package models

import (
	"github.com/baphled/kariya/internal/domain/career"
	burstfact "github.com/baphled/kariya/internal/service/career/burst_fact"
)

// ViewEventMsg is sent when the user wants to view event details
type ViewEventMsg struct {
	Event *career.CareerEvent
}

// EventActionMenuMsg is sent when an event is selected to show its action menu
type EventActionMenuMsg struct {
	Event *career.CareerEvent
}

// ConfirmBurstSuggestionMsg is sent when user confirms a burst suggestion
type ConfirmBurstSuggestionMsg struct {
	Suggestion burstfact.BurstSuggestion
	EditedName string // Optional: user-edited burst name
	EditedDesc string // Optional: user-edited burst description
}

// RejectBurstSuggestionMsg is sent when user rejects a burst suggestion
type RejectBurstSuggestionMsg struct {
	Suggestion burstfact.BurstSuggestion
}

// EditBurstSuggestionMsg is sent when user wants to edit burst name/description
type EditBurstSuggestionMsg struct {
	Suggestion burstfact.BurstSuggestion
}

