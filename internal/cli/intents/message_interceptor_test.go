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

			interceptor.Intercept(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(called).To(BeTrue())
		})

		It("should return command from back handler", func() {
			expectedCmd := tea.Quit
			interceptor.OnBack(func() tea.Cmd {
				return expectedCmd
			})

			cmd := interceptor.Intercept(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("OnQuit", func() {
		It("should call quit handler when q pressed", func() {
			called := false
			interceptor.OnQuit(func() tea.Cmd {
				called = true
				return tea.Quit
			})

			interceptor.Intercept(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(called).To(BeTrue())
		})
	})

	Describe("OnHelp", func() {
		It("should call help handler when ? pressed", func() {
			called := false
			interceptor.OnHelp(func() tea.Cmd {
				called = true
				return nil
			})

			interceptor.Intercept(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
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

		It("should return global key handler command when matched", func() {
			expectedCmd := tea.Quit
			interceptor.OnQuit(func() tea.Cmd {
				return expectedCmd
			})

			cmd := interceptor.InterceptOr(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}, func() tea.Cmd {
				return nil
			})

			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("Intercept", func() {
		It("should return nil when no handler matches", func() {
			cmd := interceptor.Intercept(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(cmd).To(BeNil())
		})

		It("should return handler command when matched", func() {
			expectedCmd := tea.Quit
			interceptor.OnQuit(func() tea.Cmd {
				return expectedCmd
			})

			cmd := interceptor.Intercept(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("Chaining", func() {
		It("should support method chaining", func() {
			backCalled := false
			quitCalled := false
			helpCalled := false

			result := interceptor.
				OnBack(func() tea.Cmd {
					backCalled = true
					return nil
				}).
				OnQuit(func() tea.Cmd {
					quitCalled = true
					return nil
				}).
				OnHelp(func() tea.Cmd {
					helpCalled = true
					return nil
				})

			Expect(result).To(Equal(interceptor))

			// Test each handler works
			interceptor.Intercept(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(backCalled).To(BeTrue())

			interceptor.Intercept(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(quitCalled).To(BeTrue())

			interceptor.Intercept(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
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
})
