package models

import (
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// DetailsModel represents the event details screen
type DetailsModel struct {
	event  *career.CareerEvent
	width  int
	height int
}

// NewDetailsModel creates a new details model
func NewDetailsModel(event *career.CareerEvent) *DetailsModel {
	return &DetailsModel{
		event:  event,
		width:  80,
		height: 24,
	}
}

// Init initializes the model
func (m *DetailsModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m *DetailsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// Signal back navigation to parent
			return m, func() tea.Msg { return BackMsg{} }
		case "ctrl+c", "q":
			// Signal quit to parent
			return m, func() tea.Msg { return QuitMsg{} }
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

// View renders the details view
func (m *DetailsModel) View() string {
	if m.event == nil {
		return styles.ErrorBox.Render("No event selected\n\nPress 'esc' to return to list")
	}

	var sb strings.Builder

	// Title
	sb.WriteString(styles.HeaderSection.Render("Event Details"))
	sb.WriteString("\n\n")

	// Event Text
	sb.WriteString(styles.InputLabel.Render("Description:"))
	sb.WriteString("\n")
	sb.WriteString(lipgloss.NewStyle().Foreground(styles.ColorTextPrimary).Render(m.event.Text))
	sb.WriteString("\n\n")

	// Date
	sb.WriteString(styles.InputLabel.Render("Date:"))
	sb.WriteString("\n")
	dateStr := m.event.Date.Format("2006-01-02")
	sb.WriteString(lipgloss.NewStyle().Foreground(styles.ColorTextPrimary).Render(dateStr))
	sb.WriteString("\n\n")

	// Company (if available)
	if m.event.Company != "" {
		sb.WriteString(styles.InputLabel.Render("Company:"))
		sb.WriteString("\n")
		sb.WriteString(lipgloss.NewStyle().Foreground(styles.ColorTextPrimary).Render(m.event.Company))
		sb.WriteString("\n\n")
	}

	// Project (if available)
	if m.event.Project != "" {
		sb.WriteString(styles.InputLabel.Render("Project:"))
		sb.WriteString("\n")
		sb.WriteString(lipgloss.NewStyle().Foreground(styles.ColorTextPrimary).Render(m.event.Project))
		sb.WriteString("\n\n")
	}

	// Tags (if available)
	if len(m.event.Tags) > 0 {
		sb.WriteString(styles.InputLabel.Render("Tags:"))
		sb.WriteString("\n")
		for _, tag := range m.event.Tags {
			sb.WriteString(styles.TagBase.Render(tag))
			sb.WriteString(" ")
		}
		sb.WriteString("\n\n")
	}

	// Categories (if available)
	if len(m.event.Categories) > 0 {
		sb.WriteString(styles.InputLabel.Render("Categories:"))
		sb.WriteString("\n")
		for _, category := range m.event.Categories {
			sb.WriteString(styles.TagBase.Render(category))
			sb.WriteString(" ")
		}
		sb.WriteString("\n\n")
	}

	// Event ID
	sb.WriteString(styles.InputLabel.Render("Event ID:"))
	sb.WriteString("\n")
	sb.WriteString(lipgloss.NewStyle().Foreground(styles.ColorTextMuted).Render(m.event.ID))
	sb.WriteString("\n\n")

	// Timestamps
	sb.WriteString(styles.InputLabel.Render("Created:"))
	sb.WriteString("\n")
	sb.WriteString(lipgloss.NewStyle().Foreground(styles.ColorTextMuted).Render(m.formatTime(m.event.CreatedAt)))
	sb.WriteString("\n\n")

	sb.WriteString(styles.InputLabel.Render("Last Updated:"))
	sb.WriteString("\n")
	sb.WriteString(lipgloss.NewStyle().Foreground(styles.ColorTextMuted).Render(m.formatTime(m.event.UpdatedAt)))
	sb.WriteString("\n\n")

	// Footer
	sb.WriteString(styles.InputHint.Render("Press 'esc' to return to list"))

	return sb.String()
}

// formatTime formats a timestamp for display
func (m *DetailsModel) formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// Event returns the event being displayed
func (m *DetailsModel) Event() *career.CareerEvent {
	return m.event
}
