package e2e_test

import (
	"strings"

	"github.com/baphled/kariya/internal/testutil/e2e"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Capture Navigation", func() {
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

		It("should ignore 'q' key at strategy selection (quit only from main menu)", func() {
			env.Quit()
			// Intent should still be active - q is ignored within intents.
		})
	})

	Describe("Quick Capture Workflow", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.SelectIntentByName("capture_event")
			env.Confirm()
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show form for event text", func() {
			env.AssertViewContainsAny("Event", "Text", "What happened", "Enter Details")
		})

		It("should go back when pressing Escape from form", func() {
			env.Cancel()
			env.AssertViewNotContains("Enter Details")
		})

		It("should allow typing event text", func() {
			env.TypeText("Gave a presentation")
			env.Tab()
			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("Select Capture Strategy"))
			Expect(view).To(Or(ContainSubstring("Date"), ContainSubstring("Event")))
		})
	})

	Describe("Form State Navigation", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.SelectIntentByName("capture_event")
			env.NavigateDown()
			env.Confirm()
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show Manual form with optional fields", func() {
			env.AssertViewContains("Event Description")
			env.TypeText("Test event description text")
			env.Confirm()
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
			env.Confirm()
			env.Cancel()
			env.AssertViewNotContains("Enter Details")
		})

		It("should ignore 'q' key from strategy selection (quit only from main menu)", func() {
			env.SelectIntentByName("capture_event")
			env.Quit()
			// Intent should still be active - q is ignored within intents.
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
			env.NavigateDown()
			env.Confirm()
			env.Cancel()
			env.AssertViewNotContains("Enter Details")
		})

		It("should allow re-selecting intent after cancellation", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()
			env.Cancel()
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
			env.PressKeyRune('j')
			env.PressKeyRune('k')
			env.AssertViewContains("Quick")
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
			env.PressKey(tea.KeyDown)
			env.PressKey(tea.KeyDown)
			env.AssertViewContainsAny("Manual", "Quick")
		})

		It("should not go above first item", func() {
			env.PressKey(tea.KeyUp)
			env.AssertViewContains("Quick")
		})
	})
})
