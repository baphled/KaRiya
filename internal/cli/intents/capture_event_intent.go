package intents

import (
	"fmt"

	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

// CaptureEventIntent implements the Intent interface for capturing career events.
// It owns the complete lifecycle of event capture, including:
// - Choosing capture strategy (manual, quick, enriched)
// - Capturing event details via form
// - Reviewing and refining inferred bursts and facts
// - Submitting the final event
type CaptureEventIntent struct {
	// context is the input context passed to the intent.
	context *CaptureEventContext

	// state represents the current state of the intent.
	state *CaptureEventModel

	// active indicates whether this intent is currently active.
	active bool

	// result is the final result of the intent (set when complete).
	result *IntentResult[*CaptureEventResult]

	// domainService is the service for interacting with the domain.
	// Injected via constructor for testability.
	domainService interface{} // TODO: Define a proper domain service interface
}

// NewCaptureEventIntent creates a new CaptureEvent intent.
func NewCaptureEventIntent(context *CaptureEventContext) (*CaptureEventIntent, error) {
	// Validate the context.
	if err := context.Validate(); err != nil {
		return nil, err
	}

	return &CaptureEventIntent{
		context: context,
		state: &CaptureEventModel{
			context:      context,
			currentState: CaptureStateChooseStrategy,
			reviewState: &ReviewInferredEventState{
				AcceptedBursts: make([]*career.Burst, 0),
				AcceptedFacts:  make([]*career.Fact, 0),
				RejectedItems:  make(map[string]string),
			},
			result: &CaptureEventResult{
				AcceptedFields: make(map[string]bool),
				RejectedFields: make(map[string]string),
			},
		},
		active: true,
	}, nil
}

// Init is called when the intent is activated.
func (i *CaptureEventIntent) Init() tea.Cmd {
	// If this is an edit operation, initialize the form with the previous event's data.
	if i.context.PreviousEvent != nil {
		return i.initializeFormForEdit()
	}

	// Otherwise, initialize for a new event.
	return i.initializeFormForNew()
}

// initializeFormForEdit initializes the form for editing an existing event.
func (i *CaptureEventIntent) initializeFormForEdit() tea.Cmd {
	// TODO: Implement form initialization with previous event data.
	// This should:
	// - Populate the form fields with existing event data
	// - Skip the strategy selection (or allow changing strategy)
	// - Return any startup commands (e.g., fetch related data)
	return nil
}

// initializeFormForNew initializes the form for capturing a new event.
func (i *CaptureEventIntent) initializeFormForNew() tea.Cmd {
	// TODO: Implement form initialization for new event.
	// This should:
	// - Create a fresh form with empty fields
	// - Initialize with the capture strategy
	// - Return any startup commands
	return nil
}

// Update processes a message in the intent.
// This is Bubble Tea's standard Update, constrained to intent-local state.
func (i *CaptureEventIntent) Update(msg tea.Msg) tea.Cmd {
	if !i.active {
		return nil
	}

	switch i.state.currentState {
	case CaptureStateChooseStrategy:
		return i.updateChooseStrategy(msg)

	case CaptureStateForm:
		return i.updateCaptureForm(msg)

	case CaptureStateReview:
		return i.updateReviewInferredEvent(msg)

	case CaptureStateSubmit:
		return i.updateSubmit(msg)

	default:
		return nil
	}
}

// updateChooseStrategy handles messages while choosing capture strategy.
func (i *CaptureEventIntent) updateChooseStrategy(msg tea.Msg) tea.Cmd {
	// TODO: Implement strategy selection UI.
	// This should:
	// - Display options: Manual, Quick, Enriched
	// - Transition to CaptureStateForm when selected
	// - Transition to CaptureStateCancelled if user quits
	return nil
}

// updateCaptureForm handles messages while capturing event details.
func (i *CaptureEventIntent) updateCaptureForm(msg tea.Msg) tea.Cmd {
	// TODO: Implement form update logic.
	// This should:
	// - Delegate to the form model's Update
	// - Transition to CaptureStateReview when submitted
	// - Transition to CaptureStateCancelled if user quits
	return nil
}

// updateReviewInferredEvent handles messages while reviewing inferred bursts and facts.
func (i *CaptureEventIntent) updateReviewInferredEvent(msg tea.Msg) tea.Cmd {
	// TODO: Implement review UI update logic.
	// This should:
	// - Display inferred bursts and facts
	// - Allow inline editing via modal sub-flows:
	//   - EditMetadata: Modify event details
	//   - EditBursts: Review and refine bursts
	//   - EditFacts: Review and refine facts
	// - Transition to CaptureStateSubmit when confirmed
	// - Transition back to CaptureStateForm if user wants to re-edit
	// - Transition to CaptureStateCancelled if user quits
	return nil
}

// updateSubmit handles messages while submitting the event.
func (i *CaptureEventIntent) updateSubmit(msg tea.Msg) tea.Cmd {
	// TODO: Implement submit logic.
	// This should:
	// - Call the domain service to save the event
	// - Transition to a completed state and return the result
	// - Handle errors and transition to a failed state if necessary
	return nil
}

// View renders the intent's current state.
func (i *CaptureEventIntent) View() string {
	if !i.active {
		return "CaptureEvent intent is not active"
	}

	switch i.state.currentState {
	case CaptureStateChooseStrategy:
		return i.viewChooseStrategy()

	case CaptureStateForm:
		return i.viewCaptureForm()

	case CaptureStateReview:
		return i.viewReviewInferredEvent()

	case CaptureStateSubmit:
		return i.viewSubmit()

	default:
		return fmt.Sprintf("Unknown state: %s", i.state.currentState)
	}
}

// viewChooseStrategy renders the strategy selection UI.
func (i *CaptureEventIntent) viewChooseStrategy() string {
	// TODO: Implement strategy selection view.
	return "Choose Capture Strategy\n\n1. Manual\n2. Quick\n3. Enriched"
}

// viewCaptureForm renders the form for capturing event details.
func (i *CaptureEventIntent) viewCaptureForm() string {
	// TODO: Implement form view.
	// This should delegate to the form model's View.
	return "Capture Event Form"
}

// viewReviewInferredEvent renders the review UI for inferred bursts and facts.
func (i *CaptureEventIntent) viewReviewInferredEvent() string {
	// TODO: Implement review view.
	// This should display:
	// - The captured event details
	// - Inferred bursts with accept/reject options
	// - Inferred facts with accept/reject options
	// - Navigation options to edit each section
	return "Review Inferred Event"
}

// viewSubmit renders the submit confirmation.
func (i *CaptureEventIntent) viewSubmit() string {
	// TODO: Implement submit view (e.g., a spinner or confirmation message).
	return "Submitting event..."
}

// IsActive returns true if this intent is currently active.
func (i *CaptureEventIntent) IsActive() bool {
	return i.active
}

// GetResult returns the final result of the intent (if complete).
// This is called by the intent router to check if the intent has completed.
func (i *CaptureEventIntent) GetResult() *IntentResult[*CaptureEventResult] {
	return i.result
}

// setCompleted marks the intent as completed with a result.
func (i *CaptureEventIntent) setCompleted(result *CaptureEventResult) {
	i.result = NewCompletedResult(result)
	i.active = false
}

// setPartial marks the intent as partially completed with a result.
func (i *CaptureEventIntent) setPartial(result *CaptureEventResult, code, message string) {
	i.result = NewPartialResult(result, code, message)
	i.active = false
}

// setCancelled marks the intent as cancelled by the user.
func (i *CaptureEventIntent) setCancelled() {
	i.result = NewCancelledResult[*CaptureEventResult]()
	i.active = false
}

// setFailed marks the intent as failed with an error.
func (i *CaptureEventIntent) setFailed(code, message string, cause error) {
	i.result = NewFailedResult[*CaptureEventResult](code, message, cause)
	i.active = false
}

// Result returns the intent's result if it has completed, or nil if still active.
// This implements the Intent interface.
func (i *CaptureEventIntent) Result() *IntentResult[interface{}] {
	if i.result == nil {
		return nil
	}
	// Convert typed result to interface result
	return &IntentResult[interface{}]{
		Status:   i.result.Status,
		Data:     i.result.Data,
		Error:    i.result.Error,
		Metadata: i.result.Metadata,
	}
}
