package intents

import (
	"fmt"

	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

// Custom message types for state transitions.
// These are used to communicate between states and modal sub-flows.

// StrategySelectedMsg indicates the user selected a capture strategy.
type StrategySelectedMsg struct {
	Strategy string
}

// FormSubmittedMsg indicates the form was submitted with event data.
type FormSubmittedMsg struct {
	Event *career.CareerEvent
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
	domainService interface{} // nolint:unused
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
// The strategy is typically pre-selected by the router, but this allows the user to confirm or change it.
func (i *CaptureEventIntent) updateChooseStrategy(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "1":
			// Manual strategy selected
			i.state.currentState = CaptureStateForm
			return nil

		case "2":
			// Quick strategy selected
			i.state.currentState = CaptureStateForm
			return nil

		case "3":
			// Enriched strategy selected
			i.state.currentState = CaptureStateForm
			return nil

		case "q", "ctrl+c":
			// User cancelled
			i.setCancelled()
			return nil

		case "enter":
			// Confirm current strategy and move to form
			i.state.currentState = CaptureStateForm
			return nil
		}

	case StrategySelectedMsg:
		// Strategy was selected (possibly by router or other component)
		i.state.currentState = CaptureStateForm
		return nil
	}

	return nil
}

// updateCaptureForm handles messages while capturing event details.
func (i *CaptureEventIntent) updateCaptureForm(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+s", "enter":
			// Submit form - transition to review state
			// In a real implementation, validate form and get data
			event := &career.CareerEvent{
				// TODO: Populate from form data
			}
			i.state.reviewState.Event = event
			i.state.currentState = CaptureStateReview
			return nil

		case "q", "ctrl+c":
			// Cancel form
			i.setCancelled()
			return nil

		case "esc":
			// Go back to strategy selection
			i.state.currentState = CaptureStateChooseStrategy
			return nil
		}

	case FormSubmittedMsg:
		// Form was submitted with event data
		i.state.reviewState.Event = msg.Event
		i.state.currentState = CaptureStateReview
		return nil

	case FormCancelledMsg:
		// Form was cancelled
		i.setCancelled()
		return nil
	}

	return nil
}

// updateReviewInferredEvent handles messages while reviewing inferred bursts and facts.
func (i *CaptureEventIntent) updateReviewInferredEvent(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+s", "enter":
			// Confirm review and submit
			i.state.currentState = CaptureStateSubmit
			return i.performSubmit()

		case "q", "ctrl+c":
			// Cancel review
			i.setCancelled()
			return nil

		case "esc":
			// Go back to form
			i.state.currentState = CaptureStateForm
			return nil

		case "e":
			// Edit metadata (modal sub-flow)
			i.state.reviewState.EditingMode = EditingModeMetadata
			return nil

		case "b":
			// Edit bursts (modal sub-flow)
			i.state.reviewState.EditingMode = EditingModeBursts
			return nil

		case "f":
			// Edit facts (modal sub-flow)
			i.state.reviewState.EditingMode = EditingModeFacts
			return nil
		}

	case ReviewConfirmedMsg:
		// Review was confirmed with accepted/rejected items
		i.state.reviewState.AcceptedBursts = msg.AcceptedBursts
		i.state.reviewState.AcceptedFacts = msg.AcceptedFacts
		i.state.reviewState.RejectedItems = msg.RejectedItems
		i.state.currentState = CaptureStateSubmit
		return i.performSubmit()

	case ReviewCancelledMsg:
		// Review was cancelled
		i.setCancelled()
		return nil

	case ReviewBackMsg:
		// User wants to go back to form
		i.state.currentState = CaptureStateForm
		return nil
	}

	return nil
}

// updateSubmit handles messages while submitting the event.
func (i *CaptureEventIntent) updateSubmit(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case SubmitCompleteMsg:
		// Submission succeeded
		result := &CaptureEventResult{
			Event:          i.state.reviewState.Event,
			Bursts:         i.state.reviewState.AcceptedBursts,
			Facts:          i.state.reviewState.AcceptedFacts,
			AcceptedFields: make(map[string]bool),
			RejectedFields: i.state.reviewState.RejectedItems,
		}
		i.setCompleted(result)
		return nil

	case SubmitErrorMsg:
		// Submission failed
		i.setFailed(msg.Code, msg.Message, msg.Cause)
		return nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			// Cancel submission (if possible)
			i.setCancelled()
			return nil

		case "r":
			// Retry submission
			return i.performSubmit()
		}
	}

	return nil
}

// performSubmit performs the actual submission of the event.
// In a real implementation, this would call the domain service.
func (i *CaptureEventIntent) performSubmit() tea.Cmd {
	return func() tea.Msg {
		// TODO: Call domain service to save the event
		// For now, simulate successful submission
		return SubmitCompleteMsg{}
	}
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

