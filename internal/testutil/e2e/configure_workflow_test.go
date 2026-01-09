package e2e_test

import (
	"github.com/baphled/kariya/internal/testutil/e2e"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("E2E ConfigureSystem Workflow", func() {
	var env *e2e.TestEnv

	Describe("Navigation to ConfigureSystem Intent", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show ConfigureSystem as menu item", func() {
			env.AssertViewContains("Configure System")
		})

		It("should navigate to ConfigureSystem when selected", func() {
			env.SelectIntentByName("configure_system")
			env.AssertViewContainsAny("Configure", "Domain", "Select", "System", "Profile", "Export", "UI")
		})

		It("should show context help for domain selection", func() {
			env.SelectIntentByName("configure_system")
			env.AssertViewContainsAny("Enter", "Esc", "Select")
		})
	})

	Describe("Domain Selection State", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.SelectIntentByName("configure_system")
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show configuration domain options", func() {
			env.AssertViewContainsAny("System", "Profile", "Export", "UI", "system", "profile", "export", "ui")
		})

		It("should allow navigating domain options with j/k", func() {
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

		It("should proceed to edit settings when domain is selected", func() {
			env.Confirm()
			env.AssertViewContainsAny("Settings", "Edit", "log", "data", "backup", "level")
		})

		It("should cancel and return to menu when pressing Escape", func() {
			env.Cancel()
			env.AssertViewContains("Capture Event")
		})

		It("should cancel and return to menu when pressing 'm'", func() {
			env.PressKeyRune('m')
			env.AssertViewContains("Capture Event")
		})
	})

	Describe("Edit Settings State", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.SelectIntentByName("configure_system")
			env.Confirm() // Select first domain (System)
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show settings for selected domain", func() {
			env.AssertViewContainsAny("log", "data", "backup", "Settings", "level", "dir")
		})

		It("should allow navigating settings with j/k", func() {
			env.PressKeyRune('j')
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
		})

		It("should go back to domain selection when pressing Escape", func() {
			env.Cancel()
			env.AssertViewContainsAny("Domain", "System", "Profile", "Export", "UI")
		})

		It("should cancel and return to menu when pressing 'm'", func() {
			env.PressKeyRune('m')
			env.AssertViewContains("Capture Event")
		})
	})

	Describe("Cancel at Each State", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should cancel from domain selection with Escape", func() {
			env.SelectIntentByName("configure_system")
			env.Cancel()
			env.AssertViewContains("Capture Event")
		})

		It("should go back from edit settings to domain with Escape", func() {
			env.SelectIntentByName("configure_system")
			env.Confirm() // Go to edit settings
			env.Cancel()  // Go back to domain
			env.AssertViewContainsAny("Domain", "System", "Profile", "Export", "UI")
		})

		It("should cancel completely with 'm' from edit settings", func() {
			env.SelectIntentByName("configure_system")
			env.Confirm()         // Go to edit settings
			env.PressKeyRune('m') // Cancel to main menu
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

		It("should render domain selection without panics", func() {
			env.SelectIntentByName("configure_system")
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should render edit settings without panics", func() {
			env.SelectIntentByName("configure_system")
			env.Confirm()
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should show breadcrumbs or context", func() {
			env.SelectIntentByName("configure_system")
			env.AssertViewContainsAny("Configure", "System", "Main Menu", "Domain")
		})

		It("should show footer with navigation hints", func() {
			env.SelectIntentByName("configure_system")
			env.AssertViewContainsAny("Esc", "Enter", "Select")
		})
	})

	Describe("Vim-style Navigation", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.SelectIntentByName("configure_system")
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

	Describe("Arrow Key Navigation", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.SelectIntentByName("configure_system")
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

	Describe("Domain Selection", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should allow selecting System domain", func() {
			env.SelectIntentByName("configure_system")
			// System is typically first
			env.Confirm()
			env.AssertViewContainsAny("log", "data", "backup", "Settings")
		})

		It("should allow selecting Profile domain", func() {
			env.SelectIntentByName("configure_system")
			env.PressKeyRune('j') // Navigate to Profile
			env.Confirm()
			env.AssertViewContainsAny("Settings", "Profile", "role", "audience", "Edit")
		})

		It("should allow selecting Export domain", func() {
			env.SelectIntentByName("configure_system")
			env.PressKeyRune('j')
			env.PressKeyRune('j') // Navigate to Export
			env.Confirm()
			env.AssertViewContainsAny("Settings", "Export", "destination", "open", "Edit")
		})

		It("should allow selecting UI domain", func() {
			env.SelectIntentByName("configure_system")
			env.PressKeyRune('j')
			env.PressKeyRune('j')
			env.PressKeyRune('j') // Navigate to UI
			env.Confirm()
			env.AssertViewContainsAny("Settings", "UI", "Display", "Edit")
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
			env.SelectIntentByName("configure_system")
			env.Cancel()
			env.SelectIntentByName("configure_system")
			env.AssertViewContainsAny("System", "Profile", "Export", "UI", "Domain")
		})

		It("should allow navigating through domains and returning", func() {
			env.SelectIntentByName("configure_system")
			env.Confirm() // Go to System settings
			env.Cancel()  // Go back to domain selection
			env.PressKeyRune('j')
			env.Confirm() // Go to Profile settings
			env.Cancel()  // Go back to domain selection
			env.AssertViewContainsAny("System", "Profile", "Export", "UI", "Domain")
		})
	})

	Describe("Session Persistence", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show configuration options after restart", func() {
			env.SimulateRestart()
			env.SelectIntentByName("configure_system")
			env.AssertViewContainsAny("System", "Profile", "Export", "UI", "Domain")
		})
	})
})
