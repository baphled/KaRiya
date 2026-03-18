package career

import (
	"time"

	"github.com/baphled/kariya/internal/testutil/fixtures"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Skill Model", func() {
	Describe("TableName", func() {
		It("returns skills", func() {
			Expect(Skill{}.TableName()).To(Equal("skills"))
		})
	})

	Describe("ToDomain", func() {
		It("converts all fields to domain Skill", func() {
			now := time.Now()
			years := 5
			lastUsed := now.AddDate(0, -1, 0)
			model := &Skill{
				ID:        "skill-1",
				Name:      "Go",
				Category:  "backend",
				Level:     "expert",
				YearsUsed: &years,
				LastUsed:  &lastUsed,
				CreatedAt: now,
				UpdatedAt: now,
			}

			domain := model.ToDomain()

			Expect(domain.ID).To(Equal("skill-1"))
			Expect(domain.Name).To(Equal("Go"))
			Expect(domain.Category).To(Equal("backend"))
			Expect(domain.Level).To(Equal("expert"))
			Expect(*domain.YearsUsed).To(Equal(5))
			Expect(*domain.LastUsed).To(Equal(lastUsed))
			Expect(domain.CreatedAt).To(Equal(now))
			Expect(domain.UpdatedAt).To(Equal(now))
		})

		It("handles nil optional fields", func() {
			model := &Skill{
				ID:        "skill-2",
				Name:      "Python",
				Category:  "backend",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			domain := model.ToDomain()

			Expect(domain.YearsUsed).To(BeNil())
			Expect(domain.LastUsed).To(BeNil())
		})
	})

	Describe("SkillFromDomain", func() {
		It("converts all fields from domain Skill", func() {
			years := 3
			lastUsed := time.Now().AddDate(0, -2, 0)
			domainSkill := fixtures.SkillWithYears("skill-d1", "Ruby", "backend", years)
			domainSkill.Level = "advanced"
			domainSkill.LastUsed = &lastUsed

			model := SkillFromDomain(domainSkill)

			Expect(model.ID).To(Equal("skill-d1"))
			Expect(model.Name).To(Equal("Ruby"))
			Expect(model.Category).To(Equal("backend"))
			Expect(model.Level).To(Equal("advanced"))
			Expect(*model.YearsUsed).To(Equal(3))
			Expect(*model.LastUsed).To(Equal(lastUsed))
			Expect(model.CreatedAt).To(Equal(domainSkill.CreatedAt))
			Expect(model.UpdatedAt).To(Equal(domainSkill.UpdatedAt))
		})

		It("handles nil optional fields", func() {
			domainSkill := fixtures.SkillWith("skill-d2", "Docker", "devops", "")

			model := SkillFromDomain(domainSkill)

			Expect(model.YearsUsed).To(BeNil())
			Expect(model.LastUsed).To(BeNil())
		})
	})

	Describe("ToDomain-FromDomain round trip", func() {
		It("preserves all Skill data", func() {
			years := 7
			lastUsed := time.Now().Truncate(time.Second).AddDate(0, -3, 0)
			original := fixtures.SkillWithYears("skill-rt", "Kubernetes", "devops", years)
			original.Level = "expert"
			original.LastUsed = &lastUsed
			original.CreatedAt = original.CreatedAt.Truncate(time.Second)
			original.UpdatedAt = original.UpdatedAt.Truncate(time.Second)

			model := SkillFromDomain(original)
			restored := model.ToDomain()

			Expect(restored.ID).To(Equal(original.ID))
			Expect(restored.Name).To(Equal(original.Name))
			Expect(restored.Category).To(Equal(original.Category))
			Expect(restored.Level).To(Equal(original.Level))
			Expect(*restored.YearsUsed).To(Equal(*original.YearsUsed))
			Expect(*restored.LastUsed).To(Equal(*original.LastUsed))
			Expect(restored.CreatedAt).To(Equal(original.CreatedAt))
			Expect(restored.UpdatedAt).To(Equal(original.UpdatedAt))
		})
	})
})
