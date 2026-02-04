package captureevent

import (
	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/screens"
	captureScreens "github.com/baphled/kariya/internal/cli/screens/capture"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
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
func NewIntent(ctx *IntentContext) (*Intent, error) {
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
	breadcrumbs := []string{"Main Menu", "Capture Event"}
	i.activeScreen = captureScreens.NewStrategySelectScreen(breadcrumbs)

	termInfo := i.GetTerminalInfo()
	width, height := 120, 40
	if termInfo != nil && termInfo.Width > 0 && termInfo.Height > 0 {
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

	switch msg := msg.(type) {
	case SubmitCompleteMsg:
		i.submitModal = feedback.NewSuccessModal("Event saved!")
		return i.submitModal.Init()

	case SubmitErrorMsg:
		i.submitModal = feedback.NewErrorModal("Save Failed", msg.Message)
		return nil

	case feedback.ModalAutoDismissMsg:
		return func() tea.Msg { return DismissModalMsg{} }

	case DismissModalMsg:
		if i.submitModal != nil {
			i.submitModal = nil
			i.postSaveReview = true
			i.currentState = StateReview

			breadcrumbs := []string{"Main Menu", "Capture Event", "Review Enrichment"}
			i.activeScreen = captureScreens.NewEventReviewScreen(
				breadcrumbs,
				i.reviewState.Event,
				i.reviewState.InferredBursts,
				i.reviewState.InferredFacts,
			)

			termInfo := i.GetTerminalInfo()
			width, height := 120, 40
			if termInfo != nil && termInfo.Width > 0 && termInfo.Height > 0 {
				width = termInfo.Width
				height = termInfo.Height
			}

			i.activeScreen.SetTerminalInfo(width, height)
			i.activeScreen.SetTheme(i.Theme())
			i.activeScreen.SetLogo(i.GetLogo(), i.GetLogoSpacing())
		}
		return nil

	case SubmitMsg:
		if i.currentState == StateForm && msg.Err == nil && msg.Event != nil {
			formData := forms.GetCaptureEventFormData(msg.Event)
			formData.SubmitConfirmed = true
			return i.handleScreenResult(&screens.SubmitResult{FormData: formData})
		}
		return nil
	}

	if i.submitModal != nil {
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

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch intents.HandleGlobalKeys(keyMsg) {
		case intents.KeyQuit:
			return tea.Quit
		case intents.KeyHelp:
			i.ToggleHelp()
			return nil
		}
	}

	if i.activeScreen != nil {
		if i.reviewState != nil && i.reviewState.EditingMode != EditingModeNone {
			return i.updateEditingModal(msg)
		}

		cmd, result := i.activeScreen.Update(msg)

		if result != nil {
			return i.handleScreenResult(result)
		}

		return cmd
	}

	return nil
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

	if i.activeScreen == nil {
		return "No active screen"
	}

	baseView := i.activeScreen.View()

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
		modalContent := i.getEditingModalContent()
		if modalContent != nil {
			return i.renderModalOverlay(baseView, modalContent)
		}
	}

	return baseView
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
