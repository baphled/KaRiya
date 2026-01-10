package models

import (
	"context"
	"fmt"
	"time"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/domain/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// FactEditorModel represents the fact editor form state using huh.
type FactEditorModelNew struct {
	*BaseStandardModel
	fact         *career.Fact
	originalFact *career.Fact // For reverting changes
	service      *careerservice.Service
	ctx          context.Context
	form         *huh.Form
	formData     *forms.FactFormData
	err          error
	submitted    bool
	cancelled    bool
	width        int
	height       int
	helpFooter   components.HelpFooterModel
}

// NewFactEditorModelNew creates a new fact editor model using huh forms.
func NewFactEditorModelNew(fact *career.Fact, service *careerservice.Service, ctx context.Context) *FactEditorModelNew {
	// Create a copy of the fact for reverting
	factCopy := *fact

	// Extract form data from fact
	formData := forms.GetFactFormData(fact)

	// Create huh form
	form := forms.NewFactEditorFormWithData(formData)

	return &FactEditorModelNew{
		BaseStandardModel: NewBaseStandardModel(),
		fact:              fact,
		originalFact:      &factCopy,
		service:           service,
		ctx:               ctx,
		form:              form,
		formData:          formData,
		err:               nil,
		submitted:         false,
		cancelled:         false,
		width:             80,
		height:            24,
		helpFooter:        components.NewHelpFooter("fact_editor", 80),
	}
}

// Init initializes the model
func (m *FactEditorModelNew) Init() tea.Cmd {
	return m.form.Init()
}

// Update handles messages
func (m *FactEditorModelNew) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		// Handle quit
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, func() tea.Msg { return QuitMsg{} }
		}
	}

	// Update the form
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	// Check form state
	if forms.IsCompleted(m.form) {
		return m.handleFormCompletion()
	}

	if forms.IsAborted(m.form) {
		m.cancelled = true
		return m, nil
	}

	return m, cmd
}

// handleFormCompletion processes the completed form and saves the fact.
func (m *FactEditorModelNew) handleFormCompletion() (tea.Model, tea.Cmd) {
	// Apply form data to fact
	err := forms.ApplyFactFormData(m.fact, m.formData)
	if err != nil {
		m.err = fmt.Errorf("failed to apply form data: %w", err)
		return m, nil
	}

	// Update timestamp
	m.fact.UpdatedAt = time.Now()

	// Save the fact using the repository
	factRepo := m.service.GetFactRepository()
	if factRepo == nil {
		m.err = fmt.Errorf("fact repository not available")
		return m, nil
	}

	err = factRepo.Update(m.ctx, m.fact)
	if err != nil {
		m.err = fmt.Errorf("failed to save fact: %w", err)
		return m, nil
	}

	m.submitted = true
	return m, nil
}

// GetFact returns the edited fact
func (m *FactEditorModelNew) GetFact() *career.Fact {
	return m.fact
}

// IsSubmitted returns true if changes were saved
func (m *FactEditorModelNew) IsSubmitted() bool {
	return m.submitted
}

// IsCancelled returns true if operation was cancelled
func (m *FactEditorModelNew) IsCancelled() bool {
	return m.cancelled
}

// Revert reverts changes to the original fact
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
func (m *FactEditorModelNew) GetError() error {
	return m.err
}

// GetTitle returns the modal title for overlay rendering.
func (m *FactEditorModelNew) GetTitle() string {
	return "Fact Editor"
}

// GetContent returns just the form content without header/footer.
// This allows parent intents to compose the modal as an overlay.
func (m *FactEditorModelNew) GetContent() string {
	formView := m.form.View()

	// Add error if present
	if m.err != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(styles.ColorError).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.ColorError).
			Padding(1, 2).
			MarginTop(1)

		formView += "\n\n" + errorStyle.Render(m.err.Error())
	}

	return formView
}

// GetFooter returns the footer instructions for the modal.
func (m *FactEditorModelNew) GetFooter() string {
	return "Enter: Confirm | Esc: Cancel | Tab: Next Field | Shift+Tab: Previous"
}

// View renders the editor UI
func (m *FactEditorModelNew) View() string {
	// Render form using huh
	formView := m.form.View()

	// Add error if present
	if m.err != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F38BA8")). // Catppuccin Red
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#F38BA8")).
			Padding(1, 2).
			MarginTop(1)

		formView += "\n\n" + errorStyle.Render(m.err.Error())
	}

	// Use header and footer components
	headerView := components.NewHeader("Fact Editor", m.width).View()
	footerView := components.NewFooter(m.width).View()

	// Render help footer
	m.helpFooter.SetWidth(m.width)
	helpFooterContent := m.helpFooter.View()

	// Combine all sections
	contentStyle := lipgloss.NewStyle().
		Width(m.width).
		Padding(1, 2)

	fullContent := lipgloss.JoinVertical(
		lipgloss.Left,
		headerView,
		"",
		contentStyle.Render(formView),
		"",
		footerView,
		"",
		helpFooterContent,
	)

	return fullContent
}
