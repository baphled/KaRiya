package captureevent_test

import (
	"strings"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/e2e"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// CaptureEvent E2E Workflow Tests
//
// These tests validate the complete CaptureEvent workflow as implemented in
// the screens architecture:
//
//	Choose Strategy -> Form -> Submit (with loading modal) -> Enrichment Review -> Complete
//
// NOTE: The screens architecture differs from the original PRD workflow:
// - PRD: Form -> Pre-Save Review -> Submit -> Enrichment -> Enrichment Review
// - Screens: Form -> Submit (direct) -> Enrichment Review
//
// The screens implementation skips the pre-save review step and goes directly
// to submit with a loading modal overlay.
//
// These tests use SubmitEvent() to bypass huh form keystroke simulation issues.
// The huh library requires command chaining that doesn't work well in E2E tests.

// createTestEvent creates a test event with the given text and date.
func createTestEvent(text string) *career.Event {
	return &career.Event{
		Text:      text,
		Date:      time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// createTestEventWithDetails creates a test event with full metadata.
func createTestEventWithDetails(text, company, project string, tags, categories []string) *career.Event {
	return &career.Event{
		Text:       text,
		Date:       time.Now(),
		Company:    company,
		Project:    project,
		Tags:       tags,
		Categories: categories,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

// isAtMainMenu checks if the view is showing the main menu (not just breadcrumb containing "Main Menu")
func isAtMainMenu(view string) bool {
	// Main menu shows the menu items like "Capture Event", "Browse Timeline", etc.
	// AND doesn't show form elements
	return strings.Contains(view, "Capture Event") &&
		strings.Contains(view, "Browse Timeline") &&
		!strings.Contains(view, "Event Description") &&
		!strings.Contains(view, "Submit") // Not a form submit button
}

// isInFormState checks if the view is showing the capture form
func isInFormState(view string) bool {
	// Form state shows the event description field
	return strings.Contains(view, "Event Description") ||
		strings.Contains(view, "Describe what you accomplished")
}

// isInStrategyState checks if the view is showing the strategy selection
func isInStrategyState(view string) bool {
	return strings.Contains(view, "Choose Strategy") ||
		(strings.Contains(view, "Quick") && strings.Contains(view, "Manual") &&
			!strings.Contains(view, "Event Description"))
}

var _ = Describe("CaptureEvent E2E Workflow", func() {
	var env *e2e.TestEnv

	BeforeEach(func() {
		env = e2e.Setup(GinkgoT())
	})

	AfterEach(func() {
		env.Cleanup()
	})

	Describe("Standard Capture Workflow (Quick Mode)", func() {
		It("should navigate from main menu to form via strategy selection", func() {
			// Starting state: Main menu
			view := env.GetView()
			Expect(isAtMainMenu(view)).To(BeTrue(),
				"Should be at main menu showing menu options")

			// Step 1: Select "Capture Event" from main menu
			env.SelectIntentByName("capture_event")

			// Step 2: Should be in "Choose Strategy" state
			view = env.GetView()
			Expect(isInStrategyState(view)).To(BeTrue(),
				"Should be in strategy selection state")

			// Step 3: Select "Quick" strategy (default is Quick, so just press Enter)
			env.Confirm()

			// Step 4: Should be in "Form" state
			view = env.GetView()
			Expect(isInFormState(view)).To(BeTrue(),
				"Should be in form state")
		})

		It("should trigger submit when event is submitted from form", func() {
			// Navigate to form
			env.SelectIntentByName("capture_event")
			env.Confirm() // Choose Quick strategy

			// Submit event using helper
			testEventText := "Built REST API with Go, PostgreSQL, JWT auth, and Redis caching"
			testEvent := createTestEvent(testEventText)
			env.SubmitEvent(testEvent)

			// Should show submit/loading state (async save triggers loading modal)
			view := env.GetView()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Saving"),
				ContainSubstring("Loading"),
				ContainSubstring("Success"),
				ContainSubstring("Review"),
			), "Should be in submit or post-submit state")
		})

		It("should complete intent without re-saving in post-save review", func() {
			// Complete workflow up to post-save review.
			env.SelectIntentByName("capture_event")
			env.Confirm() // Choose Quick strategy

			// Submit event using helper.
			testEvent := createTestEvent("Test event for post-save review")
			env.SubmitEvent(testEvent)

			// After SubmitEvent, the E2E helper processes:
			// SubmitMsg -> HandleSubmit(StateForm) -> showSubmitModal -> performSubmit
			// -> SubmitCompleteMsg -> success modal -> DismissModalMsg -> StateReview
			// We should now be on the post-save review screen.
			view := env.GetView()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Review"),
				ContainSubstring("Enrichment"),
			), "Should be on post-save review screen")

			// Pressing Enter should complete the intent and return to main menu
			// WITHOUT showing a loading/saving modal.
			env.Confirm()

			view = env.GetView()
			Expect(view).NotTo(ContainSubstring("Saving"),
				"Should NOT show saving modal in post-save review")
			Expect(isAtMainMenu(view)).To(BeTrue(),
				"Should return to main menu after confirming post-save review")

			// Verify only 1 event exists - no duplicate save.
			env.AssertEventCount(1)
		})
	})

	Describe("Quick Submit Path (Ctrl+S from Form)", func() {
		It("should submit when using Ctrl+S from form", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm() // Choose Quick strategy

			// In form state - type some text first
			env.TypeText("Quick submit test event")

			// Press Ctrl+S to submit
			env.PressKey(tea.KeyCtrlS)

			// Should trigger submit (may show loading or go to post-submit)
			view := env.GetView()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Saving"),
				ContainSubstring("Success"),
				ContainSubstring("Review"),
				And(ContainSubstring("Capture Event"), ContainSubstring("Browse Timeline")),
			), "Should be in submit, post-submit, or main menu state")
		})
	})

	Describe("Cancel Behavior", func() {
		It("should not persist event when cancelled from form", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm() // Choose strategy

			// Verify we're in form state
			view := env.GetView()
			Expect(isInFormState(view)).To(BeTrue(),
				"Should be in form state")

			// Cancel from form
			env.Cancel()

			// Should NOT have persisted any event
			env.AssertEventCount(0)
		})

		It("should return to previous state when cancelled from strategy selection", func() {
			env.SelectIntentByName("capture_event")

			// Verify we're in strategy selection
			view := env.GetView()
			Expect(isInStrategyState(view)).To(BeTrue(),
				"Should be in strategy selection state")

			// Cancel from strategy selection
			env.Cancel()

			// Should be back at main menu
			view = env.GetView()
			Expect(isAtMainMenu(view)).To(BeTrue(),
				"Should be at main menu after cancelling strategy selection")

			// Should NOT have persisted any event
			env.AssertEventCount(0)
		})
	})

	Describe("Manual Strategy Workflow", func() {
		It("should support manual strategy with all fields", func() {
			env.SelectIntentByName("capture_event")

			// Navigate down to Manual strategy
			env.NavigateDown()
			env.Confirm() // Select Manual

			// Submit event with full details using helper
			testEvent := createTestEventWithDetails(
				"Manual capture test event with rich context",
				"Test Company Inc",
				"Test Project Alpha",
				nil,
				nil,
			)
			env.SubmitEvent(testEvent)

			// Wait for save to complete
			maxAttempts := 10
			for range maxAttempts {
				view := env.GetView()
				if isAtMainMenu(view) {
					break
				}
				env.Confirm()
			}

			// Verify event saved
			env.AssertEventCount(1)

			// Verify event details were preserved
			events := env.GetEvents()
			Expect(events).To(HaveLen(1))
			Expect(events[0].Company).To(Equal("Test Company Inc"))
			Expect(events[0].Project).To(Equal("Test Project Alpha"))
		})
	})

	Describe("Escape Key Navigation", func() {
		It("should navigate back through states correctly", func() {
			// Test: Form -> Esc -> Choose Strategy -> Esc -> Main Menu

			env.SelectIntentByName("capture_event")
			env.Confirm() // Choose strategy - go to form

			// Should be in Form state
			view := env.GetView()
			Expect(isInFormState(view)).To(BeTrue(),
				"Should be in form state")

			// Press Esc - should go back to Choose Strategy
			env.Cancel()

			view = env.GetView()
			Expect(isInStrategyState(view)).To(BeTrue(),
				"Should be back in strategy selection state")

			// Press Esc again - should go back to Main Menu
			env.Cancel()

			view = env.GetView()
			Expect(isAtMainMenu(view)).To(BeTrue(),
				"Should be at main menu")
		})
	})
})
