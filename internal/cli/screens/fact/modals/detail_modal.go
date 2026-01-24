// Package modals provides modal dialogs for fact management.
package modals

import (
	"fmt"

	"github.com/baphled/kariya/internal/cli/themes"
	domain "github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

// DetailModal displays fact details.
type DetailModal struct {
	fact    *domain.Fact
	visible bool
	theme   themes.Theme
}

// NewDetailModal creates a new detail modal.
func NewDetailModal(fact *domain.Fact) *DetailModal {
	return &DetailModal{
		fact:    fact,
		visible: false,
	}
}

// SetTheme sets the theme for the modal.
func (m *DetailModal) SetTheme(theme themes.Theme) {
	m.theme = theme
}

// Show makes the modal visible.
func (m *DetailModal) Show() {
	m.visible = true
}

// Hide hides the modal.
func (m *DetailModal) Hide() {
	m.visible = false
}

// IsVisible returns whether the modal is visible.
func (m *DetailModal) IsVisible() bool {
	return m.visible
}

// Update processes messages.
func (m *DetailModal) Update(msg tea.Msg) (bool, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "esc", "q":
			m.Hide()
			return false, nil
		}
	}
	return true, nil
}

// View renders the modal.
func (m *DetailModal) View() string {
	if !m.visible || m.fact == nil {
		return ""
	}

	var content string
	content += "Fact Details\n\n"
	content += fmt.Sprintf("Text: %s\n\n", m.fact.Text)
	content += fmt.Sprintf("Categories: %v\n", m.fact.CompetencyCategories)
	content += fmt.Sprintf("Strength Signal: %s\n", m.fact.StrengthSignal)
	content += fmt.Sprintf("Role Fit: %v\n", m.fact.RoleFit)
	content += fmt.Sprintf("Audience Relevance: %v\n", m.fact.AudienceRelevance)
	if m.fact.CreatedAt.Year() > 1 {
		content += fmt.Sprintf("\nCreated: %s\n", m.fact.CreatedAt.Format("2006-01-02 15:04"))
	}

	return content
}

// GetFact returns the fact being displayed.
func (m *DetailModal) GetFact() *domain.Fact {
	return m.fact
}
