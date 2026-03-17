package steps

import (
	"context"
	"strings"

	"github.com/baphled/kariya/features/support"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cucumber/godog"
	"github.com/onsi/gomega"
)

// RegisterConfigureSteps registers configure system step definitions with Godog.
//
// Expected:
//   - scenariocontext must be valid.
//
// Side effects:
//   - None.
func RegisterConfigureSteps(sc *godog.ScenarioContext) {
	sc.Step(`^I press "," to open settings$`, iPressCommaToOpenSettings)
	sc.Step(`^I should see the settings modal$`, iShouldSeeTheSettingsModal)
	sc.Step(`^I navigate to the "([^"]*)" section$`, iNavigateToSection)
}

func iPressCommaToOpenSettings(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune(',')
	return ctx, nil
}

func iShouldSeeTheSettingsModal(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Configure System"),
		gomega.ContainSubstring("System"),
		gomega.ContainSubstring("Profile"),
	))
	return nil
}

func iNavigateToSection(ctx context.Context, section string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	sectionIndicators := map[string]string{
		"System":  "Log Level",
		"Profile": "Full Name",
		"Export":  "Default Destination",
		"UI":      "Theme",
	}
	targetIndicator, ok := sectionIndicators[section]
	if !ok {
		return ctx, godog.ErrPending
	}
	for range 4 {
		view := env.GetView()
		if strings.Contains(view, targetIndicator) {
			return ctx, nil
		}
		env.PressKey(tea.KeyDown)
	}
	return ctx, nil
}
