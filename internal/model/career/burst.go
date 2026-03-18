package career

import (
	"time"

	"github.com/baphled/kariya/internal/domain/career"
)

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
// Returns:
//   - The string "bursts".
//
// Side effects:
//   - None.
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
