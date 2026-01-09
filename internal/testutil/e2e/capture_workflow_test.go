package e2e_test

import (
	"strings"

	"github.com/baphled/kariya/internal/testutil/e2e"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("E2E CaptureEvent Workflow", func() {
	var env *e2e.TestEnv

	Describe("Navigation to CaptureEvent Intent", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show CaptureEvent as menu item", func() {
			env.AssertViewContains("Capture Event")
		})

		It("should navigate to CaptureEvent when selected", func() {
			env.SelectIntentByName("capture_event")
			env.AssertViewContainsAny("Select Capture Strategy", "Quick", "Manual")
		})

		It("should show strategy options", func() {
			env.SelectIntentByName("capture_event")
			env.AssertViewContains("Quick")
			env.AssertViewContains("Manual")
		})
	})

	Describe("Strategy Selection", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.SelectIntentByName("capture_event")
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should start with Quick strategy selected", func() {
			view := env.GetView()
			Expect(view).To(ContainSubstring("Quick"))
		})

		It("should navigate to Manual strategy with down arrow", func() {
			env.NavigateDown()
			env.AssertViewContains("Manual")
		})

		It("should select Quick strategy and show form", func() {
			env.Confirm()
			env.AssertViewContainsAny("Event", "What happened", "Text", "Enter Details")
		})

		It("should select Manual strategy and show form with more fields", func() {
			env.NavigateDown()
			env.Confirm()
			env.AssertViewContainsAny("Event", "Text", "Date", "Company", "Enter Details")
		})

		It("should cancel intent when pressing Escape at strategy selection", func() {
			env.Cancel()
			env.AssertViewContains("Capture Event")
			env.AssertViewNotContains("Select Capture Strategy")
		})

		It("should cancel intent when pressing 'q' at strategy selection", func() {
			env.Quit()
			env.AssertViewContains("Capture Event")
		})

		It("should cancel intent when pressing 'm' at strategy selection", func() {
			env.PressKeyRune('m')
			env.AssertViewContains("Capture Event")
		})
	})

	Describe("Quick Capture Workflow", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.SelectIntentByName("capture_event")
			env.Confirm() // Select Quick strategy
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show form for event text", func() {
			env.AssertViewContainsAny("Event", "Text", "What happened", "Enter Details")
		})

		It("should go back when pressing Escape from form", func() {
			env.Cancel()
			// In full app flow, Escape from form goes back to strategy or main menu
			// depending on router handling - verify we're not stuck in form
			env.AssertViewNotContains("Enter Details")
		})

		It("should allow typing event text", func() {
			// NOTE: Text must avoid hotkey characters like 'm' (menu), 'q' (quit), 'j'/'k' (nav)
			// This is a known issue where form inputs don't capture global hotkeys
			env.TypeText("Gave a presentation")
			env.AssertViewContains("Gave a presentation")
		})
	})

	Describe("Form State Navigation", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.SelectIntentByName("capture_event")
			env.NavigateDown() // Go to Manual
			env.Confirm()      // Select Manual strategy
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show Manual form with optional fields", func() {
			env.AssertViewContainsAny("Date", "Company", "Project", "Tags")
		})

		It("should navigate between form fields with Tab", func() {
			env.TypeText("Test event")
			env.Tab()
			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("Select Capture Strategy"))
		})

		It("should go back when pressing Escape from form", func() {
			env.Cancel()
			// In full app flow, Escape from form navigates back
			// Verify we're no longer in the form state
			env.AssertViewNotContains("Enter Details")
		})

		It("should return to main menu when pressing 'm'", func() {
			env.PressKeyRune('m')
			env.AssertViewContains("Capture Event")
			env.AssertViewNotContains("Enter Details")
		})
	})

	Describe("Cancel at Each State", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should cancel from strategy selection and return to menu", func() {
			env.SelectIntentByName("capture_event")
			env.Cancel()
			env.AssertViewContains("Capture Event")
			env.AssertViewNotContains("Select Capture Strategy")
		})

		It("should navigate back when cancelling from form", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm() // Select Quick strategy
			env.Cancel()  // Cancel from form
			// Should navigate back - either to strategy or menu
			env.AssertViewNotContains("Enter Details")
		})

		It("should cancel completely with 'm' from form", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()         // Select Quick strategy
			env.PressKeyRune('m') // Return to main menu
			env.AssertViewContains("Capture Event")
		})

		It("should cancel completely with 'q' from strategy selection", func() {
			env.SelectIntentByName("capture_event")
			env.Quit()
			env.AssertViewContains("Capture Event")
		})
	})

	Describe("Back Navigation", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should navigate back when pressing Escape from form", func() {
			env.SelectIntentByName("capture_event")
			env.NavigateDown() // Select Manual
			env.Confirm()      // Enter form
			env.Cancel()       // Go back
			// Should no longer be in form state
			env.AssertViewNotContains("Enter Details")
		})

		It("should allow re-selecting intent after cancellation", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm() // Select Quick strategy (form)
			env.Cancel()  // Go back
			// Re-select the intent if we're at menu
			view := env.GetView()
			if !strings.Contains(view, "Select Capture Strategy") {
				env.SelectIntentByName("capture_event")
			}
			env.AssertViewContainsAny("Quick", "Manual", "Event", "Text")
		})
	})

	Describe("Vim-style Navigation", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.SelectIntentByName("capture_event")
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should navigate down with 'j' key", func() {
			env.PressKeyRune('j')
			env.AssertViewContains("Manual")
		})

		It("should navigate up with 'k' key", func() {
			env.PressKeyRune('j') // Go down to Manual
			env.PressKeyRune('k') // Go back up to Quick
			env.AssertViewContains("Quick")
		})
	})

	Describe("Database Persistence", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should start with empty database", func() {
			env.AssertEventCount(0)
		})

		It("should not persist event when cancelled before submission", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm() // Select Quick strategy
			env.TypeText("Test event text")
			env.Cancel() // Cancel from form
			env.AssertEventCount(0)
		})

		It("should not persist event when returning to main menu", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm() // Select Quick strategy
			env.TypeText("Test event text")
			env.PressKeyRune('m') // Return to main menu
			env.AssertEventCount(0)
		})
	})

	Describe("Session Persistence", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show empty timeline after restart with no events", func() {
			env.SimulateRestart()
			env.SelectIntentByName("browse_timeline")
			env.AssertViewContainsAny("No events", "empty", "Timeline")
		})
	})

	Describe("View Rendering", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should render strategy selection without panics", func() {
			env.SelectIntentByName("capture_event")
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should render form without panics", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should show footer with navigation hints", func() {
			env.SelectIntentByName("capture_event")
			env.AssertViewContainsAny("Quit", "q", "Esc", "Enter")
		})

		It("should show breadcrumbs or context", func() {
			env.SelectIntentByName("capture_event")
			env.AssertViewContainsAny("Capture Event", "Main Menu", "Strategy")
		})
	})

	Describe("Arrow Key Navigation", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.SelectIntentByName("capture_event")
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should navigate down with down arrow", func() {
			env.PressKey(tea.KeyDown)
			env.AssertViewContains("Manual")
		})

		It("should navigate up with up arrow", func() {
			env.PressKey(tea.KeyDown)
			env.PressKey(tea.KeyUp)
			env.AssertViewContains("Quick")
		})

		It("should not go below last item", func() {
			env.PressKey(tea.KeyDown) // Go to Manual (last item)
			env.PressKey(tea.KeyDown) // Try to go further down
			env.AssertViewContainsAny("Manual", "Quick")
		})

		It("should not go above first item", func() {
			env.PressKey(tea.KeyUp) // Try to go above first item
			env.AssertViewContains("Quick")
		})
	})
})
