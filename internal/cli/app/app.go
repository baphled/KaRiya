package app

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/cli/terminal"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	cv "github.com/baphled/kariya/internal/service/career/cv"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// AppState represents the current state of the application
type AppState string

const (
	StateMenu   AppState = "menu"
	StateIntent AppState = "intent"
)

// Model is the root Bubble Tea model for the KaRiya application
type Model struct {
	// Core services
	cliService      *service.CLIEventService
	careerService   *careerservice.Service
	logger          *logger.Logger
	intentRouter    *intents.DefaultIntentRouter
	configManager   cv.ConfigManager
	cvGenService    cv.CVGenerationService
	cvExportService *cv.ExportService

	// UI state
	state       AppState
	width       int
	height      int
	showingHelp bool

	// Menu state
	selectedMenuIndex int
	menuItems         []MenuItem
	logo              *components.ASCIILogo

	// Terminal info for responsive rendering
	terminalInfo *terminal.Info

	// Context for intent creation
	ctx context.Context
}

// MenuItem represents a menu option
type MenuItem struct {
	Name   string
	Intent string
	Help   string
}

// NewModel creates and initializes a new application model
func NewModel(cliService *service.CLIEventService, careerService *careerservice.Service) *Model {
	ctx := context.Background()
	log := logger.DefaultLogger()

	// Initialize CV services
	configMgr := initConfigManager(log)
	cvGenService := initCVGenerationService(careerService, configMgr, log)
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
		{Name: "Export Artifact", Intent: "export_artifact", Help: "Export CV or data"},
		{Name: "Configure System", Intent: "configure_system", Help: "Manage settings"},
		{Name: "Manage Bursts", Intent: "burst_management", Help: "Organize career bursts"},
		{Name: "Manage Facts", Intent: "fact_management", Help: "Review extracted facts"},
		{Name: "Import Data", Intent: "import_wizard", Help: "Import from CSV"},
		{Name: "Edit Metadata", Intent: "metadata_editor", Help: "Update metadata"},
		{Name: "Bulk Operations", Intent: "bulk_operations", Help: "Perform bulk actions"},
	}

	// Create ASCII logo with animation
	logo := components.NewASCIILogo(true, 80)

	// Share logo with intent router so all intents can use it
	router.SetLogo(logo)

	return &Model{
		cliService:        cliService,
		careerService:     careerService,
		logger:            log,
		intentRouter:      router,
		configManager:     configMgr,
		cvGenService:      cvGenService,
		cvExportService:   cvExportService,
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

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	// Request initial terminal size and initialize logo animation
	return tea.Batch(
		tea.WindowSize(),
		m.logo.Init(),
	)
}

// Update handles messages - FIXED: Using correct Bubble Tea v1.3.10 signature
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Update logo animation if in menu state
	if m.state == StateMenu {
		if _, ok := msg.(components.TickMsg); ok {
			updatedLogo, cmd := m.logo.Update(msg)
			m.logo = updatedLogo.(*components.ASCIILogo)
			return m, cmd
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			if m.state == StateMenu {
				return m, tea.Quit
			}
		case "?":
			// Toggle help screen (only ? key, not 'h' which is vim-style left navigation)
			m.showingHelp = !m.showingHelp
			return m, nil
		}

		// Route messages to appropriate handler based on state
		// Note: We don't intercept escape here - intents handle their own back navigation
		// per TUI Standards (intermediate states go back one state, root states cancel intent)
		if m.state == StateMenu {
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

		// Clear screen to prevent artifacts on resize
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
			_ = m.intentRouter.RegisterIntent("capture_event_edit", func() intents.Intent {
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

	if m.state == StateMenu {
		return m.viewMenu()
	} else if m.state == StateIntent {
		activeIntent := m.intentRouter.GetActiveIntent()
		if activeIntent != nil {
			return activeIntent.View()
		}
		return "No active intent"
	}
	return ""
}

// renderHelpScreen renders the keyboard reference help screen
func (m *Model) renderHelpScreen() string {
	// Build help content
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("12")).
		Bold(true)

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("14")).
		Bold(true)

	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("10"))

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	borderStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	var lines []string

	// Title
	lines = append(lines, "")
	lines = append(lines, titleStyle.Render("  KaRiya Keyboard Reference"))
	lines = append(lines, borderStyle.Render("  ════════════════════════════════════════════════════════"))
	lines = append(lines, "")

	// Global shortcuts
	lines = append(lines, headerStyle.Render("  Global Shortcuts"))
	lines = append(lines, borderStyle.Render("  ──────────────────────────────────────────────────────────"))
	lines = append(lines, keyStyle.Render("  ?")+"         "+descStyle.Render("Toggle this help screen"))
	lines = append(lines, keyStyle.Render("  q")+"         "+descStyle.Render("Quit application (from menu)"))
	lines = append(lines, keyStyle.Render("  Ctrl+C")+"    "+descStyle.Render("Force quit application"))
	lines = append(lines, keyStyle.Render("  Esc")+"       "+descStyle.Render("Go back / Cancel / Return to menu"))
	lines = append(lines, "")

	// Navigation
	lines = append(lines, headerStyle.Render("  Navigation"))
	lines = append(lines, borderStyle.Render("  ──────────────────────────────────────────────────────────"))
	lines = append(lines, keyStyle.Render("  ↑/k")+"       "+descStyle.Render("Move up / Previous item"))
	lines = append(lines, keyStyle.Render("  ↓/j")+"       "+descStyle.Render("Move down / Next item"))
	lines = append(lines, keyStyle.Render("  ←/h")+"       "+descStyle.Render("Move left / Previous"))
	lines = append(lines, keyStyle.Render("  →/l")+"       "+descStyle.Render("Move right / Next"))
	lines = append(lines, keyStyle.Render("  Enter")+"     "+descStyle.Render("Confirm / Select"))
	lines = append(lines, keyStyle.Render("  Space")+"     "+descStyle.Render("Toggle selection"))
	lines = append(lines, "")

	// Form shortcuts
	lines = append(lines, headerStyle.Render("  Forms & Input"))
	lines = append(lines, borderStyle.Render("  ──────────────────────────────────────────────────────────"))
	lines = append(lines, keyStyle.Render("  Tab")+"       "+descStyle.Render("Next field"))
	lines = append(lines, keyStyle.Render("  Shift+Tab")+" "+descStyle.Render("Previous field"))
	lines = append(lines, keyStyle.Render("  Ctrl+O")+"    "+descStyle.Render("Toggle optional fields"))
	lines = append(lines, "")

	// List shortcuts
	lines = append(lines, headerStyle.Render("  Lists & Browse"))
	lines = append(lines, borderStyle.Render("  ──────────────────────────────────────────────────────────"))
	lines = append(lines, keyStyle.Render("  e")+"         "+descStyle.Render("Edit selected item"))
	lines = append(lines, keyStyle.Render("  d")+"         "+descStyle.Render("Delete selected item"))
	lines = append(lines, keyStyle.Render("  /")+"         "+descStyle.Render("Search"))
	lines = append(lines, keyStyle.Render("  f")+"         "+descStyle.Render("Filter"))
	lines = append(lines, "")

	// Footer
	lines = append(lines, borderStyle.Render("  ════════════════════════════════════════════════════════"))
	lines = append(lines, descStyle.Render("  Press ")+keyStyle.Render("?")+" "+descStyle.Render("or")+" "+keyStyle.Render("Esc")+" "+descStyle.Render("to close this help"))
	lines = append(lines, "")

	helpContent := lipgloss.JoinVertical(lipgloss.Left, lines...)

	// Center the help screen
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, helpContent)
}

// handleMenuInput processes menu navigation and selection
func (m *Model) handleMenuInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		// Move up in the TableListContainer
		if m.selectedMenuIndex > 0 {
			m.selectedMenuIndex--
			return m, nil
		}
	case "down", "j":
		// Move down in the TableListContainer
		if m.selectedMenuIndex < len(m.menuItems)-1 {
			m.selectedMenuIndex++
			return m, nil
		}
	case "enter", " ":
		// Select the current menu item
		selectedItem := m.menuItems[m.selectedMenuIndex]
		m.state = StateIntent
		cmd, err := m.intentRouter.ActivateIntent(selectedItem.Intent, make(map[string]interface{}))
		if err != nil {
			m.logger.Error("Failed to activate intent %s: %v", selectedItem.Intent, err)
			return m, nil
		}
		if cmd == nil {
			cmd = func() tea.Msg { return nil }
		}
		return m, cmd
	}
	return m, nil
}

// handleIntentInput forwards messages to the active intent
func (m *Model) handleIntentInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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

// viewMenu renders the main menu using lipgloss.JoinVertical for consistent centering
func (m *Model) viewMenu() string {
	// Ensure terminalInfo has current dimensions
	if !m.terminalInfo.IsValid && m.width > 0 && m.height > 0 {
		m.terminalInfo.Width = m.width
		m.terminalInfo.Height = m.height
		m.terminalInfo.IsValid = true
	}

	// Build menu components WITHOUT individual centering
	var parts []string

	// 1. Logo (static view, no animation during menu)
	logoView := m.logo.ViewStatic()
	parts = append(parts, logoView)

	// 2. Spacing between logo and menu
	parts = append(parts, "")
	parts = append(parts, "")

	// 3. Menu table with responsive columns
	tableView := m.renderResponsiveTable()
	parts = append(parts, tableView)

	// 4. Spacing between menu and help
	parts = append(parts, "")
	parts = append(parts, "")

	// 5. Help text
	helpText := "↑/k Up  ↓/j Down  Enter Select  ? Help  q Quit"
	parts = append(parts, helpText)

	// Join all parts with center alignment - aligns to widest line
	combined := lipgloss.JoinVertical(lipgloss.Center, parts...)

	// Center within terminal (both horizontal and vertical)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, combined)
}

// renderResponsiveTable creates the menu as simple text lines that can be centered
func (m *Model) renderResponsiveTable() string {
	var lines []string

	// Calculate responsive column widths based on terminal size
	category := m.terminalInfo.GetCategory()

	var actionWidth, descWidth int

	switch category {
	case terminal.SizeTiny:
		// Very small terminals: minimal widths
		actionWidth = 18
		descWidth = 28
	case terminal.SizeCompact:
		// Compact terminals: balanced widths
		actionWidth = 20
		descWidth = 35
	case terminal.SizeNormal:
		// Normal terminals: comfortable widths
		actionWidth = 22
		descWidth = 40
	case terminal.SizeLarge:
		// Large terminals: generous widths
		actionWidth = 25
		descWidth = 50
	default: // SizeXLarge
		// Extra large terminals: maximum widths
		actionWidth = 28
		descWidth = 60
	}

	// Create header row
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("12")).
		Bold(true)

	header := headerStyle.Render(
		lipgloss.NewStyle().Width(actionWidth).Render("Action") + "  " +
			lipgloss.NewStyle().Width(descWidth).Render("Description"),
	)
	lines = append(lines, header)

	// Add separator
	separator := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render(lipgloss.NewStyle().Width(actionWidth + descWidth + 2).Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	lines = append(lines, separator)

	// Create menu rows as simple text
	for i, item := range m.menuItems {
		// Style based on selection
		var rowStyle lipgloss.Style
		indicator := "  "

		if i == m.selectedMenuIndex {
			indicator = "▶ "
			rowStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("12")).
				Bold(true)
		} else {
			rowStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("252"))
		}

		// Format row with fixed widths
		actionText := lipgloss.NewStyle().Width(actionWidth).Render(indicator + item.Name)
		descText := lipgloss.NewStyle().Width(descWidth).Render(item.Help)

		row := rowStyle.Render(actionText + "  " + descText)
		lines = append(lines, row)
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

// GetMenuItems returns the menu items from the model
func (m *Model) GetMenuItems() []MenuItem {
	return m.menuItems
}

// createDefaultCVProfiles creates a set of default CV profiles for the GenerateCV intent
func createDefaultCVProfiles() []*intents.CVProfile {
	return []*intents.CVProfile{
		{
			ID:             "profile-staff-engineer",
			Name:           "Staff Engineer",
			TargetRole:     "staff",
			TargetAudience: "hiring_manager",
			Description:    "CV tailored for staff engineering roles",
		},
		{
			ID:             "profile-principal-engineer",
			Name:           "Principal Engineer",
			TargetRole:     "principal",
			TargetAudience: "hiring_manager",
			Description:    "CV tailored for principal/architect roles",
		},
		{
			ID:             "profile-engineering-manager",
			Name:           "Engineering Manager",
			TargetRole:     "em",
			TargetAudience: "hiring_manager",
			Description:    "CV tailored for engineering management roles",
		},
		{
			ID:             "profile-senior-engineer",
			Name:           "Senior Engineer",
			TargetRole:     "senior_ic",
			TargetAudience: "hiring_manager",
			Description:    "CV tailored for senior individual contributor roles",
		},
	}
}

// registerAllIntents registers all 10 intents with the router
func registerAllIntents(router *intents.DefaultIntentRouter, cliService *service.CLIEventService, careerService *careerservice.Service, log *logger.Logger, ctx context.Context, cvGenService cv.CVGenerationService, cvExportService *cv.ExportService) {
	// CaptureEvent
	_ = router.RegisterIntent("capture_event", func() intents.Intent {
		captureCtx := &intents.CaptureEventContext{
			CaptureStrategy: "manual",
			Metadata:        make(map[string]string),
			CLIEventService: cliService,
			CareerService:   careerService,
		}
		intent, err := intents.NewCaptureEventIntent(captureCtx)
		if err != nil {
			log.Error("Failed to create CaptureEvent intent: %v", err)
			return nil
		}
		return intent
	})

	// BrowseTimeline
	_ = router.RegisterIntent("browse_timeline", func() intents.Intent {
		events, err := careerService.GetEventRepository().List(ctx, careerrepo.ListFilters{
			Limit:     1000,
			SortBy:    "date",
			SortOrder: "desc",
		})
		if err != nil {
			log.Error("Failed to load events: %v", err)
			events = make([]*career.CareerEvent, 0)
		}
		browserCtx := &intents.BrowseTimelineContext{
			Events:          events,
			CLIEventService: cliService,
		}
		intent, err := intents.NewBrowseTimelineIntent(browserCtx)
		if err != nil {
			log.Error("Failed to create BrowseTimeline intent: %v", err)
			return nil
		}
		return intent
	})

	// ManageSkills
	_ = router.RegisterIntent("manage_skills", func() intents.Intent {
		skillsCtx := &intents.ManageSkillsContext{
			Ctx:             ctx,
			SkillRepository: careerService.GetSkillRepository(),
			Service:         careerService,
		}
		return intents.NewManageSkillsIntent(skillsCtx)
	})

	// GenerateCV
	_ = router.RegisterIntent("generate_cv", func() intents.Intent {
		events, err := careerService.GetEventRepository().List(ctx, careerrepo.ListFilters{Limit: 100})
		if err != nil || len(events) == 0 {
			events = []*career.CareerEvent{{ID: "ev-stub", Text: "Test event for navigation integration", Date: time.Now()}}
		}
		facts, err := careerService.GetFactRepository().List(ctx, careerrepo.FactListFilters{Limit: 100})
		if err != nil || len(facts) == 0 {
			facts = []*career.Fact{{ID: "fact-stub", Text: "Test fact for navigation integration", CompetencyCategories: []string{"technical"}, RoleFit: "staff", AudienceRelevance: []string{"peer"}, SourceEventID: "ev-stub"}}
		}
		cvCtx := &intents.GenerateCVContext{
			Events:                  events,
			Facts:                   facts,
			AvailableProfiles:       createDefaultCVProfiles(),
			DefaultProfile:          createDefaultCVProfiles()[0],
			CVGenerationService:     cvGenService,
			DataProcessingService:   cv.NewDataProcessingService(log),
			EnhancedBulletGenerator: cv.NewEnhancedBulletGenerator(log),
			ExportService:           cvExportService,
			AppContext:              ctx,
		}
		intent, err := intents.NewGenerateCVIntent(cvCtx)
		if err != nil {
			log.Error("Failed to create GenerateCV intent: %v", err)
			return nil
		}
		return intent
	})

	// ExportArtifact
	_ = router.RegisterIntent("export_artifact", func() intents.Intent {
		exportCtx := &intents.ExportArtifactContext{
			ArtifactTypes:    intents.DefaultArtifactTypes(),
			SupportedFormats: intents.DefaultSupportedFormats(),
			DefaultFormat:    intents.DefaultFormats(),
			Destinations:     intents.DefaultDestinations(),
			ExportService:    cvExportService,
			CareerService:    careerService,
			EventRepository:  careerService.GetEventRepository(),
			FactRepository:   careerService.GetFactRepository(),
			BurstRepository:  careerService.GetBurstRepository(),
			AppContext:       ctx,
		}
		intent, err := intents.NewExportArtifactIntent(exportCtx)
		if err != nil {
			log.Error("Failed to create ExportArtifact intent: %v", err)
			return nil
		}
		return intent
	})

	// ConfigureSystem
	_ = router.RegisterIntent("configure_system", func() intents.Intent {
		intent, err := intents.NewConfigureSystemIntent(ctx)
		if err != nil {
			log.Error("Failed to create ConfigureSystem intent: %v", err)
			return nil
		}
		return intent
	})

	// BurstManagement - Use helper constructor
	_ = router.RegisterIntent("burst_management", func() intents.Intent {
		burstRepo := careerService.GetBurstRepository()
		burstCtx := intents.NewBurstManagementContext(careerService, burstRepo, ctx)
		if burstCtx == nil {
			log.Error("Failed to create BurstManagement context")
			return nil
		}
		intent, err := intents.NewBurstManagementIntent(burstCtx)
		if err != nil || intent == nil {
			log.Error("Failed to create BurstManagement intent: %v", err)
			return nil
		}
		return intent
	})

	// FactManagement - Use helper constructor
	_ = router.RegisterIntent("fact_management", func() intents.Intent {
		factRepo := careerService.GetFactRepository()
		factCtx := intents.NewFactManagementContext(factRepo, ctx)
		if factCtx == nil {
			log.Error("Failed to create FactManagement context")
			return nil
		}
		intent := intents.NewFactManagementIntent(factCtx)
		if intent == nil {
			log.Error("Failed to create FactManagement intent")
			return nil
		}
		return intent
	})

	// ImportWizard - Use helper constructor
	_ = router.RegisterIntent("import_wizard", func() intents.Intent {
		importCtx := intents.NewImportWizardContext(ctx)
		if importCtx == nil {
			log.Error("Failed to create ImportWizard context")
			return nil
		}
		intent := intents.NewImportWizardIntent(importCtx)
		if intent == nil {
			log.Error("Failed to create ImportWizard intent")
			return nil
		}
		return intent
	})

	// MetadataEditor - Use helper constructor
	_ = router.RegisterIntent("metadata_editor", func() intents.Intent {
		metaCtx := intents.NewMetadataEditorContext(ctx)
		if metaCtx == nil {
			log.Error("Failed to create MetadataEditor context")
			return nil
		}
		intent := intents.NewMetadataEditorIntent(metaCtx)
		if intent == nil {
			log.Error("Failed to create MetadataEditor intent")
			return nil
		}
		return intent
	})

	// BulkOperations - Use helper constructor
	_ = router.RegisterIntent("bulk_operations", func() intents.Intent {
		bulkCtx := intents.NewBulkOperationsContext(ctx)
		if bulkCtx == nil {
			log.Error("Failed to create BulkOperations context")
			return nil
		}
		intent := intents.NewBulkOperationsIntent(bulkCtx)
		if intent == nil {
			log.Error("Failed to create BulkOperations intent")
			return nil
		}
		return intent
	})
}

// initConfigManager initializes the CV config manager
func initConfigManager(log *logger.Logger) cv.ConfigManager {
	var configMgr cv.ConfigManager
	yamlMgr, err := cv.NewYAMLConfigManager(log)
	if err != nil {
		log.Error("Failed to initialize CV config manager: %v", err)
		configMgr = cv.NewMemoryConfigManager()
	} else {
		configMgr = yamlMgr
	}
	return configMgr
}

// initCVGenerationService initializes the CV generation service
func initCVGenerationService(careerService *careerservice.Service, configMgr cv.ConfigManager, log *logger.Logger) cv.CVGenerationService {
	bulletGenerator := cv.NewBulletGenerator(
		careerService.GetEventRepository(),
		careerService.GetFactRepository(),
		log,
	)
	sectionBuilder := cv.NewSectionBuilder(log)
	return cv.NewCVGenerationService(
		careerService.GetEventRepository(),
		careerService.GetFactRepository(),
		configMgr,
		bulletGenerator,
		sectionBuilder,
		log,
	)
}

// IntentCompletedMsg signals the completion of an intent
type IntentCompletedMsg struct{}

// MainMenuSelectMsg signals a menu selection by index
type MainMenuSelectMsg struct {
	Index int
}

// CompleteIntentMsg signals that the user has completed an intent
type CompleteIntentMsg struct{}

// SetInitialScreen sets the initial screen to display
func (m *Model) SetInitialScreen(screen Screen) {
	// For now, this is a no-op since we always start with menu
	// In the future, this could be used to navigate directly to a specific intent
}

// SetInitialCaptureMode sets the initial capture mode for CaptureEvent intent
func (m *Model) SetInitialCaptureMode(mode string) {
	// This would be used to configure the CaptureEvent intent when activated
	// For now, it's a placeholder
}

// GetState returns the current application state
func (m *Model) GetState() AppState {
	return m.state
}

// GetActiveIntent returns the currently active intent
func (m *Model) GetActiveIntent() intents.Intent {
	return m.intentRouter.GetActiveIntent()
}
