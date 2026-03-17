package configure

import (
	"github.com/baphled/kariya/internal/tui/intents"
	"github.com/baphled/kariya/internal/ui/themes"
	"github.com/baphled/kariya/internal/ui/uikit/feedback"
	tea "github.com/charmbracelet/bubbletea"
)

// Compile-time check that Intent implements the intents.Intent interface.
var _ intents.Intent = (*Intent)(nil)

// NewIntent constructs a fully configured ConfigureSystem intent.
//
// Expected: ctx must contain a valid Config.
// Returns: A fully initialised Intent and nil error, or nil and an error.
//
// Side effects: None.
func NewIntent(ctx *IntentValidator) (*Intent, error) {
	if err := ctx.Validate(); err != nil {
		return nil, err
	}

	baseIntent := intents.NewBaseIntent()
	baseIntent.SetThemeManager(themes.NewThemeManager())

	intent := &Intent{
		BaseIntent:    baseIntent,
		context:       ctx,
		state:         ConfigStateSelectDomain,
		cfg:           ctx.Cfg,
		settings:      ctx.Settings,
		active:        false,
		modalRegistry: intents.NewModalRegistry(),
	}

	return intent, nil
}

// Init activates the intent, opens the settings modal, and returns its initial command.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (i *Intent) Init() tea.Cmd {
	i.active = true
	i.openSettingsModal()

	if i.settingsModal != nil {
		return i.settingsModal.Init()
	}
	return nil
}

// Update drives the configuration workflow state machine.
//
// Expected: msg must be valid.
// Returns: A tea.Cmd value.
//
// Side effects: None.
func (i *Intent) Update(msg tea.Msg) tea.Cmd {
	if !i.active {
		return nil
	}

	i.handleWindowSizeMsg(msg)

	if i.handleKeyMsg(msg) {
		return nil
	}

	if cmd := i.handleAsyncCompletion(msg); cmd != nil {
		return cmd
	}

	return i.routeToActiveComponent(msg)
}

func (i *Intent) handleWindowSizeMsg(msg tea.Msg) {
	wsMsg, ok := msg.(tea.WindowSizeMsg)
	if !ok {
		return
	}

	if i.settingsModal != nil {
		i.settingsModal.SetDimensions(wsMsg.Width, wsMsg.Height)
	}
}

func (i *Intent) handleKeyMsg(msg tea.Msg) bool {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return false
	}

	if intents.HandleGlobalKeys(keyMsg) == intents.KeyHelp {
		i.ToggleHelp()
		return true
	}

	if keyMsg.String() == "q" && i.settingsModal == nil {
		i.setCancelled()
		return true
	}

	return false
}

func (i *Intent) routeToActiveComponent(msg tea.Msg) tea.Cmd {
	// Handle modals in priority order (highest to lowest):
	// 1. resultModal (success/error feedback)
	// 2. savingModal (loading spinner)
	// 3. settingsModal (form modal)

	if i.resultModal != nil {
		return i.updateResultModal(msg)
	}
	if i.savingModal != nil {
		return i.updateSavingModal(msg)
	}
	if i.settingsModal != nil {
		return i.updateSettingsModal(msg)
	}

	return nil
}

// View composes the visible UI by layering any active modal overlay.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (i *Intent) View() string {
	if !i.active {
		return ""
	}

	view := i.CreateViewWithBreadcrumbs("Main Menu", "Configure System", i.getStateName())

	view.WithContent("")
	view.WithHelp(i.getContextHelp()).WithFooterSeparator(true)

	width, height := i.getTerminalDimensions()

	if overlay := i.activeModalOverlay(width, height); overlay != nil {
		view.ShowModalOverlay(overlay)
	}

	return view.Render()
}

func (i *Intent) activeModalOverlay(width, height int) interface{ Render(int, int) string } {
	if i.settingsModal != nil {
		return settingsModalAdapter{i.settingsModal, width, height}
	}
	if i.savingModal != nil {
		return i.savingModal
	}
	if i.resultModal != nil {
		return i.resultModal
	}
	return nil
}

// Result returns the intent's result if it has completed, or nil if still active.
//
// Returns:
//   - A fully initialized intents.IntentResult[interface{}] ready for use.
//
// Side effects:
//   - None.
func (i *Intent) Result() *intents.IntentResult[interface{}] {
	if i.active {
		return nil
	}

	if i.result != nil {
		return i.result
	}

	if i.configResult != nil {
		return &intents.IntentResult[interface{}]{
			Status: intents.Completed,
			Data:   i.configResult,
		}
	}

	return &intents.IntentResult[interface{}]{
		Status: intents.Cancelled,
	}
}

// SetState transitions the intent to a new state with appropriate modal creation.
//
// Expected:
//   - config must be a valid configuration object.
//
// Side effects:
//   - None.
func (i *Intent) SetState(state ConfigurationState) {
	i.clearAllModals()
	i.state = state

	switch state {
	case ConfigStateSelectDomain, ConfigStateEditSettings, ConfigStateReviewChanges, ConfigStateConfirm:
		i.openSettingsModal()
	case ConfigStateSaving:
		i.savingModal = feedback.NewLoadingModal("Saving...", false)
	case ConfigStateComplete:
		i.configResult = &SystemResult{Success: true}
		i.resultModal = feedback.NewSuccessModal("Complete")
	case ConfigStateFailed:
		i.resultModal = feedback.NewErrorModal("Failed", "Error occurred")
	}
}
