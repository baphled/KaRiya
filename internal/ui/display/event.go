package display

import (
	"time"

	"github.com/baphled/kariya/internal/domain/career"
)

// Event is a presentation-only view of a career event.
type Event struct {
	ID         string
	Text       string
	Date       time.Time
	Company    string
	Project    string
	Tags       []string
	Categories []string
	Skills     []string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// EventFromDomain converts a domain event to a display event.
//
// Expected:
//   - event must be valid.
//
// Returns:
//   - A Event value.
//
// Side effects:
//   - None.
func EventFromDomain(e *career.Event) Event {
	if e == nil {
		return Event{}
	}

	return Event{
		ID:         e.ID,
		Text:       e.Text,
		Date:       e.Date,
		Company:    e.Company,
		Project:    e.Project,
		Tags:       append([]string(nil), e.Tags...),
		Categories: append([]string(nil), e.Categories...),
		Skills:     append([]string(nil), e.Skills...),
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
	}
}

// EventsFromDomain converts domain events to display events.
//
// Expected:
//   - event must be valid.
//
// Returns:
//   - A []Event value.
//
// Side effects:
//   - None.
func EventsFromDomain(events []*career.Event) []Event {
	if events == nil {
		return nil
	}

	result := make([]Event, len(events))
	for i, event := range events {
		result[i] = EventFromDomain(event)
	}

	return result
}
