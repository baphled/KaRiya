package intents_test

import (
	"github.com/baphled/kariya/internal/cli/intents"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("MessageInterceptor", func() {
	var interceptor *intents.MessageInterceptor

	BeforeEach(func() {
		interceptor = intents.NewMessageInterceptor()
	})

	Describe("OnBack", func() {
		It("should call back handler when escape pressed", func() {
			called := false
			interceptor.OnBack(func() tea.Cmd {
				called = true
				return nil
			})

			interceptor.InterceptOr(tea.KeyMsg{Type: tea.KeyEsc}, func() tea.Cmd { return nil })
			Expect(called).To(BeTrue())
		})

		It("should return command from back handler", func() {
			expectedCmd := tea.Quit
			interceptor.OnBack(func() tea.Cmd {
				return expectedCmd
			})

			cmd := interceptor.InterceptOr(tea.KeyMsg{Type: tea.KeyEsc}, func() tea.Cmd { return nil })
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("OnHelp", func() {
		It("should call help handler when ? pressed", func() {
			called := false
			interceptor.OnHelp(func() tea.Cmd {
				called = true
				return nil
			})

			interceptor.InterceptOr(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}, func() tea.Cmd { return nil })
			Expect(called).To(BeTrue())
		})
	})

	Describe("InterceptOr", func() {
		It("should call fallback when no global key matched", func() {
			fallbackCalled := false
			interceptor.InterceptOr(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}, func() tea.Cmd {
				fallbackCalled = true
				return nil
			})

			Expect(fallbackCalled).To(BeTrue())
		})

		It("should NOT call fallback when global key matched", func() {
			fallbackCalled := false
			interceptor.
				OnBack(func() tea.Cmd { return nil }).
				InterceptOr(tea.KeyMsg{Type: tea.KeyEsc}, func() tea.Cmd {
					fallbackCalled = true
					return nil
				})

			Expect(fallbackCalled).To(BeFalse())
		})

		It("should return fallback command when no global key matched", func() {
			expectedCmd := tea.Quit
			cmd := interceptor.InterceptOr(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}, func() tea.Cmd {
				return expectedCmd
			})

			Expect(cmd).NotTo(BeNil())
		})

		It("should return nil from fallback when no global key matched and fallback returns nil", func() {
			cmd := interceptor.InterceptOr(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}, func() tea.Cmd {
				return nil
			})
			Expect(cmd).To(BeNil())
		})
	})

	Describe("Chaining", func() {
		It("should support method chaining for back and help handlers", func() {
			backCalled := false
			helpCalled := false

			result := interceptor.
				OnBack(func() tea.Cmd {
					backCalled = true
					return nil
				}).
				OnHelp(func() tea.Cmd {
					helpCalled = true
					return nil
				})

			Expect(result).To(Equal(interceptor))

			// Test each handler works via InterceptOr.
			interceptor.InterceptOr(tea.KeyMsg{Type: tea.KeyEsc}, func() tea.Cmd { return nil })
			Expect(backCalled).To(BeTrue())

			interceptor.InterceptOr(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}, func() tea.Cmd { return nil })
			Expect(helpCalled).To(BeTrue())
		})
	})

	Describe("Non-KeyMsg handling", func() {
		It("should call fallback for non-KeyMsg messages", func() {
			fallbackCalled := false
			interceptor.
				OnBack(func() tea.Cmd { return nil }).
				InterceptOr(tea.WindowSizeMsg{}, func() tea.Cmd {
					fallbackCalled = true
					return nil
				})

			Expect(fallbackCalled).To(BeTrue())
		})
	})

	Describe("StandardHelpHandler", func() {
		It("should toggle help on BaseIntent", func() {
			baseIntent := intents.NewBaseIntent()
			handler := intents.StandardHelpHandler(baseIntent)

			// Initially help is not visible.
			Expect(baseIntent.IsHelpVisible()).To(BeFalse())

			// Call handler.
			handler()

			// Help should now be visible.
			Expect(baseIntent.IsHelpVisible()).To(BeTrue())

			// Call again to toggle back.
			handler()

			// Help should be hidden again.
			Expect(baseIntent.IsHelpVisible()).To(BeFalse())
		})

		It("should work with MessageInterceptor", func() {
			baseIntent := intents.NewBaseIntent()
			interceptor.OnHelp(intents.StandardHelpHandler(baseIntent))

			// Initially help is not visible.
			Expect(baseIntent.IsHelpVisible()).To(BeFalse())

			// Press help key.
			interceptor.InterceptOr(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}, func() tea.Cmd { return nil })

			// Help should now be visible.
			Expect(baseIntent.IsHelpVisible()).To(BeTrue())
		})
	})
})
