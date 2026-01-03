package app

import (
	"context"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/service"
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
	state  AppState
	width  int
	height int

	// Menu state
	selectedMenuIndex int
	menuItems         []MenuItem

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
	registerAllIntents(router, careerService, log, ctx)

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
		ctx:               ctx,
	}
}

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case "home", "escape":
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
		m.width = msg.Width
		m.height = msg.Height

	case intentResultMsg:
		// Handle result from completed intent
		m.state = StateMenu
		m.selectedMenuIndex = 0
		return m, nil
	}

	return m, nil
}

// View renders the current screen
func (m *Model) View() string {
	if m.state == StateMenu {
		return m.viewMenu()
	} else if m.state == StateIntent {
		return m.intentRouter.View()
	}
	return ""
}

// handleMenuInput processes menu navigation and selection
func (m *Model) handleMenuInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.selectedMenuIndex > 0 {
			m.selectedMenuIndex--
		}
	case "down", "j":
		if m.selectedMenuIndex < len(m.menuItems)-1 {
			m.selectedMenuIndex++
		}
	case "enter", " ":
		selectedItem := m.menuItems[m.selectedMenuIndex]
		m.state = StateIntent
		return m, m.activateIntent(selectedItem.Intent)
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
		return m, cmd
	}

	return m, cmd
}

// activateIntent activates an intent by name
func (m *Model) activateIntent(intentName string) tea.Cmd {
	return func() tea.Msg {
		_, _ = m.intentRouter.ActivateIntent(intentName, make(map[string]interface{}))
		return nil
	}
}

// viewMenu renders the main menu
func (m *Model) viewMenu() string {
	var output string
	output += lipgloss.NewStyle().Bold(true).Render("KaRiya - Career Event Manager\n\n")
	output += "Select an option:\n\n"

	for i, item := range m.menuItems {
		prefix := "  "
		if i == m.selectedMenuIndex {
			prefix = "> "
			output += lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Render(prefix + item.Name + "\n")
		} else {
			output += prefix + item.Name + "\n"
		}
	}

	output += "\n" + lipgloss.NewStyle().Faint(true).Render("↑/↓: Navigate | Enter: Select | Ctrl+C: Quit\n")

	return output
}

// registerAllIntents registers all 10 intents with the router
func registerAllIntents(router *intents.DefaultIntentRouter, careerService *careerservice.Service, log *logger.Logger, ctx context.Context) {
	// CaptureEvent
	_ = router.RegisterIntent("capture_event", func() intents.Intent {
		captureCtx := &intents.CaptureEventContext{
			CaptureStrategy: "manual",
			Metadata:        make(map[string]string),
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
		events, err := careerService.GetEventRepository().List(ctx, careerrepo.ListFilters{Limit: 1000})
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
		events, err := careerService.GetEventRepository().List(ctx, careerrepo.ListFilters{Limit: 1000})
		if err != nil {
			log.Error("Failed to load events: %v", err)
			events = make([]*career.CareerEvent, 0)
		}
		facts, err := careerService.GetFactRepository().List(ctx, careerrepo.FactListFilters{Limit: 1000})
		if err != nil {
			log.Error("Failed to load facts: %v", err)
			facts = make([]*career.Fact, 0)
		}
		cvCtx := &intents.GenerateCVContext{
			Events: events,
			Facts:  facts,
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
		intent := intents.NewBurstManagementIntent(burstCtx)
		if intent == nil {
			log.Error("Failed to create BurstManagement intent")
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

// intentResultMsg is used to signal intent completion
type intentResultMsg struct{}

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
