package app

import (
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

	return m, nil
}

// View renders the current screen
func (m *Model) View() string {
	switch m.currentScreen {
	case HomeScreen:
		return m.renderHome()
	case CaptureScreen:
		return m.renderCapture()
	case ListScreen:
		return m.renderList()
	case ViewScreen:
		return m.renderView()
	case QuitScreen:
		return "Goodbye!\n"
	default:
		return m.renderHome()
	}
}

// renderHome renders the home screen
func (m *Model) renderHome() string {
	return `
╔════════════════════════════════════════╗
║       KaRiya - Career Journal CLI      ║
║                                        ║
║ Welcome! Choose an option:             ║
║                                        ║
║ [c] Capture Career Event               ║
║ [l] List Events                        ║
║ [q] Quit                               ║
║                                        ║
║ Use arrow keys to navigate             ║
║ Press 'h' to return home               ║
╚════════════════════════════════════════╝
`
}

// renderCapture renders the capture screen
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
║       Recent Career Events             ║
║                                        ║
║ Event 1: Completed important project   ║
║ Date: 2025-12-23                       ║
║ Company: TechCorp Inc.                 ║
║ Tags: [technical] [achievement]        ║
║                                        ║
║ Event 2: Led cross-functional team     ║
║ Date: 2025-12-22                       ║
║ Company: StartupXYZ                    ║
║ Tags: [leadership]                     ║
║                                        ║
║ [↑] [↓] Navigate  [Enter] View Details ║
║ Press 'backspace' to go back           ║
╚════════════════════════════════════════╝
`
}

// renderView renders the event detail screen
func (m *Model) renderView() string {
	return `
╔════════════════════════════════════════╗
║       Event Details                    ║
║                                        ║
║ Title:                                 ║
║ Completed important project            ║
║                                        ║
║ Date: 2025-12-23                       ║
║ Company: TechCorp Inc.                 ║
║ Project: Platform Migration            ║
║                                        ║
║ Tags:                                  ║
║ [technical] [achievement] [leadership] ║
║                                        ║
║ Category: Technical                    ║
║                                        ║
║ [Edit]  [Delete]  [Back]               ║
║                                        ║
║ Press 'backspace' to go back           ║
╚════════════════════════════════════════╝
`
}
