package capture

import (
	"time"

	"github.com/baphled/kariya/internal/domain/career"
)

// EventInput holds the raw data needed to create a career event.
// This struct decouples event creation from any form library.
type EventInput struct {
	Text       string
	Date       string
	Company    string
	Project    string
	Tags       []string
	Categories []string
	Skills     []string
}

// NewEventFromInput creates a career.Event from raw input data.
//
// Expected: Input.Date may be empty (defaults to now) or a valid date string.
// Returns: A fully initialised career.Event, or a DateParseError for invalid dates.
// Side effects: None.
func NewEventFromInput(input EventInput) (*career.Event, error) {
	var eventDate time.Time
	var err error

	if input.Date == "" {
		eventDate = time.Now()
	} else {
		eventDate, err = parseDateString(input.Date)
		if err != nil {
			return nil, err
		}
	}

	event := &career.Event{
		Text:       input.Text,
		Date:       eventDate,
		Company:    input.Company,
		Project:    input.Project,
		Tags:       input.Tags,
		Categories: input.Categories,
		Skills:     input.Skills,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	return event, nil
}

// EditEventInput holds the raw data needed to update an existing career event.
// This struct decouples event editing from any form library.
type EditEventInput struct {
	EventID    string
	Text       string
	Date       string
	Company    string
	Project    string
	Tags       []string
	Categories []string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// UpdateEventFromInput creates a career.Event from edit input data,
// preserving the original event ID and timestamps.
//
// Expected: Input.Date may be empty (defaults to now) or a valid date string.
// Returns: A fully initialised career.Event, or a DateParseError for invalid dates.
// Side effects: None.
func UpdateEventFromInput(input EditEventInput) (*career.Event, error) {
	var eventDate time.Time
	var err error

	if input.Date == "" {
		eventDate = time.Now()
	} else {
		eventDate, err = parseDateString(input.Date)
		if err != nil {
			return nil, err
		}
	}

	return &career.Event{
		ID:         input.EventID,
		Text:       input.Text,
		Date:       eventDate,
		Company:    input.Company,
		Project:    input.Project,
		Tags:       input.Tags,
		Categories: input.Categories,
		Skills:     []string{},
		CreatedAt:  input.CreatedAt,
		UpdatedAt:  input.UpdatedAt,
	}, nil
}

// parseDateString parses a date string supporting ISO, US, UK formats plus "today" and "yesterday".
func parseDateString(s string) (time.Time, error) {
	switch s {
	case "today":
		return time.Now(), nil
	case "yesterday":
		return time.Now().AddDate(0, 0, -1), nil
	}

	formats := []string{
		"2006-01-02",
		"01/02/2006",
		"02-01-2006",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, s); err == nil {
			return t, nil
		}
	}

	return time.Time{}, &DateParseError{Input: s}
}

// DateParseError is returned when a date string cannot be parsed.
type DateParseError struct {
	Input string
}

// Error returns the error message for the unparseable date input.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (e *DateParseError) Error() string {
	return "unable to parse date: " + e.Input
}
