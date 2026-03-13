package selectors_test

import (
	"github.com/baphled/kariya/internal/ui/uikit/selectors"
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
			selected := selector.GetSelected()
			Expect(selected).To(HaveLen(2))
			Expect(selected).To(ContainElement("hiring_manager"))
			Expect(selected).To(ContainElement("recruiter"))
		})

		It("should clear previous selections", func() {
			selector.SetSelected([]string{"hiring_manager"})
			selector.SetSelected([]string{"peer"})
			selected := selector.GetSelected()
			Expect(selected).To(HaveLen(1))
			Expect(selected).To(ContainElement("peer"))
		})

		It("should handle empty list", func() {
			selector.SetSelected([]string{"hiring_manager"})
			selector.SetSelected([]string{})
			Expect(selector.GetSelected()).To(BeEmpty())
		})

		It("should set all three audiences", func() {
			selector.SetSelected([]string{"hiring_manager", "recruiter", "peer"})
			selected := selector.GetSelected()
			Expect(selected).To(HaveLen(3))
			Expect(selected).To(ContainElement("hiring_manager"))
			Expect(selected).To(ContainElement("recruiter"))
			Expect(selected).To(ContainElement("peer"))
		})
	})

	Describe("GetSelected", func() {
		It("should return selected audiences in options order", func() {
			selector.SetSelected([]string{"peer", "hiring_manager"})
			selected := selector.GetSelected()
			Expect(selected).To(HaveLen(2))
			Expect(selected[0]).To(Equal("hiring_manager"))
			Expect(selected[1]).To(Equal("peer"))
		})

		It("should return empty when nothing selected", func() {
			Expect(selector.GetSelected()).To(BeEmpty())
		})

		It("should only return audiences in the options list", func() {
			selector.SetSelected([]string{"hiring_manager", "unknown_audience"})
			selected := selector.GetSelected()
			Expect(selected).To(HaveLen(1))
			Expect(selected).To(ContainElement("hiring_manager"))
		})

		It("should maintain consistent order across multiple calls", func() {
			selector.SetSelected([]string{"recruiter", "peer", "hiring_manager"})
			first := selector.GetSelected()
			second := selector.GetSelected()
			Expect(first).To(Equal(second))
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
			selector.SetSelected([]string{"hiring_manager"})
			Expect(selector.IsSelected("unknown")).To(BeFalse())
		})

		It("should handle all three audiences", func() {
			selector.SetSelected([]string{"hiring_manager", "recruiter", "peer"})
			Expect(selector.IsSelected("hiring_manager")).To(BeTrue())
			Expect(selector.IsSelected("recruiter")).To(BeTrue())
			Expect(selector.IsSelected("peer")).To(BeTrue())
		})
	})

	Describe("Render", func() {
		It("should render all options with checkmarks for selected", func() {
			selector.SetSelected([]string{"hiring_manager"})
			rendered := selector.Render()
			Expect(rendered).To(ContainSubstring("✓ hiring_manager"))
			Expect(rendered).To(ContainSubstring("[ recruiter ]"))
			Expect(rendered).To(ContainSubstring("[ peer ]"))
		})

		It("should render all unselected when nothing selected", func() {
			rendered := selector.Render()
			Expect(rendered).To(ContainSubstring("[ hiring_manager ]"))
			Expect(rendered).To(ContainSubstring("[ recruiter ]"))
			Expect(rendered).To(ContainSubstring("[ peer ]"))
		})

		It("should render all selected when everything selected", func() {
			selector.SetSelected([]string{"hiring_manager", "recruiter", "peer"})
			rendered := selector.Render()
			Expect(rendered).To(ContainSubstring("✓ hiring_manager"))
			Expect(rendered).To(ContainSubstring("✓ recruiter"))
			Expect(rendered).To(ContainSubstring("✓ peer"))
		})

		It("should render multiple selections correctly", func() {
			selector.SetSelected([]string{"recruiter", "peer"})
			rendered := selector.Render()
			Expect(rendered).To(ContainSubstring("[ hiring_manager ]"))
			Expect(rendered).To(ContainSubstring("✓ recruiter"))
			Expect(rendered).To(ContainSubstring("✓ peer"))
		})

		It("should render in consistent order", func() {
			selector.SetSelected([]string{"peer"})
			first := selector.Render()
			second := selector.Render()
			Expect(first).To(Equal(second))
		})
	})
})
