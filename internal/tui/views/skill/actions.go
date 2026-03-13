package skill

import "github.com/baphled/kariya/internal/ui/display"

// ActionKey constants define the available skill view actions.
type ActionKey string

// ActionKey values for skill view navigation.
const (
	ActionEdit   ActionKey = "skill.edit"
	ActionDelete ActionKey = "skill.delete"
	ActionView   ActionKey = "skill.view"
	ActionAdd    ActionKey = "skill.add"
	ActionFilter ActionKey = "skill.filter"
	ActionSearch ActionKey = "skill.search"
	ActionSort   ActionKey = "skill.sort"
	ActionInfer  ActionKey = "skill.infer"
	ActionHelp   ActionKey = "skill.help"
)

// Nav is the typed payload for skill navigation actions.
type Nav struct {
	Action ActionKey
	Skill  display.Skill
}
