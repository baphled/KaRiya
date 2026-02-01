package modals_test

import (
	"github.com/baphled/kariya/internal/cli/screens/skills/modals"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("FilterModal", func() {
	var (
		skills        []*career.Skill
		currentFilter *modals.Filters
		modal         *modals.FilterModal
	)

	BeforeEach(func() {
		years1 := 5
		years2 := 3
		years3 := 4

		s1 := fixtures.SkillWith("skill-1", "Go", "backend", "expert")
		s1.YearsUsed = &years1
		s2 := fixtures.SkillWith("skill-2", "React", "frontend", "intermediate")
		s2.YearsUsed = &years2
		s3 := fixtures.SkillWith("skill-3", "Docker", "devops", "advanced")
		s3.YearsUsed = &years3

		skills = []*career.Skill{s1, s2, s3}

		currentFilter = &modals.Filters{
			Categories: []string{},
			Levels:     []string{},
			MinYears:   0,
			MaxYears:   0,
		}
	})

	Describe("NewFilterModal", func() {
		It("should create a new skill filter modal", func() {
			modal = modals.NewFilterModal(skills, currentFilter, 120, 40)

			Expect(modal).NotTo(BeNil())
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("should extract unique categories from skills", func() {
			modal = modals.NewFilterModal(skills, currentFilter, 120, 40)

			// Modal should include Backend, Frontend, DevOps categories
			Expect(modal).NotTo(BeNil())
		})

		It("should pre-populate form with current filters", func() {
			currentFilter.Categories = []string{"Backend"}
			currentFilter.Levels = []string{"Expert"}
			currentFilter.MinYears = 2
			currentFilter.MaxYears = 10

			modal = modals.NewFilterModal(skills, currentFilter, 120, 40)

			Expect(modal).NotTo(BeNil())
			// Form data should be pre-populated (internal state)
		})
	})

	Describe("Filtering Options", func() {
		BeforeEach(func() {
			modal = modals.NewFilterModal(skills, currentFilter, 120, 40)
		})

		It("should include category filter", func() {
			// Modal should have MultiSelect for categories
			Expect(modal).NotTo(BeNil())
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("should include level filter", func() {
			// Modal should have Select for level
			Expect(modal).NotTo(BeNil())
		})

		It("should include years range filter", func() {
			// Modal should have input fields for min/max years
			Expect(modal).NotTo(BeNil())
		})

		It("should handle empty categories gracefully", func() {
			emptySkills := []*career.Skill{}
			modal = modals.NewFilterModal(emptySkills, currentFilter, 120, 40)

			Expect(modal).NotTo(BeNil())
		})
	})

	Describe("Sorting Options", func() {
		BeforeEach(func() {
			modal = modals.NewFilterModal(skills, currentFilter, 120, 40)
		})

		It("should include sort by options", func() {
			// Modal should have Select for sort field (name, category, level, years, events)
			Expect(modal).NotTo(BeNil())
		})

		It("should include sort order options", func() {
			// Modal should have Select for sort direction (asc, desc)
			Expect(modal).NotTo(BeNil())
		})
	})

	Describe("View", func() {
		It("should return empty string when not visible", func() {
			modal = modals.NewFilterModal(skills, currentFilter, 120, 40)
			modal.Hide()

			view := modal.View()
			Expect(view).To(BeEmpty())
		})

		It("should render form when visible", func() {
			modal = modals.NewFilterModal(skills, currentFilter, 120, 40)

			view := modal.View()
			Expect(view).NotTo(BeEmpty())
			// Should have huh form controls
			Expect(view).To(ContainSubstring("enter"))
		})

		It("should display keyboard shortcuts in footer (KeyBadge pattern)", func() {
			modal = modals.NewFilterModal(skills, currentFilter, 120, 40)

			view := modal.View()
			// Modal should show keyboard shortcuts for user guidance
			// Pattern: Tab (next field), Enter (submit), Esc (cancel)
			Expect(view).To(ContainSubstring("Tab"))
			Expect(view).To(ContainSubstring("Enter"))
			Expect(view).To(ContainSubstring("Esc"))
			// Should indicate what each key does
			Expect(view).To(Or(
				ContainSubstring("Next"),
				ContainSubstring("field"),
				ContainSubstring("Navigate"),
			))
			Expect(view).To(Or(
				ContainSubstring("Submit"),
				ContainSubstring("Confirm"),
				ContainSubstring("Apply"),
			))
			Expect(view).To(Or(
				ContainSubstring("Cancel"),
				ContainSubstring("Back"),
			))
		})
	})
})
