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
// Use this in E2E tests that need to verify the onboarding wizard.
func NewOnboardingTestModel(existingProfile *config.ProfileConfig) *OnboardingTestModel {
	return &OnboardingTestModel{
		model: newOnboardingModel(existingProfile),
	}
}

// Init returns the initialization command for the onboarding model.
func (m *OnboardingTestModel) Init() tea.Cmd {
	return m.model.Init()
}

// Update handles messages for the onboarding model.
func (m *OnboardingTestModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.model.Update(msg)
}

// View renders the onboarding screen.
func (m *OnboardingTestModel) View() string {
	return m.model.View()
}

// IsCompleted returns whether the wizard was completed successfully.
func (m *OnboardingTestModel) IsCompleted() bool {
	return m.model.result != nil
}

// Result returns the profile config if onboarding completed successfully.
func (m *OnboardingTestModel) Result() *config.ProfileConfig {
	return m.model.Result()
}
