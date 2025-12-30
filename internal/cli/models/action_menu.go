package models

import (
	"strings"

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
}

// NewActionMenuModel creates a new action menu for a given event
func NewActionMenuModel(event *career.CareerEvent) *ActionMenuModel {
	return &ActionMenuModel{
		event:       event,
		options:     []EventAction{EventActionView, EventActionEdit, EventActionDelete},
		selectedIdx: 0,
		width:       40,
		height:      10,
	}
}

// Init initializes the model
func (m *ActionMenuModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the action menu
func (m *ActionMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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
		case "backspace", "esc":
			// Go back to list
			return m, func() tea.Msg { return BackMsg{} }
		}
	}
	return m, nil
}

// View renders the action menu
func (m *ActionMenuModel) View() string {
	var sb strings.Builder

	// Title
	title := styles.HeaderSection.Render("Event Actions")
	sb.WriteString(title)
	sb.WriteString("\n\n")

	// Event summary
	eventSummary := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		Render(m.truncateText(m.event.Text, 40))
	sb.WriteString(eventSummary)
	sb.WriteString("\n\n")

	// Render options
	actionLabels := []string{"View Event", "Edit Event", "Delete Event"}
	for i, option := range m.options {
		var optionStyle lipgloss.Style
		if i == m.selectedIdx {
			optionStyle = styles.ListItemSelected
			sb.WriteString("▶ ")
		} else {
			optionStyle = styles.ListItem
			sb.WriteString("  ")
		}

		sb.WriteString(optionStyle.Render(actionLabels[option]))
		sb.WriteString("\n")
	}

	// Instructions
	sb.WriteString("\n")
	instructions := styles.InfoHint.Render("↑/↓ or j/k: Navigate | Enter: Select | Backspace: Cancel")
	sb.WriteString(instructions)

	// Wrap in a card
	card := styles.CardBase.
		Width(m.width).
		Render(sb.String())

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
