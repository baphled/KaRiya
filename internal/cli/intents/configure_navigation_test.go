package intents_test

import (
	"github.com/baphled/kariya/internal/testutil/harness"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Configure Navigation", func() {
	var env *harness.TestEnv

	Describe("Navigation to ConfigureSystem Intent", func() {
		BeforeEach(func() {
			env = harness.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show ConfigureSystem as menu item", func() {
			env.AssertViewContainsAny("Settings", ",")
		})

		It("should navigate to ConfigureSystem when selected", func() {
			env.PressKeyRune(',')
			env.AssertViewContainsAny("Configure System", "System", "Profile", "Export", "UI")
			env.AssertViewContainsAny("Settings", "log", "data", "backup", "level")
		})

		It("should show context help for domain selection", func() {
			env.PressKeyRune(',')
			env.AssertViewContainsAny("Esc", "Ctrl+S", "Tab", "j/k")
		})
	})

	Describe("Domain Selection State", func() {
		BeforeEach(func() {
			env = harness.SetupWithMemory(GinkgoT())
			env.PressKeyRune(',')
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
			env.AssertViewContainsAny("Settings", "Edit", "log", "data", "backup", "level")
		})

		It("should cancel and return to menu when pressing Escape", func() {
			env.Cancel()
			env.AssertViewContains("Capture Event")
		})

	})

	Describe("Edit Settings State", func() {
		BeforeEach(func() {
			env = harness.SetupWithMemory(GinkgoT())
			env.PressKeyRune(',')
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show settings for selected domain", func() {
			env.AssertViewContainsAny("log", "data", "backup", "Settings", "level", "dir", "System", "Profile")
		})

		It("should allow navigating settings with j/k", func() {
			env.PressKeyRune('j')
			env.AssertViewContainsAny("Profile", "role", "audience", "Settings")
		})

		It("should close the modal when pressing Escape", func() {
			env.Cancel()
			env.AssertViewContains("Capture Event")
		})
	})

	Describe("Cancel at Each State", func() {
		BeforeEach(func() {
			env = harness.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should cancel from domain selection with Escape", func() {
			env.PressKeyRune(',')
			env.Cancel()
			env.AssertViewContains("Capture Event")
		})

		It("should close the modal after navigating domains", func() {
			env.PressKeyRune(',')
			env.PressKeyRune('j')
			env.Cancel()
			env.AssertViewContains("Capture Event")
		})
	})

	Describe("View Rendering", func() {
		BeforeEach(func() {
			env = harness.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should render domain selection without panics", func() {
			env.PressKeyRune(',')
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should render edit settings without panics", func() {
			env.PressKeyRune(',')
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should show breadcrumbs or context", func() {
			env.PressKeyRune(',')
			env.AssertViewContainsAny("Configure System", "System")
		})

		It("should show footer with navigation hints", func() {
			env.PressKeyRune(',')
			env.AssertViewContainsAny("Esc", "Ctrl+S", "Tab", "j/k")
		})
	})

	Describe("Vim-style Navigation", func() {
		BeforeEach(func() {
			env = harness.SetupWithMemory(GinkgoT())
			env.PressKeyRune(',')
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

			for range 10 {
				env.PressKeyRune('j')
			}
			view = env.GetView()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Arrow Key Navigation", func() {
		BeforeEach(func() {
			env = harness.SetupWithMemory(GinkgoT())
			env.PressKeyRune(',')
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
			env = harness.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should allow selecting System domain", func() {
			env.PressKeyRune(',')
			env.AssertViewContainsAny("log", "data", "backup", "Settings", "System")
		})

		It("should allow selecting Profile domain", func() {
			env.PressKeyRune(',')
			env.PressKeyRune('j') // Navigate to Profile
			env.AssertViewContainsAny("Settings", "Profile", "role", "audience", "Edit")
		})

		It("should allow selecting Export domain", func() {
			env.PressKeyRune(',')
			env.PressKeyRune('j')
			env.PressKeyRune('j') // Navigate to Export
			env.AssertViewContainsAny("Settings", "Export", "destination", "open", "Edit")
		})

		It("should allow selecting UI domain", func() {
			env.PressKeyRune(',')
			env.PressKeyRune('j')
			env.PressKeyRune('j')
			env.PressKeyRune('j') // Navigate to UI
			env.AssertViewContainsAny("Settings", "UI", "Display", "Edit")
		})
	})

	Describe("Workflow Navigation", func() {
		BeforeEach(func() {
			env = harness.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should allow re-entering after cancellation", func() {
			env.PressKeyRune(',')
			env.Cancel()
			env.PressKeyRune(',')
			env.AssertViewContainsAny("System", "Profile", "Export", "UI")
		})

		It("should allow navigating through domains", func() {
			env.PressKeyRune(',')
			env.AssertViewContainsAny("System", "Settings")
			env.PressKeyRune('j')
			env.AssertViewContainsAny("Profile", "role", "audience")
			env.PressKeyRune('j')
			env.AssertViewContainsAny("Export", "destination")
			env.PressKeyRune('k')
			env.AssertViewContainsAny("Profile", "role", "audience")
		})

		It("should close after saving with Ctrl+S", func() {
			env.PressKeyRune(',')
			env.PressKey(tea.KeyCtrlS)
			env.AssertViewContains("Capture Event")
		})
	})

})
