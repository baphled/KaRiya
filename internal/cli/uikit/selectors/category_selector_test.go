package selectors_test

import (
	"github.com/baphled/kariya/internal/cli/uikit/selectors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CategorySelector", func() {
	var selector *selectors.CategorySelector

	BeforeEach(func() {
		selector = selectors.NewCategorySelector()
	})

	Describe("NewCategorySelector", func() {
		It("should create a new selector with no selected categories", func() {
			Expect(selector).NotTo(BeNil())
			Expect(selector.SelectedCategories()).To(BeEmpty())
		})
	})

	Describe("SelectCategory", func() {
		It("should select a valid category", func() {
			err := selector.SelectCategory("technical")
			Expect(err).NotTo(HaveOccurred())
			Expect(selector.IsSelected("technical")).To(BeTrue())
		})

		It("should return error for invalid category", func() {
			err := selector.SelectCategory("invalid")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("not a valid category"))
		})

		It("should return error for duplicate category", func() {
			err := selector.SelectCategory("technical")
			Expect(err).NotTo(HaveOccurred())

			err = selector.SelectCategory("technical")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("already selected"))
		})

		It("should handle case-insensitive selection", func() {
			err := selector.SelectCategory("Technical")
			Expect(err).NotTo(HaveOccurred())
			Expect(selector.IsSelected("technical")).To(BeTrue())
		})

		It("should select multiple categories", func() {
			err := selector.SelectCategory("technical")
			Expect(err).NotTo(HaveOccurred())

			err = selector.SelectCategory("leadership")
			Expect(err).NotTo(HaveOccurred())

			err = selector.SelectCategory("product")
			Expect(err).NotTo(HaveOccurred())

			selected := selector.SelectedCategories()
			Expect(selected).To(HaveLen(3))
			Expect(selected).To(ContainElement("technical"))
			Expect(selected).To(ContainElement("leadership"))
			Expect(selected).To(ContainElement("product"))
		})
	})

	Describe("DeselectCategory", func() {
		It("should deselect a selected category", func() {
			Expect(selector.SelectCategory("technical")).To(Succeed())
			err := selector.DeselectCategory("technical")
			Expect(err).NotTo(HaveOccurred())
			Expect(selector.IsSelected("technical")).To(BeFalse())
		})

		It("should return error when deselecting unselected category", func() {
			err := selector.DeselectCategory("technical")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("not selected"))
		})

		It("should handle case-insensitive deselection", func() {
			Expect(selector.SelectCategory("technical")).To(Succeed())
			err := selector.DeselectCategory("Technical")
			Expect(err).NotTo(HaveOccurred())
			Expect(selector.IsSelected("technical")).To(BeFalse())
		})
	})

	Describe("ToggleCategory", func() {
		It("should select unselected category", func() {
			err := selector.ToggleCategory("technical")
			Expect(err).NotTo(HaveOccurred())
			Expect(selector.IsSelected("technical")).To(BeTrue())
		})

		It("should deselect selected category", func() {
			Expect(selector.SelectCategory("technical")).To(Succeed())
			err := selector.ToggleCategory("technical")
			Expect(err).NotTo(HaveOccurred())
			Expect(selector.IsSelected("technical")).To(BeFalse())
		})

		It("should return error for invalid category", func() {
			err := selector.ToggleCategory("invalid")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("FilterCategories", func() {
		It("should return all categories when no prefix", func() {
			filtered := selector.FilterCategories("")
			Expect(filtered).To(HaveLen(11))
			Expect(filtered).To(ContainElement("technical"))
			Expect(filtered).To(ContainElement("leadership"))
			Expect(filtered).To(ContainElement("product"))
			Expect(filtered).To(ContainElement("consulting"))
			Expect(filtered).To(ContainElement("research"))
			Expect(filtered).To(ContainElement("mentoring"))
			Expect(filtered).To(ContainElement("communication"))
			Expect(filtered).To(ContainElement("collaboration"))
			Expect(filtered).To(ContainElement("problem-solving"))
			Expect(filtered).To(ContainElement("project-management"))
			Expect(filtered).To(ContainElement("architecture"))
		})

		It("should filter categories by prefix", func() {
			filtered := selector.FilterCategories("te")
			Expect(filtered).To(HaveLen(1))
			Expect(filtered).To(ContainElement("technical"))
		})

		It("should handle case-insensitive filtering", func() {
			filtered := selector.FilterCategories("LE")
			Expect(filtered).To(HaveLen(1))
			Expect(filtered).To(ContainElement("leadership"))
		})

		It("should return empty when no match", func() {
			filtered := selector.FilterCategories("xyz")
			Expect(filtered).To(BeEmpty())
		})

		It("should filter multiple matches", func() {
			filtered := selector.FilterCategories("m")
			Expect(filtered).To(HaveLen(1))
			Expect(filtered).To(ContainElement("mentoring"))
		})
	})

	Describe("AvailableCategories", func() {
		It("should return all available categories including soft skills", func() {
			categories := selector.AvailableCategories()
			Expect(categories).To(HaveLen(11))
			Expect(categories).To(ContainElement("technical"))
			Expect(categories).To(ContainElement("leadership"))
			Expect(categories).To(ContainElement("product"))
			Expect(categories).To(ContainElement("consulting"))
			Expect(categories).To(ContainElement("research"))
			Expect(categories).To(ContainElement("mentoring"))
			Expect(categories).To(ContainElement("communication"))
			Expect(categories).To(ContainElement("collaboration"))
			Expect(categories).To(ContainElement("problem-solving"))
			Expect(categories).To(ContainElement("project-management"))
			Expect(categories).To(ContainElement("architecture"))
		})

		It("should return sorted categories", func() {
			categories := selector.AvailableCategories()
			// Check that they're sorted
			for i := range len(categories) - 1 {
				Expect(categories[i] <= categories[i+1]).To(BeTrue())
			}
		})
	})

	Describe("SelectedCategories", func() {
		It("should return empty slice when no categories selected", func() {
			categories := selector.SelectedCategories()
			Expect(categories).To(BeEmpty())
		})

		It("should return selected categories in sorted order", func() {
			Expect(selector.SelectCategory("product")).To(Succeed())
			Expect(selector.SelectCategory("technical")).To(Succeed())
			Expect(selector.SelectCategory("mentoring")).To(Succeed())

			categories := selector.SelectedCategories()
			Expect(categories).To(HaveLen(3))
			// Verify sorted order
			Expect(categories[0]).To(Equal("mentoring"))
			Expect(categories[1]).To(Equal("product"))
			Expect(categories[2]).To(Equal("technical"))
		})
	})

	Describe("Clear", func() {
		It("should remove all selected categories", func() {
			Expect(selector.SelectCategory("technical")).To(Succeed())
			Expect(selector.SelectCategory("leadership")).To(Succeed())
			Expect(selector.SelectCategory("product")).To(Succeed())

			selector.Clear()

			Expect(selector.SelectedCategories()).To(BeEmpty())
			Expect(selector.IsSelected("technical")).To(BeFalse())
		})

		It("should work on empty selector", func() {
			selector.Clear()
			Expect(selector.SelectedCategories()).To(BeEmpty())
		})
	})

	Describe("SetSelected", func() {
		It("should set selected categories from list", func() {
			err := selector.SetSelected([]string{"technical", "leadership"})
			Expect(err).NotTo(HaveOccurred())

			selected := selector.SelectedCategories()
			Expect(selected).To(HaveLen(2))
			Expect(selected).To(ContainElement("technical"))
			Expect(selected).To(ContainElement("leadership"))
		})

		It("should clear previous selections", func() {
			Expect(selector.SelectCategory("product")).To(Succeed())
			err := selector.SetSelected([]string{"technical", "leadership"})
			Expect(err).NotTo(HaveOccurred())

			selected := selector.SelectedCategories()
			Expect(selected).To(HaveLen(2))
			Expect(selected).NotTo(ContainElement("product"))
		})

		It("should return error for invalid category in list", func() {
			err := selector.SetSelected([]string{"technical", "invalid"})
			Expect(err).To(HaveOccurred())
		})

		It("should work with empty list", func() {
			Expect(selector.SelectCategory("technical")).To(Succeed())
			err := selector.SetSelected([]string{})
			Expect(err).NotTo(HaveOccurred())
			Expect(selector.SelectedCategories()).To(BeEmpty())
		})
	})

	Describe("GetCategoryDescription", func() {
		It("should return description for technical", func() {
			desc := selectors.GetCategoryDescription("technical")
			Expect(desc).To(ContainSubstring("Technical skills"))
		})

		It("should return description for leadership", func() {
			desc := selectors.GetCategoryDescription("leadership")
			Expect(desc).To(ContainSubstring("Leadership"))
		})

		It("should return description for product", func() {
			desc := selectors.GetCategoryDescription("product")
			Expect(desc).To(ContainSubstring("Product"))
		})

		It("should return description for consulting", func() {
			desc := selectors.GetCategoryDescription("consulting")
			Expect(desc).To(ContainSubstring("Consulting"))
		})

		It("should return description for research", func() {
			desc := selectors.GetCategoryDescription("research")
			Expect(desc).To(ContainSubstring("Research"))
		})

		It("should return description for mentoring", func() {
			desc := selectors.GetCategoryDescription("mentoring")
			Expect(desc).To(ContainSubstring("Mentoring"))
		})

		It("should return description for communication", func() {
			desc := selectors.GetCategoryDescription("communication")
			Expect(desc).To(ContainSubstring("Communication"))
		})

		It("should return description for collaboration", func() {
			desc := selectors.GetCategoryDescription("collaboration")
			Expect(desc).To(ContainSubstring("collaboration"))
		})

		It("should return description for problem-solving", func() {
			desc := selectors.GetCategoryDescription("problem-solving")
			Expect(desc).To(ContainSubstring("problem-solving"))
		})

		It("should return description for project-management", func() {
			desc := selectors.GetCategoryDescription("project-management")
			Expect(desc).To(ContainSubstring("Project"))
		})

		It("should return description for architecture", func() {
			desc := selectors.GetCategoryDescription("architecture")
			Expect(desc).To(ContainSubstring("architecture"))
		})

		It("should handle case-insensitive lookup", func() {
			desc := selectors.GetCategoryDescription("TECHNICAL")
			Expect(desc).To(ContainSubstring("Technical skills"))
		})

		It("should return category name for unknown category", func() {
			desc := selectors.GetCategoryDescription("unknown")
			Expect(desc).To(Equal("unknown"))
		})
	})

	Describe("IsSelected", func() {
		It("should return true for selected category", func() {
			Expect(selector.SelectCategory("technical")).To(Succeed())
			Expect(selector.IsSelected("technical")).To(BeTrue())
		})

		It("should return false for unselected category", func() {
			Expect(selector.IsSelected("technical")).To(BeFalse())
		})

		It("should handle case-insensitive check", func() {
			Expect(selector.SelectCategory("technical")).To(Succeed())
			Expect(selector.IsSelected("Technical")).To(BeTrue())
			Expect(selector.IsSelected("TECHNICAL")).To(BeTrue())
		})
	})
})
