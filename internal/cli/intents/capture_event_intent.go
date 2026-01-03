package intents

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/cli/validation"
	"github.com/baphled/kariya/internal/domain/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
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

	// eventService is the CLI service for interacting with career events.
	// Injected via context for dependency management.
	eventService *service.CLIEventService
}

// NewCaptureEventIntent creates a new CaptureEvent intent.
func NewCaptureEventIntent(context *CaptureEventContext) (*CaptureEventIntent, error) {
	// Validate the context.
	if err := context.Validate(); err != nil {
		return nil, err
	}

	// Create the form model for capturing event details
	formModel := models.NewFormModel(context.CLIEventService)

	// Capture event intent created successfully
	return &CaptureEventIntent{
		context:      context,
		eventService: context.CLIEventService,
		state: &CaptureEventModel{
			context:      context,
			currentState: CaptureStateChooseStrategy,
			captureForm:  formModel,
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
	// Initialize the review state with the previous event data.
	// The form will be pre-populated with existing event details.
	if i.context.PreviousEvent != nil {
		i.state.reviewState.Event = i.context.PreviousEvent
	}
	// Return a no-op command to satisfy the intent lifecycle
	return func() tea.Msg { return nil }
}

// initializeFormForNew initializes the form for capturing a new event.
func (i *CaptureEventIntent) initializeFormForNew() tea.Cmd {
	// Create a fresh event with current timestamp.
	// The form will guide the user through data entry.
	i.state.reviewState.Event = &career.CareerEvent{
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Tags:       make([]string, 0),
		Categories: make([]string, 0),
	}
	// Return a no-op command to satisfy the intent lifecycle
	return func() tea.Msg { return nil }
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
// It processes form input, validates data, and transitions to review state.
// The form model handles all text input and field navigation.
func (i *CaptureEventIntent) updateCaptureForm(msg tea.Msg) tea.Cmd {
	// Delegate all messages to the form model to handle input and state
	_, formCmd := i.state.captureForm.Update(msg)

	// Check for special messages that indicate form completion or navigation
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			// User cancelled
			i.setCancelled()
			return nil

		case "esc":
			// Go back to strategy selection
			i.state.currentState = CaptureStateChooseStrategy
			return nil
		}

	case models.SubmitMsg:
		// Form submission completed
		if msg.Err != nil {
			// Form submission failed - show error
			i.state.error = &IntentError{
				Code:    "FORM_SUBMISSION_ERROR",
				Message: msg.Err.Error(),
				Cause:   msg.Err,
			}
			return nil
		}

		// Form submission succeeded
		if msg.Event == nil {
			i.setFailed("INVALID_FORM", "Form submission with nil event", nil)
			return nil
		}

		// Validate the event data
		if err := msg.Event.Validate(); err != nil {
			i.setFailed("VALIDATION_ERROR", fmt.Sprintf("Form validation failed: %v", err), err)
			return nil
		}

		i.state.reviewState.Event = msg.Event
		i.state.currentState = CaptureStateReview
		return nil
	}

	// Return the command from the form update
	return formCmd
}

// updateReviewInferredEvent handles messages while reviewing inferred bursts and facts.
// It processes review confirmations, edits, and transitions to submit state.
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
// It processes submission completion, errors, and retry attempts.
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
		// Submission failed - store error for display
		i.state.error = &IntentError{
			Code:    msg.Code,
			Message: msg.Message,
			Cause:   msg.Cause,
		}
		// Don't mark as failed yet - user can retry
		return nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			// Cancel submission
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
// It calls the domain service to persist the event to the database and optionally enriches it.
// Logs: Event submission start, validation results, service calls, and completion status.
func (i *CaptureEventIntent) performSubmit() tea.Cmd {
	return func() tea.Msg {
		// Validate event before submission
		if i.state.reviewState.Event == nil {
			return SubmitErrorMsg{
				Code:    "MISSING_EVENT",
				Message: "No event data to submit",
				Cause:   nil,
			}
		}

		// Validate event data
		if err := i.state.reviewState.Event.Validate(); err != nil {
			return SubmitErrorMsg{
				Code:    "VALIDATION_ERROR",
				Message: fmt.Sprintf("Event validation failed: %v", err),
				Cause:   err,
			}
		}

		// Ensure we have an event service
		if i.eventService == nil {
			return SubmitErrorMsg{
				Code:    "SERVICE_ERROR",
				Message: "Event service not initialized",
				Cause:   nil,
			}
		}

		event := i.state.reviewState.Event

		// Create a context with timeout for the submission
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Determine the capture mode based on the strategy
		mode := careerservice.ManualEntry
		if i.context.CaptureStrategy != "" {
			switch i.context.CaptureStrategy {
			case "quick":
				mode = careerservice.TimelineJournaling
			case "enriched":
				mode = careerservice.ManualEntry
			default:
				mode = careerservice.ManualEntry
			}
		}

		// Call the service to capture the event
		// The service handles persistence and any enrichment logic
		err := i.eventService.CaptureEvent(
			ctx,
			event.Text,
			event.Date,
			mode,
		)

		if err != nil {
			// Map service errors to intent errors
			return SubmitErrorMsg{
				Code:    "PERSISTENCE_ERROR",
				Message: fmt.Sprintf("Failed to save event: %v", err),
				Cause:   err,
			}
		}

		// For enriched strategy, perform enrichment
		if i.context.CaptureStrategy == "enriched" && i.context.CareerService != nil {
			if err := i.performEnrichment(ctx, event); err != nil {
				// Log enrichment error but don't fail the submission
				// The event is already saved successfully
				return SubmitCompleteMsg{}
			}
		}

		// Successfully submitted
		return SubmitCompleteMsg{}
	}
}

// performEnrichment performs AI-powered enrichment of the captured event.
// It suggests bursts and extracts facts from the event.
// Logs: Enrichment start, burst suggestion results, fact extraction results, and completion status.
func (i *CaptureEventIntent) performEnrichment(ctx context.Context, event *career.CareerEvent) error {
	if i.context.CareerService == nil {
		return fmt.Errorf("career service not available for enrichment")
	}

	// Suggest bursts for the event
	burstSuggestions, err := i.context.CareerService.SuggestBursts(ctx, []string{event.ID})
	if err == nil && len(burstSuggestions) > 0 {
		// Save burst suggestions and store them for review
		bursts, err := i.context.CareerService.SaveBurstSuggestions(ctx, burstSuggestions)
		if err == nil && len(bursts) > 0 {
			i.state.reviewState.InferredBursts = bursts
		}
	}

	// Extract facts from the event
	facts, err := i.context.CareerService.ExtractFactsFromEvent(ctx, event)
	if err == nil && len(facts) > 0 {
		// Store inferred facts for review
		for j := range facts {
			i.state.reviewState.InferredFacts = append(i.state.reviewState.InferredFacts, &facts[j])
		}
	}

	return nil
}


// validateEventWithDetails performs comprehensive validation of the event and provides detailed error messages.
// This uses the domain validators to check all event fields.
func (i *CaptureEventIntent) validateEventWithDetails(event *career.CareerEvent) error {
	if event == nil {
		return fmt.Errorf("event cannot be nil")
	}

	// Validate text
	eventValidator := validation.NewEventValidator()
	if err := eventValidator.ValidateText(event.Text); err != nil {
		return err
	}

	// Validate date
	if err := eventValidator.ValidateDate(event.Date); err != nil {
		return err
	}

	// Validate tags if present
	if len(event.Tags) > 0 {
		if err := eventValidator.ValidateTags(event.Tags); err != nil {
			return err
		}
	}

	// Validate using domain model
	if err := event.Validate(); err != nil {
		return err
	}

	return nil
}

// View renders the intent's current state.
func (i *CaptureEventIntent) View() string {
	if !i.active {
		return "CaptureEvent intent is not active"
	}

	// Show error view if there's an error
	if i.state.error != nil {
		return i.viewError()
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
// Displays three capture strategy options with descriptions.
func (i *CaptureEventIntent) viewChooseStrategy() string {
	strategies := []struct {
		number string
		name   string
		desc   string
	}{
		{"1", "Manual", "Manually enter event details"},
		{"2", "Quick", "Quick capture with minimal fields"},
		{"3", "Enriched", "Capture with AI-powered enrichment"},
	}

	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString("┌─ Choose Capture Strategy ─────────────────────┐\n")
	sb.WriteString("│                                                │\n")

	for _, s := range strategies {
		sb.WriteString(fmt.Sprintf("│  %s) %-40s │\n", s.number, s.name))
		sb.WriteString(fmt.Sprintf("│     %s                             │\n", s.desc))
		sb.WriteString("│                                                │\n")
	}

	sb.WriteString("│  q) Cancel                                     │\n")
	sb.WriteString("│                                                │\n")
	sb.WriteString("└────────────────────────────────────────────────┘\n")
	sb.WriteString("\nSelect strategy (1-3) or press 'q' to cancel:\n")

	return sb.String()
}

// viewCaptureForm renders the form for capturing event details.
// Delegates to the FormModel's View method to render the actual form.
func (i *CaptureEventIntent) viewCaptureForm() string {
	if i.state.captureForm == nil {
		return "Error: Form not initialized"
	}
	return i.state.captureForm.View()
}

// viewReviewInferredEvent renders the review UI for inferred bursts and facts.
// Displays the captured event details, inferred bursts, and facts with accept/reject options.
func (i *CaptureEventIntent) viewReviewInferredEvent() string {
	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString("┌─ Review Inferred Event ────────────────────────┐\n")
	sb.WriteString("│                                                │\n")

	// Event summary
	if i.state.result != nil && i.state.result.Event != nil {
		title := i.state.result.Event.Text
		if len(title) > 40 {
			title = title[:37] + "..."
		}
		sb.WriteString(fmt.Sprintf("│ Event: %s                    │\n", title))
		sb.WriteString("│                                                │\n")
	}

	// Inferred bursts
	sb.WriteString("│ Inferred Bursts:                               │\n")
	if len(i.state.reviewState.AcceptedBursts) > 0 {
		for idx, burst := range i.state.reviewState.AcceptedBursts {
			burstTitle := burst.Name
			if len(burstTitle) > 35 {
				burstTitle = burstTitle[:32] + "..."
			}
			sb.WriteString(fmt.Sprintf("│   [✓] Burst %d: %s              │\n", idx+1, burstTitle))
		}
	} else {
		sb.WriteString("│   (No bursts detected)                         │\n")
	}
	sb.WriteString("│                                                │\n")

	// Inferred facts
	sb.WriteString("│ Inferred Facts:                                │\n")
	if len(i.state.reviewState.AcceptedFacts) > 0 {
		for idx, fact := range i.state.reviewState.AcceptedFacts {
			desc := fact.Text
			if len(desc) > 35 {
				desc = desc[:32] + "..."
			}
			sb.WriteString(fmt.Sprintf("│   [✓] Fact %d: %s              │\n", idx+1, desc))
		}
	} else {
		sb.WriteString("│   (No facts detected)                          │\n")
	}
	sb.WriteString("│                                                │\n")
	sb.WriteString("└────────────────────────────────────────────────┘\n")
	sb.WriteString("\nPress Ctrl+S to submit, Esc to go back, 'e' to edit\n")

	return sb.String()
}

// viewSubmit renders the submit confirmation.
// Displays a summary of the event to be submitted with confirmation options.
func (i *CaptureEventIntent) viewSubmit() string {
	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString("┌─ Confirm Submission ───────────────────────────┐\n")
	sb.WriteString("│                                                │\n")

	if i.state.result != nil && i.state.result.Event != nil {
		title := i.state.result.Event.Text
		if len(title) > 40 {
			title = title[:37] + "..."
		}
		sb.WriteString(fmt.Sprintf("│ Event: %s                    │\n", title))
		sb.WriteString(fmt.Sprintf("│ Date: %s                      │\n", i.state.result.Event.Date))
		sb.WriteString("│                                                │\n")
		sb.WriteString(fmt.Sprintf("│ Bursts: %d                                    │\n", len(i.state.reviewState.AcceptedBursts)))
		sb.WriteString(fmt.Sprintf("│ Facts: %d                                     │\n", len(i.state.reviewState.AcceptedFacts)))
		sb.WriteString("│                                                │\n")
	}

	sb.WriteString("│ Ready to submit? Press Enter to confirm.       │\n")
	sb.WriteString("│ Press Esc to cancel.                           │\n")
	sb.WriteString("│                                                │\n")
	sb.WriteString("└────────────────────────────────────────────────┘\n")
	sb.WriteString("\nSubmitting event...\n")

	return sb.String()
}

// viewError renders an error state.
// Displays error details and recovery options.
func (i *CaptureEventIntent) viewError() string {
	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString("┌─ Error ─────────────────────────────────────────┐\n")
	sb.WriteString("│                                                │\n")

	if i.state.error != nil {
		code := i.state.error.Code
		if len(code) > 40 {
			code = code[:37] + "..."
		}
		sb.WriteString(fmt.Sprintf("│ Code: %s                        │\n", code))

		msg := i.state.error.Message
		if len(msg) > 40 {
			msg = msg[:37] + "..."
		}
		sb.WriteString(fmt.Sprintf("│ Message: %s                 │\n", msg))
		sb.WriteString("│                                                │\n")
	}

	sb.WriteString("│ Press 'r' to retry or Esc to cancel.           │\n")
	sb.WriteString("│                                                │\n")
	sb.WriteString("└────────────────────────────────────────────────┘\n")

	return sb.String()
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
