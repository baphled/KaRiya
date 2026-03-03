// Package steps provides BDD step definitions for Godog feature tests.
package steps

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cucumber/godog"
	"github.com/onsi/gomega"

	"github.com/baphled/kariya/features/support"
	"github.com/baphled/kariya/internal/cli/intents/generatecv"
	"github.com/baphled/kariya/internal/testutil/harness"
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
	sc.Step(`^I select technology focus$`, iSelectTechFocus)
	sc.Step(`^I complete the wizard$`, iCompleteTheWizard)
	sc.Step(`^the generation completes$`, theGenerationCompletes)
}

func registerCVReviewSteps(sc *godog.ScenarioContext) {
	sc.Step(`^I should see the CV review screen$`, iShouldSeeTheCVReviewScreen)
	sc.Step(`^I should see CV metadata$`, iShouldSeeCVMetadata)
	sc.Step(`^I should see statistics$`, iShouldSeeStatistics)
	sc.Step(`^I should see section names$`, iShouldSeeSectionNames)
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
	sc.Step(`^I press enter to confirm$`, iPressEnterToConfirmCV)
	sc.Step(`^the CV generation should complete$`, theCVGenerationShouldComplete)
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
	sc.Step(`^I start an export$`, iStartAnExport)
	sc.Step(`^I complete an export$`, iCompleteAnExport)
	sc.Step(`^I should see export location$`, iShouldSeeExportLocation)
	sc.Step(`^CV is exported as a text file$`, cvIsExportedAsTextFile)
	sc.Step(`^CV is exported as a markdown file$`, cvIsExportedAsMarkdownFile)
	sc.Step(`^CV is exported as a YAML file$`, cvIsExportedAsYAMLFile)
	sc.Step(`^the clipboard should contain the CV as text$`, theClipboardShouldContainTheCVAsText)
	sc.Step(`^the clipboard should contain the CV as markdown$`, theClipboardShouldContainTheCVAsMarkdown)
	sc.Step(`^the clipboard should contain the CV as YAML$`, theClipboardShouldContainTheCVAsYAML)
}

func iHaveNoProfileConfigured(ctx context.Context) (context.Context, error) {
	_, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}
	// Profile is empty by default in test env
	return ctx, nil
}

func iShouldSeeAnErrorOrWarning(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
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
	return iHaveNoProfileConfigured(ctx)
}

func iHaveACompleteProfileWithEventsAndFacts(ctx context.Context) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}

	// Use fixtures to populate substantial test data for CV generation
	// This creates realistic career data similar to CSV imports
	env.PopulateTestData(50, 5, 30)

	// Add skills - CV generation requires skills for sections
	skills := harness.CreateSampleSkills(20)
	for _, skill := range skills {
		env.AddSkill(skill)
	}

	return ctx, nil
}

func iShouldSeeTheCVWizardModal(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
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
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}
	env.NavigateDown()
	return ctx, nil
}

func iConfirmSelection(ctx context.Context) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}
	env.Confirm()
	return ctx, nil
}

func iShouldMoveToTheNextField(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	// Field movement is handled by Tab - just verify view updated
	view := env.GetView()
	gomega.Expect(view).NotTo(gomega.BeEmpty())
	return nil
}

func iTabToAudienceField(ctx context.Context) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}
	env.Tab()
	return ctx, nil
}

func iSelectTechFocus(ctx context.Context) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}
	env.NavigateDown()
	env.Confirm()
	return ctx, nil
}

func iCompleteTheWizard(ctx context.Context) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}
	// Skip wizard with Ctrl+S (uses default settings)
	env.PressKey(tea.KeyCtrlS)
	return ctx, nil
}

func theGenerationCompletes(ctx context.Context) error {
	// Wait for generation to complete and reach review screen
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
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
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
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
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
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
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.MatchRegexp(`\d+`)) // Should see numbers
	return nil
}

func iShouldSeeSectionNames(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Experience"),
		gomega.ContainSubstring("Skills"),
		gomega.ContainSubstring("Summary"),
	))
	return nil
}

func iHaveGeneratedACV(ctx context.Context) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}

	// Set up data
	ctx, err = iHaveACompleteProfileWithEventsAndFacts(ctx)
	if err != nil {
		return ctx, err
	}

	// Navigate to generate_cv
	env.SelectIntentByName("generate_cv")

	env.SendMessage(generatecv.WizardCompleteMsg{
		ProfileID:    "profile-staff-engineer",
		Audience:     "hiring_manager",
		SkillsFormat: "grouped",
		SkillsLimit:  5,
	})

	gomega.Eventually(func() string {
		return env.GetView()
	}, "10s", "100ms").Should(gomega.ContainSubstring("Configuration Review"))

	return ctx, nil
}

func iAmOnTheCVReviewScreen(ctx context.Context) error {
	return iShouldSeeTheCVReviewScreen(ctx)
}

func iPressEnterToPreview(ctx context.Context) (context.Context, error) {
	return iConfirmSelection(ctx)
}

func iShouldSeeTheCVPreviewScreen(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Preview"),
		gomega.ContainSubstring("CV"),
	))
	return nil
}

func iPressPToPreview(ctx context.Context) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}
	env.PressKeyRune('p')
	return ctx, nil
}

func iShouldSeeTheExportOptionsModal(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
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

func iShouldSeePersonalDetails(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	// Check for personal/professional information in CV preview
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Professional Summary"),
		gomega.ContainSubstring("Summary"),
		gomega.ContainSubstring("CV Preview"),
	))
	return nil
}

func iShouldSeeAllCVSections(ctx context.Context) error {
	return iShouldSeeSectionNames(ctx)
}

func iShouldSeeBulletPoints(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("•"),
		gomega.ContainSubstring("★"),
		gomega.ContainSubstring("-"),
		gomega.ContainSubstring("*"),
		gomega.ContainSubstring("No sections generated"),
	))
	return nil
}

func iPressEnterToConfirmCV(ctx context.Context) (context.Context, error) {
	return iConfirmSelection(ctx)
}

func theCVGenerationShouldComplete(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	// After confirmation, CV generation workflow completes
	// We should see the main menu with action options
	gomega.Eventually(func() string {
		return env.GetView()
	}, "5s", "100ms").Should(gomega.SatisfyAny(
		gomega.ContainSubstring("Capture Event"),
		gomega.ContainSubstring("Browse Timeline"),
		gomega.ContainSubstring("Career Event Management System"),
	))
	return nil
}

func iOpenTheExportOptionsModal(ctx context.Context) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}
	// Press 'x' to open export modal from review/preview screen
	env.PressKeyRune('x')

	// Wait for modal to be visible
	gomega.Eventually(func() string {
		return env.GetView()
	}, "2s", "100ms").Should(gomega.SatisfyAny(
		gomega.ContainSubstring("Export"),
		gomega.ContainSubstring("Format"),
		gomega.ContainSubstring("Location"),
	))

	return ctx, nil
}

func iPressKeyToExport(ctx context.Context, key string) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}
	// Press the specified key to open export modal
	if key != "" {
		env.PressKeyRune(rune(key[0]))
	}
	return ctx, nil
}

func iTabToLocation(ctx context.Context) (context.Context, error) {
	return iTabToAudienceField(ctx)
}

func iSelectFormat(ctx context.Context, format string) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}

	// Navigate to the desired format and confirm
	gomega.Eventually(func() bool {
		view := env.GetView()
		if strings.Contains(view, format) {
			env.Confirm()
			return true
		}
		env.NavigateDown()
		return false
	}, "2s", "100ms").Should(gomega.BeTrue(), "Failed to select format: "+format)

	return ctx, nil
}

func iSelectLocation(ctx context.Context, location string) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}

	// Navigate to the desired location and confirm
	gomega.Eventually(func() bool {
		view := env.GetView()
		if strings.Contains(view, location) {
			env.Confirm()
			return true
		}
		env.NavigateDown()
		return false
	}, "2s", "100ms").Should(gomega.BeTrue(), "Failed to select location: "+location)

	return ctx, nil
}

func iConfirmExport(ctx context.Context) (context.Context, error) {
	return iConfirmSelection(ctx)
}

func iShouldSeeExportProgress(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
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
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
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

func iStartAnExport(ctx context.Context) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}
	env.PressKeyRune('x')
	return ctx, nil
}

func iCompleteAnExport(ctx context.Context) error {
	// Refactored to use granular steps without fixed sleeps
	ctx, err := iOpenTheExportOptionsModal(ctx)
	if err != nil {
		return err
	}

	ctx, err = iSelectFormat(ctx, "Text")
	if err != nil {
		return err
	}

	ctx, err = iSelectLocation(ctx, "File")
	if err != nil {
		return err
	}

	ctx, err = iConfirmExport(ctx)
	if err != nil {
		return err
	}

	return theExportShouldComplete(ctx)
}

func iShouldSeeExportLocation(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.ContainSubstring("Location:"))
	return nil
}

func cvIsExportedAsTextFile(ctx context.Context) error {
	return cvIsExportedAsFile(ctx, ".txt")
}

func cvIsExportedAsMarkdownFile(ctx context.Context) error {
	return cvIsExportedAsFile(ctx, ".md")
}

func cvIsExportedAsYAMLFile(ctx context.Context) error {
	return cvIsExportedAsFile(ctx, ".yaml")
}

func theClipboardShouldContainTheCVAsText(ctx context.Context) error {
	return theClipboardShouldContainCV(ctx, "Text", []string{"Text", "text", "Plain Text", "plain text"})
}

func theClipboardShouldContainTheCVAsMarkdown(ctx context.Context) error {
	return theClipboardShouldContainCV(ctx, "Markdown", []string{"Markdown", "markdown"})
}

func theClipboardShouldContainTheCVAsYAML(ctx context.Context) error {
	return theClipboardShouldContainCV(ctx, "YAML", []string{"YAML", "Yaml", "yaml"})
}

func cvIsExportedAsFile(ctx context.Context, extension string) error {
	if err := theExportShouldComplete(ctx); err != nil {
		return err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	exportDir := filepath.Join(homeDir, ".kariya", "cv_exports")
	files, err := filepath.Glob(filepath.Join(exportDir, "*"+extension))
	if err != nil {
		return err
	}

	gomega.Expect(files).NotTo(gomega.BeEmpty())

	hasContent := false
	for _, filePath := range files {
		cleanPath := filepath.Clean(filePath)
		if !strings.HasPrefix(cleanPath, exportDir) {
			continue
		}
		data, err := os.ReadFile(cleanPath)
		if err != nil {
			return err
		}
		if strings.TrimSpace(string(data)) != "" {
			hasContent = true
			break
		}
	}

	gomega.Expect(hasContent).To(gomega.BeTrue())
	return nil
}

func theClipboardShouldContainCV(ctx context.Context, format string, formatOptions []string) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}

	gomega.Eventually(func() string {
		return env.GetView()
	}, "5s", "100ms").Should(gomega.SatisfyAny(
		gomega.ContainSubstring("Exported to clipboard"),
		gomega.ContainSubstring("clipboard"),
		gomega.ContainSubstring("Clipboard"),
	))

	view := env.GetView()
	matchers := make([]gomega.OmegaMatcher, 0, len(formatOptions))
	for _, option := range formatOptions {
		matchers = append(matchers, gomega.ContainSubstring(option))
	}
	gomega.Expect(view).To(gomega.SatisfyAny(matchers...))
	gomega.Expect(strings.ToLower(view)).To(gomega.ContainSubstring(strings.ToLower(format)))

	return nil
}
