package navigation

import (
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ListKeyHandler", func() {
	var handler *ListKeyHandler

	BeforeEach(func() {
		handler = NewListKeyHandler()
	})

	Context("Edit key (e)", func() {
		It("should handle 'e' key and return KeyEdit", func() {
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
			action := handler.HandleKey(keyMsg)

			Expect(action.IsHandled).To(BeTrue())
			Expect(action.NavigationKey).NotTo(BeNil())
			Expect(*action.NavigationKey).To(Equal(KeyEdit))
		})

		It("should recognize 'e' as a valid key string", func() {
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
			action := handler.HandleKey(keyMsg)

			Expect(action.IsHandled).To(BeTrue())
		})
	})

	Context("Delete key (d)", func() {
		It("should handle 'd' key and return KeyDelete", func() {
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}
			action := handler.HandleKey(keyMsg)

			Expect(action.IsHandled).To(BeTrue())
			Expect(action.NavigationKey).NotTo(BeNil())
			Expect(*action.NavigationKey).To(Equal(KeyDelete))
		})

		It("should recognize 'd' as a valid key string", func() {
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}
			action := handler.HandleKey(keyMsg)

			Expect(action.IsHandled).To(BeTrue())
		})
	})

	Context("Other list keys", func() {
		It("should handle up navigation with 'k'", func() {
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
			action := handler.HandleKey(keyMsg)

			Expect(action.IsHandled).To(BeTrue())
			Expect(action.NavigationKey).NotTo(BeNil())
			Expect(*action.NavigationKey).To(Equal(KeyUp))
		})

		It("should handle down navigation with 'j'", func() {
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
			action := handler.HandleKey(keyMsg)

			Expect(action.IsHandled).To(BeTrue())
			Expect(action.NavigationKey).NotTo(BeNil())
			Expect(*action.NavigationKey).To(Equal(KeyDown))
		})

		It("should handle sort with 's'", func() {
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
			action := handler.HandleKey(keyMsg)

			Expect(action.IsHandled).To(BeTrue())
			Expect(action.NavigationKey).NotTo(BeNil())
			Expect(*action.NavigationKey).To(Equal(KeySort))
		})

		It("should handle filter with 'f'", func() {
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}}
			action := handler.HandleKey(keyMsg)

			Expect(action.IsHandled).To(BeTrue())
			Expect(action.NavigationKey).NotTo(BeNil())
			Expect(*action.NavigationKey).To(Equal(KeyFilter))
		})
	})
})

var _ = Describe("DialogKeyHandler", func() {
	var handler *DialogKeyHandler

	BeforeEach(func() {
		handler = NewDialogKeyHandler()
	})

	Context("No key (n)", func() {
		It("should handle 'n' key and return dialog:no action", func() {
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
			action := handler.HandleKey(keyMsg)

			Expect(action.IsHandled).To(BeTrue())
			Expect(action.ActionType).To(Equal("dialog:no"))
		})

		It("should recognize 'n' as a valid key string", func() {
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
			action := handler.HandleKey(keyMsg)

			Expect(action.IsHandled).To(BeTrue())
		})
	})

	Context("Yes key (y)", func() {
		It("should handle 'y' key and return dialog:yes action", func() {
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}}
			action := handler.HandleKey(keyMsg)

			Expect(action.IsHandled).To(BeTrue())
			Expect(action.ActionType).To(Equal("dialog:yes"))
		})
	})

	Context("Other dialog keys", func() {
		It("should handle enter key", func() {
			keyMsg := tea.KeyMsg{Type: tea.KeyEnter}
			action := handler.HandleKey(keyMsg)

			Expect(action.IsHandled).To(BeTrue())
			Expect(action.ActionType).To(Equal("dialog:ok"))
		})

		It("should handle escape key", func() {
			keyMsg := tea.KeyMsg{Type: tea.KeyEsc}
			action := handler.HandleKey(keyMsg)

			Expect(action.IsHandled).To(BeTrue())
			Expect(action.ActionType).To(Equal("dialog:cancel"))
		})
	})
})

var _ = Describe("FormKeyHandler", func() {
	var handler *FormKeyHandler

	BeforeEach(func() {
		handler = NewFormKeyHandler()
	})

	Context("Form navigation", func() {
		It("should handle tab key for next field", func() {
			keyMsg := tea.KeyMsg{Type: tea.KeyTab}
			action := handler.HandleKey(keyMsg)

			Expect(action.IsHandled).To(BeTrue())
			Expect(action.ActionType).To(Equal("field:next"))
		})

		It("should handle shift+tab key for previous field", func() {
			keyMsg := tea.KeyMsg{Type: tea.KeyShiftTab}
			action := handler.HandleKey(keyMsg)

			Expect(action.IsHandled).To(BeTrue())
			Expect(action.ActionType).To(Equal("field:previous"))
		})

		It("should handle down key for next field", func() {
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
			action := handler.HandleKey(keyMsg)

			Expect(action.IsHandled).To(BeTrue())
			Expect(action.ActionType).To(Equal("field:next"))
		})
	})
})

var _ = Describe("ViewKeyHandler", func() {
	var handler *ViewKeyHandler

	BeforeEach(func() {
		handler = NewViewKeyHandler()
	})

	Context("View navigation", func() {
		It("should handle up key for scroll up", func() {
			keyMsg := tea.KeyMsg{Type: tea.KeyUp}
			action := handler.HandleKey(keyMsg)

			Expect(action.IsHandled).To(BeTrue())
			Expect(action.ActionType).To(Equal("view:scroll_up"))
		})

		It("should handle down key for scroll down", func() {
			keyMsg := tea.KeyMsg{Type: tea.KeyDown}
			action := handler.HandleKey(keyMsg)

			Expect(action.IsHandled).To(BeTrue())
			Expect(action.ActionType).To(Equal("view:scroll_down"))
		})
	})
})

var _ = Describe("MenuKeyHandler", func() {
	var handler *MenuKeyHandler

	BeforeEach(func() {
		handler = NewMenuKeyHandler()
	})

	Context("Menu navigation", func() {
		It("should handle up key", func() {
			keyMsg := tea.KeyMsg{Type: tea.KeyUp}
			action := handler.HandleKey(keyMsg)

			Expect(action.IsHandled).To(BeTrue())
			Expect(action.NavigationKey).NotTo(BeNil())
			Expect(*action.NavigationKey).To(Equal(KeyUp))
		})

		It("should handle down key", func() {
			keyMsg := tea.KeyMsg{Type: tea.KeyDown}
			action := handler.HandleKey(keyMsg)

			Expect(action.IsHandled).To(BeTrue())
			Expect(action.NavigationKey).NotTo(BeNil())
			Expect(*action.NavigationKey).To(Equal(KeyDown))
		})
	})
})

var _ = Describe("Key handler integration", func() {
	Context("IsNavigationKey function", func() {
		It("should recognize 'e' as KeyEdit", func() {
			Expect(IsNavigationKey("e", KeyEdit)).To(BeTrue())
		})

		It("should recognize 'd' as KeyDelete", func() {
			Expect(IsNavigationKey("d", KeyDelete)).To(BeTrue())
		})

		It("should not recognize 'e' as KeyDelete", func() {
			Expect(IsNavigationKey("e", KeyDelete)).To(BeFalse())
		})

		It("should not recognize 'd' as KeyEdit", func() {
			Expect(IsNavigationKey("d", KeyEdit)).To(BeFalse())
		})
	})

	Context("GetNavigationKeyFromString function", func() {
		It("should return KeyEdit for 'e'", func() {
			navKey := GetNavigationKeyFromString("e")
			Expect(navKey).NotTo(BeNil())
			Expect(*navKey).To(Equal(KeyEdit))
		})

		It("should return KeyDelete for 'd'", func() {
			navKey := GetNavigationKeyFromString("d")
			Expect(navKey).NotTo(BeNil())
			Expect(*navKey).To(Equal(KeyDelete))
		})

		It("should return nil for unknown key", func() {
			navKey := GetNavigationKeyFromString("unknown")
			Expect(navKey).To(BeNil())
		})
	})

	Context("MatchesNavigationKey function", func() {
		It("should match 'e' key to KeyEdit", func() {
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
			Expect(MatchesNavigationKey(keyMsg, KeyEdit)).To(BeTrue())
		})

		It("should match 'd' key to KeyDelete", func() {
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}
			Expect(MatchesNavigationKey(keyMsg, KeyDelete)).To(BeTrue())
		})

		It("should not match 'e' key to KeyDelete", func() {
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
			Expect(MatchesNavigationKey(keyMsg, KeyDelete)).To(BeFalse())
		})
	})
})
