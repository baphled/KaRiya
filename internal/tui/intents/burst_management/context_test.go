package burst_management_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/tui/intents/burst_management"
)

var _ = Describe("Context", func() {
	Describe("IntentValidator", func() {
		Describe("Construction", func() {
			It("should create an empty context", func() {
				ctx := &burst_management.IntentValidator{}
				Expect(ctx).NotTo(BeNil())
			})

			It("should create a context with bursts", func() {
				bursts := []*career.Burst{
					fixtures.Burst("1"),
					fixtures.Burst("2"),
				}
				ctx := &burst_management.IntentValidator{
					Bursts: bursts,
				}
				Expect(ctx.Bursts).To(HaveLen(2))
			})

			It("should create a context with service", func() {
				ctx := &burst_management.IntentValidator{
					Service: nil, // Service would be mocked in real tests
				}
				Expect(ctx.Service).To(BeNil())
			})
		})

		Describe("Validate", func() {
			It("should initialize nil Bursts to empty slice", func() {
				ctx := &burst_management.IntentValidator{
					Bursts: nil,
				}
				err := ctx.Validate()
				Expect(err).NotTo(HaveOccurred())
				Expect(ctx.Bursts).NotTo(BeNil())
				Expect(ctx.Bursts).To(BeEmpty())
			})

			It("should preserve existing Bursts when validating", func() {
				bursts := []*career.Burst{
					fixtures.Burst("1"),
				}
				ctx := &burst_management.IntentValidator{
					Bursts: bursts,
				}
				err := ctx.Validate()
				Expect(err).NotTo(HaveOccurred())
				Expect(ctx.Bursts).To(HaveLen(1))
				Expect(ctx.Bursts[0].ID).To(Equal("1"))
			})
		})

		Describe("LoadBursts", func() {
			It("should handle nil repository gracefully", func() {
				ctx := &burst_management.IntentValidator{
					BurstRepository: nil,
				}
				err := ctx.LoadBursts()
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Describe("CreateBurst", func() {
			It("should handle nil repository gracefully", func() {
				ctx := &burst_management.IntentValidator{
					BurstRepository: nil,
				}
				burst := fixtures.Burst("new-burst")
				err := ctx.CreateBurst(burst)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Describe("UpdateBurst", func() {
			It("should handle nil repository gracefully", func() {
				ctx := &burst_management.IntentValidator{
					BurstRepository: nil,
				}
				burst := fixtures.Burst("burst-1")
				burst.Name = "Updated Burst"
				err := ctx.UpdateBurst(burst)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Describe("DeleteBurst", func() {
			It("should handle nil repository gracefully", func() {
				ctx := &burst_management.IntentValidator{
					BurstRepository: nil,
				}
				err := ctx.DeleteBurst("burst-1")
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Describe("StartNewBurst", func() {
			It("should initialize a new burst for editing", func() {
				ctx := &burst_management.IntentValidator{}
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
				ctx := &burst_management.IntentValidator{
					EditingBurst: fixtures.Burst("burst-1"),
					IsNewBurst:   true,
				}
				ctx.CancelEdit()
				Expect(ctx.EditingBurst).To(BeNil())
				Expect(ctx.IsNewBurst).To(BeFalse())
			})
		})
	})
})
