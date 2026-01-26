package feedback_test

import (
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("InfoModal E2E", func() {
	Describe("Information Display Workflow", func() {
		It("should display and dismiss info modal on Enter", func() {
			// User is shown an informational message.
			modal := feedback.NewInfoModal("Welcome", "Welcome to KaRiya! Press Enter to continue.")

			// Modal should be visible.
			Expect(modal.IsVisible()).To(BeTrue())
			Expect(modal.GetTitle()).To(Equal("Welcome"))
			Expect(modal.GetVariant()).To(Equal(feedback.InfoModalInfo))

			// View should contain the message.
			view := modal.View()
			Expect(view).To(ContainSubstring("Welcome"))
			Expect(view).To(ContainSubstring("KaRiya"))

			// User presses Enter to dismiss.
			dismissed := modal.Update(tea.KeyMsg{Type: tea.KeyEnter})

			Expect(dismissed).To(BeTrue())
			Expect(modal.IsVisible()).To(BeFalse())
			Expect(modal.View()).To(BeEmpty())
		})

		It("should display and dismiss info modal on Space", func() {
			modal := feedback.NewInfoModal("Info", "Some information")

			dismissed := modal.Update(tea.KeyMsg{Type: tea.KeySpace})

			Expect(dismissed).To(BeTrue())
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should display and dismiss info modal on Esc", func() {
			modal := feedback.NewInfoModal("Info", "Some information")

			dismissed := modal.Update(tea.KeyMsg{Type: tea.KeyEscape})

			Expect(dismissed).To(BeTrue())
			Expect(modal.IsVisible()).To(BeFalse())
		})
	})

	Describe("Warning Display Workflow", func() {
		It("should display warning modal with warning styling", func() {
			// User triggered an action that needs a warning.
			modal := feedback.NewWarningInfoModal(
				"Configuration Missing",
				"No profile configuration found. Some features may not work correctly.",
			)

			Expect(modal.IsVisible()).To(BeTrue())
			Expect(modal.GetTitle()).To(Equal("Configuration Missing"))
			Expect(modal.GetVariant()).To(Equal(feedback.InfoModalWarning))

			// View should show the warning message.
			view := modal.View()
			Expect(view).To(ContainSubstring("Configuration Missing"))
			Expect(view).To(ContainSubstring("profile configuration"))
			Expect(view).To(ContainSubstring("Close")) // Dismiss instruction

			// User acknowledges the warning.
			dismissed := modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(dismissed).To(BeTrue())
		})
	})

	Describe("Success Display Workflow", func() {
		It("should display success modal with success styling", func() {
			// User completed an action successfully.
			modal := feedback.NewSuccessInfoModal(
				"Export Complete",
				"Your CV has been exported to cv_output.pdf",
			)

			Expect(modal.IsVisible()).To(BeTrue())
			Expect(modal.GetTitle()).To(Equal("Export Complete"))
			Expect(modal.GetVariant()).To(Equal(feedback.InfoModalSuccess))

			// View should show success message.
			view := modal.View()
			Expect(view).To(ContainSubstring("Export Complete"))
			Expect(view).To(ContainSubstring("cv_output.pdf"))

			// User acknowledges.
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(modal.IsVisible()).To(BeFalse())
		})
	})

	Describe("Modal Visibility Control", func() {
		It("should support show/hide workflow", func() {
			modal := feedback.NewInfoModal("Test", "Test message")

			// Initially visible.
			Expect(modal.IsVisible()).To(BeTrue())

			// Hide the modal.
			modal.Hide()
			Expect(modal.IsVisible()).To(BeFalse())
			Expect(modal.View()).To(BeEmpty())

			// Show again.
			modal.Show()
			Expect(modal.IsVisible()).To(BeTrue())
			Expect(modal.View()).NotTo(BeEmpty())
		})

		It("should not respond to keys when hidden", func() {
			modal := feedback.NewInfoModal("Test", "Test message")
			modal.Hide()

			// Pressing Enter when hidden should not dismiss (already hidden).
			dismissed := modal.Update(tea.KeyMsg{Type: tea.KeyEnter})

			Expect(dismissed).To(BeFalse())
			Expect(modal.IsVisible()).To(BeFalse())
		})
	})

	Describe("Window Resize Handling", func() {
		It("should handle window resize without dismissing", func() {
			modal := feedback.NewInfoModal("Test", "Test message")

			// Simulate window resize.
			dismissed := modal.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

			Expect(dismissed).To(BeFalse())
			Expect(modal.IsVisible()).To(BeTrue())
		})
	})

	Describe("Multiline Message Support", func() {
		It("should display multiline messages correctly", func() {
			modal := feedback.NewInfoModal(
				"Keyboard Shortcuts",
				"Available shortcuts:\n\n• Enter - Select\n• Esc - Cancel\n• ? - Help",
			)
			modal.SetDimensions(100, 30)

			view := modal.View()

			Expect(view).To(ContainSubstring("Keyboard Shortcuts"))
			Expect(view).To(ContainSubstring("Select"))
			Expect(view).To(ContainSubstring("Cancel"))
			Expect(view).To(ContainSubstring("Help"))
		})
	})

	Describe("Other Key Handling", func() {
		It("should ignore unrecognized keys", func() {
			modal := feedback.NewInfoModal("Test", "Test message")

			// Press random key.
			dismissed := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})

			Expect(dismissed).To(BeFalse())
			Expect(modal.IsVisible()).To(BeTrue())
		})
	})

	Describe("Init Lifecycle", func() {
		It("should return nil command from Init", func() {
			modal := feedback.NewInfoModal("Test", "Test message")

			cmd := modal.Init()

			Expect(cmd).To(BeNil())
		})
	})
})
