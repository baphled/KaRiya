package app

import (
	"github.com/baphled/kariya/internal/domain/career"
)

// Screen represents a screen in the application.
type Screen string

// FormSubmittedMsg is sent when a form is successfully submitted.
type FormSubmittedMsg struct {
	Event *career.Event
	Err   error
}

// Moved to models package to avoid duplication:
// BackMsg, QuitMsg, ConfirmBurstMsg, RejectBurstSuggestionMsg, BurstProcessingCompleteMsg
// Import from models package when needed.

// Screen constants for navigation.
const (
	ListScreen Screen = "list"
)

// IntentCompletedMsg signals the completion of an intent.
type IntentCompletedMsg struct{}
