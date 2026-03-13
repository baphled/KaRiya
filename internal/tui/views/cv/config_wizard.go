package cv

import (
	"github.com/baphled/kariya/internal/tui/forms"
	"github.com/baphled/kariya/internal/ui/behaviors"
	"github.com/baphled/kariya/internal/ui/uikit/containers"
	"github.com/baphled/kariya/internal/ui/uikit/primitives"
	"github.com/baphled/kariya/internal/ui/uikit/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ConfigWizard provides a 3-step wizard for CV generation configuration.
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
//	modal := NewConfigWizard(width, height)
//	cmd := modal.Init()
//	// In Update:
//	cmd := modal.Update(msg)
//	if modal.IsCompleted() {
//	    config := modal.GetConfigData()
//	    // Use config to generate CV
//	}
type ConfigWizard struct {
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

// NewConfigWizard creates a new CV configuration wizard modal.
//
// Expected:
//   - width and height must be valid positive integers.
//
// Returns:
//   - A fully initialized ConfigWizard ready for use.
//
// Side effects:
//   - None.
func NewConfigWizard(width, height int) *ConfigWizard {
	return NewConfigWizardWithProfiles(width, height, nil)
}

// NewConfigWizardWithProfiles creates a wizard modal with profile options.
//
// Expected:
//   - width and height must be valid positive integers.
//   - profiles can be nil or a slice of ProfileOption.
//
// Returns:
//   - A fully initialized ConfigWizard ready for use.
//
// Side effects:
//   - None.
func NewConfigWizardWithProfiles(width, height int, profiles []ProfileOption) *ConfigWizard {
	modal := &ConfigWizard{
		data:           &ConfigData{},
		width:          width,
		height:         height,
		profileOptions: profiles,
		techsAvailable: false,
	}

	modal.applyDefaults()
	modal.buildForm()

	wizard := behaviors.NewWizardBehavior[ConfigData](modal.form, modal.data)
	modal.wizard = wizard

	return modal
}

// buildForm creates the huh form with 3 steps (groups).
func (m *ConfigWizard) buildForm() {
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
func (m *ConfigWizard) Init() tea.Cmd {
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
//   - true if the message was handled, false otherwise.
//
// Side effects:
//   - May update internal state based on key messages.
func (m *ConfigWizard) handleWizardKeyMsg(msg tea.KeyMsg) bool {
	switch msg.String() {
	case "ctrl+s":
		if !m.HasRequiredFields() {
			return true
		}
		m.applyDefaults()
		m.wizard.Skip()
		return true

	case "esc":
		if m.form == nil || m.form.CurrentStep() == 0 {
			m.wizard.Hide()
			return true
		}
		newStep := m.form.CurrentStep() - 1
		if newStep == 1 && !m.techsAvailable {
			newStep = 0
		}
		m.form.SetCurrentStep(newStep)
		return true
	}
	return false
}

func (m *ConfigWizard) advanceWizardStep(msg tea.Msg) {
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

// Update processes incoming messages and updates the wizard state.
//
// Expected:
//   - msg: the incoming Bubble Tea message to process.
//
// Returns:
//   - a tea.Cmd for any side effects (or nil).
//
// Side effects:
//   - updates wizard form state, step navigation, and dimensions.
func (m *ConfigWizard) Update(msg tea.Msg) tea.Cmd {
	if !m.wizard.IsVisible() {
		return nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.handleWizardKeyMsg(msg) {
			return nil
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
		m.advanceWizardStep(msg)
	}

	return cmd
}

// View renders the wizard modal.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *ConfigWizard) View() string {
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
func (m *ConfigWizard) buildFooter() string {
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
//   - None.
func (m *ConfigWizard) Show() {
	m.wizard.Show()
}

// Reset resets the wizard state while preserving the entered data.
//
// Side effects:
//   - None.
func (m *ConfigWizard) Reset() {
	m.wizard.Reset()

	m.buildForm()
	m.wizard.SetForm(m.form)
}

// Hide makes the modal invisible.
//
// Side effects:
//   - None.
func (m *ConfigWizard) Hide() {
	m.wizard.Hide()
}

// IsVisible returns whether the modal is currently visible.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *ConfigWizard) IsVisible() bool {
	return m.wizard.IsVisible()
}

// IsCompleted returns whether the wizard has been completed.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *ConfigWizard) IsCompleted() bool {
	return m.wizard.IsCompleted()
}

// IsSkipped returns whether the wizard was skipped (Ctrl+S).
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *ConfigWizard) IsSkipped() bool {
	return m.wizard.IsSkipped()
}

// GetConfigData returns the collected configuration data.
//
// Returns:
//   - A fully initialized ConfigData ready for use.
//
// Side effects:
//   - None.
func (m *ConfigWizard) GetConfigData() *ConfigData {
	m.syncFromFormData()
	return m.data
}

// GetCurrentStep returns the current step index (0-based).
//
// Returns:
//   - A int value.
//
// Side effects:
//   - None.
func (m *ConfigWizard) GetCurrentStep() int {
	return m.wizard.CurrentStep()
}

// GetStepCount returns the total number of steps.
//
// Returns:
//   - A int value.
//
// Side effects:
//   - None.
func (m *ConfigWizard) GetStepCount() int {
	return 3
}

// AreTechsAvailable returns whether extracted technologies are available.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *ConfigWizard) AreTechsAvailable() bool {
	return m.techsAvailable
}

// GetExtractedTechnologies returns the extracted technologies.
//
// Returns:
//   - A []ExtractedTechnology value.
//
// Side effects:
//   - None.
func (m *ConfigWizard) GetExtractedTechnologies() []ExtractedTechnology {
	return m.extractedTechs
}

// SetExtractedTechnologies sets the extracted technologies and rebuilds the form.
//
// Expected:
//   - []extractedtechnology must be valid.
//
// Side effects:
//   - None.
func (m *ConfigWizard) SetExtractedTechnologies(techs []ExtractedTechnology) {
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
//   - A []ProfileOption value.
//
// Side effects:
//   - None.
func (m *ConfigWizard) GetProfileOptions() []ProfileOption {
	return m.profileOptions
}

// SetProfileID sets the profile ID.
//
// Expected:
//   - Must be a valid string.
//
// Side effects:
//   - None.
func (m *ConfigWizard) SetProfileID(id string) {
	if m.formData != nil {
		m.formData.ProfileID = id
	}
	m.data.ProfileID = id
}

// SetAudience sets the target audience.
//
// Expected:
//   - Must be a valid string.
//
// Side effects:
//   - None.
func (m *ConfigWizard) SetAudience(audience string) {
	if m.formData != nil {
		m.formData.Audience = audience
	}
	m.data.Audience = audience
}

// SetTechFocus sets the technology focus.
//
// Expected:
//   - Must be a valid string.
//
// Side effects:
//   - None.
func (m *ConfigWizard) SetTechFocus(focus string) {
	if m.formData != nil {
		m.formData.TechFocus = focus
	}
	m.data.TechFocus = focus
}

// SetTechnologies sets the selected technologies (for generalist mode).
//
// Expected:
//   - Must be a valid string.
//
// Side effects:
//   - None.
func (m *ConfigWizard) SetTechnologies(techs []string) {
	if m.formData != nil {
		m.formData.Technologies = techs
	}
	m.data.Technologies = techs
}

// SetTechnology sets the selected technology (for specialist mode).
//
// Expected:
//   - Must be a valid string.
//
// Side effects:
//   - None.
func (m *ConfigWizard) SetTechnology(tech string) {
	if m.formData != nil {
		m.formData.Technology = tech
	}
	m.data.Technology = tech
}

// SetFocusArea sets the focus area.
//
// Expected:
//   - Must be a valid string.
//
// Side effects:
//   - None.
func (m *ConfigWizard) SetFocusArea(area string) {
	if m.formData != nil {
		m.formData.FocusArea = area
	}
	m.data.FocusArea = area
}

// SetSkillsFormat sets the skills format.
//
// Expected:
//   - Must be a valid string.
//
// Side effects:
//   - None.
func (m *ConfigWizard) SetSkillsFormat(format string) {
	if m.formData != nil {
		m.formData.SkillsFormat = format
	}
	m.data.SkillsFormat = format
}

// SetSkillsLimit sets the skills limit per category/total.
//
// Expected:
//   - int must be valid.
//
// Side effects:
//   - None.
func (m *ConfigWizard) SetSkillsLimit(limit int) {
	if m.formData != nil {
		m.formData.SkillsLimit = limit
	}
	m.data.SkillsLimit = limit
}

// SetCVLength sets the CV length.
//
// Expected:
//   - Must be a valid string.
//
// Side effects:
//   - None.
func (m *ConfigWizard) SetCVLength(length string) {
	if m.formData != nil {
		m.formData.CVLength = length
	}
	m.data.CVLength = length
}

// Complete marks the wizard as completed.
//
// Side effects:
//   - None.
func (m *ConfigWizard) Complete() {
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
//   - None.
func (m *ConfigWizard) HasRequiredFields() bool {
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
func (m *ConfigWizard) GetDimensions() (width, height int) {
	return m.width, m.height
}

func (m *ConfigWizard) syncFromFormData() {
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

func (m *ConfigWizard) applyDefaults() {
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
		m.data.CVLength = "detailed"
		if m.formData != nil {
			m.formData.CVLength = "detailed"
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
