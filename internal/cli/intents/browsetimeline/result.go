// Package browsetimeline implements the BrowseTimeline intent for browsing career events.
package browsetimeline

import "github.com/baphled/kariya/internal/domain/career"

// Result is the output of a successful BrowseTimeline intent.
type Result struct {
	// SelectedEvent is the event selected by the user (may be nil if cancelled).
	SelectedEvent *career.CareerEvent

	// FinalFilters are the final filter state when exiting the intent.
	FinalFilters *Filters

	// ViewedEvents are the events viewed during the session (for analytics).
	ViewedEvents []*career.CareerEvent

	// SelectedFacts are facts selected from the event (may be empty).
	SelectedFacts []*career.Fact
}
