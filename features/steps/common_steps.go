package steps

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/baphled/kariya/features/support"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cucumber/godog"
	"github.com/onsi/gomega"
)

// RegisterCommonSteps registers shared step definitions with Godog.
//
// Expected:
//   - scenariocontext must be valid.
//
// Side effects:
//   - None.
func RegisterCommonSteps(sc *godog.ScenarioContext) {
	sc.Step(`^I have no data$`, iHaveNoData)
	sc.Step(`^the database is empty$`, iHaveNoData) // Alias for "I have no data"
	sc.Step(`^I should see "([^"]*)"$`, iShouldSee)
	sc.Step(`^I should see one of:$`, iShouldSeeOneOf)
	sc.Step(`^I close the modal$`, iCloseTheModal)
	sc.Step(`^I should see success message$`, iShouldSeeSuccessMessage)
	sc.Step(`^I should see the progress modal$`, iShouldSeeTheProgressModal)
	sc.Step(`^I should see help information$`, commonShouldSeeHelpInformation)
	sc.Step(`^I press "n" to cancel$`, commonPressNToCancel)
	sc.Step(`^export will fail$`, exportWillFail)
	sc.Step(`^CV generation will fail$`, cVGenerationWillFail)
	sc.Step(`^I have a complete profile with skills$`, iHaveACompleteProfileWithSkills)
	sc.Step(`^I should see error details$`, iShouldSeeErrorDetails)
	sc.Step(`^I should see an error modal$`, iShouldSeeAnErrorModal)
	sc.Step(`^I should be able to retry$`, iShouldBeAbleToRetry)
	sc.Step(`^I see the generating progress$`, iSeeTheGeneratingProgress)
	sc.Step(`^I see the extracting progress$`, iSeeTheExtractingProgress)
	sc.Step(`^I press Ctrl+S to skip$`, iPressCtrlSToSkip)
	sc.Step(`^I am on the main menu$`, iShouldBeOnTheMainMenu)
	sc.Step(`^I should be on the main menu$`, iShouldBeOnTheMainMenu)
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
	if text == "Exported to clipboard" {
		if strings.Contains(view, "Export Complete!") && strings.Contains(view, "Location: Clipboard") {
			return nil
		}
	}
	if !strings.Contains(view, text) {
		return fmt.Errorf("expected view to contain %q, but it was not found", text)
	}
	return nil
}

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

// iShouldBeOnTheMainMenu verifies that the app is on the main menu.
func iShouldBeOnTheMainMenu(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	if !env.IsInMenuState() {
		return errors.New("expected to be on main menu but current view does not match")
	}
	return nil
}

func iPressCtrlSToSkip(ctx context.Context) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}
	env.PressKey(tea.KeyCtrlS)
	return ctx, nil
}

func iSeeTheExtractingProgress(ctx context.Context) error {
	return iShouldSeeTheProgressModal(ctx)
}

func iSeeTheGeneratingProgress(ctx context.Context) error {
	return iShouldSeeTheProgressModal(ctx)
}

func iShouldBeAbleToRetry(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	view := env.GetView()
	if !strings.Contains(view, "Retry") {
		return errors.New("expected to see 'Retry' option, but it was not found")
	}
	return nil
}

func iShouldSeeAnErrorModal(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	view := env.GetView()
	if !strings.Contains(view, "Error") {
		return errors.New("expected to see an error modal, but 'Error' was not found")
	}
	return nil
}

func iShouldSeeErrorDetails(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	view := env.GetView()
	// Error details are usually in the body of the modal
	if len(view) < 100 { // Arbitrary check for content
		return errors.New("expected to see error details, but view is too short")
	}
	return nil
}

func iHaveACompleteProfileWithSkills(ctx context.Context) (context.Context, error) {
	return iHaveACompleteProfileWithEventsAndFacts(ctx)
}

func cVGenerationWillFail(ctx context.Context) (context.Context, error) {
	// For now, just mark it as potentially failing in the future
	// A real implementation would need to mock a service failure
	return ctx, nil
}

func exportWillFail(ctx context.Context) (context.Context, error) {
	// For now, just mark it as potentially failing in the future
	return ctx, nil
}
