package intents_test

import (
	"strings"

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

			// Submit the huh form (tab to Submit button and press Enter)
			env.SubmitHuhForm()

			// Step 5: Should be in "Pre-Save Review" state
			view = env.GetView()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Review"),
				ContainSubstring(testEventText),
			), "Should be in Pre-Save Review state")

			// At this point, bursts/facts should NOT be shown yet (pre-save review)
			// The workflow doc says enrichment happens AFTER save
			Expect(view).NotTo(ContainSubstring("Burst"), "Pre-save review should not show bursts")
			Expect(view).NotTo(ContainSubstring("Fact"), "Pre-save review should not show facts")

			// Step 6: Confirm to submit (press Enter)
			env.Confirm()

			// Step 7: Should briefly show "Submit" state with loading
			// Then auto-dismiss and show enrichment loading
			// We may not see this if it's too fast, so we'll just wait a moment
			// by checking the next state

			// Step 8: After save completes, enrichment runs automatically
			// We should eventually see Enrichment Review state
			// Note: Success modal auto-dismisses after 2 seconds
			// We need to wait for that, then wait for enrichment
			// For now, let's just check we're not at main menu yet

			// Give time for success modal + enrichment (in real app this is automatic)
			// In tests, we may need to manually advance through states
			// Let's check what state we're in after a reasonable wait

			view = env.GetView()
			// We might see: success modal, enrichment loading, or enrichment review
			// We should NOT see main menu yet
			Expect(view).NotTo(ContainSubstring("Main Menu"),
				"Should not return to main menu immediately after save")

			// Step 9: If we're seeing success modal or enrichment loading, wait/advance
			// For now, let's just verify we eventually reach Enrichment Review
			// Since success modal auto-dismisses, we may already be past it

			// Keep checking until we see enrichment review or timeout
			maxAttempts := 5
			for attempt := 0; attempt < maxAttempts; attempt++ {
				view = env.GetView()
				if strings.Contains(view, "Review") || strings.Contains(view, "Enrichment") {
					break
				}
				// If we see success modal, try to dismiss it
				if strings.Contains(view, "Success") || strings.Contains(view, "saved") {
					env.Confirm() // Dismiss success modal
				}
			}

			// Step 10: CRITICAL CHECK - Should be in Enrichment Review state
			// Per PRD_MASTER.md Section 7: Submit → Enrichment → Enrichment Review
			view = env.GetView()
			Expect(view).NotTo(ContainSubstring("Main Menu"),
				"After save and enrichment, should show Enrichment Review state, NOT main menu")

			// Should show Enrichment Review state with enriched data (or at least the review UI)
			Expect(view).To(SatisfyAny(
				ContainSubstring("Review"),
				ContainSubstring("Enrichment"),
			), "Should be in Enrichment Review state")

			// Note: Bursts/facts may be empty if enrichment didn't find any,
			// but the UI should still show the review state

			// Step 11: Press Enter to complete the intent
			env.Confirm()

			// Step 12: NOW we should be back at main menu
			view = env.GetView()
			Expect(view).To(ContainSubstring("Main Menu"),
				"After completing enrichment review, should return to main menu")

			// Step 13: Verify event was actually persisted
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
			env.SubmitHuhForm() // Submit form

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
			_ = env.GetView()
			// This behavior is unclear from docs, so let's just verify event saved
			env.AssertEventCount(1)
		})
	})

	Describe("Cancel Behavior", func() {
		It("should not persist event when cancelled from form", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm() // Choose strategy

			// Fill form (don't submit - we'll cancel)
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
			env.SubmitHuhForm() // Go to review

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

			// Submit the form (tab to Submit button and press Enter)
			env.SubmitHuhForm()

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
			env.SubmitHuhForm() // Go to review

			// Press Esc - should go back to form
			env.Cancel()

			// Should be back in form state (hard to verify exact state)
			// Let's just verify we're not at main menu
			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("Main Menu"))
		})
	})
})
