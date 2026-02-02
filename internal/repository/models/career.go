// Package models defines GORM database models and their conversions to and
// from domain types.
//
// Each model struct maps directly to a SQLite table and provides ToDomain and
// FromDomain functions for translating between the persistence layer and the
// application's domain types in the career package.
package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
)

// StringSlice represents a slice of strings that serializes to and from a
// comma-separated string for storage in SQLite text columns.
//
// The encoding preserves backward compatibility with existing data stored in
// the format "a,b,c". StringSlice satisfies both sql.Scanner and
// driver.Valuer so GORM handles the conversion transparently during reads
// and writes.
type StringSlice []string

// Scan deserializes a raw database column value into the StringSlice receiver.
//
// Expected:
//   - interface{} must be valid.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (s *StringSlice) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}
	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	default:
		return fmt.Errorf("unsupported type for StringSlice: %T", value)
	}
	if str == "" {
		*s = nil
		return nil
	}
	*s = strings.Split(str, ",")
	return nil
}

// Value serializes the StringSlice into a comma-separated string suitable for
// database storage.
//
// Returns an empty string and nil error when the slice is nil or empty.
// Returns a single comma-joined string and nil error when the slice contains
// one or more elements.
func (s StringSlice) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "", nil
	}
	return strings.Join(s, ","), nil
}

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
// Returns the string "skills".
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

// Event represents a GORM database row in the "career_events" table.
//
// Each Event captures a professional accomplishment or activity. The ID field
// holds a UUID primary key, Text holds the event description, and Date holds
// the date the event occurred. Company and Project provide optional
// organisational context. Tags and Categories are stored as comma-separated
// text via StringSlice for flexible categorisation. CreatedAt and UpdatedAt
// track record timestamps. The Skills field declares a many-to-many
// relationship with Skill through the event_skills join table.
type Event struct {
	ID         string    `gorm:"primaryKey"`
	Text       string    `gorm:"not null"`
	Date       time.Time `gorm:"index;not null"`
	Company    string    `gorm:"index"`
	Project    string
	Tags       StringSlice `gorm:"type:text"`
	Categories StringSlice `gorm:"type:text"`
	CreatedAt  time.Time   `gorm:"not null"`
	UpdatedAt  time.Time   `gorm:"not null"`
	Skills     []Skill     `gorm:"many2many:event_skills;joinForeignKey:event_id;joinReferences:skill_id"`
}

// TableName returns the SQLite table name for the Event model.
//
// Returns the string "career_events".
func (Event) TableName() string { return "career_events" }

// ToDomain converts a database Event record into the corresponding domain
//
// Returns:
//   - A fully initialized career.Event ready for use.
//
// Side effects:
//   - None.
func (m *Event) ToDomain() *career.Event {
	skillIDs := make([]string, len(m.Skills))
	for i := range m.Skills {
		skillIDs[i] = m.Skills[i].ID
	}
	return &career.Event{
		ID:         m.ID,
		Text:       m.Text,
		Date:       m.Date,
		Company:    m.Company,
		Project:    m.Project,
		Tags:       []string(m.Tags),
		Categories: []string(m.Categories),
		Skills:     skillIDs,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

// EventFromDomain creates a database Event record from a domain career.Event
//
// Expected:
//   - event must be valid.
//
// Returns:
//   - A fully initialized Event ready for use.
//
// Side effects:
//   - None.
func EventFromDomain(e *career.Event) *Event {
	return &Event{
		ID:         e.ID,
		Text:       e.Text,
		Date:       e.Date,
		Company:    e.Company,
		Project:    e.Project,
		Tags:       StringSlice(e.Tags),
		Categories: StringSlice(e.Categories),
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
	}
}

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
// Returns the string "facts".
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

// Burst represents a GORM database row in the "bursts" table.
//
// Each Burst groups related career events that occurred in a concentrated
// period, representing a theme or pattern of activity. The ID field holds a
// UUID primary key, Name holds a short label, and Description holds an
// optional longer explanation. EventIDs stores the associated event
// identifiers as comma-separated text via StringSlice. The Confirmed boolean
// tracks whether a user has validated the detected burst, and ConfirmedAt
// stores the timestamp of that validation (nil when unconfirmed). CreatedAt
// and UpdatedAt track record timestamps.
type Burst struct {
	ID          string `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Description string
	EventIDs    StringSlice `gorm:"type:text;not null"`
	Confirmed   bool        `gorm:"index;default:false"`
	ConfirmedAt *time.Time
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}

// TableName returns the SQLite table name for the Burst model.
//
// Returns the string "bursts".
func (Burst) TableName() string { return "bursts" }

// ToDomain converts a database Burst record into the corresponding domain
//
// Returns:
//   - A fully initialized career.Burst ready for use.
//
// Side effects:
//   - None.
func (m *Burst) ToDomain() *career.Burst {
	return &career.Burst{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		EventIDs:    []string(m.EventIDs),
		Confirmed:   m.Confirmed,
		ConfirmedAt: m.ConfirmedAt,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// BurstFromDomain creates a database Burst record from a domain career.Burst
//
// Expected:
//   - burst must be valid.
//
// Returns:
//   - A fully initialized Burst ready for use.
//
// Side effects:
//   - None.
func BurstFromDomain(b *career.Burst) *Burst {
	return &Burst{
		ID:          b.ID,
		Name:        b.Name,
		Description: b.Description,
		EventIDs:    StringSlice(b.EventIDs),
		Confirmed:   b.Confirmed,
		ConfirmedAt: b.ConfirmedAt,
		CreatedAt:   b.CreatedAt,
		UpdatedAt:   b.UpdatedAt,
	}
}
