package captureevent

import (
	"github.com/baphled/kariya/internal/domain/career"
	burstfact "github.com/baphled/kariya/internal/service/career/burstfact"
)

// SubmitCompleteMsg is sent when event persistence succeeds.
type SubmitCompleteMsg struct{}

// SubmitErrorMsg is sent when event persistence fails.
type SubmitErrorMsg struct {
	// Code is a machine-readable error identifier (e.g. "PERSISTENCE_ERROR").
	Code string

	// Message is a human-readable error description.
	Message string

	// Cause is the underlying error that triggered the failure, if any.
	Cause error
}

// DismissModalMsg is sent to dismiss the submit result modal (after a timed delay
// following a successful save).
type DismissModalMsg struct{}

// SubmitMsg represents a form submission result.
type SubmitMsg struct {
	Event *career.Event
	Err   error
}

// QuitMsg is sent when user wants to quit the application.
type QuitMsg struct{}

// BackMsg is sent when the user wants to go back to the previous screen.
type BackMsg struct{}

// ConfirmBurstMsg is sent when user confirms a burst.
type ConfirmBurstMsg struct {
	Burst *career.Burst
}

// RejectBurstSuggestionMsg is sent when user rejects a burst suggestion.
type RejectBurstSuggestionMsg struct {
	Suggestion burstfact.BurstSuggestion
}

// BurstProcessingCompleteMsg is sent when burst suggestion workflow is done.
type BurstProcessingCompleteMsg struct {
	ConfirmedCount int
	RejectedCount  int
}
