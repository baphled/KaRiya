package models

import (
	"context"
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/domain/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// FactProcessingCompleteMsg is sent when fact review workflow is done
type FactProcessingCompleteMsg struct {
	ConfirmedCount int
	RejectedCount  int
}

// FactsResultsModel represents the fact extraction results review and confirmation screen
// This model displays facts extracted from imported events and allows users to confirm or reject them
type FactsResultsModel struct {
	service      *careerservice.Service
	ctx          context.Context
	facts        []*career.Fact // Facts extracted from imported events
	currentIdx   int            // Currently focused fact index
	confirmed    []*career.Fact // Facts user confirmed
	rejected     []*career.Fact // Facts user rejected
	width        int
	height       int
	helpFooter   components.HelpFooterModel // Help footer for keyboard shortcuts
	scrollOffset int
}

// NewFactsResultsModel creates a new facts results model
func NewFactsResultsModel(svc *careerservice.Service, facts []*career.Fact, ctx context.Context) *FactsResultsModel {
	return &FactsResultsModel{
		service:      svc,
		ctx:          ctx,
		facts:        facts,
		currentIdx:   0,
		confirmed:    []*career.Fact{},
		rejected:     []*career.Fact{},
		width:        80,
		height:       24,
		scrollOffset: 0,
	}
}

// Init initializes the model
func (m *FactsResultsModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m *FactsResultsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			// Move to previous fact
			if m.currentIdx > 0 {
				m.currentIdx--
				if m.currentIdx < m.scrollOffset {
					m.scrollOffset = m.currentIdx
				}
			}
			return m, nil

		case "down", "j":
			// Move to next fact
			if m.currentIdx < len(m.facts)-1 {
				m.currentIdx++
				// Ensure we can see the current item
				visibleItems := m.height - 8 // Account for header, footer, spacing
				if m.currentIdx >= m.scrollOffset+visibleItems {
					m.scrollOffset = m.currentIdx - visibleItems + 1
				}
			}
			return m, nil

		case "y", "enter":
			// Confirm current fact
			if m.currentIdx < len(m.facts) {
				fact := m.facts[m.currentIdx]
				m.confirmed = append(m.confirmed, fact)
				m.removeFactAt(m.currentIdx)
				// Adjust index if at end
				if m.currentIdx >= len(m.facts) && m.currentIdx > 0 {
					m.currentIdx--
				}
				if len(m.facts) == 0 {
					return m, tea.Quit
				}
			}
			return m, nil

		case "n", "d":
			// Reject current fact
			if m.currentIdx < len(m.facts) {
				fact := m.facts[m.currentIdx]
				m.rejected = append(m.rejected, fact)
				m.removeFactAt(m.currentIdx)
				// Adjust index if at end
				if m.currentIdx >= len(m.facts) && m.currentIdx > 0 {
					m.currentIdx--
				}
				if len(m.facts) == 0 {
					return m, tea.Quit
				}
			}
			return m, nil

		case "q", "esc":
			// Skip all remaining facts
			return m, tea.Quit
		}
	}

	return m, nil
}

// View renders the model
func (m *FactsResultsModel) View() string {
	if len(m.facts) == 0 {
		// Show completion message
		title := styles.HeaderMain.Render("Fact Extraction Complete")
		summary := fmt.Sprintf("Confirmed: %d | Rejected: %d", len(m.confirmed), len(m.rejected))
		summaryStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("243")).
			Margin(1, 0)
		return lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			summaryStyle.Render(summary),
			"",
			"Press any key to continue...",
		)
	}

	// Display current fact
	var factDisplays []string

	// Header
	header := styles.HeaderMain.Render("Review Extracted Facts")
	progress := fmt.Sprintf("[%d/%d] Facts", m.currentIdx+1, len(m.facts)+len(m.confirmed)+len(m.rejected))
	progressStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		PaddingTop(1)
	factDisplays = append(factDisplays, header)
	factDisplays = append(factDisplays, progressStyle.Render(progress))

	// Current fact
	if m.currentIdx < len(m.facts) {
		fact := m.facts[m.currentIdx]
		factDisplay := m.renderFactCard(fact)
		factDisplays = append(factDisplays, factDisplay)
	}

	// Help footer
	helpText := "Y/Enter: Confirm | N/D: Reject | Q/Esc: Skip remaining"
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		PaddingTop(1)
	factDisplays = append(factDisplays, helpStyle.Render(helpText))

	return lipgloss.JoinVertical(lipgloss.Left, factDisplays...)
}

// renderFactCard renders a single fact as a card
func (m *FactsResultsModel) renderFactCard(fact *career.Fact) string {
	var cardContent []string

	// Title
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86"))
	cardContent = append(cardContent, titleStyle.Render(fact.Text))

	// Competency categories
	if len(fact.CompetencyCategories) > 0 {
		compStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")).
			MarginTop(1)
		compText := fmt.Sprintf("Competencies: %s", strings.Join(fact.CompetencyCategories, ", "))
		cardContent = append(cardContent, compStyle.Render(compText))
	}

	// Role fit
	if fact.RoleFit != "" {
		roleStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("177")).
			MarginTop(1)
		roleText := fmt.Sprintf("Role Fit: %s", fact.RoleFit)
		cardContent = append(cardContent, roleStyle.Render(roleText))
	}

	// Audience relevance
	if len(fact.AudienceRelevance) > 0 {
		audienceStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("51")).
			MarginTop(1)
		audienceText := fmt.Sprintf("Audience: %s", strings.Join(fact.AudienceRelevance, ", "))
		cardContent = append(cardContent, audienceStyle.Render(audienceText))
	}

	// Strength signal
	if fact.StrengthSignal != "" {
		strengthStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("172")).
			MarginTop(1)
		strengthText := fmt.Sprintf("Strength: %s", fact.StrengthSignal)
		cardContent = append(cardContent, strengthStyle.Render(strengthText))
	}

	// Wrap in card
	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(1).
		MarginTop(1).
		Width(m.width - 4)

	return cardStyle.Render(lipgloss.JoinVertical(lipgloss.Left, cardContent...))
}

// removeFactAt removes the fact at the given index
func (m *FactsResultsModel) removeFactAt(idx int) {
	if idx >= 0 && idx < len(m.facts) {
		m.facts = append(m.facts[:idx], m.facts[idx+1:]...)
	}
}

// GetConfirmed returns the facts user confirmed
func (m *FactsResultsModel) GetConfirmed() []*career.Fact {
	return m.confirmed
}

// GetRejected returns the facts user rejected
func (m *FactsResultsModel) GetRejected() []*career.Fact {
	return m.rejected
}

// IsDone returns true if all facts have been reviewed
func (m *FactsResultsModel) IsDone() bool {
	return len(m.facts) == 0
}
