// Package browse_timeline implements the BrowseTimeline intent for browsing career events.
package browse_timeline

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
)

// EventService defines the interface for event CRUD operations.
// This allows for mocking in tests.
type EventService interface {
	DeleteEvent(ctx context.Context, eventID string) error
	ListEvents(ctx context.Context, filters *careerrepo.ListFilters) ([]*career.CareerEvent, error)
	CaptureEvent(ctx context.Context, text string, date time.Time, mode careerservice.EventCaptureMode, opts ...service.Option) error
	UpdateEventMetadata(ctx context.Context, event *career.CareerEvent) error
	GetSkillsForEvent(ctx context.Context, eventID string) ([]*career.Skill, error)
}

// IntentContext is the minimal context passed to BrowseTimeline intent.
// It contains only what's necessary to start the intent.
type IntentContext struct {
	// Events is the list of events to browse.
	Events []*career.CareerEvent

	// InitialFilters is the initial filter state (may be empty).
	InitialFilters *Filters

	// SelectedEventID is the initially selected event (may be empty).
	SelectedEventID string

	// CLIEventService is the service for event CRUD operations (edit/delete).
	CLIEventService EventService
}

// Validate ensures the context is complete.
func (c *IntentContext) Validate() error {
	if c.Events == nil {
		c.Events = make([]*career.CareerEvent, 0)
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

// Filters represents the current filter and sort state.
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
