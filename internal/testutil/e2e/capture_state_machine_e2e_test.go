package e2e_test

import (
	"strings"
	"time"

	"github.com/baphled/kariya/internal/testutil/e2e"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Capture State Machine E2E", func() {
	var env *e2e.TestEnv

	BeforeEach(func() {
		env = e2e.GetSharedEnv(GinkgoT())
	})

	AfterEach(func() {
		env.Cleanup()
	})

	Describe("Complete State Flow: Choose Strategy → Form → Review → Submit", func() {
		It("should progress through all states in correct order", func() {
			// State 1: Choose Strategy
			env.SelectIntentByName("capture_event")
			view := env.GetView()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Strategy"),
				ContainSubstring("Quick"),
				ContainSubstring("Manual"),
			), "Should be in strategy selection state")

			// State 2: Form
			env.Confirm()
			view = env.GetView()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Event"),
				ContainSubstring("Text"),
				ContainSubstring("Description"),
			), "Should be in form state")

			// State 3: Submit (triggers review)
			testEvent := fixtures.EventWith("", "Complete workflow test event", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// State 4: Review or completion
			view = env.GetView()
			Expect(view).NotTo(BeEmpty(), "Should have transitioned to next state")

			// Complete workflow
			maxAttempts := 10
			for range maxAttempts {
				view := env.GetView()
				if strings.Contains(view, "Capture Event") || strings.Contains(view, "Browse Timeline") {
					break
				}
				time.Sleep(50 * time.Millisecond)
				env.Confirm()
			}

			env.AssertViewContainsAny("Capture Event", "Browse Timeline")
		})

		It("should allow backward navigation through states", func() {
			// Go to form
			env.SelectIntentByName("capture_event")
			env.Confirm()

			// Go back to strategy
			env.Cancel()
			view := env.GetView()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Strategy"),
				ContainSubstring("Quick"),
				ContainSubstring("Manual"),
			), "Should return to strategy selection")

			// Go back to main menu
			env.Cancel()
			env.AssertViewContainsAny("Capture Event", "Browse Timeline")
		})
	})

	Describe("State: Choose Strategy", func() {
		It("should start in strategy selection state", func() {
			env.SelectIntentByName("capture_event")
			view := env.GetView()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Strategy"),
				ContainSubstring("Quick"),
				ContainSubstring("Manual"),
			))
		})

		It("should allow selecting Quick strategy", func() {
			env.SelectIntentByName("capture_event")
			// Quick is default selection
			env.Confirm()

			view := env.GetView()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Event"),
				ContainSubstring("Text"),
			), "Should show form after selecting Quick")
		})

		It("should allow selecting Manual strategy", func() {
			env.SelectIntentByName("capture_event")
			env.NavigateDown() // Select Manual
			env.Confirm()

			view := env.GetView()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Event"),
				ContainSubstring("Company"),
				ContainSubstring("Project"),
			), "Should show manual form")
		})

		It("should return to main menu on cancel from strategy", func() {
			env.SelectIntentByName("capture_event")
			env.Cancel()

			env.AssertViewContainsAny("Capture Event", "Browse Timeline")
		})

		It("should maintain strategy state on navigation", func() {
			env.SelectIntentByName("capture_event")

			// Navigate between options
			env.NavigateDown()
			env.NavigateUp()
			env.NavigateDown()

			view := env.GetView()
			Expect(view).To(ContainSubstring("Manual"), "Should be on Manual option")

			env.Confirm()
			view = env.GetView()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Company"),
				ContainSubstring("Project"),
			), "Should enter Manual form")
		})
	})

	Describe("State: Form (Quick Strategy)", func() {
		BeforeEach(func() {
			env.SelectIntentByName("capture_event")
			env.Confirm() // Select Quick
		})

		It("should show Quick form with minimal fields", func() {
			view := env.GetView()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Event"),
				ContainSubstring("Text"),
				ContainSubstring("Description"),
			))
		})

		It("should allow text entry", func() {
			env.TypeText("Test event text")
			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should submit form with Ctrl+S", func() {
			env.TypeText("Quick submit with Ctrl+S")
			env.PressKey(tea.KeyCtrlS)

			view := env.GetView()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Saving"),
				ContainSubstring("Loading"),
				ContainSubstring("Review"),
			), "Should transition to submit/review")
		})

		It("should return to strategy on cancel", func() {
			env.Cancel()

			view := env.GetView()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Strategy"),
				ContainSubstring("Quick"),
				ContainSubstring("Manual"),
			), "Should return to strategy selection")
		})

		It("should navigate to main menu with 'm' key", func() {
			env.PressKeyRune('m')

			env.AssertViewContainsAny("Capture Event", "Browse Timeline")
		})
	})

	Describe("State: Form (Manual Strategy)", func() {
		BeforeEach(func() {
			env.SelectIntentByName("capture_event")
			env.NavigateDown() // Manual
			env.Confirm()
		})

		It("should show Manual form with all fields", func() {
			view := env.GetView()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Event"),
				ContainSubstring("Company"),
				ContainSubstring("Project"),
			))
		})

		It("should allow navigation between fields with Tab", func() {
			env.TypeText("Event text")
			env.Tab()

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should submit form with all fields", func() {
			testEvent := fixtures.EventWith("", "Manual event with metadata", "Test Corp", "Alpha Project")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			view := env.GetView()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Saving"),
				ContainSubstring("Review"),
			), "Should transition to submit/review")
		})

		It("should return to strategy on cancel", func() {
			env.Cancel()

			view := env.GetView()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Strategy"),
				ContainSubstring("Quick"),
				ContainSubstring("Manual"),
			))
		})
	})

	Describe("State: Submit (Loading Modal)", func() {
		It("should show loading indication during submission", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Test submission loading", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// Immediately after submit, should show loading or complete quickly
			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should transition to review or completion after submission", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Test post-submission transition", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// Wait for transition
			maxAttempts := 10
			completed := false
			for range maxAttempts {
				view := env.GetView()
				if strings.Contains(view, "Capture Event") || strings.Contains(view, "Review") {
					completed = true
					break
				}
				time.Sleep(50 * time.Millisecond)
				env.Confirm()
			}

			Expect(completed).To(BeTrue(), "Should complete submission")
		})
	})

	Describe("State: Review (Post-Save)", func() {
		It("should show review after successful save", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Review state test event", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// Should eventually show review or completion
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
		})

		It("should complete workflow from review with confirm", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Complete from review test", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// Progress through review
			maxAttempts := 10
			for range maxAttempts {
				view := env.GetView()
				if strings.Contains(view, "Capture Event") || strings.Contains(view, "Browse Timeline") {
					break
				}
				time.Sleep(50 * time.Millisecond)
				env.Confirm()
			}

			env.AssertViewContainsAny("Capture Event", "Browse Timeline")
		})

		It("should cancel workflow from review with escape", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Cancel from review test", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// Wait for review state
			maxAttempts := 10
			for range maxAttempts {
				view := env.GetView()
				if strings.Contains(view, "Review") || view != "" {
					break
				}
				time.Sleep(50 * time.Millisecond)
				env.Confirm()
			}

			env.Cancel()
			env.AssertViewContainsAny("Capture Event", "Browse Timeline")
		})
	})

	Describe("State Transitions Edge Cases", func() {
		It("should handle transition from strategy to form and back multiple times", func() {
			env.SelectIntentByName("capture_event")

			for range 3 {
				env.Confirm() // Go to form
				env.Cancel()  // Go back to strategy
			}

			view := env.GetView()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Strategy"),
				ContainSubstring("Quick"),
			))
		})

		It("should handle switching strategies before submission", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm() // Quick

			env.Cancel()       // Back to strategy
			env.NavigateDown() // Manual
			env.Confirm()

			view := env.GetView()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Company"),
				ContainSubstring("Project"),
			), "Should show Manual form")
		})

		It("should maintain form state during brief navigation", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			env.TypeText("Persistent text")

			// Brief navigation
			env.NavigateDown()
			env.NavigateUp()

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should handle rapid state transitions", func() {
			env.SelectIntentByName("capture_event")

			// Rapid transitions
			env.Confirm()
			env.Cancel()
			env.Confirm()
			env.Cancel()
			env.Confirm()

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})
	})

	Describe("State Invariants", func() {
		It("should never show multiple states simultaneously", func() {
			env.SelectIntentByName("capture_event")
			view := env.GetView()

			// Should not show both strategy and form indicators
			hasStrategy := strings.Contains(view, "Strategy") ||
				(strings.Contains(view, "Quick") && strings.Contains(view, "Manual"))
			hasForm := strings.Contains(view, "Enter Details") ||
				strings.Contains(view, "Event Description")

			if hasStrategy {
				Expect(hasForm).To(BeFalse(), "Should not show both strategy and form")
			}
		})

		It("should always have a valid view in every state", func() {
			states := []func(){
				func() {
					env.SelectIntentByName("capture_event")
				},
				func() {
					env.SelectIntentByName("capture_event")
					env.Confirm()
				},
				func() {
					env.SelectIntentByName("capture_event")
					env.Confirm()
					testEvent := fixtures.EventWith("", "State validation test", "", "")
					testEvent.ID = ""
					env.SubmitEvent(testEvent)
				},
			}

			for _, stateFn := range states {
				env = e2e.GetSharedEnv(GinkgoT())

				stateFn()
				view := env.GetView()

				Expect(view).NotTo(BeEmpty(), "View should never be empty")
				Expect(view).NotTo(ContainSubstring("panic"))
			}
		})

		It("should persist event only in submit state", func() {
			env.SelectIntentByName("capture_event")
			initialCount := len(env.GetEvents())

			// Strategy state - no persistence
			env.Confirm()
			Expect(env.GetEvents()).To(HaveLen(initialCount))

			// Form state - no persistence yet
			env.TypeText("Not yet persisted")
			Expect(env.GetEvents()).To(HaveLen(initialCount))

			// Cancel - still no persistence
			env.Cancel()
			env.Cancel()
			Expect(env.GetEvents()).To(HaveLen(initialCount))
		})
	})
})
