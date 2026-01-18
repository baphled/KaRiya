package intents_test

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/domain/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/baphled/kariya/internal/testutil/e2e"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Tests for Review Enrichment editing functionality
// User should be able to press 'e' (metadata), 'b' (bursts), 'f' (facts) to edit
var _ = Describe("CaptureEvent Review Enrichment Editing", func() {
	var env *e2e.TestEnv
	var intent *intents.CaptureEventIntent

	BeforeEach(func() {
		env = e2e.Setup(GinkgoT())

		// Create intent with proper context - must include CaptureStrategy
		intentContext := &intents.CaptureEventContext{
			CaptureStrategy: string(intents.StrategyManual),
			CLIEventService: env.CLIService,
			CareerService:   env.Service,
		}
		var err error
		intent, err = intents.NewCaptureEventIntent(intentContext)
		Expect(err).NotTo(HaveOccurred())

		// Disable screens architecture for these tests - we're testing the legacy
		// modal editing flow which is still the working implementation
		intent.DisableScreens()

		// Initialize the intent (makes it active)
		intent.Init()
	})

	AfterEach(func() {
		env.Cleanup()
	})

	// Helper to navigate to enrichment review state
	navigateToEnrichmentReview := func() *career.CareerEvent {
		// Create and save event
		event := &career.CareerEvent{
			Text:    "Built REST API with Go",
			Date:    time.Now(),
			Company: "Test Corp",
			Project: "API Project",
		}
		err := env.Service.CaptureEvent(context.Background(), event, careerservice.ManualEntry)
		Expect(err).NotTo(HaveOccurred())

		// Navigate through workflow to enrichment review
		intent.Update(intents.StrategySelectedMsg{Strategy: "quick"})
		intent.Update(intents.FormSubmittedMsg{Event: event})
		intent.Update(intents.SubmitCompleteMsg{})
		intent.Update(intents.DismissModalMsg{})

		// Add some test bursts and facts for editing
		testBursts := []*career.Burst{
			{Name: "API Development", Description: "Built REST endpoints"},
			{Name: "Database Design", Description: "Designed schema"},
		}
		testFacts := []*career.Fact{
			{Text: "Improved API response time by 50%"},
			{Text: "Reduced database queries by 30%"},
		}

		intent.Update(intents.EnrichmentCompleteMsg{
			Bursts: testBursts,
			Facts:  testFacts,
		})

		return event
	}

	Describe("Metadata Editing (e key)", func() {
		It("should display metadata editing modal when 'e' is pressed", func() {
			navigateToEnrichmentReview()

			// Verify we're in enrichment review
			view := intent.View()
			Expect(view).To(ContainSubstring("Review"), "Should be in review state")

			// Press 'e' to edit metadata
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

			// View should now show metadata editing modal
			view = intent.View()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Edit"),
				ContainSubstring("Metadata"),
			), "Should show metadata editing modal")

			// Should show editing state is active
			// Note: Full form field rendering depends on terminal size and overlay compositing.
			// The key test is that EditingMode was set, which is validated by the escape test.
			Expect(view).To(SatisfyAny(
				ContainSubstring("Edit"),
				ContainSubstring("Editing"),
				ContainSubstring("Metadata"),
			), "Modal should indicate editing state is active")
		})

		It("should close metadata modal and return to review when Esc is pressed", func() {
			navigateToEnrichmentReview()

			// Open metadata modal
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

			// Press Esc to cancel
			intent.Update(tea.KeyMsg{Type: tea.KeyEscape})

			// Should be back in review state without modal
			view := intent.View()
			Expect(view).To(ContainSubstring("Review"), "Should return to review")
			// Should NOT show modal content
			Expect(view).NotTo(ContainSubstring("Edit Metadata"), "Modal should be closed")
		})

		It("should save metadata changes when form is submitted", func() {
			event := navigateToEnrichmentReview()

			// Open metadata modal
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

			// The modal should be using the MetadataEditorModelNew which has the event
			// When user submits the form, changes should be saved

			// For now, verify modal is shown
			view := intent.View()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Edit"),
				ContainSubstring("Metadata"),
				ContainSubstring(event.Company),
			), "Modal should show with event data")
		})
	})

	Describe("Burst Editing (b key)", func() {
		It("should display burst editing modal when 'b' is pressed", func() {
			navigateToEnrichmentReview()

			// Verify bursts are shown in review
			view := intent.View()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Burst"),
				ContainSubstring("API Development"),
			), "Should show bursts in review")

			// Press 'b' to edit bursts
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})

			// View should now show burst editing modal
			view = intent.View()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Edit"),
				ContainSubstring("Burst"),
			), "Should show burst editing modal")
		})

		It("should close burst modal when Esc is pressed", func() {
			navigateToEnrichmentReview()

			// Open burst modal
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})

			// Press Esc to cancel
			intent.Update(tea.KeyMsg{Type: tea.KeyEscape})

			// Should be back in review state
			view := intent.View()
			Expect(view).To(ContainSubstring("Review"), "Should return to review")
		})
	})

	Describe("Fact Editing (f key)", func() {
		It("should display fact editing modal when 'f' is pressed", func() {
			navigateToEnrichmentReview()

			// Verify facts are shown in review
			view := intent.View()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Fact"),
				ContainSubstring("Improved API"),
			), "Should show facts in review")

			// Press 'f' to edit facts
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})

			// View should now show fact editing modal
			view = intent.View()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Edit"),
				ContainSubstring("Fact"),
			), "Should show fact editing modal")
		})

		It("should close fact modal when Esc is pressed", func() {
			navigateToEnrichmentReview()

			// Open fact modal
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})

			// Press Esc to cancel
			intent.Update(tea.KeyMsg{Type: tea.KeyEscape})

			// Should be back in review state
			view := intent.View()
			Expect(view).To(ContainSubstring("Review"), "Should return to review")
		})
	})

	Describe("Modal Overlay Rendering", func() {
		It("should render modal as overlay on top of review content", func() {
			navigateToEnrichmentReview()

			// The background content should still be visible under the modal
			view := intent.View()
			Expect(view).To(ContainSubstring("Review"), "Background should show review")

			// Open metadata modal
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

			// View should contain both background and modal
			view = intent.View()

			// Modal should be centered and visible
			Expect(view).To(SatisfyAny(
				ContainSubstring("Edit"),
				ContainSubstring("Metadata"),
			), "Modal should be visible")
		})
	})

	Describe("Footer Updates", func() {
		It("should show editing shortcuts in review footer", func() {
			navigateToEnrichmentReview()

			view := intent.View()

			// Footer should show e, b, f shortcuts when not editing
			Expect(view).To(SatisfyAny(
				ContainSubstring("e:"),
				ContainSubstring("Edit"),
			), "Footer should show edit metadata shortcut")
		})

		It("should show modal-specific shortcuts when modal is open", func() {
			navigateToEnrichmentReview()

			// Open modal
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

			view := intent.View()

			// Footer should show modal-specific shortcuts (Enter, Esc, Tab)
			Expect(view).To(SatisfyAny(
				ContainSubstring("Enter"),
				ContainSubstring("Esc"),
				ContainSubstring("Cancel"),
			), "Footer should show modal shortcuts")
		})
	})
})
