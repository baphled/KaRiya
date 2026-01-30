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

// NewIntent creates a new CaptureEvent intent from the given context.
// It validates the context and returns an error if required fields are missing.
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
		state: &Model{
			context:      ctx,
			currentState: StateChooseStrategy,
			captureForm:  formModel,
			reviewState: &ReviewInferredEventState{
				AcceptedBursts: make([]*career.Burst, 0),
				AcceptedFacts:  make([]*career.Fact, 0),
				RejectedItems:  make(map[string]string),
			},
			result: &Result{
				AcceptedFields: make(map[string]bool),
				RejectedFields: make(map[string]string),
			},
		},
		active:     true,
		useScreens: true,
	}, nil
}

// Init is called when the intent is activated.
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

// Update processes a message in the intent.
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

		i.state.reviewState = &ReviewInferredEventState{
			Event:          msg.Event,
			InferredBursts: make([]*career.Burst, 0),
			InferredFacts:  make([]*career.Fact, 0),
			AcceptedBursts: make([]*career.Burst, 0),
			AcceptedFacts:  make([]*career.Fact, 0),
			RejectedItems:  make(map[string]string),
		}
		i.state.currentState = StateReview
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
			i.state.currentState = StateReview

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
		switch intents.HandleGlobalKeys(keyMsg) {
		case intents.KeyQuit:
			return tea.Quit
		case intents.KeyHelp:
			i.BaseIntent.ToggleHelp()
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

// View renders the intent's current state. When screens are active, it delegates
// to the active screen and overlays any modal. Otherwise it falls back to the
// legacy breadcrumb-based view.
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

// Result returns the intent's outcome as a type-erased IntentResult.
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
