package e2e_test

import (
	"time"

	"github.com/baphled/kariya/internal/cli/app"
	"github.com/baphled/kariya/internal/testutil/e2e"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// E2E Tests for Generate CV Empty State Handling (BUG-004)
//
// PURPOSE: Verify that when a user has no career events, selecting "Generate CV"
// shows a helpful warning modal instead of stub/test data.
//
// BUG-004: Generate CV was showing fake "Test event for navigation integration"
// stub data instead of informing users they need to add events first.
//
// Related:
// - bugs/BUG-004-generate-cv-shows-stub-data-when-no-events-exist.md
// - internal/cli/app/app.go (stub data removal)
// - internal/cli/components/info_modal.go (new modal component)

var _ = Describe("E2E Generate CV Empty State (BUG-004)", func() {
	var env *e2e.TestEnv

	Describe("when no career events exist", func() {
		BeforeEach(func() {
			// Setup with empty database (no events, no facts)
			env = e2e.GetSharedEnv(GinkgoT())
			// Explicitly verify we have no events
			Expect(env.GetEvents()).To(BeEmpty(), "Test setup should have no events")
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show warning modal when selecting Generate CV", func() {
			// Navigate to Generate CV menu item and select it
			env.SelectIntentByName("generate_cv")

			// Should show warning modal with helpful message
			view := env.GetView()
			Expect(view).To(ContainSubstring("No Career Events"),
				"Should show 'No Career Events' title in modal")
			Expect(view).To(ContainSubstring("Capture Event"),
				"Should guide user to use 'Capture Event' feature")
		})

		It("should NOT show stub test data", func() {
			env.SelectIntentByName("generate_cv")

			// Must not show the old stub data that was injected
			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("Test event for navigation"),
				"Should NOT show stub test event data")
			Expect(view).NotTo(ContainSubstring("ev-stub"),
				"Should NOT show stub event ID")
			Expect(view).NotTo(ContainSubstring("fact-stub"),
				"Should NOT show stub fact ID")
		})

		It("should dismiss modal and return to menu on Esc", func() {
			env.SelectIntentByName("generate_cv")

			// Verify modal is showing
			env.AssertViewContains("No Career Events")

			// Press Esc to dismiss
			env.Cancel()

			// Should be back at menu
			Expect(env.IsInMenuState()).To(BeTrue(),
				"Should return to menu after dismissing modal with Esc")
			env.AssertViewNotContains("No Career Events")
		})

		It("should dismiss modal and return to menu on Enter", func() {
			env.SelectIntentByName("generate_cv")

			// Verify modal is showing
			env.AssertViewContains("No Career Events")

			// Press Enter to dismiss
			env.Confirm()

			// Should be back at menu
			Expect(env.IsInMenuState()).To(BeTrue(),
				"Should return to menu after dismissing modal with Enter")
			env.AssertViewNotContains("No Career Events")
		})

		It("should remain in menu state (not enter intent state)", func() {
			env.SelectIntentByName("generate_cv")

			// The modal overlays the menu, but we remain in StateMenu (not StateIntent)
			// Note: IsInMenuState() checks view content which is replaced by modal,
			// so we check the actual state instead
			Expect(env.Model.GetState()).To(Equal(app.StateMenu),
				"Should remain in menu state with modal overlay")
		})
	})

	Describe("when career events exist", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnv(GinkgoT())

			event := fixtures.EventWith("event-real-1", "Led team of 5 engineers on microservices migration project", "Test Company", "Migration Project")
			event.Date = time.Now().AddDate(0, -1, 0)
			env.AddEvent(event)

			// Verify event was added
			Expect(env.GetEvents()).To(HaveLen(1), "Should have 1 event after setup")
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should allow Generate CV and proceed to wizard (no warning modal)", func() {
			env.SelectIntentByName("generate_cv")

			// Should NOT show warning modal
			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("No Career Events"),
				"Should NOT show empty state warning when events exist")

			// Should be in the CV generation flow (intent active, not menu)
			Expect(env.IsInMenuState()).To(BeFalse(),
				"Should have left menu state to enter CV generation")
		})

		It("should show CV wizard content instead of warning", func() {
			env.SelectIntentByName("generate_cv")

			// Should show wizard content (profile selection or wizard modal)
			view := env.GetView()
			Expect(view).To(Or(
				ContainSubstring("Profile"),
				ContainSubstring("CV"),
				ContainSubstring("Generate"),
			), "Should show CV generation wizard content")
		})
	})
})
