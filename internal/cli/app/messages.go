package app

import (
	"github.com/baphled/kariya/internal/domain/career"
)

// Screen represents a screen in the application
type Screen string

// FormSubmittedMsg is sent when a form is successfully submitted
type FormSubmittedMsg struct {
	Event *career.CareerEvent
	Err   error
}

// BackMsg is sent to go back to the previous screen
type BackMsg struct{}

// QuitMsg is sent to quit the application
type QuitMsg struct{}

// ConfirmBurstMsg is sent when user confirms a burst suggestion
type ConfirmBurstMsg struct {
	Burst *career.Burst
}

// RejectBurstSuggestionMsg is sent when user rejects a burst suggestion
type RejectBurstSuggestionMsg struct {
	EventIDs []string
}

// BurstProcessingCompleteMsg is sent when burst suggestion workflow is done
type BurstProcessingCompleteMsg struct {
	ConfirmedBursts []career.Burst
	RejectedGroups  [][]string
}

// Screen constants for navigation
const (
	ListScreen Screen = "list"
)
