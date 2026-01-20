package e2e_test

import (
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/testutil/e2e"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("E2E Configure Workflow", func() {
	var env *e2e.TestEnv

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

	Describe("Modal Sequence Flow", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		Context("Domain Selection", func() {
			It("should show domain selection screen when entering ConfigureSystem", func() {
				env.SelectIntentByName("configure_system")

				// Should show all 4 domains
				env.AssertViewContainsAny("System", "system")
				env.AssertViewContainsAny("Profile", "profile")
				env.AssertViewContainsAny("Export", "export")
				env.AssertViewContainsAny("UI", "ui")
			})

			It("should navigate between domains with j/k keys", func() {
				env.SelectIntentByName("configure_system")

				// Initial selection should be first item (System)
				view1 := env.GetView()
				Expect(view1).To(ContainSubstring("System"))

				// Navigate down
				env.NavigateDown()

				// Should now highlight Profile (second item)
				view2 := env.GetView()
				Expect(view2).To(ContainSubstring("Profile"))

				// Navigate up
				env.NavigateUp()

				// Should be back at System
				view3 := env.GetView()
				Expect(view3).To(ContainSubstring("System"))
			})

			It("should cancel intent when pressing Escape at domain selection", func() {
				env.SelectIntentByName("configure_system")
				env.AssertViewContainsAny("System", "Domain")

				// Press Escape to cancel
				env.Cancel()

				// Should return to main menu
				env.AssertViewContainsAny("Capture Event", "Career Event")
			})

			It("should cancel intent when pressing q at domain selection", func() {
				env.SelectIntentByName("configure_system")
				env.AssertViewContainsAny("System", "Domain")

				// Press q to quit/cancel
				env.Quit()

				// Should return to main menu
				env.AssertViewContainsAny("Capture Event", "Career Event")
			})
		})

		Context("Edit Settings Modal", func() {
			It("should open edit modal when selecting a domain", func() {
				env.SelectIntentByName("configure_system")

				// Select System domain
				env.Confirm() // Press Enter on first domain

				// Should show edit modal with settings
				env.AssertViewContainsAny("Edit", "Settings", "Auto Save", "Backup")
			})

			It("should show domain-specific settings in edit modal", func() {
				env.SelectIntentByName("configure_system")

				// Navigate to Profile domain
				env.NavigateDown() // Profile is second

				// Select Profile domain
				env.Confirm()

				// Should show profile settings
				env.AssertViewContainsAny("Name", "Email", "Display")
			})

			It("should close edit modal and return to domain selection on Escape", func() {
				env.SelectIntentByName("configure_system")

				// Select System domain
				env.Confirm()

				// Should be in edit modal
				env.AssertViewContainsAny("Edit", "Settings", "Auto Save")

				// Press Escape to close modal
				env.Cancel()

				// Should return to domain selection
				env.AssertViewContainsAny("System", "Profile", "Export", "UI")
			})
		})

		Context("Review Changes Modal", func() {
			It("should transition from edit to review when submitting form", func() {
				env.SelectIntentByName("configure_system")

				// Select System domain
				env.Confirm()

				// Should be in edit modal
				env.AssertViewContainsAny("Edit", "Settings")

				// Submit the form (Ctrl+S or form completion)
				// The edit modal uses huh form which submits on Enter when complete
				// or we can send Ctrl+S
				env.PressKey(tea.KeyCtrlS)

				// Should show review modal
				env.AssertViewContainsAny("Review", "Changes")
			})

			It("should return to edit modal from review on Escape", func() {
				env.SelectIntentByName("configure_system")

				// Select domain and submit form
				env.Confirm()
				env.PressKey(tea.KeyCtrlS)

				// Should be in review modal
				env.AssertViewContainsAny("Review", "Changes")

				// Press Escape to go back
				env.Cancel()

				// Should return to edit modal
				env.AssertViewContainsAny("Edit", "Settings")
			})
		})

		Context("Confirm Modal", func() {
			It("should transition from review to confirm when confirming review", func() {
				env.SelectIntentByName("configure_system")

				// Navigate through: domain → edit → review → confirm
				env.Confirm()              // Select domain
				env.PressKey(tea.KeyCtrlS) // Submit edit form

				// Should be in review
				env.AssertViewContainsAny("Review", "Changes")

				// Confirm review (Enter or y)
				env.Confirm()

				// Should show confirm modal
				env.AssertViewContainsAny("Confirm", "Are you sure", "save")
			})

			It("should return to review from confirm on Escape", func() {
				env.SelectIntentByName("configure_system")

				// Navigate to confirm
				env.Confirm()              // domain
				env.PressKey(tea.KeyCtrlS) // edit → review
				env.Confirm()              // review → confirm

				// Should be in confirm
				env.AssertViewContainsAny("Confirm", "Are you sure")

				// Press Escape
				env.Cancel()

				// Should return to review
				env.AssertViewContainsAny("Review", "Changes")
			})
		})

		Context("Complete Save Flow", func() {
			It("should complete save flow and show success modal", func() {
				env.SelectIntentByName("configure_system")

				// Navigate through entire flow
				env.Confirm()              // domain
				env.PressKey(tea.KeyCtrlS) // edit → review
				env.Confirm()              // review → confirm

				// Confirm save (y or Enter)
				env.PressKeyRune('y')

				// Should show saving or success
				env.AssertViewContainsAny("Saving", "saved", "Success", "Complete")
			})

			It("should return to main menu after success confirmation", func() {
				env.SelectIntentByName("configure_system")

				// Complete the flow
				env.Confirm()              // domain
				env.PressKey(tea.KeyCtrlS) // edit → review
				env.Confirm()              // review → confirm
				env.PressKeyRune('y')      // confirm → saving → success

				// Dismiss success modal
				env.Confirm()

				// Should return to main menu
				env.AssertViewContainsAny("Capture Event", "Career Event", "Main Menu")
			})
		})

		Context("Escape Navigation at Each Step", func() {
			It("should support full escape path: confirm → review → edit → domain → cancel", func() {
				env.SelectIntentByName("configure_system")

				// Navigate to confirm state
				env.Confirm()              // domain → edit
				env.PressKey(tea.KeyCtrlS) // edit → review
				env.Confirm()              // review → confirm

				// Should be at confirm
				env.AssertViewContainsAny("Confirm", "Are you sure")

				// Escape back to review
				env.Cancel()
				env.AssertViewContainsAny("Review", "Changes")

				// Escape back to edit
				env.Cancel()
				env.AssertViewContainsAny("Edit", "Settings")

				// Escape back to domain
				env.Cancel()
				env.AssertViewContainsAny("System", "Profile", "Export", "UI")

				// Escape to cancel intent
				env.Cancel()
				env.AssertViewContainsAny("Capture Event", "Career Event")
			})
		})

		Context("Different Domains", func() {
			It("should show Export settings when selecting Export domain", func() {
				env.SelectIntentByName("configure_system")

				// Navigate to Export domain (third item)
				env.NavigateDown() // Profile
				env.NavigateDown() // Export

				// Select Export
				env.Confirm()

				// Should show export-specific settings
				env.AssertViewContainsAny("Export", "Format", "Destination", "markdown", "json")
			})

			It("should show UI settings when selecting UI domain", func() {
				env.SelectIntentByName("configure_system")

				// Navigate to UI domain (fourth item)
				env.NavigateDown() // Profile
				env.NavigateDown() // Export
				env.NavigateDown() // UI

				// Select UI
				env.Confirm()

				// Should show UI-specific settings
				env.AssertViewContainsAny("UI", "Theme", "theme", "dark", "light")
			})
		})
	})

	Describe("Direct Message Testing", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should handle ConfigCompleteMsg for success", func() {
			env.SelectIntentByName("configure_system")

			// Navigate to saving state
			env.Confirm()              // domain
			env.PressKey(tea.KeyCtrlS) // edit → review
			env.Confirm()              // review → confirm
			env.PressKeyRune('y')      // start saving

			// Manually send ConfigCompleteMsg to simulate successful save
			env.SendMessage(intents.ConfigCompleteMsg{
				Result: &intents.ConfigureSystemResult{
					Success: true,
					Domain:  intents.DomainSystem,
				},
			})

			// Should show success state
			env.AssertViewContainsAny("Success", "saved", "Complete", "Configuration")
		})

		It("should handle ConfigErrorMsg for failure", func() {
			env.SelectIntentByName("configure_system")

			// Just select a domain to be in the intent
			env.Confirm() // domain → edit

			// Send ConfigErrorMsg directly - this should trigger error modal
			// regardless of current state due to global message handling
			env.SendMessage(intents.ConfigErrorMsg{
				Error: &intents.IntentError{
					Code:    "save_failed",
					Message: "Test error message",
				},
			})

			// Should show error modal
			env.AssertViewContainsAny("Failed", "Save Failed", "Test error")
		})
	})
})
