package intents

import (
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/domain/career"
)

// BrowseTimelineContext is the minimal context passed to BrowseTimeline intent.
// It contains only what's necessary to start the intent.
type BrowseTimelineContext struct {
	// Events is the list of events to browse.
	Events []*career.CareerEvent

	// InitialFilters is the initial filter state (may be empty).
	InitialFilters *TimelineFilters

	// SelectedEventID is the initially selected event (may be empty).
	SelectedEventID string

	// CLIEventService is the service for event CRUD operations (edit/delete).
	CLIEventService *service.CLIEventService
}

// Validate ensures the context is complete.
func (c *BrowseTimelineContext) Validate() error {
	if c.Events == nil {
		c.Events = make([]*career.CareerEvent, 0)
	}
	if c.InitialFilters == nil {
		c.InitialFilters = &TimelineFilters{
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

// TimelineFilters represents the current filter and sort state.
type TimelineFilters struct {
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

// BrowseTimelineResult is the output of a successful BrowseTimeline intent.
type BrowseTimelineResult struct {
	// SelectedEvent is the event selected by the user (may be nil if cancelled).
	SelectedEvent *career.CareerEvent

	// FinalFilters are the final filter state when exiting the intent.
	FinalFilters *TimelineFilters

	// ViewedEvents are the events viewed during the session (for analytics).
	ViewedEvents []*career.CareerEvent

	// SelectedFacts are facts selected from the event (may be empty).
	SelectedFacts []*career.Fact
}

// BrowseTimelineModel represents the local state of the BrowseTimeline intent.
// This is the ONLY mutable state owned by the intent.
type BrowseTimelineModel struct {
	// context is the input context passed to the intent.
	context *BrowseTimelineContext

	// currentState tracks which view is active.
	currentState string // BrowseStateTimeline, BrowseStateEventDetail, BrowseStateDeleteConfirm

	// filteredEvents are the events after applying current filters.
	filteredEvents []*career.CareerEvent

	// selectedIndex is the index of the currently selected event.
	selectedIndex int

	// filters is the current filter and sort state.
	filters *TimelineFilters

	// filterStack tracks active filters in FIFO order for progressive clearing.
	filterStack *FilterStack

	// selectedEvent is the event currently being viewed.
	selectedEvent *career.CareerEvent

	// selectedFacts are facts selected from the event.
	selectedFacts []*career.Fact

	// viewedEvents tracks events viewed during the session.
	viewedEvents []*career.CareerEvent

	// deleteError stores any error from delete operation.
	deleteError error
}

// BrowseTimelineStates for navigation.
const (
	BrowseStateTimeline      = "timeline"
	BrowseStateEventDetail   = "event_detail"
	BrowseStateDeleteConfirm = "delete_confirm"
)
