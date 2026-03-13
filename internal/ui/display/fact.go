package display

import (
	"time"

	"github.com/baphled/kariya/internal/domain/career"
)

// Fact is a presentation-only view of a career fact.
type Fact struct {
	ID                   string
	Text                 string
	CompetencyCategories []string
	RoleFit              string
	AudienceRelevance    []string
	StrengthSignal       string
	SourceEventID        string
	SourceBurstID        string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

// FactFromDomain converts a domain fact to a display fact.
//
// Expected:
//   - fact must be valid.
//
// Returns:
//   - A Fact value.
//
// Side effects:
//   - None.
func FactFromDomain(f *career.Fact) Fact {
	if f == nil {
		return Fact{}
	}

	return Fact{
		ID:                   f.ID,
		Text:                 f.Text,
		CompetencyCategories: append([]string(nil), f.CompetencyCategories...),
		RoleFit:              string(f.RoleFit),
		AudienceRelevance:    append([]string(nil), f.AudienceRelevance...),
		StrengthSignal:       f.StrengthSignal,
		SourceEventID:        f.SourceEventID,
		SourceBurstID:        f.SourceBurstID,
		CreatedAt:            f.CreatedAt,
		UpdatedAt:            f.UpdatedAt,
	}
}

// FactsFromDomain converts domain facts to display facts.
//
// Expected:
//   - fact must be valid.
//
// Returns:
//   - A []Fact value.
//
// Side effects:
//   - None.
func FactsFromDomain(facts []*career.Fact) []Fact {
	if facts == nil {
		return nil
	}

	result := make([]Fact, len(facts))
	for i, fact := range facts {
		result[i] = FactFromDomain(fact)
	}

	return result
}
