package selectors_test

import (
	domain "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/ui/uikit/selectors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SkillSelector", func() {
	var selector *selectors.SkillSelector

	BeforeEach(func() {
		skillData := []struct {
			ID   string
			Name string
		}{
			{ID: "s1", Name: "Go"},
			{ID: "s2", Name: "Python"},
			{ID: "s3", Name: "Rust"},
			{ID: "s4", Name: "JavaScript"},
			{ID: "s5", Name: "TypeScript"},
		}

		domainSkills := make([]*domain.Skill, len(skillData))
		for i, s := range skillData {
			domainSkills[i] = fixtures.SkillWith(s.ID, s.Name, "backend", "intermediate")
		}
		selector = selectors.NewSkillSelector(domainSkills)
	})

	Describe("NewSkillSelector", func() {
		It("should create a selector with available skills", func() {
			Expect(selector).NotTo(BeNil())
			Expect(selector.SelectedSkillIDs()).To(BeEmpty())
		})

		It("should initialize with empty selection", func() {
			Expect(selector.SelectedSkillIDs()).To(BeEmpty())
		})

		It("should have all skills available", func() {
			available := selector.AvailableSkills()
			Expect(available).To(HaveLen(5))
		})
	})

	Describe("SelectedSkillIDs", func() {
		It("should return empty when no skills selected", func() {
			Expect(selector.SelectedSkillIDs()).To(BeEmpty())
		})

		It("should return selected skill IDs", func() {
			Expect(selector.SelectSkill("s1")).To(Succeed())
			Expect(selector.SelectSkill("s3")).To(Succeed())
			selected := selector.SelectedSkillIDs()
			Expect(selected).To(HaveLen(2))
			Expect(selected).To(ContainElement("s1"))
			Expect(selected).To(ContainElement("s3"))
		})

		It("should return sorted IDs", func() {
			Expect(selector.SelectSkill("s3")).To(Succeed())
			Expect(selector.SelectSkill("s1")).To(Succeed())
			Expect(selector.SelectSkill("s2")).To(Succeed())
			selected := selector.SelectedSkillIDs()
			Expect(selected).To(Equal([]string{"s1", "s2", "s3"}))
		})
	})

	Describe("AvailableSkills", func() {
		It("should return all available skills", func() {
			available := selector.AvailableSkills()
			Expect(available).To(HaveLen(5))
		})

		It("should return skills sorted by name", func() {
			available := selector.AvailableSkills()
			Expect(available[0].Name).To(Equal("Go"))
			Expect(available[1].Name).To(Equal("JavaScript"))
			Expect(available[2].Name).To(Equal("Python"))
			Expect(available[3].Name).To(Equal("Rust"))
			Expect(available[4].Name).To(Equal("TypeScript"))
		})

		It("should return a copy not affecting original", func() {
			available := selector.AvailableSkills()
			Expect(available).To(HaveLen(5))
			Expect(selector.AvailableSkills()).To(HaveLen(5))
		})
	})

	Describe("SelectSkill", func() {
		It("should select a valid skill", func() {
			err := selector.SelectSkill("s1")
			Expect(err).NotTo(HaveOccurred())
			Expect(selector.IsSelected("s1")).To(BeTrue())
		})

		It("should return error for invalid skill ID", func() {
			err := selector.SelectSkill("invalid")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("not a valid skill ID"))
		})

		It("should return error for duplicate selection", func() {
			Expect(selector.SelectSkill("s1")).To(Succeed())
			err := selector.SelectSkill("s1")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("already selected"))
		})

		It("should allow selecting multiple skills", func() {
			Expect(selector.SelectSkill("s1")).To(Succeed())
			Expect(selector.SelectSkill("s2")).To(Succeed())
			Expect(selector.SelectSkill("s3")).To(Succeed())
			selected := selector.SelectedSkillIDs()
			Expect(selected).To(HaveLen(3))
		})
	})

	Describe("DeselectSkill", func() {
		BeforeEach(func() {
			Expect(selector.SelectSkill("s1")).To(Succeed())
			Expect(selector.SelectSkill("s2")).To(Succeed())
		})

		It("should deselect a selected skill", func() {
			err := selector.DeselectSkill("s1")
			Expect(err).NotTo(HaveOccurred())
			Expect(selector.IsSelected("s1")).To(BeFalse())
		})

		It("should return error when deselecting unselected skill", func() {
			err := selector.DeselectSkill("s3")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("not selected"))
		})

		It("should maintain other selections", func() {
			Expect(selector.DeselectSkill("s1")).To(Succeed())
			Expect(selector.IsSelected("s2")).To(BeTrue())
		})
	})

	Describe("ToggleSkill", func() {
		It("should select unselected skill", func() {
			err := selector.ToggleSkill("s1")
			Expect(err).NotTo(HaveOccurred())
			Expect(selector.IsSelected("s1")).To(BeTrue())
		})

		It("should deselect selected skill", func() {
			Expect(selector.SelectSkill("s1")).To(Succeed())
			err := selector.ToggleSkill("s1")
			Expect(err).NotTo(HaveOccurred())
			Expect(selector.IsSelected("s1")).To(BeFalse())
		})

		It("should return error for invalid skill ID", func() {
			err := selector.ToggleSkill("invalid")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("IsSelected", func() {
		It("should return true for selected skill", func() {
			Expect(selector.SelectSkill("s1")).To(Succeed())
			Expect(selector.IsSelected("s1")).To(BeTrue())
		})

		It("should return false for unselected skill", func() {
			Expect(selector.IsSelected("s1")).To(BeFalse())
		})

		It("should return false for invalid skill ID", func() {
			Expect(selector.IsSelected("invalid")).To(BeFalse())
		})
	})

	Describe("FilterSkills", func() {
		It("should return all skills when prefix is empty", func() {
			filtered := selector.FilterSkills("")
			Expect(filtered).To(HaveLen(5))
		})

		It("should return all skills sorted when prefix is empty string", func() {
			filtered := selector.FilterSkills("")
			Expect(filtered).To(HaveLen(5))
			Expect(filtered[0].Name).To(Equal("Go"))
			Expect(filtered[1].Name).To(Equal("JavaScript"))
			Expect(filtered[2].Name).To(Equal("Python"))
			Expect(filtered[3].Name).To(Equal("Rust"))
			Expect(filtered[4].Name).To(Equal("TypeScript"))
		})

		It("should filter skills by name prefix", func() {
			filtered := selector.FilterSkills("Py")
			Expect(filtered).To(HaveLen(1))
			Expect(filtered[0].Name).To(Equal("Python"))
		})

		It("should be case-insensitive", func() {
			filtered := selector.FilterSkills("go")
			Expect(filtered).To(HaveLen(1))
			Expect(filtered[0].Name).To(Equal("Go"))
		})

		It("should return empty when no match", func() {
			filtered := selector.FilterSkills("xyz")
			Expect(filtered).To(BeEmpty())
		})

		It("should return sorted results", func() {
			filtered := selector.FilterSkills("T")
			Expect(filtered).To(HaveLen(1))
			Expect(filtered[0].Name).To(Equal("TypeScript"))
		})

		It("should match single skill", func() {
			filtered := selector.FilterSkills("R")
			Expect(filtered).To(HaveLen(1))
			Expect(filtered[0].Name).To(Equal("Rust"))
		})

		It("should filter with partial prefix", func() {
			filtered := selector.FilterSkills("J")
			Expect(filtered).To(HaveLen(1))
			Expect(filtered[0].Name).To(Equal("JavaScript"))
		})
	})

	Describe("GetSkillByID", func() {
		It("should return skill for valid ID", func() {
			skill, err := selector.GetSkillByID("s1")
			Expect(err).NotTo(HaveOccurred())
			Expect(skill).NotTo(BeNil())
			Expect(skill.Name).To(Equal("Go"))
		})

		It("should return error for invalid ID", func() {
			skill, err := selector.GetSkillByID("invalid")
			Expect(err).To(HaveOccurred())
			Expect(skill).To(BeNil())
			Expect(err.Error()).To(ContainSubstring("not found"))
		})

		It("should return correct skill for each ID", func() {
			skill, err := selector.GetSkillByID("s2")
			Expect(err).NotTo(HaveOccurred())
			Expect(skill.Name).To(Equal("Python"))
		})
	})

	Describe("Reset", func() {
		BeforeEach(func() {
			Expect(selector.SelectSkill("s1")).To(Succeed())
			Expect(selector.SelectSkill("s2")).To(Succeed())
			Expect(selector.SelectSkill("s3")).To(Succeed())
		})

		It("should clear all selected skills", func() {
			selector.Reset()
			Expect(selector.SelectedSkillIDs()).To(BeEmpty())
		})

		It("should allow reselecting after reset", func() {
			selector.Reset()
			err := selector.SelectSkill("s1")
			Expect(err).NotTo(HaveOccurred())
			Expect(selector.IsSelected("s1")).To(BeTrue())
		})
	})

	Describe("SetSelectedSkills", func() {
		It("should set selected skills from IDs", func() {
			selector.SetSelectedSkills([]string{"s1", "s3"})
			selected := selector.SelectedSkillIDs()
			Expect(selected).To(HaveLen(2))
			Expect(selected).To(ContainElement("s1"))
			Expect(selected).To(ContainElement("s3"))
		})

		It("should clear previous selections", func() {
			Expect(selector.SelectSkill("s1")).To(Succeed())
			selector.SetSelectedSkills([]string{"s2", "s3"})
			selected := selector.SelectedSkillIDs()
			Expect(selected).To(HaveLen(2))
			Expect(selected).NotTo(ContainElement("s1"))
		})

		It("should ignore invalid skill IDs", func() {
			selector.SetSelectedSkills([]string{"s1", "invalid", "s2"})
			selected := selector.SelectedSkillIDs()
			Expect(selected).To(HaveLen(2))
			Expect(selected).To(ContainElement("s1"))
			Expect(selected).To(ContainElement("s2"))
		})

		It("should handle empty list", func() {
			Expect(selector.SelectSkill("s1")).To(Succeed())
			selector.SetSelectedSkills([]string{})
			Expect(selector.SelectedSkillIDs()).To(BeEmpty())
		})

		It("should handle all invalid IDs", func() {
			selector.SetSelectedSkills([]string{"invalid1", "invalid2"})
			Expect(selector.SelectedSkillIDs()).To(BeEmpty())
		})
	})
})
