// Package fixtures provides test data factories for KaRiya domain objects.
// It uses factory-go for the factory pattern and gofakeit for realistic fake data.
//
// Usage:
//
//	// Simple creation with defaults
//	event := fixtures.EventFactory.MustCreate().(*career.CareerEvent)
//
//	// With overrides
//	event := fixtures.EventFactory.MustCreateWithOption(map[string]interface{}{
//	    "Company": "TechCorp",
//	    "Project": "Platform",
//	}).(*career.CareerEvent)
//
//	// Quick helpers for minimal valid objects
//	event := fixtures.Event("my-id")
//	burst := fixtures.Burst("burst-id", "evt-1", "evt-2")
//	fact := fixtures.Fact("fact-id", "evt-1")
package fixtures

import (
	"fmt"
	"time"

	"github.com/bluele/factory-go/factory"
	"github.com/brianvoe/gofakeit/v7"

	"github.com/baphled/kariya/internal/domain/career"
)

func init() {
	// Seed gofakeit for reproducible tests when needed
	// Can be overridden with SetSeed() for deterministic tests
	gofakeit.Seed(0)
}

// SetSeed sets the random seed for reproducible test data.
// Use this at the start of tests that need deterministic data.
func SetSeed(seed int64) {
	gofakeit.Seed(seed)
}

// -----------------------------------------------------------------------------
// Event Factory
// -----------------------------------------------------------------------------

// EventFactory creates CareerEvent fixtures with realistic fake data.
// Each call generates unique IDs and varied content.
var EventFactory = factory.NewFactory(
	&career.CareerEvent{},
).SeqInt("ID", func(n int) (interface{}, error) {
	return fmt.Sprintf("event-%d", n), nil
}).Attr("Text", func(args factory.Args) (interface{}, error) {
	// Generate realistic career event text
	verbs := []string{"Led", "Implemented", "Designed", "Optimized", "Built", "Mentored", "Managed", "Delivered"}
	objects := []string{
		"team of %d engineers on %s project",
		"new %s system reducing latency by %d%%",
		"architecture for %s platform",
		"database queries improving performance by %d%%",
		"CI/CD pipeline for %s deployment",
		"junior developers on %s practices",
		"cross-functional project with %d stakeholders",
		"Q%d release on schedule",
	}
	verb := verbs[gofakeit.Number(0, len(verbs)-1)]
	obj := objects[gofakeit.Number(0, len(objects)-1)]
	return fmt.Sprintf("%s "+obj, verb, gofakeit.Number(3, 12), gofakeit.BuzzWord()), nil
}).Attr("Date", func(args factory.Args) (interface{}, error) {
	// Random date within last 2 years
	return gofakeit.DateRange(
		time.Now().AddDate(-2, 0, 0),
		time.Now(),
	), nil
}).Attr("Company", func(args factory.Args) (interface{}, error) {
	return gofakeit.Company(), nil
}).Attr("Project", func(args factory.Args) (interface{}, error) {
	return gofakeit.BuzzWord() + " " + gofakeit.AppName(), nil
}).Attr("Tags", func(args factory.Args) (interface{}, error) {
	tags := []string{"technical", "leadership", "mentoring", "product", "research", "architecture"}
	count := gofakeit.Number(1, 3)
	selected := make([]string, count)
	for i := 0; i < count; i++ {
		selected[i] = tags[gofakeit.Number(0, len(tags)-1)]
	}
	return selected, nil
}).Attr("Categories", func(args factory.Args) (interface{}, error) {
	categories := []string{"technical", "leadership", "product", "consulting", "mentoring", "research"}
	count := gofakeit.Number(1, 2)
	selected := make([]string, count)
	for i := 0; i < count; i++ {
		selected[i] = categories[gofakeit.Number(0, len(categories)-1)]
	}
	return selected, nil
}).Attr("CreatedAt", func(args factory.Args) (interface{}, error) {
	return time.Now(), nil
}).Attr("UpdatedAt", func(args factory.Args) (interface{}, error) {
	return time.Now(), nil
})

// -----------------------------------------------------------------------------
// Burst Factory
// -----------------------------------------------------------------------------

// BurstFactory creates Burst fixtures with realistic fake data.
var BurstFactory = factory.NewFactory(
	&career.Burst{},
).SeqInt("ID", func(n int) (interface{}, error) {
	return fmt.Sprintf("burst-%d", n), nil
}).Attr("Name", func(args factory.Args) (interface{}, error) {
	prefixes := []string{"Q1", "Q2", "Q3", "Q4", "Annual", "Sprint"}
	suffixes := []string{"Initiative", "Project", "Migration", "Overhaul", "Optimization", "Refactoring"}
	prefix := prefixes[gofakeit.Number(0, len(prefixes)-1)]
	suffix := suffixes[gofakeit.Number(0, len(suffixes)-1)]
	return fmt.Sprintf("%s %s %s", prefix, gofakeit.BuzzWord(), suffix), nil
}).Attr("Description", func(args factory.Args) (interface{}, error) {
	return gofakeit.Sentence(gofakeit.Number(8, 15)), nil
}).Attr("EventIDs", func(args factory.Args) (interface{}, error) {
	// Default to 2 event IDs (minimum required for a burst)
	return []string{
		fmt.Sprintf("event-%d", gofakeit.Number(1, 100)),
		fmt.Sprintf("event-%d", gofakeit.Number(101, 200)),
	}, nil
}).Attr("Confirmed", func(args factory.Args) (interface{}, error) {
	return false, nil
}).Attr("CreatedAt", func(args factory.Args) (interface{}, error) {
	return time.Now(), nil
}).Attr("UpdatedAt", func(args factory.Args) (interface{}, error) {
	return time.Now(), nil
})

// -----------------------------------------------------------------------------
// Fact Factory
// -----------------------------------------------------------------------------

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
	selected := make([]string, count)
	for i := 0; i < count; i++ {
		selected[i] = categories[gofakeit.Number(0, len(categories)-1)]
	}
	return selected, nil
}).Attr("RoleFit", func(args factory.Args) (interface{}, error) {
	fits := []career.RoleFit{career.RoleFitStaff, career.RoleFitSeniorIC, career.RoleFitPrincipal, career.RoleFitEM}
	return fits[gofakeit.Number(0, len(fits)-1)], nil
}).Attr("AudienceRelevance", func(args factory.Args) (interface{}, error) {
	audiences := []string{"hiring_manager", "recruiter", "peer"}
	count := gofakeit.Number(1, 2)
	selected := make([]string, count)
	for i := 0; i < count; i++ {
		selected[i] = audiences[gofakeit.Number(0, len(audiences)-1)]
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

// -----------------------------------------------------------------------------
// Quick Helper Functions
// -----------------------------------------------------------------------------

// Event creates a minimal valid CareerEvent with the given ID.
// Use this for simple tests that just need a valid event.
func Event(id string) *career.CareerEvent {
	now := time.Now()
	return &career.CareerEvent{
		ID:        id,
		Text:      "Test event " + id,
		Date:      now,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// EventWith creates a CareerEvent with custom fields.
// Unspecified fields get sensible defaults.
func EventWith(id, text, company, project string) *career.CareerEvent {
	now := time.Now()
	return &career.CareerEvent{
		ID:        id,
		Text:      text,
		Company:   company,
		Project:   project,
		Date:      now,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Burst creates a minimal valid Burst with the given ID and event IDs.
// At least 2 event IDs are required for a valid burst.
func Burst(id string, eventIDs ...string) *career.Burst {
	if len(eventIDs) < 2 {
		eventIDs = []string{"event-1", "event-2"}
	}
	now := time.Now()
	return &career.Burst{
		ID:        id,
		Name:      "Test burst " + id,
		EventIDs:  eventIDs,
		Confirmed: false,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// BurstConfirmed creates a confirmed burst with the given ID and event IDs.
func BurstConfirmed(id string, eventIDs ...string) *career.Burst {
	burst := Burst(id, eventIDs...)
	burst.Confirmed = true
	now := time.Now()
	burst.ConfirmedAt = &now
	return burst
}

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

// -----------------------------------------------------------------------------
// Batch Creation Helpers
// -----------------------------------------------------------------------------

// Events creates n events with sequential IDs (event-1, event-2, etc.)
func Events(n int) []*career.CareerEvent {
	events := make([]*career.CareerEvent, n)
	for i := 0; i < n; i++ {
		events[i] = EventFactory.MustCreate().(*career.CareerEvent)
	}
	return events
}

// Bursts creates n bursts with sequential IDs, linked to the provided events.
// If events is empty, bursts will have placeholder event IDs.
func Bursts(n int, events []*career.CareerEvent) []*career.Burst {
	bursts := make([]*career.Burst, n)
	for i := 0; i < n; i++ {
		burst := BurstFactory.MustCreate().(*career.Burst)
		// Link to actual events if provided
		if len(events) >= 2 {
			startIdx := (i * 2) % len(events)
			endIdx := startIdx + 2
			if endIdx > len(events) {
				endIdx = len(events)
			}
			eventIDs := make([]string, 0)
			for j := startIdx; j < endIdx; j++ {
				eventIDs = append(eventIDs, events[j].ID)
			}
			if len(eventIDs) >= 2 {
				burst.EventIDs = eventIDs
			}
		}
		bursts[i] = burst
	}
	return bursts
}

// Facts creates n facts with sequential IDs, linked to the provided events.
func Facts(n int, events []*career.CareerEvent) []*career.Fact {
	facts := make([]*career.Fact, n)
	for i := 0; i < n; i++ {
		fact := FactFactory.MustCreate().(*career.Fact)
		// Link to actual events if provided
		if len(events) > 0 {
			fact.SourceEventID = events[i%len(events)].ID
		}
		facts[i] = fact
	}
	return facts
}
