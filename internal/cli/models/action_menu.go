package models

import (
	"strings"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// EventAction represents the possible actions for an event
type EventAction int

const (
	EventActionView EventAction = iota
	EventActionEdit
	EventActionDelete
)

// ActionMenuModel represents the action selection menu for an event
type ActionMenuModel struct {
	event       *career.CareerEvent
	options     []EventAction
	selectedIdx int
	width       int
	height      int
	header      components.HeaderModel
	helpFooter  components.HelpFooterModel
}

// NewActionMenuModel creates a new action menu for a given event
func NewActionMenuModel(event *career.CareerEvent) *ActionMenuModel {
	return &ActionMenuModel{
		event:       event,
		options:     []EventAction{EventActionView, EventActionEdit, EventActionDelete},
		selectedIdx: 0,
		width:       40,
		height:      10,
		header:      components.NewHeader("Event Actions", 80),
		helpFooter:  components.NewHelpFooter("action_menu", 80),
	}
}

// Init initializes the model
func (m *ActionMenuModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the action menu
func (m *ActionMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.header.SetWidth(msg.Width)
		m.helpFooter.SetWidth(msg.Width)
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selectedIdx > 0 {
				m.selectedIdx--
			}
		case "down", "j":
			if m.selectedIdx < len(m.options)-1 {
				m.selectedIdx++
			}
		case "enter":
			// Return the selected action
			return m, func() tea.Msg {
				return EventActionSelectedMsg{
					Event:  m.event,
					Action: m.options[m.selectedIdx],
				}
			}
		case "esc", "backspace":
			// Go back to list
			return m, func() tea.Msg { return BackMsg{} }
		}
	}
	return m, nil
}

// View renders the action menu
func (m *ActionMenuModel) View() string {
	// Header
	headerContent := m.header.View()

	// Event summary
	eventSummary := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		Render(m.truncateText(m.event.Text, 40))

	// Render options
	var optionsContent strings.Builder
	actionLabels := []string{"View Event", "Edit Event", "Delete Event"}
	for i, option := range m.options {
		var optionStyle lipgloss.Style
		if i == m.selectedIdx {
			optionStyle = styles.ListItemSelected
			optionsContent.WriteString("▶ ")
		} else {
			optionStyle = styles.ListItem
			optionsContent.WriteString("  ")
		}

		optionsContent.WriteString(optionStyle.Render(actionLabels[option]))
		optionsContent.WriteString("\n")
	}

	// Help footer
	m.helpFooter.SetWidth(styles.MaxWidth(m.width))
	helpFooterContent := m.helpFooter.View()

	// Combine all content
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		headerContent,
		"",
		eventSummary,
		"",
		optionsContent.String(),
		helpFooterContent,
	)

	// Wrap in a card
	card := styles.CardBase.
		Width(styles.MaxWidth(m.width) - 4).
		Render(content)

	return card
}

// truncateText helps truncate long text for display
func (m *ActionMenuModel) truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen-3] + "..."
}

// EventActionSelectedMsg is sent when an action is selected
// SelectedAction returns the currently selected action
func (m *ActionMenuModel) SelectedAction() EventAction {
	return m.options[m.selectedIdx]
}

type EventActionSelectedMsg struct {
	Event  *career.CareerEvent
	Action EventAction
}
