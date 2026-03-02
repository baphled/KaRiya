package feedback_test

import (
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type testHelpKeyMap struct {
	Quit key.Binding
	Help key.Binding
}

func (k testHelpKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Help}
}

func (k testHelpKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit, k.Help},
	}
}

func newTestKeyMap() testHelpKeyMap {
	return testHelpKeyMap{
		Quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
	}
}

var _ = Describe("HelpModal", func() {
	var (
		modal  *feedback.HelpModal
		theme  themes.Theme
		keyMap testHelpKeyMap
	)

	BeforeEach(func() {
		theme = themes.NewDefaultTheme()
		keyMap = newTestKeyMap()
		modal = feedback.NewHelpModal(keyMap)
	})

	Describe("NewHelpModal", func() {
		It("should create a modal", func() {
			Expect(modal).NotTo(BeNil())
		})

		It("should not be visible initially", func() {
			Expect(modal.IsVisible()).To(BeFalse())
		})
	})

	Describe("DefaultHelpModalKeyMap", func() {
		It("should return a keymap with toggle binding", func() {
			km := feedback.DefaultHelpModalKeyMap()
			Expect(km.Toggle.Keys()).NotTo(BeEmpty())
		})

		It("should return a keymap with close binding", func() {
			km := feedback.DefaultHelpModalKeyMap()
			Expect(km.Close.Keys()).NotTo(BeEmpty())
		})

		It("should return a keymap with full help binding", func() {
			km := feedback.DefaultHelpModalKeyMap()
			Expect(km.FullHelp.Keys()).NotTo(BeEmpty())
		})
	})

	Describe("Visibility", func() {
		It("should not be visible initially", func() {
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should be visible after Show", func() {
			modal.Show()
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("should not be visible after Hide", func() {
			modal.Show()
			modal.Hide()
			Expect(modal.IsVisible()).To(BeFalse())
		})

		Context("Toggle", func() {
			It("should make modal visible when hidden", func() {
				modal.Toggle()
				Expect(modal.IsVisible()).To(BeTrue())
			})

			It("should make modal hidden when visible", func() {
				modal.Show()
				modal.Toggle()
				Expect(modal.IsVisible()).To(BeFalse())
			})

			It("should cycle visibility on repeated toggles", func() {
				modal.Toggle()
				Expect(modal.IsVisible()).To(BeTrue())
				modal.Toggle()
				Expect(modal.IsVisible()).To(BeFalse())
			})
		})
	})

	Describe("SetSize", func() {
		It("should update dimensions without error", func() {
			modal.SetSize(120, 40)
			modal.Show()
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle small dimensions", func() {
			modal.SetSize(30, 10)
			modal.Show()
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("WithTheme", func() {
		It("should accept a custom theme", func() {
			result := modal.WithTheme(theme)
			Expect(result).NotTo(BeNil())
		})

		It("should handle nil theme gracefully", func() {
			result := modal.WithTheme(nil)
			result.Show()
			view := result.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should support method chaining", func() {
			result := feedback.NewHelpModal(keyMap).WithTheme(theme)
			Expect(result).NotTo(BeNil())
		})
	})

	Describe("View", func() {
		It("should return empty string when not visible", func() {
			view := modal.View()
			Expect(view).To(BeEmpty())
		})

		It("should return non-empty string when visible", func() {
			modal.Show()
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should contain keyboard shortcuts title", func() {
			modal.Show()
			view := modal.View()
			Expect(view).To(ContainSubstring("Keyboard Shortcuts"))
		})

		It("should contain help content from keymap", func() {
			modal.Show()
			view := modal.View()
			Expect(view).To(ContainSubstring("quit"))
		})

		It("should contain footer hints", func() {
			modal.Show()
			view := modal.View()
			Expect(view).To(ContainSubstring("close"))
		})
	})

	Describe("ToggleFullHelp", func() {
		It("should toggle between short and full help", func() {
			modal.Show()
			shortView := modal.View()

			modal.ToggleFullHelp()
			fullView := modal.View()

			Expect(shortView).NotTo(BeEmpty())
			Expect(fullView).NotTo(BeEmpty())
		})
	})

	Describe("SetKeyMap", func() {
		It("should update the displayed keymap", func() {
			newKeyMap := testHelpKeyMap{
				Quit: key.NewBinding(
					key.WithKeys("x"),
					key.WithHelp("x", "exit"),
				),
				Help: key.NewBinding(
					key.WithKeys("h"),
					key.WithHelp("h", "help"),
				),
			}
			modal.SetKeyMap(newKeyMap)
			modal.Show()
			view := modal.View()
			Expect(view).To(ContainSubstring("exit"))
		})
	})

	Describe("ShortHelp", func() {
		It("should return help text from keymap", func() {
			helpText := modal.ShortHelp()
			Expect(helpText).To(ContainSubstring("quit"))
		})
	})

	Describe("RenderOverlay", func() {
		It("should return base content when not visible", func() {
			base := "Hello World\nSecond line"
			result := modal.RenderOverlay(base)
			Expect(result).To(Equal(base))
		})

		It("should render overlay when visible", func() {
			modal.Show()
			modal.SetSize(80, 24)
			base := "Line 1\nLine 2\nLine 3\nLine 4\nLine 5\nLine 6\nLine 7\nLine 8"
			result := modal.RenderOverlay(base)
			Expect(result).NotTo(Equal(base))
		})
	})

	Describe("Integration - Full Workflow", func() {
		It("should handle help display workflow", func() {
			Expect(modal.IsVisible()).To(BeFalse())
			Expect(modal.View()).To(BeEmpty())

			modal.Show()
			Expect(modal.IsVisible()).To(BeTrue())

			view := modal.View()
			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("Keyboard Shortcuts"))

			modal.Hide()
			Expect(modal.IsVisible()).To(BeFalse())
			Expect(modal.View()).To(BeEmpty())
		})
	})

	Describe("Update", func() {
		Context("when not visible", func() {
			It("should open on toggle key press", func() {
				consumed, cmd := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
				Expect(consumed).To(BeTrue())
				Expect(cmd).To(BeNil())
				Expect(modal.IsVisible()).To(BeTrue())
			})

			It("should not consume other keys", func() {
				consumed, _ := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
				Expect(consumed).To(BeFalse())
				Expect(modal.IsVisible()).To(BeFalse())
			})

			It("should not consume non-key messages", func() {
				consumed, _ := modal.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
				Expect(consumed).To(BeFalse())
			})
		})

		Context("when visible", func() {
			BeforeEach(func() {
				modal.Show()
			})

			It("should close on esc key", func() {
				consumed, cmd := modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
				Expect(consumed).To(BeTrue())
				Expect(cmd).To(BeNil())
				Expect(modal.IsVisible()).To(BeFalse())
			})

			It("should toggle full help on f key", func() {
				consumed, cmd := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
				Expect(consumed).To(BeTrue())
				Expect(cmd).To(BeNil())
				Expect(modal.IsVisible()).To(BeTrue())
			})

			It("should consume unhandled keys", func() {
				consumed, _ := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'z'}})
				Expect(consumed).To(BeTrue())
			})
		})
	})
})
