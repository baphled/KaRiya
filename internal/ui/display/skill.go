package display

import (
	"time"

	"github.com/baphled/kariya/internal/domain/career"
)

// Skill is a presentation-only view of a career skill.
type Skill struct {
	ID        string
	Name      string
	Category  string
	Level     string
	YearsUsed *int
	LastUsed  *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

// SkillFromDomain converts a domain skill to a display skill.
//
// Expected:
//   - skill must be valid.
//
// Returns:
//   - A Skill value.
//
// Side effects:
//   - None.
func SkillFromDomain(s *career.Skill) Skill {
	if s == nil {
		return Skill{}
	}

	return Skill{
		ID:        s.ID,
		Name:      s.Name,
		Category:  s.Category,
		Level:     s.Level,
		YearsUsed: cloneInt(s.YearsUsed),
		LastUsed:  cloneTimePointer(s.LastUsed),
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

// SkillsFromDomain converts domain skills to display skills.
//
// Expected:
//   - skill must be valid.
//
// Returns:
//   - A []Skill value.
//
// Side effects:
//   - None.
func SkillsFromDomain(skills []*career.Skill) []Skill {
	if skills == nil {
		return nil
	}

	result := make([]Skill, len(skills))
	for i, skill := range skills {
		result[i] = SkillFromDomain(skill)
	}

	return result
}
