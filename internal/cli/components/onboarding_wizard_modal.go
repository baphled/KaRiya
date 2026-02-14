package components

import (
	"strings"

	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/cli/uikit/containers"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/cli/uikit/theme"
	"github.com/baphled/kariya/internal/config"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// OnboardingWizardModal provides a 3-step wizard for first-run profile setup.
// The wizard guides users through:
//   - Step 1: Welcome + Name (required)
//   - Step 2: Contact: Email (required) + Location (optional)
//   - Step 3: Professional: Title, GitHub, Portfolio (all optional)
//
// Navigation:
//   - Tab/Enter: Move to next field
//   - Shift+Tab: Move to previous field
//   - Enter on last field of group: Advance to next step
//
// Onboarding is mandatory: the Esc key is blocked and users must complete
// the required fields (Name and Email) to proceed with the application.
type OnboardingWizardModal struct {
	wizard   *behaviors.WizardBehavior[OnboardingData]
	form     *forms.WizardFormAdapter
	formData *forms.OnboardingFormData
	data     *OnboardingData
	width    int
	height   int
}

// OnboardingData holds the data collected from the onboarding wizard.
type OnboardingData struct {
	// Step 1: Welcome
	Name string // Required

	// Step 2: Contact
	Email    string // Required
	Location string // Optional

	// Step 3: Professional
	Title     string // Optional
	GitHub    string // Optional
	Portfolio string // Optional
}

// NewOnboardingWizardModal creates a new onboarding wizard modal.
//
// Expected:
//   - int must be valid.
//
// Returns:
//   - A fully initialized OnboardingWizardModal ready for use.
//
// Side effects:
//   - None.
func NewOnboardingWizardModal(width, height int) *OnboardingWizardModal {
	return NewOnboardingWizardModalWithConfig(width, height, nil)
}

// NewOnboardingWizardModalWithConfig creates an onboarding wizard with pre-filled data.
//
// Expected:
//   - int must be valid.
//   - config must be a valid configuration object.
//
// Returns:
//   - A fully initialized OnboardingWizardModal ready for use.
//
// Side effects:
//   - None.
func NewOnboardingWizardModalWithConfig(width, height int, cfg *config.ProfileConfig) *OnboardingWizardModal {
	data := &OnboardingData{}

	if cfg != nil {
		data.Name = cfg.Name
		data.Email = cfg.Email
		data.Location = cfg.Location
		data.Title = cfg.Title
		data.GitHub = cfg.GitHub
		data.Portfolio = cfg.Portfolio
	}

	formData := toOnboardingFormData(data)
	modalWidth := calcOnboardingModalWidth(width)
	huhForm := forms.NewOnboardingWizardForm(formData, modalWidth, height)
	adapter := forms.NewWizardFormAdapter(huhForm, 3)
	wizard := behaviors.NewWizardBehavior[OnboardingData](adapter, data)

	modal := &OnboardingWizardModal{
		wizard:   wizard,
		form:     adapter,
		formData: formData,
		data:     data,
		width:    width,
		height:   height,
	}

	return modal
}

// Init initializes the wizard modal and its form.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (m *OnboardingWizardModal) Init() tea.Cmd {
	return m.wizard.Init()
}

// Update handles messages for the wizard modal.
//
// Expected:
//   - msg must be valid.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (m *OnboardingWizardModal) Update(msg tea.Msg) tea.Cmd {
	if !m.wizard.IsVisible() {
		return nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc {
			return nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.form.SetDimensions(calcOnboardingModalWidth(m.width), m.height)
		return m.form.Init()
	}

	cmd := m.wizard.Update(msg)

	if m.wizard.IsCompleted() {
		m.syncFromFormData()
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
func (m *OnboardingWizardModal) View() string {
	if !m.wizard.IsVisible() && !m.wizard.IsCompleted() {
		return ""
	}

	modalWidth := calcOnboardingModalWidth(m.width)

	modalHeight := m.height - 10
	if modalHeight < 20 {
		modalHeight = 20
	}

	th := theme.Default()

	title := primitives.Title("Profile Setup", th).Width(modalWidth - 4).Align(lipgloss.Center).Render()

	var formView string
	if m.wizard.IsCompleted() {
		formView = m.buildSummary()
	} else {
		formView = m.wizard.View()
	}

	footer := m.buildFooter()
	content := lipgloss.JoinVertical(lipgloss.Left, title, "", formView, "", footer)

	box := containers.NewBox(th).
		Content(content).
		Width(modalWidth).
		MaxHeight(modalHeight).
		Background(th.BackgroundColor()).
		BorderColor(th.PrimaryColor()).
		Padding(1)

	return box.Render()
}

// buildFooter creates the keyboard shortcuts footer using UIKit primitives.
func (m *OnboardingWizardModal) buildFooter() string {
	th := theme.Default()

	badges := []*primitives.Badge{
		primitives.KeyBadge("tab", "next field", th),
		primitives.KeyBadge("shift+tab", "prev field", th),
		primitives.KeyBadge("enter", "continue", th),
	}

	return primitives.RenderHelpFooter(th, badges...)
}

// buildSummary creates a summary view of the completed profile data.
func (m *OnboardingWizardModal) buildSummary() string {
	var lines []string
	lines = append(lines, "Profile Setup Complete!")
	lines = append(lines, "")

	if m.data.Name != "" {
		lines = append(lines, "Name: "+m.data.Name)
	}
	if m.data.Email != "" {
		lines = append(lines, "Email: "+m.data.Email)
	}
	if m.data.Location != "" {
		lines = append(lines, "Location: "+m.data.Location)
	}
	if m.data.Title != "" {
		lines = append(lines, "Title: "+m.data.Title)
	}
	if m.data.GitHub != "" {
		lines = append(lines, "GitHub: "+m.data.GitHub)
	}
	if m.data.Portfolio != "" {
		lines = append(lines, "Portfolio: "+m.data.Portfolio)
	}

	return strings.Join(lines, "\n")
}

// IsVisible returns whether the modal is currently visible.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *OnboardingWizardModal) IsVisible() bool {
	return m.wizard.IsVisible()
}

// IsCompleted returns whether the wizard was completed successfully.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *OnboardingWizardModal) IsCompleted() bool {
	return m.wizard.IsCompleted()
}

// WasCancelled returns whether the wizard was cancelled by the user.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *OnboardingWizardModal) WasCancelled() bool {
	return m.wizard.IsCancelled()
}

// CurrentStep returns the current step index (0-based).
//
// Returns:
//   - A int value.
//
// Side effects:
//   - None.
func (m *OnboardingWizardModal) CurrentStep() int {
	return m.wizard.CurrentStep()
}

// TotalSteps returns the total number of steps in the wizard.
//
// Returns:
//   - A int value.
//
// Side effects:
//   - None.
func (m *OnboardingWizardModal) TotalSteps() int {
	return m.wizard.TotalSteps()
}

// Hide hides the modal.
//
// Side effects:
//   - None.
func (m *OnboardingWizardModal) Hide() {
	m.wizard.Hide()
}

// Show shows the modal.
//
// Side effects:
//   - None.
func (m *OnboardingWizardModal) Show() {
	m.wizard.Show()
}

// HasRequiredFields returns whether all required fields have values.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *OnboardingWizardModal) HasRequiredFields() bool {
	m.syncFromFormData()
	return strings.TrimSpace(m.data.Name) != "" && strings.TrimSpace(m.data.Email) != ""
}

// GetProfileConfig returns the collected data as a ProfileConfig.
//
// Returns:
//   - A fully initialized config.ProfileConfig ready for use.
//
// Side effects:
//   - None.
func (m *OnboardingWizardModal) GetProfileConfig() *config.ProfileConfig {
	if !m.wizard.IsCompleted() {
		return nil
	}

	m.syncFromFormData()
	return &config.ProfileConfig{
		Name:      strings.TrimSpace(m.data.Name),
		Email:     strings.TrimSpace(m.data.Email),
		Location:  strings.TrimSpace(m.data.Location),
		Title:     strings.TrimSpace(m.data.Title),
		GitHub:    strings.TrimSpace(m.data.GitHub),
		Portfolio: strings.TrimSpace(m.data.Portfolio),
	}
}

// GetOnboardingData returns the raw onboarding data.
//
// Returns:
//   - A fully initialized OnboardingData ready for use.
//
// Side effects:
//   - None.
func (m *OnboardingWizardModal) GetOnboardingData() *OnboardingData {
	return m.data
}

func calcOnboardingModalWidth(terminalWidth int) int {
	modalWidth := terminalWidth - 20
	if modalWidth > 80 {
		modalWidth = 80
	}
	if modalWidth < 50 {
		modalWidth = 50
	}
	return modalWidth
}

func toOnboardingFormData(data *OnboardingData) *forms.OnboardingFormData {
	return &forms.OnboardingFormData{
		Name:      data.Name,
		Email:     data.Email,
		Location:  data.Location,
		Title:     data.Title,
		GitHub:    data.GitHub,
		Portfolio: data.Portfolio,
	}
}

func (m *OnboardingWizardModal) syncFromFormData() {
	if m.formData == nil {
		return
	}
	m.data.Name = m.formData.Name
	m.data.Email = m.formData.Email
	m.data.Location = m.formData.Location
	m.data.Title = m.formData.Title
	m.data.GitHub = m.formData.GitHub
	m.data.Portfolio = m.formData.Portfolio
}
