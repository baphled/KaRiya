package burst

import "github.com/baphled/kariya/internal/ui/display"

// ActionKey is a string-typed constant for burst actions.
type ActionKey string

// ActionView and related constants define the available burst actions.
const (
	ActionView        ActionKey = "view"
	ActionEdit        ActionKey = "edit"
	ActionDelete      ActionKey = "delete"
	ActionAdd         ActionKey = "add"
	ActionSuggest     ActionKey = "suggest"
	ActionFilter      ActionKey = "filter"
	ActionSearch      ActionKey = "search"
	ActionSort        ActionKey = "sort"
	ActionConfirm     ActionKey = "confirm"
	ActionViewEvents  ActionKey = "view_events"
	ActionViewFacts   ActionKey = "view_facts"
	ActionViewSkills  ActionKey = "view_skills"
	ActionInferSkills ActionKey = "infer_skills"
	ActionHelp        ActionKey = "help"
)

// Nav is the typed payload for burst navigation actions.
type Nav struct {
	Action ActionKey
	Burst  display.Burst
}
