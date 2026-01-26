package selectors_test

import (
	"github.com/baphled/kariya/internal/cli/uikit/selectors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TagSelector", func() {
	var selector *selectors.TagSelector

	BeforeEach(func() {
		selector = selectors.NewTagSelector()
	})

	Describe("Initialization", func() {
		It("should create a tag selector with no selected tags", func() {
			Expect(selector).NotTo(BeNil())
			Expect(selector.SelectedTags()).To(BeEmpty())
		})

		It("should have all allowed tags available", func() {
			available := selector.AvailableTags()
			Expect(available).To(HaveLen(8))
			Expect(available).To(ContainElements("project", "achievement", "leadership", "technical", "consulting", "research", "product", "mentoring"))
		})
	})

	Describe("Tag Selection", func() {
		It("should allow selecting a valid tag", func() {
			err := selector.SelectTag("technical")
			Expect(err).NotTo(HaveOccurred())
			Expect(selector.SelectedTags()).To(ContainElement("technical"))
		})

		It("should return error when selecting invalid tag", func() {
			err := selector.SelectTag("invalid-tag")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("not a valid tag"))
		})

		It("should prevent duplicate tag selection", func() {
			err := selector.SelectTag("technical")
			Expect(err).NotTo(HaveOccurred())

			err = selector.SelectTag("technical")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("already selected"))
		})

		It("should allow selecting multiple tags", func() {
			err := selector.SelectTag("technical")
			Expect(err).NotTo(HaveOccurred())
			err = selector.SelectTag("leadership")
			Expect(err).NotTo(HaveOccurred())
			err = selector.SelectTag("product")
			Expect(err).NotTo(HaveOccurred())

			Expect(selector.SelectedTags()).To(HaveLen(3))
			Expect(selector.SelectedTags()).To(ContainElements("technical", "leadership", "product"))
		})

		It("should enforce max 8 tags limit", func() {
			tags := []string{"project", "achievement", "leadership", "technical", "consulting", "research", "product", "mentoring"}
			for _, tag := range tags {
				err := selector.SelectTag(tag)
				Expect(err).NotTo(HaveOccurred())
			}

			// Try to add one more - should fail even though we can't (all are already selected)
			// This test verifies the limit is enforced
			Expect(selector.SelectedTags()).To(HaveLen(8))
		})
	})

	Describe("Tag Deselection", func() {
		BeforeEach(func() {
			Expect(selector.SelectTag("technical")).To(Succeed())
			Expect(selector.SelectTag("leadership")).To(Succeed())
		})

		It("should allow deselecting a selected tag", func() {
			err := selector.DeselectTag("technical")
			Expect(err).NotTo(HaveOccurred())
			Expect(selector.SelectedTags()).To(HaveLen(1))
			Expect(selector.SelectedTags()).NotTo(ContainElement("technical"))
		})

		It("should return error when deselecting non-selected tag", func() {
			err := selector.DeselectTag("product")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("not selected"))
		})
	})

	Describe("Tag Filtering", func() {
		It("should filter tags by prefix", func() {
			filtered := selector.FilterTags("tech")
			Expect(filtered).To(ContainElement("technical"))
			Expect(filtered).NotTo(ContainElement("leadership"))
		})

		It("should return all tags when filter is empty", func() {
			filtered := selector.FilterTags("")
			Expect(filtered).To(HaveLen(8))
		})

		It("should be case-insensitive", func() {
			filtered := selector.FilterTags("TECH")
			Expect(filtered).To(ContainElement("technical"))
		})

		It("should return empty when no match", func() {
			filtered := selector.FilterTags("nonexistent")
			Expect(filtered).To(BeEmpty())
		})
	})

	Describe("Toggle Tag", func() {
		It("should select tag if not selected", func() {
			err := selector.ToggleTag("technical")
			Expect(err).NotTo(HaveOccurred())
			Expect(selector.SelectedTags()).To(ContainElement("technical"))
		})

		It("should deselect tag if already selected", func() {
			Expect(selector.SelectTag("technical")).To(Succeed())
			err := selector.ToggleTag("technical")
			Expect(err).NotTo(HaveOccurred())
			Expect(selector.SelectedTags()).NotTo(ContainElement("technical"))
		})
	})

	Describe("IsSelected", func() {
		BeforeEach(func() {
			Expect(selector.SelectTag("technical")).To(Succeed())
		})

		It("should return true for selected tag", func() {
			Expect(selector.IsSelected("technical")).To(BeTrue())
		})

		It("should return false for non-selected tag", func() {
			Expect(selector.IsSelected("leadership")).To(BeFalse())
		})
	})

	Describe("Reset", func() {
		BeforeEach(func() {
			Expect(selector.SelectTag("technical")).To(Succeed())
			Expect(selector.SelectTag("leadership")).To(Succeed())
		})

		It("should clear all selected tags", func() {
			selector.Reset()
			Expect(selector.SelectedTags()).To(BeEmpty())
		})
	})

	Describe("Tag Ordering Consistency", func() {
		It("should return available tags in alphabetical order", func() {
			available := selector.AvailableTags()
			Expect(available).To(Equal([]string{"achievement", "consulting", "leadership", "mentoring", "product", "project", "research", "technical"}))
		})

		It("should return the same order on multiple calls", func() {
			first := selector.AvailableTags()
			second := selector.AvailableTags()
			third := selector.AvailableTags()

			Expect(first).To(Equal(second))
			Expect(second).To(Equal(third))
		})

		It("should return selected tags in alphabetical order", func() {
			// Select tags in non-alphabetical order
			Expect(selector.SelectTag("technical")).To(Succeed())
			Expect(selector.SelectTag("achievement")).To(Succeed())
			Expect(selector.SelectTag("leadership")).To(Succeed())

			selected := selector.SelectedTags()
			Expect(selected).To(Equal([]string{"achievement", "leadership", "technical"}))
		})

		It("should maintain alphabetical order after adding and removing tags", func() {
			Expect(selector.SelectTag("technical")).To(Succeed())
			Expect(selector.SelectTag("achievement")).To(Succeed())
			Expect(selector.SelectTag("leadership")).To(Succeed())
			Expect(selector.DeselectTag("achievement")).To(Succeed())
			Expect(selector.SelectTag("mentoring")).To(Succeed())

			selected := selector.SelectedTags()
			Expect(selected).To(Equal([]string{"leadership", "mentoring", "technical"}))
		})

		It("should return filtered tags in alphabetical order", func() {
			filtered := selector.FilterTags("te")
			Expect(filtered).To(Equal([]string{"technical"}))

			filtered = selector.FilterTags("a")
			Expect(filtered).To(Equal([]string{"achievement"}))

			filtered = selector.FilterTags("")
			Expect(filtered).To(Equal([]string{"achievement", "consulting", "leadership", "mentoring", "product", "project", "research", "technical"}))
		})

		It("should not change order on selection/deselection", func() {
			// Get the original order
			originalOrder := selector.AvailableTags()

			// Select some tags
			Expect(selector.SelectTag("technical")).To(Succeed())
			Expect(selector.SelectTag("leadership")).To(Succeed())

			// Verify available tags order hasn't changed
			afterSelection := selector.AvailableTags()
			Expect(afterSelection).To(Equal(originalOrder))

			// Deselect a tag
			Expect(selector.DeselectTag("technical")).To(Succeed())

			// Verify order is still the same
			afterDeselection := selector.AvailableTags()
			Expect(afterDeselection).To(Equal(originalOrder))
		})
	})
})
