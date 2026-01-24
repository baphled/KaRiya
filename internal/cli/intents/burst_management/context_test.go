package burst_management_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents/burst_management"
	"github.com/baphled/kariya/internal/domain/career"
)

var _ = Describe("Context", func() {
	Describe("IntentContext", func() {
		Describe("Construction", func() {
			It("should create an empty context", func() {
				ctx := &burst_management.IntentContext{}
				Expect(ctx).NotTo(BeNil())
			})

			It("should create a context with bursts", func() {
				bursts := []*career.Burst{
					{ID: "1", Name: "Burst 1"},
					{ID: "2", Name: "Burst 2"},
				}
				ctx := &burst_management.IntentContext{
					Bursts: bursts,
				}
				Expect(ctx.Bursts).To(HaveLen(2))
			})

			It("should create a context with service", func() {
				ctx := &burst_management.IntentContext{
					Service: nil, // Service would be mocked in real tests
				}
				Expect(ctx.Service).To(BeNil())
			})
		})

		Describe("Validate", func() {
			It("should initialize nil Bursts to empty slice", func() {
				ctx := &burst_management.IntentContext{
					Bursts: nil,
				}
				err := ctx.Validate()
				Expect(err).NotTo(HaveOccurred())
				Expect(ctx.Bursts).NotTo(BeNil())
				Expect(ctx.Bursts).To(BeEmpty())
			})

			It("should initialize nil Context to Background", func() {
				ctx := &burst_management.IntentContext{
					Context: nil,
				}
				err := ctx.Validate()
				Expect(err).NotTo(HaveOccurred())
				Expect(ctx.Context).NotTo(BeNil())
			})

			It("should preserve existing Bursts when validating", func() {
				bursts := []*career.Burst{
					{ID: "1", Name: "Burst 1"},
				}
				ctx := &burst_management.IntentContext{
					Bursts: bursts,
				}
				err := ctx.Validate()
				Expect(err).NotTo(HaveOccurred())
				Expect(ctx.Bursts).To(HaveLen(1))
				Expect(ctx.Bursts[0].ID).To(Equal("1"))
			})

			It("should preserve existing Context when validating", func() {
				existingCtx := context.WithValue(context.Background(), "key", "value")
				ctx := &burst_management.IntentContext{
					Context: existingCtx,
				}
				err := ctx.Validate()
				Expect(err).NotTo(HaveOccurred())
				Expect(ctx.Context).To(Equal(existingCtx))
			})
		})

		Describe("LoadBursts", func() {
			It("should handle nil repository gracefully", func() {
				ctx := &burst_management.IntentContext{
					BurstRepository: nil,
					Context:         context.Background(),
				}
				err := ctx.LoadBursts()
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Describe("CreateBurst", func() {
			It("should handle nil repository gracefully", func() {
				ctx := &burst_management.IntentContext{
					BurstRepository: nil,
					Context:         context.Background(),
				}
				burst := &career.Burst{
					ID:   "new-burst",
					Name: "New Burst",
				}
				err := ctx.CreateBurst(burst)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Describe("UpdateBurst", func() {
			It("should handle nil repository gracefully", func() {
				ctx := &burst_management.IntentContext{
					BurstRepository: nil,
					Context:         context.Background(),
				}
				burst := &career.Burst{
					ID:   "burst-1",
					Name: "Updated Burst",
				}
				err := ctx.UpdateBurst(burst)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Describe("DeleteBurst", func() {
			It("should handle nil repository gracefully", func() {
				ctx := &burst_management.IntentContext{
					BurstRepository: nil,
					Context:         context.Background(),
				}
				err := ctx.DeleteBurst("burst-1")
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Describe("StartNewBurst", func() {
			It("should initialize a new burst for editing", func() {
				ctx := &burst_management.IntentContext{}
				ctx.StartNewBurst()
				Expect(ctx.EditingBurst).NotTo(BeNil())
				Expect(ctx.EditingBurst.ID).To(BeEmpty())
				Expect(ctx.EditingBurst.Name).To(BeEmpty())
				Expect(ctx.EditingBurst.EventIDs).To(BeEmpty())
				Expect(ctx.IsNewBurst).To(BeTrue())
			})
		})

		Describe("CancelEdit", func() {
			It("should clear editing state", func() {
				ctx := &burst_management.IntentContext{
					EditingBurst: &career.Burst{ID: "burst-1"},
					IsNewBurst:   true,
				}
				ctx.CancelEdit()
				Expect(ctx.EditingBurst).To(BeNil())
				Expect(ctx.IsNewBurst).To(BeFalse())
			})
		})
	})
})
