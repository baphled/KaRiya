// Package browsetimeline implements the BrowseTimeline intent for browsing career events.
package browsetimeline

import "github.com/baphled/kariya/internal/domain/career"

// Custom message types for BrowseTimeline state transitions.
// ALL *Msg structs MUST be in this file.
//
// RequestEditEventMsg and RequestAddEventMsg are defined in the parent
// intents package (intents/messages.go) since they are cross-intent coordination
// messages used by the app router.

// EventSelectedMsg indicates the user selected an event.
type EventSelectedMsg struct {
	Event *career.Event
	Index int
}

// FilterChangedMsg indicates the filters have changed.
type FilterChangedMsg struct {
	Filters *Filters
}

// EventDeletedMsg notifies that an event was successfully deleted.
// This is sent back to the intent after a delete operation completes.
type EventDeletedMsg struct {
	EventID string
}
