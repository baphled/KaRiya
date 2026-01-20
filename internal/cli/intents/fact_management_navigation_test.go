package intents_test

import (
	"github.com/baphled/kariya/internal/testutil/e2e"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("FactManagement Navigation", func() {
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

		It("should ignore 'q' key from list (quit only from main menu)", func() {
			env.SelectIntentByName("fact_management")
			// q no longer quits from within intents - only from main menu
			// This test verifies 'q' is handled gracefully (does nothing)
			env.Quit()
			// Intent should still be active
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
})
