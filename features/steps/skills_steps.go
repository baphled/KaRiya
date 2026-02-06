// Package steps provides BDD step definitions for Godog feature tests.
package steps

import (
	"context"
	"fmt"

	"github.com/baphled/kariya/features/support"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/cucumber/godog"
	"github.com/onsi/gomega"
)

// RegisterSkillsSteps registers skills management step definitions with Godog.
//
//nolint:dupl // Step registration functions look similar but register different steps.
func RegisterSkillsSteps(sc *godog.ScenarioContext) {
	sc.Step(`^I have (\d+) skills? in my profile$`, iHaveNSkillsInMyProfile)
	sc.Step(`^I should see a list of skills$`, iShouldSeeAListOfSkills)
	sc.Step(`^I should still be on the skills list$`, iShouldStillBeOnTheSkillsList)
	sc.Step(`^I have a skill "([^"]*)" with category "([^"]*)"$`, iHaveASkillWithCategory)
	sc.Step(`^I press "a" to add skill$`, iPressAToAddSkill)
	sc.Step(`^I should see the add skill form$`, iShouldSeeTheAddSkillForm)
	sc.Step(`^there should be (\d+) skills?$`, thereShouldBeNSkills)
	sc.Step(`^I enter skill name "([^"]*)"$`, iEnterSkillName)
	sc.Step(`^I select category "([^"]*)"$`, iSelectCategory)
	sc.Step(`^I submit the skill form$`, iSubmitTheSkillForm)
	sc.Step(`^the skill should have name "([^"]*)"$`, theSkillShouldHaveName)
	sc.Step(`^I should see the edit skill form$`, iShouldSeeTheEditSkillForm)
	sc.Step(`^I clear the skill name field$`, iClearTheSkillNameField)
	sc.Step(`^I press "i" to infer skills$`, iPressIToInferSkills)
	sc.Step(`^I should see the loading modal$`, iShouldSeeTheLoadingModal)
	sc.Step(`^the inference completes$`, theInferenceCompletes)
	sc.Step(`^I should see the skill suggestions modal$`, iShouldSeeTheSkillSuggestionsModal)
	sc.Step(`^I accept the first suggestion$`, iAcceptTheFirstSuggestion)
	sc.Step(`^I reject the first suggestion$`, iRejectTheFirstSuggestion)
	sc.Step(`^I should still be on the skill suggestions modal$`, iShouldStillBeOnSkillSuggestionsModal)
	sc.Step(`^the suggestion should be marked as rejected$`, theSuggestionShouldBeMarkedAsRejected)
	sc.Step(`^I create skills with all 15 categories$`, iCreateSkillsWithAll15Categories)
	sc.Step(`^each skill should have a unique category$`, eachSkillShouldHaveUniqueCategory)
}

func iHaveNSkillsInMyProfile(ctx context.Context, count int) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	categories := []string{"Programming", "Database", "Cloud", "DevOps", "Testing"}
	for i := range count {
		skill := &career.Skill{
			Name:     fmt.Sprintf("Skill%d", i+1),
			Category: categories[i%len(categories)],
		}
		env.AddSkill(skill)
	}
	return ctx, nil
}

func iShouldSeeAListOfSkills(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Skill"),
		gomega.ContainSubstring("Category"),
		gomega.ContainSubstring("Name"),
	))
	return nil
}

func iShouldStillBeOnTheSkillsList(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Skills"),
		gomega.ContainSubstring("Skill"),
		gomega.ContainSubstring("No skills"),
	))
	return nil
}

func iHaveASkillWithCategory(ctx context.Context, name, category string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	skill := &career.Skill{
		Name:     name,
		Category: category,
	}
	env.AddSkill(skill)
	return ctx, nil
}

func iPressAToAddSkill(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('a')
	return ctx, nil
}

func iShouldSeeTheAddSkillForm(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Add"),
		gomega.ContainSubstring("Name"),
		gomega.ContainSubstring("Category"),
		gomega.ContainSubstring("Submit"),
	))
	return nil
}

func thereShouldBeNSkills(ctx context.Context, expected int) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	env.AssertSkillCount(expected)
	return nil
}

func iEnterSkillName(ctx context.Context, name string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.TypeText(name)
	return ctx, nil
}

func iSelectCategory(ctx context.Context, _ string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.Tab()
	env.NavigateDown()
	return ctx, nil
}

func iSubmitTheSkillForm(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.SubmitHuhForm()
	return ctx, nil
}

func theSkillShouldHaveName(ctx context.Context, expected string) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	skills := env.GetSkills()
	gomega.Expect(skills).NotTo(gomega.BeEmpty())
	var found bool
	for _, s := range skills {
		if s.Name == expected {
			found = true
			break
		}
	}
	gomega.Expect(found).To(gomega.BeTrue(), "Should have skill with name %s", expected)
	return nil
}

func iShouldSeeTheEditSkillForm(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Edit"),
		gomega.ContainSubstring("Update"),
		gomega.ContainSubstring("Name"),
		gomega.ContainSubstring("Submit"),
	))
	return nil
}

func iClearTheSkillNameField(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	for range 50 {
		env.PressKeyRune('\b')
	}
	return ctx, nil
}

func iPressIToInferSkills(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('i')
	return ctx, nil
}

func iShouldSeeTheLoadingModal(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Loading"),
		gomega.ContainSubstring("Analyzing"),
		gomega.ContainSubstring("..."),
	))
	return nil
}

func theInferenceCompletes(_ context.Context) error {
	return godog.ErrPending
}

func iShouldSeeTheSkillSuggestionsModal(_ context.Context) error {
	return godog.ErrPending
}

func iAcceptTheFirstSuggestion(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('a')
	return ctx, nil
}

func iRejectTheFirstSuggestion(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('r')
	return ctx, nil
}

func iShouldStillBeOnSkillSuggestionsModal(_ context.Context) error {
	return godog.ErrPending
}

func theSuggestionShouldBeMarkedAsRejected(_ context.Context) error {
	return godog.ErrPending
}

func iCreateSkillsWithAll15Categories(_ context.Context) error {
	return godog.ErrPending
}

func eachSkillShouldHaveUniqueCategory(_ context.Context) error {
	return godog.ErrPending
}
