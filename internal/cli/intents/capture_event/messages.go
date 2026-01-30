package capture_event

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
