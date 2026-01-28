// Package models provides database models for GORM.
package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
)

// StringSlice handles comma-separated string serialization for SQLite.
// This maintains backward compatibility with existing data stored as "a,b,c".
type StringSlice []string

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

func (s StringSlice) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "", nil
	}
	return strings.Join(s, ","), nil
}

// Skill is the database model for skills.
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

func (Skill) TableName() string { return "skills" }

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

// Event is the database model for career events.
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

func (Event) TableName() string { return "career_events" }

func (m *Event) ToDomain() *career.CareerEvent {
	skillIDs := make([]string, len(m.Skills))
	for i, s := range m.Skills {
		skillIDs[i] = s.ID
	}
	return &career.CareerEvent{
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

func EventFromDomain(e *career.CareerEvent) *Event {
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

// Fact is the database model for facts.
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

func (Fact) TableName() string { return "facts" }

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

// Burst is the database model for bursts.
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

func (Burst) TableName() string { return "bursts" }

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
