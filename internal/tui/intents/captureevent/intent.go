package captureevent

import (
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/tui/forms"
	"github.com/baphled/kariya/internal/tui/intents"
	"github.com/baphled/kariya/internal/tui/views/event"
	"github.com/baphled/kariya/internal/ui/behaviors"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/uikit/feedback"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	captureEventTitle = "Capture Event"
	mainMenuTitle     = "Main Menu"
)

// NewIntent creates and returns a fully initialised CaptureEvent intent.
//
// Expected:
//   - ctx must be non-nil and pass Validate()
//   - ctx.CLIEventService should be set for form submission
//
// Returns:
//   - A ready-to-use Intent on success
//   - error if context validation fails
//
// Side effects:
//   - Allocates the form model and empty review state.
func NewIntent(ctx *IntentValidator) (*Intent, error) {
	if err := ctx.Validate(); err != nil {
		return nil, err
	}

	base := intents.NewBaseIntent()

	return &Intent{
		BaseIntent:   base,
		context:      ctx,
		eventService: ctx.CLIEventService,
		active:       true,
		currentState: StateChooseStrategy,
		reviewState: &ReviewInferredEventState{
			AcceptedBursts: make([]*career.Burst, 0),
			AcceptedFacts:  make([]*career.Fact, 0),
			AcceptedSkills: make([]*career.Skill, 0),
			RejectedItems:  make(map[string]string),
		},
	}, nil
}

// Init prepares the intent for its first render cycle.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (i *Intent) Init() tea.Cmd {
	i.activeView = event.NewStrategySelect()
	return i.activeView.Init()
}

// Update processes a single Bubble Tea message and advances the workflow.
//
// Expected:
//   - The intent must be active (returns nil otherwise).
//
// Returns:
//   - A tea.Cmd for follow-up work, or nil.
//
// Side effects:
//   - May transition between states, create/dismiss modals, or mark the
//     intent as completed/cancelled/failed.
func (i *Intent) Update(msg tea.Msg) tea.Cmd {
	if !i.active {
		return nil
	}

	if cmd, handled := i.handleUpdateMessage(msg); handled {
		return cmd
	}

	if i.submitModal != nil {
		return i.updateSubmitModal(msg)
	}

	if cmd, handled := i.handleGlobalKeyMsg(msg); handled {
		return cmd
	}

	return i.updateActiveView(msg)
}

// handleUpdateMessage handles message type branches for Update.
//
// Returns:
//   - A tea.Cmd to execute.
//   - true when the message was handled.
//
// Side effects:
//   - May update review state, create/dismiss modals, or complete the intent.
func (i *Intent) handleUpdateMessage(msg tea.Msg) (tea.Cmd, bool) {
	switch msg := msg.(type) {
	case SubmitCompleteMsg:
		return i.handleSubmitComplete(msg), true
	case InferenceCompleteMsg:
		i.handleInferenceComplete(msg)
		return nil, true
	case PostSavePersistenceCompleteMsg:
		i.handlePostSavePersistenceComplete(msg)
		return nil, true
	case SubmitErrorMsg:
		return i.handleSubmitError(msg), true
	case feedback.ModalAutoDismissMsg:
		return func() tea.Msg { return DismissModalMsg{} }, true
	case DismissModalMsg:
		i.handleDismissModal()
		return nil, true
	case SubmitMsg:
		return i.handleSubmitMsg(msg), true
	default:
		return nil, false
	}
}

// handleSubmitComplete prepares the success modal and kicks off inference.
//
// Returns:
//   - A batched tea.Cmd initialising the modal and inference.
//
// Side effects:
//   - Creates a success modal assigned to i.submitModal.
func (i *Intent) handleSubmitComplete(_ SubmitCompleteMsg) tea.Cmd {
	i.submitModal = feedback.NewSuccessModal("Event saved!")
	return tea.Batch(i.submitModal.Init(), i.performInference())
}

// handleInferenceComplete stores inference results and refreshes the review view if active.
//
// Expected:
//   - reviewState is initialised.
//
// Side effects:
//   - Updates inferred skills, bursts, suggestions, and facts.
//   - Updates the active review view when present.
func (i *Intent) handleInferenceComplete(msg InferenceCompleteMsg) {
	i.reviewState.InferredSkills = msg.InferredSkills
	i.reviewState.InferredBursts = msg.InferredBursts
	i.reviewState.InferredBurstSuggestions = msg.InferredBurstSuggestions
	i.reviewState.InferredFacts = msg.InferredFacts
	if screen, ok := i.activeView.(*event.Review); ok {
		screen.SetSuggestedSkills(display.SkillSuggestionsFromDomain(msg.InferredSkills))
		screen.SetSuggestedBursts(display.BurstsFromDomain(msg.InferredBursts))
		screen.SetSuggestedFacts(display.FactsFromDomain(msg.InferredFacts))
	}
}

// handlePostSavePersistenceComplete finalises the intent result.
//
// Side effects:
//   - Sets the intent result and marks the intent inactive.
func (i *Intent) handlePostSavePersistenceComplete(msg PostSavePersistenceCompleteMsg) {
	i.result = &intents.IntentResult[*Result]{
		Status: intents.Completed,
		Data: &Result{
			Event:  msg.Event,
			Bursts: msg.Bursts,
			Facts:  msg.Facts,
			Skills: msg.Skills,
		},
	}
	i.active = false
}

// handleSubmitError shows a failure modal for submit errors.
//
// Returns:
//   - A tea.Cmd to initialise the error modal.
//
// Side effects:
//   - Assigns i.submitModal to a failure modal.
func (i *Intent) handleSubmitError(msg SubmitErrorMsg) tea.Cmd {
	i.submitModal = feedback.NewErrorModal("Save Failed", msg.Message)
	return i.submitModal.Init()
}

// handleDismissModal clears the submit modal and enters review on success.
//
// Side effects:
//   - Clears i.submitModal.
//   - Transitions to review state and view when the modal was a success.
func (i *Intent) handleDismissModal() {
	if i.submitModal == nil {
		return
	}
	wasSuccess := i.submitModal.Type == feedback.ModalSuccess
	i.submitModal = nil
	if wasSuccess {
		i.currentState = StateReview
		i.activeView = event.NewReview(
			display.EventFromDomain(i.reviewState.Event),
			display.BurstsFromDomain(i.reviewState.InferredBursts),
			display.FactsFromDomain(i.reviewState.InferredFacts),
			display.SkillSuggestionsFromDomain(i.reviewState.InferredSkills),
		)
	}
}

// handleSubmitMsg translates submit messages into view results when applicable.
//
// Returns:
//   - A tea.Cmd from handling the submit or nil.
//
// Side effects:
//   - May present a validation error modal or submit a form result.
func (i *Intent) handleSubmitMsg(msg SubmitMsg) tea.Cmd {
	if i.currentState != StateForm {
		return nil
	}
	if msg.Error != nil {
		return i.showValidationErrorModal(msg.Error.Error())
	}
	if msg.Event == nil {
		return nil
	}
	formData := forms.GetCaptureEventFormData(msg.Event)
	formData.SubmitConfirmed = true
	return i.handleViewResult(&widgets.SubmitViewResult{FormData: formData})
}

// handleGlobalKeyMsg processes global key bindings.
//
// Returns:
//   - A tea.Cmd when handled.
//   - true when a key binding was handled.
//
// Side effects:
//   - Toggles help or quits the application.
func (i *Intent) handleGlobalKeyMsg(msg tea.Msg) (tea.Cmd, bool) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return nil, false
	}
	switch intents.HandleGlobalKeys(keyMsg) {
	case intents.KeyQuit:
		return tea.Quit, true
	case intents.KeyHelp:
		i.ToggleHelp()
		return nil, true
	default:
		return nil, false
	}
}

// updateActiveView updates the active view and handles view results.
//
// Returns:
//   - A tea.Cmd from the view update or nil.
//
// Side effects:
//   - May update the editing modal or transition intent state via view results.
func (i *Intent) updateActiveView(msg tea.Msg) tea.Cmd {
	if i.activeView == nil {
		return nil
	}
	if i.reviewState != nil && i.reviewState.EditingMode != EditingModeNone {
		return i.updateEditingModal(msg)
	}
	cmd, result := i.activeView.Update(msg)
	if result != nil {
		return i.handleViewResult(result)
	}
	return cmd
}

// updateSubmitModal handles messages while the submit modal is displayed.
//
// Returns:
//   - A tea.Cmd from the modal's update, or nil.
func (i *Intent) updateSubmitModal(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc && i.submitModal.Type != feedback.ModalLoading {
			return func() tea.Msg { return DismissModalMsg{} }
		}
		return nil
	case feedback.ModalSpinnerTickMsg:
		if i.submitModal.Type == feedback.ModalLoading {
			return i.submitModal.Update(msg)
		}
		return nil
	case feedback.ModalCountdownTickMsg:
		if i.submitModal != nil && i.submitModal.Type == feedback.ModalSuccess {
			return i.submitModal.Update(msg)
		}
		return nil
	default:
		return nil
	}
}

// View renders the intent's current visual state.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (i *Intent) View() string {
	if !i.active {
		return "CaptureEvent intent is not active"
	}

	if i.activeView == nil {
		return "No active view"
	}

	view := i.CreateViewWithBreadcrumbs(i.breadcrumbs()...)
	view.WithContent(i.activeView.RenderContent())
	view.WithHelp(i.activeView.HelpText())
	baseView := view.Render()

	if i.submitModal != nil {
		termInfo := i.GetTerminalInfo()
		width, height := 80, 24
		if termInfo != nil && termInfo.Width > 0 && termInfo.Height > 0 {
			width = termInfo.Width
			height = termInfo.Height
		}
		rendered := i.submitModal.Render(width, height)
		modal := &behaviors.StaticViewModel{Content: rendered}
		return behaviors.RenderModalOverlay(modal, baseView)
	}

	if i.reviewState != nil && i.reviewState.EditingMode != EditingModeNone {
		return i.renderEditingOverlay(baseView)
	}

	return baseView
}

// renderEditingOverlay renders a modal overlay for the active editing mode.
//
// Returns:
//   - The base view with a modal overlay applied, or the base view unchanged.
func (i *Intent) renderEditingOverlay(baseView string) string {
	if i.reviewState.EditingMode == EditingModeBursts && i.reviewState.burstModal != nil {
		rendered := i.reviewState.burstModal.View()
		modal := &behaviors.StaticViewModel{Content: rendered}
		return behaviors.RenderModalOverlay(modal, baseView)
	}

	if i.reviewState.EditingMode == EditingModeSkills && i.reviewState.skillModal != nil {
		rendered := i.reviewState.skillModal.View()
		modal := &behaviors.StaticViewModel{Content: rendered}
		return behaviors.RenderModalOverlay(modal, baseView)
	}

	if i.reviewState.EditingMode == EditingModeFacts && i.reviewState.factSuggestionModal != nil {
		rendered := i.reviewState.factSuggestionModal.View()
		modal := &behaviors.StaticViewModel{Content: rendered}
		return behaviors.RenderModalOverlay(modal, baseView)
	}

	modalContent := i.getEditingModalContent()
	if modalContent != nil {
		return i.renderModalOverlay(baseView, modalContent)
	}

	return baseView
}

func (i *Intent) breadcrumbs() []string {
	switch i.currentState {
	case StateForm:
		return []string{mainMenuTitle, captureEventTitle, "Form"}
	case StateReview:
		return []string{mainMenuTitle, captureEventTitle, "Review Enrichment"}
	default:
		return []string{mainMenuTitle, captureEventTitle}
	}
}

// Result returns the intent's outcome as a type-erased IntentResult.
//
// Returns:
//   - A fully initialized intents.IntentResult[interface{}] ready for use.
//
// Side effects:
//   - None.
func (i *Intent) Result() *intents.IntentResult[interface{}] {
	if i.result == nil {
		return nil
	}
	return &intents.IntentResult[interface{}]{
		Status:   i.result.Status,
		Data:     i.result.Data,
		Error:    i.result.Error,
		Metadata: i.result.Metadata,
	}
}
