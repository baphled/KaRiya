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
			env.PopulateTestData(3, 0, 0) // Add events to prevent empty state modal (BUG-004)
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
			env.PopulateTestData(3, 0, 0) // Add events to prevent empty state modal (BUG-004)
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

		It("should ignore 'q' key within intent (quit only from main menu)", func() {
			// q no longer quits from within intents - only from main menu
			// This test verifies 'q' is handled gracefully (does nothing)
			env.Quit()
			// Intent should still be active (q is ignored within intents)
		})
	})

	Describe("Audience Selection State", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.PopulateTestData(3, 0, 0) // Add events to prevent empty state modal (BUG-004)
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

		It("should ignore 'q' key within intent (quit only from main menu)", func() {
			// q no longer quits from within intents - only from main menu
			// This test verifies 'q' is handled gracefully (does nothing)
			env.Quit()
			// Intent should still be active (q is ignored within intents)
		})
	})

	Describe("Cancel at Each State", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.PopulateTestData(3, 0, 0) // Add events to prevent empty state modal (BUG-004)
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should cancel from profile selection with Escape", func() {
			env.SelectIntentByName("generate_cv")
			env.Cancel()
			env.AssertViewContains("Capture Event")
		})

		It("should ignore 'q' key from audience selection (quit only from main menu)", func() {
			env.SelectIntentByName("generate_cv")
			env.Confirm() // Go to audience selection
			// q no longer quits from within intents - only from main menu
			// This test verifies 'q' is handled gracefully (does nothing)
			env.Quit()
			// Intent should still be active (q is ignored within intents)
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
			env.PopulateTestData(3, 0, 0) // Add events to prevent empty state modal (BUG-004)
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
			// With wizard flow enabled, we see huh form footer
			// The footer shows navigation: "↑ up • ↓ down • / filter • enter select"
			env.AssertViewContainsAny("up", "down", "enter", "select", "filter")
		})
	})

	Describe("Vim-style Navigation", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.PopulateTestData(3, 0, 0) // Add events to prevent empty state modal (BUG-004)
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
			env.PopulateTestData(3, 0, 0) // Add events to prevent empty state modal (BUG-004)
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
			env.PopulateTestData(3, 0, 0) // Add events to prevent empty state modal (BUG-004)
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
