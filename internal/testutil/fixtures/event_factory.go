package fixtures

import (
	"fmt"
	"time"

	"github.com/bluele/factory-go/factory"
	"github.com/brianvoe/gofakeit/v7"

	"github.com/baphled/kariya/internal/domain/career"
)

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
	// Use only valid tags from career.AllowedTags
	tags := []string{"project", "achievement", "leadership", "technical", "consulting", "research", "product", "mentoring"}
	count := gofakeit.Number(1, 3)
	// Use a map to ensure uniqueness
	seen := make(map[string]bool)
	selected := make([]string, 0, count)
	for len(selected) < count {
		tag := tags[gofakeit.Number(0, len(tags)-1)]
		if !seen[tag] {
			seen[tag] = true
			selected = append(selected, tag)
		}
	}
	return selected, nil
}).Attr("Categories", func(args factory.Args) (interface{}, error) {
	// Use only valid categories from career.AllowedCategories
	categories := []string{"technical", "leadership", "product", "consulting", "mentoring", "research"}
	count := gofakeit.Number(1, 2)
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
}).Attr("CreatedAt", func(args factory.Args) (interface{}, error) {
	return time.Now(), nil
}).Attr("UpdatedAt", func(args factory.Args) (interface{}, error) {
	return time.Now(), nil
})

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

// Events creates n events with sequential IDs (event-1, event-2, etc.)
func Events(n int) []*career.CareerEvent {
	events := make([]*career.CareerEvent, n)
	for i := 0; i < n; i++ {
		events[i] = EventFactory.MustCreate().(*career.CareerEvent)
	}
	return events
}
