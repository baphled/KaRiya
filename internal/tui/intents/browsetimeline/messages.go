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
	Error   error
}

// SkillUnlinkedMsg notifies that a skill was successfully unlinked from an event.
type SkillUnlinkedMsg struct {
	EventID string
	SkillID string
	Error   error
}

// SkillCreatedMsg notifies that a new skill was successfully created.
type SkillCreatedMsg struct {
	Skill *career.Skill
	Error error
}

// SkillSuggestionsLoadedMsg is sent when skill inference completes.
type SkillSuggestionsLoadedMsg struct {
	Suggestions        []skillinference.SkillSuggestion
	ExistingSkillNames []string
}

// SkillSuggestionsErrorMsg is sent when skill inference fails.
type SkillSuggestionsErrorMsg struct {
	Error error
}

// SkillsForModalLoadedMsg is sent when skills for an event are loaded for the modal.
//
// Expected:
//   - EventID must be valid and non-empty.
//   - Skills slice may be empty if none found or error occurred.
//
// Returns:
//   - Used by intent Update to trigger modal creation.
//
// Side effects:
//   - Triggers modal display for event skills detail.
type SkillsForModalLoadedMsg struct {
	EventID string
	Skills  []*career.Skill
	Error   error
}

// SkillPickerDataLoadedMsg is sent when all skills and event skills are loaded for the picker modal.
type SkillPickerDataLoadedMsg struct {
	AllSkills   []*career.Skill
	EventSkills []*career.Skill
	Error       error
}

// SkillsRefreshedMsg is sent when skills for an event are refreshed for the modal.
type SkillsRefreshedMsg struct {
	Skills []*career.Skill
	Error  error
}
