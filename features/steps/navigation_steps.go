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
func RegisterNavigationSteps(sc *godog.ScenarioContext) {
	sc.Step(`^I start the application$`, navStartTheApplication)
	sc.Step(`^I should see the main menu$`, navShouldSeeTheMainMenu)
	sc.Step(`^I press "j" to move down$`, navPressJToMoveDown)
	sc.Step(`^I press "k" to move up$`, navPressKToMoveUp)
	sc.Step(`^I should still be on the main menu$`, navShouldStillBeOnTheMainMenu)
	sc.Step(`^I should not be on the main menu$`, navShouldNotBeOnTheMainMenu)
	sc.Step(`^I press "\?" for help$`, navPressQuestionForHelp)
	sc.Step(`^I should see help information$`, navShouldSeeHelpInformation)
	sc.Step(`^I should see "([^"]*)" or "([^"]*)"$`, navShouldSeeOr)
	sc.Step(`^the application should exit$`, navApplicationShouldExit)
	sc.Step(`^I press Ctrl\+C$`, navPressCtrlC)
	sc.Step(`^I am in the middle of capturing an event$`, navAmInTheMiddleOfCapturingAnEvent)
	sc.Step(`^I should see a confirmation dialog$`, navShouldSeeAConfirmationDialog)
	sc.Step(`^I should be able to save or discard changes$`, navShouldBeAbleToSaveOrDiscardChanges)
	sc.Step(`^pressing "([^"]*)" should show help$`, navPressingXShouldShowHelp)
	sc.Step(`^pressing "([^"]*)" should quit$`, navPressingXShouldQuit)
	sc.Step(`^pressing "([^"]*)" should open search$`, navPressingXShouldOpenSearch)
	sc.Step(`^pressing "([^"]*)" should open filter$`, navPressingXShouldOpenFilter)
	sc.Step(`^pressing "([^"]*)" should open sort$`, navPressingXShouldOpenSort)
	sc.Step(`^typing should go to the search input$`, navTypingShouldGoToTheSearchInput)
	sc.Step(`^pressing escape should close the modal$`, navPressingEscapeShouldCloseTheModal)
	sc.Step(`^pressing "j" should navigate the list$`, navPressingJShouldNavigateTheList)
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

func navPressJToMoveDown(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('j')
	return ctx, nil
}

func navPressKToMoveUp(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('k')
	return ctx, nil
}

func navShouldStillBeOnTheMainMenu(ctx context.Context) error {
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
	gomega.Expect(env.IsInMenuState()).To(gomega.BeFalse())
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

func navAmInTheMiddleOfCapturingAnEvent(_ context.Context) (context.Context, error) {
	return nil, godog.ErrPending
}

func navShouldSeeAConfirmationDialog(_ context.Context) error {
	return godog.ErrPending
}

func navShouldBeAbleToSaveOrDiscardChanges(_ context.Context) error {
	return godog.ErrPending
}

func navPressingXShouldShowHelp(_ context.Context, _ string) error {
	return godog.ErrPending
}

func navPressingXShouldQuit(_ context.Context, _ string) error {
	return godog.ErrPending
}

func navPressingXShouldOpenSearch(_ context.Context, _ string) error {
	return godog.ErrPending
}

func navPressingXShouldOpenFilter(_ context.Context, _ string) error {
	return godog.ErrPending
}

func navPressingXShouldOpenSort(_ context.Context, _ string) error {
	return godog.ErrPending
}

func navTypingShouldGoToTheSearchInput(_ context.Context) error {
	return godog.ErrPending
}

func navPressingEscapeShouldCloseTheModal(_ context.Context) error {
	return godog.ErrPending
}

func navPressingJShouldNavigateTheList(_ context.Context) error {
	return godog.ErrPending
}
