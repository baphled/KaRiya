// Package manage_skills implements the ManageSkills intent for managing user-defined skills.
package manage_skills

// State represents the current state of the ManageSkills intent.
type State string

// State constants for the ManageSkills intent.
const (
	// StateList shows all skills grouped by category.
	StateList State = "list"

	// StateDetail shows single skill details.
	StateDetail State = "detail"

	// StateDetailEvents shows events using a skill.
	StateDetailEvents State = "detail_events"

	// StateDetailEventDetail shows single event details from skill events.
	StateDetailEventDetail State = "detail_event_detail"

	// StateAdd shows the form for adding a new skill.
	StateAdd State = "add"

	// StateEdit shows the form for editing an existing skill.
	StateEdit State = "edit"

	// StateDelete shows the deletion confirmation.
	StateDelete State = "delete"

	// StateFilter shows the filter menu.
	StateFilter State = "filter"

	// StateSort shows the sort menu.
	StateSort State = "sort"
)
