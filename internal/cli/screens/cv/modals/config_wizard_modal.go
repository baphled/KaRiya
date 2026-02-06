package modals

import (
	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/cli/uikit/containers"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/cli/uikit/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ConfigWizardModal provides a 3-step wizard for CV generation configuration.
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
//	modal := modals.NewConfigWizardModal(width, height)
//	cmd := modal.Init()
//	// In Update:
//	cmd := modal.Update(msg)
//	if modal.IsCompleted() {
//	    config := modal.GetConfigData()
//	    // Use config to generate CV
//	}
type ConfigWizardModal struct {
	wizard   *behaviors.WizardBehavior[ConfigData]
	form     *forms.WizardFormAdapter
	formData *forms.CVConfigFormData
	data     *ConfigData

	width          int
	height         int
	extractedTechs []ExtractedTechnology
	techsAvailable bool
	profileOptions []ProfileOption
}

// ConfigData holds the configuration data collected from the wizard.
type ConfigData struct {
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

// NewConfigWizardModal creates a new CV configuration wizard modal.
//
// Expected:
//   - width and height must be valid positive integers.
//
// Returns:
//   - A fully initialized ConfigWizardModal ready for use.
//
// Side effects:
//   - None.
func NewConfigWizardModal(width, height int) *ConfigWizardModal {
	return NewConfigWizardModalWithProfiles(width, height, nil)
}

// NewConfigWizardModalWithProfiles creates a wizard modal with profile options.
//
// Expected:
//   - width and height must be valid positive integers.
//   - profiles can be nil or a slice of ProfileOption.
//
// Returns:
//   - A fully initialized ConfigWizardModal ready for use.
//
// Side effects:
//   - None.
func NewConfigWizardModalWithProfiles(width, height int, profiles []ProfileOption) *ConfigWizardModal {
	modal := &ConfigWizardModal{
		data:           &ConfigData{},
		width:          width,
		height:         height,
		profileOptions: profiles,
		techsAvailable: false,
	}

	modal.buildForm()

	wizard := behaviors.NewWizardBehavior[ConfigData](modal.form, modal.data)
	modal.wizard = wizard

	return modal
}

// buildForm creates the huh form with 3 steps (groups).
func (m *ConfigWizardModal) buildForm() {
	m.syncFromFormData()

	modalWidth := calcModalWidth(m.width)

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
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (m *ConfigWizardModal) Init() tea.Cmd {
	if m.form == nil {
		return nil
	}
	return m.form.Init()
}

// Update handles messages for the wizard modal.
//
// Expected:
//   - msg must be a valid tea.Msg.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - May update internal state based on key messages.
func (m *ConfigWizardModal) Update(msg tea.Msg) tea.Cmd {
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
			if m.form == nil || m.form.CurrentStep() == 0 {
				m.wizard.Hide()
				return nil
			}
			newStep := m.form.CurrentStep() - 1
			if newStep == 1 && !m.techsAvailable {
				newStep = 0
			}
			m.form.SetCurrentStep(newStep)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.form.SetDimensions(calcModalWidth(m.width), m.height)
		return m.form.Init()
	}

	wasCompleted := m.form.IsCompleted()

	cmd := m.wizard.Update(msg)

	if !wasCompleted && !m.form.IsCompleted() {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			if keyMsg.String() == "enter" && m.form.CurrentStep() < 2 {
				newStep := m.form.CurrentStep() + 1
				if newStep == 1 && !m.techsAvailable {
					newStep = 2
				}
				m.form.SetCurrentStep(newStep)
			}
		}
	}

	return cmd
}

// View renders the wizard modal.
//
// Returns:
//   - A string containing the rendered modal view.
//
// Side effects:
//   - None.
func (m *ConfigWizardModal) View() string {
	if !m.wizard.IsVisible() {
		return ""
	}

	if m.form == nil {
		return ""
	}

	modalWidth := calcModalWidth(m.width)

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
func (m *ConfigWizardModal) buildFooter() string {
	th := theme.Default()

	badges := []*primitives.Badge{
		primitives.NavigateBadge(th),
		primitives.SelectBadge(th),
	}

	if m.form.CurrentStep() > 0 {
		badges = append(badges, primitives.BackBadge(th))
	} else {
		badges = append(badges, primitives.CancelBadge(th))
	}

	badges = append(badges, primitives.SkipBadge(th))

	return primitives.RenderHelpFooter(th, badges...)
}

// Show makes the modal visible.
//
// Side effects:
//   - Updates modal visibility state.
func (m *ConfigWizardModal) Show() {
	m.wizard.Show()
}

// Reset resets the wizard state while preserving the entered data.
//
// Side effects:
//   - Resets wizard step to beginning.
func (m *ConfigWizardModal) Reset() {
	m.wizard.Reset()

	m.buildForm()
	m.wizard.SetForm(m.form)
}

// Hide makes the modal invisible.
//
// Side effects:
//   - Updates modal visibility state.
func (m *ConfigWizardModal) Hide() {
	m.wizard.Hide()
}

// IsVisible returns whether the modal is currently visible.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *ConfigWizardModal) IsVisible() bool {
	return m.wizard.IsVisible()
}

// IsCompleted returns whether the wizard has been completed.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *ConfigWizardModal) IsCompleted() bool {
	return m.wizard.IsCompleted()
}

// IsSkipped returns whether the wizard was skipped (Ctrl+S).
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *ConfigWizardModal) IsSkipped() bool {
	return m.wizard.IsSkipped()
}

// GetConfigData returns the collected configuration data.
//
// Returns:
//   - A fully initialized ConfigData ready for use.
//
// Side effects:
//   - Syncs data from form before returning.
func (m *ConfigWizardModal) GetConfigData() *ConfigData {
	m.syncFromFormData()
	return m.data
}

// GetCurrentStep returns the current step index (0-based).
//
// Returns:
//   - An int value representing the current step.
//
// Side effects:
//   - None.
func (m *ConfigWizardModal) GetCurrentStep() int {
	return m.wizard.CurrentStep()
}

// GetStepCount returns the total number of steps.
//
// Returns:
//   - An int value (always 3).
//
// Side effects:
//   - None.
func (m *ConfigWizardModal) GetStepCount() int {
	return 3
}

// AreTechsAvailable returns whether extracted technologies are available.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *ConfigWizardModal) AreTechsAvailable() bool {
	return m.techsAvailable
}

// GetExtractedTechnologies returns the extracted technologies.
//
// Returns:
//   - A slice of ExtractedTechnology.
//
// Side effects:
//   - None.
func (m *ConfigWizardModal) GetExtractedTechnologies() []ExtractedTechnology {
	return m.extractedTechs
}

// SetExtractedTechnologies sets the extracted technologies and rebuilds the form.
//
// Expected:
//   - techs can be nil or a slice of ExtractedTechnology.
//
// Side effects:
//   - Rebuilds the form with new technology options.
func (m *ConfigWizardModal) SetExtractedTechnologies(techs []ExtractedTechnology) {
	m.extractedTechs = techs
	m.techsAvailable = len(techs) > 0
	m.buildForm()
	if m.wizard != nil {
		m.wizard.SetForm(m.form)
	}
}

// GetProfileOptions returns the available profile options.
//
// Returns:
//   - A slice of ProfileOption.
//
// Side effects:
//   - None.
func (m *ConfigWizardModal) GetProfileOptions() []ProfileOption {
	return m.profileOptions
}

// SetProfileID sets the profile ID.
//
// Expected:
//   - id must be a valid string.
//
// Side effects:
//   - Updates profile ID in both data and formData.
func (m *ConfigWizardModal) SetProfileID(id string) {
	if m.formData != nil {
		m.formData.ProfileID = id
	}
	m.data.ProfileID = id
}

// SetAudience sets the target audience.
//
// Expected:
//   - audience must be a valid string.
//
// Side effects:
//   - Updates audience in both data and formData.
func (m *ConfigWizardModal) SetAudience(audience string) {
	if m.formData != nil {
		m.formData.Audience = audience
	}
	m.data.Audience = audience
}

// SetTechFocus sets the technology focus.
//
// Expected:
//   - focus must be a valid string.
//
// Side effects:
//   - Updates tech focus in both data and formData.
func (m *ConfigWizardModal) SetTechFocus(focus string) {
	if m.formData != nil {
		m.formData.TechFocus = focus
	}
	m.data.TechFocus = focus
}

// SetTechnologies sets the selected technologies (for generalist mode).
//
// Expected:
//   - techs must be a valid string slice.
//
// Side effects:
//   - Updates technologies in both data and formData.
func (m *ConfigWizardModal) SetTechnologies(techs []string) {
	if m.formData != nil {
		m.formData.Technologies = techs
	}
	m.data.Technologies = techs
}

// SetTechnology sets the selected technology (for specialist mode).
//
// Expected:
//   - tech must be a valid string.
//
// Side effects:
//   - Updates technology in both data and formData.
func (m *ConfigWizardModal) SetTechnology(tech string) {
	if m.formData != nil {
		m.formData.Technology = tech
	}
	m.data.Technology = tech
}

// SetFocusArea sets the focus area.
//
// Expected:
//   - area must be a valid string.
//
// Side effects:
//   - Updates focus area in both data and formData.
func (m *ConfigWizardModal) SetFocusArea(area string) {
	if m.formData != nil {
		m.formData.FocusArea = area
	}
	m.data.FocusArea = area
}

// SetSkillsFormat sets the skills format.
//
// Expected:
//   - format must be a valid string.
//
// Side effects:
//   - Updates skills format in both data and formData.
func (m *ConfigWizardModal) SetSkillsFormat(format string) {
	if m.formData != nil {
		m.formData.SkillsFormat = format
	}
	m.data.SkillsFormat = format
}

// SetSkillsLimit sets the skills limit per category/total.
//
// Expected:
//   - limit must be a valid non-negative integer.
//
// Side effects:
//   - Updates skills limit in both data and formData.
func (m *ConfigWizardModal) SetSkillsLimit(limit int) {
	if m.formData != nil {
		m.formData.SkillsLimit = limit
	}
	m.data.SkillsLimit = limit
}

// SetCVLength sets the CV length.
//
// Expected:
//   - length must be a valid string.
//
// Side effects:
//   - Updates CV length in both data and formData.
func (m *ConfigWizardModal) SetCVLength(length string) {
	if m.formData != nil {
		m.formData.CVLength = length
	}
	m.data.CVLength = length
}

// Complete marks the wizard as completed.
//
// Side effects:
//   - Marks wizard as completed if required fields are present.
func (m *ConfigWizardModal) Complete() {
	if m.HasRequiredFields() {
		m.wizard.Complete()
	}
}

// HasRequiredFields checks if all required fields are filled.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - Syncs data from form before checking.
func (m *ConfigWizardModal) HasRequiredFields() bool {
	m.syncFromFormData()
	return m.data.ProfileID != ""
}

// GetDimensions returns the current modal dimensions.
//
// Returns:
//   - width and height as int values.
//
// Side effects:
//   - None.
func (m *ConfigWizardModal) GetDimensions() (width, height int) {
	return m.width, m.height
}

func (m *ConfigWizardModal) syncFromFormData() {
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

func (m *ConfigWizardModal) applyDefaults() {
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

// calcModalWidth calculates the appropriate modal width based on terminal width.
func calcModalWidth(terminalWidth int) int {
	modalWidth := terminalWidth - 20
	if modalWidth > 80 {
		modalWidth = 80
	}
	if modalWidth < 50 {
		modalWidth = 50
	}
	return modalWidth
}
