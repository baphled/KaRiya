package event

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/baphled/kariya/internal/tui/navigation"
	"github.com/baphled/kariya/internal/tui/views/shared"
	"github.com/baphled/kariya/internal/ui/behaviors"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/uikit/primitives"
	"github.com/baphled/kariya/internal/ui/uikit/theme"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
)

// ListView displays a paginated list of career events.
// It implements widgets.View — the intent owns chrome (breadcrumbs, logo, footer).
type ListView struct {
	widgets.BaseView
	events     []display.Event
	list       *behaviors.ListBehavior[display.Event]
	keys       navigation.ListKeyMap
	globalKeys navigation.GlobalKeyMap
}

// NewListView creates a new ListView with the given events.
//
// Expected:
//   - event must be valid.
//
// Returns:
//   - A fully initialized ListView ready for use.
//
// Side effects:
//   - None.
func NewListView(events []display.Event) *ListView {
	columns := []behaviors.ColumnDef{
		{Title: "Date", Width: 12},
		{Title: "Event", Width: 40},
		{Title: "Company", Width: 18},
		{Title: "Project", Width: 15},
	}
	badgeFuncs := []behaviors.BadgeFunc{
		primitives.NavigateBadge,
		func(th theme.Theme) *primitives.Badge { return primitives.HelpKeyBadge("Ctrl+D/U", "Page", th) },
		func(th theme.Theme) *primitives.Badge { return primitives.HelpKeyBadge("Enter", "View", th) },
		primitives.AddBadge,
		primitives.EditBadge,
		primitives.DeleteBadge,
		primitives.BackBadge,
		primitives.QuitBadge,
		primitives.HelpBadge,
	}
	list := behaviors.NewListBehavior(columns, eventRowFormatter, badgeFuncs)
	list.PaginationPrefix("Events")
	list.EmptyMessage("No events found.")
	list.SetItems(events)
	return &ListView{
		events:     events,
		list:       list,
		keys:       navigation.DefaultListKeyMap(),
		globalKeys: navigation.DefaultGlobalKeyMap(),
	}
}

// eventRowFormatter formats a career event for table display.
func eventRowFormatter(event display.Event, _ int) []string {
	dateStr := event.Date.Format("2006-01-02")

	text := event.Text
	if len(text) > 40 {
		text = text[:40] + "..."
	}

	company := event.Company
	if company == "" {
		company = "-"
	}

	project := event.Project
	if project == "" {
		project = "-"
	}

	return []string{dateStr, text, company, project}
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

// Update handles tea.Msg and returns (cmd, ViewResult).
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

// handleKeyMessage processes keyboard input using key.Matches with keymaps.
func (v *ListView) handleKeyMessage(msg tea.KeyMsg) (tea.Cmd, widgets.ViewResult) {
	if result := v.handleGlobalKeys(msg); result != nil {
		return nil, result
	}

	if v.list.HandleKey(msg) {
		if key.Matches(msg, v.keys.Select) {
			selected := v.list.GetSelectedItem()
			if selected != nil {
				return nil, &widgets.NavigateViewResult{ResultData: Nav{Action: ActionView, Event: *selected}}
			}
		}
		return nil, nil
	}

	return nil, v.handleDomainActionKeys(msg)
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

func (v *ListView) handleDomainActionKeys(msg tea.KeyMsg) widgets.ViewResult {
	switch {
	case key.Matches(msg, v.keys.Add):
		return &widgets.NavigateViewResult{ResultData: Nav{Action: ActionAdd}}
	case key.Matches(msg, v.keys.Edit):
		return v.navigateWithSelectedEvent(ActionEdit)
	case key.Matches(msg, v.keys.Delete):
		return v.navigateWithSelectedEvent(ActionDelete)
	case key.Matches(msg, v.keys.Filter):
		return &widgets.NavigateViewResult{ResultData: Nav{Action: ActionFilter}}
	case key.Matches(msg, v.keys.Search):
		return &widgets.NavigateViewResult{ResultData: Nav{Action: ActionSearch}}
	case key.Matches(msg, v.keys.Sort):
		return &widgets.NavigateViewResult{ResultData: Nav{Action: ActionSort}}
	case key.Matches(msg, v.keys.Clear):
		return &widgets.NavigateViewResult{ResultData: Nav{Action: ActionClear}}
	case key.Matches(msg, v.keys.ShowSkills):
		return v.navigateWithSelectedSkillAction(ActionShowEventSkills)
	case key.Matches(msg, v.keys.InferSkills):
		return v.navigateWithSelectedSkillAction(ActionInferSkills)
	}
	return nil
}

func (v *ListView) navigateWithSelectedEvent(action ActionKey) widgets.ViewResult {
	selected := v.list.GetSelectedItem()
	if selected != nil {
		return &widgets.NavigateViewResult{ResultData: Nav{Action: action, Event: *selected}}
	}
	return nil
}

func (v *ListView) navigateWithSelectedSkillAction(action SkillActionKey) widgets.ViewResult {
	selected := v.list.GetSelectedItem()
	if selected != nil {
		return &widgets.NavigateViewResult{ResultData: SkillNav{Action: action, Event: *selected}}
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

// GetEvents returns the events list (for testing).
//
// Returns:
//   - A []display.Event value.
//
// Side effects:
//   - None.
func (v *ListView) GetEvents() []display.Event {
	return v.events
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
