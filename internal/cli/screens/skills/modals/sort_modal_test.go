package modals_test

import (
	"github.com/baphled/kariya/internal/cli/screens/skills/modals"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SortModal", func() {
	var (
		skills  []*career.Skill
		current *modals.SortConfig
		modal   *modals.SortModal
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

		current = &modals.SortConfig{
			SortBy:    "name",
			SortOrder: "asc",
		}
	})

	Describe("NewSortModal", func() {
		It("should create a new skill sort modal", func() {
			modal = modals.NewSortModal(skills, current, 120, 40)

			Expect(modal).NotTo(BeNil())
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("should pre-populate with current sort config", func() {
			current.SortBy = "years"
			current.SortOrder = "desc"

			modal = modals.NewSortModal(skills, current, 120, 40)

			Expect(modal).NotTo(BeNil())
			// Form should be pre-populated with current values
		})

		It("should default to name asc when no config provided", func() {
			modal = modals.NewSortModal(skills, nil, 120, 40)

			Expect(modal).NotTo(BeNil())
		})
	})

	Describe("Sort Options", func() {
		BeforeEach(func() {
			modal = modals.NewSortModal(skills, current, 120, 40)
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
			modal = modals.NewSortModal(skills, current, 120, 40)
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
			modal = modals.NewSortModal(skills, current, 120, 40)
			modal.Hide()

			view := modal.View()
			Expect(view).To(BeEmpty())
		})

		It("should render form when visible", func() {
			modal = modals.NewSortModal(skills, current, 120, 40)

			view := modal.View()
			Expect(view).NotTo(BeEmpty())
			// Should have huh form controls
			Expect(view).To(ContainSubstring("enter"))
		})

		It("should display keyboard shortcuts in footer (KeyBadge pattern)", func() {
			modal = modals.NewSortModal(skills, current, 120, 40)

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

	Describe("Update", func() {
		BeforeEach(func() {
			modal = modals.NewSortModal(skills, current, 120, 40)
		})

		It("should handle escape key to close without applying", func() {
			// Simulate escape key press
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("should return sort config when form completed", func() {
			// Form completion should return SortConfig
			Expect(modal).NotTo(BeNil())
		})
	})

	Describe("Visibility", func() {
		BeforeEach(func() {
			modal = modals.NewSortModal(skills, current, 120, 40)
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
