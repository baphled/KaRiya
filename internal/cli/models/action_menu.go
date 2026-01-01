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
	*BaseStandardModel
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
		BaseStandardModel: NewBaseStandardModel(),
		event:             event,
		options:           []EventAction{EventActionView, EventActionEdit, EventActionDelete},
		selectedIdx:       0,
		width:             40,
		height:            10,
		header:            components.NewHeader("Event Actions", 80),
		helpFooter:        components.NewHelpFooter("action_menu", 80),
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

// View renders the action menu using the ScreenContainer pattern
func (m *ActionMenuModel) View() string {
	headerContent := m.renderHeader()
	contentArea := m.renderContent()
	footerContent := m.renderFooter()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		headerContent,
		"",
		contentArea,
		"",
		footerContent,
	)
}

// renderHeader renders the header section
func (m *ActionMenuModel) renderHeader() string {
	m.header.SetWidth(m.width)
	return m.header.View()
}

// renderContent renders the main content (event summary and action options)
func (m *ActionMenuModel) renderContent() string {
	var content []string

	// Event summary
	eventSummary := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		Render(m.truncateText(m.event.Text, 40))
	content = append(content, eventSummary)
	content = append(content, "")

	// Render action options
	actionLabels := []string{"View Event", "Edit Event", "Delete Event"}
	for i, option := range m.options {
		var optionStyle lipgloss.Style
		var prefix string
		if i == m.selectedIdx {
			optionStyle = styles.ListItemSelected
			prefix = "▶ "
		} else {
			optionStyle = styles.ListItem
			prefix = "  "
		}

		optionText := prefix + actionLabels[option]
		content = append(content, optionStyle.Render(optionText))
	}

	// Wrap content in ScreenContainer for consistent padding
	screenContainer := components.NewScreenContainer(strings.Join(content, "\n")).
		WithPaddingMode(components.PaddingNormal)

	return screenContainer.Render()
}

// renderFooter renders the footer section
func (m *ActionMenuModel) renderFooter() string {
	m.helpFooter.SetWidth(m.width)
	return m.helpFooter.View()
}

// truncateText helps truncate long text for display
func (m *ActionMenuModel) truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen-3] + "..."
}

// SelectedAction returns the currently selected action
func (m *ActionMenuModel) SelectedAction() EventAction {
	return m.options[m.selectedIdx]
}

// EventActionSelectedMsg is sent when an action is selected
type EventActionSelectedMsg struct {
	Event  *career.CareerEvent
	Action EventAction
}

// FactAction represents the possible actions for a fact
type FactAction int

const (
	FactActionView FactAction = iota
	FactActionEdit
	FactActionDelete
)

// FactActionSelectedMsg is sent when an action is selected for a fact
type FactActionSelectedMsg struct {
	Fact   *career.Fact
	Action FactAction
}

// FactActionMenuModel represents the action selection menu for a fact
type FactActionMenuModel struct {
	*BaseStandardModel
	fact        *career.Fact
	options     []FactAction
	selectedIdx int
	width       int
	height      int
	header      components.HeaderModel
	helpFooter  components.HelpFooterModel
}

// NewFactActionMenuModel creates a new action menu for a given fact
func NewFactActionMenuModel(fact *career.Fact) *FactActionMenuModel {
	return &FactActionMenuModel{
		BaseStandardModel: NewBaseStandardModel(),
		fact:              fact,
		options:           []FactAction{FactActionView, FactActionEdit, FactActionDelete},
		selectedIdx:       0,
		width:             40,
		height:            10,
		header:            components.NewHeader("Fact Actions", 80),
		helpFooter:        components.NewHelpFooter("action_menu", 80),
	}
}

// Init initializes the model
func (m *FactActionMenuModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the fact action menu
func (m *FactActionMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				return FactActionSelectedMsg{
					Fact:   m.fact,
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

// View renders the fact action menu
func (m *FactActionMenuModel) View() string {
	headerContent := m.renderHeader()
	contentArea := m.renderContent()
	footerContent := m.renderFooter()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		headerContent,
		"",
		contentArea,
		"",
		footerContent,
	)
}

// renderHeader renders the header section
func (m *FactActionMenuModel) renderHeader() string {
	m.header.SetWidth(m.width)
	return m.header.View()
}

// renderContent renders the main content (fact summary and action options)
func (m *FactActionMenuModel) renderContent() string {
	var content []string

	// Fact summary
	factSummary := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		Render(m.truncateText(m.fact.Text, 40))
	content = append(content, factSummary)
	content = append(content, "")

	// Render action options
	actionLabels := []string{"View Fact", "Edit Fact", "Delete Fact"}
	for i, option := range m.options {
		var optionStyle lipgloss.Style
		var prefix string
		if i == m.selectedIdx {
			optionStyle = styles.ListItemSelected
			prefix = "▶ "
		} else {
			optionStyle = styles.ListItem
			prefix = "  "
		}

		optionText := prefix + actionLabels[option]
		content = append(content, optionStyle.Render(optionText))
	}

	// Wrap content in ScreenContainer for consistent padding
	screenContainer := components.NewScreenContainer(strings.Join(content, "\n")).
		WithPaddingMode(components.PaddingNormal)

	return screenContainer.Render()
}

// renderFooter renders the footer section
func (m *FactActionMenuModel) renderFooter() string {
	m.helpFooter.SetWidth(m.width)
	return m.helpFooter.View()
}

// truncateText helps truncate long text for display
func (m *FactActionMenuModel) truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen-3] + "..."
}

// SelectedAction returns the currently selected action
func (m *FactActionMenuModel) SelectedAction() FactAction {
	return m.options[m.selectedIdx]
}

// BurstAction represents the possible actions for a burst
type BurstAction int

const (
	BurstActionView BurstAction = iota
	BurstActionEdit
	BurstActionDelete
)

// BurstActionSelectedMsg is sent when an action is selected for a burst
type BurstActionSelectedMsg struct {
	Burst  *career.Burst
	Action BurstAction
}
