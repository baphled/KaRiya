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

// ViewFactsMsg is sent when user wants to view facts for an event
type ViewFactsMsg struct {
	EventID string
}

// EditFactMsg is sent when user wants to edit a specific fact
type EditFactMsg struct {
	Fact *career.Fact
}

// ConfirmFactMsg is sent when user confirms/accepts a fact
type ConfirmFactMsg struct {
	Fact *career.Fact
}

// RejectFactMsg is sent when user rejects a fact
type RejectFactMsg struct {
	FactID string
}

// FactConfirmedMsg is sent after a fact has been successfully confirmed
type FactConfirmedMsg struct {
	Fact *career.Fact
	Err  error
}

// FactRejectedMsg is sent after a fact has been successfully rejected
type FactRejectedMsg struct {
	FactID string
	Err    error
}

// SaveFactMsg is sent when a fact is saved from the editor
type SaveFactMsg struct {
	Fact *career.Fact
}
