package models

import (
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ActionOption represents an available action after successful event capture
type ActionOption int

const (
	CaptureAnotherOption ActionOption = iota
	ViewRecentOption
	ExitOption
)

// SuccessModel represents the success screen state after event capture
type SuccessModel struct {
	event          *career.CareerEvent
	selectedAction ActionOption
	width          int
	height         int
}

// NewSuccessModel creates a new success model with the captured event
func NewSuccessModel(event *career.CareerEvent) *SuccessModel {
	return &SuccessModel{
		event:          event,
		selectedAction: CaptureAnotherOption,
		width:          80,
		height:         24,
	}
}

// Event returns the captured event
func (m *SuccessModel) Event() *career.CareerEvent {
	return m.event
}

// Init initializes the success model
func (m *SuccessModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m *SuccessModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyLeft:
			// Navigate left through actions
			if m.selectedAction > CaptureAnotherOption {
				m.selectedAction--
			}
			return m, nil

		case tea.KeyRight:
			// Navigate right through actions
			if m.selectedAction < ExitOption {
				m.selectedAction++
			}
			return m, nil

		case tea.KeyEnter:
			// Execute selected action
			return m, m.executeAction()

		case tea.KeyEsc, tea.KeyCtrlC:
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

// executeAction returns a command based on the selected action
func (m *SuccessModel) executeAction() tea.Cmd {
	switch m.selectedAction {
	case CaptureAnotherOption:
		return func() tea.Msg {
			return CaptureAnotherMsg{}
		}
	case ViewRecentOption:
		return func() tea.Msg {
			return ViewRecentMsg{}
		}
	case ExitOption:
		return tea.Quit
	}
	return nil
}

// View renders the success screen
func (m *SuccessModel) View() string {
	var b strings.Builder

	// Create text style for help text
	textMuted := lipgloss.NewStyle().Foreground(styles.ColorTextMuted)

	// Success header
	header := styles.HeaderMain.Render("✓ Success! Event Captured")
	b.WriteString(header)
	b.WriteString("\n\n")

	// Event card
	card := m.renderEventCard()
	b.WriteString(card)
	b.WriteString("\n\n")

	// Action buttons
	actions := m.renderActions()
	b.WriteString(actions)
	b.WriteString("\n")

	// Help text
	help := textMuted.Render("← → to navigate • Enter to select • Esc to cancel")
	b.WriteString(help)
	b.WriteString("\n")

	return b.String()
}

// renderEventCard renders the event details in a styled card
func (m *SuccessModel) renderEventCard() string {
	var b strings.Builder

	// Create text styles using available colors
	textSecondary := lipgloss.NewStyle().Foreground(styles.ColorTextSecondary)
	textPrimary := lipgloss.NewStyle().Foreground(styles.ColorTextPrimary)
	textMuted := lipgloss.NewStyle().Foreground(styles.ColorTextMuted)

	// Event text
	textLabel := textSecondary.Render("Event:")
	textValue := textPrimary.Render(m.event.Text)
	b.WriteString(fmt.Sprintf("%s\n%s\n\n", textLabel, textValue))

	// Date
	dateLabel := textSecondary.Render("Date:")
	dateValue := textPrimary.Render(m.event.Date.Format("2006-01-02"))
	b.WriteString(fmt.Sprintf("%s %s\n", dateLabel, dateValue))

	// Company (if provided)
	if m.event.Company != "" {
		companyLabel := textSecondary.Render("Company:")
		companyValue := textPrimary.Render(m.event.Company)
		b.WriteString(fmt.Sprintf("%s %s\n", companyLabel, companyValue))
	}

	// Project (if provided)
	if m.event.Project != "" {
		projectLabel := textSecondary.Render("Project:")
		projectValue := textPrimary.Render(m.event.Project)
		b.WriteString(fmt.Sprintf("%s %s\n", projectLabel, projectValue))
	}

	// Tags (if provided)
	if len(m.event.Tags) > 0 {
		tagsLabel := textSecondary.Render("Tags:")
		tagsValue := m.renderTags()
		b.WriteString(fmt.Sprintf("\n%s\n%s\n", tagsLabel, tagsValue))
	}

	// Event ID
	b.WriteString("\n")
	idLabel := textMuted.Render("Event ID:")
	idValue := textMuted.Render(m.event.ID)
	b.WriteString(fmt.Sprintf("%s %s\n", idLabel, idValue))

	// Wrap in card style
	cardContent := b.String()
	card := styles.CardBase.
		Width(m.width - 4).
		Render(cardContent)

	return card
}

// renderTags renders the event tags with styling
func (m *SuccessModel) renderTags() string {
	var tagStyles []string
	for _, tag := range m.event.Tags {
		tagStyle := lipgloss.NewStyle().
			Foreground(styles.ColorAccentPurple).
			Background(styles.ColorBackgroundAlt).
			Padding(0, 1).
			MarginRight(1).
			Render(tag)
		tagStyles = append(tagStyles, tagStyle)
	}
	return strings.Join(tagStyles, "")
}

// renderActions renders the action buttons
func (m *SuccessModel) renderActions() string {
	var buttons []string

	// Capture Another button
	captureBtn := m.renderButton("Capture Another", m.selectedAction == CaptureAnotherOption)
	buttons = append(buttons, captureBtn)

	// View Recent button
	viewBtn := m.renderButton("View Recent", m.selectedAction == ViewRecentOption)
	buttons = append(buttons, viewBtn)

	// Exit button
	exitBtn := m.renderButton("Exit", m.selectedAction == ExitOption)
	buttons = append(buttons, exitBtn)

	return lipgloss.JoinHorizontal(lipgloss.Top, buttons...)
}

// renderButton renders a single action button
func (m *SuccessModel) renderButton(label string, selected bool) string {
	style := styles.ButtonSecondary
	if selected {
		style = styles.ButtonFocused
	}
	return style.Render(label)
}

// Message types for communication with parent model

// CaptureAnotherMsg signals to capture another event
type CaptureAnotherMsg struct{}

// ViewRecentMsg signals to view recent events
type ViewRecentMsg struct{}

