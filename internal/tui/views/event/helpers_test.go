package event_test

import (
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	eventview "github.com/baphled/kariya/internal/tui/views/event"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/themes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Helpers", func() {
	var theme themes.Theme

	BeforeEach(func() {
		theme = themes.NewDefaultTheme()
	})

	Describe("RenderOverlayModal", func() {
		It("renders modal over background", func() {
			result := eventview.RenderOverlayModal("Modal Content", "Background Content")

			Expect(result).NotTo(BeEmpty())
			Expect(result).To(ContainSubstring("Modal Content"))
		})

		It("handles empty modal view", func() {
			result := eventview.RenderOverlayModal("", "Background")
			Expect(result).NotTo(BeEmpty())
		})

		It("handles empty background view", func() {
			result := eventview.RenderOverlayModal("Modal", "")
			Expect(result).NotTo(BeEmpty())
			Expect(result).To(ContainSubstring("Modal"))
		})
	})

	Describe("RenderEventDetailContent", func() {
		It("renders zero-value event as no event selected", func() {
			result := eventview.RenderEventDetailContent(display.Event{}, theme)
			Expect(result).To(Equal("No event selected."))
		})

		It("renders event with all fields", func() {
			domainEvent := fixtures.EventWith("test-id", "Test event description", "Test Company", "Test Project")
			domainEvent.Date = time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
			domainEvent.Tags = []string{"tag1", "tag2"}
			domainEvent.Categories = []string{"category1"}
			domainEvent.Skills = []string{"skill1", "skill2", "skill3"}

			result := eventview.RenderEventDetailContent(display.EventFromDomain(domainEvent), theme)

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
			domainEvent := fixtures.Event("test-id")
			domainEvent.Text = "Minimal event"
			domainEvent.Date = time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

			result := eventview.RenderEventDetailContent(display.EventFromDomain(domainEvent), theme)

			Expect(result).To(ContainSubstring("Event Details"))
			Expect(result).To(ContainSubstring("2024-01-15"))
			Expect(result).To(ContainSubstring("Minimal event"))
			Expect(result).NotTo(ContainSubstring("Company:"))
			Expect(result).NotTo(ContainSubstring("Project:"))
		})

		It("omits empty company and project", func() {
			domainEvent := fixtures.EventWith("test-id", "Event text", "", "")

			result := eventview.RenderEventDetailContent(display.EventFromDomain(domainEvent), theme)

			Expect(result).NotTo(ContainSubstring("Company:"))
			Expect(result).NotTo(ContainSubstring("Project:"))
		})
	})

	Describe("RenderSkillsContent", func() {
		It("renders empty skills message", func() {
			result := eventview.RenderSkillsContent([]display.Skill{}, theme)
			Expect(result).To(ContainSubstring("No skills associated"))
		})

		It("renders nil skills slice as empty", func() {
			result := eventview.RenderSkillsContent(nil, theme)
			Expect(result).To(ContainSubstring("No skills associated"))
		})

		It("renders single skill", func() {
			years := 3
			domainSkill := fixtures.SkillWith("skill-1", "Go", "Programming", "Expert")
			domainSkill.YearsUsed = &years
			skills := display.SkillsFromDomain([]*career.Skill{domainSkill})

			result := eventview.RenderSkillsContent(skills, theme)

			Expect(result).To(ContainSubstring("Go"))
			Expect(result).To(ContainSubstring("Programming"))
			Expect(result).To(ContainSubstring("Expert"))
			Expect(result).To(ContainSubstring("3 years"))
		})

		It("renders single year correctly", func() {
			years := 1
			domainSkill := fixtures.Skill("skill-1")
			domainSkill.Name = "Rust"
			domainSkill.YearsUsed = &years
			skills := display.SkillsFromDomain([]*career.Skill{domainSkill})

			result := eventview.RenderSkillsContent(skills, theme)

			Expect(result).To(ContainSubstring("1 year"))
			Expect(result).NotTo(ContainSubstring("1 years"))
		})

		It("renders multiple skills", func() {
			domainSkills := []*career.Skill{
				fixtures.SkillWith("1", "Go", "Backend", "intermediate"),
				fixtures.SkillWith("2", "React", "Frontend", "intermediate"),
				fixtures.SkillWith("3", "PostgreSQL", "Database", "intermediate"),
			}

			result := eventview.RenderSkillsContent(display.SkillsFromDomain(domainSkills), theme)

			Expect(result).To(ContainSubstring("Go"))
			Expect(result).To(ContainSubstring("React"))
			Expect(result).To(ContainSubstring("PostgreSQL"))
		})

		It("skips zero-value skills in slice", func() {
			s1 := display.SkillFromDomain(fixtures.SkillWith("1", "Go", "Backend", "intermediate"))
			s2 := display.SkillFromDomain(fixtures.SkillWith("2", "Python", "Backend", "intermediate"))
			skills := []display.Skill{s1, {}, s2}

			result := eventview.RenderSkillsContent(skills, theme)

			Expect(result).To(ContainSubstring("Go"))
			Expect(result).To(ContainSubstring("Python"))
		})

		It("handles skill with nil years", func() {
			domainSkill := fixtures.SkillWith("skill-1", "Docker", "DevOps", "Intermediate")
			domainSkill.YearsUsed = nil
			skills := display.SkillsFromDomain([]*career.Skill{domainSkill})

			result := eventview.RenderSkillsContent(skills, theme)

			Expect(result).To(ContainSubstring("Docker"))
			Expect(result).NotTo(ContainSubstring("year"))
		})

		It("handles skill with zero years", func() {
			years := 0
			domainSkill := fixtures.Skill("skill-1")
			domainSkill.Name = "Kubernetes"
			domainSkill.YearsUsed = &years
			skills := display.SkillsFromDomain([]*career.Skill{domainSkill})

			result := eventview.RenderSkillsContent(skills, theme)

			Expect(result).To(ContainSubstring("Kubernetes"))
			Expect(result).NotTo(ContainSubstring("0 year"))
		})
	})
})
