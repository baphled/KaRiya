package capture_event

import "github.com/baphled/kariya/internal/domain/career"

// StrategySelectedMsg is sent when the user selects a capture strategy.
type StrategySelectedMsg struct {
	// Strategy is the chosen capture approach (e.g. "quick", "manual").
	Strategy string
}

// FormSubmittedMsg is sent when the capture form is submitted with valid event data.
type FormSubmittedMsg struct {
	// Event is the career event built from the form fields.
	Event *career.Event
}

// FormCancelledMsg is sent when the user cancels the capture form.
type FormCancelledMsg struct{}

// ReviewConfirmedMsg is sent when the user confirms the review and is ready to submit.
type ReviewConfirmedMsg struct {
	// AcceptedBursts are the bursts the user chose to keep.
	AcceptedBursts []*career.Burst

	// AcceptedFacts are the facts the user chose to keep.
	AcceptedFacts []*career.Fact

	// RejectedItems maps item IDs to rejection reasons.
	RejectedItems map[string]string
}

// ReviewCancelledMsg is sent when the user cancels the review.
type ReviewCancelledMsg struct{}

// ReviewBackMsg is sent when the user navigates back from review to form.
type ReviewBackMsg struct{}

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

// EnrichmentCompleteMsg is sent when enrichment (burst/fact extraction) finishes.
type EnrichmentCompleteMsg struct {
	// Bursts are the activity bursts detected from the event.
	Bursts []*career.Burst

	// Facts are the career facts extracted from the event.
	Facts []*career.Fact
}

// EnrichmentErrorMsg is sent when enrichment fails (non-fatal; the event is
// still saved).
type EnrichmentErrorMsg struct {
	// Code is a machine-readable error identifier.
	Code string

	// Message is a human-readable error description.
	Message string

	// Cause is the underlying error, if any.
	Cause error
}

// DismissModalMsg is sent to dismiss the submit result modal (after a timed delay
// following a successful save).
type DismissModalMsg struct{}
