package base_test

import (
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ConfirmScreen", func() {
	var screen *base.ConfirmScreen

	Describe("Construction", func() {
		It("should create a confirm screen with title and message", func() {
			screen = base.NewBaseConfirmScreen(
				[]string{"Test", "Confirm"},
				"Delete Item",
				"Are you sure you want to delete this item?",
			)

			Expect(screen).NotTo(BeNil())
			Expect(screen.GetTitle()).To(Equal("Delete Item"))
			Expect(screen.GetMessage()).To(Equal("Are you sure you want to delete this item?"))
		})

		It("should default to 'No' selection", func() {
			screen = base.NewBaseConfirmScreen(
				[]string{"Test"},
				"Confirm",
				"Message",
			)

			Expect(screen.GetSelection()).To(BeFalse())
		})

		It("should use default dimensions if not set", func() {
			screen = base.NewBaseConfirmScreen([]string{"Test"}, "Title", "Message")
			Expect(screen.Width()).To(Equal(120)) // Default width
			Expect(screen.Height()).To(Equal(40)) // Default height
		})
	})

	Describe("Terminal Info", func() {
		BeforeEach(func() {
			screen = base.NewBaseConfirmScreen([]string{"Test"}, "Title", "Message")
		})

		It("should update dimensions on SetTerminalInfo", func() {
			screen.SetTerminalInfo(100, 30)
			Expect(screen.Width()).To(Equal(100))
			Expect(screen.Height()).To(Equal(30))
		})

		It("should handle WindowSizeMsg", func() {
			msg := tea.WindowSizeMsg{Width: 80, Height: 24}
			cmd, result := screen.Update(msg)

			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
			Expect(screen.Width()).To(Equal(80))
			Expect(screen.Height()).To(Equal(24))
		})
	})

	Describe("Selection Toggle", func() {
		BeforeEach(func() {
			screen = base.NewBaseConfirmScreen([]string{"Test"}, "Title", "Message")
		})

		It("should toggle to Yes on left arrow", func() {
			msg := tea.KeyMsg{Type: tea.KeyLeft}
			_, result := screen.Update(msg)

			Expect(result).To(BeNil())
			Expect(screen.GetSelection()).To(BeTrue())
		})

		It("should toggle to Yes on right arrow", func() {
			msg := tea.KeyMsg{Type: tea.KeyRight}
			_, result := screen.Update(msg)

			Expect(result).To(BeNil())
			Expect(screen.GetSelection()).To(BeTrue())
		})

		It("should toggle to Yes on h key (vim-style)", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}}
			_, result := screen.Update(msg)

			Expect(result).To(BeNil())
			Expect(screen.GetSelection()).To(BeTrue())
		})

		It("should toggle to No on l key (vim-style)", func() {
			// First toggle to Yes
			screen.SetSelection(true)

			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}}
			_, result := screen.Update(msg)

			Expect(result).To(BeNil())
			Expect(screen.GetSelection()).To(BeFalse())
		})

		It("should toggle back and forth", func() {
			Expect(screen.GetSelection()).To(BeFalse()) // Start at No

			screen.Update(tea.KeyMsg{Type: tea.KeyLeft})
			Expect(screen.GetSelection()).To(BeTrue()) // Now Yes

			screen.Update(tea.KeyMsg{Type: tea.KeyRight})
			Expect(screen.GetSelection()).To(BeFalse()) // Back to No
		})
	})

	Describe("Keyboard Shortcuts", func() {
		BeforeEach(func() {
			screen = base.NewBaseConfirmScreen([]string{"Test"}, "Title", "Message")
		})

		It("should return CancelResult on escape key", func() {
			msg := tea.KeyMsg{Type: tea.KeyEsc}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultCancel))
		})

		It("should return NavigateResult on enter key with current selection", func() {
			// Default selection is No (false)
			msg := tea.KeyMsg{Type: tea.KeyEnter}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultNavigate))
			Expect(result.Data()).To(BeFalse())
		})

		It("should return NavigateResult with Yes when Yes is selected", func() {
			screen.SetSelection(true)

			msg := tea.KeyMsg{Type: tea.KeyEnter}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultNavigate))
			Expect(result.Data()).To(BeTrue())
		})

		It("should submit No on 'n' key", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultNavigate))
			Expect(result.Data()).To(BeFalse())
		})

		It("should submit Yes on 'y' key", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultNavigate))
			Expect(result.Data()).To(BeTrue())
		})

		It("should handle capital Y for Yes", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'Y'}}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Data()).To(BeTrue())
		})

		It("should handle capital N for No", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'N'}}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Data()).To(BeFalse())
		})
	})

	Describe("View Rendering", func() {
		BeforeEach(func() {
			screen = base.NewBaseConfirmScreen(
				[]string{"Test", "Confirm"},
				"Delete Item",
				"Are you sure you want to delete this item?",
			)
			screen.SetTerminalInfo(120, 40)
		})

		It("should render a view with StandardView structure", func() {
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should include title in view", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Delete Item"))
		})

		It("should include message in view", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Are you sure you want to delete this item?"))
		})

		It("should show Yes/No buttons", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Yes"))
			Expect(view).To(ContainSubstring("No"))
		})

		It("should highlight selected button", func() {
			// Default is No
			viewNo := screen.View()
			Expect(viewNo).To(ContainSubstring("No"))

			// Toggle to Yes
			screen.SetSelection(true)
			viewYes := screen.View()
			Expect(viewYes).To(ContainSubstring("Yes"))

			// Both views should contain both buttons
			Expect(viewNo).To(ContainSubstring("Yes"))
			Expect(viewNo).To(ContainSubstring("No"))
			Expect(viewYes).To(ContainSubstring("Yes"))
			Expect(viewYes).To(ContainSubstring("No"))
		})

		It("should update view when selection changes", func() {
			view1 := screen.View()
			Expect(view1).NotTo(BeEmpty())

			screen.Update(tea.KeyMsg{Type: tea.KeyLeft}) // Toggle to Yes
			view2 := screen.View()
			Expect(view2).NotTo(BeEmpty())

			// Both views should contain both buttons
			Expect(view1).To(ContainSubstring("Yes"))
			Expect(view1).To(ContainSubstring("No"))
			Expect(view2).To(ContainSubstring("Yes"))
			Expect(view2).To(ContainSubstring("No"))
		})
	})

	Describe("Customization", func() {
		It("should allow custom footer", func() {
			screen = base.NewBaseConfirmScreen([]string{"Test"}, "Title", "Message")
			screen.SetFooter("Custom footer text")

			view := screen.View()
			Expect(view).To(ContainSubstring("Custom footer text"))
		})

		It("should allow custom Yes text", func() {
			screen = base.NewBaseConfirmScreen([]string{"Test"}, "Title", "Message")
			screen.SetYesText("Confirm")

			view := screen.View()
			Expect(view).To(ContainSubstring("Confirm"))
		})

		It("should allow custom No text", func() {
			screen = base.NewBaseConfirmScreen([]string{"Test"}, "Title", "Message")
			screen.SetNoText("Cancel")

			view := screen.View()
			Expect(view).To(ContainSubstring("Cancel"))
		})

		It("should allow both custom Yes and No text", func() {
			screen = base.NewBaseConfirmScreen([]string{"Test"}, "Title", "Message")
			screen.SetYesText("Proceed")
			screen.SetNoText("Abort")

			view := screen.View()
			Expect(view).To(ContainSubstring("Proceed"))
			Expect(view).To(ContainSubstring("Abort"))
		})
	})

	Describe("Data Access", func() {
		BeforeEach(func() {
			screen = base.NewBaseConfirmScreen(
				[]string{"Test"},
				"Confirm Action",
				"Proceed with action?",
			)
		})

		It("should provide access to title", func() {
			Expect(screen.GetTitle()).To(Equal("Confirm Action"))
		})

		It("should provide access to message", func() {
			Expect(screen.GetMessage()).To(Equal("Proceed with action?"))
		})

		It("should allow getting selection state", func() {
			Expect(screen.GetSelection()).To(BeFalse()) // Default

			screen.SetSelection(true)
			Expect(screen.GetSelection()).To(BeTrue())
		})

		It("should allow setting selection state", func() {
			screen.SetSelection(true)
			Expect(screen.GetSelection()).To(BeTrue())

			screen.SetSelection(false)
			Expect(screen.GetSelection()).To(BeFalse())
		})
	})

	Describe("Screen Interface", func() {
		BeforeEach(func() {
			screen = base.NewBaseConfirmScreen([]string{"Test"}, "Title", "Message")
		})

		It("should implement Screen interface", func() {
			var _ screens.Screen = screen
		})

		It("should support SetTheme", func() {
			theme := "test-theme"
			screen.SetTheme(theme)
			Expect(screen.Theme()).To(Equal(theme))
		})
	})

	Describe("Edge Cases", func() {
		It("should handle empty title", func() {
			screen := base.NewBaseConfirmScreen([]string{"Test"}, "", "Message")
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle empty message", func() {
			screen := base.NewBaseConfirmScreen([]string{"Test"}, "Title", "")
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle empty breadcrumbs", func() {
			screen := base.NewBaseConfirmScreen([]string{}, "Title", "Message")
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle very small terminal dimensions", func() {
			screen := base.NewBaseConfirmScreen([]string{"Test"}, "Title", "Message")
			screen.SetTerminalInfo(20, 10)

			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle long title", func() {
			longTitle := "This is a very long title that might wrap or be truncated depending on terminal width"
			screen := base.NewBaseConfirmScreen([]string{"Test"}, longTitle, "Message")
			screen.SetTerminalInfo(40, 20)

			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle long message", func() {
			longMessage := "This is a very long message that contains multiple sentences and might need to wrap across multiple lines depending on the terminal width and the rendering implementation."
			screen := base.NewBaseConfirmScreen([]string{"Test"}, "Title", longMessage)
			screen.SetTerminalInfo(40, 20)

			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})
	})
})
