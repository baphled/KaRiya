package e2e_test

import (
	"fmt"

	"github.com/baphled/kariya/internal/constants"
	"github.com/baphled/kariya/internal/testutil/e2e"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("E2E SkillsManagement Workflow", func() {
	var env *e2e.TestEnv

	Describe("Empty Skill List", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnv(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show empty message when no skills exist", func() {
			env.AssertSkillCount(0)
			env.SelectIntentByName("manage_skills")
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
		})

		It("should not panic on empty skill list", func() {
			env.SelectIntentByName("manage_skills")
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should allow navigation back from empty list", func() {
			env.SelectIntentByName("manage_skills")
			env.Cancel()
			env.AssertViewContains("Capture Event")
		})
	})

	Describe("Skill Persistence with New Categories", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnv(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should persist skills with architecture category", func() {
			skill := e2e.CreateMinimalSkill("s-arch", "Microservices", string(constants.SkillCategoryArchitecture))
			env.AddSkill(skill)
			env.AssertSkillCount(1)

			skills := env.GetSkills()
			Expect(skills[0].Category).To(Equal(string(constants.SkillCategoryArchitecture)))
		})

		It("should persist skills with security category", func() {
			skill := e2e.CreateMinimalSkill("s-sec", "OAuth", string(constants.SkillCategorySecurity))
			env.AddSkill(skill)
			env.AssertSkillCount(1)

			skills := env.GetSkills()
			Expect(skills[0].Category).To(Equal(string(constants.SkillCategorySecurity)))
		})

		It("should persist skills with practices category", func() {
			skill := e2e.CreateMinimalSkill("s-prac", "Agile", string(constants.SkillCategoryPractices))
			env.AddSkill(skill)
			env.AssertSkillCount(1)

			skills := env.GetSkills()
			Expect(skills[0].Category).To(Equal(string(constants.SkillCategoryPractices)))
		})

		It("should persist skills across all 15 categories", func() {
			allCategories := constants.AllSkillCategories()
			Expect(allCategories).To(HaveLen(15))

			for i, cat := range allCategories {
				skill := e2e.CreateMinimalSkill(
					fmt.Sprintf("s-%02d", i),
					fmt.Sprintf("Skill_%s", string(cat)),
					string(cat),
				)
				env.AddSkill(skill)
			}

			env.AssertSkillCount(15)

			skills := env.GetSkills()
			persistedCategories := make(map[string]bool)
			for _, s := range skills {
				persistedCategories[s.Category] = true
			}

			for _, cat := range allCategories {
				Expect(persistedCategories).To(HaveKey(string(cat)),
					"category %q should be persisted", string(cat))
			}
		})
	})

	Describe("Sample Skills Fixture", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnv(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should create sample skills covering all categories", func() {
			allCategories := constants.AllSkillCategories()
			skills := e2e.CreateSampleSkills(len(allCategories))
			Expect(skills).To(HaveLen(len(allCategories)))

			coveredCategories := make(map[string]bool)
			for _, s := range skills {
				coveredCategories[s.Category] = true
			}

			for _, cat := range allCategories {
				Expect(coveredCategories).To(HaveKey(string(cat)),
					"CreateSampleSkills should cover category %q", string(cat))
			}
		})

		It("should persist sample skills to database", func() {
			allCategories := constants.AllSkillCategories()
			skills := e2e.CreateSampleSkills(len(allCategories))
			for _, s := range skills {
				env.AddSkill(s)
			}

			env.AssertSkillCount(len(allCategories))
		})
	})

	Describe("Skills Management Navigation", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnv(GinkgoT())
			skills := e2e.CreateSampleSkills(5)
			for _, s := range skills {
				env.AddSkill(s)
			}
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should navigate to manage skills intent", func() {
			env.SelectIntentByName("manage_skills")
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("Skills"))
		})

		It("should navigate down with j key", func() {
			env.SelectIntentByName("manage_skills")
			env.PressKeyRune('j')
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
		})

		It("should navigate with arrow keys", func() {
			env.SelectIntentByName("manage_skills")
			env.PressKey(tea.KeyDown)
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			env.PressKey(tea.KeyUp)
			view = env.GetView()
			Expect(view).NotTo(BeEmpty())
		})

		It("should return to menu on escape", func() {
			env.SelectIntentByName("manage_skills")
			env.Cancel()
			env.AssertViewContains("Capture Event")
		})
	})

	Describe("Category Validation via Constants", func() {
		It("should validate all new categories as valid", func() {
			Expect(constants.IsValidSkillCategory("architecture")).To(BeTrue())
			Expect(constants.IsValidSkillCategory("security")).To(BeTrue())
			Expect(constants.IsValidSkillCategory("practices")).To(BeTrue())
		})

		It("should validate all original categories as valid", func() {
			Expect(constants.IsValidSkillCategory("backend")).To(BeTrue())
			Expect(constants.IsValidSkillCategory("frontend")).To(BeTrue())
			Expect(constants.IsValidSkillCategory("devops")).To(BeTrue())
			Expect(constants.IsValidSkillCategory("database")).To(BeTrue())
			Expect(constants.IsValidSkillCategory("cloud")).To(BeTrue())
			Expect(constants.IsValidSkillCategory("mobile")).To(BeTrue())
			Expect(constants.IsValidSkillCategory("tooling")).To(BeTrue())
			Expect(constants.IsValidSkillCategory("testing")).To(BeTrue())
			Expect(constants.IsValidSkillCategory("data")).To(BeTrue())
			Expect(constants.IsValidSkillCategory("ml")).To(BeTrue())
			Expect(constants.IsValidSkillCategory("monitoring")).To(BeTrue())
			Expect(constants.IsValidSkillCategory("other")).To(BeTrue())
		})

		It("should reject invalid categories", func() {
			Expect(constants.IsValidSkillCategory("invalid")).To(BeFalse())
			Expect(constants.IsValidSkillCategory("")).To(BeFalse())
		})

		It("should produce string representations for all 15 categories", func() {
			strs := constants.SkillCategoryStrings()
			Expect(strs).To(HaveLen(15))
			Expect(strs).To(ContainElements(
				"architecture", "security", "practices",
			))
		})
	})
})
