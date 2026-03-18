package career

import (
	"time"

	"github.com/baphled/kariya/internal/domain/career"
)

// Fact represents a GORM database row in the "facts" table.
//
// Each Fact is a structured insight derived from a career event or burst. The
// ID field holds a UUID primary key, and Text holds the fact description.
// CompetencyCategories stores competency area labels as comma-separated text
// via StringSlice in the "competencies" column. RoleFit stores a role
// suitability level as a plain string. AudienceRelevance stores audience
// labels as comma-separated text via StringSlice. StrengthSignal is an
// optional descriptor indicating what professional strength the fact
// evidences. SourceEventID and SourceBurstID link back to the originating
// event or burst that produced this fact. CreatedAt and UpdatedAt track
// record timestamps.
type Fact struct {
	ID                   string      `gorm:"primaryKey"`
	Text                 string      `gorm:"not null"`
	CompetencyCategories StringSlice `gorm:"column:competencies;type:text;not null"`
	RoleFit              string      `gorm:"not null"`
	AudienceRelevance    StringSlice `gorm:"type:text;not null"`
	StrengthSignal       string
	SourceEventID        string    `gorm:"index"`
	SourceBurstID        string    `gorm:"index"`
	CreatedAt            time.Time `gorm:"not null"`
	UpdatedAt            time.Time `gorm:"not null"`
}

// TableName returns the SQLite table name for the Fact model.
//
// Returns:
//   - The string "facts".
//
// Side effects:
//   - None.
func (Fact) TableName() string { return "facts" }

// ToDomain converts a database Fact record into the corresponding domain
//
// Returns:
//   - A fully initialized career.Fact ready for use.
//
// Side effects:
//   - None.
func (m *Fact) ToDomain() *career.Fact {
	return &career.Fact{
		ID:                   m.ID,
		Text:                 m.Text,
		CompetencyCategories: []string(m.CompetencyCategories),
		RoleFit:              career.RoleFit(m.RoleFit),
		AudienceRelevance:    []string(m.AudienceRelevance),
		StrengthSignal:       m.StrengthSignal,
		SourceEventID:        m.SourceEventID,
		SourceBurstID:        m.SourceBurstID,
		CreatedAt:            m.CreatedAt,
		UpdatedAt:            m.UpdatedAt,
	}
}

// FactFromDomain creates a database Fact record from a domain career.Fact
//
// Expected:
//   - fact must be valid.
//
// Returns:
//   - A fully initialized Fact ready for use.
//
// Side effects:
//   - None.
func FactFromDomain(f *career.Fact) *Fact {
	return &Fact{
		ID:                   f.ID,
		Text:                 f.Text,
		CompetencyCategories: StringSlice(f.CompetencyCategories),
		RoleFit:              string(f.RoleFit),
		AudienceRelevance:    StringSlice(f.AudienceRelevance),
		StrengthSignal:       f.StrengthSignal,
		SourceEventID:        f.SourceEventID,
		SourceBurstID:        f.SourceBurstID,
		CreatedAt:            f.CreatedAt,
		UpdatedAt:            f.UpdatedAt,
	}
}
