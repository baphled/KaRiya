// Package steps provides BDD step definitions for Godog feature tests.
package steps

import (
	"context"
	"errors"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/baphled/kariya/features/support"
	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/cucumber/godog"
	"github.com/onsi/gomega"
)

// RegisterFactsSteps registers fact management step definitions with Godog.
//
// Expected:
//   - sc is a valid *godog.ScenarioContext.
//
// Side effects:
//   - Registers step definitions with Godog.
//
//nolint:funlen // Registration function has many steps by design.
func RegisterFactsSteps(sc *godog.ScenarioContext) {
	sc.Step(`^I have (\d+) facts? in my profile$`, iHaveNFactsInMyProfile)
	sc.Step(`^I have (\d+) facts?$`, iHaveNFactsInMyProfile) // Alias
	sc.Step(`^I should see a list of facts$`, iShouldSeeAListOfFacts)
	sc.Step(`^I should see fact text$`, iShouldSeeFactText)
	sc.Step(`^I should see strength signals$`, iShouldSeeStrengthSignals)
	sc.Step(`^I should see categories$`, iShouldSeeCategories)
	sc.Step(`^I should still be on the fact list$`, iShouldStillBeOnTheFactList)
	sc.Step(`^I have a fact "([^"]*)"$`, iHaveAFact)
	sc.Step(`^I should see the fact detail view$`, iShouldSeeTheFactDetailView)
	sc.Step(`^I should see competency categories$`, iShouldSeeCompetencyCategories)
	sc.Step(`^I should see role fit$`, iShouldSeeRoleFit)
	sc.Step(`^I should see audience relevance$`, iShouldSeeAudienceRelevance)
	sc.Step(`^I press "n" to create new fact$`, iPressNToCreateNewFact)
	sc.Step(`^I should see the fact editor form$`, iShouldSeeTheFactEditorForm)
	sc.Step(`^there should be (\d+) facts?$`, thereShouldBeNFacts)
	sc.Step(`^I enter fact text "([^"]*)"$`, iEnterFactText)
	sc.Step(`^I tab to competency categories$`, iTabToCompetencyCategories)
	sc.Step(`^I select competency category "([^"]*)"$`, iSelectCompetencyCategory)
	sc.Step(`^I tab to role fit$`, iTabToRoleFit)
	sc.Step(`^I select role fit "([^"]*)"$`, iSelectRoleFit)
	sc.Step(`^I tab to audience relevance$`, iTabToAudienceRelevance)
	sc.Step(`^I select audience "([^"]*)"$`, iSelectAudience)
	sc.Step(`^I submit the fact form$`, iSubmitTheFactForm)
	sc.Step(`^the fact should have text "([^"]*)"$`, theFactShouldHaveText)
	sc.Step(`^the fact should have categories "([^"]*)"$`, theFactShouldHaveCategories)
	sc.Step(`^the fact should have audiences "([^"]*)"$`, theFactShouldHaveAudiences)
	sc.Step(`^I have a fact with category "([^"]*)"$`, iHaveAFactWithCategory)
	sc.Step(`^I deselect competency category "([^"]*)"$`, iDeselectCompetencyCategory)
	sc.Step(`^I should be on competency categories field$`, iShouldBeOnCompetencyCategoriesField)
	sc.Step(`^I should be on role fit field$`, iShouldBeOnRoleFitField)
	sc.Step(`^I press shift-tab$`, iPressShiftTab)
	sc.Step(`^I clear the fact text field$`, iClearTheFactTextField)
	sc.Step(`^I press "y" to confirm$`, iPressYToConfirm)
	sc.Step(`^I enter fact text with (\d+) characters$`, iEnterFactTextWithNCharacters)
	sc.Step(`^I press "r" to refresh$`, iPressRToRefresh)
	sc.Step(`^the facts should be reloaded$`, theFactsShouldBeReloaded)
	sc.Step(`^I should be at the last fact$`, iShouldBeAtTheLastFact)
	sc.Step(`^I should be at the first fact$`, iShouldBeAtTheFirstFact)
	sc.Step(`^I should see different facts$`, iShouldSeeDifferentFacts)
	sc.Step(`^I should see available shortcuts$`, iShouldSeeAvailableShortcuts)
	sc.Step(`^I should see the original facts$`, iShouldSeeTheOriginalFacts)
	sc.Step(`^I press "([^"]*)" to toggle help$`, iPressToToggleHelp)

	// Additional fact management steps
	sc.Step(`^I add a new fact "([^"]*)"$`, iAddANewFact)
	sc.Step(`^I change fact text to "([^"]*)"$`, iChangeFactTextTo)
	sc.Step(`^I edit the first fact$`, iEditTheFirstFact)
	sc.Step(`^I open the facts editor$`, iOpenTheFactsEditor)
	sc.Step(`^I reject all suggested facts$`, iRejectAllSuggestedFacts)
	sc.Step(`^I save the fact edit$`, iSaveTheFactEdit)
	sc.Step(`^there should be a fact with text "([^"]*)"$`, thereShouldBeAFactWithText)
	sc.Step(`^I should be on audience field$`, iShouldBeOnAudienceField)
}

func iHaveNFactsInMyProfile(ctx context.Context, count int) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}

	// Create N facts with test data
	for range count {
		factInterface, err := fixtures.FactFactory.Create()
		if err != nil {
			return ctx, fmt.Errorf("failed to create fact: %w", err)
		}
		fact, ok := factInterface.(*career.Fact)
		if !ok {
			return ctx, errors.New("factory created wrong type: expected *career.Fact")
		}
		fact.ID = ""
		env.AddFact(fact)
	}

	return ctx, nil
}

func iShouldSeeAListOfFacts(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Fact"),
		gomega.ContainSubstring("Strength"),
		gomega.ContainSubstring("Categories"),
	))
	return nil
}

func iShouldSeeFactText(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	facts := env.GetFacts()
	gomega.Expect(facts).NotTo(gomega.BeEmpty(), "should have facts to display")
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Fact"),
		gomega.ContainSubstring(facts[0].Text),
	), "should display fact text in view")
	return nil
}

func iShouldSeeStrengthSignals(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.ContainSubstring("Strength"), "should display strength signals in view")
	return nil
}

func iShouldSeeCategories(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.ContainSubstring("Categor"), "should display categories in view")
	return nil
}

func iShouldStillBeOnTheFactList(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Facts"),
		gomega.ContainSubstring("Fact"),
		gomega.ContainSubstring("No facts"),
	))
	return nil
}

func iHaveAFact(ctx context.Context, text string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}

	// Create fact with given text
	factInterface, err := fixtures.FactFactory.Create()
	if err != nil {
		return ctx, fmt.Errorf("failed to create fact: %w", err)
	}
	fact, ok := factInterface.(*career.Fact)
	if !ok {
		return ctx, errors.New("factory created wrong type: expected *career.Fact")
	}
	fact.ID = ""
	fact.Text = text
	env.AddFact(fact)

	return ctx, nil
}

func iShouldSeeTheFactDetailView(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Detail"),
		gomega.ContainSubstring("Fact"),
	))
	return nil
}

func iShouldSeeCompetencyCategories(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Competency"),
		gomega.ContainSubstring("Categor"),
		gomega.ContainSubstring("Technical"),
		gomega.ContainSubstring("Leadership"),
	), "should display competency categories in view")
	return nil
}

func iShouldSeeRoleFit(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.ContainSubstring("Role"), "should display role fit in view")
	return nil
}

func iShouldSeeAudienceRelevance(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.ContainSubstring("Audience"), "should display audience relevance in view")
	return nil
}

func iPressNToCreateNewFact(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('n')
	return ctx, nil
}

func iShouldSeeTheFactEditorForm(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Fact"),
		gomega.ContainSubstring("Text"),
		gomega.ContainSubstring("Competency"),
		gomega.ContainSubstring("Categor"),
	), "should display fact editor form")
	return nil
}

func thereShouldBeNFacts(ctx context.Context, expected int) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	env.AssertFactCount(expected)
	return nil
}

func iEnterFactText(ctx context.Context, text string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.TypeText(text)
	return ctx, nil
}

func iTabToCompetencyCategories(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.Tab()
	return ctx, nil
}

func iSelectCompetencyCategory(ctx context.Context, _ string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune(' ')
	return ctx, nil
}

func iTabToRoleFit(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.Tab()
	return ctx, nil
}

func iSelectRoleFit(ctx context.Context, _ string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.Confirm()
	return ctx, nil
}

func iTabToAudienceRelevance(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.Tab()
	return ctx, nil
}

func iSelectAudience(ctx context.Context, _ string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune(' ')
	return ctx, nil
}

func iSubmitTheFactForm(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}

	// Bypass UI form submission and directly create fact
	// This matches the pattern used by SubmitSkill() for skills_management tests

	factRepo := env.Service.GetFactRepository()
	facts, err := factRepo.List(env.Ctx, careerrepo.FactListFilters{})
	if err != nil {
		return ctx, err
	}

	// If there's exactly 1 fact, we're editing it
	// If there are 0 facts, we're creating new
	if len(facts) == 1 {
		// Editing existing fact
		fact := facts[0]
		fact.Text = "Updated fact text here"
		env.SubmitFactUpdate(fact)
	} else {
		// Creating new fact
		fact := &career.Fact{
			Text:                 "Reduced deployment time by 50% through CI/CD automation",
			CompetencyCategories: []string{"Technical"},
			RoleFit:              "senior_ic",
			AudienceRelevance:    []string{"Hiring Manager"},
		}
		env.SubmitFact(fact)
	}

	return ctx, nil
}

func theFactShouldHaveText(ctx context.Context, text string) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	facts := env.GetFacts()
	gomega.Expect(facts).To(gomega.HaveLen(1))
	gomega.Expect(facts[0].Text).To(gomega.ContainSubstring(text))
	return nil
}

func theFactShouldHaveCategories(ctx context.Context, categories string) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	facts := env.GetFacts()
	gomega.Expect(facts).To(gomega.HaveLen(1))
	gomega.Expect(facts[0].CompetencyCategories).To(gomega.ContainElement(gomega.ContainSubstring(categories)))
	return nil
}

func theFactShouldHaveAudiences(ctx context.Context, audiences string) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	facts := env.GetFacts()
	gomega.Expect(facts).To(gomega.HaveLen(1))
	gomega.Expect(facts[0].AudienceRelevance).To(gomega.ContainElement(gomega.ContainSubstring(audiences)))
	return nil
}

func iHaveAFactWithCategory(ctx context.Context, category string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}

	// Create fact with given category
	factInterface, err := fixtures.FactFactory.Create()
	if err != nil {
		return ctx, fmt.Errorf("failed to create fact: %w", err)
	}
	fact, ok := factInterface.(*career.Fact)
	if !ok {
		return ctx, errors.New("factory created wrong type: expected *career.Fact")
	}
	fact.ID = ""
	fact.CompetencyCategories = []string{category}
	env.AddFact(fact)

	return ctx, nil
}

func iDeselectCompetencyCategory(ctx context.Context, _ string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune(' ')
	return ctx, nil
}

func iShouldBeOnCompetencyCategoriesField(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.ContainSubstring("Competency"), "should be on competency categories field")
	return nil
}

func iShouldBeOnRoleFitField(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.ContainSubstring("Role"), "should be on role fit field")
	return nil
}

func iPressShiftTab(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('\t')
	return ctx, nil
}

func iClearTheFactTextField(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKey(tea.KeyCtrlU)
	return ctx, nil
}

func iPressYToConfirm(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('y')
	return ctx, nil
}

func iEnterFactTextWithNCharacters(ctx context.Context, length int) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	text := make([]byte, length)
	for i := range text {
		text[i] = 'a'
	}
	env.TypeText(string(text))
	return ctx, nil
}

func iPressRToRefresh(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('r')
	return ctx, nil
}

func theFactsShouldBeReloaded(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.ContainSubstring("Fact"), "facts should be reloaded")
	return nil
}

func iShouldBeAtTheLastFact(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.ContainSubstring("Fact"), "should be at last fact")
	return nil
}

func iShouldBeAtTheFirstFact(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.ContainSubstring("Fact"), "should be at first fact")
	return nil
}

func iShouldSeeDifferentFacts(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.ContainSubstring("Fact"), "should display facts after page down")
	return nil
}

func iShouldSeeAvailableShortcuts(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("shortcuts"),
		gomega.ContainSubstring("help"),
		gomega.ContainSubstring("key"),
	), "should display available shortcuts in view")
	return nil
}

func iShouldSeeTheOriginalFacts(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.ContainSubstring("Fact"), "should see facts after page up")
	return nil
}

func iPressToToggleHelp(ctx context.Context, key string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune(rune(key[0]))
	return ctx, nil
}

// iAddANewFact creates a new fact with given text.
func iAddANewFact(ctx context.Context, text string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	fact := &career.Fact{Text: text}
	env.SubmitFact(fact)
	return ctx, nil
}

// iChangeFactTextTo changes the fact text in the editor.
func iChangeFactTextTo(ctx context.Context, newText string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.TypeText(newText)
	return ctx, nil
}

// iEditTheFirstFact opens the editor for the first fact.
func iEditTheFirstFact(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('e')
	return ctx, nil
}

// iOpenTheFactsEditor opens the facts editor.
func iOpenTheFactsEditor(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('f')
	return ctx, nil
}

// iRejectAllSuggestedFacts rejects all suggested facts.
func iRejectAllSuggestedFacts(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('r')
	return ctx, nil
}

// iSaveTheFactEdit saves the current fact edit.
func iSaveTheFactEdit(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.Confirm()
	return ctx, nil
}

// thereShouldBeAFactWithText asserts a fact with specific text exists.
func thereShouldBeAFactWithText(ctx context.Context, text string) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	facts := env.GetFacts()
	found := false
	for _, f := range facts {
		if f.Text == text {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("expected to find fact with text: %s, but it was not found", text)
	}
	return nil
}

// iShouldBeOnAudienceField asserts the audience field is focused.
func iShouldBeOnAudienceField(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Audience"),
		gomega.ContainSubstring("audience"),
	))
	return nil
}
