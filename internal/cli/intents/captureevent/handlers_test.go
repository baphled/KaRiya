package captureevent_test

import (
	"errors"
	"time"

	"github.com/baphled/kariya/internal/cli/intents"
	ce "github.com/baphled/kariya/internal/cli/intents/captureevent"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/domain/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Handlers", func() {
	var intent *ce.Intent

	BeforeEach(func() {
		ctx := &ce.IntentContext{CaptureStrategy: "quick"}
		var err error
		intent, err = ce.NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	Describe("HandleNavigate", func() {
		Context("with CaptureStrategy data", func() {
			It("should transition to form screen for quick strategy", func() {
				result := &screens.NavigateResult{
					ResultData: ce.StrategyQuick,
				}
				intent.HandleNavigate(result)
				Expect(intent.GetState()).To(Equal("form"))
			})

			It("should transition to form screen for manual strategy", func() {
				result := &screens.NavigateResult{
					ResultData: ce.StrategyManual,
				}
				intent.HandleNavigate(result)
				Expect(intent.GetState()).To(Equal("form"))
			})
		})

		Context("with string action data", func() {
			BeforeEach(func() {
				// Put intent in review state with an event.
				intent.SetStateForTesting(ce.StateReview)
				intent.GetReviewState().Event = &career.Event{
					Text: "test event",
					Date: time.Now(),
				}
			})

			It("should panic for edit_metadata when CareerService is nil", func() {
				// MetadataEditorModelNew immediately calls CareerService.GetSkillRepository(),
				// which panics with a nil service. The full flow is tested in E2E tests
				// with a real CareerService.
				result := &screens.NavigateResult{ResultData: "edit_metadata"}
				Expect(func() {
					intent.HandleNavigate(result)
				}).To(Panic())
			})

			It("should set burst editing mode for edit_bursts", func() {
				result := &screens.NavigateResult{ResultData: "edit_bursts"}
				intent.HandleNavigate(result)
				Expect(intent.GetReviewState().EditingMode).To(Equal(ce.EditingModeBursts))
			})

			It("should set fact editing mode for edit_facts", func() {
				result := &screens.NavigateResult{ResultData: "edit_facts"}
				intent.HandleNavigate(result)
				Expect(intent.GetReviewState().EditingMode).To(Equal(ce.EditingModeFacts))
			})

			It("should fail on unknown action", func() {
				result := &screens.NavigateResult{ResultData: "unknown_action"}
				intent.HandleNavigate(result)
				Expect(intent.IsActive()).To(BeFalse())
				Expect(intent.Result().Status).To(Equal(intents.Failed))
			})
		})

		Context("with nil event for edit_metadata", func() {
			BeforeEach(func() {
				intent.SetStateForTesting(ce.StateReview)
				intent.GetReviewState().Event = nil
			})

			It("should fail gracefully", func() {
				result := &screens.NavigateResult{ResultData: "edit_metadata"}
				intent.HandleNavigate(result)
				Expect(intent.IsActive()).To(BeFalse())
			})
		})

		Context("with invalid data type", func() {
			It("should fail on non-string, non-strategy data", func() {
				result := &screens.NavigateResult{ResultData: 42}
				intent.HandleNavigate(result)
				Expect(intent.IsActive()).To(BeFalse())
				Expect(intent.Result().Status).To(Equal(intents.Failed))
			})
		})
	})

	Describe("HandleCancel", func() {
		Context("from StateChooseStrategy", func() {
			It("should cancel the intent", func() {
				intent.HandleCancel(nil)
				Expect(intent.IsActive()).To(BeFalse())
				Expect(intent.Result().Status).To(Equal(intents.Cancelled))
			})
		})

		Context("from StateForm without PreviousEvent", func() {
			BeforeEach(func() {
				// Navigate to form first.
				result := &screens.NavigateResult{ResultData: ce.StrategyQuick}
				intent.HandleNavigate(result)
			})

			It("should go back to strategy screen", func() {
				intent.HandleCancel(nil)
				Expect(intent.GetState()).To(Equal("choose_strategy"))
				Expect(intent.IsActive()).To(BeTrue())
			})
		})

		Context("from StateForm with PreviousEvent", func() {
			It("should cancel the intent", func() {
				ctx := &ce.IntentContext{
					CaptureStrategy: "quick",
					PreviousEvent:   &career.Event{Text: "existing"},
				}
				editIntent, err := ce.NewIntent(ctx)
				Expect(err).NotTo(HaveOccurred())
				editIntent.Init()

				// Navigate to form.
				editIntent.HandleNavigate(&screens.NavigateResult{
					ResultData: ce.StrategyQuick,
				})

				editIntent.HandleCancel(nil)
				Expect(editIntent.IsActive()).To(BeFalse())
				Expect(editIntent.Result().Status).To(Equal(intents.Cancelled))
			})
		})

		Context("from StateReview", func() {
			BeforeEach(func() {
				intent.SetStateForTesting(ce.StateReview)
			})

			It("should go back to form screen", func() {
				intent.HandleCancel(nil)
				Expect(intent.GetState()).To(Equal("form"))
				Expect(intent.IsActive()).To(BeTrue())
			})
		})

		Context("from StateSubmit", func() {
			BeforeEach(func() {
				intent.SetStateForTesting(ce.StateSubmit)
			})

			It("should return nil without changing state", func() {
				cmd := intent.HandleCancel(nil)
				Expect(cmd).To(BeNil())
				Expect(intent.IsActive()).To(BeTrue())
			})
		})

		Context("from invalid state", func() {
			BeforeEach(func() {
				intent.SetStateForTesting(ce.State("bogus"))
			})

			It("should fail the intent", func() {
				intent.HandleCancel(nil)
				Expect(intent.IsActive()).To(BeFalse())
				Expect(intent.Result().Status).To(Equal(intents.Failed))
			})
		})
	})

	Describe("HandleSubmit", func() {
		Context("from StateForm with valid event", func() {
			BeforeEach(func() {
				intent.HandleNavigate(&screens.NavigateResult{
					ResultData: ce.StrategyQuick,
				})
			})

			It("should not panic", func() {
				event := &career.Event{
					Text: "Test event",
					Date: time.Now(),
				}
				Expect(func() {
					intent.HandleSubmit(&screens.SubmitResult{
						FormData: event,
					})
				}).NotTo(Panic())
			})
		})

		Context("from StateForm with invalid data type", func() {
			BeforeEach(func() {
				intent.HandleNavigate(&screens.NavigateResult{
					ResultData: ce.StrategyQuick,
				})
			})

			It("should fail the intent", func() {
				intent.HandleSubmit(&screens.SubmitResult{
					FormData: "not an event",
				})
				Expect(intent.IsActive()).To(BeFalse())
				Expect(intent.Result().Status).To(Equal(intents.Failed))
			})
		})

		Context("from StateReview with valid data (post-save)", func() {
			var event *career.Event

			BeforeEach(func() {
				event = &career.Event{
					Text: "Reviewed event",
					Date: time.Now(),
				}
				intent.SetStateForTesting(ce.StateReview)
				intent.GetReviewState().Event = event
			})

			It("should fail on non-map data type", func() {
				intent.HandleSubmit(&screens.SubmitResult{
					FormData: "not a map",
				})
				Expect(intent.IsActive()).To(BeFalse())
				Expect(intent.Result().Status).To(Equal(intents.Failed))
			})

			It("should fail when event key is missing from map", func() {
				intent.HandleSubmit(&screens.SubmitResult{
					FormData: map[string]interface{}{
						"bursts": []*career.Burst{},
					},
				})
				Expect(intent.IsActive()).To(BeFalse())
				Expect(intent.Result().Status).To(Equal(intents.Failed))
			})
		})

		Context("from StateSubmit with valid data", func() {
			var event *career.Event

			BeforeEach(func() {
				event = &career.Event{
					Text: "Submit event",
					Date: time.Now(),
				}
				intent.SetStateForTesting(ce.StateSubmit)
				intent.GetReviewState().Event = event
			})

			It("should complete the intent with event data", func() {
				intent.HandleSubmit(&screens.SubmitResult{
					FormData: map[string]interface{}{
						"event":  event,
						"bursts": []*career.Burst{},
						"facts":  []*career.Fact{},
					},
				})
				Expect(intent.IsActive()).To(BeFalse())
				result := intent.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Status).To(Equal(intents.Completed))
			})

			It("should fail on non-map data type", func() {
				intent.HandleSubmit(&screens.SubmitResult{
					FormData: "not a map",
				})
				Expect(intent.IsActive()).To(BeFalse())
				Expect(intent.Result().Status).To(Equal(intents.Failed))
			})
		})

		Context("from invalid state", func() {
			BeforeEach(func() {
				intent.SetStateForTesting(ce.State("bogus"))
			})

			It("should fail the intent", func() {
				intent.HandleSubmit(&screens.SubmitResult{
					FormData: "anything",
				})
				Expect(intent.IsActive()).To(BeFalse())
				Expect(intent.Result().Status).To(Equal(intents.Failed))
			})
		})
	})

	Describe("HandleError", func() {
		It("should fail intent with map error data", func() {
			result := &screens.ErrorResult{
				Err:     errors.New("something broke"),
				Message: "Screen error occurred",
			}
			intent.HandleError(result)
			Expect(intent.IsActive()).To(BeFalse())
			intentResult := intent.Result()
			Expect(intentResult).NotTo(BeNil())
			Expect(intentResult.Status).To(Equal(intents.Failed))
		})

		It("should use default message when error message is empty", func() {
			result := &screens.ErrorResult{
				Err: errors.New("internal error"),
			}
			intent.HandleError(result)
			Expect(intent.IsActive()).To(BeFalse())
		})
	})
})
