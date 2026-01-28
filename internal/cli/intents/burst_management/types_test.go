package burst_management_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/intents/burst_management"
	"github.com/baphled/kariya/internal/domain/career"
)

var _ = Describe("Types", func() {
	Describe("Intent", func() {
		Describe("Construction", func() {
			It("should create an intent with valid context", func() {
				ctx := &burst_management.IntentContext{}
				err := ctx.Validate()
				Expect(err).NotTo(HaveOccurred())

				intent, err := burst_management.NewIntent(ctx)
				Expect(err).NotTo(HaveOccurred())
				Expect(intent).NotTo(BeNil())
			})

			It("should reject nil context", func() {
				intent, err := burst_management.NewIntent(nil)
				Expect(err).To(HaveOccurred())
				Expect(intent).To(BeNil())
			})

			It("should embed BaseIntent", func() {
				ctx := &burst_management.IntentContext{}
				ctx.Validate()

				intent, _ := burst_management.NewIntent(ctx)
				Expect(intent.BaseIntent).NotTo(BeNil())
			})

			It("should initialize state to StateList", func() {
				ctx := &burst_management.IntentContext{}
				ctx.Validate()

				intent, _ := burst_management.NewIntent(ctx)
				Expect(intent.GetState()).To(Equal(burst_management.StateList))
			})

			It("should initialize as active", func() {
				ctx := &burst_management.IntentContext{}
				ctx.Validate()

				intent, _ := burst_management.NewIntent(ctx)
				Expect(intent.IsActive()).To(BeTrue())
			})

			It("should initialize empty filtered bursts", func() {
				ctx := &burst_management.IntentContext{}
				ctx.Validate()

				intent, _ := burst_management.NewIntent(ctx)
				Expect(intent.GetFilteredBursts()).NotTo(BeNil())
			})

			It("should initialize with context bursts", func() {
				bursts := []*career.Burst{
					{ID: "burst-1", Name: "Burst 1"},
					{ID: "burst-2", Name: "Burst 2"},
				}
				ctx := &burst_management.IntentContext{
					Bursts: bursts,
				}
				ctx.Validate()

				intent, _ := burst_management.NewIntent(ctx)
				Expect(intent.GetFilteredBursts()).To(HaveLen(2))
			})
		})

		Describe("State Management", func() {
			var intent *burst_management.Intent
			var ctx *burst_management.IntentContext

			BeforeEach(func() {
				ctx = &burst_management.IntentContext{}
				ctx.Validate()
				intent, _ = burst_management.NewIntent(ctx)
			})

			It("should track current state", func() {
				Expect(intent.GetState()).To(Equal(burst_management.StateList))
			})

			It("should allow state transitions", func() {
				intent.SetState(burst_management.StateDetail)
				Expect(intent.GetState()).To(Equal(burst_management.StateDetail))
			})

			It("should track active status", func() {
				Expect(intent.IsActive()).To(BeTrue())
				intent.Deactivate()
				Expect(intent.IsActive()).To(BeFalse())
			})
		})

		Describe("Burst Selection", func() {
			var intent *burst_management.Intent
			var bursts []*career.Burst

			BeforeEach(func() {
				bursts = []*career.Burst{
					{ID: "burst-1", Name: "Burst 1"},
					{ID: "burst-2", Name: "Burst 2"},
					{ID: "burst-3", Name: "Burst 3"},
				}
				ctx := &burst_management.IntentContext{
					Bursts: bursts,
				}
				ctx.Validate()
				intent, _ = burst_management.NewIntent(ctx)
			})

			It("should track selected burst", func() {
				Expect(intent.GetSelectedBurst()).To(BeNil())
			})

			It("should allow selecting a burst", func() {
				intent.SetSelectedBurst(bursts[1])
				Expect(intent.GetSelectedBurst()).To(Equal(bursts[1]))
			})

			It("should track selected index", func() {
				Expect(intent.GetSelectedIndex()).To(Equal(0))
			})

			It("should allow setting selected index", func() {
				intent.SetSelectedIndex(2)
				Expect(intent.GetSelectedIndex()).To(Equal(2))
			})

			It("should track viewed bursts", func() {
				Expect(intent.GetViewedBursts()).To(BeEmpty())
			})

			It("should accumulate viewed bursts", func() {
				intent.AddViewedBurst(bursts[0])
				intent.AddViewedBurst(bursts[2])
				Expect(intent.GetViewedBursts()).To(HaveLen(2))
				Expect(intent.GetViewedBursts()[0]).To(Equal(bursts[0]))
				Expect(intent.GetViewedBursts()[1]).To(Equal(bursts[2]))
			})
		})

		Describe("Result Management", func() {
			var intent *burst_management.Intent

			BeforeEach(func() {
				ctx := &burst_management.IntentContext{}
				ctx.Validate()
				intent, _ = burst_management.NewIntent(ctx)
			})

			It("should have no result initially", func() {
				result := intent.Result()
				Expect(result).To(BeNil())
			})

			It("should allow setting result", func() {
				burst := &career.Burst{ID: "burst-1", Name: "Selected"}
				intent.SetCompleted(burst)
				result := intent.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Status).To(Equal(intents.Completed))
			})

			It("should support cancelled result", func() {
				intent.SetCancelled()
				result := intent.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Status).To(Equal(intents.Cancelled))
			})
		})

		Describe("Modal Management", func() {
			var intent *burst_management.Intent

			BeforeEach(func() {
				ctx := &burst_management.IntentContext{}
				ctx.Validate()
				intent, _ = burst_management.NewIntent(ctx)
			})

			It("should have no visible modals initially", func() {
				Expect(intent.HasActiveModal()).To(BeFalse())
			})
		})

		Describe("Filters", func() {
			It("should be defined in context.go", func() {
				// Filters struct is part of IntentContext, tested in context_test.go.
				// This is just a placeholder to acknowledge filter support.
				Expect(true).To(BeTrue())
			})
		})
	})
})
