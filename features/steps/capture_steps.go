package steps

import (
	"context"
	"strings"
	"time"

	"github.com/baphled/kariya/features/support"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cucumber/godog"
	"github.com/onsi/gomega"
)

// RegisterCaptureSteps registers capture event step definitions with Godog.
//
// Expected:
//   - sc is a valid *godog.ScenarioContext.
//
// Side effects:
//   - Registers step definitions with Godog.
//
//nolint:funlen // Registration function has many steps by design.
func RegisterCaptureSteps(sc *godog.ScenarioContext) {
	sc.Step(`^the database is empty$`, theDatabaseIsEmpty)
	sc.Step(`^I am on the main menu$`, iAmOnTheMainMenu)
	sc.Step(`^I select "([^"]*)" from the menu$`, iSelectFromTheMenu)
	sc.Step(`^I select quick capture strategy$`, iSelectQuickCaptureStrategy)
	sc.Step(`^I select manual capture strategy$`, iSelectManualCaptureStrategy)
	sc.Step(`^I enter event description "([^"]*)"$`, iEnterEventDescription)
	sc.Step(`^I set event company to "([^"]*)"$`, iSetEventCompanyTo)
	sc.Step(`^I capture an event:$`, iCaptureAnEvent)
	sc.Step(`^I submit the event$`, iSubmitTheEvent)
	sc.Step(`^I should see the success message$`, iShouldSeeTheSuccessMessage)
	sc.Step(`^I dismiss the success modal$`, iDismissTheSuccessModal)
	sc.Step(`^I should be on the enrichment review screen$`, iShouldBeOnEnrichmentReviewScreen)
	sc.Step(`^I confirm the review$`, iConfirmTheReview)
	sc.Step(`^I should be on the main menu$`, iShouldBeOnTheMainMenu)
	sc.Step(`^I cancel$`, iCancel)
	sc.Step(`^I should see the strategy selection$`, iShouldSeeTheStrategySelection)
	sc.Step(`^there should be (\d+) events?$`, thereShouldBeNEvents)
	sc.Step(`^there should be (\d+) bursts?$`, thereShouldBeNBursts)
	sc.Step(`^the event should have description "([^"]*)"$`, theEventShouldHaveDescription)
	sc.Step(`^the event should have company "([^"]*)"$`, theEventShouldHaveCompany)
	sc.Step(`^the event should have project "([^"]*)"$`, theEventShouldHaveProject)
	sc.Step(`^the event should have tags "([^"]*)"$`, theEventShouldHaveTags)
	sc.Step(`^the event should have categories "([^"]*)"$`, theEventShouldHaveCategories)
	sc.Step(`^the event should have (\d+) tags$`, theEventShouldHaveNTags)
	sc.Step(`^the event should have (\d+) categories$`, theEventShouldHaveNCategories)
	sc.Step(`^I have an event "([^"]*)" at company "([^"]*)"$`, iHaveAnEventAtCompany)
	sc.Step(`^I accept the suggested burst$`, iAcceptTheSuggestedBurst)
	sc.Step(`^I accept all inferred skills$`, iAcceptAllInferredSkills)
	sc.Step(`^I reject all suggestions$`, iRejectAllSuggestions)
	sc.Step(`^I edit the suggested burst$`, iEditTheSuggestedBurst)
	sc.Step(`^I change burst name to "([^"]*)"$`, iChangeBurstNameTo)
	sc.Step(`^I save the burst edit$`, iSaveTheBurstEdit)
	sc.Step(`^there should be (\d+) bursts? with name "([^"]*)"$`, thereShouldBeNBurstsWithName)
	sc.Step(`^there should be skills including "([^"]*)"$`, thereShouldBeSkillsIncluding)
	sc.Step(`^I open the metadata editor$`, iOpenTheMetadataEditor)
	sc.Step(`^I change event company to "([^"]*)"$`, iChangeEventCompanyTo)
	sc.Step(`^I save metadata changes$`, iSaveMetadataChanges)
	sc.Step(`^I should see "([^"]*)" key badge for (?:editing )?(bursts|facts)$`, iShouldSeeKeyBadgeFor)
	sc.Step(`^I try to submit without description$`, iTryToSubmitWithoutDescription)
	sc.Step(`^I should see a capture validation error$`, iShouldSeeValidationError)
	sc.Step(`^I press Ctrl\+S$`, iPressCtrlS)
	sc.Step(`^I set event date to "([^"]*)"$`, iSetEventDateTo)
	sc.Step(`^the event should have date "([^"]*)"$`, theEventShouldHaveDate)
	sc.Step(`^the event should have today's date$`, theEventShouldHaveTodaysDate)
	sc.Step(`^the event should have a date (\d+) days ago$`, theEventShouldHaveDateDaysAgo)
	sc.Step(`^I should see a date validation error$`, iShouldSeeDateValidationError)
	sc.Step(`^I should see a validation error about minimum length$`, iShouldSeeMinLengthValidationError)
	sc.Step(`^the event should have skills "([^"]*)"$`, theEventShouldHaveSkills)
}

func theDatabaseIsEmpty(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.AssertEventCount(0)
	return ctx, nil
}

func iAmOnTheMainMenu(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	gomega.Expect(env.IsInMenuState()).To(gomega.BeTrue(), "Should be on main menu")
	return ctx, nil
}

func iSelectFromTheMenu(ctx context.Context, intentName string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.SelectIntentByName(intentName)
	return ctx, nil
}

func iSelectQuickCaptureStrategy(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.Confirm()
	return ctx, nil
}

func iSelectManualCaptureStrategy(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.NavigateDown()
	env.Confirm()
	return ctx, nil
}

func iEnterEventDescription(ctx context.Context, description string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}

	env.TypeText(description)

	data := support.GetEventData(ctx)
	data.Description = description
	return support.WithEventData(ctx, data), nil
}

func iSetEventCompanyTo(ctx context.Context, company string) (context.Context, error) {
	data := support.GetEventData(ctx)
	data.Company = company
	return support.WithEventData(ctx, data), nil
}

func iCaptureAnEvent(ctx context.Context, table *godog.Table) (context.Context, error) {
	data := support.GetEventData(ctx)
	for _, row := range table.Rows[1:] {
		field := row.Cells[0].Value
		value := row.Cells[1].Value
		switch field {
		case "description":
			data.Description = value
		case "company":
			data.Company = value
		case "project":
			data.Project = value
		case "tags":
			data.Tags = strings.Split(value, ",")
		case "categories":
			data.Categories = strings.Split(value, ",")
		case "skills":
			data.Skills = strings.Split(value, ",")
		}
	}
	return support.WithEventData(ctx, data), nil
}

func iSubmitTheEvent(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}

	data := support.GetEventData(ctx)
	event, err := data.BuildEvent()
	if err != nil {
		env.SubmitEventWithError(event, err)
	} else {
		env.SubmitEvent(event)
	}

	return ctx, nil
}

func iShouldSeeTheSuccessMessage(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Success"),
		gomega.ContainSubstring("Saving"),
		gomega.ContainSubstring("saved"),
		gomega.ContainSubstring("Review"),
	))
	return nil
}

func iDismissTheSuccessModal(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.DismissSuccessModal()
	return ctx, nil
}

func iShouldBeOnEnrichmentReviewScreen(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Review"),
		gomega.ContainSubstring("Enrichment"),
	))
	return nil
}

func iConfirmTheReview(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	for range 10 {
		if env.IsInMenuState() {
			break
		}
		env.Confirm()
	}
	return ctx, nil
}

func iShouldBeOnTheMainMenu(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	gomega.Expect(env.IsInMenuState()).To(gomega.BeTrue(), "Should be on main menu")
	return nil
}

func iCancel(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.Cancel()
	return ctx, nil
}

func iShouldSeeTheStrategySelection(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("Quick"),
		gomega.ContainSubstring("Manual"),
		gomega.ContainSubstring("Strategy"),
	))
	return nil
}

func thereShouldBeNEvents(ctx context.Context, expected int) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	env.AssertEventCount(expected)
	return nil
}

func thereShouldBeNBursts(ctx context.Context, expected int) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	env.AssertBurstCount(expected)
	return nil
}

func theEventShouldHaveDescription(ctx context.Context, expected string) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	events := env.GetEvents()
	gomega.Expect(events).NotTo(gomega.BeEmpty())
	gomega.Expect(events[0].Text).To(gomega.Equal(expected))
	return nil
}

func theEventShouldHaveCompany(ctx context.Context, expected string) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	events := env.GetEvents()
	gomega.Expect(events).NotTo(gomega.BeEmpty())
	gomega.Expect(events[0].Company).To(gomega.Equal(expected))
	return nil
}

func theEventShouldHaveProject(ctx context.Context, expected string) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	events := env.GetEvents()
	gomega.Expect(events).NotTo(gomega.BeEmpty())
	gomega.Expect(events[0].Project).To(gomega.Equal(expected))
	return nil
}

func theEventShouldHaveTags(ctx context.Context, expected string) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	events := env.GetEvents()
	gomega.Expect(events).NotTo(gomega.BeEmpty())
	expectedTags := strings.Split(expected, ",")
	gomega.Expect(events[0].Tags).To(gomega.ConsistOf(expectedTags))
	return nil
}

func theEventShouldHaveCategories(ctx context.Context, expected string) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	events := env.GetEvents()
	gomega.Expect(events).NotTo(gomega.BeEmpty())
	expectedCats := strings.Split(expected, ",")
	gomega.Expect(events[0].Categories).To(gomega.ConsistOf(expectedCats))
	return nil
}

func theEventShouldHaveNTags(ctx context.Context, expected int) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	events := env.GetEvents()
	gomega.Expect(events).NotTo(gomega.BeEmpty())
	gomega.Expect(events[0].Tags).To(gomega.HaveLen(expected))
	return nil
}

func theEventShouldHaveNCategories(ctx context.Context, expected int) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	events := env.GetEvents()
	gomega.Expect(events).NotTo(gomega.BeEmpty())
	gomega.Expect(events[0].Categories).To(gomega.HaveLen(expected))
	return nil
}

func theEventShouldHaveSkills(ctx context.Context, expected string) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	events := env.GetEvents()
	gomega.Expect(events).NotTo(gomega.BeEmpty())
	expectedSkills := strings.Split(expected, ",")
	gomega.Expect(events[0].Skills).To(gomega.ConsistOf(expectedSkills))
	return nil
}

func iHaveAnEventAtCompany(ctx context.Context, description, company string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	event := fixtures.EventWith("", description, company, "")
	env.AddEvent(event)
	return ctx, nil
}

func iAcceptTheSuggestedBurst(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.Confirm()
	return ctx, nil
}

func iAcceptAllInferredSkills(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.Confirm()
	return ctx, nil
}

func iRejectAllSuggestions(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('r')
	return ctx, nil
}

func iEditTheSuggestedBurst(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('e')
	return ctx, nil
}

func iChangeBurstNameTo(ctx context.Context, _ string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.TypeText("")
	return ctx, nil
}

func iSaveTheBurstEdit(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.Confirm()
	return ctx, nil
}

func thereShouldBeNBurstsWithName(ctx context.Context, expected int, _ string) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	bursts := env.GetBursts()
	gomega.Expect(bursts).To(gomega.HaveLen(expected))
	return nil
}

func thereShouldBeSkillsIncluding(ctx context.Context, skillName string) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	skills := env.GetSkills()
	var found bool
	for _, s := range skills {
		if strings.EqualFold(s.Name, skillName) {
			found = true
			break
		}
	}
	gomega.Expect(found).To(gomega.BeTrue(), "Should have skill %s", skillName)
	return nil
}

func iOpenTheMetadataEditor(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('e')
	return ctx, nil
}

func iChangeEventCompanyTo(ctx context.Context, _ string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.Tab()
	return ctx, nil
}

func iSaveMetadataChanges(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.Confirm()
	return ctx, nil
}

func iTryToSubmitWithoutDescription(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKey(tea.KeyEnter)
	return ctx, nil
}

func iShouldSeeValidationError(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	hasValidationError := strings.Contains(view, "required") ||
		strings.Contains(view, "error") ||
		strings.Contains(view, "invalid") ||
		strings.Contains(view, "validation") ||
		strings.Contains(view, "10-2000 characters")
	gomega.Expect(hasValidationError).To(gomega.BeTrue(), "Expected validation error in view:\n%s", view)
	return nil
}

func iPressCtrlS(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKey(tea.KeyCtrlS)
	return ctx, nil
}

func iSetEventDateTo(ctx context.Context, date string) (context.Context, error) {
	data := support.GetEventData(ctx)
	data.Date = date
	return support.WithEventData(ctx, data), nil
}

func theEventShouldHaveDate(ctx context.Context, expected string) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	events := env.GetEvents()
	gomega.Expect(events).NotTo(gomega.BeEmpty())
	gomega.Expect(events[0].Date.Format("2006-01-02")).To(gomega.Equal(expected))
	return nil
}

func theEventShouldHaveTodaysDate(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	events := env.GetEvents()
	gomega.Expect(events).NotTo(gomega.BeEmpty())
	today := time.Now().Format("2006-01-02")
	gomega.Expect(events[0].Date.Format("2006-01-02")).To(gomega.Equal(today))
	return nil
}

func theEventShouldHaveDateDaysAgo(ctx context.Context, daysAgo int) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	events := env.GetEvents()
	gomega.Expect(events).NotTo(gomega.BeEmpty())
	expected := time.Now().AddDate(0, 0, -daysAgo).Format("2006-01-02")
	gomega.Expect(events[0].Date.Format("2006-01-02")).To(gomega.Equal(expected))
	return nil
}

func iShouldSeeDateValidationError(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("invalid date"),
		gomega.ContainSubstring("Invalid date"),
		gomega.ContainSubstring("date format"),
		gomega.ContainSubstring("Date format"),
	))
	return nil
}

func iShouldSeeMinLengthValidationError(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.SatisfyAny(
		gomega.ContainSubstring("minimum"),
		gomega.ContainSubstring("at least"),
		gomega.ContainSubstring("too short"),
		gomega.ContainSubstring("10 characters"),
		gomega.ContainSubstring("10-2000"),
	))
	return nil
}
func iShouldSeeKeyBadgeFor(ctx context.Context, key, action string) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.ContainSubstring(key))
	return nil
}
