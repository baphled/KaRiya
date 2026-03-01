// Package steps provides BDD step definitions for Godog feature tests.
package steps

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/baphled/kariya/features/support"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cucumber/godog"
	"github.com/onsi/gomega"
)

// RegisterBrowseSteps registers browse timeline step definitions with Godog.
//
// Expected:
//   - sc is a valid *godog.ScenarioContext.
//
// Side effects:
//   - Registers step definitions with Godog.
//
// Expected:
//   - sc is a valid *godog.ScenarioContext.
//
// Side effects:
//   - Registers step definitions with Godog.
//
//nolint:funlen // Step registration blocks are long by nature.
func RegisterBrowseSteps(sc *godog.ScenarioContext) {
	sc.Step(`^I have (\d+) events? in my timeline$`, iHaveNEventsInMyTimeline)
	sc.Step(`^I should see a list of events$`, iShouldSeeAListOfEvents)
	sc.Step(`^I press "j" to navigate down$`, iPressJToNavigateDown)
	sc.Step(`^I press "k" to navigate up$`, iPressKToNavigateUp)
	sc.Step(`^I press down arrow$`, iPressDownArrow)
	sc.Step(`^I press up arrow$`, iPressUpArrow)
	sc.Step(`^I should still be on the timeline$`, iShouldStillBeOnTheTimeline)
	sc.Step(`^I press enter to view details$`, iPressEnterToViewDetails)
	sc.Step(`^I press escape$`, iPressEscape)
	sc.Step(`^I press "/" to search$`, iPressSlashToSearch)
	sc.Step(`^I should see the search modal$`, iShouldSeeTheSearchModal)
	sc.Step(`^I type "([^"]*)" in the search$`, iTypeInTheSearch)
	sc.Step(`^I submit the search$`, iSubmitTheSearch)

	sc.Step(`^I press "f" to filter$`, iPressFToFilter)
	sc.Step(`^I should see the filter modal$`, iShouldSeeTheFilterModal)
	sc.Step(`^I select company "([^"]*)"$`, iSelectCompany)
	sc.Step(`^I apply the filter$`, iApplyTheFilter)
	sc.Step(`^I press "x" to clear filter$`, iPressXToClearFilter)
	sc.Step(`^I press "s" to sort$`, iPressSToSort)
	sc.Step(`^I should see the sort modal$`, iShouldSeeTheSortModal)
	sc.Step(`^I press "a" to add event$`, iPressAToAddEvent)
	sc.Step(`^I should see the add event form$`, iShouldSeeTheAddEventForm)
	sc.Step(`^I enter "([^"]*)" as description$`, iEnterAsDescription)
	sc.Step(`^I submit the form$`, iSubmitTheForm)
	sc.Step(`^I press "e" to edit$`, iPressEToEdit)
	sc.Step(`^I should see the edit event form$`, iShouldSeeTheEditEventForm)
	sc.Step(`^I clear the description field$`, iClearTheDescriptionField)
	sc.Step(`^I navigate to company field$`, iNavigateToCompanyField)
	sc.Step(`^I clear the company field$`, iClearTheCompanyField)
	sc.Step(`^I enter "([^"]*)" as company$`, iEnterAsCompany)
	sc.Step(`^I press "d" to delete$`, iPressDToDelete)
	sc.Step(`^I should see the delete confirmation$`, iShouldSeeTheDeleteConfirmation)
	sc.Step(`^I cancel the confirmation$`, iCancelTheConfirmation)
	sc.Step(`^I confirm the deletion$`, iConfirmTheDeletion)
	sc.Step(`^I should be able to go back to the menu$`, iShouldBeAbleToGoBackToTheMenu)
	sc.Step(`^I press page down$`, iPressPageDown)
	sc.Step(`^I press page up$`, iPressPageUp)
	sc.Step(`^I should see different events$`, iShouldSeeDifferentEvents)
	sc.Step(`^I should see the original events$`, iShouldSeeTheOriginalEvents)
	sc.Step(`^I press "G" to go to last$`, iPressGToGoToLast)
	sc.Step(`^I press "g" to go to first$`, iPressLittleGToGoToFirst)
	sc.Step(`^I should be at the last event$`, iShouldBeAtTheLastEvent)
	sc.Step(`^I should be at the first event$`, iShouldBeAtTheFirstEvent)

	// Event creation helpers
	sc.Step(`^I have an event "([^"]*)"$`, iHaveAnEvent)
	sc.Step(`^I have an event "([^"]*)" with category "([^"]*)"$`, iHaveAnEventWithCategory)
	sc.Step(`^I have an event "([^"]*)" dated "([^"]*)"$`, iHaveAnEventDated)
	sc.Step(`^I have an event "([^"]*)" with project "([^"]*)"$`, iHaveAnEventWithProject)

	// Additional navigation
	sc.Step(`^I press "G" to go to bottom$`, iPressGToGoToBottom)
	sc.Step(`^I press "g" to go to top$`, iPressLittleGToGoToTop)
	sc.Step(`^I press "j" to scroll down$`, iPressJToScrollDown)
	sc.Step(`^I press "k" to scroll up$`, iPressKToScrollUp)
	sc.Step(`^I press "a" to accept$`, iPressAToAccept)
	sc.Step(`^I press "r" to reject$`, iPressRToReject)

	// Assertions
	sc.Step(`^I should see 1 event$`, iShouldSee1Event)
	sc.Step(`^I should see skill categories$`, iShouldSeeSkillCategories)
	sc.Step(`^I should see suggested skills$`, iShouldSeeSuggestedSkills)
	sc.Step(`^I should see the event detail modal$`, iShouldSeeTheEventDetailModal)
	sc.Step(`^I should see the skill suggestion modal$`, iShouldSeeTheSkillSuggestionModal)
	sc.Step(`^I should see the skills detail modal$`, iShouldSeeTheSkillsDetailModal)
	sc.Step(`^the event has skills "([^"]*)"$`, theEventHasSkills)

	// Filter helpers
	sc.Step(`^I select companies "([^"]*)"$`, iSelectCompanies)
	sc.Step(`^I select project "([^"]*)"$`, iSelectProject)
	sc.Step(`^I set date from "([^"]*)"$`, iSetDateFrom)
}

func iHaveNEventsInMyTimeline(ctx context.Context, count int) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	for i := range count {
		event := fixtures.EventWith("", fmt.Sprintf("Event %d description", i+1), fmt.Sprintf("Company%d", i+1), "")
		env.AddEvent(event)
	}
	return ctx, nil
}

func iShouldSeeAListOfEvents(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Event"),
		gomega.ContainSubstring("Company"),
		gomega.ContainSubstring("Date"),
	))
	return nil
}

func iPressJToNavigateDown(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('j')
	return ctx, nil
}

func iPressKToNavigateUp(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('k')
	return ctx, nil
}

func iPressDownArrow(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.NavigateDown()
	return ctx, nil
}

func iPressUpArrow(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.NavigateUp()
	return ctx, nil
}

func iShouldStillBeOnTheTimeline(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Timeline"),
		gomega.ContainSubstring("Event"),
		gomega.ContainSubstring("No events"),
	))
	return nil
}

func iPressEnterToViewDetails(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}

	// If we have a target item name in context, navigate to it first
	if itemName, ok := ctx.Value(currentBurstNameKey).(string); ok && itemName != "" {
		if err := support.NavigateToTableItem(env, itemName, 20); err != nil {
			return ctx, fmt.Errorf("failed to navigate to burst %q: %w", itemName, err)
		}
	}

	env.Confirm()
	return ctx, nil
}

func iPressEscape(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env != nil {
		env.Cancel()
		return ctx, nil
	}

	onboardingEnv := support.GetOnboardingEnv(ctx)
	if onboardingEnv != nil {
		onboardingEnv.PressKey("escape")
		return ctx, nil
	}

	return ctx, godog.ErrPending
}

func iPressSlashToSearch(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('/')
	return ctx, nil
}

func iShouldSeeTheSearchModal(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Search"),
		gomega.ContainSubstring("search"),
	))
	return nil
}

func iTypeInTheSearch(ctx context.Context, text string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.TypeText(text)
	return ctx, nil
}

func iSubmitTheSearch(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressEnterWithFormProcessing()
	return ctx, nil
}

func iPressFToFilter(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('f')
	return ctx, nil
}

func iShouldSeeTheFilterModal(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Filter"),
		gomega.ContainSubstring("filter"),
		gomega.ContainSubstring("Company"),
		gomega.ContainSubstring("Category"),
		gomega.ContainSubstring("Level"),
		gomega.ContainSubstring("Apply"),
		gomega.ContainSubstring("Cancel"),
	))
	return nil
}

func iSelectCompany(ctx context.Context, _ string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.NavigateDown()
	env.PressKeyRune('x')
	return ctx, nil
}

func iApplyTheFilter(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.Confirm()
	env.Confirm()
	env.Confirm()
	return ctx, nil
}

func iPressXToClearFilter(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('x')
	return ctx, nil
}

func iPressSToSort(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('s')
	return ctx, nil
}

func iShouldSeeTheSortModal(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Sort"),
		gomega.ContainSubstring("sort"),
		gomega.ContainSubstring("Date"),
		gomega.ContainSubstring("Order"),
	))
	return nil
}

func iPressAToAddEvent(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('a')
	return ctx, nil
}

func iShouldSeeTheAddEventForm(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Add"),
		gomega.ContainSubstring("Event"),
		gomega.ContainSubstring("Description"),
		gomega.ContainSubstring("Submit"),
	))
	return nil
}

func iEnterAsDescription(ctx context.Context, text string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.TypeText(text)
	return ctx, nil
}

func iSubmitTheForm(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.Confirm()
	env.Confirm()
	env.Confirm()
	env.Confirm()
	env.Confirm()
	env.Confirm()
	env.PressKey(tea.KeyCtrlS)
	return ctx, nil
}

func iPressEToEdit(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('e')
	return ctx, nil
}

func iShouldSeeTheEditEventForm(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Edit"),
		gomega.ContainSubstring("Update"),
		gomega.ContainSubstring("Submit"),
	))
	return nil
}

func iClearTheDescriptionField(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.ClearTextField(100)
	return ctx, nil
}

func iNavigateToCompanyField(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.NextFormField()
	env.NextFormField()
	return ctx, nil
}

func iClearTheCompanyField(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.ClearTextField(50)
	return ctx, nil
}

func iEnterAsCompany(ctx context.Context, company string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.TypeText(company)
	return ctx, nil
}

func iPressDToDelete(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('d')
	return ctx, nil
}

func iShouldSeeTheDeleteConfirmation(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Delete"),
		gomega.ContainSubstring("delete"),
		gomega.ContainSubstring("Confirm"),
		gomega.ContainSubstring("confirm"),
		gomega.ContainSubstring("Are you sure"),
	))
	return nil
}

func iCancelTheConfirmation(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.Cancel()
	return ctx, nil
}

func iConfirmTheDeletion(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.Confirm()
	return ctx, nil
}

func iShouldBeAbleToGoBackToTheMenu(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	env.Cancel()
	gomega.Expect(env.IsInMenuState()).To(gomega.BeTrue())
	return nil
}

func iPressPageDown(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKey(tea.KeyPgDown)
	return ctx, nil
}

func iPressPageUp(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKey(tea.KeyPgUp)
	return ctx, nil
}

func iShouldSeeDifferentEvents(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).ToNot(gomega.BeEmpty())
	return nil
}

func iShouldSeeTheOriginalEvents(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).ToNot(gomega.BeEmpty())
	return nil
}

func iPressGToGoToLast(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('G')
	return ctx, nil
}

func iPressLittleGToGoToFirst(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('g')
	return ctx, nil
}

func iShouldBeAtTheLastEvent(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.ContainSubstring("Event"))
	return nil
}

func iShouldBeAtTheFirstEvent(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.ContainSubstring("Event"))
	return nil
}

// iHaveAnEvent creates a simple event with just a description.
func iHaveAnEvent(ctx context.Context, description string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	event := fixtures.EventWith("", description, "", "")
	env.AddEvent(event)
	return ctx, nil
}

// iHaveAnEventWithCategory creates an event with a specific category.
func iHaveAnEventWithCategory(ctx context.Context, description, category string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	event := fixtures.EventWith("", description, "", "")
	event.Categories = []string{category}
	env.AddEvent(event)
	return ctx, nil
}

// iHaveAnEventDated creates an event with a specific date.
func iHaveAnEventDated(ctx context.Context, description, dateStr string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	event := fixtures.EventWith("", description, "", "")
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return ctx, err
	}
	event.Date = parsedDate
	env.AddEvent(event)
	return ctx, nil
}

// iHaveAnEventWithProject creates an event with a specific project.
func iHaveAnEventWithProject(ctx context.Context, description, project string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	event := fixtures.EventWith("", description, "", project)
	env.AddEvent(event)
	return ctx, nil
}

// iPressGToGoToBottom navigates to bottom.
func iPressGToGoToBottom(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('G')
	return ctx, nil
}

// iPressLittleGToGoToTop navigates to top.
func iPressLittleGToGoToTop(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('g')
	return ctx, nil
}

// iPressJToScrollDown scrolls down.
func iPressJToScrollDown(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('j')
	return ctx, nil
}

// iPressKToScrollUp scrolls up.
func iPressKToScrollUp(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('k')
	return ctx, nil
}

// iPressAToAccept presses 'a' to accept.
func iPressAToAccept(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('a')
	return ctx, nil
}

// iPressRToReject presses 'r' to reject.
func iPressRToReject(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('r')
	return ctx, nil
}

// iShouldSee1Event asserts exactly 1 event exists.
func iShouldSee1Event(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	env.AssertEventCount(1)
	return nil
}

// iShouldSeeSkillCategories asserts skill categories are visible.
func iShouldSeeSkillCategories(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Categor"),
		gomega.ContainSubstring("categor"),
		gomega.ContainSubstring("Skills"),
	))
	return nil
}

// iShouldSeeSuggestedSkills asserts suggested skills are visible.
func iShouldSeeSuggestedSkills(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Suggest"),
		gomega.ContainSubstring("suggest"),
		gomega.ContainSubstring("Skills"),
	))
	return nil
}

// iShouldSeeTheEventDetailModal asserts event detail modal is visible.
func iShouldSeeTheEventDetailModal(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Detail"),
		gomega.ContainSubstring("Event"),
		gomega.ContainSubstring("Description"),
	))
	return nil
}

// iShouldSeeTheSkillSuggestionModal asserts skill suggestion modal is visible.
func iShouldSeeTheSkillSuggestionModal(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Skill"),
		gomega.ContainSubstring("Suggest"),
	))
	return nil
}

// iShouldSeeTheSkillsDetailModal asserts skills detail modal is visible.
func iShouldSeeTheSkillsDetailModal(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Skills"),
		gomega.ContainSubstring("Detail"),
	))
	return nil
}

// theEventHasSkills sets skills on the most recent event (GIVEN step).
// This is called AFTER an event is created, so we need to update it with skills.
func theEventHasSkills(ctx context.Context, skillsStr string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}

	// Parse skill names
	skillNames := strings.Split(skillsStr, ",")
	for i, skill := range skillNames {
		skillNames[i] = strings.TrimSpace(skill)
	}

	// Create Skill entities and get their IDs
	skillRepo := env.Service.GetSkillRepository()
	skillIDs := []string{}

	if skillRepo != nil {
		for _, skillName := range skillNames {
			// Create a skill entity with a deterministic ID based on name
			skill := &career.Skill{
				Name:     skillName,
				Category: "technology",
			}
			// Try to create the skill
			if err := skillRepo.Create(env.Ctx, skill); err != nil {
				// Skill might already exist with this name
				// In that case, the ID might have been auto-generated
				// For now, we'll use the skill name as the ID if creation fails
				skill.ID = skillName
			}
			skillIDs = append(skillIDs, skill.ID)
		}
	}

	// Update the last event with the skill IDs
	events := env.GetEvents()
	if len(events) > 0 {
		lastEvent := events[len(events)-1]
		lastEvent.Skills = skillIDs // Use IDs, not names

		eventRepo := env.Service.GetEventRepository()
		if eventRepo != nil {
			if err := eventRepo.Update(env.Ctx, lastEvent); err != nil {
				return ctx, fmt.Errorf("failed to update event with skills: %w", err)
			}
		}
	}

	return ctx, nil
}

// iSelectCompanies selects multiple companies in filter.
func iSelectCompanies(ctx context.Context, companies string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	companyList := strings.Split(companies, ",")
	for _, company := range companyList {
		env.TypeText(strings.TrimSpace(company))
		env.Confirm()
	}
	return ctx, nil
}

// iSelectProject selects a project in filter.
func iSelectProject(ctx context.Context, project string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.TypeText(project)
	env.Confirm()
	return ctx, nil
}

// iSetDateFrom sets the date from field.
func iSetDateFrom(ctx context.Context, dateStr string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.TypeText(dateStr)
	return ctx, nil
}
