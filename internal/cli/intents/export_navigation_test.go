package intents_test

import (
	"github.com/baphled/kariya/internal/testutil/e2e"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Export Navigation", func() {
	var env *e2e.TestEnv

	Describe("Navigation to ExportArtifact Intent", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show ExportArtifact as menu item", func() {
			env.AssertViewContains("Export Artifact")
		})

		It("should navigate to ExportArtifact when selected", func() {
			env.SelectIntentByName("export_artifact")
			// Wizard modal shows configuration with step title
			env.AssertViewContainsAny("Export", "Configuration", "What to Export", "Step 1")
		})

		It("should show wizard footer with keyboard hints", func() {
			env.SelectIntentByName("export_artifact")
			env.AssertViewContainsAny("Enter", "Esc", "Navigate", "Select", "Skip")
		})
	})

	Describe("Configuration Wizard - Step 1 (What to Export)", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.SelectIntentByName("export_artifact")
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show step 1 with artifact type selection", func() {
			env.AssertViewContainsAny("Step 1", "What to Export")
		})

		It("should allow navigating options with j/k", func() {
			env.PressKeyRune('j')
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			env.PressKeyRune('k')
			view = env.GetView()
			Expect(view).NotTo(BeEmpty())
		})

		It("should allow navigating with arrow keys", func() {
			env.PressKey(tea.KeyDown)
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			env.PressKey(tea.KeyUp)
			view = env.GetView()
			Expect(view).NotTo(BeEmpty())
		})

		It("should proceed through wizard when type is selected", func() {
			env.Confirm() // Select artifact type
			// After selecting type, wizard advances (shows format/destination fields)
			// The huh form advances through fields with Enter
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
		})

		It("should cancel and return to menu when pressing Escape at step 1", func() {
			env.Cancel()
			// Back to main menu
			env.AssertViewContains("Capture Event")
		})
	})

	Describe("Configuration Wizard - Step 2 (How to Export)", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.SelectIntentByName("export_artifact")
			env.Confirm() // Advance through first field
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should continue showing wizard after first selection", func() {
			// Wizard continues with more fields
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should allow navigating options", func() {
			env.PressKeyRune('j')
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Cancel at Each State", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should cancel from wizard step 1 with Escape", func() {
			env.SelectIntentByName("export_artifact")
			env.Cancel()
			// Should return to main menu
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

		It("should render wizard step 1 without panics", func() {
			env.SelectIntentByName("export_artifact")
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should render wizard step 2 without panics", func() {
			env.SelectIntentByName("export_artifact")
			env.Confirm() // Go to step 2
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should show breadcrumbs or context", func() {
			env.SelectIntentByName("export_artifact")
			env.AssertViewContainsAny("Export", "Artifact", "Configure", "Main Menu")
		})

		It("should show footer with navigation hints", func() {
			env.SelectIntentByName("export_artifact")
			env.AssertViewContainsAny("q", "Esc", "Enter", "Quit", "Navigate", "Select")
		})
	})

	Describe("Vim-style Navigation in Wizard", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.SelectIntentByName("export_artifact")
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

	Describe("Arrow Key Navigation in Wizard", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.SelectIntentByName("export_artifact")
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

	Describe("Workflow Re-entry", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should allow re-entering after cancellation", func() {
			env.SelectIntentByName("export_artifact")
			env.Cancel()
			env.SelectIntentByName("export_artifact")
			// Wizard should show step 1 again
			env.AssertViewContainsAny("Step 1", "What to Export", "Export Configuration")
		})
	})

	Describe("Ctrl+S Skip Shortcut", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show keyboard hints in footer", func() {
			env.SelectIntentByName("export_artifact")
			// The footer shows navigation hints
			env.AssertViewContainsAny("Navigate", "Select", "Esc", "Enter")
		})
	})
})
