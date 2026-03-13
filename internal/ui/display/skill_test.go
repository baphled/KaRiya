package display_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/ui/display"
)

var _ = Describe("Skill display type", func() {
	Describe("SkillFromDomain", func() {
		It("converts a fully populated skill", func() {
			yearsUsed := 7
			lastUsed := time.Date(2024, time.February, 15, 0, 0, 0, 0, time.UTC)
			createdAt := time.Date(2023, time.January, 1, 9, 0, 0, 0, time.UTC)
			updatedAt := time.Date(2024, time.January, 1, 9, 0, 0, 0, time.UTC)

			skill := fixtures.SkillWith("skill-1", "Go", "backend", "expert")
			skill.YearsUsed = &yearsUsed
			skill.LastUsed = &lastUsed
			skill.CreatedAt = createdAt
			skill.UpdatedAt = updatedAt

			result := display.SkillFromDomain(skill)

			Expect(result).To(Equal(display.Skill{
				ID:        "skill-1",
				Name:      "Go",
				Category:  "backend",
				Level:     "expert",
				YearsUsed: &yearsUsed,
				LastUsed:  &lastUsed,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			}))
		})

		It("handles nil input gracefully", func() {
			Expect(display.SkillFromDomain(nil)).To(Equal(display.Skill{}))
		})

		It("preserves nil optional fields", func() {
			skill := fixtures.Skill("skill-1")
			result := display.SkillFromDomain(skill)

			Expect(result.YearsUsed).To(BeNil())
			Expect(result.LastUsed).To(BeNil())
		})
	})

	Describe("SkillsFromDomain", func() {
		It("converts a slice of skills", func() {
			s1 := fixtures.Skill("skill-1")
			s1.Name = "Go"
			s2 := fixtures.Skill("skill-2")
			s2.Name = "Ruby"
			skills := []*career.Skill{s1, s2}

			result := display.SkillsFromDomain(skills)

			Expect(result).To(HaveLen(2))
			Expect(result[0].ID).To(Equal("skill-1"))
			Expect(result[0].Name).To(Equal("Go"))
			Expect(result[1].ID).To(Equal("skill-2"))
			Expect(result[1].Name).To(Equal("Ruby"))
		})

		It("returns nil for nil input", func() {
			Expect(display.SkillsFromDomain(nil)).To(BeNil())
		})

		It("returns empty slice for empty input", func() {
			Expect(display.SkillsFromDomain([]*career.Skill{})).To(Equal([]display.Skill{}))
		})
	})
})
