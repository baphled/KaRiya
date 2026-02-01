package fixtures

import (
	"fmt"
	"time"

	"github.com/bluele/factory-go/factory"
	"github.com/brianvoe/gofakeit/v7"

	"github.com/baphled/kariya/internal/domain/career"
)

// FactFactory creates Fact fixtures with realistic fake data.
var FactFactory = factory.NewFactory(
	&career.Fact{},
).SeqInt("ID", func(n int) (interface{}, error) {
	return fmt.Sprintf("fact-%d", n), nil
}).Attr("Text", func(args factory.Args) (interface{}, error) {
	templates := []string{
		"Reduced %s latency by %d%% through system redesign",
		"Mentored %d junior developers on %s practices",
		"Led cross-functional team of %d to deliver %s",
		"Implemented %s architecture improving scalability by %d%%",
		"Designed and built %s system handling %dK requests/second",
		"Optimized %s queries reducing page load by %d%%",
	}
	template := templates[gofakeit.Number(0, len(templates)-1)]
	return fmt.Sprintf(template, gofakeit.BuzzWord(), gofakeit.Number(20, 80)), nil
}).Attr("CompetencyCategories", func(args factory.Args) (interface{}, error) {
	categories := []string{"technical", "leadership", "product", "consulting", "mentoring", "research"}
	count := gofakeit.Number(1, 3)
	// Use a map to ensure uniqueness
	seen := make(map[string]bool)
	selected := make([]string, 0, count)
	for len(selected) < count {
		cat := categories[gofakeit.Number(0, len(categories)-1)]
		if !seen[cat] {
			seen[cat] = true
			selected = append(selected, cat)
		}
	}
	return selected, nil
}).Attr("RoleFit", func(args factory.Args) (interface{}, error) {
	fits := []career.RoleFit{career.RoleFitStaff, career.RoleFitSeniorIC, career.RoleFitPrincipal, career.RoleFitEM}
	return fits[gofakeit.Number(0, len(fits)-1)], nil
}).Attr("AudienceRelevance", func(args factory.Args) (interface{}, error) {
	audiences := []string{"hiring_manager", "recruiter", "peer"}
	count := gofakeit.Number(1, 2)
	// Use a map to ensure uniqueness
	seen := make(map[string]bool)
	selected := make([]string, 0, count)
	for len(selected) < count {
		aud := audiences[gofakeit.Number(0, len(audiences)-1)]
		if !seen[aud] {
			seen[aud] = true
			selected = append(selected, aud)
		}
	}
	return selected, nil
}).Attr("StrengthSignal", func(args factory.Args) (interface{}, error) {
	signals := []string{"high", "medium", "strong"}
	return signals[gofakeit.Number(0, len(signals)-1)], nil
}).Attr("SourceEventID", func(args factory.Args) (interface{}, error) {
	return fmt.Sprintf("event-%d", gofakeit.Number(1, 100)), nil
}).Attr("CreatedAt", func(args factory.Args) (interface{}, error) {
	return time.Now(), nil
}).Attr("UpdatedAt", func(args factory.Args) (interface{}, error) {
	return time.Now(), nil
})

// Fact creates a minimal valid Fact with the given ID and source event ID.
func Fact(id, sourceEventID string) *career.Fact {
	now := time.Now()
	return &career.Fact{
		ID:                   id,
		Text:                 "Test fact " + id,
		CompetencyCategories: []string{"technical"},
		RoleFit:              career.RoleFitStaff,
		AudienceRelevance:    []string{"hiring_manager"},
		StrengthSignal:       "high",
		SourceEventID:        sourceEventID,
		CreatedAt:            now,
		UpdatedAt:            now,
	}
}

// FactFromBurst creates a minimal valid Fact linked to a burst instead of an event.
func FactFromBurst(id, sourceBurstID string) *career.Fact {
	now := time.Now()
	return &career.Fact{
		ID:                   id,
		Text:                 "Test fact " + id,
		CompetencyCategories: []string{"technical"},
		RoleFit:              career.RoleFitStaff,
		AudienceRelevance:    []string{"hiring_manager"},
		StrengthSignal:       "high",
		SourceBurstID:        sourceBurstID,
		CreatedAt:            now,
		UpdatedAt:            now,
	}
}

// Facts creates n facts with sequential IDs, linked to the provided events.
func Facts(n int, events []*career.Event) []*career.Fact {
	facts := make([]*career.Fact, n)
	for i := range n {
		fact, ok := FactFactory.MustCreate().(*career.Fact)
		if !ok {
			continue
		}
		// Link to actual events if provided
		if len(events) > 0 {
			fact.SourceEventID = events[i%len(events)].ID
		}
		facts[i] = fact
	}
	return facts
}

// FactWith creates a minimal Fact with just ID and text.
// WARNING: This creates an incomplete Fact without required domain fields
// (CompetencyCategories, RoleFit, AudienceRelevance, StrengthSignal, SourceEventID).
// Use only for tests that don't validate the complete Fact structure.
// For fully-populated facts, use Fact() or FactFactory instead.
func FactWith(id, text string) *career.Fact {
	return &career.Fact{
		ID:   id,
		Text: text,
	}
}

// FactWithCategories creates a Fact with specified categories and audience relevance.
// Use for tests that need specific category assignments (e.g., role-based scoring).
func FactWithCategories(id, text, sourceEventID string, categories, audienceRelevance []string) *career.Fact {
	now := time.Now()
	return &career.Fact{
		ID:                   id,
		Text:                 text,
		CompetencyCategories: categories,
		RoleFit:              career.RoleFitStaff,
		AudienceRelevance:    audienceRelevance,
		StrengthSignal:       "high",
		SourceEventID:        sourceEventID,
		CreatedAt:            now,
		UpdatedAt:            now,
	}
}

// FactForValidation creates a fully-populated Fact suitable for validation tests.
// All required fields are set with specified values, allowing targeted field overrides after creation.
func FactForValidation(id, text string, roleFit career.RoleFit, categories, audienceRelevance []string, sourceEventID string) *career.Fact {
	now := time.Now()
	return &career.Fact{
		ID:                   id,
		Text:                 text,
		CompetencyCategories: categories,
		RoleFit:              roleFit,
		AudienceRelevance:    audienceRelevance,
		StrengthSignal:       "leadership capability",
		SourceEventID:        sourceEventID,
		CreatedAt:            now,
		UpdatedAt:            now,
	}
}

// FactForSave creates a Fact without an ID, suitable for SaveFact tests.
// The service will assign the ID upon saving.
func FactForSave(text string, categories []string, roleFit career.RoleFit, audienceRelevance []string, sourceEventID string) *career.Fact {
	return &career.Fact{
		Text:                 text,
		CompetencyCategories: categories,
		RoleFit:              roleFit,
		AudienceRelevance:    audienceRelevance,
		StrengthSignal:       "technical",
		SourceEventID:        sourceEventID,
	}
}
