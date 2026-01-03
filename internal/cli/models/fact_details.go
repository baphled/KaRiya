package models

import (
	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// FactDetailsModel represents the fact details view screen
type FactDetailsModel struct {
	*BaseStandardModel
	fact       *career.Fact
	width      int
	height     int
	header     components.HeaderModel
	helpFooter components.HelpFooterModel
	factCard   *FactCard
}

// NewFactDetailsModel creates a new fact details model
func NewFactDetailsModel(fact *career.Fact) *FactDetailsModel {
	return &FactDetailsModel{
		BaseStandardModel: NewBaseStandardModel(),
		fact:              fact,
		width:             80,
		height:            24,
		header:            components.NewHeader("Fact Details", 80),
		helpFooter:        components.NewHelpFooter("fact_details", 80),
		factCard:          NewFactCard(fact),
	}
}

// Init initializes the model
func (m *FactDetailsModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m *FactDetailsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.header.SetWidth(msg.Width)
		m.helpFooter.SetWidth(msg.Width)
		if m.factCard != nil {
			m.factCard.SetWidth(msg.Width - 4)
		}
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "backspace":
			// Go back to fact list
			return m, func() tea.Msg { return BackMsg{} }
		case "e":
			// Edit the fact
			return m, func() tea.Msg { return FactActionSelectedMsg{Fact: m.fact, Action: FactActionEdit} }
		case "d", "x":
			// Delete the fact
			return m, func() tea.Msg { return FactActionSelectedMsg{Fact: m.fact, Action: FactActionDelete} }
		case "q", "ctrl+c":
			// Quit
			return m, func() tea.Msg { return QuitMsg{} }
		}
	}
	return m, nil
}

// View renders the fact details
func (m *FactDetailsModel) View() string {
	headerView := m.header.View()
	m.helpFooter.SetWidth(m.width)
	footerView := m.helpFooter.View()

	// Render fact card
	var factContent string
	if m.factCard != nil {
		m.factCard.SetWidth(m.width - 4)
		factContent = m.factCard.Render()
	}

	// Combine all sections
	fullContent := lipgloss.JoinVertical(
		lipgloss.Left,
		headerView,
		"",
		factContent,
		"",
		footerView,
	)

	return fullContent
}
