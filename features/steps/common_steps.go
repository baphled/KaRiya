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
