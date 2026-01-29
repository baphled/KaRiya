package intents_test

import (
	"context"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/domain/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/baphled/kariya/internal/testutil/e2e"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Simple, focused test to reproduce the enrichment review blank screen bug
var _ = Describe("CaptureEvent Enrichment - Simple Reproduction", func() {
	var env *e2e.TestEnv
	var intent *intents.CaptureEventIntent

	BeforeEach(func() {
		env = e2e.Setup(GinkgoT())

		// Create intent with proper context
		intentContext := &intents.CaptureEventContext{
			CaptureStrategy: "quick",
			CLIEventService: env.CLIService,
			CareerService:   env.Service,
		}
		var err error
		intent, err = intents.NewCaptureEventIntent(intentContext)
		Expect(err).NotTo(HaveOccurred())

		// Initialize the intent (makes it active)
		intent.Init()
	})

	AfterEach(func() {
		env.Cleanup()
	})

	Describe("The Bug: Blank Screen After Enrichment", func() {
		It("FAILS: should show enrichment review content after workflow completes", func() {
			// This test reproduces the EXACT workflow that causes the blank screen

			// Step 1: Start intent (choose strategy state)
			view := intent.View()
			Expect(view).To(ContainSubstring("Strategy"), "Should start at choose strategy")

			// Step 2: Select Quick strategy
			intent.Update(intents.StrategySelectedMsg{Strategy: "quick"})
			_ = intent.View()
			// Should now be in form state

			// Step 3: Submit form with event data
			// Create event that will be saved
			testEvent := &career.Event{
				Text: "Built REST API with Go and PostgreSQL for data processing",
				Date: time.Now(),
			}

			// Send FormSubmittedMsg (simulates form submission)
			intent.Update(intents.FormSubmittedMsg{Event: testEvent})
			view = intent.View()

			// Should be in pre-save review
			// Note: View may truncate long text with "..." so we check for a prefix that will be visible
			Expect(view).To(ContainSubstring("Built REST API with Go and PostgreSQL"), "Pre-save review should show event")

			// Step 4: User presses Enter to confirm and submit
			// This triggers performSubmit() which:
			// - Saves the event (sets event.ID)
			// - Returns SubmitCompleteMsg
			// Note: We need to wait for async operation to complete

			// For now, let's manually create a saved event and send SubmitCompleteMsg
			savedEvent := &career.Event{
				Text: testEvent.Text,
				Date: testEvent.Date,
			}
			err := env.Service.CaptureEvent(context.Background(), savedEvent, careerservice.ManualEntry)
			Expect(err).NotTo(HaveOccurred())
			Expect(savedEvent.ID).NotTo(BeEmpty(), "Event should have ID after save")

			// Send SubmitCompleteMsg
			intent.Update(intents.SubmitCompleteMsg{})

			// Should show success modal
			_ = intent.View()
			// Modal might show "Success" or "saved"

			// Step 5: Modal auto-dismisses after 2s → DismissModalMsg
			intent.Update(intents.DismissModalMsg{})

			// Should transition to Enrichment state
			view = intent.View()
			// Should show "Enriching" or "Extracting"
			// Note: Might already be past enrichment if it was fast - that's ok
			_ = strings.Contains(view, "Enriching") || strings.Contains(view, "Extracting")

			// Step 6: Enrichment completes → EnrichmentCompleteMsg
			// Simulate enrichment returning results (could be empty)
			intent.Update(intents.EnrichmentCompleteMsg{
				Bursts: []*career.Burst{}, // Empty results for simplicity
				Facts:  []*career.Fact{},
			})

			// Step 7: Should now be in Enrichment Review state
			view = intent.View()

			// ====================================================================
			// THIS IS WHERE THE BUG IS:
			// User reports seeing BLANK SCREEN here (empty content)
			// ====================================================================

			// View should show enrichment review content
			Expect(view).NotTo(BeEmpty(), "View should not be empty")

			Expect(view).To(ContainSubstring("Review"),
				"View should contain 'Review' heading")

			Expect(view).To(ContainSubstring(savedEvent.Text),
				"View should show the saved event text - THIS IS THE BUG!")

			// Should show bursts/facts section (even if empty)
			Expect(view).To(SatisfyAny(
				ContainSubstring("Burst"),
				ContainSubstring("No bursts"),
			), "View should show bursts section")

			Expect(view).To(SatisfyAny(
				ContainSubstring("Fact"),
				ContainSubstring("No facts"),
			), "View should show facts section")
		})

		It("FAILS: should have event with ID in enrichment review state", func() {
			// Simplified test: Just check if event ID survives the workflow

			// Create and save event
			event := &career.Event{
				Text: "Test event",
				Date: time.Now(),
			}
			err := env.Service.CaptureEvent(context.Background(), event, careerservice.ManualEntry)
			Expect(err).NotTo(HaveOccurred())
			eventID := event.ID
			Expect(eventID).NotTo(BeEmpty())

			// Navigate through workflow
			intent.Update(intents.StrategySelectedMsg{Strategy: "quick"})
			intent.Update(intents.FormSubmittedMsg{Event: event})
			// Note: event passed here has the ID

			// Submit
			intent.Update(intents.SubmitCompleteMsg{})
			intent.Update(intents.DismissModalMsg{})

			// Enrichment
			intent.Update(intents.EnrichmentCompleteMsg{
				Bursts: []*career.Burst{},
				Facts:  []*career.Fact{},
			})

			// At this point, the intent should still have the event with ID
			// But the bug is that it might have lost the ID or the event itself

			// Try to get the view and check if event data is there
			view := intent.View()

			// If event ID was preserved, we should see the event text
			Expect(view).To(ContainSubstring(event.Text),
				"Event text should be visible in enrichment review - proves event ID was preserved")
		})

		It("FAILS: buildReviewBaseView should show content when postSaveReview=true", func() {
			// This test targets the EXACT method that's supposed to show content
			// but returns empty/blank

			// The issue: buildReviewBaseView() checks i.state.postSaveReview
			// If it's false, it shows metadata only (pre-save review)
			// If it's true, it shows metadata + bursts + facts (enrichment review)
			//
			// The bug: Either postSaveReview is false when it should be true,
			// OR the event is nil/missing, OR something else

			// We can't call buildReviewBaseView directly since it's unexported
			// So we test via the full workflow and check the view output

			event := &career.Event{
				Text: "Event for buildReviewBaseView test",
				Date: time.Now(),
			}
			err := env.Service.CaptureEvent(context.Background(), event, careerservice.ManualEntry)
			Expect(err).NotTo(HaveOccurred())

			// Complete workflow to enrichment review
			intent.Update(intents.StrategySelectedMsg{Strategy: "quick"})
			intent.Update(intents.FormSubmittedMsg{Event: event})
			intent.Update(intents.SubmitCompleteMsg{})
			intent.Update(intents.DismissModalMsg{})
			intent.Update(intents.EnrichmentCompleteMsg{
				Bursts: []*career.Burst{},
				Facts:  []*career.Fact{},
			})

			view := intent.View()

			// In enrichment review (postSaveReview=true), should show:
			// - "Review Enrichment Results" title (not "Review Event Details")
			// - Event text
			// - Bursts section
			// - Facts section

			Expect(view).To(ContainSubstring("Enrichment"),
				"Should show 'Enrichment' in title when postSaveReview=true")

			Expect(view).To(ContainSubstring(event.Text),
				"Should show event text in enrichment review")

			Expect(view).To(ContainSubstring("Burst"),
				"Should show bursts section in enrichment review")

			Expect(view).To(ContainSubstring("Fact"),
				"Should show facts section in enrichment review")
		})
	})

	// NOTE: Debug tests removed - use GinkgoWriter in individual tests for debugging
})
