package components

import (
	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/cli/uikit/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// CVConfigWizardModal provides a 3-step wizard for CV generation configuration.
// The wizard guides users through:
//   - Step 1 (WHO): Profile selection + Audience targeting
//   - Step 2 (TECH): Technology focus (conditional - only if techs extracted)
//   - Step 3 (FORMAT): Skills formatting + CV length
//
// Navigation:
//   - ↑/↓ or j/k: Navigate options within Select fields
//   - Enter: Confirm selection and move to next field/step
//   - Esc: Previous step (or cancel if step 1)
//   - Ctrl+S: Skip to generate (use defaults)
//
// Usage:
//
//	modal := components.NewCVConfigWizardModal(width, height)
//	cmd := modal.Init()
//	// In Update:
//	cmd := modal.Update(msg)
//	if modal.IsCompleted() {
//	    config := modal.GetConfigData()
//	    // Use config to generate CV
//	}
type CVConfigWizardModal struct {
	form           *huh.Form
	data           *CVConfigData
	currentStep    int
	visible        bool
	completed      bool
	skipped        bool
	width          int
	height         int
	extractedTechs []ExtractedTechnology
	techsAvailable bool
	profileOptions []ProfileOption
}

// CVConfigData holds the configuration data collected from the wizard.
type CVConfigData struct {
	// Step 1: WHO
	ProfileID string // Required - CV profile to use
	Audience  string // "hiring_manager" | "recruiter" | "peer"

	// Step 2: TECH (optional - only if techs extracted)
	TechFocus    string   // "language_agnostic" | "generalist" | "specialist"
	Technologies []string // MultiSelect (if generalist/specialist)
	FocusArea    string   // "backend" | "frontend" | "fullstack" | "devops"

	// Step 3: FORMAT
	SkillsFormat string // "grouped" | "flat"
	SkillsLimit  int    // max skills per category/total (0 = no limit, default 5)
	CVLength     string // "1_page" | "2_page" | "detailed"
}

// ProfileOption represents a CV profile available for selection.
type ProfileOption struct {
	ID   string
	Name string
}

// ExtractedTechnology represents a technology extracted from career events.
type ExtractedTechnology struct {
	Name       string
	Category   string
	Confidence float64
}

// NewCVConfigWizardModal creates a new CV configuration wizard modal.
// width, height: terminal dimensions for responsive sizing
func NewCVConfigWizardModal(width, height int) *CVConfigWizardModal {
	return NewCVConfigWizardModalWithProfiles(width, height, nil)
}

// NewCVConfigWizardModalWithProfiles creates a wizard modal with profile options.
func NewCVConfigWizardModalWithProfiles(width, height int, profiles []ProfileOption) *CVConfigWizardModal {
	modal := &CVConfigWizardModal{
		data: &CVConfigData{
			// No defaults on init - only set when skipping
		},
		currentStep:    0,
		visible:        true,
		width:          width,
		height:         height,
		profileOptions: profiles,
		techsAvailable: false,
	}

	modal.buildForm()
	return modal
}

// buildForm creates the huh form with 3 steps (groups).
func (m *CVConfigWizardModal) buildForm() {
	// Calculate modal dimensions
	modalWidth := m.width - 20
	if modalWidth > 80 {
		modalWidth = 80
	}
	if modalWidth < 50 {
		modalWidth = 50
	}

	// Convert component types to forms package types
	formProfileOpts := make([]forms.ProfileOption, len(m.profileOptions))
	for i, opt := range m.profileOptions {
		formProfileOpts[i] = forms.ProfileOption{ID: opt.ID, Name: opt.Name}
	}

	formTechs := make([]forms.ExtractedTechnology, len(m.extractedTechs))
	for i, tech := range m.extractedTechs {
		formTechs[i] = forms.ExtractedTechnology{Name: tech.Name}
	}

	// Convert CVConfigData to CVConfigFormData
	formData := &forms.CVConfigFormData{
		ProfileID:    m.data.ProfileID,
		Audience:     m.data.Audience,
		TechFocus:    m.data.TechFocus,
		Technologies: m.data.Technologies,
		FocusArea:    m.data.FocusArea,
		SkillsFormat: m.data.SkillsFormat,
		SkillsLimit:  m.data.SkillsLimit,
		CVLength:     m.data.CVLength,
	}

	// Use forms package helper (applies theme automatically like QuickAddEventModal)
	m.form = forms.NewCVConfigForm(formData, formProfileOpts, formTechs, modalWidth, 0)

	// Sync back to component data (form modifies formData directly)
	m.data = &CVConfigData{
		ProfileID:    formData.ProfileID,
		Audience:     formData.Audience,
		TechFocus:    formData.TechFocus,
		Technologies: formData.Technologies,
		FocusArea:    formData.FocusArea,
		SkillsFormat: formData.SkillsFormat,
		SkillsLimit:  formData.SkillsLimit,
		CVLength:     formData.CVLength,
	}
}

// Init initializes the wizard modal and its form.
func (m *CVConfigWizardModal) Init() tea.Cmd {
	if m.form == nil {
		return nil
	}
	return m.form.Init()
}

// Update handles messages for the wizard modal.
func (m *CVConfigWizardModal) Update(msg tea.Msg) tea.Cmd {
	if !m.visible {
		return nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+s": // Skip shortcut (only ctrl+s, not plain 's')
			if !m.HasRequiredFields() {
				// Don't skip if required fields missing - return early WITHOUT updating form
				return nil
			}
			// Apply defaults for any empty fields
			m.applyDefaults()
			m.completed = true
			m.skipped = true
			m.visible = false
			return nil

		case "esc":
			// If form is nil or at first step, cancel wizard (close modal)
			if m.form == nil || m.currentStep == 0 {
				m.visible = false
				return nil
			}
			// Go back a step in the wizard
			if m.currentStep > 0 {
				m.currentStep--
				// Skip TECH step backwards if techs not available
				if m.currentStep == 1 && !m.techsAvailable {
					m.currentStep = 0
				}
			}
			// Pass Esc to form so it can navigate back between groups
			// Don't return here - let form handle it below
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
		// When form progresses, increment step counter
		if prevState == huh.StateNormal && m.form.State == huh.StateNormal {
			// Form might have advanced to next group on Enter
			// We detect this by checking if we're getting key messages
			if keyMsg, ok := msg.(tea.KeyMsg); ok {
				if keyMsg.String() == "enter" && m.currentStep < 2 {
					m.currentStep++

					// Skip TECH step if no techs available
					if m.currentStep == 1 && !m.techsAvailable {
						m.currentStep = 2
					}
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
func (m *CVConfigWizardModal) View() string {
	if !m.visible {
		return ""
	}

	if m.form == nil {
		return ""
	}

	// Calculate modal dimensions for border
	modalWidth := m.width - 20
	if modalWidth > 80 {
		modalWidth = 80
	}
	if modalWidth < 50 {
		modalWidth = 50
	}

	modalHeight := m.height - 10
	if modalHeight < 20 {
		modalHeight = 20
	}

	// Use default theme for colors
	th := theme.Default()

	// Add main title
	titleStyle := lipgloss.NewStyle().
		Foreground(th.AccentColor()).
		Bold(true).
		Align(lipgloss.Center).
		Width(modalWidth - 4)

	title := titleStyle.Render("CV Configuration")

	// Render form
	formView := m.form.View()

	// Create footer with keyboard shortcuts
	footer := m.buildFooter()

	// Combine title, form and footer
	content := lipgloss.JoinVertical(lipgloss.Left, title, "", formView, "", footer)

	// Wrap in styled container with solid background
	styledContent := lipgloss.NewStyle().
		Width(modalWidth).
		MaxHeight(modalHeight).
		Background(th.BackgroundColor()). // CRITICAL: Solid background prevents transparency
		Border(lipgloss.RoundedBorder()).
		BorderForeground(th.AccentColor()).
		Padding(1).
		Render(content)

	return styledContent
}

// buildFooter creates the keyboard shortcuts footer using UIKit primitives.
func (m *CVConfigWizardModal) buildFooter() string {
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
func (m *CVConfigWizardModal) Show() {
	m.visible = true
}

// Reset resets the wizard state while preserving the entered data.
// This is used when navigating back to the wizard from a later state
// (e.g., from CV preview). The form is rebuilt with the existing data
// so users see their previous selections.
func (m *CVConfigWizardModal) Reset() {
	// Reset state flags but keep data
	m.completed = false
	m.skipped = false
	m.currentStep = 0
	m.visible = true

	// Rebuild form with existing data - this creates a fresh form
	// but populates it with the previously entered values
	m.buildForm()
}

// Hide makes the modal invisible.
func (m *CVConfigWizardModal) Hide() {
	m.visible = false
}

// IsVisible returns whether the modal is currently visible.
func (m *CVConfigWizardModal) IsVisible() bool {
	return m.visible
}

// IsCompleted returns whether the wizard has been completed.
func (m *CVConfigWizardModal) IsCompleted() bool {
	return m.completed
}

// IsSkipped returns whether the wizard was skipped (Ctrl+S).
func (m *CVConfigWizardModal) IsSkipped() bool {
	return m.skipped
}

// GetConfigData returns the collected configuration data.
func (m *CVConfigWizardModal) GetConfigData() *CVConfigData {
	return m.data
}

// GetCurrentStep returns the current step index (0-based).
// Note: This tracks the form's internal step progression.
func (m *CVConfigWizardModal) GetCurrentStep() int {
	if m.form == nil {
		return 0
	}

	// huh.Form doesn't expose current group index directly,
	// so we track it manually based on Update() calls
	return m.currentStep
}

// GetStepCount returns the total number of steps.
func (m *CVConfigWizardModal) GetStepCount() int {
	return 3 // WHO, TECH, FORMAT
}

// AreTechsAvailable returns whether extracted technologies are available.
func (m *CVConfigWizardModal) AreTechsAvailable() bool {
	return m.techsAvailable
}

// GetExtractedTechnologies returns the extracted technologies.
func (m *CVConfigWizardModal) GetExtractedTechnologies() []ExtractedTechnology {
	return m.extractedTechs
}

// SetExtractedTechnologies sets the extracted technologies and rebuilds the form.
func (m *CVConfigWizardModal) SetExtractedTechnologies(techs []ExtractedTechnology) {
	m.extractedTechs = techs
	m.techsAvailable = len(techs) > 0
	m.buildForm()
}

// GetProfileOptions returns the available profile options.
func (m *CVConfigWizardModal) GetProfileOptions() []ProfileOption {
	return m.profileOptions
}

// SetProfileID sets the profile ID.
func (m *CVConfigWizardModal) SetProfileID(id string) {
	m.data.ProfileID = id
}

// SetAudience sets the target audience.
func (m *CVConfigWizardModal) SetAudience(audience string) {
	m.data.Audience = audience
}

// SetTechFocus sets the technology focus.
func (m *CVConfigWizardModal) SetTechFocus(focus string) {
	m.data.TechFocus = focus
}

// SetTechnologies sets the selected technologies.
func (m *CVConfigWizardModal) SetTechnologies(techs []string) {
	m.data.Technologies = techs
}

// SetFocusArea sets the focus area.
func (m *CVConfigWizardModal) SetFocusArea(area string) {
	m.data.FocusArea = area
}

// SetSkillsFormat sets the skills format.
func (m *CVConfigWizardModal) SetSkillsFormat(format string) {
	m.data.SkillsFormat = format
}

// SetSkillsLimit sets the skills limit per category/total.
func (m *CVConfigWizardModal) SetSkillsLimit(limit int) {
	m.data.SkillsLimit = limit
}

// SetCVLength sets the CV length.
func (m *CVConfigWizardModal) SetCVLength(length string) {
	m.data.CVLength = length
}

// Complete marks the wizard as completed.
func (m *CVConfigWizardModal) Complete() {
	if m.HasRequiredFields() {
		m.completed = true
		m.visible = false
	}
}

// HasRequiredFields checks if all required fields are filled.
func (m *CVConfigWizardModal) HasRequiredFields() bool {
	return m.data.ProfileID != ""
}

// GetDimensions returns the current modal dimensions.
func (m *CVConfigWizardModal) GetDimensions() (width, height int) {
	return m.width, m.height
}

// applyDefaults sets default values for any empty configuration fields.
// This is called when the user skips the wizard.
func (m *CVConfigWizardModal) applyDefaults() {
	if m.data.Audience == "" {
		m.data.Audience = "hiring_manager"
	}
	if m.data.TechFocus == "" {
		m.data.TechFocus = "language_agnostic"
	}
	if m.data.SkillsFormat == "" {
		m.data.SkillsFormat = "grouped"
	}
	if m.data.CVLength == "" {
		m.data.CVLength = "2_page"
	}
}
