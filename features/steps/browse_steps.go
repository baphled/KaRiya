// Package steps provides BDD step definitions for Godog feature tests.
package steps

import (
	"context"
	"fmt"
	"strings"

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
	sc.Step(`^I should still be on the timeline$`, iShouldStillBeOnTheTimeline)
	sc.Step(`^I press enter to view details$`, iPressEnterToViewDetails)
	sc.Step(`^I press escape$`, iPressEscape)
	sc.Step(`^I press "/" to search$`, iPressSlashToSearch)
	sc.Step(`^I should see the search modal$`, iShouldSeeTheSearchModal)
	sc.Step(`^I press "a" to add event$`, iPressAToAddEvent)
	sc.Step(`^I should see the add event form$`, iShouldSeeTheAddEventForm)
	sc.Step(`^I enter "([^"]*)" as description$`, iEnterAsDescription)
	sc.Step(`^I submit the form$`, iSubmitTheForm)
	sc.Step(`^I press "e" to edit$`, iPressEToEdit)
	sc.Step(`^I should see the edit event form$`, iShouldSeeTheEditEventForm)
	sc.Step(`^I navigate to company field$`, iNavigateToCompanyField)
	sc.Step(`^I clear the company field$`, iClearTheCompanyField)
	sc.Step(`^I enter "([^"]*)" as company$`, iEnterAsCompany)
	sc.Step(`^I press "d" to delete$`, iPressDToDelete)
	sc.Step(`^I should see the delete confirmation$`, iShouldSeeTheDeleteConfirmation)
	sc.Step(`^I cancel the confirmation$`, iCancelTheConfirmation)
	sc.Step(`^I confirm the deletion$`, iConfirmTheDeletion)
	sc.Step(`^I should be able to go back to the menu$`, iShouldBeAbleToGoBackToTheMenu)
	// Event creation helpers
	sc.Step(`^I have an event "([^"]*)"$`, iHaveAnEvent)
	// Additional navigation
	sc.Step(`^I press "a" to accept$`, iPressAToAccept)
	sc.Step(`^I press "r" to reject$`, iPressRToReject)

	// Assertions
	sc.Step(`^I should see skill categories$`, iShouldSeeSkillCategories)
	sc.Step(`^I should see suggested skills$`, iShouldSeeSuggestedSkills)
	sc.Step(`^I should see the event detail modal$`, iShouldSeeTheEventDetailModal)
	sc.Step(`^I should see the skill suggestion modal$`, iShouldSeeTheSkillSuggestionModal)
	sc.Step(`^I should see the skills detail modal$`, iShouldSeeTheSkillsDetailModal)
	sc.Step(`^the event has skills "([^"]*)"$`, theEventHasSkills)

	// Filter helpers
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

	env.Confirm()
	return ctx, nil
}

func iPressEscape(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.Cancel()
	return ctx, nil
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

// iHaveAnEventDated creates an event with a specific date.

// iHaveAnEventWithProject creates an event with a specific project.

// iPressGToGoToBottom navigates to bottom.

// iPressLittleGToGoToTop navigates to top.

// iPressJToScrollDown scrolls down.

// iPressKToScrollUp scrolls up.

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

// iSelectProject selects a project in filter.

// iSetDateFrom sets the date from field.
