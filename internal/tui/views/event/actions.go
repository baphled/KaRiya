package event

import "github.com/baphled/kariya/internal/ui/display"

// ActionKey is a string-typed constant for event actions.
type ActionKey string

// ActionEdit and related constants define the available event actions.
const (
	ActionEdit   ActionKey = "event.edit"
	ActionDelete ActionKey = "event.delete"
	ActionView   ActionKey = "event.view"
	ActionAdd    ActionKey = "event.add"
	ActionFilter ActionKey = "event.filter"
	ActionSearch ActionKey = "event.search"
	ActionSort   ActionKey = "event.sort"
	ActionClear  ActionKey = "event.clear"
	ActionHelp   ActionKey = "event.help"
)

// Nav is the typed payload for event navigation actions.
type Nav struct {
	Action ActionKey
	Event  display.Event
}

// SkillActionKey is a string-typed constant for skill-related actions.
type SkillActionKey string

// ActionShowEventSkills and related constants define skill-related actions.
const (
	ActionShowEventSkills SkillActionKey = "event.skill.show"
	ActionPickSkill       SkillActionKey = "event.skill.pick"
	ActionAddSkill        SkillActionKey = "event.skill.add"
	ActionInferSkills     SkillActionKey = "event.skill.infer"
)

// SkillNav is the typed payload for skill navigation actions.
type SkillNav struct {
	Action SkillActionKey
	Event  display.Event
}
