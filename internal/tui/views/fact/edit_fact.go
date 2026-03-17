// Package fact provides view components for fact management.
package fact

import (
	"github.com/baphled/kariya/internal/tui/forms"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/themes"
	"github.com/baphled/kariya/internal/ui/uikit/feedback"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// EditChanges captures the fields that were modified during editing.
type EditChanges struct {
	Text                 *string
	CompetencyCategories *[]string
	RoleFit              *string
	AudienceRelevance    *[]string
	StrengthSignal       *string
}

// EditResult captures the outcome of an edit operation.
type EditResult struct {
	Original display.Fact
	Modified display.Fact
	Accepted bool
	Changes  EditChanges
}

// HasChanges returns true if there are any changes.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (r *EditResult) HasChanges() bool {
	return r.Changes.Text != nil ||
		r.Changes.CompetencyCategories != nil ||
		r.Changes.RoleFit != nil ||
		r.Changes.AudienceRelevance != nil ||
		r.Changes.StrengthSignal != nil
}

// EditFact handles inline editing of fact details.
type EditFact struct {
	original display.Fact
	modified display.Fact
	result   *EditResult
	form     forms.Form
	formData *forms.FactFormData
	width    int
	height   int
}

// NewEditFact creates a new fact editing modal using huh forms.
//
// Expected:
//   - fact must be valid.
//
// Returns:
//   - A fully initialized EditFact ready for use.
//
// Side effects:
//   - None.
func NewEditFact(fact display.Fact) *EditFact {
	originalCopy := fact

	formData := factFormDataFromDisplay(originalCopy)

	defaultWidth := 80
	defaultHeight := 24

	form := forms.NewFactForm(
		formData,
		defaultWidth-4,
		forms.DefaultFormHeight(defaultHeight),
	)

	return &EditFact{
		original: originalCopy,
		modified: originalCopy,
		result:   nil,
		form:     form,
		formData: formData,
		width:    defaultWidth,
		height:   defaultHeight,
	}
}

// Init initializes the modal's form and returns the init command.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (m *EditFact) Init() tea.Cmd {
	return m.form.Init()
}

// Update handles user input for fact editing.
//
// Expected:
//   - msg must be valid.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (m *EditFact) Update(msg tea.Msg) tea.Cmd {
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
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *EditFact) View() string {
	if m.result != nil && m.result.Accepted {
		return ""
	}

	formView := m.form.View()

	title := titleStyle(nil).
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
//
// Returns:
//   - A fully initialized EditResult ready for use.
//
// Side effects:
//   - None.
func (m *EditFact) Result() *EditResult {
	return m.result
}

// IsComplete returns true if the modal has finished.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *EditFact) IsComplete() bool {
	return m.result != nil
}

// GetTitle returns the modal title for overlay rendering.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *EditFact) GetTitle() string {
	return "Edit Fact"
}

// GetContent returns just the form content without the modal container.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *EditFact) GetContent() string {
	if m.result != nil && m.result.Accepted {
		return ""
	}
	return m.form.View()
}

// GetFooter returns the footer instructions for the modal.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *EditFact) GetFooter() string {
	return "Enter: Confirm  |  Esc: Cancel  |  Tab: Next Field  |  Shift+Tab: Previous"
}

func titleStyle(theme themes.Theme) lipgloss.Style {
	if theme == nil {
		theme = themes.NewDefaultTheme()
	}
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.ForegroundColor()).
		MarginBottom(1)
}

func (m *EditFact) syncModified() {
	m.modified = display.Fact{
		ID:                   m.original.ID,
		Text:                 m.formData.Text,
		CompetencyCategories: m.formData.CompetencyCategories,
		RoleFit:              m.formData.RoleFit,
		AudienceRelevance:    m.formData.AudienceRelevance,
		StrengthSignal:       m.original.StrengthSignal,
		SourceEventID:        m.original.SourceEventID,
		SourceBurstID:        m.original.SourceBurstID,
		CreatedAt:            m.original.CreatedAt,
		UpdatedAt:            m.original.UpdatedAt,
	}
}

func (m *EditFact) createResult() {
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

func (m *EditFact) createCancelledResult() {
	m.result = &EditResult{
		Original: m.original,
		Modified: m.original,
		Accepted: false,
		Changes:  EditChanges{},
	}
}

func (m *EditFact) computeChanges() EditChanges {
	changes := EditChanges{}

	if m.original.Text != m.modified.Text {
		text := m.modified.Text
		changes.Text = &text
	}
	if !slicesEqual(m.original.CompetencyCategories, m.modified.CompetencyCategories) {
		categories := m.modified.CompetencyCategories
		changes.CompetencyCategories = &categories
	}
	if m.original.RoleFit != m.modified.RoleFit {
		roleFit := m.modified.RoleFit
		changes.RoleFit = &roleFit
	}
	if !slicesEqual(m.original.AudienceRelevance, m.modified.AudienceRelevance) {
		relevance := m.modified.AudienceRelevance
		changes.AudienceRelevance = &relevance
	}
	if m.original.StrengthSignal != m.modified.StrengthSignal {
		signal := m.modified.StrengthSignal
		changes.StrengthSignal = &signal
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

func factFormDataFromDisplay(fact display.Fact) *forms.FactFormData {
	return &forms.FactFormData{
		Text:                 fact.Text,
		CompetencyCategories: cloneStrings(fact.CompetencyCategories),
		RoleFit:              fact.RoleFit,
		AudienceRelevance:    cloneStrings(fact.AudienceRelevance),
	}
}

func cloneStrings(values []string) []string {
	if values == nil {
		return []string{}
	}

	return append([]string(nil), values...)
}
