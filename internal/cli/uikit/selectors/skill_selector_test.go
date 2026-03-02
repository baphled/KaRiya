package selectors_test

import (
	"github.com/baphled/kariya/internal/cli/uikit/selectors"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	domain "github.com/baphled/kariya/internal/domain/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SkillSelector", func() {
	var (
		selector *selectors.SkillSelector
		skills   []*domain.Skill
	)

	BeforeEach(func() {
		skills = []*domain.Skill{
			fixtures.SkillWith("go", "Go", "backend", "expert"),
			fixtures.SkillWith("ruby", "Ruby", "backend", "advanced"),
			fixtures.SkillWith("python", "Python", "backend", "advanced"),
		}
		selector = selectors.NewSkillSelector(skills)
	})

	Describe("NewSkillSelector", func() {
		It("should create a selector with provided skills", func() {
			Expect(selector).NotTo(BeNil())
		})

		It("should have no selected skills initially", func() {
			Expect(selector.SelectedSkillIDs()).To(BeEmpty())
		})

		It("should handle empty skills list", func() {
			emptySelector := selectors.NewSkillSelector([]*domain.Skill{})
			Expect(emptySelector).NotTo(BeNil())
			Expect(emptySelector.AvailableSkills()).To(BeEmpty())
		})
	})

	Describe("SelectedSkillIDs", func() {
		It("should return empty when nothing selected", func() {
			Expect(selector.SelectedSkillIDs()).To(BeEmpty())
		})

		It("should return sorted skill IDs", func() {
			Expect(selector.SelectSkill("ruby")).To(Succeed())
			Expect(selector.SelectSkill("go")).To(Succeed())

			ids := selector.SelectedSkillIDs()
			Expect(ids).To(HaveLen(2))
			Expect(ids[0]).To(Equal("go"))
			Expect(ids[1]).To(Equal("ruby"))
		})
	})

	Describe("AvailableSkills", func() {
		It("should return all skills", func() {
			available := selector.AvailableSkills()
			Expect(available).To(HaveLen(3))
		})

		It("should return skills sorted by name", func() {
			available := selector.AvailableSkills()
			Expect(available[0].Name).To(Equal("Go"))
			Expect(available[1].Name).To(Equal("Python"))
			Expect(available[2].Name).To(Equal("Ruby"))
		})

		It("should return a copy not the original slice", func() {
			available := selector.AvailableSkills()
			available[0] = fixtures.SkillWith("modified", "Modified", "backend", "expert")

			original := selector.AvailableSkills()
			Expect(original[0].Name).To(Equal("Go"))
		})
	})

	Describe("SelectSkill", func() {
		It("should select a valid skill", func() {
			err := selector.SelectSkill("go")
			Expect(err).NotTo(HaveOccurred())
			Expect(selector.IsSelected("go")).To(BeTrue())
		})

		It("should return error for invalid skill ID", func() {
			err := selector.SelectSkill("nonexistent")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("not a valid skill ID"))
		})

		It("should return error for already selected skill", func() {
			Expect(selector.SelectSkill("go")).To(Succeed())

			err := selector.SelectSkill("go")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("already selected"))
		})
	})

	Describe("DeselectSkill", func() {
		It("should deselect a selected skill", func() {
			Expect(selector.SelectSkill("go")).To(Succeed())

			err := selector.DeselectSkill("go")
			Expect(err).NotTo(HaveOccurred())
			Expect(selector.IsSelected("go")).To(BeFalse())
		})

		It("should return error when deselecting unselected skill", func() {
			err := selector.DeselectSkill("go")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("not selected"))
		})
	})

	Describe("ToggleSkill", func() {
		It("should select an unselected skill", func() {
			err := selector.ToggleSkill("go")
			Expect(err).NotTo(HaveOccurred())
			Expect(selector.IsSelected("go")).To(BeTrue())
		})

		It("should deselect a selected skill", func() {
			Expect(selector.SelectSkill("go")).To(Succeed())

			err := selector.ToggleSkill("go")
			Expect(err).NotTo(HaveOccurred())
			Expect(selector.IsSelected("go")).To(BeFalse())
		})

		It("should return error for invalid skill ID", func() {
			err := selector.ToggleSkill("nonexistent")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("IsSelected", func() {
		It("should return true for selected skill", func() {
			Expect(selector.SelectSkill("ruby")).To(Succeed())
			Expect(selector.IsSelected("ruby")).To(BeTrue())
		})

		It("should return false for unselected skill", func() {
			Expect(selector.IsSelected("ruby")).To(BeFalse())
		})
	})

	Describe("FilterSkills", func() {
		It("should return all skills when prefix is empty", func() {
			filtered := selector.FilterSkills("")
			Expect(filtered).To(HaveLen(3))
		})

		It("should filter skills by name prefix", func() {
			filtered := selector.FilterSkills("go")
			Expect(filtered).To(HaveLen(1))
			Expect(filtered[0].Name).To(Equal("Go"))
		})

		It("should be case-insensitive", func() {
			filtered := selector.FilterSkills("RU")
			Expect(filtered).To(HaveLen(1))
			Expect(filtered[0].Name).To(Equal("Ruby"))
		})

		It("should return empty when no match", func() {
			filtered := selector.FilterSkills("xyz")
			Expect(filtered).To(BeEmpty())
		})

		It("should return sorted results", func() {
			filtered := selector.FilterSkills("p")
			Expect(filtered).To(HaveLen(1))
			Expect(filtered[0].Name).To(Equal("Python"))
		})
	})

	Describe("GetSkillByID", func() {
		It("should return skill for valid ID", func() {
			skill, err := selector.GetSkillByID("go")
			Expect(err).NotTo(HaveOccurred())
			Expect(skill.Name).To(Equal("Go"))
		})

		It("should return error for invalid ID", func() {
			skill, err := selector.GetSkillByID("nonexistent")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("not found"))
			Expect(skill).To(BeNil())
		})
	})

	Describe("Reset", func() {
		It("should clear all selected skills", func() {
			Expect(selector.SelectSkill("go")).To(Succeed())
			Expect(selector.SelectSkill("ruby")).To(Succeed())

			selector.Reset()

			Expect(selector.SelectedSkillIDs()).To(BeEmpty())
			Expect(selector.IsSelected("go")).To(BeFalse())
		})

		It("should work on empty selector", func() {
			selector.Reset()
			Expect(selector.SelectedSkillIDs()).To(BeEmpty())
		})
	})

	Describe("SetSelectedSkills", func() {
		It("should set selected skills from ID list", func() {
			selector.SetSelectedSkills([]string{"go", "ruby"})

			ids := selector.SelectedSkillIDs()
			Expect(ids).To(HaveLen(2))
			Expect(ids).To(ContainElement("go"))
			Expect(ids).To(ContainElement("ruby"))
		})

		It("should replace previous selections", func() {
			Expect(selector.SelectSkill("python")).To(Succeed())
			selector.SetSelectedSkills([]string{"go"})

			ids := selector.SelectedSkillIDs()
			Expect(ids).To(HaveLen(1))
			Expect(ids).To(ContainElement("go"))
			Expect(ids).NotTo(ContainElement("python"))
		})

		It("should ignore invalid skill IDs", func() {
			selector.SetSelectedSkills([]string{"go", "nonexistent", "ruby"})

			ids := selector.SelectedSkillIDs()
			Expect(ids).To(HaveLen(2))
			Expect(ids).To(ContainElement("go"))
			Expect(ids).To(ContainElement("ruby"))
		})

		It("should handle empty list", func() {
			Expect(selector.SelectSkill("go")).To(Succeed())
			selector.SetSelectedSkills([]string{})

			Expect(selector.SelectedSkillIDs()).To(BeEmpty())
		})
	})
})
