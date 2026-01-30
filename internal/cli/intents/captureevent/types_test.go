package captureevent_test

import (
	"github.com/baphled/kariya/internal/cli/intents"
	ce "github.com/baphled/kariya/internal/cli/intents/captureevent"
	"github.com/baphled/kariya/internal/domain/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Intent Type", func() {
	var intent *ce.Intent

	BeforeEach(func() {
		ctx := &ce.IntentContext{CaptureStrategy: "quick"}
		var err error
		intent, err = ce.NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("IsActive", func() {
		It("should return true for a newly created intent", func() {
			Expect(intent.IsActive()).To(BeTrue())
		})
	})

	Describe("GetResult", func() {
		It("should return nil for a newly created intent", func() {
			Expect(intent.GetResult()).To(BeNil())
		})
	})

	Describe("GetState", func() {
		It("should return choose_strategy for a new intent", func() {
			Expect(intent.GetState()).To(Equal("choose_strategy"))
		})

		It("should reflect state set via SetStateForTesting", func() {
			intent.SetStateForTesting(ce.StateForm)
			Expect(intent.GetState()).To(Equal("form"))
		})

		It("should cycle through all valid states", func() {
			for _, state := range []ce.State{
				ce.StateChooseStrategy,
				ce.StateForm,
				ce.StateReview,
				ce.StateSubmit,
			} {
				intent.SetStateForTesting(state)
				Expect(intent.GetState()).To(Equal(string(state)))
			}
		})
	})

	Describe("GetForm", func() {
		It("should return the capture form", func() {
			Expect(intent.GetForm()).NotTo(BeNil())
		})

		It("should return nil for a nil intent receiver", func() {
			var nilIntent *ce.Intent
			Expect(nilIntent.GetForm()).To(BeNil())
		})
	})

	Describe("SetStateForTesting", func() {
		It("should set the internal state", func() {
			intent.SetStateForTesting(ce.StateReview)
			Expect(intent.GetState()).To(Equal("review"))
		})
	})

	Describe("GetReviewState", func() {
		It("should return non-nil review state for a new intent", func() {
			Expect(intent.GetReviewState()).NotTo(BeNil())
		})

		It("should have empty accepted bursts initially", func() {
			Expect(intent.GetReviewState().AcceptedBursts).To(BeEmpty())
		})

		It("should have empty accepted facts initially", func() {
			Expect(intent.GetReviewState().AcceptedFacts).To(BeEmpty())
		})

		It("should have empty rejected items initially", func() {
			Expect(intent.GetReviewState().RejectedItems).To(BeEmpty())
		})
	})

	Describe("ScreenResultHandler compliance", func() {
		It("should implement the ScreenResultHandler interface", func() {
			// compile-time check via HandleCancel; if it compiles, it implements.
			Expect(intent.HandleCancel).NotTo(BeNil())
			Expect(intent.HandleNavigate).NotTo(BeNil())
			Expect(intent.HandleSubmit).NotTo(BeNil())
			Expect(intent.HandleError).NotTo(BeNil())
		})
	})
})

var _ = Describe("ReviewInferredEventState", func() {
	It("should store event, bursts, and facts", func() {
		state := &ce.ReviewInferredEventState{
			Event:          &career.Event{Text: "test"},
			InferredBursts: []*career.Burst{{Name: "b1"}},
			InferredFacts:  []*career.Fact{{Text: "f1"}},
		}
		Expect(state.Event.Text).To(Equal("test"))
		Expect(state.InferredBursts).To(HaveLen(1))
		Expect(state.InferredFacts).To(HaveLen(1))
	})

	It("should track editing mode and index", func() {
		state := &ce.ReviewInferredEventState{
			EditingMode:  ce.EditingModeMetadata,
			EditingIndex: 2,
		}
		Expect(state.EditingMode).To(Equal(ce.EditingModeMetadata))
		Expect(state.EditingIndex).To(Equal(2))
	})

	It("should track selection state", func() {
		state := &ce.ReviewInferredEventState{
			SelectedItemType: "burst",
			SelectedIndex:    1,
		}
		Expect(state.SelectedItemType).To(Equal("burst"))
		Expect(state.SelectedIndex).To(Equal(1))
	})
})

var _ = Describe("Intent Result", func() {
	var intent *ce.Intent

	BeforeEach(func() {
		ctx := &ce.IntentContext{CaptureStrategy: "quick"}
		var err error
		intent, err = ce.NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Result()", func() {
		It("should return nil when no result is set", func() {
			Expect(intent.Result()).To(BeNil())
		})
	})

	Describe("After cancellation", func() {
		It("should mark result as cancelled", func() {
			intent.HandleCancel(nil)
			result := intent.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(intents.Cancelled))
		})

		It("should deactivate the intent", func() {
			intent.HandleCancel(nil)
			Expect(intent.IsActive()).To(BeFalse())
		})
	})
})
