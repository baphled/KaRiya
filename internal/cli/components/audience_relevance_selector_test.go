package components_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/components"
)

var _ = Describe("AudienceRelevanceSelector", func() {
	var selector *components.AudienceRelevanceSelector

	Describe("NewAudienceRelevanceSelector", func() {
		It("should create a new selector with default options", func() {
			selector = components.NewAudienceRelevanceSelector()
			Expect(selector).ToNot(BeNil())
		})

		It("should start with no selections", func() {
			selector = components.NewAudienceRelevanceSelector()
			selected := selector.GetSelected()
			Expect(selected).To(BeEmpty())
		})

		It("should have default options available", func() {
			selector = components.NewAudienceRelevanceSelector()
			// Verify by checking if known options can be selected
			selector.SetSelected([]string{"hiring_manager"})
			Expect(selector.IsSelected("hiring_manager")).To(BeTrue())
		})
	})

	Describe("SetSelected", func() {
		BeforeEach(func() {
			selector = components.NewAudienceRelevanceSelector()
		})

		It("should set single audience", func() {
			selector.SetSelected([]string{"hiring_manager"})
			Expect(selector.IsSelected("hiring_manager")).To(BeTrue())
			Expect(selector.IsSelected("recruiter")).To(BeFalse())
			Expect(selector.IsSelected("peer")).To(BeFalse())
		})

		It("should set multiple audiences", func() {
			selector.SetSelected([]string{"hiring_manager", "recruiter"})
			Expect(selector.IsSelected("hiring_manager")).To(BeTrue())
			Expect(selector.IsSelected("recruiter")).To(BeTrue())
			Expect(selector.IsSelected("peer")).To(BeFalse())
		})

		It("should set all audiences", func() {
			selector.SetSelected([]string{"hiring_manager", "recruiter", "peer"})
			Expect(selector.IsSelected("hiring_manager")).To(BeTrue())
			Expect(selector.IsSelected("recruiter")).To(BeTrue())
			Expect(selector.IsSelected("peer")).To(BeTrue())
		})

		It("should clear previous selections", func() {
			selector.SetSelected([]string{"hiring_manager"})
			Expect(selector.IsSelected("hiring_manager")).To(BeTrue())

			selector.SetSelected([]string{"recruiter"})
			Expect(selector.IsSelected("hiring_manager")).To(BeFalse())
			Expect(selector.IsSelected("recruiter")).To(BeTrue())
		})

		It("should handle empty selection", func() {
			selector.SetSelected([]string{"hiring_manager"})
			selector.SetSelected([]string{})
			Expect(selector.GetSelected()).To(BeEmpty())
		})

		It("should handle nil selection", func() {
			selector.SetSelected(nil)
			Expect(selector.GetSelected()).To(BeEmpty())
		})

		It("should ignore invalid audience types", func() {
			selector.SetSelected([]string{"invalid_audience"})
			// Invalid audiences won't be in GetSelected since options list filters them
			selected := selector.GetSelected()
			Expect(selected).ToNot(ContainElement("invalid_audience"))
		})
	})

	Describe("GetSelected", func() {
		BeforeEach(func() {
			selector = components.NewAudienceRelevanceSelector()
		})

		It("should return empty list when nothing selected", func() {
			selected := selector.GetSelected()
			Expect(selected).To(BeEmpty())
		})

		It("should return selected audiences in order", func() {
			selector.SetSelected([]string{"recruiter", "hiring_manager"})
			selected := selector.GetSelected()

			// Should return in options order, not input order
			Expect(selected).To(HaveLen(2))
			Expect(selected).To(ContainElement("hiring_manager"))
			Expect(selected).To(ContainElement("recruiter"))
		})

		It("should return only valid options", func() {
			selector.SetSelected([]string{"hiring_manager", "invalid", "peer"})
			selected := selector.GetSelected()

			Expect(selected).To(ContainElement("hiring_manager"))
			Expect(selected).To(ContainElement("peer"))
			Expect(selected).ToNot(ContainElement("invalid"))
		})
	})

	Describe("IsSelected", func() {
		BeforeEach(func() {
			selector = components.NewAudienceRelevanceSelector()
		})

		It("should return false for unselected audience", func() {
			Expect(selector.IsSelected("hiring_manager")).To(BeFalse())
		})

		It("should return true for selected audience", func() {
			selector.SetSelected([]string{"hiring_manager"})
			Expect(selector.IsSelected("hiring_manager")).To(BeTrue())
		})

		It("should return false for invalid audience", func() {
			Expect(selector.IsSelected("invalid_type")).To(BeFalse())
		})

		It("should return false for empty string", func() {
			Expect(selector.IsSelected("")).To(BeFalse())
		})

		It("should handle multiple checks", func() {
			selector.SetSelected([]string{"hiring_manager", "peer"})

			Expect(selector.IsSelected("hiring_manager")).To(BeTrue())
			Expect(selector.IsSelected("recruiter")).To(BeFalse())
			Expect(selector.IsSelected("peer")).To(BeTrue())
		})
	})

	Describe("Render", func() {
		BeforeEach(func() {
			selector = components.NewAudienceRelevanceSelector()
		})

		It("should render with no selections", func() {
			rendered := selector.Render()
			Expect(rendered).ToNot(BeEmpty())
			Expect(rendered).To(ContainSubstring("[ hiring_manager ]"))
			Expect(rendered).To(ContainSubstring("[ recruiter ]"))
			Expect(rendered).To(ContainSubstring("[ peer ]"))
		})

		It("should render with single selection", func() {
			selector.SetSelected([]string{"hiring_manager"})
			rendered := selector.Render()

			Expect(rendered).To(ContainSubstring("[✓ hiring_manager]"))
			Expect(rendered).To(ContainSubstring("[ recruiter ]"))
			Expect(rendered).To(ContainSubstring("[ peer ]"))
		})

		It("should render with multiple selections", func() {
			selector.SetSelected([]string{"hiring_manager", "peer"})
			rendered := selector.Render()

			Expect(rendered).To(ContainSubstring("[✓ hiring_manager]"))
			Expect(rendered).To(ContainSubstring("[ recruiter ]"))
			Expect(rendered).To(ContainSubstring("[✓ peer]"))
		})

		It("should render with all selections", func() {
			selector.SetSelected([]string{"hiring_manager", "recruiter", "peer"})
			rendered := selector.Render()

			Expect(rendered).To(ContainSubstring("[✓ hiring_manager]"))
			Expect(rendered).To(ContainSubstring("[✓ recruiter]"))
			Expect(rendered).To(ContainSubstring("[✓ peer]"))
		})

		It("should render all options even if not selected", func() {
			selector.SetSelected([]string{})
			rendered := selector.Render()

			// All three options should appear
			Expect(rendered).To(ContainSubstring("hiring_manager"))
			Expect(rendered).To(ContainSubstring("recruiter"))
			Expect(rendered).To(ContainSubstring("peer"))
		})

		It("should produce consistent output", func() {
			selector.SetSelected([]string{"hiring_manager"})

			rendered1 := selector.Render()
			rendered2 := selector.Render()
			rendered3 := selector.Render()

			Expect(rendered1).To(Equal(rendered2))
			Expect(rendered2).To(Equal(rendered3))
		})
	})

	Describe("Edge Cases", func() {
		BeforeEach(func() {
			selector = components.NewAudienceRelevanceSelector()
		})

		It("should handle duplicate selections", func() {
			selector.SetSelected([]string{"hiring_manager", "hiring_manager", "hiring_manager"})
			selected := selector.GetSelected()

			// Should only contain one instance
			count := 0
			for _, s := range selected {
				if s == "hiring_manager" {
					count++
				}
			}
			Expect(count).To(Equal(1))
		})

		It("should handle mixed valid and invalid audiences", func() {
			selector.SetSelected([]string{"hiring_manager", "invalid1", "recruiter", "invalid2"})
			selected := selector.GetSelected()

			Expect(selected).To(HaveLen(2))
			Expect(selected).To(ContainElement("hiring_manager"))
			Expect(selected).To(ContainElement("recruiter"))
		})

		It("should maintain state across multiple operations", func() {
			// Set initial selection
			selector.SetSelected([]string{"hiring_manager"})
			Expect(selector.IsSelected("hiring_manager")).To(BeTrue())

			// Render shouldn't change state
			selector.Render()
			Expect(selector.IsSelected("hiring_manager")).To(BeTrue())

			// Get selected shouldn't change state
			selector.GetSelected()
			Expect(selector.IsSelected("hiring_manager")).To(BeTrue())
		})

		It("should handle case-sensitive audience names", func() {
			selector.SetSelected([]string{"HIRING_MANAGER"})

			// Should not match due to case sensitivity
			Expect(selector.IsSelected("hiring_manager")).To(BeFalse())
			Expect(selector.IsSelected("HIRING_MANAGER")).To(BeTrue())
		})
	})

	Describe("Integration Scenarios", func() {
		It("should support full selection workflow", func() {
			selector = components.NewAudienceRelevanceSelector()

			// Start with nothing selected
			Expect(selector.GetSelected()).To(BeEmpty())

			// Select hiring manager
			selector.SetSelected([]string{"hiring_manager"})
			Expect(selector.GetSelected()).To(Equal([]string{"hiring_manager"}))

			// Add recruiter
			selector.SetSelected([]string{"hiring_manager", "recruiter"})
			Expect(selector.GetSelected()).To(HaveLen(2))

			// Replace with peer only
			selector.SetSelected([]string{"peer"})
			Expect(selector.GetSelected()).To(Equal([]string{"peer"}))

			// Clear all
			selector.SetSelected([]string{})
			Expect(selector.GetSelected()).To(BeEmpty())
		})

		It("should support render-check-modify workflow", func() {
			selector = components.NewAudienceRelevanceSelector()

			// Initial render
			rendered1 := selector.Render()
			Expect(rendered1).To(ContainSubstring("[ hiring_manager ]"))

			// Select and render
			selector.SetSelected([]string{"hiring_manager"})
			rendered2 := selector.Render()
			Expect(rendered2).To(ContainSubstring("[✓ hiring_manager]"))
			Expect(rendered2).ToNot(Equal(rendered1))

			// Verify selection
			Expect(selector.IsSelected("hiring_manager")).To(BeTrue())
		})
	})
})
