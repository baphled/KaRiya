// Package browsetimeline implements the BrowseTimeline intent for browsing career events.
package browsetimeline

import (
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/skillinference"
)

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

// SkillLinkedMsg notifies that a skill was successfully linked to an event.
type SkillLinkedMsg struct {
	EventID string
	SkillID string
}

// SkillUnlinkedMsg notifies that a skill was successfully unlinked from an event.
type SkillUnlinkedMsg struct {
	EventID string
	SkillID string
}

// SkillCreatedMsg notifies that a new skill was successfully created.
type SkillCreatedMsg struct {
	Skill *career.Skill
}

// SkillSuggestionsLoadedMsg is sent when skill inference completes.
type SkillSuggestionsLoadedMsg struct {
	Suggestions        []skillinference.SkillSuggestion
	ExistingSkillNames []string
}

// SkillSuggestionsErrorMsg is sent when skill inference fails.
type SkillSuggestionsErrorMsg struct {
	Err error
}
