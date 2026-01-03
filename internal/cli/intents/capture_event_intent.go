package intents

import (
	"fmt"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	// Initialize the review state with the previous event data.
	// The form will be pre-populated with existing event details.
	if i.context.PreviousEvent != nil {
		i.state.reviewState.Event = i.context.PreviousEvent
	}
	return nil
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
// It processes form input, validates data, and transitions to review state.
func (i *CaptureEventIntent) updateCaptureForm(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+s", "enter":
			// Submit form - transition to review state
			if i.state.reviewState.Event == nil {
				// Create a minimal event if none exists
				i.state.reviewState.Event = &career.CareerEvent{
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
					Tags:       make([]string, 0),
					Categories: make([]string, 0),
				}
			}

			// Transition to review state regardless of validation
			// Validation errors will be shown in the review state
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

	case FormCancelledMsg:
		// Form was cancelled
		i.setCancelled()
		return nil
	}

	return nil
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
		// Submission failed - mark as failed immediately
		i.setFailed(msg.Code, msg.Message, msg.Cause)
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
// In a real implementation, this would call the domain service to persist the event.
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

		// TODO: Call domain service to save the event.
		// For now, simulate successful submission with a small delay.
		time.Sleep(100 * time.Millisecond)
		return SubmitCompleteMsg{}
	}
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

// viewChooseStrategy renders the strategy selection UI with professional styling.
// Displays three capture strategy options with descriptions using lipgloss.
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

	// Build the content
	var content strings.Builder
	content.WriteString("\nChoose Capture Strategy\n\n")

	for _, s := range strategies {
		content.WriteString(fmt.Sprintf("  %s) %s\n", s.number, s.name))
		content.WriteString(fmt.Sprintf("     %s\n\n", s.desc))
	}

	content.WriteString("  q) Cancel\n")

	// Apply card styling
	cardStyle := lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Background(styles.ColorBackgroundCard).
		Foreground(styles.ColorTextPrimary)

	card := cardStyle.Render(content.String())

	// Add footer with instructions
	footerStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		MarginTop(1)

	footer := footerStyle.Render("Select strategy (1-3) or press 'q' to cancel")

	return lipgloss.JoinVertical(lipgloss.Left, card, footer)
}

// viewCaptureForm renders the form for capturing event details with lipgloss styling.
// Displays form fields with proper styling and validation error indication.
func (i *CaptureEventIntent) viewCaptureForm() string {
	var content strings.Builder
	content.WriteString("\nCapture Event Details\n\n")

	// Strategy info
	content.WriteString(fmt.Sprintf("Strategy: %s\n\n", i.context.CaptureStrategy))

	// Form fields
	content.WriteString("Description:\n")
	content.WriteString("[Enter event description...]\n\n")

	content.WriteString("Date: [YYYY-MM-DD]\n")
	content.WriteString("Company: [Company name]\n")
	content.WriteString("Project: [Project name]\n\n")

	content.WriteString("Tags: [Add tags...]\n")
	content.WriteString("Categories: [Select categories...]\n")

	// Apply card styling
	cardStyle := lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Background(styles.ColorBackgroundCard).
		Foreground(styles.ColorTextPrimary)

	card := cardStyle.Render(content.String())

	// Add footer with instructions
	footerStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		MarginTop(1)

	footer := footerStyle.Render("Press Tab to navigate, Ctrl+S to submit, Esc to cancel")

	return lipgloss.JoinVertical(lipgloss.Left, card, footer)
}

// viewReviewInferredEvent renders the review UI with professional styling.
// Displays captured event details, inferred bursts, and facts.
func (i *CaptureEventIntent) viewReviewInferredEvent() string {
	var content strings.Builder
	content.WriteString("\nReview Inferred Event\n\n")

	// Event summary
	if i.state.result != nil && i.state.result.Event != nil {
		title := i.state.result.Event.Text
		if len(title) > 40 {
			title = title[:37] + "..."
		}
		content.WriteString(fmt.Sprintf("Event: %s\n\n", title))
	}

	// Inferred bursts
	content.WriteString("Inferred Bursts:\n")
	if len(i.state.reviewState.AcceptedBursts) > 0 {
		for idx, burst := range i.state.reviewState.AcceptedBursts {
			burstTitle := burst.Name
			if len(burstTitle) > 35 {
				burstTitle = burstTitle[:32] + "..."
			}
			content.WriteString(fmt.Sprintf("  [✓] Burst %d: %s\n", idx+1, burstTitle))
		}
	} else {
		content.WriteString("  (No bursts detected)\n")
	}
	content.WriteString("\n")

	// Inferred facts
	content.WriteString("Inferred Facts:\n")
	if len(i.state.reviewState.AcceptedFacts) > 0 {
		for idx, fact := range i.state.reviewState.AcceptedFacts {
			desc := fact.Text
			if len(desc) > 35 {
				desc = desc[:32] + "..."
			}
			content.WriteString(fmt.Sprintf("  [✓] Fact %d: %s\n", idx+1, desc))
		}
	} else {
		content.WriteString("  (No facts detected)\n")
	}

	// Apply card styling
	cardStyle := lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Background(styles.ColorBackgroundCard).
		Foreground(styles.ColorTextPrimary)

	card := cardStyle.Render(content.String())

	// Add footer with instructions
	footerStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		MarginTop(1)

	footer := footerStyle.Render("Press Ctrl+S to submit, Esc to go back, 'e' to edit")

	return lipgloss.JoinVertical(lipgloss.Left, card, footer)
}

// viewSubmit renders the submit confirmation with professional styling.
// Displays a summary of the event to be submitted with confirmation options.
func (i *CaptureEventIntent) viewSubmit() string {
	var content strings.Builder
	content.WriteString("\nConfirm Submission\n\n")

	if i.state.result != nil && i.state.result.Event != nil {
		title := i.state.result.Event.Text
		if len(title) > 40 {
			title = title[:37] + "..."
		}
		content.WriteString(fmt.Sprintf("Event: %s\n", title))
		content.WriteString(fmt.Sprintf("Date: %s\n\n", i.state.result.Event.Date))
		content.WriteString(fmt.Sprintf("Bursts: %d\n", len(i.state.reviewState.AcceptedBursts)))
		content.WriteString(fmt.Sprintf("Facts: %d\n\n", len(i.state.reviewState.AcceptedFacts)))
	}

	content.WriteString("Ready to submit? Press Enter to confirm.\n")
	content.WriteString("Press Esc to cancel.\n")

	// Apply card styling
	cardStyle := lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Background(styles.ColorBackgroundCard).
		Foreground(styles.ColorTextPrimary)

	card := cardStyle.Render(content.String())

	// Add footer with instructions
	footerStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		MarginTop(1)

	footer := footerStyle.Render("Submitting event...")

	return lipgloss.JoinVertical(lipgloss.Left, card, footer)
}

// viewError renders an error state with professional styling.
// Displays error details and recovery options.
func (i *CaptureEventIntent) viewError() string {
	var content strings.Builder
	content.WriteString("\nError\n\n")

	if i.state.error != nil {
		code := i.state.error.Code
		if len(code) > 40 {
			code = code[:37] + "..."
		}
		content.WriteString(fmt.Sprintf("Code: %s\n\n", code))

		msg := i.state.error.Message
		if len(msg) > 40 {
			msg = msg[:37] + "..."
		}
		content.WriteString(fmt.Sprintf("Message: %s\n\n", msg))
	}

	content.WriteString("Press 'r' to retry or Esc to cancel.\n")

	// Apply card styling with error colors
	cardStyle := lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorderError).
		Background(styles.ColorBackgroundCard).
		Foreground(styles.ColorTextPrimary)

	card := cardStyle.Render(content.String())

	// Add footer with error indication
	footerStyle := lipgloss.NewStyle().
		Foreground(styles.ColorError).
		MarginTop(1)

	footer := footerStyle.Render("An error occurred")

	return lipgloss.JoinVertical(lipgloss.Left, card, footer)
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
