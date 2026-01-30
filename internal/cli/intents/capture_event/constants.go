// Package capture_event implements the CaptureEvent intent for capturing career events.
package capture_event

// State represents the current workflow state of the capture event intent.
type State string

// States for navigation within the capture event workflow.
const (
	StateChooseStrategy State = "choose_strategy"
	StateForm           State = "form"
	StateReview         State = "review"
	StateSubmit         State = "submit"
)

// EditingMode represents the active editing sub-flow within the review state.
type EditingMode string

// Editing modes for the ReviewInferredEvent sub-flow.
const (
	EditingModeNone     EditingMode = ""
	EditingModeMetadata EditingMode = "metadata"
	EditingModeBursts   EditingMode = "bursts"
	EditingModeFacts    EditingMode = "facts"
)
