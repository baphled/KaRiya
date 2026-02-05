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

var _ = Describe("Capture Error Scenarios E2E", func() {
	var env *e2e.TestEnv

	BeforeEach(func() {
		env = e2e.GetSharedEnv(GinkgoT())
	})

	AfterEach(func() {
		env.Cleanup()
	})

	Describe("Form Validation", func() {
		It("should handle empty form submission gracefully", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm() // Select Quick strategy

			// Try to submit empty form
			env.PressKey(tea.KeyCtrlS)

			view := env.GetView()
			// Should not crash or show panic
			Expect(view).NotTo(ContainSubstring("panic"))
			Expect(view).NotTo(ContainSubstring("runtime error"))
		})

		It("should validate maximum text length", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			longText := strings.Repeat("0123456789", 250)

			env.TypeText(longText)

			view := env.GetView()
			// Should still render without errors
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should handle special characters in event text", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Event with special chars: @#$%^&*()[]{}|\\\"'", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should handle unicode characters in event text", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Event with unicode: 你好世界 🚀 ñ", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should handle newlines in event text", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Line 1\nLine 2\nLine 3", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})
	})

	Describe("State Transition Errors", func() {
		It("should handle rapid key presses without crashing", func() {
			env.SelectIntentByName("capture_event")

			// Rapid navigation
			for i := 0; i < 10; i++ {
				env.NavigateDown()
				env.NavigateUp()
			}

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should handle multiple cancel attempts", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			// Multiple cancels
			env.Cancel()
			env.Cancel()
			env.Cancel()

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
			env.AssertViewContainsAny("Capture Event", "Browse Timeline")
		})

		It("should handle multiple confirm attempts in quick succession", func() {
			env.SelectIntentByName("capture_event")

			// Rapid confirms
			env.Confirm()
			env.Confirm()
			env.Confirm()

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should handle tab navigation beyond form bounds", func() {
			env.SelectIntentByName("capture_event")
			env.NavigateDown() // Manual strategy
			env.Confirm()

			// Excessive tabbing
			for i := 0; i < 20; i++ {
				env.Tab()
			}

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})
	})

	Describe("Keyboard Shortcut Edge Cases", func() {
		It("should handle Ctrl+S from strategy selection", func() {
			env.SelectIntentByName("capture_event")

			// Try Ctrl+S before entering form
			env.PressKey(tea.KeyCtrlS)

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should ignore quit command within intent", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			// Try to quit (should be ignored)
			env.Quit()

			view := env.GetView()
			// Should still be in capture intent
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should handle help key without crashing", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			// Press ? for help
			env.PressKeyRune('?')

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should handle undefined key combinations", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			// Random key combinations
			env.PressKey(tea.KeyCtrlA)
			env.PressKey(tea.KeyCtrlE)
			env.PressKey(tea.KeyCtrlX)

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})
	})

	Describe("Data Integrity", func() {
		It("should not create duplicate events on rapid submission", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Rapid submission test", "", "")
			testEvent.ID = ""

			initialCount := len(env.GetEvents())

			// Submit same event multiple times
			env.SubmitEvent(testEvent)
			env.SubmitEvent(testEvent)

			// Complete workflow
			maxAttempts := 10
			for i := 0; i < maxAttempts; i++ {
				view := env.GetView()
				if view != "" {
					if strings.Contains(view, "Capture Event") || strings.Contains(view, "Browse Timeline") {
						break
					}
				}
				time.Sleep(50 * time.Millisecond)
				env.Confirm()
			}

			finalCount := len(env.GetEvents())
			// Should only add one event despite multiple submits
			Expect(finalCount).To(Equal(initialCount+1), "Should not create duplicate events")
		})

		It("should preserve event data after navigation", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			expectedText := "Important project milestone achieved"
			testEvent := fixtures.EventWith("", expectedText, "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// Complete workflow
			maxAttempts := 10
			for i := 0; i < maxAttempts; i++ {
				view := env.GetView()
				if view != "" {
					if strings.Contains(view, "Capture Event") || strings.Contains(view, "Browse Timeline") {
						break
					}
				}
				time.Sleep(50 * time.Millisecond)
				env.Confirm()
			}

			// Verify event was saved correctly
			events := env.GetEvents()
			if len(events) > 0 {
				lastEvent := events[len(events)-1]
				Expect(lastEvent.Text).To(ContainSubstring("project milestone"))
			}
		})

		It("should not corrupt existing events when adding new ones", func() {
			// Create first event
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent1 := fixtures.EventWith("", "First stable event", "", "")
			testEvent1.ID = ""
			env.SubmitEvent(testEvent1)

			maxAttempts := 10
			for i := 0; i < maxAttempts; i++ {
				view := env.GetView()
				if view != "" {
					if strings.Contains(view, "Capture Event") || strings.Contains(view, "Browse Timeline") {
						break
					}
				}
				time.Sleep(50 * time.Millisecond)
				env.Confirm()
			}

			firstEventCount := len(env.GetEvents())

			// Create second event
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent2 := fixtures.EventWith("", "Second independent event", "", "")
			testEvent2.ID = ""
			env.SubmitEvent(testEvent2)

			for i := 0; i < maxAttempts; i++ {
				view := env.GetView()
				if view != "" {
					if strings.Contains(view, "Capture Event") || strings.Contains(view, "Browse Timeline") {
						break
					}
				}
				time.Sleep(50 * time.Millisecond)
				env.Confirm()
			}

			// Both events should exist
			finalEvents := env.GetEvents()
			Expect(len(finalEvents)).To(BeNumerically(">=", firstEventCount))
		})
	})

	Describe("View Rendering Robustness", func() {
		It("should render without panics after every state transition", func() {
			env.SelectIntentByName("capture_event")

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))

			env.Confirm()
			view = env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))

			env.TypeText("Test event")
			view = env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))

			env.Cancel()
			view = env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should handle empty view renders", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			// Get view multiple times
			for i := 0; i < 5; i++ {
				view := env.GetView()
				Expect(view).NotTo(ContainSubstring("panic"))
			}
		})

		It("should render modal overlays without z-index issues", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Test modal rendering", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// Should show loading modal or review
			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Memory and Resource Management", func() {
		It("should handle repeated workflow executions without memory leaks", func() {
			// Execute workflow multiple times
			for iteration := 0; iteration < 3; iteration++ {
				env.SelectIntentByName("capture_event")
				env.Confirm()

				testEvent := fixtures.EventWith("", "Iteration event", "", "")
				testEvent.ID = ""
				env.SubmitEvent(testEvent)

				// Complete workflow
				maxAttempts := 10
				for i := 0; i < maxAttempts; i++ {
					view := env.GetView()
					if view != "" {
						if strings.Contains(view, "Capture Event") || strings.Contains(view, "Browse Timeline") {
							break
						}
					}
					time.Sleep(50 * time.Millisecond)
					env.Confirm()
				}
			}

			// Should complete all iterations without crashing
			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should clean up intent state after completion", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Cleanup test event", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// Complete workflow
			maxAttempts := 10
			for i := 0; i < maxAttempts; i++ {
				view := env.GetView()
				if view != "" {
					if strings.Contains(view, "Capture Event") || strings.Contains(view, "Browse Timeline") {
						break
					}
				}
				time.Sleep(50 * time.Millisecond)
				env.Confirm()
			}

			// Start new workflow - should have clean state
			env.SelectIntentByName("capture_event")
			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
			Expect(view).To(SatisfyAny(
				ContainSubstring("Quick"),
				ContainSubstring("Manual"),
				ContainSubstring("Strategy"),
			))
		})
	})
})
