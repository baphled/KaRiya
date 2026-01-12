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

	Describe("OnContextAwareBack", func() {
		It("should call cancel handler when in edit mode", func() {
			cancelCalled := false
			goBackCalled := false

			interceptor.OnContextAwareBack(
				func() bool { return true },                        // isEditMode
				func() tea.Cmd { goBackCalled = true; return nil }, // goBack
				func() tea.Cmd { cancelCalled = true; return nil }, // cancel
			)

			interceptor.Intercept(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(cancelCalled).To(BeTrue())
			Expect(goBackCalled).To(BeFalse())
		})

		It("should call goBack handler when not in edit mode", func() {
			cancelCalled := false
			goBackCalled := false

			interceptor.OnContextAwareBack(
				func() bool { return false },                       // isEditMode
				func() tea.Cmd { goBackCalled = true; return nil }, // goBack
				func() tea.Cmd { cancelCalled = true; return nil }, // cancel
			)

			interceptor.Intercept(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(goBackCalled).To(BeTrue())
			Expect(cancelCalled).To(BeFalse())
		})

		It("should return command from cancel handler", func() {
			expectedCmd := tea.Quit
			interceptor.OnContextAwareBack(
				func() bool { return true },           // isEditMode
				func() tea.Cmd { return nil },         // goBack
				func() tea.Cmd { return expectedCmd }, // cancel
			)

			cmd := interceptor.Intercept(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(cmd).NotTo(BeNil())
		})

		It("should return command from goBack handler", func() {
			expectedCmd := tea.Quit
			interceptor.OnContextAwareBack(
				func() bool { return false },          // isEditMode
				func() tea.Cmd { return expectedCmd }, // goBack
				func() tea.Cmd { return nil },         // cancel
			)

			cmd := interceptor.Intercept(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("OnModalAwareBack", func() {
		It("should call closeModal handler when modal is active", func() {
			closeModalCalled := false
			goBackCalled := false

			interceptor.OnModalAwareBack(
				func() bool { return true },                            // hasActiveModal
				func() tea.Cmd { closeModalCalled = true; return nil }, // closeModal
				func() tea.Cmd { goBackCalled = true; return nil },     // goBack
			)

			interceptor.Intercept(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(closeModalCalled).To(BeTrue())
			Expect(goBackCalled).To(BeFalse())
		})

		It("should call goBack handler when no modal is active", func() {
			closeModalCalled := false
			goBackCalled := false

			interceptor.OnModalAwareBack(
				func() bool { return false },                           // hasActiveModal
				func() tea.Cmd { closeModalCalled = true; return nil }, // closeModal
				func() tea.Cmd { goBackCalled = true; return nil },     // goBack
			)

			interceptor.Intercept(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(goBackCalled).To(BeTrue())
			Expect(closeModalCalled).To(BeFalse())
		})

		It("should return command from closeModal handler", func() {
			expectedCmd := tea.Quit
			interceptor.OnModalAwareBack(
				func() bool { return true },           // hasActiveModal
				func() tea.Cmd { return expectedCmd }, // closeModal
				func() tea.Cmd { return nil },         // goBack
			)

			cmd := interceptor.Intercept(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(cmd).NotTo(BeNil())
		})

		It("should return command from goBack handler", func() {
			expectedCmd := tea.Quit
			interceptor.OnModalAwareBack(
				func() bool { return false },          // hasActiveModal
				func() tea.Cmd { return nil },         // closeModal
				func() tea.Cmd { return expectedCmd }, // goBack
			)

			cmd := interceptor.Intercept(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("StandardQuitHandler", func() {
		It("should return tea.Quit command", func() {
			handler := intents.StandardQuitHandler()
			cmd := handler()
			Expect(cmd).NotTo(BeNil())
		})

		It("should work with MessageInterceptor", func() {
			interceptor.OnQuit(intents.StandardQuitHandler())
			cmd := interceptor.Intercept(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("StandardHelpHandler", func() {
		It("should toggle help on BaseIntent", func() {
			baseIntent := intents.NewBaseIntent()
			handler := intents.StandardHelpHandler(baseIntent)

			// Initially help is not visible
			Expect(baseIntent.IsHelpVisible()).To(BeFalse())

			// Call handler
			handler()

			// Help should now be visible
			Expect(baseIntent.IsHelpVisible()).To(BeTrue())

			// Call again to toggle back
			handler()

			// Help should be hidden again
			Expect(baseIntent.IsHelpVisible()).To(BeFalse())
		})

		It("should work with MessageInterceptor", func() {
			baseIntent := intents.NewBaseIntent()
			interceptor.OnHelp(intents.StandardHelpHandler(baseIntent))

			// Initially help is not visible
			Expect(baseIntent.IsHelpVisible()).To(BeFalse())

			// Press help key
			interceptor.Intercept(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})

			// Help should now be visible
			Expect(baseIntent.IsHelpVisible()).To(BeTrue())
		})
	})
})
