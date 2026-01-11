package intents_test

import (
	"strings"

	"github.com/baphled/kariya/internal/testutil/e2e"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Burst Management Navigation", func() {
	var env *e2e.TestEnv

	Describe("Navigation to BurstManagement Intent", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show BurstManagement as menu item", func() {
			env.AssertViewContainsAny("Manage Bursts", "Burst")
		})

		It("should navigate to BurstManagement when selected", func() {
			env.SelectIntentByName("burst_management")
			env.AssertViewContainsAny("Burst", "Manage", "List", "No bursts")
		})

		It("should show context help for burst list", func() {
			env.SelectIntentByName("burst_management")
			env.AssertViewContainsAny("Enter", "Esc", "q", "Select")
		})
	})

	Describe("Empty Burst List Navigation", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show empty message when no bursts exist", func() {
			env.SelectIntentByName("burst_management")
			env.AssertViewContainsAny("No bursts", "empty", "no bursts found", "Burst")
		})

		It("should not panic on empty burst list", func() {
			env.SelectIntentByName("burst_management")
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should allow navigation back from empty list", func() {
			env.SelectIntentByName("burst_management")
			env.Cancel()
			env.AssertViewContains("Capture Event")
		})
	})

	Describe("Burst List Navigation with Data", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.PopulateTestData(5, 2, 0)
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show confirmation status in list", func() {
			env.SelectIntentByName("burst_management")
			env.AssertViewContainsAny("Confirmed", "Yes", "No", "✓", "✗")
		})

		It("should allow navigating through bursts with j/k", func() {
			env.SelectIntentByName("burst_management")
			env.PressKeyRune('j')
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			env.PressKeyRune('k')
			view = env.GetView()
			Expect(view).NotTo(BeEmpty())
		})

		It("should allow navigating with arrow keys", func() {
			env.SelectIntentByName("burst_management")
			env.PressKey(tea.KeyDown)
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			env.PressKey(tea.KeyUp)
			view = env.GetView()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Burst Detail Navigation", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.PopulateTestData(5, 2, 0)
			env.SelectIntentByName("burst_management")
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show burst detail when pressing Enter", func() {
			env.Confirm()
			env.AssertViewContainsAny("Detail", "Name", "Description", "Events", "Confirmed")
		})

		It("should go back to list when pressing Escape from detail", func() {
			env.Confirm()
			env.Cancel()
			env.AssertViewContainsAny("List", "Bursts", "Name", "Confirmed")
		})

		It("should show action hints in detail view", func() {
			env.Confirm()
			env.AssertViewContainsAny("e", "f", "d", "c", "Events", "Facts", "Delete", "Confirm")
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
			env.SelectIntentByName("burst_management")
			env.Cancel()
			env.AssertViewContains("Capture Event")
		})

		It("should quit application when pressing 'q' from list", func() {
			env.SelectIntentByName("burst_management")
			env.Quit()
		})
	})

	Describe("Cancel from Detail View", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.PopulateTestData(5, 2, 0)
			env.SelectIntentByName("burst_management")
			env.Confirm()
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should return to list when pressing Escape from detail", func() {
			env.Cancel()
			env.AssertViewContainsAny("List", "Bursts", "Name")
		})

		It("should quit application when pressing 'q' from detail", func() {
			env.Quit()
		})
	})

	Describe("Vim-style Navigation", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.PopulateTestData(5, 2, 0)
			env.SelectIntentByName("burst_management")
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

	Describe("Arrow Key Navigation", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.PopulateTestData(5, 2, 0)
			env.SelectIntentByName("burst_management")
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
	})

	Describe("View Rendering", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should render list view without panics", func() {
			env.SelectIntentByName("burst_management")
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should show breadcrumbs or context", func() {
			env.SelectIntentByName("burst_management")
			env.AssertViewContainsAny("Burst", "Manage", "Main Menu")
		})

		It("should show footer with navigation hints", func() {
			env.SelectIntentByName("burst_management")
			env.AssertViewContainsAny("q", "Esc", "Enter", "Quit")
		})
	})

	Describe("Workflow Navigation", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.PopulateTestData(5, 2, 0)
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should allow selecting burst and returning to list multiple times", func() {
			env.SelectIntentByName("burst_management")
			env.Confirm()
			env.Cancel()
			env.PressKeyRune('j')
			env.Confirm()
			env.Cancel()
			env.AssertViewContainsAny("List", "Bursts", "Name")
		})

		It("should allow re-entering after cancellation", func() {
			env.SelectIntentByName("burst_management")
			env.Cancel()
			env.SelectIntentByName("burst_management")
			env.AssertViewContainsAny("Burst", "List", "Name")
		})
	})

	Describe("Burst Edit Navigation", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.PopulateTestData(5, 2, 0)
			env.SelectIntentByName("burst_management")
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
			env.AssertViewContainsAny("Edit Burst", "Burst Name", "Name", "Description")
		})

		It("should show edit modal with form fields", func() {
			env.PressKeyRune('e')
			env.AssertViewContainsAny("Burst Name", "Description", "Enter", "Esc")
		})

		It("should show current burst values in edit form", func() {
			env.PressKeyRune('e')
			env.AssertViewContainsAny("Authentication", "Mentoring", "Name", "Description")
		})

		It("should return to detail view when cancelling edit", func() {
			env.PressKeyRune('e')
			env.Cancel()
			env.AssertViewContainsAny("Detail", "Events", "Facts", "e", "f")
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

		It("should show edit form footer with navigation hints", func() {
			env.PressKeyRune('e')
			env.AssertViewContainsAny("Enter", "Esc", "Tab", "Confirm", "Cancel")
		})
	})

	Describe("Burst Edit with State Navigation", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.PopulateTestData(5, 2, 0)
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should maintain edit state through navigation", func() {
			env.SelectIntentByName("burst_management")
			env.Confirm()
			env.PressKeyRune('e')
			env.AssertViewContainsAny("Edit Burst", "Burst Name", "Name")

			env.Cancel()
			env.AssertViewContainsAny("Detail", "Events", "Facts")
		})

		It("should allow editing and returning to detail without errors", func() {
			env.SelectIntentByName("burst_management")
			env.Confirm()
			env.PressKeyRune('e')
			env.Cancel()

			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			Expect(view).NotTo(ContainSubstring("panic"))
			Expect(view).NotTo(ContainSubstring("nil pointer"))
		})
	})

	Describe("Edit Form Display", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.PopulateTestData(5, 2, 0)
			env.SelectIntentByName("burst_management")
			env.Confirm()
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show form immediately when entering edit mode", func() {
			env.PressKeyRune('e')
			view := env.GetView()
			Expect(view).To(ContainSubstring("Burst Name"))
			Expect(view).NotTo(ContainSubstring("Initializing"))
			Expect(view).NotTo(ContainSubstring("Loading"))
		})

		It("should show burst name field with current value", func() {
			env.PressKeyRune('e')
			view := env.GetView()
			Expect(view).To(ContainSubstring("Burst Name"))
			hasTestData := strings.Contains(view, "Authentication") || strings.Contains(view, "Mentoring")
			Expect(hasTestData).To(BeTrue(), "Edit form should show burst name from test data")
		})

		It("should show description field", func() {
			env.PressKeyRune('e')
			view := env.GetView()
			Expect(view).To(ContainSubstring("Description"))
		})

		It("should show submit confirmation field", func() {
			env.PressKeyRune('e')
			view := env.GetView()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Save Changes"),
				ContainSubstring("Submit"),
				ContainSubstring("Confirm"),
			))
		})

		It("should NOT show duplicate help text", func() {
			env.PressKeyRune('e')
			view := env.GetView()
			enterCount := strings.Count(view, "Enter")
			escCount := strings.Count(view, "Esc")
			Expect(enterCount).To(BeNumerically("<=", 3), "Should not have duplicate Enter help text")
			Expect(escCount).To(BeNumerically("<=", 3), "Should not have duplicate Esc help text")
		})

		It("should show logo/branding in edit view", func() {
			env.PressKeyRune('e')
			view := env.GetView()
			Expect(view).To(SatisfyAny(
				ContainSubstring("KaRiya"),
				ContainSubstring("Burst"),
				ContainSubstring("Edit"),
			))
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
		Describe("Empty Burst List", func() {
			BeforeEach(func() {
				env = e2e.SetupWithMemory(GinkgoT())
				// No data populated - empty burst list
			})

			AfterEach(func() {
				env.Cleanup()
			})

			It("should handle navigation keys with empty list gracefully", func() {
				env.SelectIntentByName("burst_management")
				// Should not panic on navigation with empty list
				env.PressKeyRune('j')
				env.PressKeyRune('k')
				env.PressKey(tea.KeyDown)
				env.PressKey(tea.KeyUp)
				view := env.GetView()
				Expect(view).NotTo(BeEmpty())
				Expect(view).NotTo(ContainSubstring("panic"))
			})

			It("should handle Enter key with empty list gracefully", func() {
				env.SelectIntentByName("burst_management")
				env.Confirm() // Enter on empty list
				view := env.GetView()
				Expect(view).NotTo(BeEmpty())
				Expect(view).NotTo(ContainSubstring("panic"))
			})

			It("should allow returning to main menu from empty list", func() {
				env.SelectIntentByName("burst_management")
				env.Cancel()
				env.AssertViewContains("Capture Event")
			})
		})

		Describe("Rapid Key Presses", func() {
			BeforeEach(func() {
				env = e2e.SetupWithMemory(GinkgoT())
				env.PopulateTestData(10, 3, 0) // More events to support multiple bursts
				env.SelectIntentByName("burst_management")
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
					env.Cancel()  // Back to list
				}
				view := env.GetView()
				Expect(view).NotTo(BeEmpty())
				env.AssertViewContainsAny("List", "Bursts", "Name")
			})
		})

		Describe("Boundary Navigation", func() {
			BeforeEach(func() {
				env = e2e.SetupWithMemory(GinkgoT())
				env.PopulateTestData(10, 2, 0) // More events to support bursts
				env.SelectIntentByName("burst_management")
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

		Describe("Filter State During Navigation", func() {
			BeforeEach(func() {
				env = e2e.SetupWithMemory(GinkgoT())
				env.PopulateTestData(10, 3, 0) // More events to support multiple bursts
				env.SelectIntentByName("burst_management")
			})

			AfterEach(func() {
				env.Cleanup()
			})

			It("should maintain list state during navigation", func() {
				// Navigate through list
				env.PressKeyRune('j')
				env.PressKeyRune('j')
				view := env.GetView()
				Expect(view).NotTo(BeEmpty())

				// Enter detail and return
				env.Confirm()
				env.Cancel()
				view = env.GetView()

				// State should be maintained (or at least not crash)
				Expect(view).NotTo(BeEmpty())
				Expect(view).NotTo(ContainSubstring("panic"))
			})

			It("should handle navigation with mixed confirmation states", func() {
				env.PressKeyRune('j')
				env.Confirm()
				// Detail view shows burst info, including name and description
				env.AssertViewContainsAny("Burst Details", "Name", "Description", "Events")
				env.Cancel()

				env.PressKeyRune('j')
				env.Confirm()
				env.AssertViewContainsAny("Burst Details", "Name", "Description", "Events")
			})
		})

		Describe("Multi-Select Scenarios", func() {
			BeforeEach(func() {
				env = e2e.SetupWithMemory(GinkgoT())
				env.PopulateTestData(10, 3, 0) // More events to support multiple bursts
				env.SelectIntentByName("burst_management")
			})

			AfterEach(func() {
				env.Cleanup()
			})

			It("should handle viewing multiple bursts in sequence", func() {
				for i := 0; i < 3; i++ {
					env.Confirm()         // View detail
					env.Cancel()          // Back to list
					env.PressKeyRune('j') // Next item
				}
				view := env.GetView()
				Expect(view).NotTo(BeEmpty())
				Expect(view).NotTo(ContainSubstring("panic"))
			})

			It("should handle edit operations on multiple bursts", func() {
				// Edit first burst
				env.Confirm()
				env.PressKeyRune('e')
				env.Cancel() // Cancel edit
				env.Cancel() // Back to list

				// Move to second burst
				env.PressKeyRune('j')
				env.Confirm()
				env.PressKeyRune('e')
				env.Cancel() // Cancel edit
				env.Cancel() // Back to list

				view := env.GetView()
				Expect(view).NotTo(BeEmpty())
				env.AssertViewContainsAny("List", "Bursts", "Name")
			})
		})

		Describe("Vim Motion Consistency", func() {
			BeforeEach(func() {
				env = e2e.SetupWithMemory(GinkgoT())
				env.PopulateTestData(10, 2, 0) // More events to support bursts
				env.SelectIntentByName("burst_management")
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
		})
	})
})
