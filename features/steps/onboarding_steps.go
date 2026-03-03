package steps

import (
	"context"
	"errors"
	"fmt"

	"github.com/baphled/kariya/features/support"
	"github.com/cucumber/godog"
	"github.com/onsi/gomega"
)

// RegisterOnboardingSteps registers onboarding-specific step definitions with Godog.
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
func RegisterOnboardingSteps(sc *godog.ScenarioContext) {
	registerOnboardingInputSteps(sc)
	registerOnboardingNavigationSteps(sc)
	registerOnboardingAssertionSteps(sc)
}

func registerOnboardingInputSteps(sc *godog.ScenarioContext) {
	sc.Step(`^I start the onboarding wizard$`, iStartTheOnboardingWizard)
	sc.Step(`^I enter "([^"]*)" as my name$`, iEnterAsMyName)
	sc.Step(`^I enter "([^"]*)" as my email$`, iEnterAsMyEmail)
	sc.Step(`^I enter "([^"]*)" as my location$`, iEnterAsMyLocation)
	sc.Step(`^I enter "([^"]*)" as my title$`, iEnterAsMyTitle)
	sc.Step(`^I enter "([^"]*)" as my GitHub username$`, iEnterAsMyGitHubUsername)
	sc.Step(`^I enter "([^"]*)" as my portfolio$`, iEnterAsMyPortfolio)
}

func registerOnboardingNavigationSteps(sc *godog.ScenarioContext) {
	sc.Step(`^I press enter$`, iPressEnter)
	sc.Step(`^I press tab$`, iPressTab)
	sc.Step(`^I skip the optional fields$`, iSkipOptionalFields)
}

func registerOnboardingAssertionSteps(sc *godog.ScenarioContext) {
	sc.Step(`^the onboarding wizard should be complete$`, onboardingWizardShouldBeComplete)
	sc.Step(`^my profile should have name "([^"]*)"$`, profileShouldHaveName)
	sc.Step(`^my profile should have email "([^"]*)"$`, profileShouldHaveEmail)
	sc.Step(`^my profile should have location "([^"]*)"$`, profileShouldHaveLocation)
	sc.Step(`^my profile should have title "([^"]*)"$`, profileShouldHaveTitle)
	sc.Step(`^my profile should have GitHub username "([^"]*)"$`, profileShouldHaveGitHubUsername)
	sc.Step(`^my profile should have portfolio "([^"]*)"$`, profileShouldHavePortfolio)
	sc.Step(`^my profile should not have a title$`, profileShouldNotHaveTitle)
	sc.Step(`^I should still be on step (\d+)$`, iShouldStillBeOnStep)
	sc.Step(`^I should see a validation error$`, iShouldSeeValidationErrorOnboarding)
	sc.Step(`^I should still see the onboarding wizard$`, iShouldStillSeeOnboardingWizard)
}

// iStartTheOnboardingWizard initializes the onboarding wizard.
func iStartTheOnboardingWizard(ctx context.Context) (context.Context, error) {
	env := support.NewOnboardingEnv()
	return support.WithOnboardingEnv(ctx, env), nil
}

// iAmOnStep verifies the current step number.

// iEnterAsMyName types the name into the name field.
func iEnterAsMyName(ctx context.Context, name string) error {
	env := support.GetOnboardingEnv(ctx)
	if env == nil {
		return errors.New("onboarding environment not initialised: call 'I start the onboarding wizard' first")
	}
	env.TypeText(name)
	return nil
}

// iEnterAsMyEmail types the email into the email field.
func iEnterAsMyEmail(ctx context.Context, email string) error {
	env := support.GetOnboardingEnv(ctx)
	if env == nil {
		return errors.New("onboarding environment not initialised: call 'I start the onboarding wizard' first")
	}
	env.TypeText(email)
	return nil
}

// iEnterAsMyLocation types the location into the location field.
func iEnterAsMyLocation(ctx context.Context, location string) error {
	env := support.GetOnboardingEnv(ctx)
	if env == nil {
		return errors.New("onboarding environment not initialised: call 'I start the onboarding wizard' first")
	}
	env.TypeText(location)
	return nil
}

// iPressEnter presses the enter key.
func iPressEnter(ctx context.Context) error {
	if env := support.GetOnboardingEnv(ctx); env != nil {
		env.PressEnter()
		return nil
	}
	if env, err := support.RequireEnv(ctx); err == nil {
		env.Confirm()
		return nil
	}
	return errors.New("no test environment found (neither onboarding nor full app)")
}

// iPressTab presses the tab key.
func iPressTab(ctx context.Context) error {
	if env := support.GetOnboardingEnv(ctx); env != nil {
		env.PressTab()
		return nil
	}
	if env, err := support.RequireEnv(ctx); err == nil {
		env.Tab()
		return nil
	}
	return errors.New("no test environment found (neither onboarding nor full app)")
}

// iSkipOptionalFields tabs through optional fields and submits.
func iSkipOptionalFields(ctx context.Context) error {
	env := support.GetOnboardingEnv(ctx)
	if env == nil {
		return errors.New("onboarding environment not initialised: call 'I start the onboarding wizard' first")
	}
	env.PressTab()
	env.PressTab()
	env.PressEnter()
	return nil
}

// onboardingWizardShouldBeComplete verifies the wizard is complete.
func onboardingWizardShouldBeComplete(ctx context.Context) error {
	env := support.GetOnboardingEnv(ctx)
	if env == nil {
		return errors.New("onboarding environment not initialised: call 'I start the onboarding wizard' first")
	}
	gomega.Expect(env.IsCompleted()).To(gomega.BeTrue(), "Onboarding wizard should be complete")
	return nil
}

// profileShouldHaveName verifies the profile has the expected name.
func profileShouldHaveName(ctx context.Context, name string) error {
	env := support.GetOnboardingEnv(ctx)
	if env == nil {
		return errors.New("onboarding environment not initialised: call 'I start the onboarding wizard' first")
	}
	gomega.Expect(env.View()).To(gomega.ContainSubstring(name))
	return nil
}

// profileShouldHaveEmail verifies the profile has the expected email.
func profileShouldHaveEmail(ctx context.Context, email string) error {
	env := support.GetOnboardingEnv(ctx)
	if env == nil {
		return errors.New("onboarding environment not initialised: call 'I start the onboarding wizard' first")
	}
	gomega.Expect(env.View()).To(gomega.ContainSubstring(email))
	return nil
}

// profileShouldHaveLocation verifies the profile has the expected location.
func profileShouldHaveLocation(ctx context.Context, location string) error {
	env := support.GetOnboardingEnv(ctx)
	if env == nil {
		return errors.New("onboarding environment not initialised: call 'I start the onboarding wizard' first")
	}
	gomega.Expect(env.View()).To(gomega.ContainSubstring(location))
	return nil
}

// iEnterAsMyTitle types the title into the title field.
func iEnterAsMyTitle(ctx context.Context, title string) error {
	env := support.GetOnboardingEnv(ctx)
	if env == nil {
		return errors.New("onboarding environment not initialised: call 'I start the onboarding wizard' first")
	}
	env.TypeText(title)
	return nil
}

// iEnterAsMyGitHubUsername types the GitHub username into the GitHub field.
func iEnterAsMyGitHubUsername(ctx context.Context, username string) error {
	env := support.GetOnboardingEnv(ctx)
	if env == nil {
		return errors.New("onboarding environment not initialised: call 'I start the onboarding wizard' first")
	}
	env.TypeText(username)
	return nil
}

// iEnterAsMyPortfolio types the portfolio URL into the portfolio field.
func iEnterAsMyPortfolio(ctx context.Context, portfolio string) error {
	env := support.GetOnboardingEnv(ctx)
	if env == nil {
		return errors.New("onboarding environment not initialised: call 'I start the onboarding wizard' first")
	}
	env.TypeText(portfolio)
	return nil
}

// iPressShiftTabOnboarding presses shift+tab to go to previous field in onboarding.

// profileShouldHaveTitle verifies the profile has the expected title.
func profileShouldHaveTitle(ctx context.Context, title string) error {
	env := support.GetOnboardingEnv(ctx)
	if env == nil {
		return errors.New("onboarding environment not initialised: call 'I start the onboarding wizard' first")
	}
	gomega.Expect(env.View()).To(gomega.ContainSubstring(title))
	return nil
}

// profileShouldHaveGitHubUsername verifies the profile has the expected GitHub username.
func profileShouldHaveGitHubUsername(ctx context.Context, username string) error {
	env := support.GetOnboardingEnv(ctx)
	if env == nil {
		return errors.New("onboarding environment not initialised: call 'I start the onboarding wizard' first")
	}
	gomega.Expect(env.View()).To(gomega.ContainSubstring(username))
	return nil
}

// profileShouldHavePortfolio verifies the profile has the expected portfolio URL.
func profileShouldHavePortfolio(ctx context.Context, portfolio string) error {
	env := support.GetOnboardingEnv(ctx)
	if env == nil {
		return errors.New("onboarding environment not initialised: call 'I start the onboarding wizard' first")
	}
	gomega.Expect(env.View()).To(gomega.ContainSubstring(portfolio))
	return nil
}

// profileShouldNotHaveTitle verifies the profile has no title set.
func profileShouldNotHaveTitle(ctx context.Context) error {
	env := support.GetOnboardingEnv(ctx)
	if env == nil {
		return errors.New("onboarding environment not initialised: call 'I start the onboarding wizard' first")
	}
	view := env.View()
	gomega.Expect(view).NotTo(gomega.ContainSubstring("Title:"))
	return nil
}

// iShouldStillBeOnStep verifies we are still on the expected step.
func iShouldStillBeOnStep(ctx context.Context, step int) error {
	env := support.GetOnboardingEnv(ctx)
	if env == nil {
		return errors.New("onboarding environment not initialised: call 'I start the onboarding wizard' first")
	}
	expected := fmt.Sprintf("Step %d of 3", step)
	gomega.Expect(env.View()).To(gomega.ContainSubstring(expected))
	return nil
}

// iShouldSeeValidationErrorOnboarding checks for validation error in the onboarding view.
func iShouldSeeValidationErrorOnboarding(ctx context.Context) error {
	env := support.GetOnboardingEnv(ctx)
	if env == nil {
		return errors.New("onboarding environment not initialised: call 'I start the onboarding wizard' first")
	}
	view := env.View()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("error"),
		gomega.ContainSubstring("Error"),
		gomega.ContainSubstring("required"),
		gomega.ContainSubstring("Required"),
		gomega.ContainSubstring("invalid"),
		gomega.ContainSubstring("Invalid"),
		gomega.ContainSubstring("please"),
		gomega.ContainSubstring("Please"),
	), "Should see validation error")
	return nil
}

// iShouldStillSeeOnboardingWizard verifies the onboarding wizard is visible.
func iShouldStillSeeOnboardingWizard(ctx context.Context) error {
	env := support.GetOnboardingEnv(ctx)
	if env == nil {
		return errors.New("onboarding environment not initialised: call 'I start the onboarding wizard' first")
	}
	gomega.Expect(env.View()).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Step"),
		gomega.ContainSubstring("Profile Setup"),
	))
	return nil
}

// fieldShouldBeFocused verifies the specified field is focused.
