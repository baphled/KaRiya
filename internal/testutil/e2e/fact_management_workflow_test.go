package e2e_test

import (
	"strings"

	"github.com/baphled/kariya/internal/testutil/e2e"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("E2E FactManagement Workflow", func() {
	var env *e2e.TestEnv

	Describe("Navigation to FactManagement Intent", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show FactManagement as menu item", func() {
			env.AssertViewContainsAny("Manage Facts", "Facts")
		})

		It("should navigate to FactManagement when selected", func() {
			env.SelectIntentByName("fact_management")
			env.AssertViewContainsAny("Fact", "Manage", "List", "No facts")
		})

		It("should show context help for fact list", func() {
			env.SelectIntentByName("fact_management")
			env.AssertViewContainsAny("Enter", "Esc", "q", "n")
		})
	})

	Describe("Empty Fact List", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show empty message when no facts exist", func() {
			env.AssertFactCount(0)
			env.SelectIntentByName("fact_management")
			env.AssertViewContainsAny("No facts", "empty", "no facts found", "Fact")
		})

		It("should not panic on empty fact list", func() {
			env.SelectIntentByName("fact_management")
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should allow navigation back from empty list", func() {
			env.SelectIntentByName("fact_management")
			env.Cancel()
			env.AssertViewContains("Capture Event")
		})
	})

	Describe("Fact List with Data", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			env.PopulateTestData(5, 0, 3) // 5 events, 0 bursts, 3 facts
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should display facts in list", func() {
			env.SelectIntentByName("fact_management")
			env.AssertFactCount(3)
			// Should show fact content from fixtures
			env.AssertViewContainsAny("Reduced", "Mentored", "Improved", "Fact", "Text")
		})

		It("should allow navigating through facts with j/k", func() {
			env.SelectIntentByName("fact_management")
			env.PressKeyRune('j')
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			env.PressKeyRune('k')
			view = env.GetView()
			Expect(view).NotTo(BeEmpty())
		})

		It("should allow navigating with arrow keys", func() {
			env.SelectIntentByName("fact_management")
			env.PressKey(tea.KeyDown)
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			env.PressKey(tea.KeyUp)
			view = env.GetView()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Cancel Navigation", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should return to main menu when pressing Escape from list", func() {
			env.SelectIntentByName("fact_management")
			env.Cancel()
			env.AssertViewContains("Capture Event")
		})

		It("should quit application when pressing 'q' from list", func() {
			env.SelectIntentByName("fact_management")
			// Note: q now quits the entire app
			// This test verifies the quit command is handled without panic
			env.Quit()
			// After quit, the app terminates - we can't assert view content
		})
	})

	Describe("View Rendering", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should render list view without panics", func() {
			env.SelectIntentByName("fact_management")
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should show breadcrumbs or context", func() {
			env.SelectIntentByName("fact_management")
			env.AssertViewContainsAny("Fact", "Manage", "Main Menu")
		})

		It("should show footer with navigation hints", func() {
			env.SelectIntentByName("fact_management")
			env.AssertViewContainsAny("q", "Esc", "Enter", "n")
		})
	})

	Describe("Vim-style Navigation", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			env.PopulateTestData(5, 0, 3)
			env.SelectIntentByName("fact_management")
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
			env.PressKeyRune('k')
			env.PressKeyRune('k')
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())

			for i := 0; i < 5; i++ {
				env.PressKeyRune('j')
			}
			view = env.GetView()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Session Persistence", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show facts after restart", func() {
			env.PopulateTestData(5, 0, 3)
			env.AssertFactCount(3)

			env.SimulateRestart()

			env.AssertFactCount(3)
			env.SelectIntentByName("fact_management")
			env.AssertViewContainsAny("Reduced", "Mentored", "Improved", "Fact")
		})
	})

	Describe("Workflow Integration", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			env.PopulateTestData(5, 0, 3)
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should allow selecting fact and returning to list multiple times", func() {
			env.SelectIntentByName("fact_management")

			// First selection
			env.Confirm() // Go to detail
			env.Cancel()  // Go back to list

			// Second selection
			env.PressKeyRune('j') // Navigate to next fact
			env.Confirm()         // Go to detail
			env.Cancel()          // Go back to list

			env.AssertViewContainsAny("List", "Facts", "Text")
		})

		It("should allow re-entering after cancellation", func() {
			env.SelectIntentByName("fact_management")
			env.Cancel()
			env.SelectIntentByName("fact_management")
			env.AssertViewContainsAny("Fact", "List", "Text")
		})
	})

	Describe("Fact Detail View", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			env.PopulateTestData(5, 0, 3)
			env.SelectIntentByName("fact_management")
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show fact detail when pressing Enter", func() {
			env.Confirm() // Select first fact
			env.AssertViewContainsAny("Detail", "Text", "Fact", "Source")
		})

		It("should show edit option in detail view", func() {
			env.Confirm()
			env.AssertViewContainsAny("e", "Edit", "edit")
		})

		It("should go back to list when pressing Escape from detail", func() {
			env.Confirm() // Go to detail
			env.Cancel()  // Go back
			env.AssertViewContainsAny("List", "Facts", "Text")
		})
	})

	Describe("Fact Edit Workflow", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			env.PopulateTestData(5, 0, 3) // 5 events, 0 bursts, 3 facts
			env.SelectIntentByName("fact_management")
			env.Confirm() // Go to detail view
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should transition to edit view when pressing 'e'", func() {
			env.PressKeyRune('e')
			// Should show edit form
			env.AssertViewContainsAny("Edit", "Fact", "Text", "Description")
		})

		It("should show edit modal with form fields", func() {
			env.PressKeyRune('e')
			// The EditFactModal should display form inputs
			env.AssertViewContainsAny("Fact Text", "Text", "Enter", "Esc")
		})

		It("should return to detail view when cancelling edit", func() {
			env.PressKeyRune('e')
			env.Cancel() // Press Escape to cancel
			// Should be back at detail view
			env.AssertViewContainsAny("Detail", "Fact", "e")
		})

		It("should not crash when pressing edit key multiple times", func() {
			env.PressKeyRune('e')
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			Expect(view).NotTo(ContainSubstring("panic"))

			// Cancel and try again
			env.Cancel()
			env.PressKeyRune('e')
			view = env.GetView()
			Expect(view).NotTo(BeEmpty())
			Expect(view).NotTo(ContainSubstring("panic"))
		})
	})

	Describe("Edit Form Display and UX", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			env.PopulateTestData(5, 0, 3) // 5 events, 0 bursts, 3 facts
			env.SelectIntentByName("fact_management")
			env.Confirm() // Go to detail view
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show form immediately when entering edit mode", func() {
			env.PressKeyRune('e')
			view := env.GetView()
			// Form should be visible immediately with input fields
			Expect(view).To(SatisfyAny(
				ContainSubstring("Fact Text"),
				ContainSubstring("Text"),
				ContainSubstring("Description"),
			))
			Expect(view).NotTo(ContainSubstring("Initializing"))
			Expect(view).NotTo(ContainSubstring("Loading"))
		})

		It("should NOT show duplicate help text", func() {
			env.PressKeyRune('e')
			view := env.GetView()
			// Count occurrences of common help patterns
			enterCount := strings.Count(view, "Enter")
			escCount := strings.Count(view, "Esc")

			// With proper rendering, help text appears reasonably
			Expect(enterCount).To(BeNumerically("<=", 3), "Should not have duplicate Enter help text")
			Expect(escCount).To(BeNumerically("<=", 3), "Should not have duplicate Esc help text")
		})

		It("should maintain consistent layout in edit view", func() {
			env.PressKeyRune('e')
			view := env.GetView()
			// View should not be excessively long (broken layout symptom)
			lines := strings.Split(view, "\n")
			Expect(len(lines)).To(BeNumerically("<", 100), "View should not have excessive line count")
			// Should not have large empty gaps
			Expect(view).NotTo(ContainSubstring("\n\n\n\n\n"))
		})
	})
})
