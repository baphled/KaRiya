package steps

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/baphled/kariya/features/support"
	"github.com/cucumber/godog"
	"github.com/onsi/gomega"
)

// RegisterCommonSteps registers shared step definitions with Godog.
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
func RegisterCommonSteps(sc *godog.ScenarioContext) {
	sc.Step(`^I have no data$`, iHaveNoData)
	sc.Step(`^I should see "([^"]*)"$`, iShouldSee)
	sc.Step(`^I should see one of:$`, iShouldSeeOneOf)
	sc.Step(`^I close the modal$`, iCloseTheModal)
	sc.Step(`^I should see success message$`, iShouldSeeSuccessMessage)
	sc.Step(`^I should see the progress modal$`, iShouldSeeTheProgressModal)
	sc.Step(`^I should see help information$`, commonShouldSeeHelpInformation)
	sc.Step(`^I press "n" to cancel$`, commonPressNToCancel)
}

// getViewFromContext returns the current view from either app or onboarding env.
func getViewFromContext(ctx context.Context) (string, bool) {
	if appEnv := support.GetAppEnv(ctx); appEnv != nil {
		return appEnv.GetView(), true
	}
	if onboardingEnv := support.GetOnboardingEnv(ctx); onboardingEnv != nil {
		return onboardingEnv.View(), true
	}
	return "", false
}

// iShouldSee asserts that the view contains the given text.
func iShouldSee(ctx context.Context, text string) error {
	view, ok := getViewFromContext(ctx)
	if !ok {
		return errors.New("no test environment in context: check scenario tags and BeforeScenario hook")
	}
	if !strings.Contains(view, text) {
		return fmt.Errorf("expected view to contain %q, but it was not found", text)
	}
	return nil
}

// iShouldNotSee asserts that the view does not contain the given text.

// iShouldSeeOneOf asserts that the view contains at least one of the given texts.
func iShouldSeeOneOf(ctx context.Context, table *godog.Table) error {
	view, ok := getViewFromContext(ctx)
	if !ok {
		return errors.New("no test environment in context: check scenario tags and BeforeScenario hook")
	}

	for _, row := range table.Rows {
		if len(row.Cells) > 0 {
			text := row.Cells[0].Value
			if strings.Contains(view, text) {
				return nil
			}
		}
	}

	texts := make([]string, 0, len(table.Rows))
	for _, row := range table.Rows {
		if len(row.Cells) > 0 {
			texts = append(texts, row.Cells[0].Value)
		}
	}
	gomega.Expect(view).To(gomega.ContainSubstring(texts[0]), "View should contain at least one of: %v", texts)
	return nil
}

// iCloseTheModal closes the currently open modal.
func iCloseTheModal(ctx context.Context) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}
	env.Cancel()
	return ctx, nil
}

// iShouldSeeSuccessMessage asserts a success message is visible.
func iShouldSeeSuccessMessage(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Success"),
		gomega.ContainSubstring("success"),
		gomega.ContainSubstring("saved"),
		gomega.ContainSubstring("Saved"),
	))
	return nil
}

// iShouldSeeTheProgressModal asserts a progress modal is visible.
func iShouldSeeTheProgressModal(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Progress"),
		gomega.ContainSubstring("Loading"),
		gomega.ContainSubstring("Processing"),
		gomega.ContainSubstring("Generating"),
		gomega.ContainSubstring("Extracting"),
		gomega.ContainSubstring("Exporting"),
	))
	return nil
}

// iShouldSeeDifferentContent asserts the view has changed.

// iHaveNoData asserts that the system has no data (empty state).
func iHaveNoData(ctx context.Context) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}
	env.AssertEventCount(0)
	return ctx, nil
}

// commonShouldSeeHelpInformation asserts that help information is visible.
func commonShouldSeeHelpInformation(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
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

// commonPressNToCancel presses "n" to cancel.
func commonPressNToCancel(ctx context.Context) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}
	env.PressKeyRune('n')
	return ctx, nil
}
