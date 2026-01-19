package components

import (
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/cli/uikit/theme"
	"github.com/baphled/kariya/internal/config"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
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
//   - Enter on last field: Complete step
//   - Esc: Cancel wizard (step 1) or go back (steps 2-3)
//
// Usage:
//
//	modal := components.NewOnboardingWizardModal(width, height)
//	cmd := modal.Init()
//	// In Update:
//	cmd := modal.Update(msg)
//	if modal.IsCompleted() {
//	    cfg := modal.GetProfileConfig()
//	    // Save config
//	}
type OnboardingWizardModal struct {
	form        *huh.Form
	data        *OnboardingData
	currentStep int
	totalSteps  int
	visible     bool
	completed   bool
	cancelled   bool
	width       int
	height      int
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
func NewOnboardingWizardModal(width, height int) *OnboardingWizardModal {
	return NewOnboardingWizardModalWithConfig(width, height, nil)
}

// NewOnboardingWizardModalWithConfig creates an onboarding wizard with pre-filled data.
func NewOnboardingWizardModalWithConfig(width, height int, cfg *config.ProfileConfig) *OnboardingWizardModal {
	data := &OnboardingData{}

	// Pre-fill from existing config if provided
	if cfg != nil {
		data.Name = cfg.Name
		data.Email = cfg.Email
		data.Location = cfg.Location
		data.Title = cfg.Title
		data.GitHub = cfg.GitHub
		data.Portfolio = cfg.Portfolio
	}

	modal := &OnboardingWizardModal{
		data:        data,
		currentStep: 0,
		totalSteps:  3,
		visible:     true,
		width:       width,
		height:      height,
	}

	modal.buildForm()
	return modal
}

// buildForm creates the huh form with 3 steps (groups).
func (m *OnboardingWizardModal) buildForm() {
	// Calculate modal dimensions
	modalWidth := m.width - 20
	if modalWidth > 70 {
		modalWidth = 70
	}
	if modalWidth < 50 {
		modalWidth = 50
	}

	// Get theme for consistent styling
	huhTheme := forms.Theme()

	// Step 1: Welcome + Name
	step1 := huh.NewGroup(
		huh.NewNote().
			Title("Welcome to KaRiya!").
			Description("Let's set up your profile for CV generation.\nThis information will appear on your CVs."),
		huh.NewInput().
			Key("name").
			Title("Your Name").
			Description("Required - appears at the top of your CV").
			Placeholder("e.g., Jane Doe").
			Validate(func(s string) error {
				if strings.TrimSpace(s) == "" {
					return fmt.Errorf("name is required")
				}
				return nil
			}).
			Value(&m.data.Name),
	).Title("Step 1 of 3: Welcome")

	// Step 2: Contact
	step2 := huh.NewGroup(
		huh.NewInput().
			Key("email").
			Title("Email Address").
			Description("Required - contact information for your CV").
			Placeholder("e.g., jane@example.com").
			Validate(func(s string) error {
				if strings.TrimSpace(s) == "" {
					return fmt.Errorf("email is required")
				}
				if !strings.Contains(s, "@") {
					return fmt.Errorf("please enter a valid email address")
				}
				return nil
			}).
			Value(&m.data.Email),
		huh.NewInput().
			Key("location").
			Title("Location").
			Description("Optional - e.g., city, country, or 'Remote'").
			Placeholder("e.g., London, UK").
			Value(&m.data.Location),
	).Title("Step 2 of 3: Contact")

	// Step 3: Professional
	step3 := huh.NewGroup(
		huh.NewInput().
			Key("title").
			Title("Professional Title").
			Description("Optional - your current role or target role").
			Placeholder("e.g., Senior Software Engineer").
			Value(&m.data.Title),
		huh.NewInput().
			Key("github").
			Title("GitHub Profile").
			Description("Optional - link to your GitHub").
			Placeholder("e.g., https://github.com/username").
			Value(&m.data.GitHub),
		huh.NewInput().
			Key("portfolio").
			Title("Portfolio/Website").
			Description("Optional - personal website or portfolio").
			Placeholder("e.g., https://janedoe.dev").
			Value(&m.data.Portfolio),
	).Title("Step 3 of 3: Professional Details")

	m.form = huh.NewForm(step1, step2, step3).
		WithTheme(huhTheme).
		WithWidth(modalWidth).
		WithShowHelp(true).
		WithShowErrors(true)
}

// Init initializes the wizard modal and its form.
func (m *OnboardingWizardModal) Init() tea.Cmd {
	if m.form == nil {
		return nil
	}
	return m.form.Init()
}

// Update handles messages for the wizard modal.
func (m *OnboardingWizardModal) Update(msg tea.Msg) tea.Cmd {
	if !m.visible {
		return nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// If at first step, cancel wizard
			if m.currentStep == 0 {
				m.cancelled = true
				m.visible = false
				return nil
			}
			// Go back a step
			if m.currentStep > 0 {
				m.currentStep--
			}
			// Let form handle the back navigation
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
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
			// Check if Enter was pressed to advance
			if keyMsg, ok := msg.(tea.KeyMsg); ok {
				if keyMsg.String() == "enter" && m.currentStep < m.totalSteps-1 {
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
func (m *OnboardingWizardModal) View() string {
	if !m.visible {
		return ""
	}

	if m.form == nil {
		return ""
	}

	// Calculate modal dimensions
	modalWidth := m.width - 10
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

	// Title style
	titleStyle := lipgloss.NewStyle().
		Foreground(styles.ColorAccentTeal).
		Bold(true).
		Align(lipgloss.Center).
		Width(modalWidth - 4)

	title := titleStyle.Render("Profile Setup")

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
		Background(styles.ColorBackground).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorAccentTeal).
		Padding(1).
		Render(content)

	// Center the modal
	centeredStyle := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center)

	return centeredStyle.Render(styledContent)
}

// buildFooter creates the keyboard shortcuts footer using UIKit primitives.
func (m *OnboardingWizardModal) buildFooter() string {
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

	return primitives.RenderHelpFooter(th, badges...)
}

// IsVisible returns whether the modal is currently visible.
func (m *OnboardingWizardModal) IsVisible() bool {
	return m.visible
}

// IsCompleted returns whether the wizard was completed successfully.
func (m *OnboardingWizardModal) IsCompleted() bool {
	return m.completed
}

// WasCancelled returns whether the wizard was cancelled by the user.
func (m *OnboardingWizardModal) WasCancelled() bool {
	return m.cancelled
}

// CurrentStep returns the current step index (0-based).
func (m *OnboardingWizardModal) CurrentStep() int {
	return m.currentStep
}

// TotalSteps returns the total number of steps in the wizard.
func (m *OnboardingWizardModal) TotalSteps() int {
	return m.totalSteps
}

// Hide hides the modal.
func (m *OnboardingWizardModal) Hide() {
	m.visible = false
}

// Show shows the modal.
func (m *OnboardingWizardModal) Show() {
	m.visible = true
}

// HasRequiredFields returns whether all required fields have values.
func (m *OnboardingWizardModal) HasRequiredFields() bool {
	return strings.TrimSpace(m.data.Name) != "" && strings.TrimSpace(m.data.Email) != ""
}

// GetProfileConfig returns the collected data as a ProfileConfig.
// Returns nil if the wizard was not completed.
func (m *OnboardingWizardModal) GetProfileConfig() *config.ProfileConfig {
	if !m.completed {
		return nil
	}

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
func (m *OnboardingWizardModal) GetOnboardingData() *OnboardingData {
	return m.data
}
