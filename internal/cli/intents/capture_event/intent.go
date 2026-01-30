package capture_event

import (
	"fmt"
	"time"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/models"
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

	formModel := models.NewCaptureForm(ctx.CLIEventService)
	base := intents.NewBaseIntent()

	return &Intent{
		BaseIntent:   base,
		context:      ctx,
		eventService: ctx.CLIEventService,
		active:       true,
		useScreens:   true,
		currentState: StateChooseStrategy,
		captureForm:  formModel,
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
//   - A tea.Cmd to initialise the active screen, or nil.
//
// Side effects:
//   - Creates and configures the strategy selection screen when screens are enabled.
//   - Falls back to initialising the legacy form when screens are disabled.
func (i *Intent) Init() tea.Cmd {
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
	case FormSubmittedMsg:
		if msg.Event == nil {
			i.setFailed("INVALID_FORM", "Form submission with nil event", nil)
			return nil
		}

		if err := msg.Event.Validate(); err != nil {
			i.setFailed("VALIDATION_ERROR", fmt.Sprintf("Form validation failed: %v", err), err)
			return nil
		}

		i.reviewState = &ReviewInferredEventState{
			Event:          msg.Event,
			InferredBursts: make([]*career.Burst, 0),
			InferredFacts:  make([]*career.Fact, 0),
			AcceptedBursts: make([]*career.Burst, 0),
			AcceptedFacts:  make([]*career.Fact, 0),
			RejectedItems:  make(map[string]string),
		}
		i.currentState = StateReview
		i.activeScreen = nil
		return nil

	case SubmitCompleteMsg:
		i.submitModal = feedback.NewSuccessModal("Event saved!")
		return tea.Tick(2*time.Second, func(_ time.Time) tea.Msg {
			return DismissModalMsg{}
		})

	case SubmitErrorMsg:
		i.submitModal = feedback.NewErrorModal("Save Failed", msg.Message)
		return nil

	case DismissModalMsg:
		if i.submitModal != nil {
			i.submitModal = nil
			i.postSaveReview = true
			i.currentState = StateReview

			if i.useScreens {
				breadcrumbs := []string{"Main Menu", "Capture Event", "Review Enrichment"}
				i.activeScreen = captureScreens.NewEventReviewScreen(
					breadcrumbs,
					i.reviewState.Event,
					i.reviewState.InferredBursts,
					i.reviewState.InferredFacts,
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

	if i.submitModal != nil {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.Type == tea.KeyEsc {
				i.submitModal = nil
				return nil
			}
			return nil
		case feedback.ModalSpinnerTickMsg:
			if i.submitModal.Type == feedback.ModalLoading {
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
			i.BaseIntent.ToggleHelp()
			return nil
		}
	}

	if i.useScreens && i.activeScreen != nil {
		if i.reviewState != nil && i.reviewState.EditingMode != EditingModeNone {
			return i.updateEditingModal(msg)
		}

		cmd, result := i.activeScreen.Update(msg)

		if result != nil {
			return i.handleScreenResult(result)
		}

		return cmd
	}

	switch i.currentState {
	case StateChooseStrategy:
		return i.updateChooseStrategy(msg)

	case StateForm:
		return i.updateCaptureForm(msg)

	case StateReview:
		return i.updateReviewInferredEvent(msg)

	case StateSubmit:
		return i.updateSubmit(msg)

	default:
		return nil
	}
}

// View renders the intent's current visual state.
//
// Returns:
//   - A string containing the full terminal output for the current frame.
//
// When screens are enabled, View delegates to the active screen and overlays
// any visible modal. When screens are disabled, it uses the legacy
// breadcrumb-based view with inline content rendering.
func (i *Intent) View() string {
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

		if i.submitModal != nil {
			modalContent := i.submitModal.Render(width, height)
			return i.overlayModal(baseView, modalContent, width, height)
		}

		if i.reviewState != nil && i.reviewState.EditingMode != EditingModeNone {
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

	if i.intentError != nil {
		i.SetError(i.intentError)
	}

	content := i.getStateContent()
	view.WithContent(content)

	help := i.getContextHelp()
	view.WithHelp(help).WithFooterSeparator(true)

	return view.Render()
}

// Result returns the intent's outcome as a type-erased IntentResult.
//
// Returns:
//   - An IntentResult[interface{}] wrapping the typed result, or nil if
//     the intent has not yet completed.
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
