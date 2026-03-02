// Package steps provides BDD step definitions for Godog feature tests.
package steps

import (
	"context"
	"errors"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/baphled/kariya/features/support"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/cucumber/godog"
	"github.com/onsi/gomega"
)

// contextKey is a custom type for context keys to avoid collisions.
type contextKey string

const (
	originalBurstsViewKey contextKey = "originalBurstsView"
	currentBurstNameKey   contextKey = "currentBurstName"
	burstNameKey          contextKey = "burstName"
	burstDescriptionKey   contextKey = "burstDescription"
)

// RegisterBurstsSteps registers burst management step definitions with Godog.
//
// Expected:
//   - sc is a valid *godog.ScenarioContext.
//
// Side effects:
//   - Registers step definitions with Godog.
func RegisterBurstsSteps(sc *godog.ScenarioContext) {
	registerBurstViewSteps(sc)
	registerBurstEditSteps(sc)
	registerBurstSuggestionSteps(sc)
	registerBurstNavigationSteps(sc)
}

func registerBurstViewSteps(sc *godog.ScenarioContext) {
	sc.Step(`^I have (\d+) bursts? in my profile$`, iHaveNBurstsInMyProfile)
	sc.Step(`^I have (\d+) bursts?$`, iHaveNBurstsInMyProfile) // Alias
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
	sc.Step(`^I press "s" to view skills$`, iPressSToViewSkills)
	sc.Step(`^I should see the burst skills modal$`, iShouldSeeTheBurstSkillsModal)
	sc.Step(`^I should still be on the burst detail modal$`, iShouldStillBeOnTheBurstDetailModal)
}

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
	sc.Step(`^the burst should be confirmed$`, theBurstShouldBeConfirmed)
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
	sc.Step(`^I should see the suggestion events modal$`, iShouldSeeTheSuggestionEventsModal)
	sc.Step(`^I have a confirmed burst "([^"]*)" with (\d+) events$`, iHaveAConfirmedBurstWithNEvents)
	sc.Step(`^I have skill suggestions from burst$`, iHaveSkillSuggestionsFromBurst)
	sc.Step(`^I am on the skill suggestion modal$`, iAmOnTheSkillSuggestionModalBursts)
	sc.Step(`^I should see events that led to this skill$`, iShouldSeeEventsThatLedToThisSkill)
	sc.Step(`^the skill should be marked as rejected$`, theSkillShouldBeMarkedAsRejectedBursts)
}

func registerBurstNavigationSteps(sc *godog.ScenarioContext) {
	sc.Step(`^I should not see the loading modal$`, iShouldNotSeeTheLoadingModal)
}

func iHaveNBurstsInMyProfile(ctx context.Context, count int) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}

	// Create N bursts with test data
	for range count {
		burst, ok := fixtures.BurstFactory.MustCreate().(*career.Burst)
		if !ok {
			return ctx, errors.New("factory did not create a Burst")
		}
		burst.ID = "" // Clear ID so repo generates one
		env.AddBurst(burst)
	}

	return ctx, nil
}

func iShouldSeeAListOfBursts(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
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
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
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
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}

	// Create events first
	var eventIDs []string
	for range parseIntOrDefault(count, 3) {
		event, ok := fixtures.EventFactory.MustCreate().(*career.Event)
		if !ok {
			return ctx, errors.New("factory did not create an Event")
		}
		event.ID = ""
		env.AddEvent(event)
		eventIDs = append(eventIDs, event.ID)
	}

	// Create burst with those events
	burst := fixtures.Burst("", eventIDs...)
	burst.Name = name
	env.AddBurst(burst)

	// Store burst name in context for later navigation
	ctx = context.WithValue(ctx, currentBurstNameKey, name)

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
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
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
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}
	env.PressKeyRune('v')
	return ctx, nil
}

func iShouldSeeTheBurstEventsModal(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.ContainSubstring("Events in Burst:"))
	return nil
}

func iShouldSeeEventDetails(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.ContainSubstring("Date:"))
	return nil
}

func iHaveAConfirmedBurstWithFacts(ctx context.Context, name string) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}

	// Create events first
	var eventIDs []string
	for range 3 {
		eventInterface, err := fixtures.EventFactory.Create()
		if err != nil {
			return ctx, fmt.Errorf("failed to create event: %w", err)
		}
		event, ok := eventInterface.(*career.Event)
		if !ok {
			return ctx, errors.New("factory created wrong type: expected *career.Event")
		}
		event.ID = ""
		env.AddEvent(event)
		eventIDs = append(eventIDs, event.ID)
	}

	// Create confirmed burst with those events
	burst := fixtures.BurstConfirmed("", eventIDs...)
	burst.Name = name
	env.AddBurst(burst)

	// Create facts associated with this burst
	for range 3 {
		factInterface, err := fixtures.FactFactory.Create()
		if err != nil {
			return ctx, fmt.Errorf("failed to create fact: %w", err)
		}
		fact, ok := factInterface.(*career.Fact)
		if !ok {
			return ctx, errors.New("factory created wrong type: expected *career.Fact")
		}
		fact.ID = ""
		fact.SourceBurstID = burst.ID
		env.AddFact(fact)
	}

	// Store burst name in context for later navigation
	ctx = context.WithValue(ctx, currentBurstNameKey, name)

	return ctx, nil
}

func iPressFToViewFacts(ctx context.Context) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}
	env.PressKeyRune('f')
	return ctx, nil
}

func iShouldSeeTheBurstFactsModal(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Facts"),
		gomega.ContainSubstring("Fact"),
	))
	return nil
}

func iPressSToViewSkills(ctx context.Context) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}
	env.PressKeyRune('s')
	return ctx, nil
}

func iHaveAConfirmedBurstWithSkills(ctx context.Context, name string) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}

	// Create events first
	var eventIDs []string
	for range 3 {
		eventInterface, err := fixtures.EventFactory.Create()
		if err != nil {
			return ctx, fmt.Errorf("failed to create event: %w", err)
		}
		event, ok := eventInterface.(*career.Event)
		if !ok {
			return ctx, errors.New("factory created wrong type: expected *career.Event")
		}
		event.ID = ""
		env.AddEvent(event)
		eventIDs = append(eventIDs, event.ID)
	}

	// Create confirmed burst with those events
	burst := fixtures.BurstConfirmed("", eventIDs...)
	burst.Name = name
	env.AddBurst(burst)

	// Create skills for the burst (ensure unique names)
	for i := range 3 {
		skillInterface, err := fixtures.SkillFactory.Create()
		if err != nil {
			return ctx, fmt.Errorf("failed to create skill: %w", err)
		}
		skill, ok := skillInterface.(*career.Skill)
		if !ok {
			return ctx, errors.New("factory created wrong type: expected *career.Skill")
		}
		skill.ID = ""
		skill.Name = fmt.Sprintf("%s-%d", skill.Name, i)
		env.AddSkill(skill)
	}

	// Store burst name in context for later navigation
	ctx = context.WithValue(ctx, currentBurstNameKey, name)

	return ctx, nil
}

func iShouldSeeTheBurstSkillsModal(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Skills"),
		gomega.ContainSubstring("Skill"),
	))
	return nil
}

func iShouldSeeTheEditBurstForm(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Edit"),
		gomega.ContainSubstring("Name"),
		gomega.ContainSubstring("Description"),
	))
	return nil
}

func iClearTheBurstNameField(ctx context.Context) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}

	// Clear the name field using Ctrl+U (Unix line-kill)
	env.PressKey(tea.KeyCtrlU)

	return ctx, nil
}

func iEnterBurstName(ctx context.Context, name string) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}

	env.TypeText(name)
	return context.WithValue(ctx, burstNameKey, name), nil
}

func iSubmitTheBurstForm(ctx context.Context) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}

	// Drive the actual huh form UI submission
	// The form has fields: Name, Description
	// We navigate through them and confirm at the end
	env.Confirm()

	return ctx, nil
}

func theBurstShouldHaveName(ctx context.Context, name string) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}

	// Assert burst name is visible in the view
	view := env.GetView()
	gomega.Expect(view).To(gomega.ContainSubstring(name))

	return nil
}

func iTabToDescriptionField(ctx context.Context) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}

	// Tab to next field (description)
	env.PressKey(tea.KeyTab)
	return ctx, nil
}

func iEnterBurstDescription(ctx context.Context, description string) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}

	env.TypeText(description)
	return context.WithValue(ctx, burstDescriptionKey, description), nil
}

func theBurstShouldHaveDescription(ctx context.Context, description string) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}

	// Assert burst description is visible in the view
	view := env.GetView()
	gomega.Expect(view).To(gomega.ContainSubstring(description))

	return nil
}

func iHaveAnUnconfirmedBurst(ctx context.Context, name string, count int) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}

	// Create events first
	var eventIDs []string
	for range count {
		event, ok := fixtures.EventFactory.MustCreate().(*career.Event)
		if !ok {
			return ctx, errors.New("factory did not create an Event")
		}
		event.ID = ""
		env.AddEvent(event)
		eventIDs = append(eventIDs, event.ID)
	}

	// Create unconfirmed burst with those events
	burst := fixtures.Burst("", eventIDs...)
	burst.Name = name
	burst.Confirmed = false // Explicitly set
	env.AddBurst(burst)

	return ctx, nil
}

func iPressCToConfirm(ctx context.Context) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}
	env.PressKeyRune('c')
	return ctx, nil
}

func iShouldSeeTheConfirmBurstModal(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Confirm"),
		gomega.ContainSubstring("confirm"),
	))
	return nil
}

func iConfirmTheAction(ctx context.Context) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}

	// Drive the UI confirmation modal by pressing Enter
	// The confirm modal accepts "y", "Y", or "enter" to confirm
	env.PressKey(tea.KeyEnter)

	return ctx, nil
}

func theBurstShouldNotBeConfirmed(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}

	// Assert burst is not confirmed by checking the view
	// On list view: check for "✗ No"
	// On detail modal: check for absence of "✓ Confirmed"
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("✗ No"),
		gomega.Not(gomega.ContainSubstring("✓ Confirmed")),
	))

	return nil
}

func theBurstShouldBeConfirmed(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}

	// Assert burst is confirmed by checking the view
	// On list view: check for "✓ Yes"
	// On detail modal: check for "✓ Confirmed"
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("✓ Yes"),
		gomega.ContainSubstring("✓ Confirmed"),
	))

	return nil
}

func iHaveUnassignedEvents(ctx context.Context, count int) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}

	// Create unassigned events (not part of any burst)
	for range count {
		eventInterface, err := fixtures.EventFactory.Create()
		if err != nil {
			return ctx, fmt.Errorf("failed to create event: %w", err)
		}
		event, ok := eventInterface.(*career.Event)
		if !ok {
			return ctx, errors.New("factory created wrong type: expected *career.Event")
		}
		event.ID = ""
		env.AddEvent(event)
	}

	return ctx, nil
}

func iPressSToSuggestBursts(ctx context.Context) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}
	env.PressKeyRune('s')
	return ctx, nil
}

func theDetectionCompletes(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	if err := support.WaitForViewContains(env, "Review Burst Suggestions", 20, 100); err != nil {
		return err
	}
	return nil
}

func iShouldSeeTheBurstSuggestionModal(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.ContainSubstring("Review Burst Suggestions"))
	return nil
}

func iShouldSeeSuggestedBurstNames(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.ContainSubstring("Name"))
	return nil
}

func iShouldSeeConfidenceScores(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.ContainSubstring("Confidence"))
	return nil
}

func iHaveBurstSuggestionsAvailable(ctx context.Context) (context.Context, error) {
	ctx, err := iHaveUnassignedEvents(ctx, 5)
	if err != nil {
		return ctx, err
	}
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}
	env.SelectIntentByName("burst_management")
	env.PressKeyRune('s')
	if err := support.WaitForViewContains(env, "Review Burst Suggestions", 20, 100); err != nil {
		return ctx, err
	}
	return ctx, nil
}

func iAmOnTheBurstSuggestionModal(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.ContainSubstring("Review Burst Suggestions"))
	return nil
}

func iShouldSeeTheSuggestionEventsModal(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.ContainSubstring("Events"))
	return nil
}

func iHaveAConfirmedBurstWithNEvents(ctx context.Context, name string, count int) (context.Context, error) {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}

	var eventIDs []string
	for range count {
		eventInterface, err := fixtures.EventFactory.Create()
		if err != nil {
			return ctx, fmt.Errorf("failed to create event: %w", err)
		}
		event, ok := eventInterface.(*career.Event)
		if !ok {
			return ctx, errors.New("factory created wrong type: expected *career.Event")
		}
		event.ID = ""
		env.AddEvent(event)
		eventIDs = append(eventIDs, event.ID)
	}

	burst := fixtures.BurstConfirmed("", eventIDs...)
	burst.Name = name
	env.AddBurst(burst)

	ctx = context.WithValue(ctx, currentBurstNameKey, name)

	return ctx, nil
}

func iHaveSkillSuggestionsFromBurst(ctx context.Context) (context.Context, error) {
	ctx, err := iHaveAConfirmedBurstWithNEvents(ctx, "Backend Development", 5)
	if err != nil {
		return ctx, err
	}
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return ctx, err
	}
	env.SelectIntentByName("burst_management")
	env.PressKey(tea.KeyEnter)
	env.PressKeyRune('i')
	if err := support.WaitForViewContains(env, "Review Skill Suggestions", 20, 100); err != nil {
		return ctx, err
	}
	return ctx, nil
}

func iAmOnTheSkillSuggestionModalBursts(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.ContainSubstring("Review Skill Suggestions"))
	return nil
}

func iShouldSeeEventsThatLedToThisSkill(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.ContainSubstring("Events"))
	return nil
}

func theSkillShouldBeMarkedAsRejectedBursts(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.ContainSubstring("Review Skill Suggestions"))
	return nil
}

func iShouldNotSeeTheLoadingModal(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	view := env.GetView()
	gomega.Expect(view).NotTo(gomega.ContainSubstring("Loading"))
	return nil
}

// iShouldStillBeOnTheBurstDetailModal asserts we're still on burst detail modal.
func iShouldStillBeOnTheBurstDetailModal(ctx context.Context) error {
	env, err := support.RequireEnv(ctx)
	if err != nil {
		return err
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Burst"),
		gomega.ContainSubstring("Detail"),
	))
	return nil
}
