// Package captureevent implements the CaptureEvent intent for capturing career events.
//
// Responsibilities:
//   - Orchestrating the event capture workflow (strategy -> form -> review -> submit)
//   - Managing editing modals for metadata, bursts, and facts during review
//   - Persisting captured events and associated enrichment data
//
// This package does NOT:
//   - Define UI components (those live in views/event/)
//   - Perform enrichment logic (delegated to service/career/)
//   - Handle form rendering (delegated to views/event/ and forms/)
package captureevent

// State represents the current workflow step of the capture event intent.
//
// The intent progresses through states linearly:
//
//	StateChooseStrategy -> StateForm -> StateReview -> StateSubmit
type State string

// Workflow states for the capture event intent.
const (
	// StateChooseStrategy is the initial state where the user picks quick or manual capture.
	StateChooseStrategy State = "choose_strategy"

	// StateForm is the state where the user fills in event details via a form.
	StateForm State = "form"

	// StateReview is the state where the user reviews inferred bursts and facts.
	StateReview State = "review"

	// StateSubmit is the final state where the event is persisted.
	StateSubmit State = "submit"
)

// EditingMode represents which editing sub-flow is active within the review state.
//
// When a mode is active, the corresponding modal is shown as an overlay.
type EditingMode string

// Editing modes for the review state.
const (
	// EditingModeNone means no editing modal is active.
	EditingModeNone EditingMode = ""

	// EditingModeMetadata means the metadata editor modal is active.
	EditingModeMetadata EditingMode = "metadata"

	// EditingModeBursts means the burst editor modal is active.
	EditingModeBursts EditingMode = "bursts"

	// EditingModeFacts means the fact editor modal is active.
	EditingModeFacts EditingMode = "facts"

	// EditingModeSkills means the skill editor modal is active.
	EditingModeSkills EditingMode = "skills"
)
