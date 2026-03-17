// Package skill provides View implementations for skill management.
package skill

import (
	"strconv"

	"github.com/baphled/kariya/internal/tui/navigation"
	"github.com/baphled/kariya/internal/tui/views/shared"
	"github.com/baphled/kariya/internal/ui/behaviors"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/uikit/primitives"
	"github.com/baphled/kariya/internal/ui/uikit/theme"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

var badgeFuncs = []behaviors.BadgeFunc{
	primitives.NavigateBadge,
	func(th theme.Theme) *primitives.Badge { return primitives.HelpKeyBadge("Enter", "View", th) },
	primitives.AddBadge,
	primitives.EditBadge,
	primitives.DeleteBadge,
	func(th theme.Theme) *primitives.Badge { return primitives.HelpKeyBadge("f", "Filter", th) },
	func(th theme.Theme) *primitives.Badge { return primitives.HelpKeyBadge("s", "Sort", th) },
	func(th theme.Theme) *primitives.Badge { return primitives.HelpKeyBadge("/", "Search", th) },
	primitives.BackBadge,
	primitives.QuitBadge,
	primitives.HelpBadge,
}

// ListView implements widgets.View for displaying a paginated list of skills.
type ListView struct {
	widgets.BaseView
	skills      []display.Skill
	list        *behaviors.ListBehavior[display.Skill]
	keys        navigation.ListKeyMap
	globalKeys  navigation.GlobalKeyMap
	eventCounts map[string]int
}

// NewListView creates a new skill list view.
//
// Expected:
//   - skill must be valid.
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized ListView ready for use.
//
// Side effects:
//   - None.
func NewListView(skills []display.Skill, eventCounts map[string]int) *ListView {
	columns := []behaviors.ColumnDef{
		{Title: "Name", Width: 25},
		{Title: "Category", Width: 15},
		{Title: "Level", Width: 12},
		{Title: "Years", Width: 8},
		{Title: "Events", Width: 8},
	}
	rowFormatter := func(skill display.Skill, _ int) []string {
		name := skill.Name
		category := skill.Category
		if category == "" {
			category = "-"
		}
		level := skill.Level
		if level == "" {
			level = "-"
		}
		years := "-"
		if skill.YearsUsed != nil && *skill.YearsUsed > 0 {
			years = strconv.Itoa(*skill.YearsUsed)
		}
		eventCount := "-"
		if count, ok := eventCounts[skill.ID]; ok && count > 0 {
			eventCount = strconv.Itoa(count)
		}
		return []string{name, category, level, years, eventCount}
	}
	list := behaviors.NewListBehavior(columns, rowFormatter, badgeFuncs)
	list.PaginationPrefix("Skills")
	list.EmptyMessage("No skills found. Press 'a' to add your first skill.")
	list.SetItems(skills)
	return &ListView{
		skills:      skills,
		list:        list,
		keys:        navigation.DefaultListKeyMap(),
		globalKeys:  navigation.DefaultGlobalKeyMap(),
		eventCounts: eventCounts,
	}
}

// Init returns nil (no async loading needed).
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (v *ListView) Init() tea.Cmd {
	return nil
}

// Update handles messages and returns a ViewResult on user action.
//
// Expected:
//   - msg is a valid tea.Msg (key press, window resize, or custom message).
//
// Returns:
//   - A tea.Cmd and widgets.ViewResult pair.
//
// Side effects:
//   - Updates internal view state based on message type.
func (v *ListView) Update(msg tea.Msg) (tea.Cmd, widgets.ViewResult) {
	return shared.HandleListViewUpdate(&v.BaseView, msg, v.handleKeyMessage)
}

func (v *ListView) handleKeyMessage(msg tea.KeyMsg) (tea.Cmd, widgets.ViewResult) {
	if result := v.handleGlobalKeys(msg); result != nil {
		return nil, result
	}
	if result := v.handleListKeys(msg); result != nil {
		return nil, result
	}
	return nil, v.handleDomainKeys(msg)
}

func (v *ListView) handleGlobalKeys(msg tea.KeyMsg) widgets.ViewResult {
	switch {
	case key.Matches(msg, v.globalKeys.Help):
		return &widgets.NavigateViewResult{ResultData: Nav{Action: ActionHelp}}
	case key.Matches(msg, v.globalKeys.Back):
		return &widgets.CancelViewResult{}
	}
	return nil
}

func (v *ListView) handleListKeys(msg tea.KeyMsg) widgets.ViewResult {
	if !v.list.HandleKey(msg) {
		return nil
	}
	if key.Matches(msg, v.keys.Select) {
		selected := v.list.GetSelectedItem()
		if selected != nil {
			return &widgets.NavigateViewResult{ResultData: Nav{Action: ActionView, Skill: *selected}}
		}
	}
	return nil
}

func (v *ListView) handleDomainKeys(msg tea.KeyMsg) widgets.ViewResult {
	switch msg.String() {
	case "a":
		return &widgets.NavigateViewResult{ResultData: Nav{Action: ActionAdd}}
	case "e":
		selected := v.list.GetSelectedItem()
		if selected != nil {
			return &widgets.NavigateViewResult{ResultData: Nav{Action: ActionEdit, Skill: *selected}}
		}
	case "d":
		selected := v.list.GetSelectedItem()
		if selected != nil {
			return &widgets.NavigateViewResult{ResultData: Nav{Action: ActionDelete, Skill: *selected}}
		}
	case "f":
		return &widgets.NavigateViewResult{ResultData: Nav{Action: ActionFilter}}
	case "s":
		return &widgets.NavigateViewResult{ResultData: Nav{Action: ActionSort}}
	case "/":
		return &widgets.NavigateViewResult{ResultData: Nav{Action: ActionSearch}}
	case "i":
		return &widgets.NavigateViewResult{ResultData: Nav{Action: ActionInfer}}
	case "?":
		return &widgets.NavigateViewResult{ResultData: Nav{Action: ActionHelp}}
	}
	return nil
}

// RenderContent returns the table content string — no chrome.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (v *ListView) RenderContent() string {
	return v.list.RenderContent()
}

// HelpText returns the rendered key binding footer for this view.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (v *ListView) HelpText() string {
	var th theme.Theme
	if t, ok := v.GetTheme().(theme.Theme); ok && t != nil {
		th = t
	}
	v.list.SetTheme(th)
	return v.list.HelpText()
}

// SetEventCounts updates the eventCounts map and rebuilds the list formatter.
//
// Expected:
//   - Must be a valid string.
//
// Side effects:
//   - None.
func (v *ListView) SetEventCounts(counts map[string]int) {
	for k, val := range counts {
		v.eventCounts[k] = val
	}
	rowFormatter := func(skill display.Skill, _ int) []string {
		name := skill.Name
		category := skill.Category
		if category == "" {
			category = "-"
		}
		level := skill.Level
		if level == "" {
			level = "-"
		}
		years := "-"
		if skill.YearsUsed != nil && *skill.YearsUsed > 0 {
			years = strconv.Itoa(*skill.YearsUsed)
		}
		eventCount := "-"
		if count, ok := v.eventCounts[skill.ID]; ok && count > 0 {
			eventCount = strconv.Itoa(count)
		}
		return []string{name, category, level, years, eventCount}
	}
	columns := []behaviors.ColumnDef{
		{Title: "Name", Width: 25},
		{Title: "Category", Width: 15},
		{Title: "Level", Width: 12},
		{Title: "Years", Width: 8},
		{Title: "Events", Width: 8},
	}
	v.list = behaviors.NewListBehavior(columns, rowFormatter, badgeFuncs)
	v.list.PaginationPrefix("Skills")
	v.list.EmptyMessage("No skills found. Press 'a' to add your first skill.")
	v.list.SetItems(v.skills)
}

// SetItems updates the skills slice and list items.
//
// Expected:
//   - skill must be valid.
//
// Side effects:
//   - None.
func (v *ListView) SetItems(skills []display.Skill) {
	v.skills = skills
	v.list.SetItems(skills)
}

// GetSelectedIndex returns the current table selection index (for testing).
//
// Returns:
//   - A int value.
//
// Side effects:
//   - None.
func (v *ListView) GetSelectedIndex() int {
	return v.list.GetSelectedIndex()
}
