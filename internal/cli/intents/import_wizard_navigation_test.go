package intents_test

import (
	"github.com/baphled/kariya/internal/testutil/e2e"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// ImportWizard is a UI shell only - testing navigation and rendering only
var _ = Describe("ImportWizard Navigation", func() {
	var env *e2e.TestEnv

	Describe("Navigation to ImportWizard Intent", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show ImportWizard as menu item", func() {
			env.AssertViewContainsAny("Import Data", "Import")
		})

		It("should navigate to ImportWizard when selected", func() {
			env.SelectIntentByName("import_wizard")
			env.AssertViewContainsAny("Import", "Wizard", "CSV", "File")
		})
	})

	Describe("State Navigation", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.SelectIntentByName("import_wizard")
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show initial import state", func() {
			env.AssertViewContainsAny("Import", "Select", "File", "CSV")
		})

		It("should return to main menu when pressing Escape", func() {
			env.Cancel()
			env.AssertViewContains("Capture Event")
		})

		It("should quit application when pressing 'q'", func() {
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

		It("should render without panics", func() {
			env.SelectIntentByName("import_wizard")
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should show breadcrumbs or context", func() {
			env.SelectIntentByName("import_wizard")
			env.AssertViewContainsAny("Import", "Wizard", "Main Menu")
		})

		It("should show footer with navigation hints", func() {
			env.SelectIntentByName("import_wizard")
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
			env.SelectIntentByName("import_wizard")
			env.Cancel()
			env.SelectIntentByName("import_wizard")
			env.AssertViewContainsAny("Import", "File", "CSV")
		})
	})
})
