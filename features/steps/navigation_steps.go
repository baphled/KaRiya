// Package steps provides BDD step definitions for Godog feature tests.
package steps

import (
	"context"

	"github.com/baphled/kariya/features/support"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cucumber/godog"
	"github.com/onsi/gomega"
)

// RegisterNavigationSteps registers navigation step definitions with Godog.
//
// Expected: sc is a valid ScenarioContext.
// Returns: None.
// Side effects: Registers step definitions with the scenario context.
//
// Expected:
//   - sc is a valid *godog.ScenarioContext.
//
// Side effects:
//   - Registers step definitions with Godog.
func RegisterNavigationSteps(sc *godog.ScenarioContext) {
	sc.Step(`^I start the application$`, navStartTheApplication)
	sc.Step(`^I should see the main menu$`, navShouldSeeTheMainMenu)
	sc.Step(`^I should not be on the main menu$`, navShouldNotBeOnTheMainMenu)
	sc.Step(`^I press "\?" for help$`, navPressQuestionForHelp)
	sc.Step(`^I should see help information$`, navShouldSeeHelpInformation)
	sc.Step(`^I should see "([^"]*)" or "([^"]*)"$`, navShouldSeeOr)
	sc.Step(`^the application should exit$`, navApplicationShouldExit)
	sc.Step(`^the application should not exit$`, navApplicationShouldNotExit)
	sc.Step(`^I press Ctrl\+C$`, navPressCtrlC)
	sc.Step(`^I am in the middle of capturing an event$`, navAmInTheMiddleOfCapturingAnEvent)
	sc.Step(`^I select "([^"]*)" from the menu$`, navSelectFeatureFromMenu)
}

func navStartTheApplication(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	return ctx, nil
}

func navShouldSeeTheMainMenu(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	gomega.Expect(env.IsInMenuState()).To(gomega.BeTrue())
	return nil
}

func navShouldNotBeOnTheMainMenu(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	gomega.Expect(env.IsInMenuState()).To(gomega.BeFalse())
	return nil
}

func navPressQuestionForHelp(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('?')
	return ctx, nil
}

func navShouldSeeHelpInformation(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Help"),
		gomega.ContainSubstring("help"),
		gomega.ContainSubstring("Keyboard"),
		gomega.ContainSubstring("Shortcuts"),
	))
	return nil
}

func navShouldSeeOr(ctx context.Context, text1, text2 string) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring(text1),
		gomega.ContainSubstring(text2),
	))
	return nil
}

func navApplicationShouldExit(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	gomega.Expect(env.QuitRequested).To(gomega.BeTrue(), "Application should have requested quit")
	return nil
}

func navApplicationShouldNotExit(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	gomega.Expect(env.QuitRequested).To(gomega.BeFalse(), "Application should NOT have requested quit")
	return nil
}

func navPressCtrlC(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKey(tea.KeyCtrlC)
	return ctx, nil
}

func navAmInTheMiddleOfCapturingAnEvent(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.SelectIntentByName("capture_event")
	return ctx, nil
}

func navSelectFeatureFromMenu(ctx context.Context, feature string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}

	env.SelectIntentByName(feature)
	return ctx, nil
}
