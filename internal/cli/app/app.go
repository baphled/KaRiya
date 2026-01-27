package app

import (
	"context"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/cli/terminal"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/display"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/logger"
	careerservice "github.com/baphled/kariya/internal/service/career"
	cv "github.com/baphled/kariya/internal/service/career/cv"
	tea "github.com/charmbracelet/bubbletea"
)

// NewModel creates and initializes a new application model
func NewModel(cliService *service.CLIEventService, careerService *careerservice.Service) *Model {
	ctx := context.Background()
	log := logger.DefaultLogger()

	// Load application config (for profile check)
	appCfg, err := config.LoadConfig()
	if err != nil {
		log.Error("Failed to load config: %v, using defaults", err)
		appCfg = config.DefaultConfig()
	}

	// Initialize CV services
	configMgr := initConfigManager(log)
	cvGenService := initCVGenerationService(careerService, configMgr, &appCfg.Scoring, log)
	cvExportService := cv.NewExportService(log)

	// Initialize intent router
	router := intents.NewDefaultIntentRouter()
	registerAllIntents(router, cliService, careerService, log, ctx, cvGenService, cvExportService)

	// Create menu items for all intents
	menuItems := []MenuItem{
		{Name: "Capture Event", Intent: "capture_event", Help: "Record a new career event"},
		{Name: "Browse Timeline", Intent: "browse_timeline", Help: "View your career events"},
		{Name: "Manage Skills", Intent: "manage_skills", Help: "Manage your skills"},
		{Name: "Generate CV", Intent: "generate_cv", Help: "Create a new CV"},
		{Name: "Configure System", Intent: "configure_system", Help: "Manage settings"},
		{Name: "Manage Bursts", Intent: "burst_management", Help: "Organize career bursts"},
		{Name: "Manage Facts", Intent: "fact_management", Help: "Review extracted facts"},
	}

	// Create ASCII logo with animation
	logo := display.NewLogo(true, 80)

	// Share logo with intent router so all intents can use it
	router.SetLogo(logo)

	// Determine initial state - show onboarding if required profile fields are missing
	// Both Name and Email are required for CV generation
	initialState := StateMenu
	var onboardingWizard *components.OnboardingWizardModal
	if appCfg.Profile.Name == "" || appCfg.Profile.Email == "" {
		initialState = StateOnboarding
		onboardingWizard = components.NewOnboardingWizardModalWithConfig(80, 24, &appCfg.Profile)
	}

	return &Model{
		cliService:        cliService,
		careerService:     careerService,
		logger:            log,
		intentRouter:      router,
		configManager:     configMgr,
		cvGenService:      cvGenService,
		cvExportService:   cvExportService,
		theme:             themes.NewDefaultTheme(),
		state:             initialState,
		selectedMenuIndex: 0,
		menuItems:         menuItems,
		logo:              logo,
		terminalInfo:      terminal.NewInfo(),
		ctx:               ctx,
		width:             80,
		height:            24,
		onboardingWizard:  onboardingWizard,
		appConfig:         appCfg,
	}
}

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	cmds := []tea.Cmd{
		tea.WindowSize(),
		m.logo.Init(),
	}

	// Initialize onboarding wizard if in onboarding state
	if m.state == StateOnboarding && m.onboardingWizard != nil {
		cmds = append(cmds, m.onboardingWizard.Init())
		return tea.Batch(cmds...)
	}

	// Handle initial screen/mode settings (skip if onboarding is active)
	if m.initialScreen == ListScreen {
		// Navigate directly to browse_timeline intent
		cmd, err := m.intentRouter.ActivateIntent("browse_timeline", make(map[string]interface{}))
		if err == nil {
			m.state = StateIntent
			cmds = append(cmds, cmd)
		}
	} else if m.initialCaptureMode != "" {
		// Navigate directly to capture_event intent with specified mode
		cmd, err := m.intentRouter.ActivateIntent("capture_event", map[string]interface{}{
			"mode": m.initialCaptureMode,
		})
		if err == nil {
			m.state = StateIntent
			cmds = append(cmds, cmd)
		}
	}

	return tea.Batch(cmds...)
}

// Update handles messages - FIXED: Using correct Bubble Tea v1.3.10 signature
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
		// Handle info modal dismissal first (highest priority)
		// BUG-004: Info modal shows when user tries to generate CV without events
		if m.infoModal != nil && m.infoModal.IsVisible() {
			if m.infoModal.Update(msg) {
				m.infoModal = nil // Dismissed
			}
			return m, nil // Consume all keys when modal is showing
		}

		switch msg.String() {
		case keyCtrlC:
			return m, tea.Quit
		case keyQuit:
			if m.state == StateMenu {
				return m, tea.Quit
			}
		case keyHelp:
			// Toggle help screen (only ? key, not 'h' which is vim-style left navigation).
			m.showingHelp = !m.showingHelp
			return m, nil
		}

		// Route messages to appropriate handler based on state
		// Note: We don't intercept escape here - intents handle their own back navigation
		// per TUI Standards (intermediate states go back one state, root states cancel intent)
		if m.state == StateOnboarding {
			return m.handleOnboardingInput(msg)
		} else if m.state == StateMenu {
			return m.handleMenuInput(msg)
		} else if m.state == StateIntent {
			return m.handleIntentInput(msg)
		}

	case tea.WindowSizeMsg:
		// Update dimensions
		m.width = msg.Width
		m.height = msg.Height

		// Update terminal info
		m.terminalInfo.Update(msg)

		// Propagate terminal info to intent router
		m.intentRouter.UpdateTerminalInfo(m.terminalInfo)

		// Update logo width for centering
		m.logo.SetWidth(msg.Width)

		// Update onboarding wizard dimensions if active
		// IMPORTANT: Execute returned cmd - wizard rebuilds form and returns form.Init()
		var wizardCmd tea.Cmd
		if m.onboardingWizard != nil {
			wizardCmd = m.onboardingWizard.Update(msg)
		}

		// Clear screen to prevent artifacts on resize
		// Batch with wizard cmd if present
		if wizardCmd != nil {
			return m, tea.Batch(tea.ClearScreen, wizardCmd)
		}
		return m, tea.ClearScreen

	case IntentCompletedMsg:
		// Handle result from completed intent
		m.state = StateMenu
		m.selectedMenuIndex = 0
		return m, nil

	default:
		// Check for RequestEditEventMsg before routing to intent
		if editMsg, ok := msg.(intents.RequestEditEventMsg); ok {
			// User wants to edit an event - activate CaptureEvent intent with PreviousEvent
			captureCtx := &intents.CaptureEventContext{
				CaptureStrategy: "manual",
				PreviousEvent:   editMsg.Event,
				Metadata:        make(map[string]string),
				CLIEventService: m.cliService,
				CareerService:   m.careerService,
			}

			// Temporarily register the edit intent
			//nolint:errcheck // RegisterIntent only errors on duplicate registration which cannot happen here
			m.intentRouter.RegisterIntent("capture_event_edit", func() intents.Intent {
				intent, err := intents.NewCaptureEventIntent(captureCtx)
				if err != nil {
					m.logger.Error("Failed to create CaptureEvent intent for editing: %v", err)
					return nil
				}
				return intent
			})

			// Activate the edit intent
			cmd, err := m.intentRouter.ActivateIntent("capture_event_edit", make(map[string]interface{}))
			if err != nil {
				m.logger.Error("Failed to activate CaptureEvent for editing: %v", err)
				return m, nil
			}
			m.state = StateIntent
			return m, cmd
		}

		// Route all other messages to the onboarding wizard (e.g., huh internal messages like nextGroupMsg)
		// This is essential for huh forms to advance between groups
		if m.state == StateOnboarding && m.onboardingWizard != nil {
			cmd := m.onboardingWizard.Update(msg)

			// Check if wizard completed
			if m.onboardingWizard.IsCompleted() {
				// Get the profile config from the wizard
				profileCfg := m.onboardingWizard.GetProfileConfig()
				if profileCfg != nil {
					// Update the app config with the new profile
					m.appConfig.Profile = *profileCfg

					// Save the config
					if err := config.SaveConfig(m.appConfig); err != nil {
						m.logger.Error("Failed to save config: %v", err)
					} else {
						m.logger.Info("Profile saved successfully")
					}
				}

				// Transition to main menu
				m.state = StateMenu
				m.onboardingWizard = nil
				return m, nil
			}

			return m, cmd
		}

		// Route all other messages to the active intent (e.g., SubmitMsg from form commands)
		if m.state == StateIntent {
			cmd, result := m.intentRouter.HandleMessage(msg)

			// Check if intent has completed
			if result != nil {
				m.state = StateMenu
				m.selectedMenuIndex = 0
				// Return command that will trigger the intent completed message
				return m, tea.Batch(
					cmd,
					func() tea.Msg { return IntentCompletedMsg{} },
				)
			}

			return m, cmd
		}
	}

	return m, nil
}

// View renders the current screen
func (m *Model) View() string {
	// Show help overlay if active
	if m.showingHelp {
		return m.renderHelpScreen()
	}

	if m.state == StateOnboarding {
		return m.viewOnboarding()
	} else if m.state == StateMenu {
		menuView := m.viewMenu()

		// BUG-004: Overlay info modal if active (e.g., empty state warning).
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
