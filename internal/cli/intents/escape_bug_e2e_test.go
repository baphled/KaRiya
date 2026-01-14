package intents_test

import (
	"github.com/baphled/kariya/internal/testutil/e2e"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// E2E Test: Escape Key Navigation per TUI Standards
// Reference: docs/TUI_STANDARDS.md lines 147-166
//
// Expected behavior:
// - Escape from Form (intermediate state) → Back to ChooseStrategy
// - Escape from ChooseStrategy (root state) → Cancel intent, return to main menu

var _ = Describe("E2E - Escape Key Navigation (TUI Standards Validation)", func() {
	var env *e2e.TestEnv

	BeforeEach(func() {
		env = e2e.SetupWithMemory(GinkgoT())
	})

	AfterEach(func() {
		env.Cleanup()
	})

	Describe("CaptureEvent Intent - Per TUI Standards", func() {
		Context("Escape from intermediate state (Form)", func() {
			It("should navigate back to ChooseStrategy (NOT main menu)", func() {
				By("Starting at main menu")
				view := env.GetView()
				Expect(view).To(ContainSubstring("Capture Event"))

				By("Selecting CaptureEvent intent")
				env.SelectIntentByName("capture_event")

				By("Verifying we're at strategy selection (root state)")
				view = env.GetView()
				Expect(view).To(Or(
					ContainSubstring("Quick"),
					ContainSubstring("Manual"),
					ContainSubstring("Strategy"),
				))

				By("Selecting Quick strategy to enter form state (intermediate)")
				env.Confirm()

				By("Verifying we're now in form state")
				view = env.GetView()
				// Form should show some input fields
				Expect(view).To(Or(
					ContainSubstring("Event"),
					ContainSubstring("Title"),
					ContainSubstring("Company"),
				))

				By("Pressing ESCAPE from intermediate state")
				env.Cancel() // Escape key

				By("Verifying we went back to ChooseStrategy (NOT main menu)")
				view = env.GetView()

				// TUI Standards expectation: Should be back at strategy selection
				Expect(view).To(Or(
					ContainSubstring("Quick"),
					ContainSubstring("Manual"),
					ContainSubstring("Strategy"),
				), "Expected: Back at strategy selection (intermediate → previous state)\nActual: %s", view)

				// Should NOT be at main menu
				Expect(view).NotTo(ContainSubstring("Browse Timeline"),
					"Should NOT be at main menu (only ChooseStrategy state)")
			})
		})

		Context("Escape from root state (ChooseStrategy)", func() {
			It("should cancel intent and return to main menu", func() {
				By("Ensuring we're at main menu first")
				env.PressKeyRune('m') // Force main menu
				view := env.GetView()
				Expect(view).To(ContainSubstring("Capture Event"))

				By("Selecting CaptureEvent intent")
				env.SelectIntentByName("capture_event")

				By("Verifying we're at strategy selection (root state)")
				view = env.GetView()
				Expect(view).To(Or(
					ContainSubstring("Quick"),
					ContainSubstring("Manual"),
					ContainSubstring("Strategy"),
				))

				By("Pressing ESCAPE from root state")
				env.Cancel()

				By("Verifying we're back at main menu (intent cancelled)")
				view = env.GetView()
				Expect(view).To(ContainSubstring("Capture Event"), "Should be at main menu")
				Expect(view).To(ContainSubstring("Browse Timeline"), "Should show other menu items")
				Expect(view).NotTo(ContainSubstring("Strategy"), "Should NOT be in CaptureEvent intent")
			})
		})
	})

	// Note: 'm' key behavior is tested separately in existing test suites
	// This bug focuses specifically on escape key navigation per TUI Standards
})
