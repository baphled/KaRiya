package intents_test

import (
	"github.com/baphled/kariya/internal/testutil/e2e"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// CaptureEvent E2E Workflow Tests
//
// These tests validate the complete CaptureEvent workflow as documented in:
// - docs/workflows/EVENT_CAPTURE_WORKFLOW.md
// - docs/PRD_MASTER.md (Section 7: UI Canvas, Section 9a: Flow Diagram)
//
// Expected workflow per PRD:
//   Choose Strategy → Form → Pre-Save Review → Submit → Save →
//   Enrichment → Enrichment Review → Complete
//
// Key points:
// - Enrichment is AUTOMATIC after save (not optional/user-triggered)
// - User MUST see Enrichment Review state after save (not main menu)
// - Bursts/facts shown in Enrichment Review (not Pre-Save Review)
// - User confirms enriched data before completing workflow
//
// Status: Currently FAILING due to form submission issues in E2E test environment.
// Once form submission is fixed, these tests validate the correct PRD workflow.

var _ = Describe("CaptureEvent E2E Workflow", func() {
	var env *e2e.TestEnv

	BeforeEach(func() {
		env = e2e.Setup(GinkgoT())
	})

	AfterEach(func() {
		env.Cleanup()
	})

	Describe("Standard Capture Workflow (Quick Mode)", func() {
		It("should complete the full workflow: Choose → Form → Pre-Save Review → Submit → Enrichment → Enrichment Review → Complete", func() {
			// Starting state: Main menu
			view := env.GetView()
			Expect(view).To(ContainSubstring("Capture Event"),
				"Should be at main menu showing Capture Event option")

			// Step 1: Select "Capture Event" from main menu
			env.SelectIntentByName("capture_event")

			// Step 2: Should be in "Choose Strategy" state
			view = env.GetView()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Choose Strategy"),
				ContainSubstring("Quick"),
				ContainSubstring("Manual"),
			))

			// Step 3: Select "Quick" strategy (default is Quick, so just press Enter)
			env.Confirm()

			// Step 4: Should be in "Form" state - fill out event details
			// Wait for form to render (may need a tick for huh form initialization)
			view = env.GetView()
			// Form should show event text field or form elements
			// Note: Huh forms may not show exact field names in view
			Expect(view).To(SatisfyAny(
				ContainSubstring("Event"),
				ContainSubstring("Capture"),
				ContainSubstring("Form"),
				ContainSubstring("text"),
			))

			// Type event text (this should go into the first field)
			testEventText := "Built REST API with Go, PostgreSQL, JWT auth, and Redis caching"
			env.TypeText(testEventText)

			// Tab to next field (date)
			env.Tab()

			// Type date
			env.TypeText("today")

			// Note: In Huh forms, we need to navigate to the submit button
			// Usually this is done by tabbing through fields until we reach the button
			// For now, let's try pressing Enter which should submit the form
			env.Confirm()

			// Step 5: Should be in "Review" state (pre-save)
			view = env.GetView()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Review"),
				ContainSubstring(testEventText),
				ContainSubstring("Confirm"),
			))

			// At this point, bursts/facts should NOT be shown yet (pre-save review)
			// The workflow doc says enrichment happens AFTER save

			// Step 6: Confirm to submit (press Enter or Ctrl+S)
			env.Confirm()

			// Step 7: Should be in "Submit" state - wait for save to complete
			// The submit modal should appear with "Saving..." or similar
			view = env.GetView()
			// Submit state may show modal or progress indicator
			// Let's be flexible and check for various submit indicators
			Expect(view).To(SatisfyAny(
				ContainSubstring("Saving"),
				ContainSubstring("Submitting"),
				ContainSubstring("Processing"),
				ContainSubstring("Success"),
				// In case submit is instant, we might immediately see success
			))

			// Step 8: After save completes, should show success modal
			// Press Enter to dismiss the modal
			env.Confirm()

			// Step 9: CRITICAL CHECK - Should go to Enrichment Review state (post-save)
			// Per PRD_MASTER.md Section 7: Submit → Enrichment → Enrichment Review
			// This is the CORRECT expected behavior per PRD, not a bug
			view = env.GetView()

			// Expected: Enrichment Review state with inferred bursts and facts
			// NOT expected: Main menu (that would be incorrect per PRD)
			Expect(view).NotTo(ContainSubstring("Main Menu"),
				"After save and enrichment, should show Enrichment Review state, NOT main menu")

			// Should show Enrichment Review state with enriched data
			Expect(view).To(SatisfyAny(
				ContainSubstring("Review"),
				ContainSubstring("Enrichment"),
				ContainSubstring("Burst"),
				ContainSubstring("Fact"),
				ContainSubstring("Inferred"),
			), "Enrichment Review should show inferred bursts and facts per PRD")

			// Step 10: Press Enter to complete the intent
			env.Confirm()

			// Step 11: NOW we should be back at main menu
			view = env.GetView()
			Expect(view).To(ContainSubstring("Main Menu"),
				"After completing post-save review, should return to main menu")

			// Step 12: Verify event was actually persisted
			env.AssertEventCount(1)
		})

		It("should NOT double-save when pressing Enter in post-save review", func() {
			// This tests for the double-save bug that postSaveReview flag prevents

			// Complete workflow up to post-save review
			env.SelectIntentByName("capture_event")
			env.Confirm() // Choose Quick strategy

			// Fill form
			env.TypeText("Test event for double-save check")
			env.Tab()
			env.TypeText("today")
			env.Confirm() // Submit form

			// Pre-save review
			env.Confirm() // Confirm to submit

			// Dismiss success modal
			env.Confirm()

			// Should be in post-save review - verify we're not at main menu
			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("Main Menu"))

			// Verify only 1 event exists so far
			env.AssertEventCount(1)

			// Press Enter in post-save review - should complete, NOT re-submit
			env.Confirm()

			// Should be at main menu now
			view = env.GetView()
			Expect(view).To(ContainSubstring("Main Menu"))

			// CRITICAL: Should still only have 1 event (not 2 from double-save)
			env.AssertEventCount(1)
		})
	})

	Describe("Quick Submit Path (Ctrl+S from Form)", func() {
		It("should skip pre-save review when using Ctrl+S from form", func() {
			// This tests the documented quick submit path:
			// Form → Ctrl+S → Submit (skip Review)

			env.SelectIntentByName("capture_event")
			env.Confirm() // Choose Quick strategy

			// Fill form
			env.TypeText("Quick submit test event")
			env.Tab()
			env.TypeText("today")

			// Press Ctrl+S to quick submit (skip review)
			env.PressKey(tea.KeyCtrlS)

			// Should go directly to Submit state (not Review)
			view := env.GetView()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Saving"),
				ContainSubstring("Success"),
			))

			// Dismiss success modal if shown
			env.Confirm()

			// Should either be in post-save review OR main menu
			// (Depends on if quick submit also skips post-save review)
			view = env.GetView()
			// This behavior is unclear from docs, so let's just verify event saved
			env.AssertEventCount(1)
		})
	})

	Describe("Cancel Behavior", func() {
		It("should not persist event when cancelled from form", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm() // Choose strategy

			// Fill form
			env.TypeText("Event to be cancelled")
			env.Tab()
			env.TypeText("today")

			// Cancel before submitting
			env.Cancel()

			// Should NOT have persisted the event
			env.AssertEventCount(0)
		})

		It("should not persist event when cancelled from pre-save review", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm() // Choose strategy

			// Fill form
			env.TypeText("Event to be cancelled in review")
			env.Tab()
			env.TypeText("today")
			env.Confirm() // Go to review

			// Cancel from review
			env.Cancel()

			// Should NOT have persisted the event
			env.AssertEventCount(0)
		})
	})

	Describe("Manual Strategy Workflow", func() {
		It("should support manual strategy with all fields", func() {
			env.SelectIntentByName("capture_event")

			// Navigate down to Manual strategy
			env.NavigateDown()
			env.Confirm() // Select Manual

			// Fill form with all fields
			env.TypeText("Manual capture test event with rich context")
			env.Tab() // Move to date
			env.TypeText("today")
			env.Tab() // Move to company
			env.TypeText("Test Company Inc")
			env.Tab() // Move to project
			env.TypeText("Test Project Alpha")

			// Continue through remaining fields and submit
			env.Confirm()

			// Review
			view := env.GetView()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Review"),
				ContainSubstring("Manual"),
			))

			// Submit
			env.Confirm()

			// Dismiss success modal
			env.Confirm()

			// Should be in post-save review (if enrichment enabled)
			// Or back at main menu (if enrichment skipped)
			// Let's just verify event saved
			env.AssertEventCount(1)
		})
	})

	Describe("Escape Key Navigation", func() {
		It("should navigate back through states correctly", func() {
			// Test: Form → Esc → Choose Strategy → Esc → Main Menu

			env.SelectIntentByName("capture_event")
			env.Confirm() // Choose strategy

			// Should be in Form state
			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("Main Menu"))

			// Press Esc - should go back to Choose Strategy
			env.Cancel()

			view = env.GetView()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Strategy"),
				ContainSubstring("Quick"),
				ContainSubstring("Manual"),
			))

			// Press Esc again - should go back to Main Menu
			env.Cancel()

			view = env.GetView()
			Expect(view).To(ContainSubstring("Main Menu"))
		})

		It("should navigate back from review to form", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm() // Choose strategy

			// Fill form
			env.TypeText("Test navigation back")
			env.Tab()
			env.TypeText("today")
			env.Confirm() // Go to review

			// Press Esc - should go back to form
			env.Cancel()

			// Should be back in form state (hard to verify exact state)
			// Let's just verify we're not at main menu
			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("Main Menu"))
		})
	})
})
