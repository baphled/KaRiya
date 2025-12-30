package app

import (
	"context"
	"strings"

	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/cli/styles"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Screen represents the different screens in the application
type Screen string

const (
	HomeScreen       Screen = "home"
	CaptureScreen    Screen = "capture"
	ListScreen       Screen = "list"
	ViewScreen       Screen = "view"
	QuitScreen       Screen = "quit"
	SuccessScreen    Screen = "success"
	ActionMenuScreen Screen = "action_menu"
)

// Model represents the main application state
type Model struct {
	cliService      *service.CLIEventService
	service         *careerservice.Service
	currentScreen   Screen
	previousScreen  Screen
	width           int
	height          int
	formModel       *models.FormModel
	successModel    *models.SuccessModel
	listModel       *models.ListModel
	detailsModel    *models.DetailsModel
	actionMenuModel *models.ActionMenuModel
}

// NewModel creates a new application model
func NewModel(cliService *service.CLIEventService, careerService *careerservice.Service) *Model {
	ctx := context.Background()
	return &Model{
		cliService:      cliService,
		service:         careerService,
		currentScreen:   HomeScreen,
		previousScreen:  HomeScreen,
		width:           80,
		height:          24,
		formModel:       models.NewFormModel(cliService),
		successModel:    nil,
		listModel:       models.NewListModel(careerService, ctx),
		detailsModel:    nil,
		actionMenuModel: nil,
	}
}

// Init initializes the application
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model state
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle quit and back messages first
	switch msg.(type) {
	case models.BackMsg:
		m.currentScreen = m.previousScreen
		return m, nil
	case models.QuitMsg:
		return m, tea.Quit
	}

	// Handle ViewEventMsg
	if viewMsg, ok := msg.(ViewEventMsg); ok {
		m.detailsModel = models.NewDetailsModel(viewMsg.Event)
		m.previousScreen = ListScreen
		m.currentScreen = ViewScreen
		return m, nil
	}

	// Handle EventActionMenuMsg
	if actionMenuMsg, ok := msg.(EventActionMenuMsg); ok {
		m.actionMenuModel = models.NewActionMenuModel(actionMenuMsg.Event)
		m.previousScreen = m.currentScreen
		m.currentScreen = ActionMenuScreen
		return m, nil
	}

	// Handle EventActionSelectedMsg
	if actionMsg, ok := msg.(models.EventActionSelectedMsg); ok {
		switch actionMsg.Action {
		case models.EventActionView:
			m.detailsModel = models.NewDetailsModel(actionMsg.Event)
			m.currentScreen = ViewScreen
		case models.EventActionEdit:
			m.previousScreen = m.currentScreen
			m.currentScreen = CaptureScreen
			m.formModel = models.NewFormModel(m.cliService)
		case models.EventActionDelete:
			m.previousScreen = m.currentScreen
			m.currentScreen = ListScreen
		}
		return m, nil
	}

	// Handle FormSubmittedMsg
	if submitMsg, ok := msg.(FormSubmittedMsg); ok {
		m.successModel = models.NewSuccessModel(submitMsg.Event)
		m.previousScreen = m.currentScreen
		m.currentScreen = SuccessScreen
		return m, nil
	}

	// Handle SuccessModel messages
	switch msg.(type) {
	case models.CaptureAnotherMsg:
		m.previousScreen = m.currentScreen
		m.currentScreen = CaptureScreen
		m.formModel = models.NewFormModel(m.cliService)
		m.successModel = nil
		return m, nil
	case models.ViewRecentMsg:
		m.previousScreen = m.currentScreen
		m.currentScreen = ListScreen
		m.successModel = nil
		return m, nil
	}

	// Delegate to active model based on current screen
	switch m.currentScreen {
	case ViewScreen:
		if m.detailsModel != nil {
			updatedDetailsModel, cmd := m.detailsModel.Update(msg)
			m.detailsModel = updatedDetailsModel.(*models.DetailsModel)
			return m, cmd
		}

	case CaptureScreen:
		if m.formModel != nil {
			updatedFormModel, cmd := m.formModel.Update(msg)
			m.formModel = updatedFormModel.(*models.FormModel)
			if m.formModel.Submitted() {
				m.successModel = models.NewSuccessModel(m.formModel.Event())
				m.currentScreen = SuccessScreen
			}
			return m, cmd
		}

	case SuccessScreen:
		if m.successModel != nil {
			updatedSuccessModel, cmd := m.successModel.Update(msg)
			m.successModel = updatedSuccessModel.(*models.SuccessModel)
			return m, cmd
		}

	case ListScreen:
		if m.listModel != nil {
			updatedListModel, cmd := m.listModel.Update(msg)
			m.listModel = updatedListModel.(*models.ListModel)
			return m, cmd
		}

	case ActionMenuScreen:
		if m.actionMenuModel != nil {
			updatedActionMenuModel, cmd := m.actionMenuModel.Update(msg)
			m.actionMenuModel = updatedActionMenuModel.(*models.ActionMenuModel)
			return m, cmd
		}
	}

	// Handle global navigation shortcuts
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "h":
			m.previousScreen = m.currentScreen
			m.currentScreen = HomeScreen
		case "c":
			m.previousScreen = m.currentScreen
			m.currentScreen = CaptureScreen
			m.formModel = models.NewFormModel(m.cliService)
		case "l":
			m.previousScreen = m.currentScreen
			m.currentScreen = ListScreen
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
		if m.formModel != nil {
			return m.formModel.View()
		}
		return "Error: Form model not initialized\n"
	case ListScreen:
		if m.listModel != nil {
			return m.listModel.View()
		}
		return "Error: List model not initialized\n"
	case ViewScreen:
		if m.detailsModel != nil {
			return m.detailsModel.View()
		}
		return "Error: Details model not initialized\n"
	case SuccessScreen:
		if m.successModel != nil {
			return m.successModel.View()
		}
		return "Success!\n"
	case ActionMenuScreen:
		if m.actionMenuModel != nil {
			return m.actionMenuModel.View()
		}
		return "Error: Action menu not initialized\n"
	case QuitScreen:
		return "Goodbye!\n"
	default:
		return m.renderHome()
	}
}

// renderHome renders the home screen
func (m *Model) renderHome() string {
	title := styles.HeaderMain.Render("KaRiya - Career Journal CLI")
	commandsHeader := styles.HeaderSection.Render("Commands:")
	commands := []string{
		styles.InfoText.Render("c") + " - Capture Career Event",
		styles.InfoText.Render("l") + " - List Events",
		styles.InfoText.Render("h") + " - Home",
		styles.InfoText.Render("q") + " - Quit",
	}
	commandsList := strings.Join(commands, "\n")
	cta := styles.SuccessBox.Width(styles.MaxWidth(m.width) - 4).Render("Press 'c' to get started!")
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		commandsHeader,
		commandsList,
		"",
		cta,
	)
	card := styles.ResponsiveCard(m.width).Render(content)
	return styles.Center(card, m.width, m.height)
}

// SetInitialScreen sets the initial screen to display on startup
func (m *Model) SetInitialScreen(screen Screen) {
	m.currentScreen = screen
	m.previousScreen = screen
}

// SetInitialCaptureMode sets the initial capture mode for the form
func (m *Model) SetInitialCaptureMode(mode string) {
	if m.formModel != nil {
		m.formModel.SetInitialMode(mode)
	}
}
