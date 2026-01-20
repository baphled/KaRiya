package fixtures

import (
	"fmt"
	"time"

	"github.com/bluele/factory-go/factory"
	"github.com/brianvoe/gofakeit/v7"

	"github.com/baphled/kariya/internal/domain/career"
)

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

// Burst creates a minimal valid Burst with the given ID and event IDs.
// At least 2 event IDs are required for a valid burst.
// If fewer than 2 event IDs are provided, defaults to ["event-1", "event-2"].
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

// Bursts creates n bursts with sequential IDs, linked to the provided events.
// If events is empty or has fewer than 2 elements, bursts will have placeholder event IDs.
// Events are distributed using wrap-around to ensure each burst has exactly 2 event IDs.
func Bursts(n int, events []*career.CareerEvent) []*career.Burst {
	bursts := make([]*career.Burst, n)
	for i := 0; i < n; i++ {
		burst := BurstFactory.MustCreate().(*career.Burst)
		// Link to actual events if provided (need at least 2 for valid bursts)
		if len(events) >= 2 {
			startIdx := (i * 2) % len(events)
			// Always assign two event IDs, wrapping around if needed
			eventIDs := []string{
				events[startIdx].ID,
				events[(startIdx+1)%len(events)].ID,
			}
			burst.EventIDs = eventIDs
		}
		bursts[i] = burst
	}
	return bursts
}
