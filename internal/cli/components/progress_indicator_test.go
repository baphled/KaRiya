package components_test

import (
	"github.com/baphled/kariya/internal/cli/components"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ProgressIndicator", func() {
	Describe("NewProgressIndicator", func() {
		It("should create progress indicator with current and total", func() {
			pi := components.NewProgressIndicator(3, 10)
			Expect(pi).NotTo(BeNil())
		})
	})

	Describe("SetCurrent", func() {
		It("should update current value", func() {
			pi := components.NewProgressIndicator(3, 10)
			pi.SetCurrent(5)
			Expect(pi.GetCurrent()).To(Equal(5))
		})
	})

	Describe("SetTotal", func() {
		It("should update total value", func() {
			pi := components.NewProgressIndicator(3, 10)
			pi.SetTotal(20)
			Expect(pi.GetTotal()).To(Equal(20))
		})
	})

	Describe("SetLabel", func() {
		It("should update label text", func() {
			pi := components.NewProgressIndicator(3, 10)
			pi.SetLabel("Processing events")
			Expect(pi.GetLabel()).To(Equal("Processing events"))
		})
	})

	Describe("SetStatus", func() {
		It("should update status type", func() {
			pi := components.NewProgressIndicator(3, 10)
			pi.SetStatus("success")
			Expect(pi.GetStatus()).To(Equal("success"))
		})
	})

	Describe("View", func() {
		Context("with basic progress", func() {
			It("should render step indicator", func() {
				pi := components.NewProgressIndicator(2, 5)
				view := pi.View()
				Expect(view).To(ContainSubstring("2"))
				Expect(view).To(ContainSubstring("5"))
			})
		})

		Context("with label", func() {
			It("should render label with step indicator", func() {
				pi := components.NewProgressIndicator(3, 10)
				pi.SetLabel("Processing")
				view := pi.View()
				Expect(view).To(ContainSubstring("Processing"))
				Expect(view).To(ContainSubstring("3"))
				Expect(view).To(ContainSubstring("10"))
			})
		})

		Context("with progress bar enabled", func() {
			It("should render progress bar", func() {
				pi := components.NewProgressIndicator(5, 10)
				pi.EnableProgressBar(true)
				view := pi.View()
				Expect(view).NotTo(BeEmpty())
			})
		})

		Context("with status colors", func() {
			It("should apply success color", func() {
				pi := components.NewProgressIndicator(10, 10)
				pi.SetStatus("success")
				view := pi.View()
				Expect(view).NotTo(BeEmpty())
			})

			It("should apply error color", func() {
				pi := components.NewProgressIndicator(3, 10)
				pi.SetStatus("error")
				view := pi.View()
				Expect(view).NotTo(BeEmpty())
			})

			It("should apply warning color", func() {
				pi := components.NewProgressIndicator(5, 10)
				pi.SetStatus("warning")
				view := pi.View()
				Expect(view).NotTo(BeEmpty())
			})

			It("should apply info color", func() {
				pi := components.NewProgressIndicator(2, 10)
				pi.SetStatus("info")
				view := pi.View()
				Expect(view).NotTo(BeEmpty())
			})
		})

		Context("with percentage display", func() {
			It("should show percentage when enabled", func() {
				pi := components.NewProgressIndicator(5, 10)
				pi.ShowPercentage(true)
				view := pi.View()
				Expect(view).To(ContainSubstring("50%"))
			})
		})
	})

	Describe("GetPercentage", func() {
		It("should calculate percentage correctly", func() {
			pi := components.NewProgressIndicator(5, 10)
			Expect(pi.GetPercentage()).To(Equal(float64(50)))
		})

		It("should handle zero total", func() {
			pi := components.NewProgressIndicator(0, 0)
			Expect(pi.GetPercentage()).To(Equal(float64(0)))
		})

		It("should handle 100% completion", func() {
			pi := components.NewProgressIndicator(10, 10)
			Expect(pi.GetPercentage()).To(Equal(float64(100)))
		})
	})

	Describe("IsComplete", func() {
		It("should return true when current equals total", func() {
			pi := components.NewProgressIndicator(10, 10)
			Expect(pi.IsComplete()).To(BeTrue())
		})

		It("should return false when current less than total", func() {
			pi := components.NewProgressIndicator(5, 10)
			Expect(pi.IsComplete()).To(BeFalse())
		})
	})
})
