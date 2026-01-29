package intents

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/cli/screens"
	captureScreens "github.com/baphled/kariya/internal/cli/screens/capture"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/domain/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	burstfact "github.com/baphled/kariya/internal/service/career/burstfact"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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

// Ensure CaptureEventIntent implements ScreenResultHandler interface
var _ behaviors.ScreenResultHandler = (*CaptureEventIntent)(nil)

// CaptureEventIntent implements the Intent interface for capturing career events.
// It owns the complete lifecycle of event capture, including:
// - Choosing capture strategy (manual, quick, enriched)
// - Capturing event details via form
// - Reviewing and refining inferred bursts and facts
// - Submitting the final event
type CaptureEventIntent struct {
	*BaseIntent

	context      *CaptureEventContext
	state        *CaptureEventModel
	active       bool
	result       *IntentResult[*CaptureEventResult]
	eventService *service.CLIEventService

	activeScreen screens.Screen
	useScreens   bool
}

// NewCaptureEventIntent creates a new CaptureEvent intent.
func NewCaptureEventIntent(context *CaptureEventContext) (*CaptureEventIntent, error) {
	if err := context.Validate(); err != nil {
		return nil, err
	}

	formModel := models.NewCaptureForm(context.CLIEventService)
	base := NewBaseIntent()

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
		useScreens: true,
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
	if i.useScreens {
		breadcrumbs := []string{"Main Menu", "Capture Event"}
		i.activeScreen = captureScreens.NewStrategySelectScreen(breadcrumbs)

		termInfo := i.GetTerminalInfo()
		width, height := 120, 40
		if termInfo != nil {
			width = termInfo.Width
			height = termInfo.Height
		}

		i.activeScreen.SetTerminalInfo(width, height)
		i.activeScreen.SetTheme(i.Theme())
		i.activeScreen.SetLogo(i.GetLogo(), i.GetLogoSpacing())

		if initable, ok := i.activeScreen.(interface{ Init() tea.Cmd }); ok {
			return initable.Init()
		}
		return nil
	}

	i.initializeFormForNew()
	return func() tea.Msg { return nil }
}

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
	i.state.reviewState.Event = &career.Event{
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Tags:       make([]string, 0),
		Categories: make([]string, 0),
	}
	return func() tea.Msg { return nil }
}

// Update processes a message in the intent.
// This is Bubble Tea's standard Update, constrained to intent-local state.
func (i *CaptureEventIntent) Update(msg tea.Msg) tea.Cmd {
	if !i.active {
		return nil
	}

	switch msg := msg.(type) {
	case FormSubmittedMsg:
		if msg.Event == nil {
			i.setFailed("INVALID_FORM", "Form submission with nil event", nil)
			return nil
		}

		if err := msg.Event.Validate(); err != nil {
			i.setFailed("VALIDATION_ERROR", fmt.Sprintf("Form validation failed: %v", err), err)
			return nil
		}

		i.state.reviewState = &ReviewInferredEventState{
			Event:          msg.Event,
			InferredBursts: make([]*career.Burst, 0),
			InferredFacts:  make([]*career.Fact, 0),
			AcceptedBursts: make([]*career.Burst, 0),
			AcceptedFacts:  make([]*career.Fact, 0),
			RejectedItems:  make(map[string]string),
		}
		i.state.currentState = CaptureStateReview
		i.activeScreen = nil
		return nil

	case SubmitCompleteMsg:
		i.state.submitModal = feedback.NewSuccessModal("Event saved!")
		return tea.Tick(2*time.Second, func(_ time.Time) tea.Msg {
			return DismissModalMsg{}
		})

	case SubmitErrorMsg:
		i.state.submitModal = feedback.NewErrorModal("Save Failed", msg.Message)
		return nil

	case DismissModalMsg:
		if i.state.submitModal != nil {
			i.state.submitModal = nil
			i.state.postSaveReview = true
			i.state.currentState = CaptureStateReview

			if i.useScreens {
				breadcrumbs := []string{"Main Menu", "Capture Event", "Review Enrichment"}
				i.activeScreen = captureScreens.NewEventReviewScreen(
					breadcrumbs,
					i.state.reviewState.Event,
					i.state.reviewState.InferredBursts,
					i.state.reviewState.InferredFacts,
				)

				termInfo := i.GetTerminalInfo()
				width, height := 120, 40
				if termInfo != nil {
					width = termInfo.Width
					height = termInfo.Height
				}

				i.activeScreen.SetTerminalInfo(width, height)
				i.activeScreen.SetTheme(i.Theme())
				i.activeScreen.SetLogo(i.GetLogo(), i.GetLogoSpacing())
			}
		}
		return nil
	}

	if i.state.submitModal != nil {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.Type == tea.KeyEsc {
				i.state.submitModal = nil
				return nil
			}
			return nil
		case feedback.ModalSpinnerTickMsg:
			if i.state.submitModal.Type == feedback.ModalLoading {
				return i.state.submitModal.Update(msg)
			}
			return nil
		default:
			return nil
		}
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch HandleGlobalKeys(keyMsg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			i.helpModal.Toggle()
			return nil
		}
	}

	if i.useScreens && i.activeScreen != nil {
		if i.state.reviewState != nil && i.state.reviewState.EditingMode != EditingModeNone {
			return i.updateEditingModal(msg)
		}

		cmd, result := i.activeScreen.Update(msg)

		if result != nil {
			return i.handleScreenResult(result)
		}

		return cmd
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
			if i.state.selectedStrategyIndex > 0 {
				i.state.selectedStrategyIndex--
			}
			return nil

		case "down", "j":
			if i.state.selectedStrategyIndex < 1 {
				i.state.selectedStrategyIndex++
			}
			return nil

		case "enter":
			strategies := []CaptureStrategy{StrategyQuick, StrategyManual}
			i.state.strategy = strategies[i.state.selectedStrategyIndex]
			i.state.captureForm.SetStrategy(string(i.state.strategy))
			i.state.currentState = CaptureStateForm
			return i.state.captureForm.Init()
		}

		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			i.ToggleHelp()
			return nil
		case KeyBack:
			i.setCancelled()
			return nil
		}

	case StrategySelectedMsg:
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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			i.ToggleHelp()
			return nil
		case KeyBack:
			if i.context.PreviousEvent != nil {
				i.setCancelled()
				return nil
			}
			i.state.currentState = CaptureStateChooseStrategy
			return nil
		}

		switch msg.String() {
		case "ctrl+s":
			return i.state.captureForm.SubmitForm()
		}

	case models.SubmitMsg:
		if msg.Err != nil {
			i.state.error = &IntentError{
				Code:    "FORM_SUBMISSION_ERROR",
				Message: msg.Err.Error(),
				Cause:   msg.Err,
			}
			return nil
		}

		if msg.Event == nil {
			i.setFailed("INVALID_FORM", "Form submission with nil event", nil)
			return nil
		}

		if err := msg.Event.Validate(); err != nil {
			i.setFailed("VALIDATION_ERROR", fmt.Sprintf("Form validation failed: %v", err), err)
			return nil
		}

		i.state.reviewState.Event = msg.Event
		i.state.currentState = CaptureStateReview
		return nil

	case FormSubmittedMsg:
		if msg.Event == nil {
			i.setFailed("INVALID_FORM", "Form submission with nil event", nil)
			return nil
		}

		if err := msg.Event.Validate(); err != nil {
			i.setFailed("VALIDATION_ERROR", fmt.Sprintf("Form validation failed: %v", err), err)
			return nil
		}

		i.state.reviewState.Event = msg.Event
		i.state.currentState = CaptureStateReview
		return nil
	}

	_, formCmd := i.state.captureForm.Update(msg)

	return formCmd
}

// updateReviewInferredEvent handles messages while reviewing inferred bursts and facts.
// It processes review confirmations, edits, and transitions to submit state.
func (i *CaptureEventIntent) updateReviewInferredEvent(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			i.ToggleHelp()
			return nil
		case KeyBack:
			if i.state.reviewState.EditingMode != EditingModeNone {
				i.state.reviewState.metadataModal = nil
				i.state.reviewState.burstModal = nil
				i.state.reviewState.factModal = nil
				i.state.reviewState.EditingMode = EditingModeNone
				return nil
			}
			i.state.currentState = CaptureStateForm
			return nil
		}
	}

	switch i.state.reviewState.EditingMode {
	case EditingModeMetadata:
		if i.state.reviewState.metadataModal != nil {
			modal, cmd := i.state.reviewState.metadataModal.Update(msg)
			//nolint:errcheck // Type assertion is safe - Update always returns same type.
			i.state.reviewState.metadataModal = modal.(*models.MetadataEditorModelNew)

			if i.state.reviewState.metadataModal.IsSubmitted() {
				i.state.reviewState.Event = i.state.reviewState.metadataModal.GetEvent()
				i.state.reviewState.metadataModal = nil
				i.state.reviewState.EditingMode = EditingModeNone
			} else if i.state.reviewState.metadataModal.IsCancelled() {
				i.state.reviewState.metadataModal = nil
				i.state.reviewState.EditingMode = EditingModeNone
			}
			return cmd
		}

	case EditingModeBursts:
		if i.state.reviewState.burstModal != nil {
			modal, cmd := i.state.reviewState.burstModal.Update(msg)
			//nolint:errcheck // Type assertion is safe - Update always returns same type.
			i.state.reviewState.burstModal = modal.(*models.BurstSuggestionModelNew)

			return cmd
		}

	case EditingModeFacts:
		if i.state.reviewState.factModal != nil {
			modal, cmd := i.state.reviewState.factModal.Update(msg)
			//nolint:errcheck // Type assertion is safe - Update always returns same type.
			i.state.reviewState.factModal = modal.(*models.FactEditorModelNew)

			if i.state.reviewState.factModal.IsSubmitted() {
				i.state.reviewState.factModal = nil
				i.state.reviewState.EditingMode = EditingModeNone
			} else if i.state.reviewState.factModal.IsCancelled() {
				i.state.reviewState.factModal = nil
				i.state.reviewState.EditingMode = EditingModeNone
			}
			return cmd
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+s", "enter":
			if i.state.postSaveReview {
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
			i.state.currentState = CaptureStateSubmit
			return i.performSubmit()

		case "e":
			i.state.reviewState.EditingMode = EditingModeMetadata
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
			i.state.reviewState.EditingMode = EditingModeBursts
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
			i.state.reviewState.EditingMode = EditingModeFacts
			if i.state.reviewState.factModal == nil {
				i.state.reviewState.factModal = models.NewFactEditorModelNew(
					&career.Fact{},
					i.context.CareerService,
					context.Background(),
				)
				return i.state.reviewState.factModal.Init()
			}
			return nil

		case "a":
			i.acceptCurrentItem()
			return nil

		case "r":
			i.rejectCurrentItem()
			return nil

		case "j", "down":
			totalItems := len(i.state.reviewState.InferredBursts) + len(i.state.reviewState.InferredFacts)
			if totalItems > 0 {
				i.state.reviewState.SelectedIndex++
				if i.state.reviewState.SelectedIndex >= totalItems {
					i.state.reviewState.SelectedIndex = 0
				}
				if i.state.reviewState.SelectedIndex < len(i.state.reviewState.InferredBursts) {
					i.state.reviewState.SelectedItemType = "burst"
				} else {
					i.state.reviewState.SelectedItemType = "fact"
					i.state.reviewState.SelectedIndex -= len(i.state.reviewState.InferredBursts)
				}
			}
			return nil

		case "k", "up":
			totalItems := len(i.state.reviewState.InferredBursts) + len(i.state.reviewState.InferredFacts)
			if totalItems > 0 {
				i.state.reviewState.SelectedIndex--
				if i.state.reviewState.SelectedIndex < 0 {
					i.state.reviewState.SelectedIndex = totalItems - 1
				}
				if i.state.reviewState.SelectedIndex < len(i.state.reviewState.InferredBursts) {
					i.state.reviewState.SelectedItemType = "burst"
				} else {
					i.state.reviewState.SelectedItemType = "fact"
					i.state.reviewState.SelectedIndex -= len(i.state.reviewState.InferredBursts)
				}
			}
			return nil
		}

	case ReviewConfirmedMsg:
		i.state.reviewState.AcceptedBursts = msg.AcceptedBursts
		i.state.reviewState.AcceptedFacts = msg.AcceptedFacts
		i.state.reviewState.RejectedItems = msg.RejectedItems
		i.state.currentState = CaptureStateSubmit
		return i.performSubmit()

	case ReviewCancelledMsg:
		i.setCancelled()
		return nil

	case ReviewBackMsg:
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
		i.state.error = &IntentError{
			Code:    msg.Code,
			Message: msg.Message,
			Cause:   msg.Cause,
		}
		return nil

	case tea.KeyMsg:
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			i.ToggleHelp()
			return nil
		case KeyBack:
			i.state.currentState = CaptureStateReview
			return nil
		}

		switch msg.String() {
		case "r":
			return i.performSubmit()
		}
	}

	return nil
}

// performSubmit performs the actual submission of the event.
// It calls the domain service to persist the event to the database and optionally enriches it.
// Logs: Event submission start, validation results, service calls, and completion status.
func (i *CaptureEventIntent) performSubmit() tea.Cmd {
	event := i.state.reviewState.Event
	acceptedFacts := i.state.reviewState.AcceptedFacts
	strategy := i.state.strategy
	careerService := i.state.context.CareerService
	eventService := i.eventService

	return func() tea.Msg {
		if event == nil {
			return SubmitErrorMsg{
				Code:    "MISSING_EVENT",
				Message: "No event data to submit",
				Cause:   nil,
			}
		}

		if err := event.Validate(); err != nil {
			return SubmitErrorMsg{
				Code:    "VALIDATION_ERROR",
				Message: fmt.Sprintf("Event validation failed: %v", err),
				Cause:   err,
			}
		}

		if eventService == nil {
			return SubmitErrorMsg{
				Code:    "SERVICE_ERROR",
				Message: "Event service not initialized",
				Cause:   nil,
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if strategy == StrategyQuick && event.Date.IsZero() {
			event.Date = time.Now()
		}

		mode := careerservice.ManualEntry

		if careerService == nil {
			return SubmitErrorMsg{
				Code:    "SERVICE_ERROR",
				Message: "Career service not initialized",
				Cause:   nil,
			}
		}

		err := careerService.CaptureEvent(ctx, event, mode)

		if err != nil {
			return SubmitErrorMsg{
				Code:    "PERSISTENCE_ERROR",
				Message: fmt.Sprintf("Failed to save event: %v", err),
				Cause:   err,
			}
		}

		if careerService != nil && len(acceptedFacts) > 0 {
			for _, fact := range acceptedFacts {
				if fact.ID == "" {
					if fact.SourceEventID == "" {
						fact.SourceEventID = event.ID
					}

					if err := careerService.SaveFact(ctx, fact); err != nil {
						continue
					}
				}
			}
		}

		return SubmitCompleteMsg{}
	}
}

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
		return i.state.currentState
	}
}

// getStateContent returns the content for the current state.
func (i *CaptureEventIntent) getStateContent() string {
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

	if i.useScreens && i.activeScreen != nil {
		baseView := i.activeScreen.View()

		termInfo := i.GetTerminalInfo()
		width, height := 80, 24
		if termInfo != nil {
			width = termInfo.Width
			height = termInfo.Height
		}

		if i.state.submitModal != nil {
			modalContent := i.state.submitModal.Render(width, height)
			return i.overlayModal(baseView, modalContent, width, height)
		}

		if i.state.reviewState != nil && i.state.reviewState.EditingMode != EditingModeNone {
			modalContent := i.getEditingModalContent()
			if modalContent != nil {
				return i.renderModalOverlay(baseView, modalContent)
			}
		}

		return baseView
	}

	view := i.CreateViewWithBreadcrumbs("Main Menu", "Capture Event", i.getStateName())

	if info := i.GetTerminalInfo(); info != nil && info.Height < 30 {
		if logo := i.GetLogo(); logo != nil {
			view.WithLogo(logo, 0)
		}
	}

	if i.state.error != nil {
		i.SetError(i.state.error)
	}

	content := i.getStateContent()
	view.WithContent(content)

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

		optStyle := lipgloss.NewStyle().Foreground(i.getPrimaryColor())
		if idx == i.state.selectedStrategyIndex {
			optStyle = optStyle.Foreground(i.getAccentColor()).Bold(true)
		}

		line := fmt.Sprintf("%s%s - %s", prefix, s.label, s.description)
		content.WriteString(optStyle.Render(line) + "\n")
	}

	return i.getCardStyle().Render(content.String())
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
// When in editing mode, the modal is rendered as an overlay on top of the review content.
func (i *CaptureEventIntent) viewReviewInferredEvent() string {
	baseView := i.buildReviewBaseView()

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

	if i.state.reviewState.Event != nil {
		title := i.state.reviewState.Event.Text
		if len(title) > 40 {
			title = title[:37] + "..."
		}
		sb.WriteString(fmt.Sprintf("│ Event: %s                    │\n", title))
		sb.WriteString("│                                                │\n")
	}

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
	return sb.String()
}

// renderModalOverlay renders a modal overlay on top of the background content.
func (i *CaptureEventIntent) renderModalOverlay(background string, modalContent *modalContentData) string {
	if modalContent == nil {
		return background
	}

	info := i.GetTerminalInfo()
	width := 80
	height := 24
	if info != nil {
		width = info.Width
		height = info.Height
	}

	overlay := feedback.NewOverlayModal(modalContent.title, modalContent.content)
	overlay.SetFooter(modalContent.footer)
	overlay.SetWidth(80)

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
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.String() == "esc" {
			i.state.reviewState.metadataModal = nil
			i.state.reviewState.burstModal = nil
			i.state.reviewState.factModal = nil
			i.state.reviewState.EditingMode = EditingModeNone
			return nil
		}
	}

	switch i.state.reviewState.EditingMode {
	case EditingModeMetadata:
		if i.state.reviewState.metadataModal != nil {
			modal, cmd := i.state.reviewState.metadataModal.Update(msg)
			//nolint:errcheck // Type assertion is safe - Update always returns same type.
			i.state.reviewState.metadataModal = modal.(*models.MetadataEditorModelNew)

			if i.state.reviewState.metadataModal.IsSubmitted() {
				i.state.reviewState.Event = i.state.reviewState.metadataModal.GetEvent()
				i.state.reviewState.metadataModal = nil
				i.state.reviewState.EditingMode = EditingModeNone
			} else if i.state.reviewState.metadataModal.IsCancelled() {
				i.state.reviewState.metadataModal = nil
				i.state.reviewState.EditingMode = EditingModeNone
			}
			return cmd
		}

	case EditingModeBursts:
		if i.state.reviewState.burstModal != nil {
			modal, cmd := i.state.reviewState.burstModal.Update(msg)
			//nolint:errcheck // Type assertion is safe - Update always returns same type.
			i.state.reviewState.burstModal = modal.(*models.BurstSuggestionModelNew)
			return cmd
		}

	case EditingModeFacts:
		if i.state.reviewState.factModal != nil {
			modal, cmd := i.state.reviewState.factModal.Update(msg)
			//nolint:errcheck // Type assertion is safe - Update always returns same type.
			i.state.reviewState.factModal = modal.(*models.FactEditorModelNew)

			if i.state.reviewState.factModal.IsSubmitted() {
				i.state.reviewState.factModal = nil
				i.state.reviewState.EditingMode = EditingModeNone
			} else if i.state.reviewState.factModal.IsCancelled() {
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
			i.state.reviewState.InferredBursts = append(
				i.state.reviewState.InferredBursts[:idx],
				i.state.reviewState.InferredBursts[idx+1:]...)
			if i.state.reviewState.SelectedIndex >= len(i.state.reviewState.InferredBursts) && len(i.state.reviewState.InferredBursts) > 0 {
				i.state.reviewState.SelectedIndex = len(i.state.reviewState.InferredBursts) - 1
			}
		}
	} else if i.state.reviewState.SelectedItemType == "fact" {
		idx := i.state.reviewState.SelectedIndex
		if idx >= 0 && idx < len(i.state.reviewState.InferredFacts) {
			fact := i.state.reviewState.InferredFacts[idx]
			i.state.reviewState.AcceptedFacts = append(i.state.reviewState.AcceptedFacts, fact)
			i.state.reviewState.InferredFacts = append(
				i.state.reviewState.InferredFacts[:idx],
				i.state.reviewState.InferredFacts[idx+1:]...)
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
			if i.state.reviewState.RejectedItems == nil {
				i.state.reviewState.RejectedItems = make(map[string]string)
			}
			i.state.reviewState.RejectedItems[burst.ID] = "user_rejected"
			i.state.reviewState.InferredBursts = append(
				i.state.reviewState.InferredBursts[:idx],
				i.state.reviewState.InferredBursts[idx+1:]...)
			if i.state.reviewState.SelectedIndex >= len(i.state.reviewState.InferredBursts) && len(i.state.reviewState.InferredBursts) > 0 {
				i.state.reviewState.SelectedIndex = len(i.state.reviewState.InferredBursts) - 1
			}
		}
	} else if i.state.reviewState.SelectedItemType == "fact" {
		idx := i.state.reviewState.SelectedIndex
		if idx >= 0 && idx < len(i.state.reviewState.InferredFacts) {
			fact := i.state.reviewState.InferredFacts[idx]
			if i.state.reviewState.RejectedItems == nil {
				i.state.reviewState.RejectedItems = make(map[string]string)
			}
			i.state.reviewState.RejectedItems[fact.ID] = "user_rejected"
			i.state.reviewState.InferredFacts = append(
				i.state.reviewState.InferredFacts[:idx],
				i.state.reviewState.InferredFacts[idx+1:]...)
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

// GetState returns the current workflow step of the capture event intent.
//
// The method takes no parameters beyond the receiver.
//
// Returns a string matching one of the CaptureState constants:
// "choose_strategy", "form", "review", or "submit". Tests and the intent
// router use this value to inspect which sub-flow screen the intent is
// currently displaying.
func (i *CaptureEventIntent) GetState() string {
	return i.state.currentState
}

// GetForm returns the current form model instance (for test and debug)
func (i *CaptureEventIntent) GetForm() *models.CaptureForm {
	if i == nil || i.state == nil {
		return nil
	}
	return i.state.captureForm
}

// Result returns the intent's outcome as a type-erased IntentResult.
//
// The method takes no parameters beyond the receiver.
//
// Returns nil if the intent is still active and has not yet completed,
// been cancelled, or failed. Returns a non-nil *IntentResult[interface{}]
// once the intent reaches a terminal state. The returned value is
// converted from the strongly-typed IntentResult[*CaptureEventResult],
// preserving the Status, Data, Error, and Metadata fields. The intent
// router calls this method after every Update cycle to detect whether the
// intent has finished.
func (i *CaptureEventIntent) Result() *IntentResult[interface{}] {
	if i.result == nil {
		return nil
	}
	return &IntentResult[interface{}]{
		Status:   i.result.Status,
		Data:     i.result.Data,
		Error:    i.result.Error,
		Metadata: i.result.Metadata,
	}
}

// handleScreenResult processes results from screen updates.
// This is the central hub for all screen-to-intent communication.
//
// Uses ScreenResultDispatcher pattern to eliminate repetitive type switching.
// CaptureEventIntent implements ScreenResultHandler interface for compile-time safety.
func (i *CaptureEventIntent) handleScreenResult(result screens.ScreenResult) tea.Cmd {
	return behaviors.NewScreenResultDispatcher(i).Dispatch(result)
}

// HandleNavigate handles navigation actions from screens.
//
// Implements ScreenResultHandler interface.
func (i *CaptureEventIntent) HandleNavigate(result *screens.NavigateResult) tea.Cmd {
	data := result.Data()

	if action, ok := data.(string); ok {
		switch action {
		case "edit_metadata":
			if i.state.reviewState == nil || i.state.reviewState.Event == nil {
				return i.setFailedCmd("NO_EVENT", "No event to edit", nil)
			}
			i.state.reviewState.EditingMode = EditingModeMetadata
			i.state.reviewState.metadataModal = models.NewMetadataEditorModelNew(
				i.state.reviewState.Event,
				i.context.CareerService,
				i.context.CLIEventService,
				context.Background(),
			)
			return i.state.reviewState.metadataModal.Init()

		case "edit_bursts":
			i.state.reviewState.EditingMode = EditingModeBursts
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
			i.state.reviewState.EditingMode = EditingModeFacts
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

	if strategy, ok := data.(CaptureStrategy); ok {
		i.state.strategy = strategy
		return i.transitionToFormScreen(strategy)
	}

	return i.setFailedCmd("INVALID_NAVIGATION_DATA", fmt.Sprintf("Invalid navigation data type: %T", data), nil)
}

// HandleCancel handles cancellation from screens.
//
// Implements ScreenResultHandler interface.
func (i *CaptureEventIntent) HandleCancel(_ *screens.CancelResult) tea.Cmd {
	switch i.state.currentState {
	case CaptureStateChooseStrategy:
		i.setCancelled()
		return nil

	case CaptureStateForm:
		if i.state.context.PreviousEvent != nil {
			i.setCancelled()
			return nil
		}
		return i.transitionToStrategyScreen()

	case CaptureStateReview:
		return i.transitionToFormScreen(i.state.strategy)

	case CaptureStateSubmit:
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
		if event, ok := data.(*career.Event); ok {
			if err := event.Validate(); err != nil {
				return i.setFailedCmd("VALIDATION_ERROR", fmt.Sprintf("Event validation failed: %v", err), err)
			}

			i.state.reviewState = &ReviewInferredEventState{
				Event:          event,
				InferredBursts: make([]*career.Burst, 0),
				InferredFacts:  make([]*career.Fact, 0),
				AcceptedBursts: make([]*career.Burst, 0),
				AcceptedFacts:  make([]*career.Fact, 0),
				RejectedItems:  make(map[string]string),
			}

			i.state.submitModal = feedback.NewLoadingModal("Saving event...", false)
			return tea.Batch(i.performSubmit(), i.state.submitModal.Init())
		}
		return i.setFailedCmd("INVALID_FORM_DATA", fmt.Sprintf("Invalid form data type: %T", data), nil)

	case CaptureStateReview:
		if reviewData, ok := data.(map[string]interface{}); ok {
			//nolint:errcheck // Type assertions are safe for map data extraction.
			event, _ := reviewData["event"].(*career.Event)
			//nolint:errcheck // Type assertions are safe for map data extraction.
			bursts, _ := reviewData["bursts"].([]*career.Burst)
			//nolint:errcheck // Type assertions are safe for map data extraction.
			facts, _ := reviewData["facts"].([]*career.Fact)

			i.state.reviewState.Event = event
			i.state.reviewState.AcceptedBursts = bursts
			i.state.reviewState.AcceptedFacts = facts

			i.state.submitModal = feedback.NewLoadingModal("Saving event...", false)
			return tea.Batch(i.performSubmit(), i.state.submitModal.Init())
		}
		return i.setFailedCmd("INVALID_REVIEW_DATA", fmt.Sprintf("Invalid review data type: %T", data), nil)

	case CaptureStateSubmit:
		if submitData, ok := data.(map[string]interface{}); ok {
			//nolint:errcheck // Type assertion is safe for map data extraction.
			event, _ := submitData["event"].(*career.Event)
			//nolint:errcheck // Type assertion is safe for map data extraction.
			bursts, _ := submitData["bursts"].([]*career.Burst)
			//nolint:errcheck // Type assertion is safe for map data extraction.
			facts, _ := submitData["facts"].([]*career.Fact)

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
	data := result.Data()
	if errorData, ok := data.(map[string]interface{}); ok {
		//nolint:errcheck // Type assertion is safe for map data extraction.
		err, _ := errorData["error"].(error)
		//nolint:errcheck // Type assertion is safe for map data extraction.
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

// transitionToStrategyScreen transitions to the strategy selection screen.
func (i *CaptureEventIntent) transitionToStrategyScreen() tea.Cmd {
	i.state.currentState = CaptureStateChooseStrategy
	breadcrumbs := []string{"Main Menu", "Capture Event"}
	i.activeScreen = captureScreens.NewStrategySelectScreen(breadcrumbs)

	termInfo := i.GetTerminalInfo()
	width, height := 120, 40
	if termInfo != nil {
		width = termInfo.Width
		height = termInfo.Height
	}

	i.activeScreen.SetTerminalInfo(width, height)
	i.activeScreen.SetTheme(i.Theme())
	i.activeScreen.SetLogo(i.GetLogo(), i.GetLogoSpacing())

	return nil
}

// transitionToFormScreen transitions to the event form screen.
func (i *CaptureEventIntent) transitionToFormScreen(strategy CaptureStrategy) tea.Cmd {
	i.state.currentState = CaptureStateForm
	i.state.strategy = strategy
	breadcrumbs := []string{"Main Menu", "Capture Event", "Form"}
	i.activeScreen = captureScreens.NewEventFormScreen(
		i.eventService,
		breadcrumbs,
		strategy,
	)

	termInfo := i.GetTerminalInfo()
	width, height := 120, 40
	if termInfo != nil {
		width = termInfo.Width
		height = termInfo.Height
	}

	i.activeScreen.SetTerminalInfo(width, height)
	i.activeScreen.SetTheme(i.Theme())
	i.activeScreen.SetLogo(i.GetLogo(), i.GetLogoSpacing())

	if initable, ok := i.activeScreen.(interface{ Init() tea.Cmd }); ok {
		return initable.Init()
	}
	return nil
}

// overlayModal overlays modal content on top of background content (centered).
func (i *CaptureEventIntent) overlayModal(background, modal string, width, _ int) string {
	bgLines := strings.Split(background, "\n")
	modalLines := strings.Split(modal, "\n")

	bgHeight := len(bgLines)
	modalHeight := len(modalLines)
	startLine := (bgHeight - modalHeight) / 2
	if startLine < 0 {
		startLine = 0
	}

	result := make([]string, len(bgLines))
	copy(result, bgLines)

	for i, modalLine := range modalLines {
		lineIndex := startLine + i
		if lineIndex >= 0 && lineIndex < len(result) {
			centeredModalLine := lipgloss.PlaceHorizontal(width, lipgloss.Center, modalLine)
			result[lineIndex] = centeredModalLine
		}
	}

	return strings.Join(result, "\n")
}
