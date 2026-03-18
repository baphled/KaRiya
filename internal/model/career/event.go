package career

import (
	"time"

	"github.com/baphled/kariya/internal/domain/career"
)

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
// Returns:
//   - The string "career_events".
//
// Side effects:
//   - None.
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
