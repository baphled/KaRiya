package capture_event

import "github.com/baphled/kariya/internal/domain/career"

// StrategySelectedMsg indicates the user selected a capture strategy.
type StrategySelectedMsg struct {
	// Strategy is the chosen capture approach (e.g. "quick", "manual").
	Strategy string
}

// FormSubmittedMsg indicates the form was submitted with event data.
type FormSubmittedMsg struct {
	// Event is the career event built from the form fields.
	Event *career.Event
}

// FormCancelledMsg indicates the form was cancelled.
type FormCancelledMsg struct{}

// ReviewConfirmedMsg indicates the user confirmed the review and is ready to submit.
type ReviewConfirmedMsg struct {
	// AcceptedBursts are the bursts the user chose to keep.
	AcceptedBursts []*career.Burst

	// AcceptedFacts are the facts the user chose to keep.
	AcceptedFacts []*career.Fact

	// RejectedItems maps item IDs to rejection reasons.
	RejectedItems map[string]string
}

// ReviewCancelledMsg indicates the user cancelled the review.
type ReviewCancelledMsg struct{}

// ReviewBackMsg indicates the user wants to go back to the form.
type ReviewBackMsg struct{}

// SubmitCompleteMsg indicates submission succeeded.
type SubmitCompleteMsg struct{}

// SubmitErrorMsg indicates submission failed.
type SubmitErrorMsg struct {
	// Code is a machine-readable error identifier (e.g. "PERSISTENCE_ERROR").
	Code string

	// Message is a human-readable error description.
	Message string

	// Cause is the underlying error that triggered the failure, if any.
	Cause error
}

// EnrichmentCompleteMsg indicates enrichment (burst/fact extraction) completed.
type EnrichmentCompleteMsg struct {
	// Bursts are the activity bursts detected from the event.
	Bursts []*career.Burst

	// Facts are the career facts extracted from the event.
	Facts []*career.Fact
}

// EnrichmentErrorMsg indicates enrichment failed (non-fatal).
type EnrichmentErrorMsg struct {
	// Code is a machine-readable error identifier.
	Code string

	// Message is a human-readable error description.
	Message string

	// Cause is the underlying error, if any.
	Cause error
}

// DismissModalMsg indicates the submit modal should be dismissed (after success).
type DismissModalMsg struct{}
