// Package steps provides BDD step definitions for Godog feature tests.
package steps

import (
	"context"
	"fmt"

	"github.com/baphled/kariya/features/support"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/cucumber/godog"
	"github.com/onsi/gomega"
)

// RegisterBurstsSteps registers burst management step definitions with Godog.
func RegisterBurstsSteps(sc *godog.ScenarioContext) {
	registerBurstViewSteps(sc)
	registerBurstEditSteps(sc)
	registerBurstSuggestionSteps(sc)
	registerBurstNavigationSteps(sc)
}

//nolint:dupl // Step registration functions look similar but register different steps.
func registerBurstViewSteps(sc *godog.ScenarioContext) {
	sc.Step(`^I have (\d+) bursts? in my profile$`, iHaveNBurstsInMyProfile)
	sc.Step(`^I should see a list of bursts$`, iShouldSeeAListOfBursts)
	sc.Step(`^I should still be on the burst list$`, iShouldStillBeOnTheBurstList)
	sc.Step(`^I have a burst "([^"]*)" with (\d+) events$`, iHaveABurstWithEvents)
	sc.Step(`^I should see the burst detail modal$`, iShouldSeeTheBurstDetailModal)
	sc.Step(`^I press "v" to view events$`, iPressVToViewEvents)
	sc.Step(`^I should see the burst events modal$`, iShouldSeeTheBurstEventsModal)
	sc.Step(`^I should see event details$`, iShouldSeeEventDetails)
	sc.Step(`^I have a confirmed burst "([^"]*)" with facts$`, iHaveAConfirmedBurstWithFacts)
	sc.Step(`^I press "f" to view facts$`, iPressFToViewFacts)
	sc.Step(`^I should see the burst facts modal$`, iShouldSeeTheBurstFactsModal)
	sc.Step(`^I have a confirmed burst "([^"]*)" with skills$`, iHaveAConfirmedBurstWithSkills)
	sc.Step(`^I should see the burst skills modal$`, iShouldSeeTheBurstSkillsModal)
}

//nolint:dupl // Step registration functions look similar but register different steps.
func registerBurstEditSteps(sc *godog.ScenarioContext) {
	sc.Step(`^I should see the edit burst form$`, iShouldSeeTheEditBurstForm)
	sc.Step(`^I clear the burst name field$`, iClearTheBurstNameField)
	sc.Step(`^I enter burst name "([^"]*)"$`, iEnterBurstName)
	sc.Step(`^I submit the burst form$`, iSubmitTheBurstForm)
	sc.Step(`^the burst should have name "([^"]*)"$`, theBurstShouldHaveName)
	sc.Step(`^I tab to description field$`, iTabToDescriptionField)
	sc.Step(`^I enter burst description "([^"]*)"$`, iEnterBurstDescription)
	sc.Step(`^the burst should have description "([^"]*)"$`, theBurstShouldHaveDescription)
	sc.Step(`^I have an unconfirmed burst "([^"]*)" with (\d+) events$`, iHaveAnUnconfirmedBurst)
	sc.Step(`^I press "c" to confirm$`, iPressCToConfirm)
	sc.Step(`^I should see the confirm burst modal$`, iShouldSeeTheConfirmBurstModal)
	sc.Step(`^I confirm the action$`, iConfirmTheAction)
	sc.Step(`^the burst should not be confirmed$`, theBurstShouldNotBeConfirmed)
}

func registerBurstSuggestionSteps(sc *godog.ScenarioContext) {
	sc.Step(`^I have (\d+) unassigned events$`, iHaveUnassignedEvents)
	sc.Step(`^I press "s" to suggest bursts$`, iPressSToSuggestBursts)
	sc.Step(`^the detection completes$`, theDetectionCompletes)
	sc.Step(`^I should see the burst suggestion modal$`, iShouldSeeTheBurstSuggestionModal)
	sc.Step(`^I should see suggested burst names$`, iShouldSeeSuggestedBurstNames)
	sc.Step(`^I should see confidence scores$`, iShouldSeeConfidenceScores)
	sc.Step(`^I have burst suggestions available$`, iHaveBurstSuggestionsAvailable)
	sc.Step(`^I am on the burst suggestion modal$`, iAmOnTheBurstSuggestionModal)
	sc.Step(`^I should see different suggestions highlighted$`, iShouldSeeDifferentSuggestionsHighlighted)
	sc.Step(`^I should see the suggestion events modal$`, iShouldSeeTheSuggestionEventsModal)
	sc.Step(`^I have a confirmed burst "([^"]*)" with (\d+) events$`, iHaveAConfirmedBurstWithNEvents)
	sc.Step(`^I have skill suggestions from burst$`, iHaveSkillSuggestionsFromBurst)
	sc.Step(`^I am on the skill suggestion modal$`, iAmOnTheSkillSuggestionModalBursts)
	sc.Step(`^I should see different skills highlighted$`, iShouldSeeDifferentSkillsHighlighted)
	sc.Step(`^I should see events that led to this skill$`, iShouldSeeEventsThatLedToThisSkill)
	sc.Step(`^the skill should be marked as rejected$`, theSkillShouldBeMarkedAsRejectedBursts)
}

func registerBurstNavigationSteps(sc *godog.ScenarioContext) {
	sc.Step(`^I should be at the last burst$`, iShouldBeAtTheLastBurst)
	sc.Step(`^I should be at the first burst$`, iShouldBeAtTheFirstBurst)
	sc.Step(`^I should see different bursts$`, iShouldSeeDifferentBursts)
	sc.Step(`^I should not see the loading modal$`, iShouldNotSeeTheLoadingModal)
}

func iHaveNBurstsInMyProfile(ctx context.Context, count int) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}

	// Create N bursts with test data
	for i := 0; i < count; i++ {
		burst := fixtures.BurstFactory.MustCreate().(*career.Burst)
		burst.ID = "" // Clear ID so repo generates one
		env.AddBurst(burst)
	}

	return ctx, nil
}

func iShouldSeeAListOfBursts(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Burst"),
		gomega.ContainSubstring("Name"),
		gomega.ContainSubstring("Events"),
	))
	return nil
}

func iShouldStillBeOnTheBurstList(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Bursts"),
		gomega.ContainSubstring("Burst"),
		gomega.ContainSubstring("No bursts"),
	))
	return nil
}

func iHaveABurstWithEvents(ctx context.Context, name string, count string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}

	// Create events first
	var eventIDs []string
	for i := 0; i < parseIntOrDefault(count, 3); i++ {
		event := fixtures.EventFactory.MustCreate().(*career.Event)
		event.ID = ""
		env.AddEvent(event)
		eventIDs = append(eventIDs, event.ID)
	}

	// Create burst with those events
	burst := fixtures.Burst("", eventIDs...)
	burst.Name = name
	env.AddBurst(burst)

	return ctx, nil
}

func parseIntOrDefault(s string, def int) int {
	var i int
	if _, err := fmt.Sscanf(s, "%d", &i); err == nil {
		return i
	}
	return def
}

func iShouldSeeTheBurstDetailModal(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Detail"),
		gomega.ContainSubstring("Name"),
		gomega.ContainSubstring("Events"),
	))
	return nil
}

func iPressVToViewEvents(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('v')
	return ctx, nil
}

func iShouldSeeTheBurstEventsModal(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Events"),
		gomega.ContainSubstring("Date"),
	))
	return nil
}

func iShouldSeeEventDetails(_ context.Context) error {
	return godog.ErrPending
}

func iHaveAConfirmedBurstWithFacts(_ context.Context, _ string) (context.Context, error) {
	return nil, godog.ErrPending
}

func iPressFToViewFacts(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('f')
	return ctx, nil
}

func iShouldSeeTheBurstFactsModal(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Facts"),
		gomega.ContainSubstring("Fact"),
	))
	return nil
}

func iHaveAConfirmedBurstWithSkills(_ context.Context, _ string) (context.Context, error) {
	return nil, godog.ErrPending
}

func iShouldSeeTheBurstSkillsModal(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Skills"),
		gomega.ContainSubstring("Skill"),
	))
	return nil
}

func iShouldSeeTheEditBurstForm(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Edit"),
		gomega.ContainSubstring("Name"),
		gomega.ContainSubstring("Description"),
	))
	return nil
}

func iClearTheBurstNameField(_ context.Context) (context.Context, error) {
	return nil, godog.ErrPending
}

func iEnterBurstName(_ context.Context, _ string) (context.Context, error) {
	return nil, godog.ErrPending
}

func iSubmitTheBurstForm(_ context.Context) (context.Context, error) {
	return nil, godog.ErrPending
}

func theBurstShouldHaveName(_ context.Context, _ string) error {
	return godog.ErrPending
}

func iTabToDescriptionField(_ context.Context) (context.Context, error) {
	return nil, godog.ErrPending
}

func iEnterBurstDescription(_ context.Context, _ string) (context.Context, error) {
	return nil, godog.ErrPending
}

func theBurstShouldHaveDescription(_ context.Context, _ string) error {
	return godog.ErrPending
}

func iHaveAnUnconfirmedBurst(_ context.Context, _ string, _ int) (context.Context, error) {
	return nil, godog.ErrPending
}

func iPressCToConfirm(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('c')
	return ctx, nil
}

func iShouldSeeTheConfirmBurstModal(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Confirm"),
		gomega.ContainSubstring("confirm"),
	))
	return nil
}

func iConfirmTheAction(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.Confirm()
	return ctx, nil
}

func theBurstShouldNotBeConfirmed(_ context.Context) error {
	return godog.ErrPending
}

func iHaveUnassignedEvents(_ context.Context, _ int) (context.Context, error) {
	return nil, godog.ErrPending
}

func iPressSToSuggestBursts(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('s')
	return ctx, nil
}

func theDetectionCompletes(_ context.Context) error {
	return godog.ErrPending
}

func iShouldSeeTheBurstSuggestionModal(_ context.Context) error {
	return godog.ErrPending
}

func iShouldSeeSuggestedBurstNames(_ context.Context) error {
	return godog.ErrPending
}

func iShouldSeeConfidenceScores(_ context.Context) error {
	return godog.ErrPending
}

func iHaveBurstSuggestionsAvailable(_ context.Context) (context.Context, error) {
	return nil, godog.ErrPending
}

func iAmOnTheBurstSuggestionModal(_ context.Context) error {
	return godog.ErrPending
}

func iShouldSeeDifferentSuggestionsHighlighted(_ context.Context) error {
	return godog.ErrPending
}

func iShouldSeeTheSuggestionEventsModal(_ context.Context) error {
	return godog.ErrPending
}

func iHaveAConfirmedBurstWithNEvents(_ context.Context, _ string, _ int) (context.Context, error) {
	return nil, godog.ErrPending
}

func iHaveSkillSuggestionsFromBurst(_ context.Context) (context.Context, error) {
	return nil, godog.ErrPending
}

func iAmOnTheSkillSuggestionModalBursts(_ context.Context) error {
	return godog.ErrPending
}

func iShouldSeeDifferentSkillsHighlighted(_ context.Context) error {
	return godog.ErrPending
}

func iShouldSeeEventsThatLedToThisSkill(_ context.Context) error {
	return godog.ErrPending
}

func theSkillShouldBeMarkedAsRejectedBursts(_ context.Context) error {
	return godog.ErrPending
}

func iShouldBeAtTheLastBurst(_ context.Context) error {
	return godog.ErrPending
}

func iShouldBeAtTheFirstBurst(_ context.Context) error {
	return godog.ErrPending
}

func iShouldSeeDifferentBursts(_ context.Context) error {
	return godog.ErrPending
}

func iShouldNotSeeTheLoadingModal(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).NotTo(gomega.ContainSubstring("Loading"))
	return nil
}
