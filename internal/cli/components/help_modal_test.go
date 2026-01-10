package components

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/navigation"
)

var _ = Describe("HelpModal", func() {
	var modal *HelpModal
	var keyMap navigation.GlobalKeyMap

	BeforeEach(func() {
		keyMap = navigation.DefaultGlobalKeyMap()
		modal = NewHelpModal(keyMap)
	})

	Describe("NewHelpModal", func() {
		It("should create a new help modal", func() {
			Expect(modal).NotTo(BeNil())
		})

		It("should start hidden", func() {
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should have default size", func() {
			Expect(modal.width).To(Equal(80))
			Expect(modal.height).To(Equal(24))
		})
	})

	Describe("Show/Hide/Toggle", func() {
		It("should show the modal", func() {
			modal.Show()
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("should hide the modal", func() {
			modal.Show()
			modal.Hide()
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should toggle visibility", func() {
			Expect(modal.IsVisible()).To(BeFalse())
			modal.Toggle()
			Expect(modal.IsVisible()).To(BeTrue())
			modal.Toggle()
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should reset showingAll when hidden", func() {
			modal.Show()
			modal.ToggleFullHelp()
			Expect(modal.showingAll).To(BeTrue())
			modal.Hide()
			Expect(modal.showingAll).To(BeFalse())
		})
	})

	Describe("SetKeyMap", func() {
		It("should update the keymap", func() {
			listKeyMap := navigation.DefaultListKeyMap()
			modal.SetKeyMap(listKeyMap)
			// Verify by checking view contains list keys
			modal.Show()
			view := modal.View()
			Expect(view).To(ContainSubstring("up/k"))
		})
	})

	Describe("SetSize", func() {
		It("should update the size", func() {
			modal.SetSize(100, 50)
			Expect(modal.width).To(Equal(100))
			Expect(modal.height).To(Equal(50))
		})
	})

	Describe("ToggleFullHelp", func() {
		It("should toggle between short and full help", func() {
			Expect(modal.showingAll).To(BeFalse())
			modal.ToggleFullHelp()
			Expect(modal.showingAll).To(BeTrue())
			modal.ToggleFullHelp()
			Expect(modal.showingAll).To(BeFalse())
		})
	})

	Describe("Update", func() {
		Context("when modal is hidden", func() {
			It("should open on ? key", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")}
				consumed, _ := modal.Update(msg)
				Expect(consumed).To(BeTrue())
				Expect(modal.IsVisible()).To(BeTrue())
			})

			It("should not consume other keys", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")}
				consumed, _ := modal.Update(msg)
				Expect(consumed).To(BeFalse())
			})
		})

		Context("when modal is visible", func() {
			BeforeEach(func() {
				modal.Show()
			})

			It("should close on ? key", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")}
				consumed, _ := modal.Update(msg)
				Expect(consumed).To(BeTrue())
				Expect(modal.IsVisible()).To(BeFalse())
			})

			It("should close on esc key", func() {
				msg := tea.KeyMsg{Type: tea.KeyEsc}
				consumed, _ := modal.Update(msg)
				Expect(consumed).To(BeTrue())
				Expect(modal.IsVisible()).To(BeFalse())
			})

			It("should toggle full help on f key", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("f")}
				consumed, _ := modal.Update(msg)
				Expect(consumed).To(BeTrue())
				Expect(modal.showingAll).To(BeTrue())
			})

			It("should consume all events", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")}
				consumed, _ := modal.Update(msg)
				Expect(consumed).To(BeTrue())
			})
		})
	})

	Describe("View", func() {
		It("should return empty string when hidden", func() {
			Expect(modal.View()).To(BeEmpty())
		})

		It("should return content when visible", func() {
			modal.Show()
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("Keyboard Shortcuts"))
		})

		It("should include key bindings in short help view", func() {
			modal.Show()
			view := modal.View()
			Expect(view).To(ContainSubstring("q/ctrl+c"))
			Expect(view).To(ContainSubstring("quit"))
		})

		It("should include footer with toggle hint", func() {
			modal.Show()
			view := modal.View()
			Expect(view).To(ContainSubstring("f for full help"))
		})

		It("should show different footer when in full help mode", func() {
			modal.Show()
			modal.ToggleFullHelp()
			view := modal.View()
			Expect(view).To(ContainSubstring("f for short help"))
		})
	})

	Describe("ShortHelp", func() {
		It("should return short help string", func() {
			help := modal.ShortHelp()
			Expect(help).To(ContainSubstring("q/ctrl+c"))
		})

		It("should return empty string with nil keymap", func() {
			modal.keyMap = nil
			help := modal.ShortHelp()
			Expect(help).To(BeEmpty())
		})
	})

	Describe("RenderOverlay", func() {
		It("should return base content when hidden", func() {
			base := "Line 1\nLine 2\nLine 3"
			result := modal.RenderOverlay(base)
			Expect(result).To(Equal(base))
		})

		It("should overlay modal on base content when visible", func() {
			modal.Show()
			base := "Line 1\nLine 2\nLine 3\nLine 4\nLine 5"
			result := modal.RenderOverlay(base)
			Expect(result).To(ContainSubstring("Keyboard Shortcuts"))
		})
	})

	Describe("HelpModalKeyMap", func() {
		var helpKeys HelpModalKeyMap

		BeforeEach(func() {
			helpKeys = DefaultHelpModalKeyMap()
		})

		It("should have toggle binding for ?", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")}
			Expect(key.Matches(msg, helpKeys.Toggle)).To(BeTrue())
		})

		It("should have close binding for esc", func() {
			msg := tea.KeyMsg{Type: tea.KeyEsc}
			Expect(key.Matches(msg, helpKeys.Close)).To(BeTrue())
		})

		It("should have full help binding for f", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("f")}
			Expect(key.Matches(msg, helpKeys.FullHelp)).To(BeTrue())
		})
	})
})
