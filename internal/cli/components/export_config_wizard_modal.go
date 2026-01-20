package components

import (
	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/cli/types"
	"github.com/baphled/kariya/internal/cli/uikit/containers"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/cli/uikit/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// ExportConfigWizardModal provides a 2-step wizard for export configuration.
// The wizard guides users through:
//   - Step 1 (WHAT): Artifact type selection (events, facts, bursts, cv, profile)
//   - Step 2 (HOW): Format and destination selection
//
// Navigation:
//   - ↑/↓ or j/k: Navigate options within Select fields
//   - Enter: Confirm selection and move to next field/step
//   - Esc: Previous step (or cancel if step 1)
//   - Ctrl+S: Skip to export (use defaults)
//
// Usage:
//
//	modal := components.NewExportConfigWizardModal(width, height)
//	cmd := modal.Init()
//	// In Update:
//	cmd := modal.Update(msg)
//	if modal.IsCompleted() {
//	    config := modal.GetConfigData()
//	    // Use config to perform export
//	}
type ExportConfigWizardModal struct {
	form        *huh.Form
	data        *forms.ExportConfigFormData
	currentStep int
	visible     bool
	completed   bool
	cancelled   bool
	skipped     bool
	width       int
	height      int
}

// NewExportConfigWizardModal creates a new export configuration wizard modal.
// width, height: terminal dimensions for responsive sizing
func NewExportConfigWizardModal(width, height int) *ExportConfigWizardModal {
	modal := &ExportConfigWizardModal{
		data:        forms.NewExportConfigFormData(),
		currentStep: 0,
		visible:     true,
		width:       width,
		height:      height,
	}

	modal.buildForm()
	return modal
}

// buildForm creates the huh form with 2 steps (groups).
func (m *ExportConfigWizardModal) buildForm() {
	// Calculate modal dimensions
	modalWidth := m.width - 20
	if modalWidth > 70 {
		modalWidth = 70
	}
	if modalWidth < 50 {
		modalWidth = 50
	}

	// Use forms package to create the form
	m.form = forms.NewExportConfigForm(m.data, modalWidth, 0)
}

// Init initializes the wizard modal and its form.
func (m *ExportConfigWizardModal) Init() tea.Cmd {
	if m.form == nil {
		return nil
	}
	return m.form.Init()
}

// Update handles messages for the wizard modal.
func (m *ExportConfigWizardModal) Update(msg tea.Msg) tea.Cmd {
	if !m.visible {
		return nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+s": // Skip shortcut - export with defaults
			m.applyDefaults()
			m.completed = true
			m.skipped = true
			m.visible = false
			return nil

		case "esc":
			// If at first step, cancel wizard
			if m.currentStep == 0 {
				m.cancelled = true
				m.visible = false
				return nil
			}
			// Go back a step in the wizard
			if m.currentStep > 0 {
				m.currentStep--
				m.form.PrevGroup() // Navigate huh form to previous group
				return nil         // Don't let Esc propagate (would cancel form)
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Rebuild form with new dimensions
		m.buildForm()
		return m.form.Init()
	}

	// Update form
	if m.form != nil {
		// Track form state before update
		prevState := m.form.State

		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
		}

		// Track step progression based on state changes
		if prevState == huh.StateNormal && m.form.State == huh.StateNormal {
			if keyMsg, ok := msg.(tea.KeyMsg); ok {
				if keyMsg.String() == "enter" && m.currentStep < 1 {
					m.currentStep++
				}
			}
		}

		// Check if form completed
		if m.form.State == huh.StateCompleted {
			m.completed = true
			m.visible = false
		}

		return cmd
	}

	return nil
}

// View renders the wizard modal.
func (m *ExportConfigWizardModal) View() string {
	if !m.visible {
		return ""
	}

	if m.form == nil {
		return ""
	}

	// Calculate modal dimensions
	modalWidth := m.width - 20
	if modalWidth > 70 {
		modalWidth = 70
	}
	if modalWidth < 50 {
		modalWidth = 50
	}

	modalHeight := m.height - 10
	if modalHeight < 15 {
		modalHeight = 15
	}

	// Use default theme for colors
	th := theme.Default()

	// Add main title using UIKit Text with centering
	title := primitives.Title("Export Configuration", th).
		Width(modalWidth - 4).
		Center().
		Render()

	// Render form
	formView := m.form.View()

	// Create footer with keyboard shortcuts
	footer := m.buildFooter()

	// Combine title, form and footer
	content := lipgloss.JoinVertical(lipgloss.Left, title, "", formView, "", footer)

	// Wrap in styled container with solid background using UIKit Box
	return containers.NewBox(th).
		Content(content).
		Width(modalWidth).
		MaxHeight(modalHeight).
		Padding(1).
		Background(th.BackgroundColor()).
		Variant(containers.BoxInfo). // Use accent color border
		Render()
}

// buildFooter creates the keyboard shortcuts footer using UIKit primitives.
func (m *ExportConfigWizardModal) buildFooter() string {
	th := theme.Default()

	badges := []*primitives.Badge{
		primitives.NavigateBadge(th),
		primitives.SelectBadge(th),
	}

	if m.currentStep > 0 {
		badges = append(badges, primitives.BackBadge(th))
	} else {
		badges = append(badges, primitives.CancelBadge(th))
	}

	badges = append(badges, primitives.SkipBadge(th))

	return primitives.RenderHelpFooter(th, badges...)
}

// Show makes the modal visible.
func (m *ExportConfigWizardModal) Show() {
	m.visible = true
}

// Hide makes the modal invisible.
func (m *ExportConfigWizardModal) Hide() {
	m.visible = false
}

// Reset resets the wizard state while preserving the entered data.
func (m *ExportConfigWizardModal) Reset() {
	m.completed = false
	m.cancelled = false
	m.skipped = false
	m.currentStep = 0
	m.visible = true
	m.buildForm()
}

// IsVisible returns whether the modal is currently visible.
func (m *ExportConfigWizardModal) IsVisible() bool {
	return m.visible
}

// IsCompleted returns whether the wizard has been completed.
func (m *ExportConfigWizardModal) IsCompleted() bool {
	return m.completed
}

// IsCancelled returns whether the wizard was cancelled.
func (m *ExportConfigWizardModal) IsCancelled() bool {
	return m.cancelled
}

// IsSkipped returns whether the wizard was skipped (Ctrl+S).
func (m *ExportConfigWizardModal) IsSkipped() bool {
	return m.skipped
}

// GetConfigData returns the collected configuration data.
func (m *ExportConfigWizardModal) GetConfigData() *forms.ExportConfigFormData {
	return m.data
}

// GetArtifactType returns the selected artifact type.
func (m *ExportConfigWizardModal) GetArtifactType() types.ExportArtifactType {
	return types.ExportArtifactType(m.data.ArtifactType)
}

// GetFormat returns the selected export format.
func (m *ExportConfigWizardModal) GetFormat() types.ExportFormat {
	return types.ExportFormat(m.data.Format)
}

// GetDestination returns the selected export destination.
func (m *ExportConfigWizardModal) GetDestination() types.ExportDestination {
	return types.ExportDestination(m.data.Destination)
}

// GetCurrentStep returns the current step index (0-based).
func (m *ExportConfigWizardModal) GetCurrentStep() int {
	return m.currentStep
}

// GetStepCount returns the total number of steps.
func (m *ExportConfigWizardModal) GetStepCount() int {
	return 2 // WHAT, HOW
}

// GoToStep sets the current step (0-based). Useful when going back from preview.
func (m *ExportConfigWizardModal) GoToStep(step int) {
	if step >= 0 && step < m.GetStepCount() {
		m.currentStep = step
		m.buildForm()
		// Navigate huh form to the correct group (form starts at group 0)
		for i := 0; i < step; i++ {
			m.form.NextGroup()
		}
	}
}

// GoToLastStep goes to the last step of the wizard.
func (m *ExportConfigWizardModal) GoToLastStep() {
	m.GoToStep(m.GetStepCount() - 1)
}

// ResumeAtLastStep resumes the wizard at the last step, preserving all data.
// Use this when returning from preview to allow incremental back-navigation.
func (m *ExportConfigWizardModal) ResumeAtLastStep() {
	m.completed = false
	m.cancelled = false
	m.skipped = false
	m.visible = true
	m.GoToLastStep()
}

// SetArtifactType sets the artifact type.
func (m *ExportConfigWizardModal) SetArtifactType(artifactType string) {
	m.data.ArtifactType = artifactType
}

// SetFormat sets the export format.
func (m *ExportConfigWizardModal) SetFormat(format string) {
	m.data.Format = format
}

// SetDestination sets the export destination.
func (m *ExportConfigWizardModal) SetDestination(destination string) {
	m.data.Destination = destination
}

// GetDimensions returns the current modal dimensions.
func (m *ExportConfigWizardModal) GetDimensions() (width, height int) {
	return m.width, m.height
}

// applyDefaults sets default values for any empty configuration fields.
// This is called when the user skips the wizard.
func (m *ExportConfigWizardModal) applyDefaults() {
	if m.data.ArtifactType == "" {
		m.data.ArtifactType = "events"
	}
	if m.data.Format == "" {
		m.data.Format = "json"
	}
	if m.data.Destination == "" {
		m.data.Destination = "file"
	}
}
