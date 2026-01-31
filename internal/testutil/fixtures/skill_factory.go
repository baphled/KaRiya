package fixtures

import (
	"fmt"
	"time"

	"github.com/bluele/factory-go/factory"
	"github.com/brianvoe/gofakeit/v7"

	"github.com/baphled/kariya/internal/domain/career"
)

// SkillFactory creates Skill fixtures with realistic fake data.
var SkillFactory = factory.NewFactory(
	&career.Skill{},
).SeqInt("ID", func(n int) (interface{}, error) {
	return fmt.Sprintf("skill-%d", n), nil
}).Attr("Name", func(args factory.Args) (interface{}, error) {
	skills := []string{"Go", "Python", "Ruby", "Kubernetes", "React", "PostgreSQL", "Docker", "TypeScript", "Terraform", "Redis"}
	return skills[gofakeit.Number(0, len(skills)-1)], nil
}).Attr("Category", func(args factory.Args) (interface{}, error) {
	categories := []string{"backend", "frontend", "devops", "database", "cloud", "tooling"}
	return categories[gofakeit.Number(0, len(categories)-1)], nil
}).Attr("Level", func(args factory.Args) (interface{}, error) {
	levels := []string{"beginner", "intermediate", "advanced", "expert"}
	return levels[gofakeit.Number(0, len(levels)-1)], nil
}).Attr("CreatedAt", func(args factory.Args) (interface{}, error) {
	return time.Now(), nil
}).Attr("UpdatedAt", func(args factory.Args) (interface{}, error) {
	return time.Now(), nil
})

// Skill creates a minimal valid Skill with the given ID.
func Skill(id string) *career.Skill {
	now := time.Now()
	return &career.Skill{
		ID:        id,
		Name:      "Test Skill " + id,
		Category:  "backend",
		Level:     "intermediate",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// SkillWith creates a Skill with custom key fields.
func SkillWith(id, name, category, level string) *career.Skill {
	now := time.Now()
	return &career.Skill{
		ID:        id,
		Name:      name,
		Category:  category,
		Level:     level,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// SkillWithYears creates a Skill with years of experience set.
func SkillWithYears(id, name, category string, years int) *career.Skill {
	now := time.Now()
	return &career.Skill{
		ID:        id,
		Name:      name,
		Category:  category,
		Level:     "advanced",
		YearsUsed: &years,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Skills creates n skills with sequential IDs.
func Skills(n int) []*career.Skill {
	skills := make([]*career.Skill, n)
	for i := range n {
		skill, ok := SkillFactory.MustCreate().(*career.Skill)
		if !ok {
			continue
		}
		skills[i] = skill
	}
	return skills
}
