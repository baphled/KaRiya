package capture_event

import "github.com/baphled/kariya/internal/domain/career"

// StrategySelectedMsg indicates the user selected a capture strategy.
type StrategySelectedMsg struct {
	Strategy string
}

// FormSubmittedMsg indicates the form was submitted with event data.
type FormSubmittedMsg struct {
	Event *career.Event
}

// FormCancelledMsg indicates the form was cancelled.
type FormCancelledMsg struct{}

// ReviewConfirmedMsg indicates the user confirmed the review and is ready to submit.
type ReviewConfirmedMsg struct {
	AcceptedBursts []*career.Burst
	AcceptedFacts  []*career.Fact
	RejectedItems  map[string]string
}

// ReviewCancelledMsg indicates the user cancelled the review.
type ReviewCancelledMsg struct{}

// ReviewBackMsg indicates the user wants to go back to the form.
type ReviewBackMsg struct{}

// SubmitCompleteMsg indicates submission succeeded.
type SubmitCompleteMsg struct{}

// SubmitErrorMsg indicates submission failed.
type SubmitErrorMsg struct {
	Code    string
	Message string
	Cause   error
}

// EnrichmentCompleteMsg indicates enrichment (burst/fact extraction) completed.
type EnrichmentCompleteMsg struct {
	Bursts []*career.Burst
	Facts  []*career.Fact
}

// EnrichmentErrorMsg indicates enrichment failed (non-fatal).
type EnrichmentErrorMsg struct {
	Code    string
	Message string
	Cause   error
}

// DismissModalMsg indicates the submit modal should be dismissed (after success).
type DismissModalMsg struct{}
