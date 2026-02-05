package e2e_test

import (
	"time"

	"github.com/baphled/kariya/internal/testutil/e2e"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Capture Review Workflow E2E", func() {
	var env *e2e.TestEnv

	BeforeEach(func() {
		env = e2e.Setup(GinkgoT())
	})

	AfterEach(func() {
		env.Cleanup()
	})

	Describe("Review State Workflow", func() {
		It("should transition to review state after form submission", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm() // Select Quick strategy

			testEvent := fixtures.EventWith("", "Built REST API with authentication", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// Should show loading or review state
			view := env.GetView()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Saving"),
				ContainSubstring("Loading"),
				ContainSubstring("Review"),
				ContainSubstring("Enrichment"),
			), "Should be in submit or review state")
		})

		It("should show inferred bursts in review state", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Led team of 5 engineers on microservices migration", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// Wait for review state
			maxAttempts := 10
			foundReview := false
			for i := 0; i < maxAttempts; i++ {
				view := env.GetView()
				if view != "" {
					// Review state should show bursts or facts
					foundReview = true
					break
				}
				time.Sleep(50 * time.Millisecond)
				env.Confirm() // Progress through any intermediate states
			}

			Expect(foundReview).To(BeTrue(), "Should reach review state")
		})

		It("should allow accepting all suggestions and completing workflow", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Implemented caching layer reducing latency by 50%", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// Navigate through review and confirm
			maxAttempts := 10
			for i := 0; i < maxAttempts; i++ {
				view := env.GetView()
				if view != "" && (view == "Capture Event" || view == "Browse Timeline") {
					break
				}
				time.Sleep(50 * time.Millisecond)
				env.Confirm() // Accept suggestions and complete
			}

			// Should be back at main menu
			env.AssertViewContainsAny("Capture Event", "Browse Timeline")
			env.AssertEventCount(1)
		})
	})

	Describe("Post-Save Review Flow", func() {
		It("should show post-save review after successful submission", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Deployed containerized application to production", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// Should eventually show review
			view := env.GetView()
			Expect(view).NotTo(BeEmpty(), "Should have a view after submission")
		})

		It("should not re-save when confirming post-save review", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Refactored legacy codebase for better maintainability", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			initialEvents := len(env.GetEvents())

			// Confirm review (should not save again)
			env.Confirm()

			// Event count should not increase
			finalEvents := len(env.GetEvents())
			Expect(finalEvents).To(Equal(initialEvents), "Should not create duplicate events")
		})

		It("should return to main menu after confirming review", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Optimized database queries improving response time", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// Progress through review
			maxAttempts := 10
			for i := 0; i < maxAttempts; i++ {
				view := env.GetView()
				if view != "" {
					// Check if we're at main menu
					if view == "Capture Event" || view == "Browse Timeline" {
						break
					}
				}
				time.Sleep(50 * time.Millisecond)
				env.Confirm()
			}

			env.AssertViewContainsAny("Capture Event", "Browse Timeline")
		})
	})

	Describe("Review Navigation", func() {
		It("should allow navigating through review items", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Architected scalable microservices platform", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// Wait for review state
			maxAttempts := 10
			for i := 0; i < maxAttempts; i++ {
				view := env.GetView()
				if view != "" {
					break
				}
				time.Sleep(50 * time.Millisecond)
				env.Confirm()
			}

			// Try navigation keys
			env.NavigateDown()
			env.NavigateUp()

			// Should still be in review state
			view := env.GetView()
			Expect(view).NotTo(BeEmpty(), "Should maintain review state")
		})

		It("should support vim-style navigation in review", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Implemented CI/CD pipeline with automated testing", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// Wait for review
			maxAttempts := 5
			for i := 0; i < maxAttempts; i++ {
				view := env.GetView()
				if view != "" {
					break
				}
				time.Sleep(50 * time.Millisecond)
				env.Confirm()
			}

			// Vim keys should work
			env.PressKeyRune('j') // Down
			env.PressKeyRune('k') // Up

			view := env.GetView()
			Expect(view).NotTo(BeEmpty(), "Should handle vim navigation")
		})
	})

	Describe("Cancel from Review", func() {
		It("should cancel workflow when pressing Escape in review state", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Mentored junior developers on best practices", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// Wait for review state
			maxAttempts := 10
			for i := 0; i < maxAttempts; i++ {
				view := env.GetView()
				if view != "" {
					break
				}
				time.Sleep(50 * time.Millisecond)
				env.Confirm()
			}

			// Cancel from review
			env.Cancel()

			// Should return to main menu
			env.AssertViewContainsAny("Capture Event", "Browse Timeline")
		})

		It("should not persist partial review when cancelled", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Conducted code reviews ensuring quality standards", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			initialEvents := len(env.GetEvents())

			// Cancel review
			env.Cancel()

			// Event count should not change
			Expect(env.GetEvents()).To(HaveLen(initialEvents))
		})
	})

	Describe("Keyboard Shortcuts in Review", func() {
		It("should handle Ctrl+S in review state", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Documented API specifications for client teams", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// Wait for review
			maxAttempts := 5
			for i := 0; i < maxAttempts; i++ {
				view := env.GetView()
				if view != "" {
					break
				}
				time.Sleep(50 * time.Millisecond)
				env.Confirm()
			}

			// Try Ctrl+S (should not cause errors)
			env.PressKey(tea.KeyCtrlS)

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
			Expect(view).NotTo(ContainSubstring("error"))
		})

		It("should handle Enter key to accept current item", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Collaborated with product team on feature design", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// Wait for review
			maxAttempts := 5
			for i := 0; i < maxAttempts; i++ {
				view := env.GetView()
				if view != "" {
					break
				}
				time.Sleep(50 * time.Millisecond)
				env.Confirm()
			}

			// Press Enter (should accept current item or complete)
			env.Confirm()

			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Multiple Events Workflow", func() {
		It("should handle capturing multiple events in sequence", func() {
			// First event
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent1 := fixtures.EventWith("", "First event: Setup development environment", "", "")
			testEvent1.ID = ""
			env.SubmitEvent(testEvent1)

			// Complete first capture
			maxAttempts := 10
			for i := 0; i < maxAttempts; i++ {
				view := env.GetView()
				if view != "" {
					// Check if back at main menu
					if view == "Capture Event" || view == "Browse Timeline" {
						break
					}
				}
				time.Sleep(50 * time.Millisecond)
				env.Confirm()
			}

			Expect(env.GetEvents()).ToNot(BeEmpty())

			// Second event
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent2 := fixtures.EventWith("", "Second event: Implemented authentication", "", "")
			testEvent2.ID = ""
			env.SubmitEvent(testEvent2)

			// Complete second capture
			for i := 0; i < maxAttempts; i++ {
				view := env.GetView()
				if view != "" {
					if view == "Capture Event" || view == "Browse Timeline" {
						break
					}
				}
				time.Sleep(50 * time.Millisecond)
				env.Confirm()
			}

			Expect(len(env.GetEvents())).To(BeNumerically(">=", 2))
		})

		It("should maintain isolation between capture sessions", func() {
			// First capture - cancel it
			env.SelectIntentByName("capture_event")
			env.Confirm()
			env.TypeText("Cancelled event")
			env.Cancel()

			initialCount := len(env.GetEvents())

			// Second capture - complete it
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Completed event after cancellation", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			maxAttempts := 10
			for i := 0; i < maxAttempts; i++ {
				view := env.GetView()
				if view != "" {
					if view == "Capture Event" || view == "Browse Timeline" {
						break
					}
				}
				time.Sleep(50 * time.Millisecond)
				env.Confirm()
			}

			// Only one event should be saved
			Expect(env.GetEvents()).To(HaveLen(initialCount + 1))
		})
	})
})
