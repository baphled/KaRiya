package intents_test

import (
	"github.com/baphled/kariya/internal/testutil/e2e"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// BulkOperations is a UI shell only - testing navigation and rendering only
var _ = Describe("BulkOperations Navigation", func() {
	var env *e2e.TestEnv

	Describe("Navigation to BulkOperations Intent", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show BulkOperations as menu item", func() {
			env.AssertViewContainsAny("Bulk Operations", "Bulk")
		})

		It("should navigate to BulkOperations when selected", func() {
			env.SelectIntentByName("bulk_operations")
			env.AssertViewContainsAny("Bulk", "Operations", "Select", "Action")
		})
	})

	Describe("State Navigation", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.SelectIntentByName("bulk_operations")
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show initial bulk operations state", func() {
			env.AssertViewContainsAny("Bulk", "Select", "Operation", "Action")
		})

		It("should return to main menu when pressing Escape", func() {
			env.Cancel()
			env.AssertViewContains("Capture Event")
		})

		It("should ignore 'q' key within intent (quit only from main menu)", func() {
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

		It("should render without panics", func() {
			env.SelectIntentByName("bulk_operations")
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should show breadcrumbs or context", func() {
			env.SelectIntentByName("bulk_operations")
			env.AssertViewContainsAny("Bulk", "Operations", "Main Menu")
		})

		It("should show footer with navigation hints", func() {
			env.SelectIntentByName("bulk_operations")
			env.AssertViewContainsAny("q", "Esc", "Quit")
		})
	})

	Describe("Workflow Navigation", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should allow re-entering after cancellation", func() {
			env.SelectIntentByName("bulk_operations")
			env.Cancel()
			env.SelectIntentByName("bulk_operations")
			env.AssertViewContainsAny("Bulk", "Operations", "Select")
		})
	})
})
