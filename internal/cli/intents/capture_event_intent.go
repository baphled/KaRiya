package intents

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/domain/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	burstfact "github.com/baphled/kariya/internal/service/career/burst_fact"
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
	// Embed BaseIntent for terminal awareness, logo, and state management
	*BaseIntent

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

	// Create BaseIntent for terminal awareness and state management
	base := NewBaseIntent()

	// Capture event intent created successfully
	return &CaptureEventIntent{
		BaseIntent:   base,
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
		i.initializeFormForEdit()
		return nil
	}

	// Otherwise, initialize for a new event.
	i.initializeFormForNew()
	return func() tea.Msg { return nil }
}

// initializeFormForEdit initializes the form for editing an existing event.
func (i *CaptureEventIntent) initializeFormForEdit() tea.Cmd {
	// Initialize the review state with the previous event data.
	// The form will be pre-populated with existing event details.
	if i.context.PreviousEvent != nil {
		i.state.reviewState.Event = i.context.PreviousEvent

		// Edit mode always uses Manual strategy with all fields shown
		i.state.strategy = StrategyManual
		i.state.captureForm.SetStrategy(string(StrategyManual))
		i.state.showOptionalFields = true

		// Skip strategy selection and go straight to form
		i.state.currentState = CaptureStateForm
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
// Supports arrow key/vim navigation to select between Quick and Manual strategies.
func (i *CaptureEventIntent) updateChooseStrategy(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			// Navigate up in strategy list
			if i.state.selectedStrategyIndex > 0 {
				i.state.selectedStrategyIndex--
			}
			return nil

		case "down", "j":
			// Navigate down in strategy list (0=Quick, 1=Manual)
			if i.state.selectedStrategyIndex < 1 {
				i.state.selectedStrategyIndex++
			}
			return nil

		case "enter":
			// Confirm selected strategy and transition to form
			strategies := []CaptureStrategy{StrategyQuick, StrategyManual}
			i.state.strategy = strategies[i.state.selectedStrategyIndex]

			// Configure form based on selected strategy
			i.state.captureForm.SetStrategy(string(i.state.strategy))

			// For quick mode, pre-fill date with today
			if i.state.strategy == StrategyQuick {
				// Date will be set to today automatically in the submit handler
			}

			i.state.currentState = CaptureStateForm
			return nil

		case "q", "ctrl+c":
			// User cancelled
			i.setCancelled()
			return nil

		case "esc":
			// User cancelled (this is the root state)
			i.setCancelled()
			return nil

		case "m":
			// Return to main menu
			i.setCancelled()
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
// Users can press Ctrl+S to submit the form, or Tab+Enter to submit via the button.
func (i *CaptureEventIntent) updateCaptureForm(msg tea.Msg) tea.Cmd {
	// Delegate all messages to the form model to handle input and state
	_, formCmd := i.state.captureForm.Update(msg)

	// Check for special messages that indicate form completion or navigation
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+s":
			// User pressed Ctrl+S to submit the form
			// Trigger form submission
			return i.state.captureForm.SubmitForm()

		case "q", "ctrl+c":
			// User cancelled
			i.setCancelled()
			return nil

		case "esc":
			// Go back to strategy selection
			i.state.currentState = CaptureStateChooseStrategy
			return nil

		case "m":
			// Return to main menu
			i.setCancelled()
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

	case FormSubmittedMsg:
		// Handle test/legacy FormSubmittedMsg
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
	// If a modal is active, pass updates to it
	switch i.state.reviewState.EditingMode {
	case EditingModeMetadata:
		if i.state.reviewState.metadataModal != nil {
			modal, cmd := i.state.reviewState.metadataModal.Update(msg)
			i.state.reviewState.metadataModal = modal.(*models.MetadataEditorModelNew)

			// Check if modal completed
			if i.state.reviewState.metadataModal.IsSubmitted() {
				// Apply changes to event
				i.state.reviewState.Event = i.state.reviewState.metadataModal.GetEvent()
				// Clear modal and editing mode
				i.state.reviewState.metadataModal = nil
				i.state.reviewState.EditingMode = EditingModeNone
			} else if i.state.reviewState.metadataModal.IsCancelled() {
				// Clear modal without applying changes
				i.state.reviewState.metadataModal = nil
				i.state.reviewState.EditingMode = EditingModeNone
			}
			return cmd
		}

	case EditingModeBursts:
		if i.state.reviewState.burstModal != nil {
			modal, cmd := i.state.reviewState.burstModal.Update(msg)
			i.state.reviewState.burstModal = modal.(*models.BurstSuggestionModelNew)

			// Check if modal completed
			// TODO: Add completion check when BurstSuggestionModelNew has IsComplete/IsCancelled methods
			// For now, allow Esc to exit
			if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "esc" {
				i.state.reviewState.burstModal = nil
				i.state.reviewState.EditingMode = EditingModeNone
			}
			return cmd
		}

	case EditingModeFacts:
		if i.state.reviewState.factModal != nil {
			modal, cmd := i.state.reviewState.factModal.Update(msg)
			i.state.reviewState.factModal = modal.(*models.FactEditorModelNew)

			// Check if modal completed
			if i.state.reviewState.factModal.IsSubmitted() {
				// Apply changes to fact
				// TODO: Update the inferred facts list with edited fact
				i.state.reviewState.factModal = nil
				i.state.reviewState.EditingMode = EditingModeNone
			} else if i.state.reviewState.factModal.IsCancelled() {
				// Clear modal without applying changes
				i.state.reviewState.factModal = nil
				i.state.reviewState.EditingMode = EditingModeNone
			}
			return cmd
		}
	}

	// Normal review handling (no modal active)
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

		case "m":
			// Return to main menu
			i.setCancelled()
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

		case "a":
			// Accept currently selected item
			i.acceptCurrentItem()
			return nil

		case "r":
			// Reject currently selected item
			i.rejectCurrentItem()
			return nil

		case "j", "down":
			// Navigate down through items
			totalItems := len(i.state.reviewState.InferredBursts) + len(i.state.reviewState.InferredFacts)
			if totalItems > 0 {
				i.state.reviewState.SelectedIndex++
				if i.state.reviewState.SelectedIndex >= totalItems {
					i.state.reviewState.SelectedIndex = 0
				}
				// Update SelectedItemType based on new index
				if i.state.reviewState.SelectedIndex < len(i.state.reviewState.InferredBursts) {
					i.state.reviewState.SelectedItemType = "burst"
				} else {
					i.state.reviewState.SelectedItemType = "fact"
					i.state.reviewState.SelectedIndex = i.state.reviewState.SelectedIndex - len(i.state.reviewState.InferredBursts)
				}
			}
			return nil

		case "k", "up":
			// Navigate up through items
			totalItems := len(i.state.reviewState.InferredBursts) + len(i.state.reviewState.InferredFacts)
			if totalItems > 0 {
				i.state.reviewState.SelectedIndex--
				if i.state.reviewState.SelectedIndex < 0 {
					i.state.reviewState.SelectedIndex = totalItems - 1
				}
				// Update SelectedItemType based on new index
				if i.state.reviewState.SelectedIndex < len(i.state.reviewState.InferredBursts) {
					i.state.reviewState.SelectedItemType = "burst"
				} else {
					i.state.reviewState.SelectedItemType = "fact"
					i.state.reviewState.SelectedIndex = i.state.reviewState.SelectedIndex - len(i.state.reviewState.InferredBursts)
				}
			}
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

		case "esc":
			// Go back to review state (keep error visible per user preference)
			i.state.currentState = CaptureStateReview
			return nil

		case "m":
			// Return to main menu
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

		// For quick mode, default date to today if not set
		if i.state.strategy == StrategyQuick && event.Date.IsZero() {
			event.Date = time.Now()
		}

		// Always use ManualEntry mode (mode selector has been removed from UI)
		mode := careerservice.ManualEntry

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

		// Save any accepted facts from review that might have been manually edited/added
		// Note: Facts from enrichment are already saved in performEnrichment()
		// This is a safety check for any facts that might have been added during review
		if i.context.CareerService != nil && len(i.state.reviewState.AcceptedFacts) > 0 {
			for _, fact := range i.state.reviewState.AcceptedFacts {
				// Only save facts that don't have an ID yet (haven't been saved)
				// Facts from enrichment already have IDs
				if fact.ID == "" {
					// Ensure fact is linked to the saved event
					if fact.SourceEventID == "" {
						fact.SourceEventID = event.ID
					}

					// Save the fact
					if err := i.context.CareerService.SaveFact(ctx, fact); err != nil {
						// Log error but don't fail the entire submission
						// Event is already saved successfully
						continue
					}
				}
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
		// Persist each extracted fact to the database
		for j := range facts {
			fact := &facts[j]

			// Set source event ID (linking fact to this event)
			fact.SourceEventID = event.ID

			// Save fact to repository
			if err := i.context.CareerService.SaveFact(ctx, fact); err != nil {
				// Log warning but continue with other facts
				// Fact extraction is an enhancement, not critical to event capture
				continue
			}

			// Store saved fact for review
			i.state.reviewState.InferredFacts = append(i.state.reviewState.InferredFacts, fact)
		}
	}

	return nil
}

// validateEventWithDetails performs comprehensive validation of the event and provides detailed error messages.

// getStateName returns a human-readable name for the current state.
func (i *CaptureEventIntent) getStateName() string {
	switch i.state.currentState {
	case CaptureStateChooseStrategy:
		return "Choose Strategy"
	case CaptureStateForm:
		return "Enter Details"
	case CaptureStateReview:
		return "Review"
	case CaptureStateSubmit:
		return "Submit"
	default:
		return string(i.state.currentState)
	}
}

// getStateContent returns the content for the current state.
func (i *CaptureEventIntent) getStateContent() string {
	// Show error content if there's an error and not already showing modal
	if i.state.error != nil && !i.HasError() {
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

// getContextHelp returns context-aware help text for the current state.
func (i *CaptureEventIntent) getContextHelp() string {
	base := "q Quit  m Main Menu"

	switch i.state.currentState {
	case CaptureStateChooseStrategy:
		return CombineFooters(NavigationFooter(), base)
	case CaptureStateForm:
		// Show different help based on strategy
		if i.state.strategy == StrategyManual {
			return CombineFooters(FormFooter(), "Ctrl+O Toggle fields", base)
		}
		return CombineFooters(FormFooter(), base)
	case CaptureStateReview:
		// Show different help when modal is active
		if i.state.reviewState.EditingMode != EditingModeNone {
			return "Editing... | Esc Cancel  Enter Save"
		}
		return CombineFooters(NavigationFooter(), "e Edit  b Bursts  f Facts  a Accept  r Reject  j/k Navigate", base)
	case CaptureStateSubmit:
		return CombineFooters("Enter Continue  Esc Back", base)
	default:
		return base
	}
}

// View renders the intent's current state using StandardView.
func (i *CaptureEventIntent) View() string {
	if !i.active {
		return "CaptureEvent intent is not active"
	}

	// Create standard view with breadcrumbs
	view := i.CreateViewWithBreadcrumbs("Main Menu", "Capture Event", i.getStateName())

	// Reduce logo spacing on small terminals to maximize form visibility
	if info := i.GetTerminalInfo(); info != nil && info.Height < 30 {
		if logo := i.GetLogo(); logo != nil {
			view.WithLogo(logo, 0) // No spacing above logo for small terminals
		}
	}

	// Sync state from CaptureEventModel to BaseIntent for modal display
	if i.state.error != nil {
		i.SetError(i.state.error)
	}

	// Get content for current state
	content := i.getStateContent()
	view.WithContent(content)

	// Get context-aware help
	help := i.getContextHelp()
	view.WithHelp(help).WithFooterSeparator(true)

	return view.Render()
}

// viewChooseStrategy renders the strategy selection UI.
// Displays two capture strategy options with descriptions using modern card-based UI.
func (i *CaptureEventIntent) viewChooseStrategy() string {
	var content strings.Builder
	content.WriteString("\n📝 Select Capture Strategy\n\n")

	strategies := []struct {
		value       CaptureStrategy
		label       string
		description string
	}{
		{StrategyQuick, "Quick", "Capture with minimal fields (event text only)"},
		{StrategyManual, "Manual", "Full form with optional fields (date, company, project, tags)"},
	}

	for idx, s := range strategies {
		prefix := "  "
		if idx == i.state.selectedStrategyIndex {
			prefix = "▶ "
		}

		// Apply highlighting to selected item
		optStyle := lipgloss.NewStyle().Foreground(styles.ColorTextPrimary)
		if idx == i.state.selectedStrategyIndex {
			optStyle = optStyle.Foreground(styles.ColorAccentTeal).Bold(true)
		}

		line := fmt.Sprintf("%s%s - %s", prefix, s.label, s.description)
		content.WriteString(optStyle.Render(line) + "\n")
	}

	// Card styling (matching GenerateCV pattern)
	cardStyle := lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Background(styles.ColorBackgroundCard).
		Foreground(styles.ColorTextPrimary)

	return cardStyle.Render(content.String())
}

// viewCaptureForm renders the form for capturing event details.
// Delegates to the FormModel's View method to render the actual form.
func (i *CaptureEventIntent) viewCaptureForm() string {
	if i.state.captureForm == nil {
		return "Error: Form not initialized"
	}
	// Return just the form view - StandardView handles title and navigation
	return i.state.captureForm.View()
}

// viewReviewInferredEvent renders the review UI for inferred bursts and facts.
// Displays the captured event details, inferred bursts, and facts with accept/reject options.
func (i *CaptureEventIntent) viewReviewInferredEvent() string {
	// Check if editing mode is active and display appropriate modal
	switch i.state.reviewState.EditingMode {
	case EditingModeMetadata:
		return i.viewMetadataEditModal()
	case EditingModeBursts:
		return i.viewBurstEditModal()
	case EditingModeFacts:
		return i.viewFactEditModal()
	}

	// Normal review view
	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString("┌─ Review Inferred Event ────────────────────────┐\n")
	sb.WriteString("│                                                │\n")

	// Event summary
	if i.state.reviewState.Event != nil {
		title := i.state.reviewState.Event.Text
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
	// Footer now handled by StandardView
	return sb.String()
}

// viewSubmit renders the submit confirmation.
// Displays a summary of the event to be submitted with confirmation options.
func (i *CaptureEventIntent) viewSubmit() string {
	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString("┌─ Confirm Submission ───────────────────────────┐\n")
	sb.WriteString("│                                                │\n")

	if i.state.reviewState.Event != nil {
		title := i.state.reviewState.Event.Text
		if len(title) > 40 {
			title = title[:37] + "..."
		}
		sb.WriteString(fmt.Sprintf("│ Event: %s                    │\n", title))
		sb.WriteString(fmt.Sprintf("│ Date: %s                      │\n", i.state.reviewState.Event.Date))
		sb.WriteString("│                                                │\n")
		sb.WriteString(fmt.Sprintf("│ Bursts: %d                                    │\n", len(i.state.reviewState.AcceptedBursts)))
		sb.WriteString(fmt.Sprintf("│ Facts: %d                                     │\n", len(i.state.reviewState.AcceptedFacts)))
		sb.WriteString("│                                                │\n")
	}

	sb.WriteString("│ Ready to submit? Press Enter to confirm.       │\n")
	sb.WriteString("│                                                │\n")
	sb.WriteString("└────────────────────────────────────────────────┘\n")
	// Footer now handled by StandardView
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

// viewMetadataEditModal displays the metadata editor modal.
func (i *CaptureEventIntent) viewMetadataEditModal() string {
	if i.state.reviewState.metadataModal == nil {
		// Initialize metadata modal with event
		i.state.reviewState.metadataModal = models.NewMetadataEditorModelNew(
			i.state.reviewState.Event,
			i.context.CareerService,
			i.context.CLIEventService,
			context.Background(),
		)
	}
	return i.state.reviewState.metadataModal.View()
}

// viewBurstEditModal displays the burst suggestion modal.
func (i *CaptureEventIntent) viewBurstEditModal() string {
	if i.state.reviewState.burstModal == nil {
		// Convert inferred bursts to suggestions for the modal
		// For now, we'll work with an empty list - in a full implementation,
		// we'd convert i.state.reviewState.InferredBursts to suggestions
		var suggestions []burstfact.BurstSuggestion
		i.state.reviewState.burstModal = models.NewBurstSuggestionModelNew(
			i.context.CareerService,
			suggestions,
			context.Background(),
		)
	}
	return i.state.reviewState.burstModal.View()
}

// viewFactEditModal displays the fact editor modal.
func (i *CaptureEventIntent) viewFactEditModal() string {
	if i.state.reviewState.factModal == nil {
		// Use the first inferred fact, or create a new empty fact
		var fact *career.Fact
		if len(i.state.reviewState.InferredFacts) > 0 && i.state.reviewState.EditingIndex < len(i.state.reviewState.InferredFacts) {
			fact = i.state.reviewState.InferredFacts[i.state.reviewState.EditingIndex]
		} else {
			// Create a new empty fact
			fact = &career.Fact{}
		}
		i.state.reviewState.factModal = models.NewFactEditorModelNew(
			fact,
			i.context.CareerService,
			context.Background(),
		)
	}
	return i.state.reviewState.factModal.View()
}

// acceptCurrentItem accepts the currently selected burst or fact.
// Moves the item from InferredBursts/InferredFacts to AcceptedBursts/AcceptedFacts.
func (i *CaptureEventIntent) acceptCurrentItem() {
	if i.state.reviewState.SelectedItemType == "burst" {
		idx := i.state.reviewState.SelectedIndex
		if idx >= 0 && idx < len(i.state.reviewState.InferredBursts) {
			burst := i.state.reviewState.InferredBursts[idx]
			i.state.reviewState.AcceptedBursts = append(i.state.reviewState.AcceptedBursts, burst)
			// Remove from inferred list
			i.state.reviewState.InferredBursts = append(
				i.state.reviewState.InferredBursts[:idx],
				i.state.reviewState.InferredBursts[idx+1:]...)
			// Adjust selection if needed
			if i.state.reviewState.SelectedIndex >= len(i.state.reviewState.InferredBursts) && len(i.state.reviewState.InferredBursts) > 0 {
				i.state.reviewState.SelectedIndex = len(i.state.reviewState.InferredBursts) - 1
			}
		}
	} else if i.state.reviewState.SelectedItemType == "fact" {
		idx := i.state.reviewState.SelectedIndex
		if idx >= 0 && idx < len(i.state.reviewState.InferredFacts) {
			fact := i.state.reviewState.InferredFacts[idx]
			i.state.reviewState.AcceptedFacts = append(i.state.reviewState.AcceptedFacts, fact)
			// Remove from inferred list
			i.state.reviewState.InferredFacts = append(
				i.state.reviewState.InferredFacts[:idx],
				i.state.reviewState.InferredFacts[idx+1:]...)
			// Adjust selection if needed
			if i.state.reviewState.SelectedIndex >= len(i.state.reviewState.InferredFacts) && len(i.state.reviewState.InferredFacts) > 0 {
				i.state.reviewState.SelectedIndex = len(i.state.reviewState.InferredFacts) - 1
			}
		}
	}
}

// rejectCurrentItem rejects the currently selected burst or fact.
// Removes the item from InferredBursts/InferredFacts without adding to accepted.
func (i *CaptureEventIntent) rejectCurrentItem() {
	if i.state.reviewState.SelectedItemType == "burst" {
		idx := i.state.reviewState.SelectedIndex
		if idx >= 0 && idx < len(i.state.reviewState.InferredBursts) {
			burst := i.state.reviewState.InferredBursts[idx]
			// Track rejection reason (optional - could add a modal for this)
			if i.state.reviewState.RejectedItems == nil {
				i.state.reviewState.RejectedItems = make(map[string]string)
			}
			i.state.reviewState.RejectedItems[burst.ID] = "user_rejected"
			// Remove from inferred list
			i.state.reviewState.InferredBursts = append(
				i.state.reviewState.InferredBursts[:idx],
				i.state.reviewState.InferredBursts[idx+1:]...)
			// Adjust selection if needed
			if i.state.reviewState.SelectedIndex >= len(i.state.reviewState.InferredBursts) && len(i.state.reviewState.InferredBursts) > 0 {
				i.state.reviewState.SelectedIndex = len(i.state.reviewState.InferredBursts) - 1
			}
		}
	} else if i.state.reviewState.SelectedItemType == "fact" {
		idx := i.state.reviewState.SelectedIndex
		if idx >= 0 && idx < len(i.state.reviewState.InferredFacts) {
			fact := i.state.reviewState.InferredFacts[idx]
			// Track rejection reason
			if i.state.reviewState.RejectedItems == nil {
				i.state.reviewState.RejectedItems = make(map[string]string)
			}
			i.state.reviewState.RejectedItems[fact.ID] = "user_rejected"
			// Remove from inferred list
			i.state.reviewState.InferredFacts = append(
				i.state.reviewState.InferredFacts[:idx],
				i.state.reviewState.InferredFacts[idx+1:]...)
			// Adjust selection if needed
			if i.state.reviewState.SelectedIndex >= len(i.state.reviewState.InferredFacts) && len(i.state.reviewState.InferredFacts) > 0 {
				i.state.reviewState.SelectedIndex = len(i.state.reviewState.InferredFacts) - 1
			}
		}
	}
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
func (i *CaptureEventIntent) GetState() string {
	return string(i.state.currentState)
}

// GetForm returns the current form model instance (for test and debug)
func (i *CaptureEventIntent) GetForm() *models.FormModel {
	if i == nil || i.state == nil {
		return nil
	}
	return i.state.captureForm
}

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
