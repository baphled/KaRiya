package display

import (
	"time"

	"github.com/baphled/kariya/internal/domain/career"
)

// Burst is a presentation-only view of a career burst.
type Burst struct {
	ID          string
	Name        string
	Description string
	EventIDs    []string
	Confirmed   bool
	ConfirmedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// BurstFromDomain converts a domain burst to a display burst.
//
// Expected:
//   - burst must be valid.
//
// Returns:
//   - A Burst value.
//
// Side effects:
//   - None.
func BurstFromDomain(b *career.Burst) Burst {
	if b == nil {
		return Burst{}
	}

	return Burst{
		ID:          b.ID,
		Name:        b.Name,
		Description: b.Description,
		EventIDs:    append([]string(nil), b.EventIDs...),
		Confirmed:   b.Confirmed,
		ConfirmedAt: cloneTimePointer(b.ConfirmedAt),
		CreatedAt:   b.CreatedAt,
		UpdatedAt:   b.UpdatedAt,
	}
}

// BurstsFromDomain converts domain bursts to display bursts.
//
// Expected:
//   - burst must be valid.
//
// Returns:
//   - A []Burst value.
//
// Side effects:
//   - None.
func BurstsFromDomain(bursts []*career.Burst) []Burst {
	if bursts == nil {
		return nil
	}

	result := make([]Burst, len(bursts))
	for i, burst := range bursts {
		result[i] = BurstFromDomain(burst)
	}

	return result
}
