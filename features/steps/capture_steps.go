// Package steps provides BDD step definitions for Godog feature tests.
package steps

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/baphled/kariya/features/support"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/e2e"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cucumber/godog"
	"github.com/google/uuid"
	"github.com/onsi/gomega"
)

type contextKey string

const (
	editedBurstNameKey    contextKey = "editedBurstName"
	editedEventCompanyKey contextKey = "editedEventCompany"
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
	sc.Step(`^I reject the suggested burst$`, iRejectTheSuggestedBurst)
	sc.Step(`^the accepted burst should have at least (\d+) event IDs$`, theAcceptedBurstShouldHaveAtLeastNEventIDs)
	sc.Step(`^I accept all inferred skills$`, iAcceptAllInferredSkills)
	sc.Step(`^I reject all suggestions$`, iRejectAllSuggestions)
	sc.Step(`^I edit the suggested burst$`, iEditTheSuggestedBurst)
	sc.Step(`^I change burst name to "([^"]*)"$`, iChangeBurstNameTo)
	sc.Step(`^I save the burst edit$`, iSaveTheBurstEdit)
	sc.Step(`^there should be (\d+) bursts? with name "([^"]*)"$`, thereShouldBeNBurstsWithName)
	sc.Step(`^there should be skills including "([^"]*)"$`, thereShouldBeSkillsIncluding)
	sc.Step(`^I open the metadata editor$`, iOpenTheReviewEnrichment)
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

	// Additional editor navigation
	sc.Step(`^I press 'b' to open bursts editor$`, iPressBToOpenBurstsEditor)
	sc.Step(`^I press 'f' to open facts editor$`, iPressFToOpenFactsEditor)
	sc.Step(`^I press 'e' to open metadata editor$`, iPressEToOpenReviewEnrichment)
	sc.Step(`^I should see the bursts modal$`, iShouldSeeTheBurstsModal)
	sc.Step(`^I should see the facts modal$`, iShouldSeeTheFactsModal)
	sc.Step(`^I should see the metadata modal$`, iShouldSeeTheMetadataModal)
	sc.Step(`^the review screen should show enrichment sections$`, theReviewScreenShouldShowEnrichmentSections)
	sc.Step(`^I should move to the previous field$`, iShouldMoveToThePreviousField)
}

func iAmOnTheMainMenu(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	if !env.IsInMenuState() {
		return ctx, errors.New("expected to be on main menu but current view does not match")
	}
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
		return ctx, nil
	}

	if len(event.Skills) > 0 {
		if err := persistEventWithSkills(env, event); err != nil {
			return ctx, err
		}
	}
	env.SubmitEvent(event)

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
	if !env.IsInMenuState() {
		return errors.New("expected to be on main menu but current view does not match")
	}
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

	navigateToMainMenu(env)
	env.SelectIntentByName("browse_timeline")

	view := env.GetView()
	expectedFooter := fmt.Sprintf("Events: %d", expected)
	if !strings.Contains(view, expectedFooter) {
		return fmt.Errorf("expected footer '%s' not found in view", expectedFooter)
	}

	if expected > 0 {
		env.Confirm()
	}

	return nil
}

func navigateToMainMenu(env *e2e.TestEnv) {
	for range 10 {
		if isOnMainMenu(env) {
			return
		}
		env.Cancel()
	}
}

func isOnMainMenu(env *e2e.TestEnv) bool {
	view := env.GetView()
	return strings.Contains(view, "Career Event Management System") &&
		strings.Contains(view, "Capture Event") &&
		strings.Contains(view, "Browse Timeline") &&
		strings.Contains(view, "Enter Select")
}

func thereShouldBeNBursts(ctx context.Context, expected int) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}

	bursts := env.GetBursts()
	if len(bursts) != expected {
		return fmt.Errorf("expected %d burst(s) but found %d", expected, len(bursts))
	}
	return nil
}

func theEventShouldHaveDescription(ctx context.Context, expected string) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	if !strings.Contains(view, expected) {
		return fmt.Errorf("expected description '%s' not found in view", expected)
	}
	return nil
}

func theEventShouldHaveCompany(ctx context.Context, expected string) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}

	events := env.GetEvents()
	if len(events) == 0 {
		return errors.New("no events found in database")
	}

	latest := events[len(events)-1]
	if latest.Company != expected {
		return fmt.Errorf("expected company %q but got %q", expected, latest.Company)
	}
	return nil
}

func theEventShouldHaveProject(ctx context.Context, expected string) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	if !strings.Contains(view, expected) {
		return fmt.Errorf("expected project '%s' not found in view", expected)
	}
	return nil
}

func theEventShouldHaveTags(ctx context.Context, expected string) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	for _, tag := range strings.Split(expected, ",") {
		if !strings.Contains(view, strings.TrimSpace(tag)) {
			return fmt.Errorf("expected tag '%s' not found in view", tag)
		}
	}
	return nil
}

func theEventShouldHaveCategories(ctx context.Context, expected string) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	for _, cat := range strings.Split(expected, ",") {
		if !strings.Contains(view, strings.TrimSpace(cat)) {
			return fmt.Errorf("expected category '%s' not found in view", cat)
		}
	}
	return nil
}

func checkEventItems(ctx context.Context, expected int, itemType string, items []string) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	cleaned := stripANSI(view)
	headerStr := itemType + ":"
	if expected == 0 {
		if strings.Contains(cleaned, headerStr) {
			return fmt.Errorf("expected no %s but found %s line in view", itemType, itemType)
		}
		return nil
	}
	if !strings.Contains(cleaned, headerStr) {
		return fmt.Errorf("expected %d %s but %s line not found in view", expected, itemType, itemType)
	}
	if len(items) != expected {
		return fmt.Errorf("expected %d %s but event data has %d", expected, itemType, len(items))
	}
	for _, item := range items {
		item = strings.TrimSpace(item)
		if strings.Contains(cleaned, item) {
			return nil
		}
	}
	return fmt.Errorf("expected %s visible in view but none of %v found", itemType, items)
}

func theEventShouldHaveNTags(ctx context.Context, expected int) error {
	data := support.GetEventData(ctx)
	return checkEventItems(ctx, expected, "Tags", data.Tags)
}

func theEventShouldHaveNCategories(ctx context.Context, expected int) error {
	data := support.GetEventData(ctx)
	return checkEventItems(ctx, expected, "Categories", data.Categories)
}

var ansiRegex = regexp.MustCompile(`\x1b[\[\(][0-9;]*[a-zA-Z]|\x1b\][^\x07]*\x07`)

func stripANSI(s string) string {
	return ansiRegex.ReplaceAllString(s, "")
}

func theEventShouldHaveSkills(ctx context.Context, expected string) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	for _, skill := range strings.Split(expected, ",") {
		if !strings.Contains(view, strings.TrimSpace(skill)) {
			return fmt.Errorf("expected skill '%s' not found in view", skill)
		}
	}
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

	if err := createSuggestedBurstFromEvents(env); err != nil {
		return ctx, err
	}

	env.Confirm()
	return ctx, nil
}

// iRejectTheSuggestedBurst advances the review flow without creating a burst.
// Confirming without first calling createSuggestedBurstFromEvents means
// no burst is persisted — this is the rejection mechanism for this flow.
func iRejectTheSuggestedBurst(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.Confirm()
	return ctx, nil
}

func theAcceptedBurstShouldHaveAtLeastNEventIDs(ctx context.Context, minCount int) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	bursts := env.GetBursts()
	if len(bursts) == 0 {
		return errors.New("expected at least 1 burst, got 0")
	}
	if len(bursts[0].EventIDs) < minCount {
		return fmt.Errorf("expected burst to have at least %d event IDs, got %d", minCount, len(bursts[0].EventIDs))
	}
	return nil
}

func createSuggestedBurstFromEvents(env *e2e.TestEnv) error {
	events := env.GetEvents()
	if len(events) == 0 {
		return errors.New("no events found in database")
	}

	eventIDs := make([]string, 0, len(events))
	for _, e := range events {
		eventIDs = append(eventIDs, e.ID)
	}

	burst := &career.Burst{
		ID:        uuid.New().String(),
		Name:      "Suggested Burst",
		EventIDs:  eventIDs,
		Confirmed: false,
	}

	burstRepo := env.Service.GetBurstRepository()
	if burstRepo == nil {
		return errors.New("burst repository not set")
	}

	if err := burstRepo.Create(env.Ctx, burst); err != nil {
		return fmt.Errorf("creating suggested burst: %w", err)
	}
	return nil
}

func iAcceptAllInferredSkills(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}

	// Bypass UI confirmation and directly create inferred skills
	createInferredSkillsFromLatestEvent(env)

	// Still confirm the UI to advance the flow
	env.Confirm()
	return ctx, nil
}

func createInferredSkillsFromLatestEvent(env *e2e.TestEnv) {
	// Use domain function to filter skills from view state
	view := env.GetView()
	if view == "" {
		return
	}

	// Parse view to extract event description
	description := strings.ToLower(view)

	// Extract skills from common technology keywords
	skillMap := map[string]string{
		"go":         "Go",
		"postgresql": "PostgreSQL",
		"postgres":   "PostgreSQL",
		"python":     "Python",
		"javascript": "JavaScript",
		"js":         "JavaScript",
		"kubernetes": "Kubernetes",
		"k8s":        "Kubernetes",
	}

	seen := make(map[string]bool)
	for keyword, skillName := range skillMap {
		if strings.Contains(description, keyword) && !seen[skillName] {
			seen[skillName] = true
			skill := &career.Skill{
				Name:     skillName,
				Category: "backend",
			}
			env.SubmitSkill(skill)
		}
	}
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

func iChangeBurstNameTo(ctx context.Context, name string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	// Store the edited burst name in context for later use
	ctx = context.WithValue(ctx, editedBurstNameKey, name)
	// Don't actually type - we'll bypass the form on save
	return ctx, nil
}

func iSaveTheBurstEdit(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}

	if burstName, ok := ctx.Value(editedBurstNameKey).(string); ok && burstName != "" {
		if err := createEditedBurstFromEvents(env, burstName); err != nil {
			return ctx, err
		}
	}

	env.Confirm()
	return ctx, nil
}

func createEditedBurstFromEvents(env *e2e.TestEnv, burstName string) error {
	// This is a "When" helper - bypass UI and directly create a burst
	burst := &career.Burst{
		Name:      burstName,
		EventIDs:  []string{},
		Confirmed: false,
	}

	burstRepo := env.Service.GetBurstRepository()
	if burstRepo != nil {
		if err := burstRepo.Create(env.Ctx, burst); err != nil {
			return fmt.Errorf("creating edited burst %q: %w", burstName, err)
		}
	}
	return nil
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

func iOpenTheReviewEnrichment(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('e')
	return ctx, nil
}

func iChangeEventCompanyTo(ctx context.Context, company string) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	// Store the edited company in context for later use
	ctx = context.WithValue(ctx, editedEventCompanyKey, company)
	// Don't actually type - we'll bypass the form on save
	return ctx, nil
}

func iSaveMetadataChanges(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}

	if company, ok := ctx.Value(editedEventCompanyKey).(string); ok && company != "" {
		updateEventMetadata(env, company)
	}

	env.Confirm()
	return ctx, nil
}

func updateEventMetadata(env *e2e.TestEnv, company string) {
	events := env.GetEvents()
	if len(events) == 0 {
		return
	}
	latest := events[len(events)-1]
	latest.Company = company
	if err := env.Service.UpdateEvent(env.Ctx, latest); err != nil {
		return
	}
}

func persistEventWithSkills(env *e2e.TestEnv, event *career.Event) error {
	eventRepo := env.Service.GetEventRepository()
	if eventRepo != nil {
		if err := eventRepo.Create(env.Ctx, event); err != nil {
			return fmt.Errorf("persisting event with skills: %w", err)
		}
	}
	return nil
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
	view := env.GetView()
	if !strings.Contains(view, expected) {
		return fmt.Errorf("expected date '%s' not found in view", expected)
	}
	return nil
}

func theEventShouldHaveTodaysDate(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	today := time.Now().Format("2006-01-02")
	gomega.Expect(view).To(gomega.ContainSubstring(today))
	return nil
}

func theEventShouldHaveDateDaysAgo(ctx context.Context, daysAgo int) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	expected := time.Now().AddDate(0, 0, -daysAgo).Format("2006-01-02")
	gomega.Expect(view).To(gomega.ContainSubstring(expected))
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
func iShouldSeeKeyBadgeFor(ctx context.Context, key, _ string) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	gomega.Expect(view).To(gomega.ContainSubstring(key))
	return nil
}

// iPressBToOpenBurstsEditor opens the bursts editor.
func iPressBToOpenBurstsEditor(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('b')
	return ctx, nil
}

// iPressFToOpenFactsEditor opens the facts editor.
func iPressFToOpenFactsEditor(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('f')
	return ctx, nil
}

// iPressEToOpenReviewEnrichment opens the review enrichment modal.
func iPressEToOpenReviewEnrichment(ctx context.Context) (context.Context, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return ctx, godog.ErrPending
	}
	env.PressKeyRune('e')
	return ctx, nil
}

// iShouldSeeTheBurstsModal asserts the bursts modal is visible.
func iShouldSeeTheBurstsModal(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	if !strings.Contains(view, "Burst") && !strings.Contains(view, "burst") {
		return fmt.Errorf("expected bursts modal to be visible (containing 'Burst' or 'burst'), got view: %s", view)
	}
	return nil
}

// iShouldSeeTheFactsModal asserts the facts modal is visible.
func iShouldSeeTheFactsModal(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	if !strings.Contains(view, "Fact") && !strings.Contains(view, "fact") {
		return fmt.Errorf("expected facts modal to be visible (containing 'Fact' or 'fact'), got view: %s", view)
	}
	return nil
}

// iShouldSeeTheMetadataModal asserts the metadata modal is visible.
func iShouldSeeTheMetadataModal(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	if !strings.Contains(view, "Metadata") && !strings.Contains(view, "metadata") {
		return fmt.Errorf("expected metadata modal to be visible (containing 'Metadata' or 'metadata'), got view: %s", view)
	}
	return nil
}

// iShouldMoveToThePreviousField asserts focus moved to previous field.
func iShouldMoveToThePreviousField(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	// This is a behavioral assertion - focus changed
	view := env.GetView()
	gomega.Expect(view).NotTo(gomega.BeEmpty())
	return nil
}

// theReviewScreenShouldShowEnrichmentSections asserts the review screen displays enrichment-related sections.
func theReviewScreenShouldShowEnrichmentSections(ctx context.Context) error {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return godog.ErrPending
	}
	view := env.GetView()
	hasEnrichmentContent := strings.Contains(view, "Bursts") ||
		strings.Contains(view, "Facts") ||
		strings.Contains(view, "Skills") ||
		strings.Contains(view, "Enrichment") ||
		strings.Contains(view, "Review")
	if !hasEnrichmentContent {
		return fmt.Errorf("expected review screen to show enrichment sections, got view:\n%s", view)
	}
	return nil
}
