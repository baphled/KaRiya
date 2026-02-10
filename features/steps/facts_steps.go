// Package steps provides BDD step definitions for Godog feature tests.
package steps

import (
	"context"
	"errors"
	"fmt"

	"github.com/baphled/kariya/features/support"
	"github.com/baphled/kariya/internal/domain/career"
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
func RegisterFactsSteps(sc *godog.ScenarioContext) {
	sc.Step(`^I have (\d+) facts? in my profile$`, iHaveNFactsInMyProfile)
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

func iShouldSeeFactText(_ context.Context) error {
	return godog.ErrPending
}

func iShouldSeeStrengthSignals(_ context.Context) error {
	return godog.ErrPending
}

func iShouldSeeCategories(_ context.Context) error {
	return godog.ErrPending
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

func iHaveAFact(_ context.Context, _ string) (context.Context, error) {
	return nil, godog.ErrPending
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

func iShouldSeeCompetencyCategories(_ context.Context) error {
	return godog.ErrPending
}

func iShouldSeeRoleFit(_ context.Context) error {
	return godog.ErrPending
}

func iShouldSeeAudienceRelevance(_ context.Context) error {
	return godog.ErrPending
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
		gomega.ContainSubstring("Role"),
		gomega.ContainSubstring("Audience"),
	))
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

func iSelectCompetencyCategory(_ context.Context, _ string) (context.Context, error) {
	return nil, godog.ErrPending
}

func iTabToRoleFit(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.Tab()
	return ctx, nil
}

func iSelectRoleFit(_ context.Context, _ string) (context.Context, error) {
	return nil, godog.ErrPending
}

func iTabToAudienceRelevance(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.Tab()
	return ctx, nil
}

func iSelectAudience(_ context.Context, _ string) (context.Context, error) {
	return nil, godog.ErrPending
}

func iSubmitTheFactForm(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.SubmitHuhForm()
	return ctx, nil
}

func theFactShouldHaveText(_ context.Context, _ string) error {
	return godog.ErrPending
}

func theFactShouldHaveCategories(_ context.Context, _ string) error {
	return godog.ErrPending
}

func theFactShouldHaveAudiences(_ context.Context, _ string) error {
	return godog.ErrPending
}

func iHaveAFactWithCategory(_ context.Context, _ string) (context.Context, error) {
	return nil, godog.ErrPending
}

func iDeselectCompetencyCategory(_ context.Context, _ string) (context.Context, error) {
	return nil, godog.ErrPending
}

func iShouldBeOnCompetencyCategoriesField(_ context.Context) error {
	return godog.ErrPending
}

func iShouldBeOnRoleFitField(_ context.Context) error {
	return godog.ErrPending
}

func iPressShiftTab(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('\t')
	return ctx, nil
}

func iClearTheFactTextField(_ context.Context) (context.Context, error) {
	return nil, godog.ErrPending
}

func iPressYToConfirm(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('y')
	return ctx, nil
}

func iEnterFactTextWithNCharacters(_ context.Context, _ int) (context.Context, error) {
	return nil, godog.ErrPending
}

func iPressRToRefresh(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('r')
	return ctx, nil
}

func theFactsShouldBeReloaded(_ context.Context) error {
	return godog.ErrPending
}

func iShouldBeAtTheLastFact(_ context.Context) error {
	return godog.ErrPending
}

func iShouldBeAtTheFirstFact(_ context.Context) error {
	return godog.ErrPending
}

func iShouldSeeDifferentFacts(_ context.Context) error {
	return godog.ErrPending
}

func iShouldSeeAvailableShortcuts(_ context.Context) error {
	return godog.ErrPending
}
