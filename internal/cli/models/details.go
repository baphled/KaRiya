package models

import (
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// DetailsModel represents the event details screen
type DetailsModel struct {
	event      *career.CareerEvent
	width      int
	height     int
	header     components.HeaderModel
	helpFooter components.HelpFooterModel
}

// NewDetailsModel creates a new details model
func NewDetailsModel(event *career.CareerEvent) *DetailsModel {
	return &DetailsModel{
		event:      event,
		width:      80,
		height:     24,
		header:     components.NewHeader("Event Details", 80),
		helpFooter: components.NewHelpFooter("details", 80),
	}
}

// Init initializes the model
func (m *DetailsModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m *DetailsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.header.SetWidth(msg.Width)
		m.helpFooter.SetWidth(msg.Width)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "backspace":
			// Signal back navigation to parent
			return m, func() tea.Msg { return BackMsg{} }
		case "ctrl+c", "q":
			// Signal quit to parent
			return m, func() tea.Msg { return QuitMsg{} }
		}
	}
	return m, nil
}

// View renders the details view
func (m *DetailsModel) View() string {
	if m.event == nil {
		return styles.ErrorBox.Render("No event selected\n\nPress 'esc' to return")
	}

	// Header
	headerContent := m.header.View()

	var content []string

	// Event Text
	content = append(content,
		styles.InputLabel.Render("Description:"),
		lipgloss.NewStyle().Foreground(styles.ColorTextPrimary).Render(m.event.Text),
		"",
	)

	// Date
	content = append(content,
		styles.InputLabel.Render("Date:"),
		lipgloss.NewStyle().Foreground(styles.ColorTextPrimary).Render(m.event.Date.Format("2006-01-02")),
		"",
	)

	// Company (if available)
	if m.event.Company != "" {
		content = append(content,
			styles.InputLabel.Render("Company:"),
			lipgloss.NewStyle().Foreground(styles.ColorTextPrimary).Render(m.event.Company),
			"",
		)
	}

	// Project (if available)
	if m.event.Project != "" {
		content = append(content,
			styles.InputLabel.Render("Project:"),
			lipgloss.NewStyle().Foreground(styles.ColorTextPrimary).Render(m.event.Project),
			"",
		)
	}

	// Tags (if available)
	if len(m.event.Tags) > 0 {
		tagLine := styles.InputLabel.Render("Tags:")
		var tagContent strings.Builder
		tagContent.WriteString(tagLine)
		tagContent.WriteString("\n")
		for _, tag := range m.event.Tags {
			tagContent.WriteString(styles.TagBase.Render(tag))
			tagContent.WriteString(" ")
		}
		content = append(content, tagContent.String(), "")
	}

	// Categories (if available)
	if len(m.event.Categories) > 0 {
		categoryLine := styles.InputLabel.Render("Categories:")
		var categoryContent strings.Builder
		categoryContent.WriteString(categoryLine)
		categoryContent.WriteString("\n")
		for _, category := range m.event.Categories {
			categoryContent.WriteString(styles.TagBase.Render(category))
			categoryContent.WriteString(" ")
		}
		content = append(content, categoryContent.String(), "")
	}

	// Event ID
	content = append(content,
		styles.InputLabel.Render("Event ID:"),
		lipgloss.NewStyle().Foreground(styles.ColorTextMuted).Render(m.event.ID),
		"",
	)

	// Timestamps
	content = append(content,
		styles.InputLabel.Render("Created:"),
		lipgloss.NewStyle().Foreground(styles.ColorTextMuted).Render(m.formatTime(m.event.CreatedAt)),
		"",
	)

	content = append(content,
		styles.InputLabel.Render("Last Updated:"),
		lipgloss.NewStyle().Foreground(styles.ColorTextMuted).Render(m.formatTime(m.event.UpdatedAt)),
		"",
	)

	// Help footer
	m.helpFooter.SetWidth(styles.MaxWidth(m.width))
	helpFooterContent := m.helpFooter.View()
	content = append(content, helpFooterContent)

	// Combine all content
	detailsContent := lipgloss.JoinVertical(
		lipgloss.Left,
		content...,
	)

	// Wrap in card with header
	card := styles.CardBase.
		Width(styles.MaxWidth(m.width) - 4).
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
			headerContent,
			"",
			detailsContent,
		))

	return card
}

// formatTime formats a timestamp for display
func (m *DetailsModel) formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// Event returns the event being displayed
func (m *DetailsModel) Event() *career.CareerEvent {
	return m.event
}
