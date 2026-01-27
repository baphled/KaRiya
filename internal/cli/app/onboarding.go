package app

import (
	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/config"
	tea "github.com/charmbracelet/bubbletea"
)

// handleOnboardingInput handles input during the onboarding wizard.
// Note: Onboarding is mandatory - users cannot skip or cancel this wizard.
// They must provide Name and Email to proceed with the application.
func (m *Model) handleOnboardingInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.onboardingWizard == nil {
		m.state = StateMenu
		return m, nil
	}

	cmd := m.onboardingWizard.Update(msg)

	// Check if wizard completed.
	if m.onboardingWizard.IsCompleted() {
		return m.completeOnboarding()
	}

	// Note: Onboarding cannot be cancelled - users must complete it.
	// The wizard ignores Esc key presses.

	return m, cmd
}

// completeOnboarding handles the completion of the onboarding wizard.
func (m *Model) completeOnboarding() (tea.Model, tea.Cmd) {
	// Get the profile config from the wizard.
	profileCfg := m.onboardingWizard.GetProfileConfig()
	if profileCfg != nil {
		// Update the app config with the new profile.
		m.appConfig.Profile = *profileCfg

		// Save the config.
		if err := config.SaveConfig(m.appConfig); err != nil {
			m.logger.Error("Failed to save config: %v", err)
		} else {
			m.logger.Info("Profile saved successfully")
		}
	}

	// Transition to main menu.
	m.state = StateMenu
	m.onboardingWizard = nil
	return m, nil
}

// viewOnboarding renders the onboarding wizard.
func (m *Model) viewOnboarding() string {
	if m.onboardingWizard == nil {
		return ""
	}

	// Center the wizard modal in the terminal.
	wizardView := m.onboardingWizard.View()
	return primitives.CenterInTerminal(wizardView, m.width, m.height)
}

// SkipOnboarding skips the onboarding wizard and goes directly to menu.
// This is primarily used by tests to avoid the onboarding flow.
func (m *Model) SkipOnboarding() {
	if m.state == StateOnboarding {
		m.state = StateMenu
		m.onboardingWizard = nil
	}
}

// ForceOnboarding forces the onboarding wizard to appear, regardless of config.
// This is primarily used by tests to verify the onboarding flow.
// It creates a fresh wizard with empty data (not the user's existing config).
func (m *Model) ForceOnboarding() {
	m.state = StateOnboarding
	m.onboardingWizard = components.NewOnboardingWizardModal(m.width, m.height)
}
