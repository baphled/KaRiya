package e2e_test

import (
	"github.com/baphled/kariya/internal/testutil/e2e"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("E2E ExportArtifact Workflow", func() {
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
			env.AssertViewContainsAny("Export", "Type", "Select", "events", "facts", "bursts")
		})

		It("should show context help for type selection", func() {
			env.SelectIntentByName("export_artifact")
			env.AssertViewContainsAny("Enter", "Esc", "q", "Select")
		})
	})

	Describe("Type Selection State", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.SelectIntentByName("export_artifact")
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show artifact type options", func() {
			env.AssertViewContainsAny("events", "facts", "bursts", "Events", "Facts", "Bursts")
		})

		It("should allow navigating type options with j/k", func() {
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

		It("should proceed to format selection when type is selected", func() {
			env.Confirm()
			env.AssertViewContainsAny("Format", "Select", "JSON", "CSV", "YAML", "TXT", "json", "csv")
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

	Describe("Format Selection State", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.SelectIntentByName("export_artifact")
			env.Confirm() // Select first type (events)
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show format options", func() {
			env.AssertViewContainsAny("Format", "JSON", "CSV", "YAML", "TXT", "json", "csv", "yaml", "txt")
		})

		It("should allow navigating format options", func() {
			env.PressKeyRune('j')
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
		})

		It("should go back to type selection when pressing Escape", func() {
			env.Cancel()
			env.AssertViewContainsAny("Type", "events", "facts", "bursts", "Events", "Facts", "Bursts")
		})

		It("should proceed to destination selection when format is selected", func() {
			env.Confirm()
			env.AssertViewContainsAny("Destination", "file", "clipboard", "File", "Clipboard")
		})

		It("should cancel and return to menu when pressing 'm'", func() {
			env.PressKeyRune('m')
			env.AssertViewContains("Capture Event")
		})
	})

	Describe("Destination Selection State", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
			env.SelectIntentByName("export_artifact")
			env.Confirm() // Select type
			env.Confirm() // Select format
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show destination options", func() {
			env.AssertViewContainsAny("Destination", "file", "clipboard", "File", "Clipboard")
		})

		It("should go back to format selection when pressing Escape", func() {
			env.Cancel()
			env.AssertViewContainsAny("Format", "JSON", "CSV", "json", "csv")
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

		It("should cancel from type selection with Escape", func() {
			env.SelectIntentByName("export_artifact")
			env.Cancel()
			env.AssertViewContains("Capture Event")
		})

		It("should go back from format to type with Escape", func() {
			env.SelectIntentByName("export_artifact")
			env.Confirm() // Go to format selection
			env.Cancel()  // Go back to type
			env.AssertViewContainsAny("Type", "events", "facts", "bursts")
		})

		It("should go back from destination to format with Escape", func() {
			env.SelectIntentByName("export_artifact")
			env.Confirm() // Go to format selection
			env.Confirm() // Go to destination selection
			env.Cancel()  // Go back to format
			env.AssertViewContainsAny("Format", "JSON", "CSV", "json", "csv")
		})

		It("should cancel completely with 'm' from destination", func() {
			env.SelectIntentByName("export_artifact")
			env.Confirm()         // Type
			env.Confirm()         // Format
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

		It("should render type selection without panics", func() {
			env.SelectIntentByName("export_artifact")
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should render format selection without panics", func() {
			env.SelectIntentByName("export_artifact")
			env.Confirm()
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should render destination selection without panics", func() {
			env.SelectIntentByName("export_artifact")
			env.Confirm() // Type
			env.Confirm() // Format
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should show breadcrumbs or context", func() {
			env.SelectIntentByName("export_artifact")
			env.AssertViewContainsAny("Export", "Artifact", "Main Menu", "Type")
		})

		It("should show footer with navigation hints", func() {
			env.SelectIntentByName("export_artifact")
			env.AssertViewContainsAny("q", "Esc", "Enter", "Quit")
		})
	})

	Describe("Vim-style Navigation", func() {
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

	Describe("Arrow Key Navigation", func() {
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

	Describe("Export Type Selection", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should allow selecting events type", func() {
			env.SelectIntentByName("export_artifact")
			// Events is typically first
			env.Confirm()
			env.AssertViewContainsAny("Format", "JSON", "CSV")
		})

		It("should allow selecting facts type", func() {
			env.SelectIntentByName("export_artifact")
			env.PressKeyRune('j') // Navigate to facts
			env.Confirm()
			env.AssertViewContainsAny("Format", "JSON", "CSV")
		})

		It("should allow selecting bursts type", func() {
			env.SelectIntentByName("export_artifact")
			env.PressKeyRune('j')
			env.PressKeyRune('j') // Navigate to bursts
			env.Confirm()
			env.AssertViewContainsAny("Format", "JSON", "CSV")
		})
	})

	Describe("With Career Data", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			env.PopulateTestData(5, 2, 3) // 5 events, 2 bursts, 3 facts
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should have data available for export", func() {
			env.AssertEventCount(5)
			env.AssertBurstCount(2)
			env.AssertFactCount(3)
		})

		It("should allow navigating through export workflow with data", func() {
			env.SelectIntentByName("export_artifact")
			env.AssertViewContainsAny("events", "facts", "bursts", "Events", "Facts", "Bursts")
			env.Confirm() // Select type
			env.AssertViewContainsAny("Format", "JSON", "CSV")
			env.Confirm() // Select format
			env.AssertViewContainsAny("Destination", "file", "clipboard")
		})
	})

	Describe("Session Persistence", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should have data available after restart", func() {
			// Use 5 events to ensure enough for 2 bursts (each needs at least 2 events)
			env.PopulateTestData(5, 2, 2)
			env.SimulateRestart()

			env.AssertEventCount(5)
			env.AssertBurstCount(2)
			env.AssertFactCount(2)

			env.SelectIntentByName("export_artifact")
			env.AssertViewContainsAny("events", "facts", "bursts", "Type")
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
			env.SelectIntentByName("export_artifact")
			env.Cancel()
			env.SelectIntentByName("export_artifact")
			env.AssertViewContainsAny("events", "facts", "bursts", "Type")
		})

		It("should allow full navigation through states", func() {
			env.SelectIntentByName("export_artifact")
			env.Confirm() // Type
			env.Confirm() // Format
			// At destination - can go back
			env.Cancel()
			env.AssertViewContainsAny("Format", "JSON", "CSV")
			env.Cancel()
			env.AssertViewContainsAny("Type", "events", "facts")
		})
	})
})
