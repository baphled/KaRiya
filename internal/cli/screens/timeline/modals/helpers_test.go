package modals_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/screens/timeline/modals"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
)

var _ = Describe("Helpers", func() {
	var theme themes.Theme

	BeforeEach(func() {
		theme = themes.NewDefaultTheme()
	})

	Describe("RenderOverlayModal", func() {
		It("renders modal over background", func() {
			modalView := "Modal Content"
			backgroundView := "Background Content"

			result := modals.RenderOverlayModal(modalView, backgroundView)

			Expect(result).NotTo(BeEmpty())
			Expect(result).To(ContainSubstring("Modal Content"))
		})

		It("handles empty modal view", func() {
			result := modals.RenderOverlayModal("", "Background")

			Expect(result).NotTo(BeEmpty())
		})

		It("handles empty background view", func() {
			result := modals.RenderOverlayModal("Modal", "")

			Expect(result).NotTo(BeEmpty())
			Expect(result).To(ContainSubstring("Modal"))
		})
	})

	Describe("RenderEventDetailContent", func() {
		It("renders nil event as no event selected", func() {
			result := modals.RenderEventDetailContent(nil, theme)

			Expect(result).To(Equal("No event selected."))
		})

		It("renders event with all fields", func() {
			event := fixtures.EventWith("test-id", "Test event description", "Test Company", "Test Project")
			event.Date = time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
			event.Tags = []string{"tag1", "tag2"}
			event.Categories = []string{"category1"}
			event.Skills = []string{"skill1", "skill2", "skill3"}

			result := modals.RenderEventDetailContent(event, theme)

			Expect(result).To(ContainSubstring("Event Details"))
			Expect(result).To(ContainSubstring("2024-01-15"))
			Expect(result).To(ContainSubstring("Test Company"))
			Expect(result).To(ContainSubstring("Test Project"))
			Expect(result).To(ContainSubstring("Test event description"))
			Expect(result).To(ContainSubstring("tag1, tag2"))
			Expect(result).To(ContainSubstring("category1"))
			Expect(result).To(ContainSubstring("3 associated"))
		})

		It("renders event with minimal fields", func() {
			event := fixtures.Event("test-id")
			event.Text = "Minimal event"
			event.Date = time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

			result := modals.RenderEventDetailContent(event, theme)

			Expect(result).To(ContainSubstring("Event Details"))
			Expect(result).To(ContainSubstring("2024-01-15"))
			Expect(result).To(ContainSubstring("Minimal event"))
			Expect(result).NotTo(ContainSubstring("Company:"))
			Expect(result).NotTo(ContainSubstring("Project:"))
		})

		It("omits empty company and project", func() {
			event := fixtures.EventWith("test-id", "Event text", "", "")

			result := modals.RenderEventDetailContent(event, theme)

			Expect(result).NotTo(ContainSubstring("Company:"))
			Expect(result).NotTo(ContainSubstring("Project:"))
		})
	})

	Describe("RenderSkillsContent", func() {
		It("renders empty skills message", func() {
			result := modals.RenderSkillsContent([]*career.Skill{}, theme)

			Expect(result).To(ContainSubstring("No skills associated"))
		})

		It("renders nil skills slice as empty", func() {
			result := modals.RenderSkillsContent(nil, theme)

			Expect(result).To(ContainSubstring("No skills associated"))
		})

		It("renders single skill", func() {
			skill := fixtures.SkillWith("skill-1", "Go", "Programming", "Expert")
			years := 3
			skill.YearsUsed = &years
			skills := []*career.Skill{skill}

			result := modals.RenderSkillsContent(skills, theme)

			Expect(result).To(ContainSubstring("Go"))
			Expect(result).To(ContainSubstring("Programming"))
			Expect(result).To(ContainSubstring("Expert"))
			Expect(result).To(ContainSubstring("3 years"))
		})

		It("renders single year correctly", func() {
			years := 1
			skill := fixtures.Skill("skill-1")
			skill.Name = "Rust"
			skill.YearsUsed = &years
			skills := []*career.Skill{skill}

			result := modals.RenderSkillsContent(skills, theme)

			Expect(result).To(ContainSubstring("1 year"))
			Expect(result).NotTo(ContainSubstring("1 years"))
		})

		It("renders multiple skills", func() {
			skills := []*career.Skill{
				fixtures.SkillWith("1", "Go", "Backend", "intermediate"),
				fixtures.SkillWith("2", "React", "Frontend", "intermediate"),
				fixtures.SkillWith("3", "PostgreSQL", "Database", "intermediate"),
			}

			result := modals.RenderSkillsContent(skills, theme)

			Expect(result).To(ContainSubstring("Go"))
			Expect(result).To(ContainSubstring("React"))
			Expect(result).To(ContainSubstring("PostgreSQL"))
		})

		It("skips nil skills in slice", func() {
			s1 := fixtures.Skill("1")
			s1.Name = "Go"
			s2 := fixtures.Skill("2")
			s2.Name = "Python"
			skills := []*career.Skill{s1, nil, s2}

			result := modals.RenderSkillsContent(skills, theme)

			Expect(result).To(ContainSubstring("Go"))
			Expect(result).To(ContainSubstring("Python"))
		})

		It("handles skill with nil years", func() {
			skill := fixtures.SkillWith("skill-1", "Docker", "DevOps", "Intermediate")
			skill.YearsUsed = nil
			skills := []*career.Skill{skill}

			result := modals.RenderSkillsContent(skills, theme)

			Expect(result).To(ContainSubstring("Docker"))
			Expect(result).NotTo(ContainSubstring("year"))
		})

		It("handles skill with zero years", func() {
			years := 0
			skill := fixtures.Skill("skill-1")
			skill.Name = "Kubernetes"
			skill.YearsUsed = &years
			skills := []*career.Skill{skill}

			result := modals.RenderSkillsContent(skills, theme)

			Expect(result).To(ContainSubstring("Kubernetes"))
			Expect(result).NotTo(ContainSubstring("0 year"))
		})
	})
})
