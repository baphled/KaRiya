// Package modals provides modal dialogs for fact management.
package modals

import (
	"github.com/baphled/kariya/internal/cli/themes"
	domain "github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

// EditResult contains the result of an edit operation.
type EditResult struct {
	Accepted bool
	Fact     *domain.Fact
}

// EditModal provides a form for editing facts.
type EditModal struct {
	fact    *domain.Fact
	isNew   bool
	visible bool
	theme   themes.Theme
}

// NewEditModal creates a new edit modal.
func NewEditModal(fact *domain.Fact, isNew bool) *EditModal {
	// Create a copy of the fact for editing.
	editFact := &domain.Fact{
		ID:                   fact.ID,
		Text:                 fact.Text,
		CompetencyCategories: append([]string{}, fact.CompetencyCategories...),
		RoleFit:              fact.RoleFit,
		AudienceRelevance:    append([]string{}, fact.AudienceRelevance...),
		StrengthSignal:       fact.StrengthSignal,
		SourceEventID:        fact.SourceEventID,
		SourceBurstID:        fact.SourceBurstID,
		CreatedAt:            fact.CreatedAt,
		UpdatedAt:            fact.UpdatedAt,
	}

	return &EditModal{
		fact:    editFact,
		isNew:   isNew,
		visible: false,
	}
}

// SetTheme sets the theme for the modal.
func (m *EditModal) SetTheme(theme themes.Theme) {
	m.theme = theme
}

// Show makes the modal visible.
func (m *EditModal) Show() {
	m.visible = true
}

// Hide hides the modal.
func (m *EditModal) Hide() {
	m.visible = false
}

// IsVisible returns whether the modal is visible.
func (m *EditModal) IsVisible() bool {
	return m.visible
}

// Update processes messages and returns the result if completed.
func (m *EditModal) Update(msg tea.Msg) (tea.Cmd, *EditResult) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "esc":
			m.Hide()
			return nil, &EditResult{Accepted: false}
		case "ctrl+s":
			m.Hide()
			return nil, &EditResult{
				Accepted: true,
				Fact:     m.fact,
			}
		}
	}
	return nil, nil
}

// View renders the modal.
func (m *EditModal) View() string {
	if !m.visible {
		return ""
	}

	title := "Edit Fact"
	if m.isNew {
		title = "New Fact"
	}

	return title + "\n\n[Form placeholder - will be implemented with huh forms]"
}

// GetContent returns the content for embedding.
func (m *EditModal) GetContent() string {
	return m.View()
}

// GetFact returns the fact being edited.
func (m *EditModal) GetFact() *domain.Fact {
	return m.fact
}

// IsNew returns whether this is a new fact.
func (m *EditModal) IsNew() bool {
	return m.isNew
}
