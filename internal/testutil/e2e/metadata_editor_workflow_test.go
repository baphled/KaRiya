package e2e_test

import (
	"github.com/baphled/kariya/internal/testutil/e2e"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// MetadataEditor is a UI shell only - testing navigation and rendering only
var _ = Describe("E2E MetadataEditor Workflow (Navigation Only)", func() {
	var env *e2e.TestEnv

	Describe("Navigation to MetadataEditor Intent", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show MetadataEditor as menu item", func() {
			env.AssertViewContainsAny("Edit Metadata", "Metadata")
		})

		It("should navigate to MetadataEditor when selected", func() {
			env.SelectIntentByName("metadata_editor")
			env.AssertViewContainsAny("Metadata", "Edit", "Entity", "Select")
		})
	})

	Describe("State Navigation", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.SelectIntentByName("metadata_editor")
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show initial metadata state", func() {
			env.AssertViewContainsAny("Metadata", "Select", "Entity", "Type")
		})

		It("should return to main menu when pressing Escape", func() {
			env.Cancel()
			env.AssertViewContains("Capture Event")
		})

		It("should return to main menu when pressing 'q'", func() {
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

		It("should render without panics", func() {
			env.SelectIntentByName("metadata_editor")
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should show breadcrumbs or context", func() {
			env.SelectIntentByName("metadata_editor")
			env.AssertViewContainsAny("Metadata", "Editor", "Main Menu")
		})

		It("should show footer with navigation hints", func() {
			env.SelectIntentByName("metadata_editor")
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
			env.SelectIntentByName("metadata_editor")
			env.Cancel()
			env.SelectIntentByName("metadata_editor")
			env.AssertViewContainsAny("Metadata", "Edit", "Entity")
		})
	})
})
