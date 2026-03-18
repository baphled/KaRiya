package career

import (
	"time"

	"github.com/baphled/kariya/internal/domain/career"
)

// Skill represents a GORM database row in the "skills" table.
//
// Each Skill has a UUID primary key in the ID field, a unique Name, and a
// Category used for grouping. The optional Level field stores a proficiency
// descriptor, YearsUsed stores the number of years the skill has been
// practised (nil when unknown), and LastUsed stores the date the skill was
// most recently applied (nil when unknown). CreatedAt and UpdatedAt track
// record timestamps. The Events field declares a many-to-many relationship
// with Event through the event_skills join table.
type Skill struct {
	ID        string `gorm:"primaryKey"`
	Name      string `gorm:"uniqueIndex;not null"`
	Category  string `gorm:"index;not null"`
	Level     string
	YearsUsed *int
	LastUsed  *time.Time
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	Events    []Event   `gorm:"many2many:event_skills;joinForeignKey:skill_id;joinReferences:event_id"`
}

// TableName returns the SQLite table name for the Skill model.
//
// Returns:
//   - The string "skills".
//
// Side effects:
//   - None.
func (Skill) TableName() string { return "skills" }

// ToDomain converts a database Skill record into the corresponding domain
//
// Returns:
//   - A fully initialized career.Skill ready for use.
//
// Side effects:
//   - None.
func (m *Skill) ToDomain() *career.Skill {
	return &career.Skill{
		ID:        m.ID,
		Name:      m.Name,
		Category:  m.Category,
		Level:     m.Level,
		YearsUsed: m.YearsUsed,
		LastUsed:  m.LastUsed,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// SkillFromDomain creates a database Skill record from a domain career.Skill
//
// Expected:
//   - skill must be valid.
//
// Returns:
//   - A fully initialized Skill ready for use.
//
// Side effects:
//   - None.
func SkillFromDomain(s *career.Skill) *Skill {
	return &Skill{
		ID:        s.ID,
		Name:      s.Name,
		Category:  s.Category,
		Level:     s.Level,
		YearsUsed: s.YearsUsed,
		LastUsed:  s.LastUsed,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}
