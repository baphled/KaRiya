package models

import (
	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// BurstDetailsModel represents the burst details view screen
type BurstDetailsModel struct {
	*BaseStandardModel
	burst      *career.Burst
	width      int
	height     int
	header     components.HeaderModel
	helpFooter components.HelpFooterModel
	burstCard  *BurstCard
}

// NewBurstDetailsModel creates a new burst details model
func NewBurstDetailsModel(burst *career.Burst) *BurstDetailsModel {
	return &BurstDetailsModel{
		BaseStandardModel: NewBaseStandardModel(),
		burst:             burst,
		width:             80,
		height:            24,
		header:            components.NewHeader("Burst Details", 80),
		helpFooter:        components.NewHelpFooter("burst_details", 80),
		burstCard:         NewBurstCard(burst),
	}
}

// Init initializes the model
func (m *BurstDetailsModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m *BurstDetailsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.header.SetWidth(msg.Width)
		m.helpFooter.SetWidth(msg.Width)
		if m.burstCard != nil {
			m.burstCard.SetWidth(msg.Width - 4)
		}
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "backspace":
			// Go back to burst list
			return m, func() tea.Msg { return BackMsg{} }
		case "e":
			// Edit the burst
			return m, func() tea.Msg { return BurstActionSelectedMsg{Burst: m.burst, Action: BurstActionEdit} }
		case "d", "x":
			// Delete the burst
			return m, func() tea.Msg { return BurstActionSelectedMsg{Burst: m.burst, Action: BurstActionDelete} }
		case "q", "ctrl+c":
			// Quit
			return m, func() tea.Msg { return QuitMsg{} }
		}
	}
	return m, nil
}

// View renders the burst details
func (m *BurstDetailsModel) View() string {
	headerView := m.header.View()
	m.helpFooter.SetWidth(m.width)
	footerView := m.helpFooter.View()

	// Render burst card
	var burstContent string
	if m.burstCard != nil {
		m.burstCard.SetWidth(m.width - 4)
		burstContent = m.burstCard.Render()
	}

	// Combine all sections
	fullContent := lipgloss.JoinVertical(
		lipgloss.Left,
		headerView,
		"",
		burstContent,
		"",
		footerView,
	)

	return fullContent
}

