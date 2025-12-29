package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/cli/styles"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	// Handle messages from SuccessModel first
	switch msg.(type) {
	case models.CaptureAnotherMsg:
		// Reset form and return to capture screen
		m.previousScreen = m.currentScreen
		m.currentScreen = CaptureScreen
		m.formModel = models.NewFormModel(m.cliService)
		m.successModel = nil
		return m, nil

	case models.ViewRecentMsg:
		// Navigate to list screen
		m.previousScreen = m.currentScreen
		m.currentScreen = ListScreen
		m.successModel = nil
		return m, nil
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

	// Handle global navigation shortcuts (only on non-input screens)
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
	// Title
	title := styles.HeaderMain.
		Render("KaRiya - Career Journal CLI")

	// Commands section
	commandsHeader := styles.HeaderSection.
		Render("Commands:")

	commands := []string{
		styles.InfoText.Render("c") + " - Capture Career Event",
		styles.InfoText.Render("l") + " - List Events",
		styles.InfoText.Render("h") + " - Home",
		styles.InfoText.Render("q") + " - Quit",
	}
	commandsList := strings.Join(commands, "\n")

	// Call to action
	cta := styles.SuccessBox.
		Width(styles.MaxWidth(m.width) - 4).
		Render("Press 'c' to get started!")

	// Combine all sections
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		commandsHeader,
		commandsList,
		"",
		cta,
	)

	// Wrap in a card
	card := styles.ResponsiveCard(m.width).
		Render(content)

	// Center the card on screen
	return styles.Center(card, m.width, m.height)
}

// renderCapture renders the capture screen (placeholder)
func (m *Model) renderCapture() string {
	// Title
	title := styles.HeaderMain.
		Render("Capture Career Event")

	// Form fields (placeholder representation)
	fields := []string{
		styles.InputLabel.Render("Event Text (required):"),
		styles.InputBase.Width(60).Render("_______________________________________________________"),
		"",
		styles.InputLabel.Render("Date (optional):"),
		styles.InputBase.Width(60).Render("_______________________________________________________"),
		"",
		styles.InputLabel.Render("Company (optional):"),
		styles.InputBase.Width(60).Render("_______________________________________________________"),
		"",
		styles.InputLabel.Render("Project (optional):"),
		styles.InputBase.Width(60).Render("_______________________________________________________"),
	}

	// Buttons
	submitBtn := styles.ButtonPrimary.Render("Submit")
	cancelBtn := styles.ButtonSecondary.Render("Cancel")
	buttons := lipgloss.JoinHorizontal(lipgloss.Left, submitBtn, cancelBtn)

	// Navigation hint
	navHint := styles.InfoHint.Render("Press 'backspace' to go back")

	// Combine all sections
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		strings.Join(fields, "\n"),
		"",
		buttons,
		"",
		navHint,
	)

	// Wrap in a card
	card := styles.ResponsiveCard(m.width).
		Render(content)

	// Center the card on screen
	return styles.Center(card, m.width, m.height)
}

// renderList renders the list screen
func (m *Model) renderList() string {
	// Title
	title := styles.HeaderMain.
		Render("Recent Career Events")

	// Fetch events from service
	ctx := context.Background()
	filters := &careerrepo.ListFilters{
		SortBy:    "date",
		SortOrder: "desc",
		Limit:     10,
		Offset:    0,
	}

	events, err := m.cliService.ListEvents(ctx, filters)

	// Build event cards
	var eventCards []string
	if err != nil {
		errorMsg := styles.ErrorText.Render(fmt.Sprintf("Error loading events: %v", err))
		eventCards = append(eventCards, errorMsg)
	} else if len(events) == 0 {
		noEventsMsg := styles.InfoText.Render("No events found. Press 'c' to capture your first event!")
		eventCards = append(eventCards, noEventsMsg)
	} else {
		for i, event := range events {
			// Event header with date
			eventHeader := styles.CardHeader.Render(fmt.Sprintf("Event %d - %s", i+1, event.Date.Format("2006-01-02")))

			// Event text (truncate if too long)
			eventText := event.Text
			if len(eventText) > 100 {
				eventText = eventText[:97] + "..."
			}
			eventContent := styles.CardContent.Render(eventText)

			// Company info if available
			var companyInfo string
			if event.Company != "" {
				companyInfo = styles.CardFooter.Render(fmt.Sprintf("Company: %s", event.Company))
			}

			// Combine event parts
			var eventParts []string
			eventParts = append(eventParts, eventHeader, eventContent)
			if companyInfo != "" {
				eventParts = append(eventParts, companyInfo)
			}

			// Create event card
			eventCard := styles.CardBase.Copy().
				Width(styles.MaxWidth(m.width) - 4).
				Render(lipgloss.JoinVertical(lipgloss.Left, eventParts...))

			eventCards = append(eventCards, eventCard)
		}
	}

	// Navigation hint
	navHint := styles.InfoHint.Render("Press 'backspace' to go back | 'c' to capture event")

	// Combine all sections
	sections := []string{title, ""}
	sections = append(sections, eventCards...)
	sections = append(sections, "", navHint)

	content := lipgloss.JoinVertical(lipgloss.Left, sections...)

	// Wrap in a card
	card := styles.ResponsiveCard(m.width).
		Render(content)

	// Center the card on screen
	return styles.Center(card, m.width, m.height)
}

// renderView renders the view screen
func (m *Model) renderView() string {
	// Title
	title := styles.HeaderMain.
		Render("Event Details")

	// Placeholder content
	detailsHeader := styles.CardHeader.Render("Event Information")

	details := []string{
		styles.CardContent.Render("Title: [Event Title]"),
		styles.CardContent.Render("Date: [Event Date]"),
		styles.CardContent.Render("Description: [Event Description]"),
	}
	detailsContent := strings.Join(details, "\n")

	// Event details card
	eventCard := styles.CardBase.
		Width(styles.MaxWidth(m.width) - 4).
		Render(lipgloss.JoinVertical(lipgloss.Left, detailsHeader, "", detailsContent))

	// Navigation hint
	navHint := styles.InfoHint.Render("Press 'backspace' to go back")

	// Combine all sections
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		eventCard,
		"",
		navHint,
	)

	// Wrap in a card
	card := styles.ResponsiveCard(m.width).
		Render(content)

	// Center the card on screen
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
