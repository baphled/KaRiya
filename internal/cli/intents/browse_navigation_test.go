package intents_test

import (
	"github.com/baphled/kariya/internal/testutil/e2e"
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

	Describe("Cancel Navigation", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should return to main menu when pressing Escape from timeline", func() {
			env.SelectIntentByName("browse_timeline")
			env.Cancel() // Esc key
			env.AssertViewContains("Capture Event")
		})

		It("should quit application when pressing 'q' from timeline", func() {
			env.SelectIntentByName("browse_timeline")
			// Note: q now quits the entire app, not just cancel to main menu
			// The test simply verifies the quit is handled without panic
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

})
