package intents_test

import (
	"errors"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/screens"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Mock handler for testing
type mockScreenResultHandler struct {
	navigateCalled bool
	cancelCalled   bool
	submitCalled   bool
	errorCalled    bool
	lastData       interface{}
	lastError      error
	cmdToReturn    tea.Cmd
}

func (m *mockScreenResultHandler) HandleNavigate(result *screens.NavigateResult) tea.Cmd {
	m.navigateCalled = true
	m.lastData = result.Data()
	return m.cmdToReturn
}

func (m *mockScreenResultHandler) HandleCancel(result *screens.CancelResult) tea.Cmd {
	m.cancelCalled = true
	return m.cmdToReturn
}

func (m *mockScreenResultHandler) HandleSubmit(result *screens.SubmitResult) tea.Cmd {
	m.submitCalled = true
	m.lastData = result.Data()
	return m.cmdToReturn
}

func (m *mockScreenResultHandler) HandleError(result *screens.ErrorResult) tea.Cmd {
	m.errorCalled = true
	m.lastError = result.Err
	return m.cmdToReturn
}

var _ = Describe("ScreenResultBehavior", func() {
	Describe("ScreenResultHandler Interface", func() {
		var handler *mockScreenResultHandler

		BeforeEach(func() {
			handler = &mockScreenResultHandler{}
		})

		It("should handle NavigateResult", func() {
			result := &screens.NavigateResult{
				ResultData: "test data",
			}

			cmd := handler.HandleNavigate(result)

			Expect(handler.navigateCalled).To(BeTrue())
			Expect(handler.lastData).To(Equal("test data"))
			Expect(cmd).To(BeNil())
		})

		It("should handle CancelResult", func() {
			result := &screens.CancelResult{}

			cmd := handler.HandleCancel(result)

			Expect(handler.cancelCalled).To(BeTrue())
			Expect(cmd).To(BeNil())
		})

		It("should handle SubmitResult", func() {
			result := &screens.SubmitResult{
				FormData: map[string]string{"key": "value"},
			}

			cmd := handler.HandleSubmit(result)

			Expect(handler.submitCalled).To(BeTrue())
			Expect(handler.lastData).To(Equal(map[string]string{"key": "value"}))
			Expect(cmd).To(BeNil())
		})

		It("should handle ErrorResult", func() {
			testErr := errors.New("test error")
			result := &screens.ErrorResult{
				Err:     testErr,
				Message: "Something went wrong",
			}

			cmd := handler.HandleError(result)

			Expect(handler.errorCalled).To(BeTrue())
			Expect(handler.lastError).To(Equal(testErr))
			Expect(cmd).To(BeNil())
		})
	})

	Describe("ScreenResultDispatcher", func() {
		var (
			dispatcher *intents.ScreenResultDispatcher
			handler    *mockScreenResultHandler
		)

		BeforeEach(func() {
			handler = &mockScreenResultHandler{}
			dispatcher = intents.NewScreenResultDispatcher(handler)
		})

		Context("Dispatch", func() {
			It("should dispatch NavigateResult to HandleNavigate", func() {
				result := &screens.NavigateResult{
					ResultData: "navigate data",
				}

				cmd := dispatcher.Dispatch(result)

				Expect(handler.navigateCalled).To(BeTrue())
				Expect(handler.lastData).To(Equal("navigate data"))
				Expect(cmd).To(BeNil())
			})

			It("should dispatch CancelResult to HandleCancel", func() {
				result := &screens.CancelResult{}

				cmd := dispatcher.Dispatch(result)

				Expect(handler.cancelCalled).To(BeTrue())
				Expect(cmd).To(BeNil())
			})

			It("should dispatch SubmitResult to HandleSubmit", func() {
				result := &screens.SubmitResult{
					FormData: "submit data",
				}

				cmd := dispatcher.Dispatch(result)

				Expect(handler.submitCalled).To(BeTrue())
				Expect(handler.lastData).To(Equal("submit data"))
				Expect(cmd).To(BeNil())
			})

			It("should dispatch ErrorResult to HandleError", func() {
				testErr := errors.New("dispatch error")
				result := &screens.ErrorResult{
					Err:     testErr,
					Message: "Error occurred",
				}

				cmd := dispatcher.Dispatch(result)

				Expect(handler.errorCalled).To(BeTrue())
				Expect(handler.lastError).To(Equal(testErr))
				Expect(cmd).To(BeNil())
			})

			It("should return nil for nil result", func() {
				cmd := dispatcher.Dispatch(nil)

				Expect(cmd).To(BeNil())
				Expect(handler.navigateCalled).To(BeFalse())
				Expect(handler.cancelCalled).To(BeFalse())
				Expect(handler.submitCalled).To(BeFalse())
				Expect(handler.errorCalled).To(BeFalse())
			})

			It("should return nil for unknown result type", func() {
				// Create a custom ScreenResult that's not one of the standard types
				type unknownResult struct {
					screens.ScreenResult
				}
				result := &unknownResult{}

				cmd := dispatcher.Dispatch(result)

				Expect(cmd).To(BeNil())
				Expect(handler.navigateCalled).To(BeFalse())
				Expect(handler.cancelCalled).To(BeFalse())
				Expect(handler.submitCalled).To(BeFalse())
				Expect(handler.errorCalled).To(BeFalse())
			})
		})

		Context("Return Command", func() {
			It("should return the command from HandleNavigate", func() {
				testCmd := func() tea.Msg { return nil }
				handler.cmdToReturn = testCmd

				result := &screens.NavigateResult{
					ResultData: "data",
				}

				cmd := dispatcher.Dispatch(result)

				Expect(cmd).NotTo(BeNil())
				// Can't compare function pointers directly, but we verify it's not nil
			})

			It("should return the command from HandleCancel", func() {
				testCmd := func() tea.Msg { return nil }
				handler.cmdToReturn = testCmd

				result := &screens.CancelResult{}

				cmd := dispatcher.Dispatch(result)

				Expect(cmd).NotTo(BeNil())
			})

			It("should return the command from HandleSubmit", func() {
				testCmd := func() tea.Msg { return nil }
				handler.cmdToReturn = testCmd

				result := &screens.SubmitResult{
					FormData: "data",
				}

				cmd := dispatcher.Dispatch(result)

				Expect(cmd).NotTo(BeNil())
			})

			It("should return the command from HandleError", func() {
				testCmd := func() tea.Msg { return nil }
				handler.cmdToReturn = testCmd

				result := &screens.ErrorResult{
					Err:     errors.New("error"),
					Message: "Error",
				}

				cmd := dispatcher.Dispatch(result)

				Expect(cmd).NotTo(BeNil())
			})
		})

		Context("Type Safety", func() {
			It("should ensure correct result types are passed to handlers", func() {
				// NavigateResult
				navigateResult := &screens.NavigateResult{ResultData: "nav"}
				dispatcher.Dispatch(navigateResult)
				Expect(handler.navigateCalled).To(BeTrue())
				Expect(handler.cancelCalled).To(BeFalse())
				Expect(handler.submitCalled).To(BeFalse())
				Expect(handler.errorCalled).To(BeFalse())

				// Reset
				handler = &mockScreenResultHandler{}
				dispatcher = intents.NewScreenResultDispatcher(handler)

				// CancelResult
				cancelResult := &screens.CancelResult{}
				dispatcher.Dispatch(cancelResult)
				Expect(handler.navigateCalled).To(BeFalse())
				Expect(handler.cancelCalled).To(BeTrue())
				Expect(handler.submitCalled).To(BeFalse())
				Expect(handler.errorCalled).To(BeFalse())

				// Reset
				handler = &mockScreenResultHandler{}
				dispatcher = intents.NewScreenResultDispatcher(handler)

				// SubmitResult
				submitResult := &screens.SubmitResult{FormData: "submit"}
				dispatcher.Dispatch(submitResult)
				Expect(handler.navigateCalled).To(BeFalse())
				Expect(handler.cancelCalled).To(BeFalse())
				Expect(handler.submitCalled).To(BeTrue())
				Expect(handler.errorCalled).To(BeFalse())

				// Reset
				handler = &mockScreenResultHandler{}
				dispatcher = intents.NewScreenResultDispatcher(handler)

				// ErrorResult
				errorResult := &screens.ErrorResult{Err: errors.New("err"), Message: "Error"}
				dispatcher.Dispatch(errorResult)
				Expect(handler.navigateCalled).To(BeFalse())
				Expect(handler.cancelCalled).To(BeFalse())
				Expect(handler.submitCalled).To(BeFalse())
				Expect(handler.errorCalled).To(BeTrue())
			})
		})

		Context("Metadata Preservation", func() {
			It("should not affect metadata on results", func() {
				result := &screens.NavigateResult{
					ResultData: "data",
					Meta:       map[string]interface{}{"key": "value"},
				}

				dispatcher.Dispatch(result)

				// Verify metadata is still present
				Expect(result.Metadata()).To(HaveKeyWithValue("key", "value"))
			})
		})
	})

	Describe("Integration with Intent Pattern", func() {
		It("should simplify handleScreenResult implementation", func() {
			// This test demonstrates the pattern usage
			// Before: 20+ lines of type switching
			// After: 1 line delegation to dispatcher

			handler := &mockScreenResultHandler{}
			dispatcher := intents.NewScreenResultDispatcher(handler)

			// Simulate various screen results
			results := []screens.ScreenResult{
				&screens.NavigateResult{ResultData: "nav"},
				&screens.CancelResult{},
				&screens.SubmitResult{FormData: "submit"},
				&screens.ErrorResult{Err: errors.New("err"), Message: "Error"},
			}

			for _, result := range results {
				dispatcher.Dispatch(result)
			}

			// All handlers should have been called once
			Expect(handler.navigateCalled).To(BeTrue())
			Expect(handler.cancelCalled).To(BeTrue())
			Expect(handler.submitCalled).To(BeTrue())
			Expect(handler.errorCalled).To(BeTrue())
		})
	})
})
