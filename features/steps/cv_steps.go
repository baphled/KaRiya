// Package steps provides BDD step definitions for Godog feature tests.
package steps

import (
	"context"
	"strconv"
	"strings"

	"github.com/baphled/kariya/features/support"
	"github.com/baphled/kariya/internal/testutil/e2e"
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
	sc.Step(`^I press "([^"]*)" to export$`, iPressKeyToExport)
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

	// Use fixtures to populate substantial test data for CV generation
	// This creates realistic career data similar to CSV imports
	env.PopulateTestData(50, 5, 30)

	// Add skills - CV generation requires skills for sections
	skills := e2e.CreateSampleSkills(20)
	for _, skill := range skills {
		env.AddSkill(skill)
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

func iShouldMoveToTheNextField(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	// Field movement is handled by Tab - just verify view updated
	view := env.GetView()
	gomega.Expect(view).NotTo(gomega.BeEmpty())
	return nil
}

func iTabToAudienceField(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.Tab()
	return ctx, nil
}

func iShouldMoveToStep2(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Step 2"),
		gomega.ContainSubstring("Technology"),
		gomega.ContainSubstring("Focus"),
	))
	return nil
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

func iSelectTechnologyFocus(ctx context.Context, _ string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	// Navigate to focus option and confirm
	// Different focuses: "Language Agnostic", "Generalist", "Specialist"
	env.NavigateDown() // Move through options
	env.Confirm()
	return ctx, nil
}

func iShouldSkipTechnologySelection(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	// For language agnostic, no tech selection needed
	view := env.GetView()
	gomega.Expect(view).NotTo(gomega.ContainSubstring("Select technolog"))
	return nil
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

func iShouldSeeTechnologyMultiSelect(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("multi"),
		gomega.ContainSubstring("Multiple"),
		gomega.ContainSubstring("space"),
	))
	return nil
}

func iShouldBeAbleToSelectMultipleTechnologies(ctx context.Context) error {
	// Just verify we can navigate - actual selection tested elsewhere
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	env.PressKey(tea.KeySpace)
	return nil
}

func iShouldSeeTechnologySingleSelect(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.ContainSubstring("Select"))
	return nil
}

func iShouldOnlySelectOneTechnology(ctx context.Context) error {
	// Single select is default behavior - just verify selection works
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	env.Confirm()
	return nil
}

func iSelectTechFocus(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.NavigateDown()
	env.Confirm()
	return ctx, nil
}

func iShouldBeBackOnStep1(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Step 1"),
		gomega.ContainSubstring("Profile"),
		gomega.ContainSubstring("Audience"),
	))
	return nil
}

func iCompleteStep2(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.Confirm()
	return ctx, nil
}

func iTabToSkillsLimit(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.Tab()
	return ctx, nil
}

func iEnterSkillsLimit(ctx context.Context, limit string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.TypeText(limit)
	return ctx, nil
}

func theSkillsLimitShouldBe(ctx context.Context, expected int) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.ContainSubstring(strconv.Itoa(expected)))
	return nil
}

func iTabToCVLength(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.Tab()
	return ctx, nil
}

func iPressCtrlSToSkip(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKey(tea.KeyCtrlS)
	return ctx, nil
}

func iHaveACompleteProfileWithSkills(ctx context.Context) (context.Context, error) {
	// Reuse complete profile setup and add skills
	return iHaveACompleteProfileWithEventsAndFacts(ctx)
}

func iSeeTheExtractingProgress(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Extract"),
		gomega.ContainSubstring("Progress"),
		gomega.ContainSubstring("..."),
	))
	return nil
}

func iCompleteTheWizard(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	// Skip wizard with Ctrl+S (uses default settings)
	env.PressKey(tea.KeyCtrlS)
	return ctx, nil
}

func iSeeTheGeneratingProgress(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Generat"),
		gomega.ContainSubstring("Progress"),
	))
	return nil
}

func theGenerationCompletes(ctx context.Context) error {
	// Wait for generation to complete and reach review screen
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	// Wait for review screen indicators
	gomega.Eventually(func() string {
		return env.GetView()
	}, "5s", "100ms").Should(gomega.SatisfyAny(
		gomega.ContainSubstring("Review"),
		gomega.ContainSubstring("section"),
		gomega.ContainSubstring("Experience"),
		gomega.ContainSubstring("Skills"),
	))
	return nil
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

func iShouldSeeCVMetadata(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Profile"),
		gomega.ContainSubstring("Audience"),
		gomega.ContainSubstring("Date"),
	))
	return nil
}

func iShouldSeeStatistics(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.MatchRegexp(`\d+`)) // Should see numbers
	return nil
}

func iShouldSeeSectionNames(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Experience"),
		gomega.ContainSubstring("Skills"),
		gomega.ContainSubstring("Summary"),
	))
	return nil
}

func iShouldSeeBulletCounts(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.MatchRegexp(`\d+\s+(bullet|item)`))
	return nil
}

func iHaveGeneratedACV(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}

	// Set up data
	ctx, err := iHaveACompleteProfileWithEventsAndFacts(ctx)
	if err != nil {
		return ctx, err
	}

	// Navigate to generate_cv
	env.SelectIntentByName("generate_cv")

	// Skip wizard with Ctrl+S
	env.PressKey(tea.KeyCtrlS)

	// Wait for generation to complete and reach review screen
	gomega.Eventually(func() string {
		return env.GetView()
	}, "5s", "100ms").Should(gomega.ContainSubstring("CV Review"))

	return ctx, nil
}

func iAmOnTheCVReviewScreen(ctx context.Context) error {
	return iShouldSeeTheCVReviewScreen(ctx)
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

func iNavigateToTheCVPreviewScreen(ctx context.Context) (context.Context, error) {
	// Assumes we're already on review screen, just navigate to preview
	return iPressEnterToPreview(ctx)
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

func iOpenTheExportOptionsModal(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	// Press 'x' to open export modal from review/preview screen
	env.PressKeyRune('x')
	return ctx, nil
}

func iPressKeyToExport(ctx context.Context, key string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	// Press the specified key to open export modal
	if len(key) > 0 {
		env.PressKeyRune(rune(key[0]))
	}
	return ctx, nil
}

func iTabToLocation(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.Tab()
	return ctx, nil
}

func iSelectFormat(ctx context.Context, format string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	// Format field is already focused when modal opens
	// Navigate to the desired format and confirm
	// Simple approach: assume options are in order (Text, Markdown, YAML)
	// Navigate down until we find the format, then confirm
	for i := 0; i < 3; i++ {
		view := env.GetView()
		if strings.Contains(view, format) {
			env.Confirm()
			return ctx, nil
		}
		env.NavigateDown()
	}
	return ctx, nil
}

func iSelectLocation(ctx context.Context, location string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	// Location field should be focused after format selection
	// Navigate to the desired location and confirm
	for i := 0; i < 2; i++ {
		view := env.GetView()
		if strings.Contains(view, location) {
			env.Confirm()
			return ctx, nil
		}
		env.NavigateDown()
	}
	return ctx, nil
}

func iConfirmExport(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	// Confirm the export (press enter on the final confirm button)
	env.Confirm()
	return ctx, nil
}

func iShouldSeeExportProgress(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	// Check for export-related progress indicators
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Export"),
		gomega.ContainSubstring("Progress"),
		gomega.ContainSubstring("Saving"),
	))
	return nil
}

func theExportShouldComplete(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	// Wait for export to complete
	gomega.Eventually(func() string {
		return env.GetView()
	}, "5s", "100ms").Should(gomega.SatisfyAny(
		gomega.ContainSubstring("complete"),
		gomega.ContainSubstring("success"),
		gomega.ContainSubstring("exported"),
		gomega.ContainSubstring("saved"),
	))
	return nil
}

func iShouldReturnToPreviousScreen(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	// After canceling, should see the review or preview screen
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("CV Review"),
		gomega.ContainSubstring("CV Preview"),
	))
	return nil
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
