package e2e_test

import (
	"github.com/baphled/kariya/internal/testutil/e2e"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("E2E Export Workflow", func() {
	var env *e2e.TestEnv

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
			// Wizard Step 1: Type selection (events is default)
			env.AssertViewContainsAny("events", "facts", "bursts", "Events", "Facts", "Bursts", "Career Events", "Export")
			// Use Ctrl+S to skip wizard and go directly to preview with defaults
			env.PressKey(tea.KeyCtrlS)
			// Preview screen shows the export content
			env.AssertViewContainsAny("Preview", "Export", "preview", "Confirm", "[", "{")
		})
	})

	Describe("Escape Navigation", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should return to main menu when esc is pressed at wizard step 1", func() {
			env.SelectIntentByName("export_artifact")
			// Wizard Step 1 is shown
			env.AssertViewContainsAny("Career Events", "Export", "What to Export")
			// Press Escape to cancel
			env.Cancel()
			// Should return to main menu
			env.AssertViewContainsAny("Capture Event", "Browse Timeline", "Main Menu")
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
			// Wizard Step 1 shows type selection
			env.AssertViewContainsAny("events", "facts", "bursts", "Career Events", "Export")
		})
	})

})
