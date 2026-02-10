// Package steps provides BDD step definitions for Godog feature tests.
package steps

import (
	"context"
	"fmt"

	"github.com/baphled/kariya/features/support"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cucumber/godog"
	"github.com/onsi/gomega"
)

// RegisterCVSteps registers CV generation step definitions with Godog.
//
// Expected:
//   - sc is a valid *godog.ScenarioContext.
//
// Side effects:
//   - Registers step definitions with Godog.
func RegisterCVSteps(sc *godog.ScenarioContext) {
	registerCVPrerequisiteSteps(sc)
	registerCVWizardSteps(sc)
	registerCVReviewSteps(sc)
	registerCVExportSteps(sc)
	registerCVErrorSteps(sc)
}

func registerCVPrerequisiteSteps(sc *godog.ScenarioContext) {
	sc.Step(`^I have no profile configured$`, iHaveNoProfileConfigured)
	sc.Step(`^I should see an error or warning$`, iShouldSeeAnErrorOrWarning)
	sc.Step(`^I have a profile configured$`, iHaveAProfileConfigured)
	sc.Step(`^I have a complete profile with events and facts$`, iHaveACompleteProfileWithEventsAndFacts)
	sc.Step(`^I should see the CV wizard modal$`, iShouldSeeTheCVWizardModal)
	sc.Step(`^I have multiple profiles$`, iHaveMultipleProfiles)
}

func registerCVWizardSteps(sc *godog.ScenarioContext) {
	sc.Step(`^I navigate down in the profile selector$`, iNavigateDownInTheProfileSelector)
	sc.Step(`^I confirm selection$`, iConfirmSelection)
	sc.Step(`^I should move to the next field$`, iShouldMoveToTheNextField)
	sc.Step(`^I tab to audience field$`, iTabToAudienceField)
	sc.Step(`^I should move to step 2$`, iShouldMoveToStep2)
	sc.Step(`^I complete step 1$`, iCompleteStep1)
	sc.Step(`^I select "([^"]*)" technology focus$`, iSelectTechnologyFocus)
	sc.Step(`^I should skip technology selection$`, iShouldSkipTechnologySelection)
	sc.Step(`^I should see focus area options$`, iShouldSeeFocusAreaOptions)
	sc.Step(`^I should see technology multi-select$`, iShouldSeeTechnologyMultiSelect)
	sc.Step(`^I should be able to select multiple technologies$`, iShouldBeAbleToSelectMultipleTechnologies)
	sc.Step(`^I should see technology single-select$`, iShouldSeeTechnologySingleSelect)
	sc.Step(`^I should only select one technology$`, iShouldOnlySelectOneTechnology)
	sc.Step(`^I select technology focus$`, iSelectTechFocus)
	sc.Step(`^I should be back on step 1$`, iShouldBeBackOnStep1)
	sc.Step(`^I complete step 2$`, iCompleteStep2)
	sc.Step(`^I tab to skills limit$`, iTabToSkillsLimit)
	sc.Step(`^I enter skills limit "([^"]*)"$`, iEnterSkillsLimit)
	sc.Step(`^the skills limit should be (\d+)$`, theSkillsLimitShouldBe)
	sc.Step(`^I tab to CV length$`, iTabToCVLength)
	sc.Step(`^I press Ctrl\+S to skip$`, iPressCtrlSToSkip)
	sc.Step(`^I have a complete profile with skills$`, iHaveACompleteProfileWithSkills)
	sc.Step(`^I see the extracting progress$`, iSeeTheExtractingProgress)
	sc.Step(`^I complete the wizard$`, iCompleteTheWizard)
	sc.Step(`^I see the generating progress$`, iSeeTheGeneratingProgress)
	sc.Step(`^the generation completes$`, theGenerationCompletes)
}

func registerCVReviewSteps(sc *godog.ScenarioContext) {
	sc.Step(`^I should see the CV review screen$`, iShouldSeeTheCVReviewScreen)
	sc.Step(`^I should see CV metadata$`, iShouldSeeCVMetadata)
	sc.Step(`^I should see statistics$`, iShouldSeeStatistics)
	sc.Step(`^I should see section names$`, iShouldSeeSectionNames)
	sc.Step(`^I should see bullet counts$`, iShouldSeeBulletCounts)
	sc.Step(`^I have generated a CV$`, iHaveGeneratedACV)
	sc.Step(`^I am on the CV review screen$`, iAmOnTheCVReviewScreen)
	sc.Step(`^I press enter to preview$`, iPressEnterToPreview)
	sc.Step(`^I should see the CV preview screen$`, iShouldSeeTheCVPreviewScreen)
	sc.Step(`^I press "p" to preview$`, iPressPToPreview)
	sc.Step(`^I should see the export options modal$`, iShouldSeeTheExportOptionsModal)
	sc.Step(`^I navigate to the CV preview screen$`, iNavigateToTheCVPreviewScreen)
	sc.Step(`^I should see personal details$`, iShouldSeePersonalDetails)
	sc.Step(`^I should see all CV sections$`, iShouldSeeAllCVSections)
	sc.Step(`^I should see bullet points$`, iShouldSeeBulletPoints)
	sc.Step(`^I have generated a long CV$`, iHaveGeneratedALongCV)
	sc.Step(`^I should scroll a full page$`, iShouldScrollAFullPage)
	sc.Step(`^I should scroll back$`, iShouldScrollBack)
	sc.Step(`^I should be at the bottom$`, iShouldBeAtTheBottom)
	sc.Step(`^I should be at the top$`, iShouldBeAtTheTop)
	sc.Step(`^I press enter to confirm$`, iPressEnterToConfirmCV)
	sc.Step(`^the CV generation should complete$`, theCVGenerationShouldComplete)
	sc.Step(`^I press "y" to confirm$`, iPressYToConfirmCV)
}

func registerCVExportSteps(sc *godog.ScenarioContext) {
	sc.Step(`^I open the export options modal$`, iOpenTheExportOptionsModal)
	sc.Step(`^I tab to location$`, iTabToLocation)
	sc.Step(`^I select format "([^"]*)"$`, iSelectFormat)
	sc.Step(`^I select location "([^"]*)"$`, iSelectLocation)
	sc.Step(`^I confirm export$`, iConfirmExport)
	sc.Step(`^I should see export progress$`, iShouldSeeExportProgress)
	sc.Step(`^the export should complete$`, theExportShouldComplete)
	sc.Step(`^I should return to previous screen$`, iShouldReturnToPreviousScreen)
	sc.Step(`^I start an export$`, iStartAnExport)
	sc.Step(`^I complete an export$`, iCompleteAnExport)
	sc.Step(`^I should see export location$`, iShouldSeeExportLocation)
}

func registerCVErrorSteps(sc *godog.ScenarioContext) {
	sc.Step(`^CV generation will fail$`, cvGenerationWillFail)
	sc.Step(`^I should see an error modal$`, iShouldSeeAnErrorModal)
	sc.Step(`^I should see error details$`, iShouldSeeErrorDetails)
	sc.Step(`^export will fail$`, exportWillFail)
	sc.Step(`^I should be able to retry$`, iShouldBeAbleToRetry)
}

func iHaveNoProfileConfigured(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	// Profile is empty by default in test env
	return ctx, nil
}

func iShouldSeeAnErrorOrWarning(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Error"),
		gomega.ContainSubstring("error"),
		gomega.ContainSubstring("Warning"),
		gomega.ContainSubstring("warning"),
		gomega.ContainSubstring("required"),
		gomega.ContainSubstring("need"),
	))
	return nil
}

func iHaveAProfileConfigured(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	// Profile setup handled by config system - assume configured
	return ctx, nil
}

func iHaveACompleteProfileWithEventsAndFacts(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	// Add events
	for i := range 3 {
		event := fixtures.EventWith("", fmt.Sprintf("Event %d", i+1), fmt.Sprintf("Company%d", i+1), "")
		env.AddEvent(event)
	}
	// Add facts
	for i := range 3 {
		fact := fixtures.FactWith("", fmt.Sprintf("Fact %d text", i+1))
		env.AddFact(fact)
	}
	return ctx, nil
}

func iShouldSeeTheCVWizardModal(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("CV"),
		gomega.ContainSubstring("Profile"),
		gomega.ContainSubstring("Audience"),
		gomega.ContainSubstring("Generate"),
	))
	return nil
}

func iHaveMultipleProfiles(ctx context.Context) (context.Context, error) {
	// Multiple profiles not supported in current implementation
	// Fall back to single profile
	return iHaveAProfileConfigured(ctx)
}

func iNavigateDownInTheProfileSelector(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.NavigateDown()
	return ctx, nil
}

func iConfirmSelection(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.Confirm()
	return ctx, nil
}

func iShouldMoveToTheNextField(_ context.Context) error {
	return godog.ErrPending
}

func iTabToAudienceField(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.Tab()
	return ctx, nil
}

func iShouldMoveToStep2(_ context.Context) error {
	return godog.ErrPending
}

func iCompleteStep1(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.Confirm()
	env.Tab()
	env.Confirm()
	return ctx, nil
}

func iSelectTechnologyFocus(_ context.Context, _ string) (context.Context, error) {
	return nil, godog.ErrPending
}

func iShouldSkipTechnologySelection(_ context.Context) error {
	return godog.ErrPending
}

func iShouldSeeFocusAreaOptions(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Backend"),
		gomega.ContainSubstring("Frontend"),
		gomega.ContainSubstring("Fullstack"),
		gomega.ContainSubstring("DevOps"),
	))
	return nil
}

func iShouldSeeTechnologyMultiSelect(_ context.Context) error {
	return godog.ErrPending
}

func iShouldBeAbleToSelectMultipleTechnologies(_ context.Context) error {
	return godog.ErrPending
}

func iShouldSeeTechnologySingleSelect(_ context.Context) error {
	return godog.ErrPending
}

func iShouldOnlySelectOneTechnology(_ context.Context) error {
	return godog.ErrPending
}

func iSelectTechFocus(_ context.Context) (context.Context, error) {
	return nil, godog.ErrPending
}

func iShouldBeBackOnStep1(_ context.Context) error {
	return godog.ErrPending
}

func iCompleteStep2(_ context.Context) (context.Context, error) {
	return nil, godog.ErrPending
}

func iTabToSkillsLimit(_ context.Context) (context.Context, error) {
	return nil, godog.ErrPending
}

func iEnterSkillsLimit(_ context.Context, _ string) (context.Context, error) {
	return nil, godog.ErrPending
}

func theSkillsLimitShouldBe(_ context.Context, _ int) error {
	return godog.ErrPending
}

func iTabToCVLength(_ context.Context) (context.Context, error) {
	return nil, godog.ErrPending
}

func iPressCtrlSToSkip(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKey(tea.KeyCtrlS)
	return ctx, nil
}

func iHaveACompleteProfileWithSkills(_ context.Context) (context.Context, error) {
	return nil, godog.ErrPending
}

func iSeeTheExtractingProgress(_ context.Context) error {
	return godog.ErrPending
}

func iCompleteTheWizard(_ context.Context) (context.Context, error) {
	return nil, godog.ErrPending
}

func iSeeTheGeneratingProgress(_ context.Context) error {
	return godog.ErrPending
}

func theGenerationCompletes(_ context.Context) error {
	return godog.ErrPending
}

func iShouldSeeTheCVReviewScreen(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Review"),
		gomega.ContainSubstring("CV"),
		gomega.ContainSubstring("Metadata"),
	))
	return nil
}

func iShouldSeeCVMetadata(_ context.Context) error {
	return godog.ErrPending
}

func iShouldSeeStatistics(_ context.Context) error {
	return godog.ErrPending
}

func iShouldSeeSectionNames(_ context.Context) error {
	return godog.ErrPending
}

func iShouldSeeBulletCounts(_ context.Context) error {
	return godog.ErrPending
}

func iHaveGeneratedACV(_ context.Context) (context.Context, error) {
	return nil, godog.ErrPending
}

func iAmOnTheCVReviewScreen(_ context.Context) error {
	return godog.ErrPending
}

func iPressEnterToPreview(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.Confirm()
	return ctx, nil
}

func iShouldSeeTheCVPreviewScreen(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Preview"),
		gomega.ContainSubstring("CV"),
	))
	return nil
}

func iPressPToPreview(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('p')
	return ctx, nil
}

func iShouldSeeTheExportOptionsModal(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Export"),
		gomega.ContainSubstring("Format"),
		gomega.ContainSubstring("Location"),
	))
	return nil
}

func iNavigateToTheCVPreviewScreen(_ context.Context) (context.Context, error) {
	return nil, godog.ErrPending
}

func iShouldSeePersonalDetails(_ context.Context) error {
	return godog.ErrPending
}

func iShouldSeeAllCVSections(_ context.Context) error {
	return godog.ErrPending
}

func iShouldSeeBulletPoints(_ context.Context) error {
	return godog.ErrPending
}

func iHaveGeneratedALongCV(_ context.Context) (context.Context, error) {
	return nil, godog.ErrPending
}

func iShouldScrollAFullPage(_ context.Context) error {
	return godog.ErrPending
}

func iShouldScrollBack(_ context.Context) error {
	return godog.ErrPending
}

func iShouldBeAtTheBottom(_ context.Context) error {
	return godog.ErrPending
}

func iShouldBeAtTheTop(_ context.Context) error {
	return godog.ErrPending
}

func iPressEnterToConfirmCV(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.Confirm()
	return ctx, nil
}

func theCVGenerationShouldComplete(_ context.Context) error {
	return godog.ErrPending
}

func iPressYToConfirmCV(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('y')
	return ctx, nil
}

func iOpenTheExportOptionsModal(_ context.Context) (context.Context, error) {
	return nil, godog.ErrPending
}

func iTabToLocation(_ context.Context) (context.Context, error) {
	return nil, godog.ErrPending
}

func iSelectFormat(_ context.Context, _ string) (context.Context, error) {
	return nil, godog.ErrPending
}

func iSelectLocation(_ context.Context, _ string) (context.Context, error) {
	return nil, godog.ErrPending
}

func iConfirmExport(_ context.Context) (context.Context, error) {
	return nil, godog.ErrPending
}

func iShouldSeeExportProgress(_ context.Context) error {
	return godog.ErrPending
}

func theExportShouldComplete(_ context.Context) error {
	return godog.ErrPending
}

func iShouldReturnToPreviousScreen(_ context.Context) error {
	return godog.ErrPending
}

func iStartAnExport(_ context.Context) (context.Context, error) {
	return nil, godog.ErrPending
}

func iCompleteAnExport(_ context.Context) error {
	return godog.ErrPending
}

func iShouldSeeExportLocation(_ context.Context) error {
	return godog.ErrPending
}

func cvGenerationWillFail(_ context.Context) (context.Context, error) {
	return nil, godog.ErrPending
}

func iShouldSeeAnErrorModal(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Error"),
		gomega.ContainSubstring("error"),
		gomega.ContainSubstring("Failed"),
		gomega.ContainSubstring("failed"),
	))
	return nil
}

func iShouldSeeErrorDetails(_ context.Context) error {
	return godog.ErrPending
}

func exportWillFail(_ context.Context) (context.Context, error) {
	return nil, godog.ErrPending
}

func iShouldBeAbleToRetry(_ context.Context) error {
	return godog.ErrPending
}
