package components_test

import (
	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/themes"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeleteConfirmModal", func() {
	var (
		modal  *components.DeleteConfirmModal
		theme  themes.Theme
		entity string
		title  string
		msg    string
	)

	BeforeEach(func() {
		theme = themes.NewDefaultTheme()
		entity = "Go Programming"
		title = "Delete Skill"
		msg = "Are you sure you want to delete 'Go Programming'?"
		modal = components.NewDeleteConfirmModal(entity, title, msg)
		modal.SetTheme(theme)
		modal.SetDimensions(100, 24)
	})

	Describe("NewDeleteConfirmModal", func() {
		It("should create a visible modal", func() {
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("should not be confirmed initially", func() {
			Expect(modal.WasConfirmed()).To(BeFalse())
		})
	})

	Describe("Init", func() {
		It("should return nil command", func() {
			cmd := modal.Init()
			Expect(cmd).To(BeNil())
		})
	})

	Describe("Update - Confirmation", func() {
		It("should confirm on 'y' key", func() {
			cmd, confirmed := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
			Expect(cmd).To(BeNil())
			Expect(confirmed).To(BeTrue())
			Expect(modal.IsVisible()).To(BeFalse())
			Expect(modal.WasConfirmed()).To(BeTrue())
		})

		It("should confirm on 'Y' key", func() {
			cmd, confirmed := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'Y'}})
			Expect(cmd).To(BeNil())
			Expect(confirmed).To(BeTrue())
			Expect(modal.WasConfirmed()).To(BeTrue())
		})

		It("should confirm on Enter key", func() {
			cmd, confirmed := modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(cmd).To(BeNil())
			Expect(confirmed).To(BeTrue())
			Expect(modal.WasConfirmed()).To(BeTrue())
		})
	})

	Describe("Update - Cancellation", func() {
		It("should cancel on 'n' key", func() {
			cmd, confirmed := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			Expect(cmd).To(BeNil())
			Expect(confirmed).To(BeFalse())
			Expect(modal.IsVisible()).To(BeFalse())
			Expect(modal.WasConfirmed()).To(BeFalse())
		})

		It("should cancel on 'N' key", func() {
			cmd, confirmed := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'N'}})
			Expect(cmd).To(BeNil())
			Expect(confirmed).To(BeFalse())
			Expect(modal.WasConfirmed()).To(BeFalse())
		})

		It("should cancel on Esc key", func() {
			cmd, confirmed := modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(cmd).To(BeNil())
			Expect(confirmed).To(BeFalse())
			Expect(modal.WasConfirmed()).To(BeFalse())
		})
	})

	Describe("Update - WindowSizeMsg", func() {
		It("should update dimensions on window resize", func() {
			modal.Update(tea.WindowSizeMsg{Width: 120, Height: 30})
			// Dimensions are updated internally (verified by View not panicking)
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("View", func() {
		It("should render modal when visible", func() {
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring(title))
			Expect(view).To(ContainSubstring("Confirm"))
			Expect(view).To(ContainSubstring("Cancel"))
		})

		It("should return empty string when not visible", func() {
			modal.Hide()
			view := modal.View()
			Expect(view).To(BeEmpty())
		})

		It("should show confirmation message", func() {
			view := modal.View()
			// Message might be wrapped, so check for key words
			Expect(view).To(ContainSubstring("delete"))
		})

		It("should use KeyBadge components in footer", func() {
			view := modal.View()
			// KeyBadges render with key and hint
			Expect(view).To(ContainSubstring("y/Enter"))
			Expect(view).To(ContainSubstring("n/Esc"))
		})
	})

	Describe("Show/Hide", func() {
		It("should show modal", func() {
			modal.Hide()
			Expect(modal.IsVisible()).To(BeFalse())
			modal.Show()
			Expect(modal.IsVisible()).To(BeTrue())
			Expect(modal.WasConfirmed()).To(BeFalse()) // Reset confirmed state
		})

		It("should hide modal", func() {
			Expect(modal.IsVisible()).To(BeTrue())
			modal.Hide()
			Expect(modal.IsVisible()).To(BeFalse())
			Expect(modal.WasConfirmed()).To(BeFalse()) // Reset confirmed state
		})
	})

	Describe("SetTheme", func() {
		It("should accept custom theme", func() {
			customTheme := themes.NewDefaultTheme()
			modal.SetTheme(customTheme)
			// Theme applied (verified by View not panicking)
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Integration - Full Workflow", func() {
		It("should handle confirm workflow", func() {
			// Initial state
			Expect(modal.IsVisible()).To(BeTrue())
			Expect(modal.WasConfirmed()).To(BeFalse())

			// View should render
			view := modal.View()
			Expect(view).NotTo(BeEmpty())

			// Confirm deletion
			cmd, confirmed := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
			Expect(cmd).To(BeNil())
			Expect(confirmed).To(BeTrue())

			// Modal should be hidden and confirmed
			Expect(modal.IsVisible()).To(BeFalse())
			Expect(modal.WasConfirmed()).To(BeTrue())

			// View should be empty
			view = modal.View()
			Expect(view).To(BeEmpty())
		})

		It("should handle cancel workflow", func() {
			// Initial state
			Expect(modal.IsVisible()).To(BeTrue())

			// View should render
			view := modal.View()
			Expect(view).NotTo(BeEmpty())

			// Cancel deletion
			cmd, confirmed := modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(cmd).To(BeNil())
			Expect(confirmed).To(BeFalse())

			// Modal should be hidden and not confirmed
			Expect(modal.IsVisible()).To(BeFalse())
			Expect(modal.WasConfirmed()).To(BeFalse())
		})
	})
})
