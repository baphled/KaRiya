package burst_management_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents/burst_management"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
)

var _ = Describe("Result", func() {
	Describe("Construction", func() {
		It("should create an empty result", func() {
			result := &burst_management.Result{}
			Expect(result).NotTo(BeNil())
		})

		It("should store action", func() {
			result := &burst_management.Result{
				Action: "created",
			}
			Expect(result.Action).To(Equal("created"))
		})

		It("should store selected burst", func() {
			burst := fixtures.Burst("burst-1")
			result := &burst_management.Result{
				Burst: burst,
			}
			Expect(result.Burst).To(Equal(burst))
			Expect(result.Burst.ID).To(Equal("burst-1"))
		})

		It("should store all bursts", func() {
			bursts := []*career.Burst{
				fixtures.Burst("burst-1"),
				fixtures.Burst("burst-2"),
			}
			result := &burst_management.Result{
				Bursts: bursts,
			}
			Expect(result.Bursts).To(HaveLen(2))
		})

		It("should store viewed bursts", func() {
			viewedBursts := []*career.Burst{
				fixtures.Burst("burst-1"),
				fixtures.Burst("burst-2"),
			}
			result := &burst_management.Result{
				ViewedBursts: viewedBursts,
			}
			Expect(result.ViewedBursts).To(HaveLen(2))
		})

		It("should store selected index", func() {
			result := &burst_management.Result{
				SelectedIndex: 3,
			}
			Expect(result.SelectedIndex).To(Equal(3))
		})
	})

	Describe("Nil Handling", func() {
		It("should handle nil burst", func() {
			result := &burst_management.Result{
				Burst: nil,
			}
			Expect(result.Burst).To(BeNil())
		})

		It("should handle nil bursts slice", func() {
			result := &burst_management.Result{
				Bursts: nil,
			}
			Expect(result.Bursts).To(BeNil())
		})

		It("should handle nil viewed bursts", func() {
			result := &burst_management.Result{
				ViewedBursts: nil,
			}
			Expect(result.ViewedBursts).To(BeNil())
		})
	})

	Describe("Empty Slices", func() {
		It("should handle empty bursts slice", func() {
			result := &burst_management.Result{
				Bursts: []*career.Burst{},
			}
			Expect(result.Bursts).NotTo(BeNil())
			Expect(result.Bursts).To(BeEmpty())
		})

		It("should handle empty viewed bursts", func() {
			result := &burst_management.Result{
				ViewedBursts: []*career.Burst{},
			}
			Expect(result.ViewedBursts).NotTo(BeNil())
			Expect(result.ViewedBursts).To(BeEmpty())
		})
	})

	Describe("Complete Result", func() {
		It("should store all fields together", func() {
			burst := fixtures.Burst("selected")
			bursts := []*career.Burst{
				fixtures.Burst("burst-1"),
				fixtures.Burst("burst-2"),
			}
			viewedBursts := []*career.Burst{
				fixtures.Burst("burst-1"),
			}

			result := &burst_management.Result{
				Action:        "selected",
				Burst:         burst,
				Bursts:        bursts,
				ViewedBursts:  viewedBursts,
				SelectedIndex: 1,
			}

			Expect(result.Action).To(Equal("selected"))
			Expect(result.Burst.ID).To(Equal("selected"))
			Expect(result.Bursts).To(HaveLen(2))
			Expect(result.ViewedBursts).To(HaveLen(1))
			Expect(result.SelectedIndex).To(Equal(1))
		})
	})

	Describe("Action Types", func() {
		It("should support 'selected' action", func() {
			result := &burst_management.Result{Action: "selected"}
			Expect(result.Action).To(Equal("selected"))
		})

		It("should support 'created' action", func() {
			result := &burst_management.Result{Action: "created"}
			Expect(result.Action).To(Equal("created"))
		})

		It("should support 'updated' action", func() {
			result := &burst_management.Result{Action: "updated"}
			Expect(result.Action).To(Equal("updated"))
		})

		It("should support 'deleted' action", func() {
			result := &burst_management.Result{Action: "deleted"}
			Expect(result.Action).To(Equal("deleted"))
		})

		It("should support 'confirmed' action", func() {
			result := &burst_management.Result{Action: "confirmed"}
			Expect(result.Action).To(Equal("confirmed"))
		})

		It("should support 'cancelled' action", func() {
			result := &burst_management.Result{Action: "cancelled"}
			Expect(result.Action).To(Equal("cancelled"))
		})
	})

	Describe("Cancelled Result Pattern", func() {
		It("should represent cancellation with nil burst", func() {
			result := &burst_management.Result{
				Action: "cancelled",
				Burst:  nil,
			}
			Expect(result.Action).To(Equal("cancelled"))
			Expect(result.Burst).To(BeNil())
		})
	})
})
