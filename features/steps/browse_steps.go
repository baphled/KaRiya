// Package steps provides BDD step definitions for Godog feature tests.
package steps

import (
	"context"
	"fmt"

	"github.com/baphled/kariya/features/support"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cucumber/godog"
	"github.com/onsi/gomega"
)

// RegisterBrowseSteps registers browse timeline step definitions with Godog.
//
//nolint:dupl // Step registration blocks share similar structure but different content.
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
	env.Tab()
	env.Confirm()
	return ctx, nil
}

func iApplyTheFilter(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
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
	env.PressKey(tea.KeyCtrlU)
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

func iShouldSeeDifferentEvents(_ context.Context) error {
	return godog.ErrPending
}

func iShouldSeeTheOriginalEvents(_ context.Context) error {
	return godog.ErrPending
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

func iShouldBeAtTheLastEvent(_ context.Context) error {
	return godog.ErrPending
}

func iShouldBeAtTheFirstEvent(_ context.Context) error {
	return godog.ErrPending
}
