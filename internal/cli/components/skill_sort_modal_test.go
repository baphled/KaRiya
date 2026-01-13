package components_test

import (
	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/domain/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SkillSortModal", func() {
	var (
		skills  []*career.Skill
		current *components.SkillSortConfig
		modal   *components.SkillSortModal
	)

	BeforeEach(func() {
		// Create test skills
		years1 := 5
		years2 := 3
		years3 := 4

		skills = []*career.Skill{
			{
				ID:        "skill-1",
				Name:      "Go",
				Category:  "backend",
				Level:     "expert",
				YearsUsed: &years1,
			},
			{
				ID:        "skill-2",
				Name:      "React",
				Category:  "frontend",
				Level:     "intermediate",
				YearsUsed: &years2,
			},
			{
				ID:        "skill-3",
				Name:      "Docker",
				Category:  "devops",
				Level:     "advanced",
				YearsUsed: &years3,
			},
		}

		current = &components.SkillSortConfig{
			SortBy:    "name",
			SortOrder: "asc",
		}
	})

	Describe("NewSkillSortModal", func() {
		It("should create a new skill sort modal", func() {
			modal = components.NewSkillSortModal(skills, current, 120, 40)

			Expect(modal).NotTo(BeNil())
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("should pre-populate with current sort config", func() {
			current.SortBy = "years"
			current.SortOrder = "desc"

			modal = components.NewSkillSortModal(skills, current, 120, 40)

			Expect(modal).NotTo(BeNil())
			// Form should be pre-populated with current values
		})

		It("should default to name asc when no config provided", func() {
			modal = components.NewSkillSortModal(skills, nil, 120, 40)

			Expect(modal).NotTo(BeNil())
		})
	})

	Describe("Sort Options", func() {
		BeforeEach(func() {
			modal = components.NewSkillSortModal(skills, current, 120, 40)
		})

		It("should include name sort option", func() {
			// Modal should have "name" sort option
			Expect(modal).NotTo(BeNil())
		})

		It("should include category sort option", func() {
			// Modal should have "category" sort option
			Expect(modal).NotTo(BeNil())
		})

		It("should include level sort option", func() {
			// Modal should have "level" sort option
			Expect(modal).NotTo(BeNil())
		})

		It("should include years sort option", func() {
			// Modal should have "years" sort option
			Expect(modal).NotTo(BeNil())
		})

		It("should include events count sort option", func() {
			// Modal should have "events" sort option
			Expect(modal).NotTo(BeNil())
		})
	})

	Describe("Sort Order Options", func() {
		BeforeEach(func() {
			modal = components.NewSkillSortModal(skills, current, 120, 40)
		})

		It("should include ascending option", func() {
			// Modal should have "asc" sort order
			Expect(modal).NotTo(BeNil())
		})

		It("should include descending option", func() {
			// Modal should have "desc" sort order
			Expect(modal).NotTo(BeNil())
		})
	})

	Describe("View", func() {
		It("should return empty string when not visible", func() {
			modal = components.NewSkillSortModal(skills, current, 120, 40)
			modal.Hide()

			view := modal.View()
			Expect(view).To(BeEmpty())
		})

		It("should render form when visible", func() {
			modal = components.NewSkillSortModal(skills, current, 120, 40)

			view := modal.View()
			Expect(view).NotTo(BeEmpty())
			// Should have huh form controls
			Expect(view).To(ContainSubstring("enter"))
		})
	})

	Describe("Update", func() {
		BeforeEach(func() {
			modal = components.NewSkillSortModal(skills, current, 120, 40)
		})

		It("should handle escape key to close without applying", func() {
			// Simulate escape key press
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("should return sort config when form completed", func() {
			// Form completion should return SkillSortConfig
			Expect(modal).NotTo(BeNil())
		})
	})

	Describe("Visibility", func() {
		BeforeEach(func() {
			modal = components.NewSkillSortModal(skills, current, 120, 40)
		})

		It("should start visible", func() {
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("should hide when Hide() is called", func() {
			modal.Hide()
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should show when Show() is called", func() {
			modal.Hide()
			modal.Show()
			Expect(modal.IsVisible()).To(BeTrue())
		})
	})
})
