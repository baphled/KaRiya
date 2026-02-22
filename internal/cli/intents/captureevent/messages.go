package captureevent

import (
	"github.com/baphled/kariya/internal/domain/career"
	burstfact "github.com/baphled/kariya/internal/service/career/burstfact"
	"github.com/baphled/kariya/internal/service/career/skillinference"
)

// SubmitCompleteMsg is sent when event persistence succeeds.
// Inference results arrive separately via InferenceCompleteMsg.
type SubmitCompleteMsg struct{}

// InferenceCompleteMsg is sent when background LLM inference completes.
// This is decoupled from the save path so that event persistence returns
// promptly without waiting for potentially slow LLM calls.
type InferenceCompleteMsg struct {
	InferredBursts []*career.Burst
	// InferredBurstSuggestions carries the raw suggestions from the inference
	// service, preserving ConfidenceScore which career.Burst does not hold.
	InferredBurstSuggestions []burstfact.BurstSuggestion
	InferredFacts            []*career.Fact
	InferredSkills           []skillinference.SkillSuggestion
}

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

// PostSavePersistenceCompleteMsg is sent when post-save review persistence
// (skills, facts, bursts) completes successfully. This replaces the inline
// persistence in HandleSubmit so that the UI thread is not blocked.
type PostSavePersistenceCompleteMsg struct {
	Event  *career.Event
	Bursts []*career.Burst
	Facts  []*career.Fact
	Skills []*career.Skill
}

// BurstProcessingCompleteMsg is sent when burst suggestion workflow is done.
type BurstProcessingCompleteMsg struct {
	ConfirmedCount int
	RejectedCount  int
}
