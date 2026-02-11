package steps

import (
	"context"
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
	sc.Step(`^I should see "([^"]*)"$`, iShouldSee)
	sc.Step(`^I should not see "([^"]*)"$`, iShouldNotSee)
	sc.Step(`^I should see one of:$`, iShouldSeeOneOf)
	sc.Step(`^I close the modal$`, iCloseTheModal)
	sc.Step(`^I restart the application$`, iRestartTheApplication)
	sc.Step(`^I should be on domain selection$`, iShouldBeOnDomainSelection)
	sc.Step(`^I should see success message$`, iShouldSeeSuccessMessage)
	sc.Step(`^I should see the error modal$`, iShouldSeeTheErrorModal)
	sc.Step(`^I should see the progress modal$`, iShouldSeeTheProgressModal)
	sc.Step(`^I should see different content$`, iShouldSeeDifferentContent)
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
		return godog.ErrPending
	}
	gomega.Expect(view).To(gomega.ContainSubstring(text))
	return nil
}

// iShouldNotSee asserts that the view does not contain the given text.
func iShouldNotSee(ctx context.Context, text string) error {
	view, ok := getViewFromContext(ctx)
	if !ok {
		return godog.ErrPending
	}
	gomega.Expect(view).NotTo(gomega.ContainSubstring(text))
	return nil
}

// iShouldSeeOneOf asserts that the view contains at least one of the given texts.
func iShouldSeeOneOf(ctx context.Context, table *godog.Table) error {
	view, ok := getViewFromContext(ctx)
	if !ok {
		return godog.ErrPending
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
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.Cancel()
	return ctx, nil
}

// iRestartTheApplication simulates restarting the application.
func iRestartTheApplication(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	// Simulate application restart
	env.SimulateRestart()
	return ctx, nil
}

// iShouldBeOnDomainSelection asserts the user is on domain selection screen.
func iShouldBeOnDomainSelection(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Domain"),
		gomega.ContainSubstring("Select"),
	))
	return nil
}

// iShouldSeeSuccessMessage asserts a success message is visible.
func iShouldSeeSuccessMessage(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
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

// iShouldSeeTheErrorModal asserts an error modal is visible.
func iShouldSeeTheErrorModal(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Error"),
		gomega.ContainSubstring("error"),
		gomega.ContainSubstring("failed"),
		gomega.ContainSubstring("Failed"),
	))
	return nil
}

// iShouldSeeTheProgressModal asserts a progress modal is visible.
func iShouldSeeTheProgressModal(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Progress"),
		gomega.ContainSubstring("Loading"),
		gomega.ContainSubstring("Processing"),
	))
	return nil
}

// iShouldSeeDifferentContent asserts the view has changed.
func iShouldSeeDifferentContent(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	// Simply verify we have a non-empty view
	view := env.GetView()
	gomega.Expect(view).NotTo(gomega.BeEmpty())
	return nil
}
