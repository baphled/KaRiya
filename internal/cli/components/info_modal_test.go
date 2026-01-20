package components_test

import (
	"github.com/baphled/kariya/internal/cli/components"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("InfoModal", func() {
	Describe("NewInfoModal", func() {
		It("should create a visible modal with info variant", func() {
			modal := components.NewInfoModal("Test Title", "Test message")

			Expect(modal.IsVisible()).To(BeTrue())
			Expect(modal.GetTitle()).To(Equal("Test Title"))
			Expect(modal.GetMessage()).To(Equal("Test message"))
			Expect(modal.GetVariant()).To(Equal(components.InfoModalInfo))
		})
	})

	Describe("NewWarningInfoModal", func() {
		It("should create a visible modal with warning variant", func() {
			modal := components.NewWarningInfoModal("Warning Title", "Warning message")

			Expect(modal.IsVisible()).To(BeTrue())
			Expect(modal.GetTitle()).To(Equal("Warning Title"))
			Expect(modal.GetMessage()).To(Equal("Warning message"))
			Expect(modal.GetVariant()).To(Equal(components.InfoModalWarning))
		})
	})

	Describe("NewSuccessInfoModal", func() {
		It("should create a visible modal with success variant", func() {
			modal := components.NewSuccessInfoModal("Success Title", "Success message")

			Expect(modal.IsVisible()).To(BeTrue())
			Expect(modal.GetTitle()).To(Equal("Success Title"))
			Expect(modal.GetMessage()).To(Equal("Success message"))
			Expect(modal.GetVariant()).To(Equal(components.InfoModalSuccess))
		})
	})

	Describe("Update", func() {
		var modal *components.InfoModal

		BeforeEach(func() {
			modal = components.NewInfoModal("Test", "Test message")
		})

		Context("when modal is visible", func() {
			It("should dismiss on Enter key", func() {
				dismissed := modal.Update(tea.KeyMsg{Type: tea.KeyEnter})

				Expect(dismissed).To(BeTrue())
				Expect(modal.IsVisible()).To(BeFalse())
			})

			It("should dismiss on Space key", func() {
				dismissed := modal.Update(tea.KeyMsg{Type: tea.KeySpace})

				Expect(dismissed).To(BeTrue())
				Expect(modal.IsVisible()).To(BeFalse())
			})

			It("should dismiss on Esc key", func() {
				dismissed := modal.Update(tea.KeyMsg{Type: tea.KeyEscape})

				Expect(dismissed).To(BeTrue())
				Expect(modal.IsVisible()).To(BeFalse())
			})

			It("should NOT dismiss on other keys", func() {
				dismissed := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

				Expect(dismissed).To(BeFalse())
				Expect(modal.IsVisible()).To(BeTrue())
			})

			It("should handle WindowSizeMsg without dismissing", func() {
				dismissed := modal.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

				Expect(dismissed).To(BeFalse())
				Expect(modal.IsVisible()).To(BeTrue())
			})
		})

		Context("when modal is hidden", func() {
			BeforeEach(func() {
				modal.Hide()
			})

			It("should not dismiss on Enter when already hidden", func() {
				dismissed := modal.Update(tea.KeyMsg{Type: tea.KeyEnter})

				Expect(dismissed).To(BeFalse())
				Expect(modal.IsVisible()).To(BeFalse())
			})
		})
	})

	Describe("View", func() {
		It("should render title and message", func() {
			modal := components.NewInfoModal("My Title", "My message content")
			modal.SetDimensions(100, 30)

			view := modal.View()

			Expect(view).To(ContainSubstring("My Title"))
			Expect(view).To(ContainSubstring("My message content"))
		})

		It("should render dismiss instructions", func() {
			modal := components.NewInfoModal("Title", "Message")
			modal.SetDimensions(100, 30)

			view := modal.View()

			Expect(view).To(ContainSubstring("Close"))
		})

		It("should return empty string when hidden", func() {
			modal := components.NewInfoModal("Title", "Message")
			modal.Hide()

			view := modal.View()

			Expect(view).To(BeEmpty())
		})

		It("should handle multiline messages", func() {
			modal := components.NewInfoModal("Title", "Line 1\n\nLine 2\n\nLine 3")
			modal.SetDimensions(100, 30)

			view := modal.View()

			Expect(view).To(ContainSubstring("Line 1"))
			Expect(view).To(ContainSubstring("Line 2"))
			Expect(view).To(ContainSubstring("Line 3"))
		})
	})

	Describe("Visibility Controls", func() {
		var modal *components.InfoModal

		BeforeEach(func() {
			modal = components.NewInfoModal("Test", "Test message")
		})

		It("should start visible", func() {
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("should hide when Hide is called", func() {
			modal.Hide()
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should show when Show is called after hide", func() {
			modal.Hide()
			modal.Show()
			Expect(modal.IsVisible()).To(BeTrue())
		})
	})

	Describe("Init", func() {
		It("should return nil command", func() {
			modal := components.NewInfoModal("Test", "Test message")

			cmd := modal.Init()

			Expect(cmd).To(BeNil())
		})
	})
})
