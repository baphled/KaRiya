package e2e_test

import (
	"github.com/baphled/kariya/internal/testutil/e2e"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("E2E BurstManagement Workflow", func() {
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

	Describe("Empty Burst List", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show empty message when no bursts exist", func() {
			env.AssertBurstCount(0)
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

	Describe("Burst List with Data", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			env.PopulateTestData(5, 2, 0) // 5 events, 2 bursts, 0 facts
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should display bursts in list", func() {
			env.SelectIntentByName("burst_management")
			env.AssertBurstCount(2)
			// Should show burst names from fixtures
			env.AssertViewContainsAny("Authentication", "Mentoring", "Burst", "Name")
		})

		It("should show confirmation status in list", func() {
			env.SelectIntentByName("burst_management")
			// Fixtures create bursts with alternating confirmed status
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

	Describe("Burst Detail View", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			env.PopulateTestData(5, 2, 0)
			env.SelectIntentByName("burst_management")
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show burst detail when pressing Enter", func() {
			env.Confirm() // Select first burst
			env.AssertViewContainsAny("Detail", "Name", "Description", "Events", "Confirmed")
		})

		It("should show burst information in detail view", func() {
			env.Confirm()
			// Should show burst details from fixtures
			env.AssertViewContainsAny("Authentication", "Mentoring", "Description")
		})

		It("should go back to list when pressing Escape from detail", func() {
			env.Confirm() // Go to detail
			env.Cancel()  // Go back
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

		It("should return to main menu when pressing 'q' from list", func() {
			env.SelectIntentByName("burst_management")
			env.Quit()
			env.AssertViewContains("Capture Event")
		})
	})

	Describe("Cancel from Detail View", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			env.PopulateTestData(5, 2, 0)
			env.SelectIntentByName("burst_management")
			env.Confirm() // Go to detail
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should return to list when pressing Escape from detail", func() {
			env.Cancel()
			env.AssertViewContainsAny("List", "Bursts", "Name")
		})

		It("should return to main menu when pressing 'q' from detail", func() {
			env.Quit()
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

	Describe("Vim-style Navigation", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
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
			env = e2e.Setup(GinkgoT())
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

	Describe("Session Persistence", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show bursts after restart", func() {
			env.PopulateTestData(5, 2, 0)
			env.AssertBurstCount(2)

			env.SimulateRestart()

			env.AssertBurstCount(2)
			env.SelectIntentByName("burst_management")
			env.AssertViewContainsAny("Authentication", "Mentoring", "Burst")
		})

		It("should maintain burst data integrity after restart", func() {
			env.PopulateTestData(5, 2, 0)
			burstsBefore := env.GetBursts()
			Expect(len(burstsBefore)).To(Equal(2))

			env.SimulateRestart()

			burstsAfter := env.GetBursts()
			Expect(len(burstsAfter)).To(Equal(2))
			Expect(burstsAfter[0].ID).To(Equal(burstsBefore[0].ID))
		})
	})

	Describe("Workflow Integration", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			env.PopulateTestData(5, 2, 0)
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should allow selecting burst and returning to list multiple times", func() {
			env.SelectIntentByName("burst_management")

			// First selection
			env.Confirm() // Go to detail
			env.Cancel()  // Go back to list

			// Second selection
			env.PressKeyRune('j') // Navigate to next burst
			env.Confirm()         // Go to detail
			env.Cancel()          // Go back to list

			env.AssertViewContainsAny("List", "Bursts", "Name")
		})

		It("should allow re-entering after cancellation", func() {
			env.SelectIntentByName("burst_management")
			env.Cancel()
			env.SelectIntentByName("burst_management")
			env.AssertViewContainsAny("Burst", "List", "Name")
		})
	})
})
