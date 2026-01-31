package components

import (
	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/cli/uikit/containers"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/cli/uikit/theme"
	tea "github.com/charmbracelet/bubbletea"
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
	wizard   *behaviors.WizardBehavior[CVConfigData]
	form     *forms.WizardFormAdapter
	formData *forms.CVConfigFormData
	data     *CVConfigData

	currentStep    int
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
	Technologies []string // MultiSelect (if generalist mode)
	Technology   string   // Single-select (if specialist mode)
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

// NewCVConfigWizardModal creates a new CV configuration wizard modal
// with the given terminal dimensions for responsive sizing.
func NewCVConfigWizardModal(width, height int) *CVConfigWizardModal {
	return NewCVConfigWizardModalWithProfiles(width, height, nil)
}

// NewCVConfigWizardModalWithProfiles creates a wizard modal with profile options.
func NewCVConfigWizardModalWithProfiles(width, height int, profiles []ProfileOption) *CVConfigWizardModal {
	modal := &CVConfigWizardModal{
		data:           &CVConfigData{},
		currentStep:    0,
		width:          width,
		height:         height,
		profileOptions: profiles,
		techsAvailable: false,
	}

	modal.buildForm()

	wizard := behaviors.NewWizardBehavior[CVConfigData](modal.form, modal.data)
	modal.wizard = wizard

	return modal
}

// buildForm creates the huh form with 3 steps (groups).
func (m *CVConfigWizardModal) buildForm() {
	m.syncFromFormData()

	modalWidth := calcCVModalWidth(m.width)

	formProfileOpts := make([]forms.ProfileOption, len(m.profileOptions))
	for i, opt := range m.profileOptions {
		formProfileOpts[i] = forms.ProfileOption{ID: opt.ID, Name: opt.Name}
	}

	formTechs := make([]forms.ExtractedTechnology, len(m.extractedTechs))
	for i, tech := range m.extractedTechs {
		formTechs[i] = forms.ExtractedTechnology{Name: tech.Name}
	}

	singleTechSelect := m.data.TechFocus == "specialist"

	m.formData = &forms.CVConfigFormData{
		ProfileID:    m.data.ProfileID,
		Audience:     m.data.Audience,
		TechFocus:    m.data.TechFocus,
		Technologies: m.data.Technologies,
		Technology:   m.data.Technology,
		FocusArea:    m.data.FocusArea,
		SkillsFormat: m.data.SkillsFormat,
		SkillsLimit:  m.data.SkillsLimit,
		CVLength:     m.data.CVLength,
	}

	huhForm := forms.NewCVConfigForm(m.formData, formProfileOpts, formTechs, modalWidth, 0, singleTechSelect)
	m.form = forms.NewWizardFormAdapter(huhForm, 3)
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
	if !m.wizard.IsVisible() {
		return nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+s":
			if !m.HasRequiredFields() {
				return nil
			}
			m.applyDefaults()
			m.wizard.Skip()
			return nil

		case "esc":
			if m.form == nil || m.currentStep == 0 {
				m.wizard.Hide()
				return nil
			}
			if m.currentStep > 0 {
				m.currentStep--
				if m.currentStep == 1 && !m.techsAvailable {
					m.currentStep = 0
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.buildForm()
		if m.wizard != nil {
			m.wizard.SetForm(m.form)
		}
		return m.form.Init()
	}

	if m.form != nil {
		wasCompleted := m.form.IsCompleted()

		cmd := m.form.Update(msg)

		if !wasCompleted && !m.form.IsCompleted() {
			if keyMsg, ok := msg.(tea.KeyMsg); ok {
				if keyMsg.String() == "enter" && m.currentStep < 2 {
					m.currentStep++

					if m.currentStep == 1 && !m.techsAvailable {
						m.currentStep = 2
					}
				}
			}
		}

		if m.form.IsCompleted() {
			m.wizard.Complete()
		}

		return cmd
	}

	return nil
}

// View renders the wizard modal.
func (m *CVConfigWizardModal) View() string {
	if !m.wizard.IsVisible() {
		return ""
	}

	if m.form == nil {
		return ""
	}

	modalWidth := calcCVModalWidth(m.width)

	modalHeight := m.height - 10
	if modalHeight < 20 {
		modalHeight = 20
	}

	th := theme.Default()

	title := primitives.Title("CV Configuration", th).
		Width(modalWidth - 4).
		Center().
		Render()

	formView := m.form.View()
	footer := m.buildFooter()
	content := lipgloss.JoinVertical(lipgloss.Left, title, "", formView, "", footer)

	return containers.NewBox(th).
		Content(content).
		Width(modalWidth).
		MaxHeight(modalHeight).
		Padding(1).
		Background(th.BackgroundColor()).
		Variant(containers.BoxInfo).
		Render()
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
	m.wizard.Show()
}

// Reset resets the wizard state while preserving the entered data.
func (m *CVConfigWizardModal) Reset() {
	m.wizard.Reset()
	m.currentStep = 0

	m.buildForm()
	m.wizard.SetForm(m.form)
}

// Hide makes the modal invisible.
func (m *CVConfigWizardModal) Hide() {
	m.wizard.Hide()
}

// IsVisible returns whether the modal is currently visible.
func (m *CVConfigWizardModal) IsVisible() bool {
	return m.wizard.IsVisible()
}

// IsCompleted returns whether the wizard has been completed.
func (m *CVConfigWizardModal) IsCompleted() bool {
	return m.wizard.IsCompleted()
}

// IsSkipped returns whether the wizard was skipped (Ctrl+S).
func (m *CVConfigWizardModal) IsSkipped() bool {
	return m.wizard.IsSkipped()
}

// GetConfigData returns the collected configuration data.
func (m *CVConfigWizardModal) GetConfigData() *CVConfigData {
	m.syncFromFormData()
	return m.data
}

// GetCurrentStep returns the current step index (0-based).
func (m *CVConfigWizardModal) GetCurrentStep() int {
	if m.form == nil {
		return 0
	}
	return m.currentStep
}

// GetStepCount returns the total number of steps.
func (m *CVConfigWizardModal) GetStepCount() int {
	return 3
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
	if m.wizard != nil {
		m.wizard.SetForm(m.form)
	}
}

// GetProfileOptions returns the available profile options.
func (m *CVConfigWizardModal) GetProfileOptions() []ProfileOption {
	return m.profileOptions
}

// SetProfileID sets the profile ID.
func (m *CVConfigWizardModal) SetProfileID(id string) {
	if m.formData != nil {
		m.formData.ProfileID = id
	}
	m.data.ProfileID = id
}

// SetAudience sets the target audience.
func (m *CVConfigWizardModal) SetAudience(audience string) {
	if m.formData != nil {
		m.formData.Audience = audience
	}
	m.data.Audience = audience
}

// SetTechFocus sets the technology focus.
func (m *CVConfigWizardModal) SetTechFocus(focus string) {
	if m.formData != nil {
		m.formData.TechFocus = focus
	}
	m.data.TechFocus = focus
}

// SetTechnologies sets the selected technologies (for generalist mode).
func (m *CVConfigWizardModal) SetTechnologies(techs []string) {
	if m.formData != nil {
		m.formData.Technologies = techs
	}
	m.data.Technologies = techs
}

// SetTechnology sets the selected technology (for specialist mode).
func (m *CVConfigWizardModal) SetTechnology(tech string) {
	if m.formData != nil {
		m.formData.Technology = tech
	}
	m.data.Technology = tech
}

// SetFocusArea sets the focus area.
func (m *CVConfigWizardModal) SetFocusArea(area string) {
	if m.formData != nil {
		m.formData.FocusArea = area
	}
	m.data.FocusArea = area
}

// SetSkillsFormat sets the skills format.
func (m *CVConfigWizardModal) SetSkillsFormat(format string) {
	if m.formData != nil {
		m.formData.SkillsFormat = format
	}
	m.data.SkillsFormat = format
}

// SetSkillsLimit sets the skills limit per category/total.
func (m *CVConfigWizardModal) SetSkillsLimit(limit int) {
	if m.formData != nil {
		m.formData.SkillsLimit = limit
	}
	m.data.SkillsLimit = limit
}

// SetCVLength sets the CV length.
func (m *CVConfigWizardModal) SetCVLength(length string) {
	if m.formData != nil {
		m.formData.CVLength = length
	}
	m.data.CVLength = length
}

// Complete marks the wizard as completed.
func (m *CVConfigWizardModal) Complete() {
	if m.HasRequiredFields() {
		m.wizard.Complete()
	}
}

// HasRequiredFields checks if all required fields are filled.
func (m *CVConfigWizardModal) HasRequiredFields() bool {
	m.syncFromFormData()
	return m.data.ProfileID != ""
}

// GetDimensions returns the current modal dimensions.
func (m *CVConfigWizardModal) GetDimensions() (width, height int) {
	return m.width, m.height
}

func (m *CVConfigWizardModal) syncFromFormData() {
	if m.formData == nil {
		return
	}
	m.data.ProfileID = m.formData.ProfileID
	m.data.Audience = m.formData.Audience
	m.data.TechFocus = m.formData.TechFocus
	m.data.Technologies = m.formData.Technologies
	m.data.Technology = m.formData.Technology
	m.data.FocusArea = m.formData.FocusArea
	m.data.SkillsFormat = m.formData.SkillsFormat
	m.data.SkillsLimit = m.formData.SkillsLimit
	m.data.CVLength = m.formData.CVLength
}

func (m *CVConfigWizardModal) applyDefaults() {
	m.syncFromFormData()

	if m.data.Audience == "" {
		m.data.Audience = "hiring_manager"
		if m.formData != nil {
			m.formData.Audience = "hiring_manager"
		}
	}
	if m.data.TechFocus == "" {
		m.data.TechFocus = "language_agnostic"
		if m.formData != nil {
			m.formData.TechFocus = "language_agnostic"
		}
	}
	if m.data.SkillsFormat == "" {
		m.data.SkillsFormat = "grouped"
		if m.formData != nil {
			m.formData.SkillsFormat = "grouped"
		}
	}
	if m.data.CVLength == "" {
		m.data.CVLength = "2_page"
		if m.formData != nil {
			m.formData.CVLength = "2_page"
		}
	}
}

func calcCVModalWidth(terminalWidth int) int {
	modalWidth := terminalWidth - 20
	if modalWidth > 80 {
		modalWidth = 80
	}
	if modalWidth < 50 {
		modalWidth = 50
	}
	return modalWidth
}
