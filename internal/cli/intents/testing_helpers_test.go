package intents_test

import (
	"errors"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/themes"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Testing Helpers", func() {
	Describe("MockIntent", func() {
		var mock *intents.MockIntent

		BeforeEach(func() {
			mock = intents.NewMockIntent()
		})

		It("should track Init calls", func() {
			Expect(mock.InitCalled).To(BeFalse())
			mock.Init()
			Expect(mock.InitCalled).To(BeTrue())
		})

		It("should track Update calls", func() {
			Expect(mock.UpdateCalled).To(Equal(0))
			mock.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(mock.UpdateCalled).To(Equal(1))
			mock.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(mock.UpdateCalled).To(Equal(2))
		})

		It("should store messages", func() {
			Expect(mock.Messages).To(HaveLen(0))
			mock.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(mock.Messages).To(HaveLen(1))
		})

		It("should track View calls", func() {
			Expect(mock.ViewCalled).To(BeFalse())
			_ = mock.View()
			Expect(mock.ViewCalled).To(BeTrue())
		})

		It("should return set result", func() {
			result := intents.NewCompletedResult[interface{}]("test data")
			mock.SetResult(result)
			Expect(mock.Result()).To(Equal(result))
		})
	})

	Describe("ThemeAwareMockIntent", func() {
		var mock *intents.ThemeAwareMockIntent

		BeforeEach(func() {
			mock = intents.NewThemeAwareMockIntent()
		})

		It("should embed MockIntent functionality", func() {
			mock.Init()
			Expect(mock.InitCalled).To(BeTrue())
		})

		It("should store and return theme manager", func() {
			tm := themes.NewThemeManager()
			mock.SetThemeManager(tm)
			Expect(mock.GetThemeManager()).To(Equal(tm))
		})

		It("should return nil when theme manager not set", func() {
			Expect(mock.GetThemeManager()).To(BeNil())
		})
	})

	Describe("MockIntentWithSelection", func() {
		var mock *intents.MockIntentWithSelection

		BeforeEach(func() {
			mock = intents.NewMockIntentWithSelection()
		})

		It("should embed MockIntent functionality", func() {
			mock.Init()
			Expect(mock.InitCalled).To(BeTrue())
		})

		It("should start with zero selection", func() {
			Expect(mock.GetSelectedIndex()).To(Equal(0))
		})

		It("should store and return selection index", func() {
			mock.SetSelectedIndex(5)
			Expect(mock.GetSelectedIndex()).To(Equal(5))
		})
	})

	Describe("MockIntentWithStateMachine", func() {
		var mock *intents.MockIntentWithStateMachine

		BeforeEach(func() {
			mock = intents.NewMockIntentWithStateMachine("initial")
		})

		It("should start with initial state", func() {
			Expect(mock.CurrentState).To(Equal("initial"))
		})

		It("should record state history", func() {
			Expect(mock.GetStateHistory()).To(ConsistOf("initial"))
		})

		It("should transition to new state", func() {
			mock.TransitionTo("working")
			Expect(mock.CurrentState).To(Equal("working"))
		})

		It("should record transition history", func() {
			mock.TransitionTo("working")
			mock.TransitionTo("complete")
			Expect(mock.GetStateHistory()).To(Equal([]string{"initial", "working", "complete"}))
		})

		It("should track transitions with details", func() {
			mock.TransitionToWithTrigger("working", "enter_key")
			transitions := mock.GetTransitions()
			Expect(transitions).To(HaveLen(1))
			Expect(transitions[0].From).To(Equal("initial"))
			Expect(transitions[0].To).To(Equal("working"))
			Expect(transitions[0].Trigger).To(Equal("enter_key"))
			Expect(transitions[0].Timestamp).NotTo(BeZero())
		})

		It("should track multiple transitions", func() {
			mock.TransitionToWithTrigger("working", "enter_key")
			mock.TransitionToWithTrigger("complete", "confirm")
			transitions := mock.GetTransitions()
			Expect(transitions).To(HaveLen(2))
			Expect(transitions[0].To).To(Equal("working"))
			Expect(transitions[1].To).To(Equal("complete"))
		})

		It("should reset state history", func() {
			mock.TransitionTo("working")
			mock.TransitionTo("complete")
			mock.Reset("initial")
			Expect(mock.CurrentState).To(Equal("initial"))
			Expect(mock.GetStateHistory()).To(ConsistOf("initial"))
			Expect(mock.GetTransitions()).To(BeEmpty())
		})
	})

	Describe("MockIntentWithErrors", func() {
		var mock *intents.MockIntentWithErrors

		BeforeEach(func() {
			mock = intents.NewMockIntentWithErrors()
		})

		It("should work normally without errors configured", func() {
			cmd := mock.Init()
			Expect(cmd).To(BeNil())
			Expect(mock.InitCalled).To(BeTrue())
			Expect(mock.GetErrorCount()).To(Equal(0))
		})

		It("should return error on Init when configured", func() {
			expectedErr := errors.New("init error")
			mock.SetInitError(expectedErr)
			cmd := mock.Init()
			Expect(cmd).NotTo(BeNil())
			// Execute the command to get the error message
			msg := cmd()
			errMsg, ok := msg.(intents.ErrorMsg)
			Expect(ok).To(BeTrue())
			Expect(errMsg.Err).To(Equal(expectedErr))
			Expect(mock.GetErrorCount()).To(Equal(1))
		})

		It("should return error on Update when configured", func() {
			expectedErr := errors.New("update error")
			mock.SetUpdateError(expectedErr, 0) // Fail immediately
			cmd := mock.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			errMsg, ok := msg.(intents.ErrorMsg)
			Expect(ok).To(BeTrue())
			Expect(errMsg.Err).To(Equal(expectedErr))
			Expect(mock.GetErrorCount()).To(Equal(1))
		})

		It("should fail on Nth update when configured", func() {
			expectedErr := errors.New("delayed error")
			mock.SetUpdateError(expectedErr, 2) // Fail on 3rd update (0-indexed)

			// First two updates should succeed
			cmd := mock.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(cmd).To(BeNil())
			cmd = mock.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(cmd).To(BeNil())

			// Third update should fail
			cmd = mock.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			errMsg, ok := msg.(intents.ErrorMsg)
			Expect(ok).To(BeTrue())
			Expect(errMsg.Err).To(Equal(expectedErr))
			Expect(mock.GetErrorCount()).To(Equal(1))
		})

		It("should track error count across multiple errors", func() {
			mock.SetInitError(errors.New("init error"))
			mock.SetUpdateError(errors.New("update error"), 0)

			mock.Init()
			Expect(mock.GetErrorCount()).To(Equal(1))

			mock.Update(tea.KeyMsg{})
			Expect(mock.GetErrorCount()).To(Equal(2))
		})

		It("should clear errors when reset", func() {
			mock.SetInitError(errors.New("init error"))
			mock.Init()
			Expect(mock.GetErrorCount()).To(Equal(1))

			mock.ClearErrors()
			Expect(mock.GetErrorCount()).To(Equal(0))
			Expect(mock.ErrorOnInit).To(BeNil())
			Expect(mock.ErrorOnUpdate).To(BeNil())
		})
	})
})
