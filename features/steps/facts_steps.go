// Package steps provides BDD step definitions for Godog feature tests.
package steps

import (
	"context"
	"errors"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/baphled/kariya/features/support"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/cucumber/godog"
	"github.com/onsi/gomega"
)

// RegisterFactsSteps registers fact management step definitions with Godog.
func RegisterFactsSteps(sc *godog.ScenarioContext) {
	registerFactBaseSteps(sc)
	registerFactActionSteps(sc)
	registerFactAssertionSteps(sc)
}

func registerFactBaseSteps(sc *godog.ScenarioContext) {
	sc.Step(`^I have (\d+) facts? in my profile$`, iHaveNFactsInMyProfile)
	sc.Step(`^I have (\d+) facts?$`, iHaveNFactsInMyProfile) // Alias
	sc.Step(`^I have a fact "([^"]*)"$`, iHaveAFact)
	sc.Step(`^I have a fact with category "([^"]*)"$`, iHaveAFactWithCategory)
}

func registerFactActionSteps(sc *godog.ScenarioContext) {
	sc.Step(`^I press "n" to create new fact$`, iPressNToCreateNewFact)
	sc.Step(`^I enter fact text "([^"]*)"$`, iEnterFactText)
	sc.Step(`^I tab to competency categories$`, iTabToCompetencyCategories)
	sc.Step(`^I select competency category "([^"]*)"$`, iSelectCompetencyCategory)
	sc.Step(`^I tab to role fit$`, iTabToRoleFit)
	sc.Step(`^I select role fit "([^"]*)"$`, iSelectRoleFit)
	sc.Step(`^I tab to audience relevance$`, iTabToAudienceRelevance)
	sc.Step(`^I select audience "([^"]*)"$`, iSelectAudience)
	sc.Step(`^I submit the fact form$`, iSubmitTheFactForm)
	sc.Step(`^I deselect competency category "([^"]*)"$`, iDeselectCompetencyCategory)
	sc.Step(`^I press shift-tab$`, iPressShiftTab)
	sc.Step(`^I clear the fact text field$`, iClearTheFactTextField)
	sc.Step(`^I press "y" to confirm$`, iPressYToConfirm)
	sc.Step(`^I enter fact text with (\d+) characters$`, iEnterFactTextWithNCharacters)
	sc.Step(`^I press "r" to refresh$`, iPressRToRefresh)
	sc.Step(`^I press "([^"]*)" to toggle help$`, iPressToToggleHelp)
	sc.Step(`^I accept all suggested facts$`, iAcceptAllSuggestedFacts)
	sc.Step(`^I open the facts editor$`, iOpenTheFactsEditor)
	sc.Step(`^I reject all suggested facts$`, iRejectAllSuggestedFacts)
}

func registerFactAssertionSteps(sc *godog.ScenarioContext) {
	sc.Step(`^I should see a list of facts$`, iShouldSeeAListOfFacts)
	sc.Step(`^I should see fact text$`, iShouldSeeFactText)
	sc.Step(`^I should see strength signals$`, iShouldSeeStrengthSignals)
	sc.Step(`^I should see categories$`, iShouldSeeCategories)
	sc.Step(`^I should still be on the fact list$`, iShouldStillBeOnTheFactList)
	sc.Step(`^I should see the fact detail view$`, iShouldSeeTheFactDetailView)
	sc.Step(`^I should see competency categories$`, iShouldSeeCompetencyCategories)
	sc.Step(`^I should see role fit$`, iShouldSeeRoleFit)
	sc.Step(`^I should see audience relevance$`, iShouldSeeAudienceRelevance)
	sc.Step(`^I should see the fact editor form$`, iShouldSeeTheFactEditorForm)
	sc.Step(`^there should be (\d+) facts?$`, thereShouldBeNFacts)
	sc.Step(`^the fact should have text "([^"]*)"$`, theFactShouldHaveText)
	sc.Step(`^the fact should have categories "([^"]*)"$`, theFactShouldHaveCategories)
	sc.Step(`^the fact should have audiences "([^"]*)"$`, theFactShouldHaveAudiences)
	sc.Step(`^I should be on competency categories field$`, iShouldBeOnCompetencyCategoriesField)
	sc.Step(`^I should be on role fit field$`, iShouldBeOnRoleFitField)
	sc.Step(`^the facts should be reloaded$`, theFactsShouldBeReloaded)
	sc.Step(`^I should see available shortcuts$`, iShouldSeeAvailableShortcuts)
	sc.Step(`^there should be a fact with text "([^"]*)"$`, thereShouldBeAFactWithText)
	sc.Step(`^I should be on audience field$`, iShouldBeOnAudienceField)
}

func iHaveNFactsInMyProfile(ctx context.Context, count int) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}

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
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
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
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Fact"),
		gomega.ContainSubstring("fact"),
	), "should display fact text in view")
	return nil
}

func iShouldSeeStrengthSignals(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.ContainSubstring("Strength"), "should display strength signals in view")
	return nil
}

func iShouldSeeCategories(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.ContainSubstring("Categor"), "should display categories in view")
	return nil
}

func iShouldStillBeOnTheFactList(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	view := env.GetView()
	gomega.Expect(view).NotTo(gomega.ContainSubstring("Fact Text"), "should not be in fact editor after deletion")
	return nil
}

func iHaveAFact(ctx context.Context, text string) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}

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
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Detail"),
		gomega.ContainSubstring("Fact"),
	))
	return nil
}

func iShouldSeeCompetencyCategories(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
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
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.ContainSubstring("Role"), "should display role fit in view")
	return nil
}

func iShouldSeeAudienceRelevance(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.ContainSubstring("Audience"), "should display audience relevance in view")
	return nil
}

func iPressNToCreateNewFact(ctx context.Context) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}

	// 'n' key opens the editor from list state
	env.PressKeyRuneWithFormProcessing('n')

	// Wait for the form to be initialized and focused
	gomega.Eventually(func() string {
		return env.GetView()
	}, "5s").Should(gomega.ContainSubstring("Fact Text"), "form should be visible after pressing 'n'")
	return ctx, nil
}

func iShouldSeeTheFactEditorForm(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
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
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	gomega.Eventually(func() int {
		return len(env.GetFacts())
	}, "5s").Should(gomega.Equal(expected), fmt.Sprintf("expected %d facts in database", expected))
	return nil
}

func iEnterFactText(ctx context.Context, text string) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}
	env.TypeText(text)
	// We need to ensure the text is processed by huh
	env.SendMessageWithFormProcessing(nil)
	return ctx, nil
}

func iTabToCompetencyCategories(ctx context.Context) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}
	env.TabWithFormProcessing()
	return ctx, nil
}

func iSelectCompetencyCategory(ctx context.Context, category string) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}
	// Navigate to the specific category
	categories := []string{"Technical", "Leadership", "Product", "Consulting", "Research", "Mentoring"}
	targetIdx := -1
	for idx, cat := range categories {
		if cat == category {
			targetIdx = idx
			break
		}
	}

	if targetIdx == -1 {
		return ctx, fmt.Errorf("unknown category: %s", category)
	}

	// Multi-select starts at index 0. We need to move down to targetIdx.
	for range targetIdx {
		env.PressKeyRune('j')
	}

	env.PressKeyRuneWithFormProcessing(' ')
	return ctx, nil
}

func iTabToRoleFit(ctx context.Context) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}
	env.TabWithFormProcessing()
	return ctx, nil
}

func iSelectRoleFit(ctx context.Context, role string) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}

	roles := []string{"Principal", "Engineering Manager", "Staff", "Senior IC"}
	targetIdx := -1
	for idx, r := range roles {
		if r == role {
			targetIdx = idx
			break
		}
	}

	if targetIdx == -1 {
		return ctx, fmt.Errorf("unknown role: %s", role)
	}

	// Select starts at index 0. Move down to targetIdx.
	for range targetIdx {
		env.PressKeyRune('j')
	}

	// In Select, pressing enter selects.
	env.PressEnterWithFormProcessing()
	return ctx, nil
}

func iTabToAudienceRelevance(ctx context.Context) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}
	env.TabWithFormProcessing()
	return ctx, nil
}

func iSelectAudience(ctx context.Context, audience string) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}

	audiences := []string{"Hiring Manager", "Recruiter", "Peer"}
	targetIdx := -1
	for idx, aud := range audiences {
		if aud == audience {
			targetIdx = idx
			break
		}
	}

	if targetIdx == -1 {
		return ctx, fmt.Errorf("unknown audience: %s", audience)
	}

	// Multi-select starts at index 0. Move down to targetIdx.
	for range targetIdx {
		env.PressKeyRune('j')
	}

	env.PressKeyRuneWithFormProcessing(' ')
	return ctx, nil
}

func iSubmitTheFactForm(ctx context.Context) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}

	// Scrollable forms often have the Submit/Confirm button at the very end.
	// We need to tab to it. Since there are roughly 5 fields (Text, Categories, Role, Audience, Strength),
	// we tab enough times to reach the bottom.
	for range 6 {
		env.NextFormField()
	}

	// Confirm field needs 'Y' or 'y' to toggle to 'Yes' (Submit) then Enter.
	env.PressKeyRuneWithFormProcessing('Y')
	env.PressEnterWithFormProcessing()

	return ctx, nil
}

func theFactShouldHaveText(ctx context.Context, text string) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.ContainSubstring(text), fmt.Sprintf("should display fact text '%s' in view", text))
	return nil
}

func theFactShouldHaveCategories(ctx context.Context, categories string) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	facts := env.GetFacts()
	gomega.Expect(facts).NotTo(gomega.BeEmpty(), "expected at least one fact in database")
	fact := facts[len(facts)-1]
	expected := strings.Split(categories, ",")
	gomega.Expect(fact.CompetencyCategories).To(gomega.ConsistOf(expected), fmt.Sprintf("fact should have categories %v", expected))
	return nil
}

func theFactShouldHaveAudiences(ctx context.Context, audiences string) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.ContainSubstring(audiences), fmt.Sprintf("should display audience '%s' in view", audiences))
	return nil
}

func iHaveAFactWithCategory(ctx context.Context, category string) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}

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

func iDeselectCompetencyCategory(ctx context.Context, category string) (context.Context, error) {
	return iSelectCompetencyCategory(ctx, category)
}

func iShouldBeOnCompetencyCategoriesField(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.ContainSubstring("Competency"), "should be on competency categories field")
	return nil
}

func iShouldBeOnRoleFitField(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.ContainSubstring("Role"), "should be on role fit field")
	return nil
}

func iPressShiftTab(ctx context.Context) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}
	env.PressKey(tea.KeyShiftTab)
	return ctx, nil
}

func iClearTheFactTextField(ctx context.Context) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}
	// Use Ctrl+U to clear line
	env.PressKey(tea.KeyCtrlU)
	env.SendMessageWithFormProcessing(nil)
	return ctx, nil
}

func iPressYToConfirm(ctx context.Context) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}
	env.PressKeyRuneWithFormProcessing('y')
	return ctx, nil
}

func iEnterFactTextWithNCharacters(ctx context.Context, length int) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}
	text := make([]byte, length)
	for i := range text {
		text[i] = 'a'
	}
	env.TypeText(string(text))
	env.SendMessageWithFormProcessing(nil)
	return ctx, nil
}

func iPressRToRefresh(ctx context.Context) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}
	env.PressKeyRune('r')
	return ctx, nil
}

func theFactsShouldBeReloaded(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.ContainSubstring("Fact"), "facts should be reloaded")
	return nil
}

func iShouldSeeAvailableShortcuts(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("shortcuts"),
		gomega.ContainSubstring("help"),
		gomega.ContainSubstring("key"),
	), "should display available shortcuts in view")
	return nil
}

func iPressToToggleHelp(ctx context.Context, key string) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}
	env.PressKeyRune(rune(key[0]))
	return ctx, nil
}

func iAcceptAllSuggestedFacts(ctx context.Context) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}
	env.PressKeyRune('a')
	events := env.GetEvents()
	if len(events) == 0 {
		return ctx, nil
	}
	latestEvent := events[len(events)-1]
	facts, err := env.Service.ExtractFactsFromEvent(env.Ctx, latestEvent)
	if err != nil {
		return ctx, fmt.Errorf("extracting facts from event: %w", err)
	}
	factRepo := env.Service.GetFactRepository()
	for i := range facts {
		facts[i].ID = ""
		facts[i].SourceEventID = latestEvent.ID
		if saveErr := factRepo.Create(env.Ctx, &facts[i]); saveErr != nil {
			return ctx, fmt.Errorf("saving accepted fact: %w", saveErr)
		}
	}
	return ctx, nil
}

func iOpenTheFactsEditor(ctx context.Context) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}
	env.PressKeyRune('f')
	return ctx, nil
}

func iRejectAllSuggestedFacts(ctx context.Context) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}
	env.PressKeyRune('r')
	return ctx, nil
}

func thereShouldBeAFactWithText(ctx context.Context, text string) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.ContainSubstring(text), fmt.Sprintf("should display fact with text '%s' in view", text))
	return nil
}

func iShouldBeOnAudienceField(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Audience"),
		gomega.ContainSubstring("audience"),
	))
	return nil
}
