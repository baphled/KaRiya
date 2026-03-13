package onboarding

import (
	"strings"

	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/tui/forms"
	"github.com/baphled/kariya/internal/ui/behaviors"
	"github.com/baphled/kariya/internal/ui/uikit/containers"
	"github.com/baphled/kariya/internal/ui/uikit/primitives"
	"github.com/baphled/kariya/internal/ui/uikit/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Wizard provides a 3-step wizard for first-run profile setup.
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
type Wizard struct {
	wizard   *behaviors.WizardBehavior[Data]
	form     *forms.WizardFormAdapter
	formData *forms.OnboardingFormData
	data     *Data
	width    int
	height   int
}

// Data holds the data collected from the onboarding wizard.
type Data struct {
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

// NewWizard creates a new onboarding wizard.
//
// Expected:
//   - int must be valid.
//
// Returns:
//   - A fully initialized Wizard ready for use.
//
// Side effects:
//   - None.
func NewWizard(width, height int) *Wizard {
	return NewWizardWithConfig(width, height, nil)
}

// NewWizardWithConfig creates an onboarding wizard with pre-filled data.
//
// Expected:
//   - int must be valid.
//   - config must be a valid configuration object.
//
// Returns:
//   - A fully initialized Wizard ready for use.
//
// Side effects:
//   - None.
func NewWizardWithConfig(width, height int, cfg *config.ProfileConfig) *Wizard {
	data := &Data{}

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
	wizard := behaviors.NewWizardBehavior[Data](adapter, data)

	modal := &Wizard{
		wizard:   wizard,
		form:     adapter,
		formData: formData,
		data:     data,
		width:    width,
		height:   height,
	}

	return modal
}

// Init initializes the wizard and its form.
//
// Returns:
//   - A tea.Cmd to execute.
//
// Side effects:
//   - Initializes the underlying wizard behavior.
func (m *Wizard) Init() tea.Cmd {
	return m.wizard.Init()
}

// Update handles messages for the wizard.
//
// Expected:
//   - msg is a Bubble Tea message.
//
// Returns:
//   - A tea.Cmd to execute.
//
// Side effects:
//   - Updates wizard form state and dimensions on window resize.
//   - Syncs form data when wizard is completed.
func (m *Wizard) Update(msg tea.Msg) tea.Cmd {
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

// View renders the wizard.
//
// Returns:
//   - A string containing the rendered wizard UI.
//
// Side effects:
//   - None.
func (m *Wizard) View() string {
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
func (m *Wizard) buildFooter() string {
	th := theme.Default()
	badges := []*primitives.Badge{
		primitives.KeyBadge("tab", "next field", th),
		primitives.KeyBadge("shift+tab", "prev field", th),
		primitives.KeyBadge("enter", "continue", th),
	}
	return primitives.RenderHelpFooter(th, badges...)
}

// buildSummary creates a summary view of the completed profile data.
func (m *Wizard) buildSummary() string {
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

// IsVisible returns whether the wizard is currently visible.
//
// Returns:
//   - true if the wizard is visible, false otherwise.
//
// Side effects:
//   - None.
func (m *Wizard) IsVisible() bool {
	return m.wizard.IsVisible()
}

// IsCompleted returns whether the wizard was completed successfully.
//
// Returns:
//   - true if the wizard was completed, false otherwise.
//
// Side effects:
//   - None.
func (m *Wizard) IsCompleted() bool {
	return m.wizard.IsCompleted()
}

// WasCancelled returns whether the wizard was cancelled by the user.
//
// Returns:
//   - true if the wizard was cancelled, false otherwise.
//
// Side effects:
//   - None.
func (m *Wizard) WasCancelled() bool {
	return m.wizard.IsCancelled()
}

// CurrentStep returns the current step index (0-based).
//
// Returns:
//   - The current step index as an integer.
//
// Side effects:
//   - None.
func (m *Wizard) CurrentStep() int {
	return m.wizard.CurrentStep()
}

// TotalSteps returns the total number of steps in the wizard.
//
// Returns:
//   - The total number of steps as an integer.
//
// Side effects:
//   - None.
func (m *Wizard) TotalSteps() int {
	return m.wizard.TotalSteps()
}

// Hide hides the wizard.
//
// Side effects:
//   - Hides the wizard from view.
func (m *Wizard) Hide() {
	m.wizard.Hide()
}

// Show shows the wizard.
//
// Side effects:
//   - Shows the wizard in view.
func (m *Wizard) Show() {
	m.wizard.Show()
}

// HasRequiredFields returns whether all required fields have values.
//
// Returns:
//   - true if Name and Email are non-empty, false otherwise.
//
// Side effects:
//   - Syncs form data to internal state.
func (m *Wizard) HasRequiredFields() bool {
	m.syncFromFormData()
	return strings.TrimSpace(m.data.Name) != "" && strings.TrimSpace(m.data.Email) != ""
}

// GetProfileConfig returns the collected data as a ProfileConfig.
//
// Returns:
//   - A ProfileConfig with trimmed field values, or nil if wizard is not completed.
//
// Side effects:
//   - Syncs form data to internal state.
func (m *Wizard) GetProfileConfig() *config.ProfileConfig {
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

// GetData returns the raw onboarding data.
//
// Returns:
//   - A pointer to the Data struct containing collected profile information.
//
// Side effects:
//   - None.
func (m *Wizard) GetData() *Data {
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

func toOnboardingFormData(data *Data) *forms.OnboardingFormData {
	return &forms.OnboardingFormData{
		Name:      data.Name,
		Email:     data.Email,
		Location:  data.Location,
		Title:     data.Title,
		GitHub:    data.GitHub,
		Portfolio: data.Portfolio,
	}
}

func (m *Wizard) syncFromFormData() {
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
