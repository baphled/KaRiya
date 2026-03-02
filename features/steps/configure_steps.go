// Package steps provides BDD step definitions for Godog feature tests.
package steps

import (
	"context"

	"github.com/baphled/kariya/features/support"
	"github.com/cucumber/godog"
	"github.com/onsi/gomega"
)

// RegisterConfigureSteps registers configure system step definitions with Godog.
//
// Expected:
//   - sc is a valid *godog.ScenarioContext.
//
// Side effects:
//   - Registers step definitions with Godog.
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
	sectionOrder := []string{"System", "Profile", "Export", "UI"}
	for i, s := range sectionOrder {
		if s == section {
			for range i {
				env.NavigateDown()
			}
			break
		}
	}
	return ctx, nil
}
