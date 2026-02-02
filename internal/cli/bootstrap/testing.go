package bootstrap

import (
	"github.com/baphled/kariya/internal/config"
	tea "github.com/charmbracelet/bubbletea"
)

// OnboardingTestModel wraps the onboarding model for testing.
// This allows E2E tests to verify the onboarding flow.
type OnboardingTestModel struct {
	model *onboardingModel
}

// NewOnboardingTestModel creates an onboarding model for testing.
//
// Expected:
//   - config must be a valid configuration object.
//
// Returns:
//   - A fully initialized OnboardingTestModel ready for use.
//
// Side effects:
//   - None.
func NewOnboardingTestModel(existingProfile *config.ProfileConfig) *OnboardingTestModel {
	return &OnboardingTestModel{
		model: newOnboardingModel(existingProfile),
	}
}

// Init returns the initialization command for the onboarding model.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (m *OnboardingTestModel) Init() tea.Cmd {
	return m.model.Init()
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
//   - Delegates to underlying model.
func (m *OnboardingTestModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.model.Update(msg)
}

// View renders the onboarding screen.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *OnboardingTestModel) View() string {
	return m.model.View()
}

// IsCompleted returns whether the wizard was completed successfully.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *OnboardingTestModel) IsCompleted() bool {
	return m.model.result != nil
}

// Result returns the profile config if onboarding completed successfully.
//
// Returns:
//   - A fully initialized config.ProfileConfig ready for use.
//
// Side effects:
//   - None.
func (m *OnboardingTestModel) Result() *config.ProfileConfig {
	return m.model.Result()
}
