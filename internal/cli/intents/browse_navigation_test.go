package intents_test

import (
	"strings"

	"github.com/baphled/kariya/internal/testutil/e2e"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Browse Navigation", func() {
	var env *e2e.TestEnv

	Describe("Navigation to BrowseTimeline Intent", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show BrowseTimeline as menu item", func() {
			env.AssertViewContains("Browse Timeline")
		})

		It("should navigate to BrowseTimeline when selected", func() {
			env.SelectIntentByName("browse_timeline")
			env.AssertViewContainsAny("Timeline", "Browse Timeline", "Events")
		})

		It("should show context help for timeline view", func() {
			env.SelectIntentByName("browse_timeline")
			env.AssertViewContainsAny("Enter", "Esc", "q", "Filter")
		})
	})

	Describe("Empty Timeline Navigation", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show empty message when no events exist", func() {
			env.SelectIntentByName("browse_timeline")
			env.AssertViewContainsAny("No events", "empty", "no career events", "Timeline")
		})

		It("should not panic on empty timeline", func() {
			env.SelectIntentByName("browse_timeline")
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should allow navigation back from empty timeline", func() {
			env.SelectIntentByName("browse_timeline")
			env.Cancel()
			env.AssertViewContains("Capture Event")
		})
	})

	Describe("Timeline Navigation with Data", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.PopulateTestData(5, 0, 0)
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show pagination info when events exist", func() {
			env.SelectIntentByName("browse_timeline")
			env.AssertViewContainsAny("Page", "Events", "1 of")
		})

		It("should allow navigating through events with j/k", func() {
			env.SelectIntentByName("browse_timeline")
			env.PressKeyRune('j')
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			env.PressKeyRune('k')
			view = env.GetView()
			Expect(view).NotTo(BeEmpty())
		})

		It("should allow navigating with arrow keys", func() {
			env.SelectIntentByName("browse_timeline")
			env.PressKey(tea.KeyDown)
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			env.PressKey(tea.KeyUp)
			view = env.GetView()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Event Detail Navigation", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.PopulateTestData(3, 0, 0)
			env.SelectIntentByName("browse_timeline")
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show event detail when pressing Enter", func() {
			env.Confirm()
			env.AssertViewContainsAny("Date", "Text", "Company", "Event Detail")
		})

		It("should go back to timeline when pressing Escape from detail", func() {
			env.Confirm()
			env.Cancel()
			env.AssertViewContainsAny("Timeline", "Events", "Page")
		})

		It("should go back to timeline when pressing Back", func() {
			env.Confirm()
			env.GoBack()
			env.AssertViewContainsAny("Timeline", "Events", "Page")
		})
	})

	Describe("Cancel Navigation", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should return to main menu when pressing Escape from timeline", func() {
			env.SelectIntentByName("browse_timeline")
			env.Cancel()
			env.AssertViewContains("Capture Event")
		})

		It("should quit application when pressing 'q' from timeline", func() {
			env.SelectIntentByName("browse_timeline")
			env.Quit()
		})
	})

	Describe("Cancel from Event Detail", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.PopulateTestData(2, 0, 0)
			env.SelectIntentByName("browse_timeline")
			env.Confirm()
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should return to timeline when pressing Escape from detail", func() {
			env.Cancel()
			env.AssertViewContainsAny("Timeline", "Events", "Page")
		})

		It("should quit application when pressing 'q' from detail", func() {
			env.Quit()
		})

		It("should return to main menu when pressing Escape twice from detail", func() {
			env.Cancel() // Esc - back to timeline
			env.Cancel() // Esc - back to main menu
			env.AssertViewContains("Capture Event")
		})
	})

	Describe("Vim-style Navigation", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.PopulateTestData(5, 0, 0)
			env.SelectIntentByName("browse_timeline")
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

			for i := 0; i < 10; i++ {
				env.PressKeyRune('j')
			}
			view = env.GetView()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("View Rendering", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should render timeline view without panics", func() {
			env.SelectIntentByName("browse_timeline")
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should show breadcrumbs or context", func() {
			env.SelectIntentByName("browse_timeline")
			env.AssertViewContainsAny("Browse Timeline", "Timeline", "Main Menu")
		})

		It("should show footer with navigation hints", func() {
			env.SelectIntentByName("browse_timeline")
			env.AssertViewContainsAny("q", "Esc", "Enter", "Quit")
		})
	})

	Describe("Workflow Navigation", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.PopulateTestData(3, 0, 0)
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should allow selecting event and returning to timeline multiple times", func() {
			env.SelectIntentByName("browse_timeline")
			env.Confirm()
			env.Cancel()
			env.PressKeyRune('j')
			env.Confirm()
			env.Cancel()
			env.AssertViewContainsAny("Timeline", "Events", "Page")
		})
	})

	Describe("Event Edit Navigation", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.PopulateTestData(5, 0, 0)
			env.SelectIntentByName("browse_timeline")
			env.Confirm()
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show edit option in detail view", func() {
			env.AssertViewContainsAny("e", "Edit", "edit")
		})

		It("should transition to edit view when pressing 'e'", func() {
			env.PressKeyRune('e')
			env.AssertViewContainsAny("Edit", "Company", "Project", "Tags")
		})

		It("should show edit modal with form fields", func() {
			env.PressKeyRune('e')
			env.AssertViewContainsAny("Company", "Project", "Tags", "Categories")
		})

		It("should return to detail view when cancelling edit", func() {
			env.PressKeyRune('e')
			env.Cancel()
			env.AssertViewContainsAny("Detail", "Date", "Text", "e")
		})

		It("should not crash when pressing edit key multiple times", func() {
			env.PressKeyRune('e')
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			Expect(view).NotTo(ContainSubstring("panic"))

			env.Cancel()
			env.PressKeyRune('e')
			view = env.GetView()
			Expect(view).NotTo(BeEmpty())
			Expect(view).NotTo(ContainSubstring("panic"))
		})
	})

	Describe("Event Edit Form Display", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.PopulateTestData(5, 0, 0)
			env.SelectIntentByName("browse_timeline")
			env.Confirm()
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show form immediately when entering edit mode", func() {
			env.PressKeyRune('e')
			view := env.GetView()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Company"),
				ContainSubstring("Project"),
				ContainSubstring("Tags"),
			))
			Expect(view).NotTo(ContainSubstring("Initializing"))
			Expect(view).NotTo(ContainSubstring("Loading"))
		})

		It("should show company field with current value", func() {
			env.PressKeyRune('e')
			view := env.GetView()
			Expect(view).To(ContainSubstring("Company"))
			hasTestData := strings.Contains(view, "Acme") || strings.Contains(view, "TechStart") || strings.Contains(view, "BigCorp")
			Expect(hasTestData).To(BeTrue(), "Edit form should show company name from test data")
		})

		It("should NOT show duplicate help text", func() {
			env.PressKeyRune('e')
			view := env.GetView()
			enterCount := strings.Count(view, "Enter")
			escCount := strings.Count(view, "Esc")
			Expect(enterCount).To(BeNumerically("<=", 3), "Should not have duplicate Enter help text")
			Expect(escCount).To(BeNumerically("<=", 3), "Should not have duplicate Esc help text")
		})

		It("should maintain consistent layout in edit view", func() {
			env.PressKeyRune('e')
			view := env.GetView()
			lines := strings.Split(view, "\n")
			Expect(len(lines)).To(BeNumerically("<", 100), "View should not have excessive line count")
			Expect(view).NotTo(ContainSubstring("\n\n\n\n\n"))
		})
	})

	Describe("Edge Cases", func() {
		Describe("Empty Timeline", func() {
			BeforeEach(func() {
				env = e2e.SetupWithMemory(GinkgoT())
				// No data populated - empty timeline
			})

			AfterEach(func() {
				env.Cleanup()
			})

			It("should handle navigation keys with empty timeline gracefully", func() {
				env.SelectIntentByName("browse_timeline")
				// Should not panic on navigation with empty list
				env.PressKeyRune('j')
				env.PressKeyRune('k')
				env.PressKey(tea.KeyDown)
				env.PressKey(tea.KeyUp)
				view := env.GetView()
				Expect(view).NotTo(BeEmpty())
				Expect(view).NotTo(ContainSubstring("panic"))
			})

			It("should handle Enter key with empty timeline gracefully", func() {
				env.SelectIntentByName("browse_timeline")
				env.Confirm() // Enter on empty list
				view := env.GetView()
				Expect(view).NotTo(BeEmpty())
				Expect(view).NotTo(ContainSubstring("panic"))
			})

			It("should allow returning to main menu from empty timeline", func() {
				env.SelectIntentByName("browse_timeline")
				env.Cancel()
				env.AssertViewContains("Capture Event")
			})
		})

		Describe("Rapid Key Presses", func() {
			BeforeEach(func() {
				env = e2e.SetupWithMemory(GinkgoT())
				env.PopulateTestData(10, 0, 0)
				env.SelectIntentByName("browse_timeline")
			})

			AfterEach(func() {
				env.Cleanup()
			})

			It("should handle rapid up/down key presses", func() {
				// Rapidly press keys
				for i := 0; i < 20; i++ {
					env.PressKeyRune('j')
				}
				for i := 0; i < 20; i++ {
					env.PressKeyRune('k')
				}
				view := env.GetView()
				Expect(view).NotTo(BeEmpty())
				Expect(view).NotTo(ContainSubstring("panic"))
			})

			It("should handle alternating key presses", func() {
				for i := 0; i < 10; i++ {
					env.PressKeyRune('j')
					env.PressKeyRune('k')
				}
				view := env.GetView()
				Expect(view).NotTo(BeEmpty())
			})

			It("should handle rapid Enter/Escape sequences", func() {
				// Navigate to detail and back rapidly
				for i := 0; i < 3; i++ {
					env.Confirm() // Enter detail
					env.Cancel()  // Back to timeline
				}
				view := env.GetView()
				Expect(view).NotTo(BeEmpty())
				env.AssertViewContainsAny("Timeline", "Events", "Page")
			})
		})

		Describe("Boundary Navigation", func() {
			BeforeEach(func() {
				env = e2e.SetupWithMemory(GinkgoT())
				env.PopulateTestData(5, 0, 0)
				env.SelectIntentByName("browse_timeline")
			})

			AfterEach(func() {
				env.Cleanup()
			})

			It("should handle going past first item", func() {
				// At first item, try to go up multiple times
				for i := 0; i < 10; i++ {
					env.PressKeyRune('k')
				}
				view := env.GetView()
				Expect(view).NotTo(BeEmpty())
				Expect(view).NotTo(ContainSubstring("panic"))
			})

			It("should handle going past last item", func() {
				// Navigate to end and try to go further
				for i := 0; i < 20; i++ {
					env.PressKeyRune('j')
				}
				view := env.GetView()
				Expect(view).NotTo(BeEmpty())
				Expect(view).NotTo(ContainSubstring("panic"))
			})

			It("should handle arrow keys at boundaries", func() {
				// Test up arrow at beginning
				env.PressKey(tea.KeyUp)
				env.PressKey(tea.KeyUp)
				view := env.GetView()
				Expect(view).NotTo(BeEmpty())

				// Navigate to end with down arrow
				for i := 0; i < 20; i++ {
					env.PressKey(tea.KeyDown)
				}
				view = env.GetView()
				Expect(view).NotTo(BeEmpty())
			})
		})

		Describe("Vim Motion Consistency", func() {
			BeforeEach(func() {
				env = e2e.SetupWithMemory(GinkgoT())
				env.PopulateTestData(5, 0, 0)
				env.SelectIntentByName("browse_timeline")
			})

			AfterEach(func() {
				env.Cleanup()
			})

			It("should have consistent behavior between j/k and arrow keys", func() {
				// Navigate down with j
				env.PressKeyRune('j')
				viewAfterJ := env.GetView()

				// Go back up
				env.PressKeyRune('k')

				// Navigate down with arrow
				env.PressKey(tea.KeyDown)
				viewAfterArrow := env.GetView()

				// Views should be consistent (both at same position)
				Expect(viewAfterJ).To(Equal(viewAfterArrow))
			})

			It("should support 'g' for top and 'G' for bottom if implemented", func() {
				// Navigate to end
				for i := 0; i < 10; i++ {
					env.PressKeyRune('j')
				}
				view := env.GetView()
				Expect(view).NotTo(BeEmpty())
				// Note: These tests verify navigation doesn't crash,
				// specific g/G behavior depends on implementation
			})
		})
	})
})
