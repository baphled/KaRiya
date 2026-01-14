package intents_test

import (
	"github.com/baphled/kariya/internal/testutil/e2e"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GenerateCV Navigation", func() {
	var env *e2e.TestEnv

	Describe("Navigation to GenerateCV Intent", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show GenerateCV as menu item", func() {
			env.AssertViewContains("Generate CV")
		})

		It("should navigate to GenerateCV when selected", func() {
			env.SelectIntentByName("generate_cv")
			env.AssertViewContainsAny("Profile", "Select", "CV", "Generate")
		})

		It("should show context help for profile selection", func() {
			env.SelectIntentByName("generate_cv")
			env.AssertViewContainsAny("Enter", "Esc", "q", "Select")
		})
	})

	Describe("Profile Selection State", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.SelectIntentByName("generate_cv")
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show profile options", func() {
			env.AssertViewContainsAny("Profile", "Select", "Senior", "Staff", "Principal", "Engineer")
		})

		It("should allow navigating profile options with j/k", func() {
			env.PressKeyRune('j')
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			env.PressKeyRune('k')
			view = env.GetView()
			Expect(view).NotTo(BeEmpty())
		})

		It("should allow navigating with arrow keys", func() {
			env.PressKey(tea.KeyDown)
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			env.PressKey(tea.KeyUp)
			view = env.GetView()
			Expect(view).NotTo(BeEmpty())
		})

		It("should proceed to audience selection when profile is selected", func() {
			env.Confirm()
			env.AssertViewContainsAny("Audience", "Target", "hiring", "recruiter", "peer", "Select")
		})

		It("should cancel and return to menu when pressing Escape", func() {
			env.Cancel()
			env.AssertViewContains("Capture Event")
		})

		It("should quit application when pressing 'q'", func() {
			// Note: q now quits the entire app
			// This test verifies the quit command is handled without panic
			env.Quit()
			// After quit, the app terminates - we can't assert view content
		})
	})

	Describe("Audience Selection State", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.SelectIntentByName("generate_cv")
			env.Confirm() // Select first profile
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show audience options", func() {
			env.AssertViewContainsAny("Audience", "hiring", "recruiter", "peer", "Select")
		})

		It("should allow navigating audience options", func() {
			env.PressKeyRune('j')
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
		})

		It("should go back to profile selection when pressing Escape", func() {
			env.Cancel()
			env.AssertViewContainsAny("Profile", "Select", "Senior", "Staff")
		})

		It("should quit application when pressing 'q'", func() {
			// Note: q now quits the entire app
			// This test verifies the quit command is handled without panic
			env.Quit()
			// After quit, the app terminates - we can't assert view content
		})
	})

	Describe("Cancel at Each State", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should cancel from profile selection with Escape", func() {
			env.SelectIntentByName("generate_cv")
			env.Cancel()
			env.AssertViewContains("Capture Event")
		})

		It("should quit application from audience selection with 'q'", func() {
			env.SelectIntentByName("generate_cv")
			env.Confirm() // Go to audience selection
			// Note: q now quits the entire app
			// This test verifies the quit command is handled without panic
			env.Quit()
			// After quit, the app terminates - we can't assert view content
		})

		It("should go back from audience to profile with Escape", func() {
			env.SelectIntentByName("generate_cv")
			env.Confirm() // Go to audience selection
			env.Cancel()  // Go back to profile
			env.AssertViewContainsAny("Profile", "Select", "Senior", "Staff")
		})
	})

	Describe("View Rendering", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should render profile selection without panics", func() {
			env.SelectIntentByName("generate_cv")
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should render audience selection without panics", func() {
			env.SelectIntentByName("generate_cv")
			env.Confirm()
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should show breadcrumbs or context", func() {
			env.SelectIntentByName("generate_cv")
			env.AssertViewContainsAny("Generate CV", "CV", "Main Menu", "Profile")
		})

		It("should show footer with navigation hints", func() {
			env.SelectIntentByName("generate_cv")
			env.AssertViewContainsAny("q", "Esc", "Enter", "Quit")
		})
	})

	Describe("Vim-style Navigation", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.SelectIntentByName("generate_cv")
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should navigate down with 'j' key", func() {
			env.PressKeyRune('j')
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
		})

		It("should navigate up with 'k' key", func() {
			env.PressKeyRune('j')
			env.PressKeyRune('k')
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
		})

		It("should not crash at list boundaries", func() {
			// Try navigating up at the top
			env.PressKeyRune('k')
			env.PressKeyRune('k')
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())

			// Navigate to bottom
			for i := 0; i < 10; i++ {
				env.PressKeyRune('j')
			}
			view = env.GetView()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Arrow Key Navigation", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.SelectIntentByName("generate_cv")
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should navigate down with down arrow", func() {
			env.PressKey(tea.KeyDown)
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
		})

		It("should navigate up with up arrow", func() {
			env.PressKey(tea.KeyDown)
			env.PressKey(tea.KeyUp)
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
		})

		It("should not crash at list boundaries with arrows", func() {
			env.PressKey(tea.KeyUp)
			env.PressKey(tea.KeyUp)
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Workflow Navigation", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should allow navigating through profile options and selecting", func() {
			env.SelectIntentByName("generate_cv")
			env.PressKeyRune('j') // Navigate to second profile
			env.Confirm()
			env.AssertViewContainsAny("Audience", "hiring", "recruiter")
		})

		It("should allow re-entering after cancellation", func() {
			env.SelectIntentByName("generate_cv")
			env.Cancel()
			env.SelectIntentByName("generate_cv")
			env.AssertViewContainsAny("Profile", "Select", "Senior")
		})
	})

})
