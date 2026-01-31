// Package modals provides modal components for the fact management screens.
package modals

import (
	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// EditResult captures the outcome of an edit operation.
type EditResult struct {
	Original *career.Fact
	Modified *career.Fact
	Accepted bool
	Changes  map[string]interface{}
}

// HasChanges returns true if there are any changes.
func (r *EditResult) HasChanges() bool {
	return len(r.Changes) > 0
}

// EditFactModal handles inline editing of fact details.
type EditFactModal struct {
	original *career.Fact
	modified *career.Fact
	result   *EditResult
	form     forms.Form
	formData *forms.FactFormData
	width    int
	height   int
}

// NewEditFactModal creates a new fact editing modal using huh forms.
func NewEditFactModal(fact *career.Fact) *EditFactModal {
	originalCopy := *fact

	formData := forms.GetFactFormData(&originalCopy)

	defaultWidth := 80
	defaultHeight := 24

	form := forms.NewFactEditorFormWithDataAndDimensions(
		formData,
		defaultWidth-4,
		forms.DefaultFormHeight(defaultHeight),
	)

	return &EditFactModal{
		original: &originalCopy,
		modified: &originalCopy,
		result:   nil,
		form:     form,
		formData: formData,
		width:    defaultWidth,
		height:   defaultHeight,
	}
}

// Init initializes the modal's form and returns the init command.
func (m *EditFactModal) Init() tea.Cmd {
	return m.form.Init()
}

// Update handles user input for fact editing.
func (m *EditFactModal) Update(msg tea.Msg) tea.Cmd {
	if wsm, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = wsm.Width
		m.height = wsm.Height
		m.form = m.form.
			WithHeight(forms.DefaultFormHeight(m.height)).
			WithWidth(m.width - 4)
		return nil
	}

	var cmd tea.Cmd
	m.form, cmd = forms.Update(m.form, msg)

	if forms.IsCompleted(m.form) {
		m.syncModified()
		m.createResult()
		return nil
	}

	if forms.IsAborted(m.form) {
		m.createCancelledResult()
		return nil
	}

	return cmd
}

// View renders the fact editing modal with professional styling.
func (m *EditFactModal) View() string {
	if m.result != nil && m.result.Accepted {
		return ""
	}

	formView := m.form.View()

	title := modalTitleStyle(nil).
		Render("Edit Fact")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		formView,
	)

	modal := feedback.NewModalContainer().
		SetTitle("").
		SetMessage(content).
		SetInstructions("Tab: Next  |  Shift+Tab: Prev  |  Enter: Confirm  |  Esc: Cancel").
		WithWidth(m.width - 4).
		WithScrollHint(true)

	return modal.Render()
}

// Result returns the modal result when editing is complete.
func (m *EditFactModal) Result() *EditResult {
	return m.result
}

// IsComplete returns true if the modal has finished.
func (m *EditFactModal) IsComplete() bool {
	return m.result != nil
}

// GetTitle returns the modal title for overlay rendering.
func (m *EditFactModal) GetTitle() string {
	return "Edit Fact"
}

// GetContent returns just the form content without the modal container.
func (m *EditFactModal) GetContent() string {
	if m.result != nil && m.result.Accepted {
		return ""
	}
	return m.form.View()
}

// GetFooter returns the footer instructions for the modal.
func (m *EditFactModal) GetFooter() string {
	return "Enter: Confirm  |  Esc: Cancel  |  Tab: Next Field  |  Shift+Tab: Previous"
}

func modalTitleStyle(theme themes.Theme) lipgloss.Style {
	if theme == nil {
		theme = themes.NewDefaultTheme()
	}
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.ForegroundColor()).
		MarginBottom(1)
}

func (m *EditFactModal) syncModified() {
	m.modified = &career.Fact{
		ID:                   m.original.ID,
		Text:                 m.formData.Text,
		CompetencyCategories: m.formData.CompetencyCategories,
		RoleFit:              career.RoleFit(m.formData.RoleFit),
		AudienceRelevance:    m.formData.AudienceRelevance,
		StrengthSignal:       m.original.StrengthSignal,
		SourceEventID:        m.original.SourceEventID,
		SourceBurstID:        m.original.SourceBurstID,
		CreatedAt:            m.original.CreatedAt,
		UpdatedAt:            m.original.UpdatedAt,
	}
}

func (m *EditFactModal) createResult() {
	if !m.formData.SubmitConfirmed {
		m.createCancelledResult()
		return
	}

	m.result = &EditResult{
		Original: m.original,
		Modified: m.modified,
		Accepted: true,
		Changes:  m.computeChanges(),
	}
}

func (m *EditFactModal) createCancelledResult() {
	m.result = &EditResult{
		Original: m.original,
		Modified: m.original,
		Accepted: false,
		Changes:  make(map[string]interface{}),
	}
}

func (m *EditFactModal) computeChanges() map[string]interface{} {
	changes := make(map[string]interface{})

	if m.original.Text != m.modified.Text {
		changes["text"] = m.modified.Text
	}
	if !slicesEqual(m.original.CompetencyCategories, m.modified.CompetencyCategories) {
		changes["competency_categories"] = m.modified.CompetencyCategories
	}
	if m.original.RoleFit != m.modified.RoleFit {
		changes["role_fit"] = m.modified.RoleFit
	}
	if !slicesEqual(m.original.AudienceRelevance, m.modified.AudienceRelevance) {
		changes["audience_relevance"] = m.modified.AudienceRelevance
	}
	if m.original.StrengthSignal != m.modified.StrengthSignal {
		changes["strength_signal"] = m.modified.StrengthSignal
	}

	return changes
}

func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
