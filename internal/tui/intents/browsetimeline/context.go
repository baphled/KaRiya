// Package browsetimeline implements the BrowseTimeline intent for browsing career events.
package browsetimeline

import (
	"github.com/baphled/kariya/internal/domain/career"
)

// IntentValidator holds the input parameters required to initialize and operate
// the BrowseTimeline intent, including event data, initial filters, and
// the service dependency for CRUD operations.
type IntentValidator struct {
	// Events is the list of events to browse.
	Events []*career.Event

	// InitialFilters is the initial filter state (may be empty).
	InitialFilters *Filters

	// SelectedEventID is the initially selected event (may be empty).
	SelectedEventID string

	// CLIEventService is the service for event CRUD operations (edit/delete).
	CLIEventService EventService

	// CLISkillCreator is the service for skill operations (create).
	CLISkillCreator SkillCreator

	// SkillInferenceService is the service for inferring skills from event text.
	SkillInferenceService SkillInferenceService
}

// Validate ensures the context is usable by initializing nil fields to
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (c *IntentValidator) Validate() error {
	if c.Events == nil {
		c.Events = make([]*career.Event, 0)
	}
	if c.InitialFilters == nil {
		c.InitialFilters = &Filters{
			Tags:       make([]string, 0),
			Companies:  make([]string, 0),
			Categories: make([]string, 0),
			Projects:   make([]string, 0),
			SortBy:     "date",
			SortOrder:  "desc",
		}
	}
	return nil
}

// Filters encapsulates all user-configurable criteria for narrowing and
// ordering the timeline event list, including text search, field-based
// filters, date ranges, and sort preferences.
type Filters struct {
	// SearchText is the text to search for in event descriptions.
	SearchText string

	// Tags filters events by tags.
	Tags []string

	// Companies filters events by company.
	Companies []string

	// Categories filters events by category.
	Categories []string

	// Projects filters events by project.
	Projects []string

	// DateFrom filters events from this date (optional).
	DateFrom string

	// DateTo filters events up to this date (optional).
	DateTo string

	// SortBy specifies the sort field (date, text, relevance).
	SortBy string

	// SortOrder specifies the sort direction (asc, desc).
	SortOrder string
}
