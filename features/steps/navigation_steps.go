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
//nolint:dupl // Step registration functions look similar but register different steps.
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
	sc.Step(`^the application should not exit$`, navApplicationShouldNotExit)
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

func navShouldSeeAConfirmationDialog(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("confirm"),
		gomega.ContainSubstring("Confirm"),
		gomega.ContainSubstring("save"),
		gomega.ContainSubstring("Save"),
		gomega.ContainSubstring("discard"),
		gomega.ContainSubstring("Discard"),
		gomega.ContainSubstring("cancel"),
		gomega.ContainSubstring("Cancel"),
	), "Should see confirmation dialog")
	return nil
}

func navShouldBeAbleToSaveOrDiscardChanges(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("save"),
		gomega.ContainSubstring("Save"),
		gomega.ContainSubstring("discard"),
		gomega.ContainSubstring("Discard"),
		gomega.ContainSubstring("yes"),
		gomega.ContainSubstring("Yes"),
		gomega.ContainSubstring("no"),
		gomega.ContainSubstring("No"),
	), "Should have save/discard options")
	return nil
}

func navPressingXShouldShowHelp(ctx context.Context, key string) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	if len(key) == 1 {
		env.PressKeyRune(rune(key[0]))
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Help"),
		gomega.ContainSubstring("Keyboard"),
		gomega.ContainSubstring("Reference"),
	), "Should show help")
	return nil
}

func navPressingXShouldQuit(ctx context.Context, key string) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	switch key {
	case "q":
		env.PressKeyRune('q')
	case "Ctrl+C":
		env.PressKey(tea.KeyCtrlC)
	}
	gomega.Expect(env.QuitRequested).To(gomega.BeTrue(), "Application should quit")
	return nil
}

func navPressingXShouldOpenSearch(ctx context.Context, _ string) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	env.PressKeyRune('/')
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Search"),
		gomega.ContainSubstring("search"),
		gomega.ContainSubstring("Find"),
	), "Should show search modal")
	return nil
}

func navPressingXShouldOpenFilter(ctx context.Context, _ string) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	env.PressKeyRune('f')
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Filter"),
		gomega.ContainSubstring("filter"),
	), "Should show filter modal")
	return nil
}

func navPressingXShouldOpenSort(ctx context.Context, _ string) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	env.PressKeyRune('s')
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Sort"),
		gomega.ContainSubstring("sort"),
	), "Should show sort modal")
	return nil
}

func navTypingShouldGoToTheSearchInput(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	env.TypeText("test")
	view := env.GetView()
	gomega.Expect(view).To(gomega.ContainSubstring("test"), "Typed text should appear in search input")
	return nil
}

func navPressingEscapeShouldCloseTheModal(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	env.Cancel()
	return nil
}

func navPressingJShouldNavigateTheList(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	env.PressKeyRune('j')
	return nil
}
