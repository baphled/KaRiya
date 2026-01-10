package navigation

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Helper to create a tea.KeyMsg for testing
func keyMsg(k string) tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(k)}
}

// Helper to create special key messages
func specialKeyMsg(keyType tea.KeyType) tea.KeyMsg {
	return tea.KeyMsg{Type: keyType}
}

var _ = Describe("KeyMaps", func() {
	Describe("GlobalKeyMap", func() {
		var keyMap GlobalKeyMap

		BeforeEach(func() {
			keyMap = DefaultGlobalKeyMap()
		})

		Describe("Quit binding", func() {
			It("should match 'q' key", func() {
				Expect(key.Matches(keyMsg("q"), keyMap.Quit)).To(BeTrue())
			})

			It("should match 'ctrl+c' key", func() {
				Expect(key.Matches(specialKeyMsg(tea.KeyCtrlC), keyMap.Quit)).To(BeTrue())
			})

			It("should not match unrelated keys", func() {
				Expect(key.Matches(keyMsg("x"), keyMap.Quit)).To(BeFalse())
			})

			It("should have help text", func() {
				help := keyMap.Quit.Help()
				Expect(help.Key).To(Equal("q/ctrl+c"))
				Expect(help.Desc).To(Equal("quit"))
			})
		})

		Describe("Help binding", func() {
			It("should match '?' key", func() {
				Expect(key.Matches(keyMsg("?"), keyMap.Help)).To(BeTrue())
			})

			It("should not match unrelated keys", func() {
				Expect(key.Matches(keyMsg("h"), keyMap.Help)).To(BeFalse())
			})

			It("should have help text", func() {
				help := keyMap.Help.Help()
				Expect(help.Key).To(Equal("?"))
				Expect(help.Desc).To(Equal("help"))
			})
		})

		Describe("Back binding", func() {
			It("should match 'esc' key", func() {
				Expect(key.Matches(specialKeyMsg(tea.KeyEsc), keyMap.Back)).To(BeTrue())
			})

			It("should not match unrelated keys", func() {
				Expect(key.Matches(keyMsg("b"), keyMap.Back)).To(BeFalse())
			})

			It("should have help text", func() {
				help := keyMap.Back.Help()
				Expect(help.Key).To(Equal("esc"))
				Expect(help.Desc).To(Equal("back"))
			})
		})

		Describe("Home binding", func() {
			It("should match 'h' key when in global context", func() {
				Expect(key.Matches(keyMsg("h"), keyMap.Home)).To(BeTrue())
			})

			It("should have help text", func() {
				help := keyMap.Home.Help()
				Expect(help.Key).To(Equal("h"))
				Expect(help.Desc).To(Equal("home"))
			})
		})

		Describe("help.KeyMap interface", func() {
			It("should implement ShortHelp", func() {
				shortHelp := keyMap.ShortHelp()
				Expect(shortHelp).To(HaveLen(4))
				// Should include quit, help, back, home
				keys := make([]string, len(shortHelp))
				for i, binding := range shortHelp {
					keys[i] = binding.Help().Key
				}
				Expect(keys).To(ContainElements("q/ctrl+c", "?", "esc", "h"))
			})

			It("should implement FullHelp", func() {
				fullHelp := keyMap.FullHelp()
				Expect(fullHelp).To(HaveLen(1)) // One column of global keys
				Expect(fullHelp[0]).To(HaveLen(4))
			})
		})
	})

	Describe("ListKeyMap", func() {
		var keyMap ListKeyMap

		BeforeEach(func() {
			keyMap = DefaultListKeyMap()
		})

		Describe("Up binding", func() {
			It("should match 'up' arrow key", func() {
				Expect(key.Matches(specialKeyMsg(tea.KeyUp), keyMap.Up)).To(BeTrue())
			})

			It("should match 'k' vim key", func() {
				Expect(key.Matches(keyMsg("k"), keyMap.Up)).To(BeTrue())
			})

			It("should have help text", func() {
				help := keyMap.Up.Help()
				Expect(help.Key).To(Equal("up/k"))
				Expect(help.Desc).To(Equal("up"))
			})
		})

		Describe("Down binding", func() {
			It("should match 'down' arrow key", func() {
				Expect(key.Matches(specialKeyMsg(tea.KeyDown), keyMap.Down)).To(BeTrue())
			})

			It("should match 'j' vim key", func() {
				Expect(key.Matches(keyMsg("j"), keyMap.Down)).To(BeTrue())
			})

			It("should have help text", func() {
				help := keyMap.Down.Help()
				Expect(help.Key).To(Equal("down/j"))
				Expect(help.Desc).To(Equal("down"))
			})
		})

		Describe("Select binding", func() {
			It("should match 'enter' key", func() {
				Expect(key.Matches(specialKeyMsg(tea.KeyEnter), keyMap.Select)).To(BeTrue())
			})

			It("should have help text", func() {
				help := keyMap.Select.Help()
				Expect(help.Key).To(Equal("enter"))
				Expect(help.Desc).To(Equal("select"))
			})
		})

		Describe("Delete binding", func() {
			It("should match 'd' key", func() {
				Expect(key.Matches(keyMsg("d"), keyMap.Delete)).To(BeTrue())
			})

			It("should have help text", func() {
				help := keyMap.Delete.Help()
				Expect(help.Key).To(Equal("d"))
				Expect(help.Desc).To(Equal("delete"))
			})
		})

		Describe("Edit binding", func() {
			It("should match 'e' key", func() {
				Expect(key.Matches(keyMsg("e"), keyMap.Edit)).To(BeTrue())
			})

			It("should have help text", func() {
				help := keyMap.Edit.Help()
				Expect(help.Key).To(Equal("e"))
				Expect(help.Desc).To(Equal("edit"))
			})
		})

		Describe("Filter binding", func() {
			It("should match 'f' key", func() {
				Expect(key.Matches(keyMsg("f"), keyMap.Filter)).To(BeTrue())
			})

			It("should have help text", func() {
				help := keyMap.Filter.Help()
				Expect(help.Key).To(Equal("f"))
				Expect(help.Desc).To(Equal("filter"))
			})
		})

		Describe("Search binding", func() {
			It("should match '/' key", func() {
				Expect(key.Matches(keyMsg("/"), keyMap.Search)).To(BeTrue())
			})

			It("should have help text", func() {
				help := keyMap.Search.Help()
				Expect(help.Key).To(Equal("/"))
				Expect(help.Desc).To(Equal("search"))
			})
		})

		Describe("PageUp binding", func() {
			It("should match 'pgup' key", func() {
				Expect(key.Matches(specialKeyMsg(tea.KeyPgUp), keyMap.PageUp)).To(BeTrue())
			})

			It("should match 'ctrl+u' key", func() {
				Expect(key.Matches(specialKeyMsg(tea.KeyCtrlU), keyMap.PageUp)).To(BeTrue())
			})
		})

		Describe("PageDown binding", func() {
			It("should match 'pgdown' key", func() {
				Expect(key.Matches(specialKeyMsg(tea.KeyPgDown), keyMap.PageDown)).To(BeTrue())
			})

			It("should match 'ctrl+d' key", func() {
				Expect(key.Matches(specialKeyMsg(tea.KeyCtrlD), keyMap.PageDown)).To(BeTrue())
			})
		})

		Describe("GoToStart binding", func() {
			It("should match 'home' key", func() {
				Expect(key.Matches(specialKeyMsg(tea.KeyHome), keyMap.GoToStart)).To(BeTrue())
			})

			It("should match 'g' key", func() {
				Expect(key.Matches(keyMsg("g"), keyMap.GoToStart)).To(BeTrue())
			})
		})

		Describe("GoToEnd binding", func() {
			It("should match 'end' key", func() {
				Expect(key.Matches(specialKeyMsg(tea.KeyEnd), keyMap.GoToEnd)).To(BeTrue())
			})

			It("should match 'G' key", func() {
				Expect(key.Matches(keyMsg("G"), keyMap.GoToEnd)).To(BeTrue())
			})
		})

		Describe("help.KeyMap interface", func() {
			It("should implement ShortHelp", func() {
				shortHelp := keyMap.ShortHelp()
				Expect(len(shortHelp)).To(BeNumerically(">=", 4))
			})

			It("should implement FullHelp", func() {
				fullHelp := keyMap.FullHelp()
				Expect(len(fullHelp)).To(BeNumerically(">=", 1))
			})
		})
	})

	Describe("FormKeyMap", func() {
		var keyMap FormKeyMap

		BeforeEach(func() {
			keyMap = DefaultFormKeyMap()
		})

		Describe("NextField binding", func() {
			It("should match 'tab' key", func() {
				Expect(key.Matches(specialKeyMsg(tea.KeyTab), keyMap.NextField)).To(BeTrue())
			})

			It("should have help text", func() {
				help := keyMap.NextField.Help()
				Expect(help.Key).To(Equal("tab"))
				Expect(help.Desc).To(Equal("next field"))
			})
		})

		Describe("PrevField binding", func() {
			It("should match 'shift+tab' key", func() {
				Expect(key.Matches(specialKeyMsg(tea.KeyShiftTab), keyMap.PrevField)).To(BeTrue())
			})

			It("should have help text", func() {
				help := keyMap.PrevField.Help()
				Expect(help.Key).To(Equal("shift+tab"))
				Expect(help.Desc).To(Equal("prev field"))
			})
		})

		Describe("Submit binding", func() {
			It("should match 'enter' key", func() {
				Expect(key.Matches(specialKeyMsg(tea.KeyEnter), keyMap.Submit)).To(BeTrue())
			})

			It("should have help text", func() {
				help := keyMap.Submit.Help()
				Expect(help.Key).To(Equal("enter"))
				Expect(help.Desc).To(Equal("submit"))
			})
		})

		Describe("Cancel binding", func() {
			It("should match 'esc' key", func() {
				Expect(key.Matches(specialKeyMsg(tea.KeyEsc), keyMap.Cancel)).To(BeTrue())
			})

			It("should have help text", func() {
				help := keyMap.Cancel.Help()
				Expect(help.Key).To(Equal("esc"))
				Expect(help.Desc).To(Equal("cancel"))
			})
		})

		Describe("Toggle binding", func() {
			It("should match 'space' key", func() {
				Expect(key.Matches(specialKeyMsg(tea.KeySpace), keyMap.Toggle)).To(BeTrue())
			})

			It("should have help text", func() {
				help := keyMap.Toggle.Help()
				Expect(help.Key).To(Equal("space"))
				Expect(help.Desc).To(Equal("toggle"))
			})
		})

		Describe("ToggleOptional binding", func() {
			It("should match 'ctrl+o' key", func() {
				Expect(key.Matches(specialKeyMsg(tea.KeyCtrlO), keyMap.ToggleOptional)).To(BeTrue())
			})

			It("should have help text", func() {
				help := keyMap.ToggleOptional.Help()
				Expect(help.Key).To(Equal("ctrl+o"))
				Expect(help.Desc).To(Equal("toggle optional"))
			})
		})

		Describe("help.KeyMap interface", func() {
			It("should implement ShortHelp", func() {
				shortHelp := keyMap.ShortHelp()
				Expect(len(shortHelp)).To(BeNumerically(">=", 4))
			})

			It("should implement FullHelp", func() {
				fullHelp := keyMap.FullHelp()
				Expect(len(fullHelp)).To(BeNumerically(">=", 1))
			})
		})
	})

	Describe("CombinedKeyMap", func() {
		var combined CombinedKeyMap

		BeforeEach(func() {
			combined = NewCombinedKeyMap(
				DefaultGlobalKeyMap(),
				DefaultListKeyMap(),
			)
		})

		It("should include global bindings", func() {
			Expect(key.Matches(keyMsg("q"), combined.Global.Quit)).To(BeTrue())
		})

		It("should include list bindings", func() {
			Expect(key.Matches(keyMsg("j"), combined.List.Down)).To(BeTrue())
		})

		Describe("help.KeyMap interface", func() {
			It("should implement ShortHelp combining both keymaps", func() {
				shortHelp := combined.ShortHelp()
				// Should include keys from both global and list
				Expect(len(shortHelp)).To(BeNumerically(">=", 6))
			})

			It("should implement FullHelp with multiple columns", func() {
				fullHelp := combined.FullHelp()
				// Should have at least 2 columns (global + list)
				Expect(len(fullHelp)).To(BeNumerically(">=", 2))
			})
		})
	})

	Describe("Enabled/Disabled bindings", func() {
		It("should allow disabling a binding", func() {
			keyMap := DefaultGlobalKeyMap()
			keyMap.Quit.SetEnabled(false)
			Expect(keyMap.Quit.Enabled()).To(BeFalse())
		})

		It("should allow re-enabling a binding", func() {
			keyMap := DefaultGlobalKeyMap()
			keyMap.Quit.SetEnabled(false)
			keyMap.Quit.SetEnabled(true)
			Expect(keyMap.Quit.Enabled()).To(BeTrue())
		})
	})
})
