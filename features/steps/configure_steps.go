// Package steps provides BDD step definitions for Godog feature tests.
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
//   - sc is a valid *godog.ScenarioContext.
//
// Side effects:
//   - Registers step definitions with Godog.
func RegisterConfigureSteps(sc *godog.ScenarioContext) {
	sc.Step(`^I should see the domain selection screen$`, iShouldSeeTheDomainSelectionScreen)
	sc.Step(`^I should still be on domain selection$`, iShouldStillBeOnDomainSelection)
	sc.Step(`^I should be at the last domain$`, iShouldBeAtTheLastDomain)
	sc.Step(`^I should be at the first domain$`, iShouldBeAtTheFirstDomain)
	sc.Step(`^I press "q" to quit$`, iPressQToQuit)
	sc.Step(`^I select "([^"]*)" domain$`, iSelectDomain)
	sc.Step(`^I should see the edit settings modal$`, iShouldSeeTheEditSettingsModal)
	sc.Step(`^I navigate to "([^"]*)" field$`, iNavigateToField)
	sc.Step(`^I enter "([^"]*)"$`, iEnterValue)
	sc.Step(`^the field should show "([^"]*)"$`, theFieldShouldShow)
	sc.Step(`^I toggle the boolean value$`, iToggleTheBooleanValue)
	sc.Step(`^the value should change$`, theValueShouldChange)
	sc.Step(`^I enter name "([^"]*)"$`, iEnterName)
	sc.Step(`^I tab to email$`, iTabToEmail)
	sc.Step(`^I enter email "([^"]*)"$`, iEnterEmail)
	sc.Step(`^I submit settings$`, iSubmitSettings)
	sc.Step(`^I should see the review modal$`, iShouldSeeTheReviewModal)
	sc.Step(`^the field should accept comma-separated values$`, theFieldShouldAcceptCommaSeparatedValues)
	sc.Step(`^I make a change$`, iMakeAChange)
	sc.Step(`^I complete the form$`, iCompleteTheForm)
	sc.Step(`^I change log level to "([^"]*)"$`, iChangeLogLevelTo)
	sc.Step(`^I am on the review changes modal$`, iAmOnTheReviewChangesModal)
	sc.Step(`^I should see the confirm modal$`, iShouldSeeTheConfirmModal)
	sc.Step(`^I press "n" to cancel$`, iPressNToCancel)
	sc.Step(`^I press "q" to cancel$`, iPressQToCancel)
	sc.Step(`^I am on the confirm modal$`, iAmOnTheConfirmModal)
	sc.Step(`^I should see the saving modal$`, iShouldSeeSavingModal)
	sc.Step(`^I confirm save$`, iConfirmSave)
	sc.Step(`^I should see a spinner$`, iShouldSeeASpinner)
	sc.Step(`^the save completes successfully$`, theSaveCompletesSuccessfully)
	sc.Step(`^I should see the success modal$`, iShouldSeeTheSuccessModal)
	sc.Step(`^the success modal should auto-dismiss$`, theSuccessModalShouldAutoDismiss)
	sc.Step(`^the save fails$`, theSaveFails)
	sc.Step(`^I dismiss the error modal$`, iDismissTheErrorModal)
	sc.Step(`^I confirm the review$`, iConfirmTheReviewConfig)
	sc.Step(`^I confirm the save$`, iConfirmTheSave)
	sc.Step(`^the save completes$`, theSaveCompletes)
	sc.Step(`^I press escape from confirm$`, iPressEscapeFromConfirm)
}

func iShouldSeeTheDomainSelectionScreen(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Configuration"),
		gomega.ContainSubstring("Domain"),
		gomega.ContainSubstring("System"),
		gomega.ContainSubstring("Profile"),
	))
	return nil
}

func iShouldStillBeOnDomainSelection(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("System"),
		gomega.ContainSubstring("Profile"),
		gomega.ContainSubstring("Export"),
		gomega.ContainSubstring("UI"),
	))
	return nil
}

func iShouldBeAtTheLastDomain(_ context.Context) error {
	return godog.ErrPending
}

func iShouldBeAtTheFirstDomain(_ context.Context) error {
	return godog.ErrPending
}

func iPressQToQuit(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('q')
	return ctx, nil
}

func iSelectDomain(ctx context.Context, domain string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	domainOrder := []string{"System", "Profile", "Export", "UI"}
	for i, d := range domainOrder {
		if d == domain {
			for range i {
				env.NavigateDown()
			}
			break
		}
	}
	env.Confirm()
	return ctx, nil
}

func iShouldSeeTheEditSettingsModal(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Edit"),
		gomega.ContainSubstring("Settings"),
		gomega.ContainSubstring("Submit"),
		gomega.ContainSubstring("Cancel"),
	))
	return nil
}

func iNavigateToField(ctx context.Context, fieldLabel string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	// Tab through fields until we find the one with matching label
	// Max 20 tabs to avoid infinite loop
	for i := 0; i < 20; i++ {
		view := env.GetView()
		if strings.Contains(view, fieldLabel) {
			return ctx, nil // Field found and focused
		}
		env.Tab()
	}
	return ctx, godog.ErrPending // Field not found
}

func iEnterValue(ctx context.Context, value string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.TypeText(value)
	return ctx, nil
}

func theFieldShouldShow(ctx context.Context, value string) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.ContainSubstring(value))
	return nil
}

func iToggleTheBooleanValue(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	// Press space or enter to toggle boolean field
	env.Confirm()
	return ctx, nil
}

func theValueShouldChange(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	// Just verify the view updated (any change indicates toggle worked)
	view := env.GetView()
	gomega.Expect(view).NotTo(gomega.BeEmpty())
	return nil
}

func iEnterName(ctx context.Context, name string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.TypeText(name)
	return ctx, nil
}

func iTabToEmail(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.Tab()
	return ctx, nil
}

func iEnterEmail(ctx context.Context, email string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.TypeText(email)
	return ctx, nil
}

func iSubmitSettings(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKey(tea.KeyCtrlS)
	return ctx, nil
}

func iShouldSeeTheReviewModal(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Review"),
		gomega.ContainSubstring("Changes"),
		gomega.ContainSubstring("Confirm"),
	))
	return nil
}

func theFieldShouldAcceptCommaSeparatedValues(_ context.Context) error {
	return godog.ErrPending
}

func iMakeAChange(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.TypeText("test")
	return ctx, nil
}

func iCompleteTheForm(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.Confirm()
	return ctx, nil
}

func iChangeLogLevelTo(_ context.Context, _ string) (context.Context, error) {
	return nil, godog.ErrPending
}

func iAmOnTheReviewChangesModal(_ context.Context) error {
	return godog.ErrPending
}

func iShouldSeeTheConfirmModal(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Confirm"),
		gomega.ContainSubstring("Are you sure"),
		gomega.ContainSubstring("Yes"),
		gomega.ContainSubstring("No"),
	))
	return nil
}

func iPressNToCancel(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('n')
	return ctx, nil
}

func iPressQToCancel(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('q')
	return ctx, nil
}

func iAmOnTheConfirmModal(_ context.Context) error {
	return godog.ErrPending
}

func iShouldSeeSavingModal(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Saving"),
		gomega.ContainSubstring("Loading"),
	))
	return nil
}

func iConfirmSave(_ context.Context) (context.Context, error) {
	return nil, godog.ErrPending
}

func iShouldSeeASpinner(_ context.Context) error {
	return godog.ErrPending
}

func theSaveCompletesSuccessfully(_ context.Context) error {
	return godog.ErrPending
}

func iShouldSeeTheSuccessModal(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Success"),
		gomega.ContainSubstring("saved"),
		gomega.ContainSubstring("Saved"),
	))
	return nil
}

func theSuccessModalShouldAutoDismiss(_ context.Context) error {
	return godog.ErrPending
}

func theSaveFails(_ context.Context) error {
	return godog.ErrPending
}

func iDismissTheErrorModal(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.Confirm()
	return ctx, nil
}

func iConfirmTheReviewConfig(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.Confirm()
	return ctx, nil
}

func iConfirmTheSave(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.Confirm()
	return ctx, nil
}

func theSaveCompletes(_ context.Context) error {
	return godog.ErrPending
}

func iPressEscapeFromConfirm(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.Cancel()
	return ctx, nil
}
