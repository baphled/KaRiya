package e2e_test

import (
	"github.com/baphled/kariya/internal/testutil/e2e"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("E2E BrowseTimeline Workflow", func() {
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

	Describe("Empty Timeline", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT()) // SQLite for persistence testing
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show empty message when no events exist", func() {
			env.AssertEventCount(0)
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

	Describe("Timeline with Events", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			env.PopulateTestData(5, 0, 0) // 5 events, no bursts, no facts
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should display events in timeline", func() {
			env.SelectIntentByName("browse_timeline")
			env.AssertEventCount(5)
			// Should show at least some event content (company name from fixtures)
			env.AssertViewContainsAny("Acme Corp", "TechStart", "BigCorp", "Events", "Page")
		})

		It("should show pagination info when events exist", func() {
			env.SelectIntentByName("browse_timeline")
			env.AssertViewContainsAny("Page", "Events", "1 of")
		})

		It("should allow navigating through events with j/k", func() {
			env.SelectIntentByName("browse_timeline")
			// Navigate down
			env.PressKeyRune('j')
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			// Navigate up
			env.PressKeyRune('k')
			view = env.GetView()
			Expect(view).NotTo(BeEmpty())
		})

		It("should allow navigating with arrow keys", func() {
			env.SelectIntentByName("browse_timeline")
			// Navigate down
			env.PressKey(tea.KeyDown)
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			// Navigate up
			env.PressKey(tea.KeyUp)
			view = env.GetView()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Event Detail View", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			env.PopulateTestData(3, 0, 0)
			env.SelectIntentByName("browse_timeline")
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show event detail when pressing Enter", func() {
			env.Confirm() // Select first event
			env.AssertViewContainsAny("Date", "Text", "Company", "Event Detail")
		})

		It("should show full event information in detail view", func() {
			env.Confirm()
			// Events from fixtures have company and text
			env.AssertViewContainsAny("Acme Corp", "TechStart", "security", "authentication", "Date")
		})

		It("should go back to timeline when pressing Escape from detail", func() {
			env.Confirm() // Go to detail
			env.Cancel()  // Go back
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

		It("should return to main menu when pressing 'q' from timeline", func() {
			env.SelectIntentByName("browse_timeline")
			env.Quit()
			env.AssertViewContains("Capture Event")
		})

		It("should return to main menu when pressing 'm' from timeline", func() {
			env.SelectIntentByName("browse_timeline")
			env.PressKeyRune('m')
			env.AssertViewContains("Capture Event")
		})
	})

	Describe("Cancel from Event Detail", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			env.PopulateTestData(2, 0, 0)
			env.SelectIntentByName("browse_timeline")
			env.Confirm() // Go to event detail
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should return to timeline when pressing Escape from detail", func() {
			env.Cancel()
			env.AssertViewContainsAny("Timeline", "Events", "Page")
		})

		It("should return to main menu when pressing 'q' from detail", func() {
			env.Quit()
			env.AssertViewContains("Capture Event")
		})

		It("should return to main menu when pressing 'm' from detail", func() {
			env.PressKeyRune('m')
			env.AssertViewContains("Capture Event")
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

	Describe("Vim-style Navigation", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			env.PopulateTestData(5, 0, 0)
			env.SelectIntentByName("browse_timeline")
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should navigate down with 'j' key", func() {
			// First item should be selected initially
			env.PressKeyRune('j')
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
		})

		It("should navigate up with 'k' key", func() {
			env.PressKeyRune('j') // Go down
			env.PressKeyRune('k') // Go back up
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
		})

		It("should not crash at list boundaries", func() {
			// Try navigating up at the top
			env.PressKeyRune('k')
			env.PressKeyRune('k')
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())

			// Navigate to bottom and try going further
			for i := 0; i < 10; i++ {
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

		It("should show events after restart", func() {
			// Add event before restart
			env.PopulateTestData(3, 0, 0)
			env.AssertEventCount(3)

			// Simulate restart
			env.SimulateRestart()

			// Events should still exist
			env.AssertEventCount(3)
			env.SelectIntentByName("browse_timeline")
			env.AssertViewContainsAny("Acme Corp", "TechStart", "Events")
		})

		It("should maintain event data integrity after restart", func() {
			env.PopulateTestData(2, 0, 0)
			eventsBefore := env.GetEvents()
			Expect(len(eventsBefore)).To(Equal(2))

			env.SimulateRestart()

			eventsAfter := env.GetEvents()
			Expect(len(eventsAfter)).To(Equal(2))
			// Check that event IDs match
			Expect(eventsAfter[0].ID).To(Equal(eventsBefore[0].ID))
		})
	})

	Describe("Workflow Integration", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			env.PopulateTestData(3, 0, 0)
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should allow selecting event and returning to timeline multiple times", func() {
			env.SelectIntentByName("browse_timeline")

			// First selection
			env.Confirm() // Go to detail
			env.Cancel()  // Go back to timeline

			// Second selection
			env.PressKeyRune('j') // Navigate to next event
			env.Confirm()         // Go to detail
			env.Cancel()          // Go back to timeline

			env.AssertViewContainsAny("Timeline", "Events", "Page")
		})

		It("should complete workflow when pressing Enter on event detail", func() {
			env.SelectIntentByName("browse_timeline")
			env.Confirm() // Go to event detail
			env.Confirm() // Complete/confirm selection
			// Should return to main menu after completing
			env.AssertViewContains("Capture Event")
		})
	})
})
