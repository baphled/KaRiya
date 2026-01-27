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

// Note: BackMsg, QuitMsg, ConfirmBurstMsg, RejectBurstSuggestionMsg, and
// BurstProcessingCompleteMsg have been moved to internal/cli/models/messages.go
// to avoid duplication. Import from models package when needed.

// Screen constants for navigation
const (
	ListScreen Screen = "list"
)

// IntentCompletedMsg signals the completion of an intent.
type IntentCompletedMsg struct{}
