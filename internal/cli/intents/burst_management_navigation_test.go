package intents_test

import (
	"github.com/baphled/kariya/internal/testutil/e2e"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Burstmanagement Navigation", func() {
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

})
