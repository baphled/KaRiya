package app

import (
	"context"

	"github.com/baphled/kariya/internal/cli/bootstrap"
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/intents/captureevent"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/cli/terminal"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/display"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/logger"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	tea "github.com/charmbracelet/bubbletea"
)

// ModelOption is a functional option for configuring the Model.
type ModelOption func(*modelOptions)

type modelOptions struct {
	registrar IntentRegistrar
}

// WithIntentRegistrar sets a custom intent registrar (useful for testing).
//
// Expected:
//   - intentregistrar must be valid.
//
// Returns:
//   - A ModelOption value.
//
// Side effects:
//   - None.
func WithIntentRegistrar(registrar IntentRegistrar) ModelOption {
	return func(o *modelOptions) {
		o.registrar = registrar
	}
}

// NewModel creates and initializes a new application model.
// Bootstrap must be run first to handle onboarding and service initialization.
//
// Expected:
//   - cliservice must be a valid CLIEventService instance.
//   - careerservice must be a valid career Service instance.
//   - bootstrapresult must be a valid bootstrap.Result with initialized services.
//
// Returns:
//   - A fully initialized Model ready for use.
//
// Side effects:
//   - Registers all intents with the router.
//   - Initializes logger.
func NewModel(
	cliService *service.CLIEventService,
	careerService *careerservice.Service,
	bootstrapResult *bootstrap.Result,
	opts ...ModelOption,
) *Model {
	ctx := context.Background()
	log := logger.DefaultLogger()

	// Apply options.
	options := &modelOptions{}
	for _, opt := range opts {
		opt(options)
	}

	// Initialize intent router.
	router := intents.NewDefaultIntentRouter()

	// Create skill inference service.
	skillInferenceService := skillinference.NewSkillInferenceService(
		careerService.GetSkillRepository(),
		careerService.GetEventRepository(),
	)

	// Use provided registrar or create default.
	registrar := options.registrar
	if registrar == nil {
		registrar = NewDefaultIntentRegistrar(&RegistrarConfig{
			CLIService:            cliService,
			CareerService:         careerService,
			SkillInferenceService: skillInferenceService,
			Log:                   log,
			CVGenService:          bootstrapResult.Services.CVGenService,
			CVExportService:       bootstrapResult.Services.CVExportService,
		})
	}

	// Register all intents.
	if err := registrar.RegisterAll(ctx, router); err != nil {
		log.Error("Failed to register intents: %v", err)
	}

	// Create menu items for all intents.
	menuItems := []MenuItem{
		{Name: "Capture Event", Intent: "capture_event", Help: "Record a new career event"},
		{Name: "Browse Timeline", Intent: "browse_timeline", Help: "View your career events"},
		{Name: "Manage Skills", Intent: "manage_skills", Help: "Manage your skills"},
		{Name: "Generate CV", Intent: "generate_cv", Help: "Create a new CV"},
		{Name: "Configure System", Intent: "configure_system", Help: "Manage settings"},
		{Name: "Manage Bursts", Intent: "burst_management", Help: "Organize career bursts"},
		{Name: "Manage Facts", Intent: "fact_management", Help: "Review extracted facts"},
	}

	// Create ASCII logo with animation.
	logo := display.NewLogo(true, 80)

	// Share logo with intent router so all intents can use it.
	router.SetLogo(logo)

	return &Model{
		cliService:        cliService,
		careerService:     careerService,
		logger:            log,
		intentRouter:      router,
		configManager:     bootstrapResult.Services.ConfigManager,
		cvGenService:      bootstrapResult.Services.CVGenService,
		cvExportService:   bootstrapResult.Services.CVExportService,
		theme:             themes.NewDefaultTheme(),
		state:             StateMenu,
		selectedMenuIndex: 0,
		menuItems:         menuItems,
		logo:              logo,
		terminalInfo:      terminal.NewInfo(),
		ctx:               ctx,
		width:             80,
		height:            24,
	}
}

// Init initializes the model.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (m *Model) Init() tea.Cmd {
	cmds := []tea.Cmd{
		tea.WindowSize(),
		m.logo.Init(),
	}

	// Handle initial screen/mode settings.
	if m.initialScreen == ListScreen {
		// Navigate directly to browse_timeline intent.
		cmd, err := m.intentRouter.ActivateIntent("browse_timeline", make(map[string]interface{}))
		if err != nil {
			m.logger.Error("Failed to navigate to initial list screen: %v", err)
		} else {
			m.state = StateIntent
			cmds = append(cmds, cmd)
		}
	} else if m.initialCaptureMode != "" {
		// Navigate directly to capture_event intent with specified mode.
		cmd, err := m.intentRouter.ActivateIntent("capture_event", map[string]interface{}{
			"mode": m.initialCaptureMode,
		})
		if err != nil {
			m.logger.Error("Failed to navigate to initial capture mode '%s': %v", m.initialCaptureMode, err)
		} else {
			m.state = StateIntent
			cmds = append(cmds, cmd)
		}
	}

	return tea.Batch(cmds...)
}

// Update handles messages and state transitions.
//
// Expected:
//   - msg must be a valid tea.Msg.
//
// Returns:
//   - Updated Model and command to execute.
//
// Side effects:
//   - May update internal state, navigate between intents.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Update logo animation if in menu state.
	if m.state == StateMenu {
		if _, ok := msg.(display.TickMsg); ok {
			updatedLogo, cmd := m.logo.Update(msg)
			//nolint:errcheck // Type assertion is safe - Logo.Update always returns *display.Logo.
			m.logo = updatedLogo.(*display.Logo)
			return m, cmd
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)

	case tea.WindowSizeMsg:
		return m.handleWindowResize(msg)

	case IntentCompletedMsg:
		return m.handleIntentCompleted()

	default:
		return m.handleDefaultMsg(msg)
	}
}

// handleKeyMsg processes keyboard input.
func (m *Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle info modal dismissal first (highest priority).
	if m.infoModal != nil && m.infoModal.IsVisible() {
		if m.infoModal.Update(msg) {
			m.infoModal = nil
		}
		return m, nil
	}

	switch msg.String() {
	case keyCtrlC:
		return m, tea.Quit
	case keyQuit:
		if m.state == StateMenu {
			return m, tea.Quit
		}
	case keyHelp:
		m.showingHelp = !m.showingHelp
		return m, nil
	}

	// Route to appropriate handler based on state.
	if m.state == StateMenu {
		return m.handleMenuInput(msg)
	} else if m.state == StateIntent {
		return m.handleIntentInput(msg)
	}

	return m, nil
}

// handleWindowResize processes window resize events.
func (m *Model) handleWindowResize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.width = msg.Width
	m.height = msg.Height

	m.terminalInfo.Update(msg)
	m.intentRouter.UpdateTerminalInfo(m.terminalInfo)
	m.logo.SetWidth(msg.Width)

	return m, tea.ClearScreen
}

// handleIntentCompleted processes intent completion.
func (m *Model) handleIntentCompleted() (tea.Model, tea.Cmd) {
	m.state = StateMenu
	m.selectedMenuIndex = 0
	return m, nil
}

// handleDefaultMsg processes all other message types.
func (m *Model) handleDefaultMsg(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Check for RequestEditEventMsg.
	if editMsg, ok := msg.(intents.RequestEditEventMsg); ok {
		return m.handleEditEventRequest(editMsg)
	}

	// Route to active intent.
	if m.state == StateIntent {
		cmd, result := m.intentRouter.HandleMessage(msg)

		if result != nil {
			m.state = StateMenu
			m.selectedMenuIndex = 0
			return m, tea.Batch(
				cmd,
				func() tea.Msg { return IntentCompletedMsg{} },
			)
		}

		return m, cmd
	}

	return m, nil
}

// handleEditEventRequest handles cross-intent navigation for editing events.
func (m *Model) handleEditEventRequest(editMsg intents.RequestEditEventMsg) (tea.Model, tea.Cmd) {
	captureCtx := &captureevent.IntentContext{
		CaptureStrategy: "manual",
		PreviousEvent:   editMsg.Event,
		Metadata:        make(map[string]string),
		CLIEventService: m.cliService,
		CareerService:   m.careerService,
	}

	// #nosec G104 -- RegisterIntent only errors on duplicate registration which cannot happen here
	m.intentRouter.RegisterIntent("capture_event_edit", func() intents.Intent {
		intent, err := captureevent.NewIntent(captureCtx)
		if err != nil {
			m.logger.Error("Failed to create CaptureEvent intent for editing: %v", err)
			return nil
		}
		return intent
	})

	cmd, err := m.intentRouter.ActivateIntent("capture_event_edit", make(map[string]interface{}))
	if err != nil {
		m.logger.Error("Failed to activate CaptureEvent for editing: %v", err)
		return m, nil
	}
	m.state = StateIntent
	return m, cmd
}

// View renders the current screen.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *Model) View() string {
	if m.showingHelp {
		return m.renderHelpScreen()
	}

	if m.state == StateMenu {
		menuView := m.viewMenu()

		if m.infoModal != nil && m.infoModal.IsVisible() {
			modalView := m.infoModal.View()
			return primitives.CenterInTerminal(modalView, m.width, m.height)
		}

		return menuView
	} else if m.state == StateIntent {
		activeIntent := m.intentRouter.GetActiveIntent()
		if activeIntent != nil {
			return activeIntent.View()
		}
		return "No active intent"
	}

	return ""
}
