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
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
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
	state  AppState
	width  int
	height int

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
	logo.SetExternalCentering(true) // Let container handle centering

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
			// Help screen could be implemented here
			return m, nil
		case "home", "esc", "escape":
			if m.state == StateIntent {
				m.state = StateMenu
				m.selectedMenuIndex = 0
				return m, nil
			}
		}

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

// viewMenu renders the main menu using SmartContainer for proper centering
func (m *Model) viewMenu() string {
	// Ensure terminalInfo has current dimensions
	if !m.terminalInfo.IsValid && m.width > 0 && m.height > 0 {
		m.terminalInfo.Width = m.width
		m.terminalInfo.Height = m.height
		m.terminalInfo.IsValid = true
	}

	// Use SmartContainer with terminal info for intelligent centering
	container := components.NewSmartContainer(m.terminalInfo)
	container.SetCenteringMode(components.CenterBoth)

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

	// Let SmartContainer handle ALL centering
	content := ""
	for i, part := range parts {
		if i > 0 {
			content += "\n"
		}
		content += part
	}
	return container.SetContent(content).Render()
}

// renderResponsiveTable creates the menu table with responsive column widths
func (m *Model) renderResponsiveTable() string {
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

	// Create table rows
	rows := make([]table.Row, len(m.menuItems))
	for i, item := range m.menuItems {
		// Add selection indicator
		indicator := "  "
		if i == m.selectedMenuIndex {
			indicator = "▶ "
		}
		rows[i] = table.Row{indicator + item.Name, item.Help}
	}

	// Create table with responsive columns
	tableModel := table.New(
		table.WithColumns([]table.Column{
			{Title: "Action", Width: actionWidth},
			{Title: "Description", Width: descWidth},
		}),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(m.menuItems)),
	)

	// Set cursor position
	if m.selectedMenuIndex >= 0 && m.selectedMenuIndex < len(rows) {
		tableModel.SetCursor(m.selectedMenuIndex)
	}

	return tableModel.View()
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
			Events: events,
		}
		intent, err := intents.NewBrowseTimelineIntent(browserCtx)
		if err != nil {
			log.Error("Failed to create BrowseTimeline intent: %v", err)
			return nil
		}
		return intent
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
		intent, err := intents.NewExportArtifactIntent(ctx)
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
