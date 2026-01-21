package intents

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/cli/screens"
	captureScreens "github.com/baphled/kariya/internal/cli/screens/capture"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
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

// Ensure CaptureEventIntent implements ScreenResultHandler interface
var _ ScreenResultHandler = (*CaptureEventIntent)(nil)

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

	// Screen orchestration (NEW - screens architecture)
	activeScreen screens.Screen // Currently active screen
	useScreens   bool           // Toggle between screens and legacy code
}

// NewCaptureEventIntent creates a new CaptureEvent intent.
func NewCaptureEventIntent(context *CaptureEventContext) (*CaptureEventIntent, error) {
	// Validate the context.
	if err := context.Validate(); err != nil {
		return nil, err
	}

	// Create the form model for capturing event details
	formModel := models.NewCaptureForm(context.CLIEventService)

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
		active:     true,
		useScreens: true, // Enable screens architecture
	}, nil
}

// DisableScreens disables the screens architecture for testing legacy flows.
// This is primarily used in tests that need to test legacy message handling.
func (i *CaptureEventIntent) DisableScreens() {
	i.useScreens = false
	i.activeScreen = nil
}

// Init is called when the intent is activated.
func (i *CaptureEventIntent) Init() tea.Cmd {
	// NEW: Check if screens architecture is enabled
	if i.useScreens {
		// Create breadcrumbs for navigation
		breadcrumbs := []string{"Main Menu", "Capture Event"}

		// Create the initial screen (StrategySelectScreen)
		i.activeScreen = captureScreens.NewStrategySelectScreen(breadcrumbs)

		// Get terminal info
		termInfo := i.GetTerminalInfo()
		width, height := 120, 40 // defaults
		if termInfo != nil {
			width = termInfo.Width
			height = termInfo.Height
		}

		// Pass terminal info, theme, and logo to screen
		i.activeScreen.SetTerminalInfo(width, height)
		i.activeScreen.SetTheme(i.Theme())
		i.activeScreen.SetLogo(i.GetLogo(), i.GetLogoSpacing())

		// Initialize the screen if it has an Init method (needed for forms)
		// This is critical for screens that wrap forms - without Init(),
		// the underlying form won't be able to accept input
		if initable, ok := i.activeScreen.(interface{ Init() tea.Cmd }); ok {
			return initable.Init()
		}
		return nil
	}

	// Otherwise, initialize for a new event.
	i.initializeFormForNew()
	return func() tea.Msg { return nil }
}

// Theme helper methods for consistent themed styling.

// getCardStyle returns a themed card style, with fallback to default styling.
// getTheme returns the theme or a default.
func (i *CaptureEventIntent) getTheme() themes.Theme {
	if theme := i.Theme(); theme != nil {
		return theme
	}
	return themes.NewDefaultTheme()
}

func (i *CaptureEventIntent) getCardStyle() lipgloss.Style {
	return i.getTheme().Styles().CardBase
}

// getPrimaryColor returns the primary text color from theme.
func (i *CaptureEventIntent) getPrimaryColor() lipgloss.Color {
	return i.getTheme().ForegroundColor()
}

// getAccentColor returns the accent color from theme.
func (i *CaptureEventIntent) getAccentColor() lipgloss.Color {
	return i.getTheme().PrimaryColor()
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

	// Handle async submission messages (modal overlay pattern)
	// These messages must be handled BEFORE delegating to screens
	switch msg := msg.(type) {
	case FormSubmittedMsg:
		// Handle form submission - works with both screens and legacy mode
		if msg.Event == nil {
			i.setFailed("INVALID_FORM", "Form submission with nil event", nil)
			return nil
		}

		// Validate the event data
		if err := msg.Event.Validate(); err != nil {
			i.setFailed("VALIDATION_ERROR", fmt.Sprintf("Form validation failed: %v", err), err)
			return nil
		}

		// Initialize review state with the event
		i.state.reviewState = &ReviewInferredEventState{
			Event:          msg.Event,
			InferredBursts: make([]*career.Burst, 0),
			InferredFacts:  make([]*career.Fact, 0),
			AcceptedBursts: make([]*career.Burst, 0),
			AcceptedFacts:  make([]*career.Fact, 0),
			RejectedItems:  make(map[string]string),
		}
		i.state.currentState = CaptureStateReview

		// Clear activeScreen so legacy review View() is used for pre-save review
		// Post-save review will create EventReviewScreen in DismissModalMsg handler
		i.activeScreen = nil
		return nil

	case SubmitCompleteMsg:
		// Submission succeeded - show success modal briefly, then complete intent
		i.state.submitModal = feedback.NewSuccessModal("Event saved!")
		// Auto-dismiss after 2 seconds
		return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
			return DismissModalMsg{}
		})

	case SubmitErrorMsg:
		// Submission failed - show error modal (user can press Esc to dismiss)
		i.state.submitModal = feedback.NewErrorModal("Save Failed", msg.Message)
		return nil

	case DismissModalMsg:
		// Modal auto-dismissed after success - return to Review for enrichment review
		if i.state.submitModal != nil {
			i.state.submitModal = nil
			// Instead of completing, return to Review state
			// This allows user to review inferred bursts/facts after save
			i.state.postSaveReview = true // Mark as post-save review
			i.state.currentState = CaptureStateReview

			// CRITICAL: If using screens architecture, create EventReviewScreen
			if i.useScreens {
				breadcrumbs := []string{"Main Menu", "Capture Event", "Review Enrichment"}
				i.activeScreen = captureScreens.NewEventReviewScreen(
					breadcrumbs,
					i.state.reviewState.Event,
					i.state.reviewState.InferredBursts,
					i.state.reviewState.InferredFacts,
				)

				// Get terminal info
				termInfo := i.GetTerminalInfo()
				width, height := 120, 40 // defaults
				if termInfo != nil {
					width = termInfo.Width
					height = termInfo.Height
				}

				// Pass terminal info, theme, and logo to screen
				i.activeScreen.SetTerminalInfo(width, height)
				i.activeScreen.SetTheme(i.Theme())
				i.activeScreen.SetLogo(i.GetLogo(), i.GetLogoSpacing())
			}
		}
		return nil
	}

	// Handle Esc key when error modal is showing (allows dismissing error and returning to form)
	if i.state.submitModal != nil {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			if keyMsg.String() == "esc" {
				// Dismiss error modal and stay in current state
				i.state.submitModal = nil
				return nil
			}
		}
	}

	// Handle global keys BEFORE delegating to screen
	// This ensures q (quit), ? (help) are always processed first
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch HandleGlobalKeys(keyMsg) {
		case KeyQuit:
			// 'q' pressed - quit the application
			return tea.Quit
		case KeyHelp:
			// '?' pressed - toggle help modal
			i.helpModal.Toggle()
			return nil
		}
		// KeyBack (esc) is handled by the screen as CancelResult or modal dismissal above
	}

	// NEW: Check if screens architecture is enabled
	if i.useScreens && i.activeScreen != nil {
		// Check if editing modal is active - delegate to modal first
		if i.state.reviewState != nil && i.state.reviewState.EditingMode != EditingModeNone {
			return i.updateEditingModal(msg)
		}

		// Delegate to active screen
		cmd, result := i.activeScreen.Update(msg)

		// If screen returned a result, handle it
		if result != nil {
			// Delegate result handling to ScreenResultHandler methods
			return i.handleScreenResult(result)
		}

		// Otherwise return the command from screen
		return cmd
	}

	// LEGACY: Fall back to old state machine
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

			// Note: For quick mode, date will be set to today automatically in the submit handler

			i.state.currentState = CaptureStateForm

			// CRITICAL: Initialize the form so it can accept input
			return i.state.captureForm.Init()

		}

		// Handle global keys (q=quit, ?=help, esc=back)
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			i.ToggleHelp()
			return nil
		case KeyBack:
			// At root state, back means cancel and return to main menu
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
	// Check for special messages that indicate form completion or navigation
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle global keys FIRST (before form processes them)
		// This ensures esc, q, ?, m keys work even when form has focus
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			i.ToggleHelp()
			return nil
		case KeyBack:
			// Determine where to go back based on context
			if i.context.PreviousEvent != nil {
				// Editing existing event - cancel and return to caller (e.g., BrowseTimeline)
				i.setCancelled()
				return nil
			}
			// New event capture - go back to strategy selection
			i.state.currentState = CaptureStateChooseStrategy
			return nil
		}

		switch msg.String() {
		case "ctrl+s":
			// User pressed Ctrl+S to submit the form
			// Trigger form submission
			return i.state.captureForm.SubmitForm()
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

	// Delegate all messages to the form model to handle input and state
	// This happens AFTER global keys are checked, so form doesn't consume them
	_, formCmd := i.state.captureForm.Update(msg)

	// Return the command from the form update
	return formCmd

}

// updateReviewInferredEvent handles messages while reviewing inferred bursts and facts.
// It processes review confirmations, edits, and transitions to submit state.
func (i *CaptureEventIntent) updateReviewInferredEvent(msg tea.Msg) tea.Cmd {
	// Check global keys BEFORE routing to modals
	// This ensures esc, q, ?, m keys work even when modal has focus
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			i.ToggleHelp()
			return nil
		case KeyBack:
			// If modal is active, close it
			if i.state.reviewState.EditingMode != EditingModeNone {
				i.state.reviewState.metadataModal = nil
				i.state.reviewState.burstModal = nil
				i.state.reviewState.factModal = nil
				i.state.reviewState.EditingMode = EditingModeNone
				return nil
			}
			// Otherwise go back to form
			i.state.currentState = CaptureStateForm
			return nil
		}
	}

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
			// Check if this is post-save review or pre-save review
			if i.state.postSaveReview {
				// Post-save review: user is done reviewing enriched data, complete intent
				result := &CaptureEventResult{
					Event:          i.state.reviewState.Event,
					Bursts:         i.state.reviewState.InferredBursts,
					Facts:          i.state.reviewState.InferredFacts,
					AcceptedFields: make(map[string]bool),
					RejectedFields: i.state.reviewState.RejectedItems,
				}
				i.setCompleted(result)
				return nil
			}
			// Pre-save review: proceed to submit
			i.state.currentState = CaptureStateSubmit
			return i.performSubmit()

		case "e":
			// Edit metadata (modal sub-flow)
			i.state.reviewState.EditingMode = EditingModeMetadata
			// Initialize modal if not already created
			if i.state.reviewState.metadataModal == nil && i.state.reviewState.Event != nil {
				i.state.reviewState.metadataModal = models.NewMetadataEditorModelNew(
					i.state.reviewState.Event,
					i.context.CareerService,
					i.context.CLIEventService,
					context.Background(),
				)
				return i.state.reviewState.metadataModal.Init()
			}
			return nil

		case "b":
			// Edit bursts (modal sub-flow)
			i.state.reviewState.EditingMode = EditingModeBursts
			// Initialize modal if not already created
			if i.state.reviewState.burstModal == nil {
				var suggestions []burstfact.BurstSuggestion
				i.state.reviewState.burstModal = models.NewBurstSuggestionModelNew(
					i.context.CareerService,
					suggestions,
					context.Background(),
				)
				return i.state.reviewState.burstModal.Init()
			}
			return nil

		case "f":
			// Edit facts (modal sub-flow)
			i.state.reviewState.EditingMode = EditingModeFacts
			// Initialize modal if not already created
			if i.state.reviewState.factModal == nil {
				// Create with an empty fact for now - in a real implementation,
				// we'd pass the selected fact for editing
				i.state.reviewState.factModal = models.NewFactEditorModelNew(
					&career.Fact{},
					i.context.CareerService,
					context.Background(),
				)
				return i.state.reviewState.factModal.Init()
			}
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
		// Handle global keys first (q=quit, ?=help, esc=back)
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			i.ToggleHelp()
			return nil
		case KeyBack:
			// Go back to review state (keep error visible per user preference)
			i.state.currentState = CaptureStateReview
			return nil
		}

		switch msg.String() {
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

		// CRITICAL: Use CareerService directly (not CLIEventService wrapper)
		// CareerService.CaptureEvent modifies the event in-place, setting its ID
		// CLIEventService.CaptureEvent creates a new event internally, leaving our event without an ID
		if i.state.context.CareerService == nil {
			return SubmitErrorMsg{
				Code:    "SERVICE_ERROR",
				Message: "Career service not initialized",
				Cause:   nil,
			}
		}

		// Call the service to capture the event
		// The service handles persistence and populates event.ID
		err := i.state.context.CareerService.CaptureEvent(ctx, event, mode)

		if err != nil {
			// Map service errors to intent errors
			return SubmitErrorMsg{
				Code:    "PERSISTENCE_ERROR",
				Message: fmt.Sprintf("Failed to save event: %v", err),
				Cause:   err,
			}
		}

		// Perform enrichment for all strategies (if CareerService is available)
		// This extracts bursts and facts from the saved event
		if i.state.context.CareerService != nil {
			// Enrichment error is logged internally but doesn't fail submission
			// The event is already saved successfully - enrichment is optional
			//nolint:errcheck // Intentional: enrichment failure should not block event submission
			i.performEnrichment(ctx, event)
		}

		// Save any accepted facts from review that might have been manually edited/added
		// Note: Facts from enrichment are already saved in performEnrichment()
		// This is a safety check for any facts that might have been added during review
		if i.state.context.CareerService != nil && len(i.state.reviewState.AcceptedFacts) > 0 {
			for _, fact := range i.state.reviewState.AcceptedFacts {
				// Only save facts that don't have an ID yet (haven't been saved)
				// Facts from enrichment already have IDs
				if fact.ID == "" {
					// Ensure fact is linked to the saved event
					if fact.SourceEventID == "" {
						fact.SourceEventID = event.ID
					}

					// Save the fact
					if err := i.state.context.CareerService.SaveFact(ctx, fact); err != nil {
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
func (i *CaptureEventIntent) performEnrichment(ctx context.Context, event *career.CareerEvent) error {
	if i.state.context.CareerService == nil {
		return fmt.Errorf("career service not available for enrichment")
	}

	// Suggest bursts for the event
	burstSuggestions, err := i.state.context.CareerService.SuggestBursts(ctx, []string{event.ID})
	if err == nil && len(burstSuggestions) > 0 {
		// Save burst suggestions and store them for review
		bursts, err := i.state.context.CareerService.SaveBurstSuggestions(ctx, burstSuggestions)
		if err == nil && len(bursts) > 0 {
			i.state.reviewState.InferredBursts = bursts
		}
	}

	// Extract facts from the event
	facts, err := i.state.context.CareerService.ExtractFactsFromEvent(ctx, event)
	if err == nil && len(facts) > 0 {
		// Persist each extracted fact to the database
		for j := range facts {
			fact := &facts[j]

			// Set source event ID (linking fact to this event)
			fact.SourceEventID = event.ID

			// Save fact to repository
			if err := i.state.context.CareerService.SaveFact(ctx, fact); err != nil {
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
	theme := i.Theme()

	switch i.state.currentState {
	case CaptureStateChooseStrategy:
		return CombineThemedFooters(
			ThemedNavigationFooter(theme),
			ThemedGlobalBadges(theme),
		)
	case CaptureStateForm:
		// Show different help based on strategy
		if i.state.strategy == StrategyManual {
			return CombineThemedFooters(
				ThemedFormFooter(theme),
				ThemedCustomFooter(theme,
					primitives.HelpKeyBadge("Ctrl+O", "Toggle fields", theme),
				),
				ThemedGlobalBadges(theme),
			)
		}
		return CombineThemedFooters(
			ThemedFormFooter(theme),
			ThemedGlobalBadges(theme),
		)
	case CaptureStateReview:
		// Show different help when modal is active
		if i.state.reviewState.EditingMode != EditingModeNone {
			return ThemedCustomFooter(theme,
				primitives.HelpKeyBadge("Editing", "...", theme),
				primitives.CancelBadge(theme),
				primitives.SaveBadge(theme),
			)
		}
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				primitives.NavigateBadge(theme),
				primitives.EditBadge(theme),
				primitives.HelpKeyBadge("b", "Bursts", theme),
				primitives.HelpKeyBadge("f", "Facts", theme),
				primitives.HelpKeyBadge("a", "Accept", theme),
				primitives.HelpKeyBadge("r", "Reject", theme),
				primitives.BackBadge(theme),
			),
			ThemedGlobalBadges(theme),
		)
	case CaptureStateSubmit:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				primitives.HelpKeyBadge("Enter", "Continue", theme),
				primitives.BackBadge(theme),
			),
			ThemedGlobalBadges(theme),
		)
	default:
		return ThemedGlobalBadges(theme)
	}
}

// View renders the intent's current state using StandardView.
func (i *CaptureEventIntent) View() string {
	if !i.active {
		return "CaptureEvent intent is not active"
	}

	// NEW: Check if screens architecture is enabled
	if i.useScreens && i.activeScreen != nil {
		// Get base view from active screen
		baseView := i.activeScreen.View()

		// Get terminal dimensions
		termInfo := i.GetTerminalInfo()
		width, height := 80, 24 // defaults
		if termInfo != nil {
			width = termInfo.Width
			height = termInfo.Height
		}

		// If submit modal is visible, overlay it on the base view
		if i.state.submitModal != nil {
			// Render modal content
			modalContent := i.state.submitModal.Render(width, height)

			// Overlay modal on background (centered)
			return i.overlayModal(baseView, modalContent, width, height)
		}

		// Check if editing modal is active (for review state)
		if i.state.reviewState != nil && i.state.reviewState.EditingMode != EditingModeNone {
			modalContent := i.getEditingModalContent()
			if modalContent != nil {
				return i.renderModalOverlay(baseView, modalContent)
			}
		}

		// No modal - return base view
		return baseView
	}

	// LEGACY: Fall back to old view rendering
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

		// Apply themed highlighting to selected item
		optStyle := lipgloss.NewStyle().Foreground(i.getPrimaryColor())
		if idx == i.state.selectedStrategyIndex {
			optStyle = optStyle.Foreground(i.getAccentColor()).Bold(true)
		}

		line := fmt.Sprintf("%s%s - %s", prefix, s.label, s.description)
		content.WriteString(optStyle.Render(line) + "\n")
	}

	// Apply themed card styling
	return i.getCardStyle().Render(content.String())
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
// When in editing mode, the modal is rendered as an overlay on top of the review content.
func (i *CaptureEventIntent) viewReviewInferredEvent() string {
	// Build the base review view
	baseView := i.buildReviewBaseView()

	// Check if editing mode is active and render modal as overlay
	switch i.state.reviewState.EditingMode {
	case EditingModeMetadata:
		return i.renderModalOverlay(baseView, i.getMetadataModalContent())
	case EditingModeBursts:
		return i.renderModalOverlay(baseView, i.getBurstModalContent())
	case EditingModeFacts:
		return i.renderModalOverlay(baseView, i.getFactModalContent())
	}

	return baseView
}

// buildReviewBaseView builds the base review view content.
func (i *CaptureEventIntent) buildReviewBaseView() string {
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

	// Inferred bursts - show both inferred and accepted
	sb.WriteString("│ Inferred Bursts:                               │\n")
	bursts := i.state.reviewState.InferredBursts
	if len(bursts) == 0 {
		bursts = i.state.reviewState.AcceptedBursts
	}
	if len(bursts) > 0 {
		for idx, burst := range bursts {
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

	// Inferred facts - show both inferred and accepted
	sb.WriteString("│ Inferred Facts:                                │\n")
	facts := i.state.reviewState.InferredFacts
	if len(facts) == 0 {
		facts = i.state.reviewState.AcceptedFacts
	}
	if len(facts) > 0 {
		for idx, fact := range facts {
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

// renderModalOverlay renders a modal overlay on top of the background content.
func (i *CaptureEventIntent) renderModalOverlay(background string, modalContent *modalContentData) string {
	if modalContent == nil {
		return background
	}

	// Get terminal dimensions
	info := i.GetTerminalInfo()
	width := 80
	height := 24
	if info != nil {
		width = info.Width
		height = info.Height
	}

	// Create overlay modal
	overlay := feedback.NewOverlayModal(modalContent.title, modalContent.content)
	overlay.SetFooter(modalContent.footer)
	overlay.SetWidth(80) // Use a standard modal width

	return overlay.RenderCentered(background, width, height)
}

// modalContentData holds modal content for overlay rendering.
type modalContentData struct {
	title   string
	content string
	footer  string
}

// getMetadataModalContent returns the modal content for metadata editing.
func (i *CaptureEventIntent) getMetadataModalContent() *modalContentData {
	if i.state.reviewState.metadataModal == nil {
		// Initialize metadata modal with event
		i.state.reviewState.metadataModal = models.NewMetadataEditorModelNew(
			i.state.reviewState.Event,
			i.context.CareerService,
			i.context.CLIEventService,
			context.Background(),
		)
	}
	return &modalContentData{
		title:   i.state.reviewState.metadataModal.GetTitle(),
		content: i.state.reviewState.metadataModal.GetContent(),
		footer:  i.state.reviewState.metadataModal.GetFooter(),
	}
}

// getBurstModalContent returns the modal content for burst editing.
func (i *CaptureEventIntent) getBurstModalContent() *modalContentData {
	if i.state.reviewState.burstModal == nil {
		// Convert inferred bursts to suggestions for the modal
		var suggestions []burstfact.BurstSuggestion
		i.state.reviewState.burstModal = models.NewBurstSuggestionModelNew(
			i.context.CareerService,
			suggestions,
			context.Background(),
		)
	}
	return &modalContentData{
		title:   i.state.reviewState.burstModal.GetTitle(),
		content: i.state.reviewState.burstModal.GetContent(),
		footer:  i.state.reviewState.burstModal.GetFooter(),
	}
}

// getFactModalContent returns the modal content for fact editing.
func (i *CaptureEventIntent) getFactModalContent() *modalContentData {
	if i.state.reviewState.factModal == nil {
		// Use the first inferred fact, or create a new empty fact
		var fact *career.Fact
		if len(i.state.reviewState.InferredFacts) > 0 && i.state.reviewState.EditingIndex < len(i.state.reviewState.InferredFacts) {
			fact = i.state.reviewState.InferredFacts[i.state.reviewState.EditingIndex]
		} else {
			fact = &career.Fact{}
		}
		i.state.reviewState.factModal = models.NewFactEditorModelNew(
			fact,
			i.context.CareerService,
			context.Background(),
		)
	}
	return &modalContentData{
		title:   i.state.reviewState.factModal.GetTitle(),
		content: i.state.reviewState.factModal.GetContent(),
		footer:  i.state.reviewState.factModal.GetFooter(),
	}
}

// getEditingModalContent returns the modal content based on current editing mode.
// Used by screens architecture to render editing modals as overlays.
func (i *CaptureEventIntent) getEditingModalContent() *modalContentData {
	if i.state.reviewState == nil {
		return nil
	}
	switch i.state.reviewState.EditingMode {
	case EditingModeMetadata:
		return i.getMetadataModalContent()
	case EditingModeBursts:
		return i.getBurstModalContent()
	case EditingModeFacts:
		return i.getFactModalContent()
	default:
		return nil
	}
}

// updateEditingModal handles modal updates when using screens architecture.
// This is called when an editing modal (metadata, bursts, facts) is active.
func (i *CaptureEventIntent) updateEditingModal(msg tea.Msg) tea.Cmd {
	// Handle escape key to close modal
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.String() == "esc" {
			// Close modal and return to review
			i.state.reviewState.metadataModal = nil
			i.state.reviewState.burstModal = nil
			i.state.reviewState.factModal = nil
			i.state.reviewState.EditingMode = EditingModeNone
			return nil
		}
	}

	// Delegate to the active modal
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
			// TODO: Add completion check when BurstSuggestionModelNew supports it
			return cmd
		}

	case EditingModeFacts:
		if i.state.reviewState.factModal != nil {
			modal, cmd := i.state.reviewState.factModal.Update(msg)
			i.state.reviewState.factModal = modal.(*models.FactEditorModelNew)

			// Check if modal completed
			if i.state.reviewState.factModal.IsSubmitted() {
				// Apply changes to fact
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

	return nil
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
func (i *CaptureEventIntent) GetForm() *models.CaptureForm {
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

// ============================================================================
// ScreenResultHandler Interface Implementation (NEW - screens architecture)
// ============================================================================

// handleScreenResult processes results from screen updates.
// This is the central hub for all screen-to-intent communication.
//
// Uses ScreenResultDispatcher pattern to eliminate repetitive type switching.
// CaptureEventIntent implements ScreenResultHandler interface for compile-time safety.
func (i *CaptureEventIntent) handleScreenResult(result screens.ScreenResult) tea.Cmd {
	return NewScreenResultDispatcher(i).Dispatch(result)
}

// HandleNavigate handles navigation actions from screens.
//
// Implements ScreenResultHandler interface.
func (i *CaptureEventIntent) HandleNavigate(result *screens.NavigateResult) tea.Cmd {
	data := result.Data()

	// Check if data is a string (simple action)
	if action, ok := data.(string); ok {
		switch action {
		case "edit_metadata":
			// User wants to edit event metadata
			// Check if we have an event to edit
			if i.state.reviewState == nil || i.state.reviewState.Event == nil {
				return i.setFailedCmd("NO_EVENT", "No event to edit", nil)
			}
			// Create metadata modal and set editing mode
			i.state.reviewState.EditingMode = EditingModeMetadata
			i.state.reviewState.metadataModal = models.NewMetadataEditorModelNew(
				i.state.reviewState.Event,
				i.context.CareerService,
				i.context.CLIEventService,
				context.Background(),
			)
			// Initialize the modal form
			return i.state.reviewState.metadataModal.Init()

		case "edit_bursts":
			// User wants to edit bursts
			i.state.reviewState.EditingMode = EditingModeBursts
			// Create burst suggestion modal with inferred bursts
			var suggestions []burstfact.BurstSuggestion
			for _, b := range i.state.reviewState.InferredBursts {
				suggestions = append(suggestions, burstfact.BurstSuggestion{
					Name:        b.Name,
					Description: b.Description,
				})
			}
			i.state.reviewState.burstModal = models.NewBurstSuggestionModelNew(
				i.context.CareerService,
				suggestions,
				context.Background(),
			)
			return i.state.reviewState.burstModal.Init()

		case "edit_facts":
			// User wants to edit facts
			i.state.reviewState.EditingMode = EditingModeFacts
			// Create fact editor modal with first inferred fact (or new)
			var fact *career.Fact
			if len(i.state.reviewState.InferredFacts) > 0 {
				fact = i.state.reviewState.InferredFacts[0]
			} else {
				fact = &career.Fact{Text: ""}
			}
			i.state.reviewState.factModal = models.NewFactEditorModelNew(
				fact,
				i.context.CareerService,
				context.Background(),
			)
			return i.state.reviewState.factModal.Init()

		default:
			return i.setFailedCmd("INVALID_NAVIGATION", fmt.Sprintf("Unknown navigation action: %s", action), nil)
		}
	}

	// Check if data is CaptureStrategy (from strategy selection)
	if strategy, ok := data.(CaptureStrategy); ok {
		// User selected a strategy - transition to form screen
		i.state.strategy = strategy
		return i.transitionToFormScreen(strategy)
	}

	return i.setFailedCmd("INVALID_NAVIGATION_DATA", fmt.Sprintf("Invalid navigation data type: %T", data), nil)
}

// HandleCancel handles cancellation from screens.
//
// Implements ScreenResultHandler interface.
func (i *CaptureEventIntent) HandleCancel(result *screens.CancelResult) tea.Cmd {
	// Determine which screen we're cancelling from based on current state
	switch i.state.currentState {
	case CaptureStateChooseStrategy:
		// Root state - cancel the entire intent
		i.setCancelled()
		return nil

	case CaptureStateForm:
		// Check if we're in edit mode (editing existing event from another intent like BrowseTimeline)
		if i.state.context.PreviousEvent != nil {
			// Edit mode - cancel the entire intent and return to caller
			i.setCancelled()
			return nil
		}
		// New event mode - go back to strategy selection
		return i.transitionToStrategyScreen()

	case CaptureStateReview:
		// Cancel review - go back to form
		return i.transitionToFormScreen(i.state.strategy)

	case CaptureStateSubmit:
		// Cannot cancel during submission
		return nil

	default:
		return i.setFailedCmd("INVALID_CANCEL_STATE", fmt.Sprintf("Cannot cancel from state: %s", i.state.currentState), nil)
	}
}

// HandleSubmit handles form/data submission from screens.
//
// Implements ScreenResultHandler interface.
func (i *CaptureEventIntent) HandleSubmit(result *screens.SubmitResult) tea.Cmd {
	data := result.Data()

	switch i.state.currentState {
	case CaptureStateForm:
		// Form submitted with event data
		if event, ok := data.(*career.CareerEvent); ok {
			// Validate event
			if err := event.Validate(); err != nil {
				return i.setFailedCmd("VALIDATION_ERROR", fmt.Sprintf("Event validation failed: %v", err), err)
			}

			// Store event and transition to review (if enrichment enabled)
			// For now, go directly to submit with modal overlay
			i.state.reviewState = &ReviewInferredEventState{
				Event:          event,
				InferredBursts: make([]*career.Burst, 0),
				InferredFacts:  make([]*career.Fact, 0),
				AcceptedBursts: make([]*career.Burst, 0),
				AcceptedFacts:  make([]*career.Fact, 0),
				RejectedItems:  make(map[string]string),
			}

			// Show loading modal and perform async submit
			i.state.submitModal = feedback.NewLoadingModal("Saving event...", false)
			return i.performSubmit()
		}
		return i.setFailedCmd("INVALID_FORM_DATA", fmt.Sprintf("Invalid form data type: %T", data), nil)

	case CaptureStateReview:
		// Review confirmed - extract event, bursts, facts from result
		if reviewData, ok := data.(map[string]interface{}); ok {
			event, _ := reviewData["event"].(*career.CareerEvent)
			bursts, _ := reviewData["bursts"].([]*career.Burst)
			facts, _ := reviewData["facts"].([]*career.Fact)

			// Update review state
			i.state.reviewState.Event = event
			i.state.reviewState.AcceptedBursts = bursts
			i.state.reviewState.AcceptedFacts = facts

			// Show loading modal and perform async submit
			i.state.submitModal = feedback.NewLoadingModal("Saving event...", false)
			return i.performSubmit()
		}
		return i.setFailedCmd("INVALID_REVIEW_DATA", fmt.Sprintf("Invalid review data type: %T", data), nil)

	case CaptureStateSubmit:
		// Submission complete - extract results
		if submitData, ok := data.(map[string]interface{}); ok {
			event, _ := submitData["event"].(*career.CareerEvent)
			bursts, _ := submitData["bursts"].([]*career.Burst)
			facts, _ := submitData["facts"].([]*career.Fact)

			// Create successful result
			i.result = &IntentResult[*CaptureEventResult]{
				Status: Completed,
				Data: &CaptureEventResult{
					Event:  event,
					Bursts: bursts,
					Facts:  facts,
				},
			}
			i.active = false
			return nil
		}
		return i.setFailedCmd("INVALID_SUBMIT_DATA", fmt.Sprintf("Invalid submit data type: %T", data), nil)

	default:
		return i.setFailedCmd("INVALID_SUBMIT_STATE", fmt.Sprintf("Cannot submit from state: %s", i.state.currentState), nil)
	}
}

// HandleError handles errors from screens.
//
// Implements ScreenResultHandler interface.
func (i *CaptureEventIntent) HandleError(result *screens.ErrorResult) tea.Cmd {
	// Screen encountered an error - propagate to intent
	data := result.Data()
	if errorData, ok := data.(map[string]interface{}); ok {
		err, _ := errorData["error"].(error)
		msg, _ := errorData["message"].(string)
		return i.setFailedCmd("SCREEN_ERROR", msg, err)
	}
	return i.setFailedCmd("SCREEN_ERROR", "Unknown screen error", fmt.Errorf("%v", data))
}

// setFailedCmd is a helper to set failed result and return nil command.
func (i *CaptureEventIntent) setFailedCmd(code, message string, cause error) tea.Cmd {
	i.setFailed(code, message, cause)
	return nil
}

// ============================================================================
// Screen Transition Helpers (NEW - screens architecture)
// ============================================================================

// transitionToStrategyScreen transitions to the strategy selection screen.
func (i *CaptureEventIntent) transitionToStrategyScreen() tea.Cmd {
	// Update intent state
	i.state.currentState = CaptureStateChooseStrategy

	// Create breadcrumbs
	breadcrumbs := []string{"Main Menu", "Capture Event"}

	// Create strategy selection screen
	i.activeScreen = captureScreens.NewStrategySelectScreen(breadcrumbs)

	// Get terminal info
	termInfo := i.GetTerminalInfo()
	width, height := 120, 40 // defaults
	if termInfo != nil {
		width = termInfo.Width
		height = termInfo.Height
	}

	// Pass context to screen
	i.activeScreen.SetTerminalInfo(width, height)
	i.activeScreen.SetTheme(i.Theme())
	i.activeScreen.SetLogo(i.GetLogo(), i.GetLogoSpacing())

	return nil
}

// transitionToFormScreen transitions to the event form screen.
func (i *CaptureEventIntent) transitionToFormScreen(strategy CaptureStrategy) tea.Cmd {
	// Update intent state
	i.state.currentState = CaptureStateForm
	i.state.strategy = strategy

	// Create breadcrumbs
	breadcrumbs := []string{"Main Menu", "Capture Event", "Form"}

	// Create event form screen
	i.activeScreen = captureScreens.NewEventFormScreen(
		i.eventService,
		breadcrumbs,
		strategy,
	)

	// Get terminal info
	termInfo := i.GetTerminalInfo()
	width, height := 120, 40 // defaults
	if termInfo != nil {
		width = termInfo.Width
		height = termInfo.Height
	}

	// Pass context to screen
	i.activeScreen.SetTerminalInfo(width, height)
	i.activeScreen.SetTheme(i.Theme())
	i.activeScreen.SetLogo(i.GetLogo(), i.GetLogoSpacing())

	// Initialize the screen if it has an Init method (critical for forms!)
	// Without this, the underlying huh form won't accept input
	if initable, ok := i.activeScreen.(interface{ Init() tea.Cmd }); ok {
		return initable.Init()
	}
	return nil
}

// TODO: Implement transitionToReviewScreen when review step is enabled.
// Currently form goes directly to submit (see HandleSubmit line 1444).
// Will use captureScreens.NewEventReviewScreen() when implemented.

// ============================================================================
// Modal Overlay Rendering
// ============================================================================

// overlayModal overlays modal content on top of background content (centered).
// This follows the StandardView modal overlay pattern for consistent modal rendering.
func (i *CaptureEventIntent) overlayModal(background, modal string, width, height int) string {
	bgLines := strings.Split(background, "\n")
	modalLines := strings.Split(modal, "\n")

	// Calculate vertical position to center modal
	bgHeight := len(bgLines)
	modalHeight := len(modalLines)
	startLine := (bgHeight - modalHeight) / 2
	if startLine < 0 {
		startLine = 0
	}

	// Overlay modal lines onto background
	result := make([]string, len(bgLines))
	copy(result, bgLines)

	for i, modalLine := range modalLines {
		lineIndex := startLine + i
		if lineIndex >= 0 && lineIndex < len(result) {
			// Center modal line horizontally
			centeredModalLine := lipgloss.PlaceHorizontal(width, lipgloss.Center, modalLine)
			result[lineIndex] = centeredModalLine
		}
	}

	return strings.Join(result, "\n")
}
