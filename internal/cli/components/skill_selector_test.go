package components_test

import (
	"testing"

	"github.com/baphled/kariya/internal/cli/components"
	domain "github.com/baphled/kariya/internal/domain/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSkillSelector(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SkillSelector Suite")
}

var _ = Describe("SkillSelector", func() {
	var (
		selector        *components.SkillSelector
		availableSkills []*domain.Skill
	)

	BeforeEach(func() {
		availableSkills = []*domain.Skill{
			{ID: "sk1", Name: "Go", Category: "backend"},
			{ID: "sk2", Name: "React", Category: "frontend"},
			{ID: "sk3", Name: "PostgreSQL", Category: "database"},
			{ID: "sk4", Name: "Kubernetes", Category: "devops"},
			{ID: "sk5", Name: "Docker", Category: "devops"},
		}
		selector = components.NewSkillSelector(availableSkills)
	})

	Describe("NewSkillSelector", func() {
		It("should create a new skill selector with empty selection", func() {
			Expect(selector).NotTo(BeNil())
			Expect(selector.SelectedSkillIDs()).To(BeEmpty())
		})
	})

	Describe("SelectSkill", func() {
		It("should select a skill by ID", func() {
			err := selector.SelectSkill("sk1")
			Expect(err).NotTo(HaveOccurred())
			Expect(selector.SelectedSkillIDs()).To(ContainElement("sk1"))
		})

		It("should return error for invalid skill ID", func() {
			err := selector.SelectSkill("invalid")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("not a valid skill"))
		})

		It("should return error when selecting already selected skill", func() {
			_ = selector.SelectSkill("sk1")
			err := selector.SelectSkill("sk1")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("already selected"))
		})

		It("should allow selecting multiple skills", func() {
			_ = selector.SelectSkill("sk1")
			_ = selector.SelectSkill("sk2")
			_ = selector.SelectSkill("sk3")
			Expect(selector.SelectedSkillIDs()).To(HaveLen(3))
			Expect(selector.SelectedSkillIDs()).To(ContainElements("sk1", "sk2", "sk3"))
		})
	})

	Describe("DeselectSkill", func() {
		BeforeEach(func() {
			_ = selector.SelectSkill("sk1")
			_ = selector.SelectSkill("sk2")
		})

		It("should deselect a selected skill", func() {
			err := selector.DeselectSkill("sk1")
			Expect(err).NotTo(HaveOccurred())
			Expect(selector.SelectedSkillIDs()).NotTo(ContainElement("sk1"))
			Expect(selector.SelectedSkillIDs()).To(ContainElement("sk2"))
		})

		It("should return error when deselecting unselected skill", func() {
			err := selector.DeselectSkill("sk3")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("not selected"))
		})
	})

	Describe("ToggleSkill", func() {
		It("should select skill if not selected", func() {
			err := selector.ToggleSkill("sk1")
			Expect(err).NotTo(HaveOccurred())
			Expect(selector.IsSelected("sk1")).To(BeTrue())
		})

		It("should deselect skill if already selected", func() {
			_ = selector.SelectSkill("sk1")
			err := selector.ToggleSkill("sk1")
			Expect(err).NotTo(HaveOccurred())
			Expect(selector.IsSelected("sk1")).To(BeFalse())
		})
	})

	Describe("AvailableSkills", func() {
		It("should return all available skills sorted by name", func() {
			skills := selector.AvailableSkills()
			Expect(skills).To(HaveLen(5))

			// Check alphabetical order
			Expect(skills[0].Name).To(Equal("Docker"))
			Expect(skills[1].Name).To(Equal("Go"))
			Expect(skills[2].Name).To(Equal("Kubernetes"))
			Expect(skills[3].Name).To(Equal("PostgreSQL"))
			Expect(skills[4].Name).To(Equal("React"))
		})
	})

	Describe("FilterSkills", func() {
		It("should return all skills when prefix is empty", func() {
			filtered := selector.FilterSkills("")
			Expect(filtered).To(HaveLen(5))
		})

		It("should filter skills by prefix (case-insensitive)", func() {
			filtered := selector.FilterSkills("go")
			Expect(filtered).To(HaveLen(1))
			Expect(filtered[0].Name).To(Equal("Go"))
		})

		It("should return multiple matches", func() {
			filtered := selector.FilterSkills("k")
			Expect(filtered).To(HaveLen(1))
			Expect(filtered[0].Name).To(Equal("Kubernetes"))
		})

		It("should return empty list for no matches", func() {
			filtered := selector.FilterSkills("xyz")
			Expect(filtered).To(BeEmpty())
		})
	})

	Describe("GetSkillByID", func() {
		It("should return skill for valid ID", func() {
			skill, err := selector.GetSkillByID("sk1")
			Expect(err).NotTo(HaveOccurred())
			Expect(skill.Name).To(Equal("Go"))
		})

		It("should return error for invalid ID", func() {
			_, err := selector.GetSkillByID("invalid")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("SetSelectedSkills", func() {
		It("should set selected skills from ID list", func() {
			selector.SetSelectedSkills([]string{"sk1", "sk3", "sk5"})
			Expect(selector.SelectedSkillIDs()).To(HaveLen(3))
			Expect(selector.SelectedSkillIDs()).To(ContainElements("sk1", "sk3", "sk5"))
		})

		It("should ignore invalid skill IDs", func() {
			selector.SetSelectedSkills([]string{"sk1", "invalid", "sk2"})
			Expect(selector.SelectedSkillIDs()).To(HaveLen(2))
			Expect(selector.SelectedSkillIDs()).To(ContainElements("sk1", "sk2"))
		})

		It("should clear previous selections", func() {
			_ = selector.SelectSkill("sk1")
			selector.SetSelectedSkills([]string{"sk2", "sk3"})
			Expect(selector.SelectedSkillIDs()).To(HaveLen(2))
			Expect(selector.SelectedSkillIDs()).NotTo(ContainElement("sk1"))
		})
	})

	Describe("Reset", func() {
		It("should clear all selections", func() {
			_ = selector.SelectSkill("sk1")
			_ = selector.SelectSkill("sk2")
			selector.Reset()
			Expect(selector.SelectedSkillIDs()).To(BeEmpty())
		})
	})
})
