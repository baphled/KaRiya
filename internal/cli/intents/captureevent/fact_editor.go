package captureevent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/domain/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// FactEditorModelNew represents the fact editor form state using huh.
type FactEditorModelNew struct {
	forms.EditorFields
	fact         *career.Fact
	originalFact *career.Fact
	service      *careerservice.Service
	ctx          context.Context
	formData     *forms.FactFormData
	submitted    bool
}

// NewFactEditorModelNew creates a new fact editor model using huh forms.
//
// Expected:
//   - fact must be valid.
//   - service must be valid.
//
// Returns:
//   - A fully initialized FactEditorModelNew ready for use.
//
// Side effects:
//   - None.
func NewFactEditorModelNew(ctx context.Context, fact *career.Fact, service *careerservice.Service) *FactEditorModelNew {
	// Create a copy of the fact for reverting
	factCopy := *fact

	// Extract form data from fact
	formData := forms.GetFactFormData(fact)

	// Create huh form
	form := forms.NewFactForm(formData, 0, 0)

	return &FactEditorModelNew{
		EditorFields: forms.EditorFields{
			Form:   form,
			Width:  80,
			Height: 24,
			Theme:  themes.NewDefaultTheme(),
		},
		fact:         fact,
		originalFact: &factCopy,
		service:      service,
		ctx:          ctx,
		formData:     formData,
	}
}

// Init initializes the model
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (m *FactEditorModelNew) Init() tea.Cmd {
	return m.Form.Init()
}

// Update handles messages.
func (m *FactEditorModelNew) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return forms.EditorUpdate(&m.EditorFields, m, msg, m.handleFactFormCompletion, func() tea.Msg { return QuitMsg{} })
}

// handleFactFormCompletion processes the completed form and saves the fact.
func (m *FactEditorModelNew) handleFactFormCompletion() (tea.Model, tea.Cmd) {
	// Check if user confirmed via the submit button
	// If they selected "Cancel" on the confirm, treat as cancelled
	if !m.formData.SubmitConfirmed {
		m.Cancelled = true
		return m, nil
	}

	// Apply form data to fact
	err := forms.ApplyFactFormData(m.fact, m.formData)
	if err != nil {
		m.Err = fmt.Errorf("failed to apply form data: %w", err)
		return m, nil
	}

	// Update timestamp
	m.fact.UpdatedAt = time.Now()

	// Save the fact using the repository
	factRepo := m.service.GetFactRepository()
	if factRepo == nil {
		m.Err = errors.New("fact repository not available")
		return m, nil
	}

	err = factRepo.Update(m.ctx, m.fact)
	if err != nil {
		m.Err = fmt.Errorf("failed to save fact: %w", err)
		return m, nil
	}

	m.submitted = true
	return m, nil
}

// GetFact returns the edited fact
//
// Returns:
//   - A fully initialized career.Fact ready for use.
//
// Side effects:
//   - None.
func (m *FactEditorModelNew) GetFact() *career.Fact {
	return m.fact
}

// IsSubmitted returns true if changes were saved
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *FactEditorModelNew) IsSubmitted() bool {
	return m.submitted
}

// IsCancelled returns true if operation was cancelled
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *FactEditorModelNew) IsCancelled() bool {
	return m.Cancelled
}

// Revert reverts changes to the original fact
//
// Side effects:
//   - None.
func (m *FactEditorModelNew) Revert() {
	m.fact.Text = m.originalFact.Text
	m.fact.CompetencyCategories = m.originalFact.CompetencyCategories
	m.fact.RoleFit = m.originalFact.RoleFit
	m.fact.AudienceRelevance = m.originalFact.AudienceRelevance
	m.fact.StrengthSignal = m.originalFact.StrengthSignal

	// Update form data
	m.formData = forms.GetFactFormData(m.originalFact)
}

// GetError returns the current error
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (m *FactEditorModelNew) GetError() error {
	return m.Err
}

// GetTitle returns the modal title for overlay rendering.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *FactEditorModelNew) GetTitle() string {
	return "Fact Editor"
}

// GetContent returns just the form content without header/footer.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *FactEditorModelNew) GetContent() string {
	formView := m.Form.View()

	// Add error if present
	if m.Err != nil {
		errorColor := m.Theme.ErrorColor()
		errorStyle := lipgloss.NewStyle().
			Foreground(errorColor).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(errorColor).
			Padding(1, 2).
			MarginTop(1)

		formView += "\n\n" + errorStyle.Render(m.Err.Error())
	}

	return formView
}

// GetFooter returns the footer instructions for the modal.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *FactEditorModelNew) GetFooter() string {
	return "Enter: Confirm | Esc: Cancel | Tab: Next Field | Shift+Tab: Previous"
}

// View renders the editor UI.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *FactEditorModelNew) View() string {
	return forms.RenderEditorView(&m.EditorFields, "Fact Editor", "fact_editor")
}
