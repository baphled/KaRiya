package captureevent_test

import (
	"github.com/baphled/kariya/internal/cli/intents"
	ce "github.com/baphled/kariya/internal/cli/intents/captureevent"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Intent", func() {
	Describe("NewIntent", func() {
		It("should create an intent with a valid context", func() {
			ctx := &ce.IntentContext{CaptureStrategy: "quick"}
			intent, err := ce.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(intent).NotTo(BeNil())
		})

		It("should fail with an empty CaptureStrategy", func() {
			ctx := &ce.IntentContext{}
			intent, err := ce.NewIntent(ctx)
			Expect(err).To(HaveOccurred())
			Expect(intent).To(BeNil())
		})

		It("should start in StateChooseStrategy", func() {
			ctx := &ce.IntentContext{CaptureStrategy: "quick"}
			intent, _ := ce.NewIntent(ctx)
			Expect(intent.GetState()).To(Equal("choose_strategy"))
		})

		It("should be active after creation", func() {
			ctx := &ce.IntentContext{CaptureStrategy: "quick"}
			intent, _ := ce.NewIntent(ctx)
			Expect(intent.IsActive()).To(BeTrue())
		})

		It("should have nil form screen before transition", func() {
			ctx := &ce.IntentContext{CaptureStrategy: "quick"}
			intent, _ := ce.NewIntent(ctx)
			Expect(intent.GetFormScreen()).To(BeNil())
		})

		It("should initialise review state with empty slices", func() {
			ctx := &ce.IntentContext{CaptureStrategy: "quick"}
			intent, _ := ce.NewIntent(ctx)
			review := intent.GetReviewState()
			Expect(review).NotTo(BeNil())
			Expect(review.AcceptedBursts).To(BeEmpty())
			Expect(review.AcceptedFacts).To(BeEmpty())
			Expect(review.RejectedItems).To(BeEmpty())
		})
	})

	Describe("Init", func() {
		var intent *ce.Intent

		BeforeEach(func() {
			ctx := &ce.IntentContext{CaptureStrategy: "quick"}
			var err error
			intent, err = ce.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should not panic", func() {
			Expect(func() { intent.Init() }).NotTo(Panic())
		})
	})

	Describe("Update", func() {
		var intent *ce.Intent

		BeforeEach(func() {
			ctx := &ce.IntentContext{CaptureStrategy: "quick"}
			var err error
			intent, err = ce.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		Context("when intent is not active", func() {
			BeforeEach(func() {
				// Cancel to deactivate.
				intent.HandleCancel(nil)
			})

			It("should return nil", func() {
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
				Expect(cmd).To(BeNil())
			})
		})

		Context("when receiving SubmitCompleteMsg", func() {
			It("should not panic", func() {
				Expect(func() {
					intent.Update(ce.SubmitCompleteMsg{})
				}).NotTo(Panic())
			})
		})

		Context("when receiving SubmitErrorMsg", func() {
			It("should not panic", func() {
				Expect(func() {
					intent.Update(ce.SubmitErrorMsg{
						Code:    "TEST_ERROR",
						Message: "test",
					})
				}).NotTo(Panic())
			})
		})

		Context("when receiving DismissModalMsg", func() {
			It("should not panic when no submit modal exists", func() {
				Expect(func() {
					intent.Update(ce.DismissModalMsg{})
				}).NotTo(Panic())
			})
		})

		Context("global key handling", func() {
			It("should return quit command on '?' key without panic", func() {
				Expect(func() {
					intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
				}).NotTo(Panic())
			})
		})

		Context("screen result dispatch", func() {
			It("should process key events without panic", func() {
				Expect(func() {
					intent.Update(tea.KeyMsg{Type: tea.KeyDown})
					intent.Update(tea.KeyMsg{Type: tea.KeyUp})
					intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
				}).NotTo(Panic())
			})
		})
	})

	Describe("View", func() {
		var intent *ce.Intent

		BeforeEach(func() {
			ctx := &ce.IntentContext{CaptureStrategy: "quick"}
			var err error
			intent, err = ce.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should return a non-empty view", func() {
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should not contain panic text", func() {
			view := intent.View()
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		Context("when intent is inactive", func() {
			BeforeEach(func() {
				intent.HandleCancel(nil)
			})

			It("should return an inactive message", func() {
				view := intent.View()
				Expect(view).To(ContainSubstring("not active"))
			})
		})
	})

	Describe("Result", func() {
		var intent *ce.Intent

		BeforeEach(func() {
			ctx := &ce.IntentContext{CaptureStrategy: "quick"}
			var err error
			intent, err = ce.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should return nil before completion", func() {
			Expect(intent.Result()).To(BeNil())
		})

		Context("after cancellation", func() {
			It("should return a cancelled result", func() {
				intent.HandleCancel(nil)
				result := intent.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Status).To(Equal(intents.Cancelled))
			})
		})

		Context("after submit completion", func() {
			It("should return completed result with event data", func() {
				event := fixtures.EventWith("", "Test event", "", "")

				// Transition to review state with data.
				intent.SetStateForTesting(ce.StateSubmit)
				intent.GetReviewState().Event = event

				submitData := map[string]interface{}{
					"event":  event,
					"bursts": []*career.Burst{},
					"facts":  []*career.Fact{},
				}
				intent.HandleSubmit(&screens.SubmitResult{FormData: submitData})

				Expect(intent.IsActive()).To(BeFalse())
				result := intent.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Status).To(Equal(intents.Completed))
			})
		})
	})
})
