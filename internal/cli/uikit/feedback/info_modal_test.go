package feedback_test

import (
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("InfoModal", func() {
	var (
		modal *feedback.InfoModal
		theme themes.Theme
	)

	BeforeEach(func() {
		theme = themes.NewDefaultTheme()
		modal = feedback.NewInfoModal("Test Title", "Test message content")
		modal.SetDimensions(100, 24)
	})

	Describe("NewInfoModal", func() {
		It("should create a modal with the given title", func() {
			Expect(modal.GetTitle()).To(Equal("Test Title"))
		})

		It("should create a modal with the given message", func() {
			Expect(modal.GetMessage()).To(Equal("Test message content"))
		})

		It("should create a modal with info variant", func() {
			Expect(modal.GetVariant()).To(Equal(feedback.InfoModalInfo))
		})

		It("should be visible by default", func() {
			Expect(modal.IsVisible()).To(BeTrue())
		})
	})

	Describe("NewWarningInfoModal", func() {
		BeforeEach(func() {
			modal = feedback.NewWarningInfoModal("Warning Title", "Warning message")
			modal.SetDimensions(100, 24)
		})

		It("should create a modal with warning variant", func() {
			Expect(modal.GetVariant()).To(Equal(feedback.InfoModalWarning))
		})

		It("should store the title", func() {
			Expect(modal.GetTitle()).To(Equal("Warning Title"))
		})

		It("should store the message", func() {
			Expect(modal.GetMessage()).To(Equal("Warning message"))
		})

		It("should be visible by default", func() {
			Expect(modal.IsVisible()).To(BeTrue())
		})
	})

	Describe("NewSuccessInfoModal", func() {
		BeforeEach(func() {
			modal = feedback.NewSuccessInfoModal("Success Title", "Success message")
			modal.SetDimensions(100, 24)
		})

		It("should create a modal with success variant", func() {
			Expect(modal.GetVariant()).To(Equal(feedback.InfoModalSuccess))
		})

		It("should store the title", func() {
			Expect(modal.GetTitle()).To(Equal("Success Title"))
		})

		It("should store the message", func() {
			Expect(modal.GetMessage()).To(Equal("Success message"))
		})
	})

	Describe("Visibility", func() {
		It("should be visible initially", func() {
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("should hide when Hide is called", func() {
			modal.Hide()
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should show when Show is called after hiding", func() {
			modal.Hide()
			Expect(modal.IsVisible()).To(BeFalse())
			modal.Show()
			Expect(modal.IsVisible()).To(BeTrue())
		})
	})

	Describe("WithTheme", func() {
		It("should accept a custom theme", func() {
			result := modal.WithTheme(theme)
			Expect(result).NotTo(BeNil())
			view := result.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle nil theme gracefully", func() {
			result := modal.WithTheme(nil)
			view := result.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should support method chaining", func() {
			result := feedback.NewInfoModal("Title", "Message").WithTheme(theme)
			Expect(result).NotTo(BeNil())
		})
	})

	Describe("SetDimensions", func() {
		It("should update dimensions without error", func() {
			modal.SetDimensions(80, 24)
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle small dimensions", func() {
			modal.SetDimensions(20, 10)
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Init", func() {
		It("should return nil command", func() {
			cmd := modal.Init()
			Expect(cmd).To(BeNil())
		})
	})

	Describe("View", func() {
		It("should return non-empty string when visible", func() {
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should return empty string when not visible", func() {
			modal.Hide()
			view := modal.View()
			Expect(view).To(BeEmpty())
		})

		It("should contain the title", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("Test Title"))
		})

		It("should contain the message", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("Test message content"))
		})

		It("should contain dismiss hint", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("Enter/Esc"))
		})

		Context("with warning variant", func() {
			It("should render with warning styling", func() {
				modal = feedback.NewWarningInfoModal("Warn", "Be careful")
				modal.SetDimensions(100, 24)
				view := modal.View()
				Expect(view).NotTo(BeEmpty())
				Expect(view).To(ContainSubstring("Warn"))
			})
		})

		Context("with success variant", func() {
			It("should render with success styling", func() {
				modal = feedback.NewSuccessInfoModal("Done", "All good")
				modal.SetDimensions(100, 24)
				view := modal.View()
				Expect(view).NotTo(BeEmpty())
				Expect(view).To(ContainSubstring("Done"))
			})
		})
	})

	Describe("Integration - Full Workflow", func() {
		It("should handle info display and dismiss workflow", func() {
			Expect(modal.IsVisible()).To(BeTrue())

			view := modal.View()
			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("Test Title"))

			modal.Hide()
			Expect(modal.IsVisible()).To(BeFalse())
			Expect(modal.View()).To(BeEmpty())
		})
	})

	Describe("Update", func() {
		It("should return false when not visible", func() {
			modal.Hide()
			result := modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(result).To(BeFalse())
		})

		It("should dismiss on enter key", func() {
			result := modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(result).To(BeTrue())
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should dismiss on space key", func() {
			result := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
			Expect(result).To(BeTrue())
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should dismiss on esc key", func() {
			result := modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(result).To(BeTrue())
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should handle window size message", func() {
			result := modal.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			Expect(result).To(BeFalse())
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("should ignore unhandled keys", func() {
			result := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
			Expect(result).To(BeFalse())
			Expect(modal.IsVisible()).To(BeTrue())
		})
	})
})
