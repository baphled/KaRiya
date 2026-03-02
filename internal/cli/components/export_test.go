package components

// CompleteWizardForTesting marks the wizard as completed for test purposes.
// This allows external test packages to trigger completion-dependent paths
// like buildSummary and GetProfileConfig without driving the full huh form.
func (m *OnboardingWizardModal) CompleteWizardForTesting() {
	m.wizard.Complete()
}

// SetDataForTesting populates the OnboardingData directly for test purposes.
// This bypasses form syncing so tests can verify buildSummary output.
func (m *OnboardingWizardModal) SetDataForTesting(data *OnboardingData) {
	m.data = data
}

// CalcOnboardingModalWidthForTesting exposes calcOnboardingModalWidth for testing.
func CalcOnboardingModalWidthForTesting(terminalWidth int) int {
	return calcOnboardingModalWidth(terminalWidth)
}
