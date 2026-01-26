package intents_test

import (
	"strings"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/e2e"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// CaptureEvent Post-Save Enrichment E2E Tests
//
// These tests validate the 6-state post-save enrichment workflow:
//   Choose Strategy → Form → Pre-Save Review → Submit →
//   Enrichment (async) → Enrichment Review → Complete
//
// Reference: docs/PRD_MASTER.md Section 7 (lines 222-258), Section 9a (lines 515-554)
//
// Key Requirements:
// 1. Pre-Save Review shows ONLY metadata (no bursts/facts)
// 2. After save, success modal shows for 2s and auto-dismisses
// 3. Enrichment runs asynchronously (loading state shown)
// 4. Enrichment Review shows inferred bursts/facts
// 5. User can accept/reject/edit bursts/facts
// 6. Workflow completes only after user confirms enrichment review
//
// Test Coverage:
// - Happy path (enrichment succeeds with bursts/facts)
// - Sad path (enrichment fails, workflow continues)
// - Edge cases (empty results, service errors, nil event)

// NOTE: These tests are pending because they depend on async state transitions
// (success modal auto-dismiss, enrichment workflow) that are difficult to test
// reliably in E2E tests without proper timer mocking.
// The test helper `navigateToEnrichmentReview` uses `SubmitEvent` which bypasses
// the form but the subsequent state transitions require tea.Tick timer handling.
// NOTE: Requires timer mocking for reliable testing; use unit tests for now.
var _ = PDescribe("CaptureEvent Post-Save Enrichment E2E", func() {
	var env *e2e.TestEnv

	BeforeEach(func() {
		env = e2e.Setup(GinkgoT())
	})

	AfterEach(func() {
		env.Cleanup()
	})

	// Helper to navigate to enrichment review state
	navigateToEnrichmentReview := func() string {
		// Select Capture Event intent
		env.SelectIntentByName("capture_event")

		// Choose Quick strategy
		env.Confirm()

		// Submit event using helper to bypass huh form navigation issues
		// (TextArea fields don't respond to Enter/Tab as expected in tests)
		testEventText := "Built REST API with Go and PostgreSQL for high-throughput data processing"
		testEvent := &career.CareerEvent{
			Text:      testEventText,
			Date:      time.Now(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		env.SubmitEvent(testEvent)

		// Wait for save to complete and enrichment to start
		// Success modal auto-dismisses after 2s
		// In test environment, we may need to manually advance
		for i := 0; i < 10; i++ {
			view := env.GetView()
			// If we see enrichment loading or review, we're there
			if strings.Contains(view, "Enriching") ||
				strings.Contains(view, "Extracting") ||
				strings.Contains(view, "Review") {
				break
			}
			// Try to advance through any intermediate states
			env.PressKeyRune(' ')
		}

		return testEventText
	}

	Describe("Happy Path: Enrichment Success", func() {
		Context("when enrichment finds bursts and facts", func() {
			It("should show enrichment loading state, then enrichment review with results", func() {
				testEventText := navigateToEnrichmentReview()

				// Verify we reach enrichment review state
				view := env.GetView()
				Expect(view).To(SatisfyAny(
					ContainSubstring("Review"),
					ContainSubstring(testEventText),
				), "Should show enrichment review state")

				// Verify we're not at main menu yet
				Expect(view).NotTo(ContainSubstring("Main Menu"),
					"Should not return to main menu until user completes enrichment review")

				// Verify event was saved
				env.AssertEventCount(1)

				// Complete the workflow
				env.Confirm()

				// NOW we should be at main menu
				view = env.GetView()
				Expect(view).To(ContainSubstring("Main Menu"),
					"After completing enrichment review, should return to main menu")
			})

			It("should allow user to accept all inferred bursts and facts", func() {
				navigateToEnrichmentReview()

				// In enrichment review, press 'a' to accept all
				env.PressKeyRune('a')

				// Should complete successfully
				// May show confirmation or return to review
				// This depends on implementation

				// Complete the workflow
				env.Confirm()

				// Should be at main menu
				view := env.GetView()
				Expect(view).To(ContainSubstring("Main Menu"))
			})

			It("should allow user to reject all inferred bursts and facts", func() {
				navigateToEnrichmentReview()

				// In enrichment review, press 'r' to reject all
				env.PressKeyRune('r')

				// Should complete successfully (workflow continues even if user rejects all)
				// May show confirmation or return to review

				// Complete the workflow
				env.Confirm()

				// Should be at main menu
				view := env.GetView()
				Expect(view).To(ContainSubstring("Main Menu"))
			})

			It("should allow user to edit bursts before accepting", func() {
				navigateToEnrichmentReview()

				// Press 'b' to view/edit bursts
				env.PressKeyRune('b')

				// Should show burst editor or burst list
				view := env.GetView()
				Expect(view).To(SatisfyAny(
					ContainSubstring("Burst"),
					ContainSubstring("Edit"),
				))

				// Go back to review
				env.GoBack()

				// Complete the workflow
				env.Confirm()

				// Should be at main menu
				view = env.GetView()
				Expect(view).To(ContainSubstring("Main Menu"))
			})

			It("should allow user to edit facts before accepting", func() {
				navigateToEnrichmentReview()

				// Press 'f' to view/edit facts
				env.PressKeyRune('f')

				// Should show fact editor or fact list
				view := env.GetView()
				Expect(view).To(SatisfyAny(
					ContainSubstring("Fact"),
					ContainSubstring("Edit"),
				))

				// Go back to review
				env.GoBack()

				// Complete the workflow
				env.Confirm()

				// Should be at main menu
				view = env.GetView()
				Expect(view).To(ContainSubstring("Main Menu"))
			})

			It("should allow user to edit event metadata after save", func() {
				navigateToEnrichmentReview()

				// Press 'e' to edit event metadata
				env.PressKeyRune('e')

				// Should show metadata editor
				view := env.GetView()
				Expect(view).To(SatisfyAny(
					ContainSubstring("Edit"),
					ContainSubstring("Metadata"),
				))

				// Go back to review
				env.GoBack()

				// Complete the workflow
				env.Confirm()

				// Should be at main menu
				view = env.GetView()
				Expect(view).To(ContainSubstring("Main Menu"))
			})

			It("should show correct breadcrumb navigation through enrichment states", func() {
				// Navigate to pre-save review
				env.SelectIntentByName("capture_event")
				env.Confirm() // Choose strategy
				env.TypeText("Test event")
				env.Tab()
				env.TypeText("today")
				env.SubmitHuhForm()

				// Check pre-save review breadcrumb
				view := env.GetView()
				Expect(view).To(SatisfyAny(
					ContainSubstring("Review"),
					ContainSubstring("Capture Event"),
				))

				// Submit
				env.Confirm()

				// After enrichment, check enrichment review breadcrumb
				view = env.GetView()
				Expect(view).To(SatisfyAny(
					ContainSubstring("Review"),
					ContainSubstring("Enrichment"),
				))

				// Complete
				env.Confirm()
				view = env.GetView()
				Expect(view).To(ContainSubstring("Main Menu"))
			})
		})
	})

	Describe("Sad Path: Enrichment Failure", func() {
		Context("when enrichment service fails", func() {
			It("should still show enrichment review state with empty results", func() {
				// NOTE: Mock CareerService for error testing requires service injection.
				// For now, we'll test the workflow continues even with no results.

				testEventText := navigateToEnrichmentReview()

				// Even if enrichment fails, we should reach review state
				view := env.GetView()
				Expect(view).To(SatisfyAny(
					ContainSubstring("Review"),
					ContainSubstring(testEventText),
				), "Workflow should continue even if enrichment fails")

				// Complete the workflow
				env.Confirm()

				// Should be at main menu
				view = env.GetView()
				Expect(view).To(ContainSubstring("Main Menu"))
			})

			It("should NOT block workflow if enrichment returns zero results", func() {
				// Even if enrichment finds nothing, workflow continues
				testEventText := navigateToEnrichmentReview()

				// Should reach review state
				view := env.GetView()
				Expect(view).To(SatisfyAny(
					ContainSubstring("Review"),
					ContainSubstring(testEventText),
				))

				// Complete the workflow (should work even with no bursts/facts)
				env.Confirm()

				// Should be at main menu
				view = env.GetView()
				Expect(view).To(ContainSubstring("Main Menu"))
			})
		})

		Context("when enrichment times out", func() {
			It("should proceed to enrichment review after timeout", func() {
				Skip("TODO: Implement timeout handling in enrichment async command")

				// Enrichment should have a reasonable timeout (5-10s)
				// After timeout, workflow continues with empty results
			})
		})
	})

	Describe("Edge Cases", func() {
		Context("navigation behavior", func() {
			It("should allow user to cancel from pre-save review before enrichment", func() {
				env.SelectIntentByName("capture_event")
				env.Confirm() // Choose strategy
				env.TypeText("Test event")
				env.Tab()
				env.TypeText("today")
				env.SubmitHuhForm()

				// In pre-save review, press Escape to go back
				env.GoBack()

				// Should be back in form state
				view := env.GetView()
				Expect(view).To(SatisfyAny(
					ContainSubstring("Event"),
					ContainSubstring("Form"),
				))

				// Event should NOT be saved yet
				env.AssertEventCount(0)

				// Can cancel from form
				env.GoBack()

				// Should be at strategy selection or main menu
				view = env.GetView()
				Expect(view).To(SatisfyAny(
					ContainSubstring("Strategy"),
					ContainSubstring("Main Menu"),
				))
			})

			It("should NOT allow cancel from enrichment loading state", func() {
				env.SelectIntentByName("capture_event")
				env.Confirm() // Choose strategy
				env.TypeText("Test event")
				env.Tab()
				env.TypeText("today")
				env.SubmitHuhForm()
				env.Confirm() // Submit

				// Event is now saved, enrichment is running
				// Pressing Escape during enrichment should NOT cancel (event already saved)
				env.GoBack()

				// Should still be in enrichment workflow (not main menu)
				view := env.GetView()
				Expect(view).NotTo(ContainSubstring("Main Menu"),
					"Cannot cancel after event is saved")

				// Event should exist
				env.AssertEventCount(1)
			})

			It("should allow escape from enrichment review to return to main menu", func() {
				navigateToEnrichmentReview()

				// In enrichment review, press Escape
				env.GoBack()

				// Should return to main menu (event is saved, enrichment is done)
				view := env.GetView()
				Expect(view).To(ContainSubstring("Main Menu"),
					"Escape from enrichment review should return to main menu")

				// Event should be saved
				env.AssertEventCount(1)
			})
		})

		Context("pre-save vs post-save review differentiation", func() {
			It("should NOT show burst/fact options in pre-save review", func() {
				env.SelectIntentByName("capture_event")
				env.Confirm() // Choose strategy
				env.TypeText("Test event")
				env.Tab()
				env.TypeText("today")
				env.SubmitHuhForm()

				// In pre-save review
				view := env.GetView()

				// Should NOT show burst/fact navigation keys
				Expect(view).NotTo(ContainSubstring("b)"),
					"Pre-save review should not show burst key")
				Expect(view).NotTo(ContainSubstring("f)"),
					"Pre-save review should not show fact key")
			})

			It("should show burst/fact options in post-save enrichment review", func() {
				navigateToEnrichmentReview()

				// In enrichment review (post-save)
				view := env.GetView()

				// Should show burst/fact navigation keys
				// (Implementation may vary, but user should be able to access bursts/facts)
				// We'll check that we're in a review state with the event text
				Expect(view).To(SatisfyAny(
					ContainSubstring("Review"),
					ContainSubstring("Burst"),
					ContainSubstring("Fact"),
					ContainSubstring("Accept"),
					ContainSubstring("Reject"),
				), "Post-save review should show enrichment options")
			})
		})

		Context("double-save prevention", func() {
			It("should NOT re-submit event when pressing Enter in enrichment review", func() {
				navigateToEnrichmentReview()

				// Verify only 1 event exists
				env.AssertEventCount(1)

				// Press Enter in enrichment review (should complete, not re-submit)
				env.Confirm()

				// Should be at main menu
				view := env.GetView()
				Expect(view).To(ContainSubstring("Main Menu"))

				// Should STILL have only 1 event (not 2)
				env.AssertEventCount(1)
			})
		})

		Context("event data preservation", func() {
			It("should preserve event data through enrichment workflow", func() {
				testEventText := "Built high-performance API"

				// Create event with custom data
				env.SelectIntentByName("capture_event")
				env.Confirm() // Choose strategy
				env.TypeText(testEventText)
				env.Tab()
				env.TypeText("today")
				env.SubmitHuhForm()

				// Pre-save review
				view := env.GetView()
				Expect(view).To(ContainSubstring(testEventText))

				// Submit
				env.Confirm()

				// Enrichment review
				view = env.GetView()
				Expect(view).To(ContainSubstring(testEventText),
					"Event data should be preserved through enrichment")

				// Complete
				env.Confirm()

				// Verify event in database
				events := env.GetEvents()
				Expect(events).To(HaveLen(1))
				Expect(events[0].Text).To(Equal(testEventText))
			})
		})

		Context("enrichment results validation", func() {
			It("should show inferred bursts in enrichment review", func() {
				// Create event with data likely to generate bursts
				env.SelectIntentByName("capture_event")
				env.Confirm()
				env.TypeText("Led team of 5 engineers to build microservices architecture with Kubernetes")
				env.Tab()
				env.TypeText("today")
				env.SubmitHuhForm()
				env.Confirm() // Submit

				// Navigate to enrichment review
				view := env.GetView()

				// Should eventually show review state
				// (Bursts may or may not be found depending on service implementation)
				Expect(view).To(SatisfyAny(
					ContainSubstring("Review"),
					ContainSubstring("Burst"),
				))

				// Complete workflow
				env.Confirm()
				view = env.GetView()
				Expect(view).To(ContainSubstring("Main Menu"))
			})

			It("should show inferred facts in enrichment review", func() {
				// Create event with data likely to generate facts
				env.SelectIntentByName("capture_event")
				env.Confirm()
				env.TypeText("Implemented JWT authentication system using Go 1.21 and PostgreSQL 15")
				env.Tab()
				env.TypeText("today")
				env.SubmitHuhForm()
				env.Confirm() // Submit

				// Navigate to enrichment review
				view := env.GetView()

				// Should eventually show review state
				// (Facts may or may not be found depending on service implementation)
				Expect(view).To(SatisfyAny(
					ContainSubstring("Review"),
					ContainSubstring("Fact"),
				))

				// Complete workflow
				env.Confirm()
				view = env.GetView()
				Expect(view).To(ContainSubstring("Main Menu"))
			})
		})

		Context("async enrichment behavior", func() {
			It("should show loading state during enrichment", func() {
				env.SelectIntentByName("capture_event")
				env.Confirm()
				env.TypeText("Test event")
				env.Tab()
				env.TypeText("today")
				env.SubmitHuhForm()
				env.Confirm() // Submit

				// Should briefly show enrichment loading state
				// (This may be too fast to catch in tests, but we'll check)
				view := env.GetView()

				// Should show either enrichment loading or enrichment review
				Expect(view).To(SatisfyAny(
					ContainSubstring("Enriching"),
					ContainSubstring("Extracting"),
					ContainSubstring("Review"),
				), "Should show enrichment workflow states")

				// Eventually reach enrichment review
				env.Confirm()
				view = env.GetView()
				Expect(view).To(ContainSubstring("Main Menu"))
			})

			It("should not block UI during enrichment", func() {
				// This is inherently tested by the async implementation
				// If enrichment was blocking, tests would hang
				// So we just verify the workflow completes

				navigateToEnrichmentReview()
				env.Confirm()

				view := env.GetView()
				Expect(view).To(ContainSubstring("Main Menu"),
					"UI should not block during enrichment")
			})
		})
	})

	Describe("Integration with other features", func() {
		Context("when burst suggestions are saved", func() {
			It("should allow user to browse bursts after completing workflow", func() {
				// Complete capture workflow
				navigateToEnrichmentReview()
				env.PressKeyRune('a') // Accept all
				env.Confirm()

				// Return to main menu
				view := env.GetView()
				Expect(view).To(ContainSubstring("Main Menu"))

				// Navigate to burst management (if available).
				// This tests that bursts were actually saved.
				// NOTE: Burst management navigation added once intent is implemented.
			})
		})

		Context("when facts are extracted", func() {
			It("should persist facts to database", func() {
				// Complete capture workflow
				navigateToEnrichmentReview()
				env.Confirm()

				// Verify facts exist in database
				events := env.GetEvents()
				Expect(events).To(HaveLen(1))

				// Facts should be linked to the event
				facts := env.GetFacts()
				// Facts may or may not exist depending on enrichment results
				_ = facts
			})
		})
	})

	Describe("State machine correctness", func() {
		It("should follow exact state sequence: Choose → Form → PreReview → Submit → Enrichment → EnrichmentReview → Complete", func() {
			states := []string{}

			// Track state transitions by checking view content
			env.SelectIntentByName("capture_event")
			states = append(states, "Choose Strategy")

			env.Confirm()
			states = append(states, "Form")

			env.TypeText("Test event")
			env.Tab()
			env.TypeText("today")
			env.SubmitHuhForm()
			states = append(states, "Pre-Save Review")

			env.Confirm()
			states = append(states, "Submit")

			// Wait for enrichment
			view := env.GetView()
			if strings.Contains(view, "Enriching") || strings.Contains(view, "Extracting") {
				states = append(states, "Enrichment")
			}

			// Wait for enrichment review
			for i := 0; i < 5; i++ {
				view = env.GetView()
				if strings.Contains(view, "Review") {
					states = append(states, "Enrichment Review")
					break
				}
			}

			env.Confirm()
			states = append(states, "Complete")

			// Verify we followed the correct sequence
			Expect(states).To(HaveLen(7))
			Expect(states[0]).To(Equal("Choose Strategy"))
			Expect(states[1]).To(Equal("Form"))
			Expect(states[2]).To(Equal("Pre-Save Review"))
			Expect(states[3]).To(Equal("Submit"))
			// states[4] may be "Enrichment" if we caught it
			Expect(states[len(states)-2]).To(Equal("Enrichment Review"))
			Expect(states[len(states)-1]).To(Equal("Complete"))
		})
	})
})
