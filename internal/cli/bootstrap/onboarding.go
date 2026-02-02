package bootstrap

import (
	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/config"
	tea "github.com/charmbracelet/bubbletea"
)

// onboardingModel wraps the OnboardingWizardModal for standalone execution.
// This runs as a separate Bubble Tea program before the main app starts.
type onboardingModel struct {
	wizard *components.OnboardingWizardModal
	width  int
	height int
	result *config.ProfileConfig
}

// newOnboardingModel creates a new onboarding model.
func newOnboardingModel(cfg *config.ProfileConfig) *onboardingModel {
	// Default dimensions - will be updated on first WindowSizeMsg.
	width, height := 80, 24

	return &onboardingModel{
		wizard: components.NewOnboardingWizardModalWithConfig(width, height, cfg),
		width:  width,
		height: height,
	}
}

// Init initializes the onboarding model.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (m *onboardingModel) Init() tea.Cmd {
	return tea.Batch(
		tea.WindowSize(),
		m.wizard.Init(),
	)
}

// Update handles messages for the onboarding model.
//
// Expected:
//   - msg must be a valid tea.Msg type.
//
// Returns:
//   - Updated model.
//   - Command to execute.
//
// Side effects:
//   - May update wizard dimensions.
//   - May complete onboarding and return result.
//   - May quit on Ctrl+C.
func (m *onboardingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Forward to wizard for dimension updates.
		cmd := m.wizard.Update(msg)
		return m, cmd

	case tea.KeyMsg:
		// Ctrl+C is the only way to abort onboarding.
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	// Forward all messages to wizard.
	cmd := m.wizard.Update(msg)

	// Check if wizard completed.
	if m.wizard.IsCompleted() {
		m.result = m.wizard.GetProfileConfig()
		return m, tea.Quit
	}

	return m, cmd
}

// View renders the onboarding screen.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *onboardingModel) View() string {
	wizardView := m.wizard.View()
	return primitives.CenterInTerminal(wizardView, m.width, m.height)
}

// Result returns the profile config if onboarding completed successfully.
//
// Returns:
//   - A fully initialized config.ProfileConfig ready for use.
//
// Side effects:
//   - None.
func (m *onboardingModel) Result() *config.ProfileConfig {
	return m.result
}

// runOnboarding runs the onboarding wizard as a standalone Bubble Tea program.
// Returns the completed profile config, or ErrUserAborted if user cancelled.
func runOnboarding(existingProfile *config.ProfileConfig) (*config.ProfileConfig, error) {
	model := newOnboardingModel(existingProfile)

	p := tea.NewProgram(model, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return nil, err
	}

	// Extract result from final model.
	if m, ok := finalModel.(*onboardingModel); ok {
		result := m.Result()
		if result == nil {
			return nil, ErrUserAborted
		}
		return result, nil
	}

	return nil, ErrUserAborted
}
