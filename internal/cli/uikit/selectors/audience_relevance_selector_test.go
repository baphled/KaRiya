package selectors_test

import (
	"github.com/baphled/kariya/internal/cli/uikit/selectors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("AudienceRelevanceSelector", func() {
	var selector *selectors.AudienceRelevanceSelector

	BeforeEach(func() {
		selector = selectors.NewAudienceRelevanceSelector()
	})

	Describe("NewAudienceRelevanceSelector", func() {
		It("should create a selector with no selected audiences", func() {
			Expect(selector).NotTo(BeNil())
			Expect(selector.GetSelected()).To(BeEmpty())
		})
	})

	Describe("SetSelected", func() {
		It("should set selected audiences", func() {
			selector.SetSelected([]string{"hiring_manager", "recruiter"})
			Expect(selector.GetSelected()).To(HaveLen(2))
			Expect(selector.GetSelected()).To(ContainElement("hiring_manager"))
			Expect(selector.GetSelected()).To(ContainElement("recruiter"))
		})

		It("should replace previous selections", func() {
			selector.SetSelected([]string{"hiring_manager"})
			selector.SetSelected([]string{"peer"})

			selected := selector.GetSelected()
			Expect(selected).To(HaveLen(1))
			Expect(selected).To(ContainElement("peer"))
			Expect(selected).NotTo(ContainElement("hiring_manager"))
		})

		It("should handle empty list", func() {
			selector.SetSelected([]string{"hiring_manager"})
			selector.SetSelected([]string{})

			Expect(selector.GetSelected()).To(BeEmpty())
		})
	})

	Describe("GetSelected", func() {
		It("should return empty when nothing selected", func() {
			Expect(selector.GetSelected()).To(BeEmpty())
		})

		It("should return selected audiences in options order", func() {
			selector.SetSelected([]string{"peer", "hiring_manager"})

			selected := selector.GetSelected()
			Expect(selected).To(HaveLen(2))
			Expect(selected[0]).To(Equal("hiring_manager"))
			Expect(selected[1]).To(Equal("peer"))
		})

		It("should return all three when all selected", func() {
			selector.SetSelected([]string{"hiring_manager", "recruiter", "peer"})

			selected := selector.GetSelected()
			Expect(selected).To(HaveLen(3))
			Expect(selected[0]).To(Equal("hiring_manager"))
			Expect(selected[1]).To(Equal("recruiter"))
			Expect(selected[2]).To(Equal("peer"))
		})
	})

	Describe("IsSelected", func() {
		It("should return true for selected audience", func() {
			selector.SetSelected([]string{"recruiter"})
			Expect(selector.IsSelected("recruiter")).To(BeTrue())
		})

		It("should return false for unselected audience", func() {
			Expect(selector.IsSelected("recruiter")).To(BeFalse())
		})

		It("should return false for unknown audience", func() {
			Expect(selector.IsSelected("unknown")).To(BeFalse())
		})
	})

	Describe("Render", func() {
		It("should return non-empty string with no selections", func() {
			rendered := selector.Render()
			Expect(rendered).NotTo(BeEmpty())
		})

		It("should include all option labels", func() {
			rendered := selector.Render()
			Expect(rendered).To(ContainSubstring("hiring_manager"))
			Expect(rendered).To(ContainSubstring("recruiter"))
			Expect(rendered).To(ContainSubstring("peer"))
		})

		It("should show check mark for selected audiences", func() {
			selector.SetSelected([]string{"recruiter"})
			rendered := selector.Render()
			Expect(rendered).To(ContainSubstring("✓ recruiter"))
		})
	})
})
