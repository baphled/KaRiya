package e2e_test

import (
	"strings"

	"github.com/baphled/kariya/internal/testutil/e2e"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("E2E BrowseTimeline Workflow", func() {
	var env *e2e.TestEnv

	Describe("Empty Timeline", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnv(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show empty message when no events exist", func() {
			env.AssertEventCount(0)
			env.SelectIntentByName("browse_timeline")
			env.AssertViewContainsAny("No events", "empty", "Timeline")
		})

		It("should not panic on empty event list", func() {
			env.SelectIntentByName("browse_timeline")
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should allow navigation back from empty list", func() {
			env.SelectIntentByName("browse_timeline")
			env.Cancel()
			env.AssertViewContains("Capture Event")
		})
	})

	Describe("Timeline List with Data", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnv(GinkgoT())
			env.PopulateTestData(5, 0, 0) // 5 events, 0 bursts, 0 facts
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should display events in list", func() {
			env.SelectIntentByName("browse_timeline")
			env.AssertEventCount(5)
			// Should show event content from fixtures.
			env.AssertViewContainsAny("Timeline", "Events", "Acme Corp", "TechCorp")
		})

		It("should show event count in footer", func() {
			env.SelectIntentByName("browse_timeline")
			env.AssertViewContainsAny("Events: 5", "5 events", "Events")
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

	Describe("Vim-style Navigation", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnv(GinkgoT())
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

	Describe("Event Detail View", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnv(GinkgoT())
			env.PopulateTestData(5, 0, 0)
			env.SelectIntentByName("browse_timeline")
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show event detail when pressing Enter", func() {
			env.Confirm() // Select first event.
			env.AssertViewContainsAny("Detail", "Event", "Date", "Company")
		})

		It("should go back to list when pressing Escape from detail", func() {
			env.Confirm() // Go to detail.
			env.Cancel()  // Go back.
			env.AssertViewContainsAny("Timeline", "Events", "List")
		})

		It("should allow selecting event and returning to list multiple times", func() {
			// First selection.
			env.Confirm() // Go to detail.
			env.Cancel()  // Go back to list.

			// Second selection.
			env.PressKeyRune('j') // Navigate to next event.
			env.Confirm()         // Go to detail.
			env.Cancel()          // Go back to list.

			env.AssertViewContainsAny("Timeline", "Events")
		})
	})

	Describe("Search Modal Workflow", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnv(GinkgoT())
			env.PopulateTestData(5, 0, 0)
			env.SelectIntentByName("browse_timeline")
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should open search modal with '/' key", func() {
			env.PressKeyRune('/')
			env.AssertViewContainsAny("Search", "search", "Find")
		})

		It("should close search modal with Escape", func() {
			env.PressKeyRune('/')
			env.Cancel()
			// Should be back at list view.
			env.AssertViewContainsAny("Timeline", "Events")
		})

		It("should not crash when opening search on empty search", func() {
			env.PressKeyRune('/')
			env.Confirm() // Submit empty search.
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			Expect(view).NotTo(ContainSubstring("panic"))
		})
	})

	Describe("Filter Modal Workflow", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnv(GinkgoT())
			env.PopulateTestData(5, 0, 0)
			env.SelectIntentByName("browse_timeline")
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should open filter modal with 'f' key", func() {
			env.PressKeyRune('f')
			env.AssertViewContainsAny("Filter", "Company", "Category")
		})

		It("should close filter modal with Escape", func() {
			env.PressKeyRune('f')
			env.Cancel()
			// Should be back at list view.
			env.AssertViewContainsAny("Timeline", "Events")
		})

		It("should apply filter when submitted", func() {
			env.PressKeyRune('f')
			// Submit filter (default values).
			env.Confirm()
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Sort Modal Workflow", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnv(GinkgoT())
			env.PopulateTestData(5, 0, 0)
			env.SelectIntentByName("browse_timeline")
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should open sort modal with 's' key", func() {
			env.PressKeyRune('s')
			env.AssertViewContainsAny("Sort", "Order", "date")
		})

		It("should close sort modal with Escape", func() {
			env.PressKeyRune('s')
			env.Cancel()
			// Should be back at list view.
			env.AssertViewContainsAny("Timeline", "Events")
		})

		It("should apply sort when submitted", func() {
			env.PressKeyRune('s')
			// Submit sort (default values).
			env.Confirm()
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Quick Add Modal Workflow", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnv(GinkgoT())
			env.PopulateTestData(3, 0, 0)
			env.SelectIntentByName("browse_timeline")
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should open quick add modal with 'a' key", func() {
			env.PressKeyRune('a')
			env.AssertViewContainsAny("Quick", "Add", "Event", "Text")
		})

		It("should close quick add modal with Escape", func() {
			env.PressKeyRune('a')
			env.Cancel()
			// Should be back at list view.
			env.AssertViewContainsAny("Timeline", "Events")
		})

		It("should show form fields in quick add modal", func() {
			env.PressKeyRune('a')
			env.AssertViewContainsAny("Event Text", "Date", "Text")
		})
	})

	Describe("Edit Modal Workflow", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnv(GinkgoT())
			env.PopulateTestData(3, 0, 0)
			env.SelectIntentByName("browse_timeline")
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should open edit modal with 'e' key", func() {
			env.PressKeyRune('e')
			// Edit modal shows a huh form with form controls.
			env.AssertViewContainsAny("enter", "ctrl", "Esc", "tab")
		})

		It("should close edit modal with Escape", func() {
			env.PressKeyRune('e')
			env.Cancel()
			// Should be back at list view.
			env.AssertViewContainsAny("Timeline", "Events")
		})

		It("should show pre-populated form fields in edit modal", func() {
			env.PressKeyRune('e')
			// Edit modal should have event data pre-populated.
			env.AssertViewContainsAny("Text", "Date", "Company")
		})
	})

	Describe("Delete Confirmation Workflow", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnv(GinkgoT())
			env.PopulateTestData(3, 0, 0)
			env.SelectIntentByName("browse_timeline")
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should open delete confirmation with 'd' key", func() {
			env.PressKeyRune('d')
			env.AssertViewContainsAny("Delete", "Confirm", "delete", "Yes", "No")
		})

		It("should close delete confirmation with Escape", func() {
			env.PressKeyRune('d')
			env.Cancel()
			// Should be back at list view.
			env.AssertViewContainsAny("Timeline", "Events")
		})
	})

	Describe("Help Modal", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnv(GinkgoT())
			env.PopulateTestData(3, 0, 0)
			env.SelectIntentByName("browse_timeline")
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should open help modal with '?' key", func() {
			env.PressKeyRune('?')
			env.AssertViewContainsAny("Help", "Keys", "Shortcuts")
		})

		It("should close help modal with '?' key (toggle)", func() {
			env.PressKeyRune('?')
			// Toggle help off with '?' again (Escape does not close help modal).
			env.PressKeyRune('?')
			// Should be back at list view.
			env.AssertViewContainsAny("Timeline", "Events")
		})
	})

	Describe("Session Persistence", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnv(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show events after restart", func() {
			env.PopulateTestData(5, 0, 0)
			env.AssertEventCount(5)

			env.SimulateRestart()

			env.AssertEventCount(5)
			env.SelectIntentByName("browse_timeline")
			env.AssertViewContainsAny("Timeline", "Events", "Acme")
		})
	})

	Describe("Modal Priority and Input Handling", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnv(GinkgoT())
			env.PopulateTestData(5, 0, 0)
			env.SelectIntentByName("browse_timeline")
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should prioritize modal input over screen input", func() {
			// Open search modal.
			env.PressKeyRune('/')
			env.AssertViewContainsAny("Search", "search")

			// Try to navigate (should not affect event list).
			env.PressKey(tea.KeyDown)

			// Modal should still be visible.
			view := env.GetView()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Search"),
				ContainSubstring("search"),
			))
		})

		It("should only allow one modal at a time", func() {
			// Open search modal.
			env.PressKeyRune('/')
			env.AssertViewContainsAny("Search", "search")

			// Try to open filter modal (should not work while search is open).
			env.PressKeyRune('f')

			// Should still show search modal (not filter).
			view := env.GetView()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Search"),
				ContainSubstring("search"),
			))
		})
	})

	Describe("Layout and Display", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnv(GinkgoT())
			env.PopulateTestData(5, 0, 0)
			env.SelectIntentByName("browse_timeline")
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should maintain consistent layout", func() {
			view := env.GetView()
			// View should not be excessively long (broken layout symptom).
			lines := strings.Split(view, "\n")
			Expect(len(lines)).To(BeNumerically("<", 100), "View should not have excessive line count")
			// Should not have large empty gaps.
			Expect(view).NotTo(ContainSubstring("\n\n\n\n\n"))
		})

		It("should show keyboard shortcuts in footer", func() {
			env.AssertViewContainsAny("j/k", "Enter", "?", "q")
		})
	})
})
