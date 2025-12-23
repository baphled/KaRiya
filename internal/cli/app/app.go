package app

import (
	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/cli/service"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
)

// Screen represents the different screens in the application
type Screen string

const (
	// HomeScreen is the main menu screen
	HomeScreen Screen = "home"
	// CaptureScreen is the event capture screen
	CaptureScreen Screen = "capture"
	// ListScreen is the event list screen
	ListScreen Screen = "list"
	// ViewScreen is the event detail screen
	ViewScreen Screen = "view"
	// QuitScreen is the quit screen
	QuitScreen Screen = "quit"
	// SuccessScreen is the success screen after event capture
	SuccessScreen Screen = "success"
)

// Model represents the main application state
type Model struct {
	cliService     *service.CLIEventService
	service        *careerservice.Service
	currentScreen  Screen
	previousScreen Screen
	width          int
	height         int
	err            error
	formModel      *models.FormModel
	successModel   *models.SuccessModel
}

// NewModel creates a new application model
func NewModel(cliService *service.CLIEventService, careerService *careerservice.Service) *Model {
	return &Model{
		cliService:     cliService,
		service:        careerService,
		currentScreen:  HomeScreen,
		previousScreen: HomeScreen,
		width:          80,
		height:         24,
		err:            nil,
		formModel:      models.NewFormModel(cliService),
		successModel:   nil, // Will be created after form submission
	}
}

// Init initializes the application
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model state
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "h":
			m.previousScreen = m.currentScreen
			m.currentScreen = HomeScreen
		case "c":
			m.previousScreen = m.currentScreen
			m.currentScreen = CaptureScreen
			m.formModel = models.NewFormModel(m.cliService) // Reset form when entering capture screen
		case "l":
			m.previousScreen = m.currentScreen
			m.currentScreen = ListScreen
		case "backspace":
			m.currentScreen = m.previousScreen
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	// Delegate to FormModel when on CaptureScreen
	if m.currentScreen == CaptureScreen && m.formModel != nil {
		updatedFormModel, cmd := m.formModel.Update(msg)
		m.formModel = updatedFormModel.(*models.FormModel)

		// Check if form was submitted
		if m.formModel.Submitted() {
			m.successModel = models.NewSuccessModel(m.formModel.Event())
			m.currentScreen = SuccessScreen
		}

		return m, cmd
	}

	// Delegate to SuccessModel when on SuccessScreen
	if m.currentScreen == SuccessScreen && m.successModel != nil {
		updatedSuccessModel, cmd := m.successModel.Update(msg)
		m.successModel = updatedSuccessModel.(*models.SuccessModel)
		return m, cmd
	}

	return m, nil
}

// View renders the current screen
func (m *Model) View() string {
	switch m.currentScreen {
	case HomeScreen:
		return m.renderHome()
	case CaptureScreen:
		if m.formModel != nil {
			return m.formModel.View()
		}
		return m.renderCapture()
	case ListScreen:
		return m.renderList()
	case ViewScreen:
		return m.renderView()
	case SuccessScreen:
		if m.successModel != nil {
			return m.successModel.View()
		}
		return "Success!\n"
	case QuitScreen:
		return "Goodbye!\n"
	default:
		return m.renderHome()
	}
}

// renderHome renders the home screen
func (m *Model) renderHome() string {
	return `
╔═══════════════════════════════════════════╗
║     KaRiya - Career Journal CLI           ║
║                                           ║
║  Commands:                                ║
║    c - Capture Career Event               ║
║    l - List Events                        ║
║    h - Home                               ║
║    q - Quit                               ║
║                                           ║
║  Press 'c' to get started!                ║
╚═══════════════════════════════════════════╝
`
}

// renderCapture renders the capture screen (placeholder)
func (m *Model) renderCapture() string {
	return `
╔════════════════════════════════════════╗
║       Capture Career Event             ║
║                                        ║
║ Event Text (required):                 ║
║ [_____________________________________] ║
║                                        ║
║ Date (optional):                       ║
║ [_____________________________________] ║
║                                        ║
║ Company (optional):                    ║
║ [_____________________________________] ║
║                                        ║
║ Project (optional):                    ║
║ [_____________________________________] ║
║                                        ║
║ [Submit]  [Cancel]                     ║
║                                        ║
║ Press 'backspace' to go back           ║
╚════════════════════════════════════════╝
`
}

// renderList renders the list screen
func (m *Model) renderList() string {
	return `
╔════════════════════════════════════════╗
║     Recent Career Events               ║
║                                        ║
║ Event 1                                ║
║ [Some career event description here]   ║
║                                        ║
║ Event 2                                ║
║ [Another career event description]     ║
║                                        ║
║ Press 'backspace' to go back           ║
╚════════════════════════════════════════╝
`
}

// renderView renders the view screen
func (m *Model) renderView() string {
	return `
╔════════════════════════════════════════╗
║        Event Details                   ║
║                                        ║
║ Title: [Event Title]                   ║
║ Date: [Event Date]                     ║
║ Description: [Event Description]       ║
║                                        ║
║ Press 'backspace' to go back           ║
╚════════════════════════════════════════╝
`
}
