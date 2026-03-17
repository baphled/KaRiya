package burst

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/baphled/kariya/internal/tui/navigation"
	"github.com/baphled/kariya/internal/ui/behaviors"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/uikit/primitives"
	"github.com/baphled/kariya/internal/ui/uikit/theme"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
)

// List displays a paginated list of career bursts.
// It implements widgets.View — the intent owns chrome (breadcrumbs, logo, footer).
type List struct {
	widgets.BaseView
	bursts     []display.Burst
	list       *behaviors.ListBehavior[display.Burst]
	keys       navigation.ListKeyMap
	globalKeys navigation.GlobalKeyMap
}

// NewList creates a new List with the given bursts.
//
// Expected:
//   - bursts must be a valid slice (may be empty or nil).
//
// Returns:
//   - A fully initialized List ready for use.
//
// Side effects:
//   - None.
func NewList(bursts []display.Burst) *List {
	columns := []behaviors.ColumnDef{
		{Title: "Name", Width: 30},
		{Title: "Description", Width: 35},
		{Title: "Confirmed", Width: 10},
		{Title: "Events", Width: 8},
		{Title: "Created", Width: 12},
	}
	badgeFuncs := []behaviors.BadgeFunc{
		primitives.NavigateBadge,
		func(th theme.Theme) *primitives.Badge { return primitives.HelpKeyBadge("Ctrl+D/U", "Page", th) },
		func(th theme.Theme) *primitives.Badge { return primitives.HelpKeyBadge("Enter", "View", th) },
		primitives.AddBadge,
		primitives.EditBadge,
		primitives.DeleteBadge,
		func(th theme.Theme) *primitives.Badge { return primitives.HelpKeyBadge("s", "Suggest", th) },
		primitives.BackBadge,
	}
	list := behaviors.NewListBehavior(columns, burstRowFormatter, badgeFuncs)
	list.PaginationPrefix("Bursts")
	list.EmptyMessage("No bursts found. Press 'a' to add or 's' to suggest.")
	list.SetItems(bursts)
	return &List{
		bursts:     bursts,
		list:       list,
		keys:       navigation.DefaultListKeyMap(),
		globalKeys: navigation.DefaultGlobalKeyMap(),
	}
}

// burstRowFormatter formats a burst for table display.
func burstRowFormatter(burst display.Burst, _ int) []string {
	nameStr := burst.Name
	if len(nameStr) > 27 {
		nameStr = nameStr[:27] + "..."
	}

	descStr := strings.TrimSpace(burst.Description)
	descStr = strings.ReplaceAll(descStr, "\n", " ")
	descStr = strings.ReplaceAll(descStr, "\r", " ")
	if descStr == "" {
		descStr = "-"
	} else if len(descStr) > 32 {
		descStr = descStr[:32] + "..."
	}

	confirmedStr := "✗ No"
	if burst.Confirmed {
		confirmedStr = "✓ Yes"
	}

	eventCount := strconv.Itoa(len(burst.EventIDs))
	createdStr := burst.CreatedAt.Format("2006-01-02")

	return []string{nameStr, descStr, confirmedStr, eventCount, createdStr}
}

// Init returns nil (no async loading needed).
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (v *List) Init() tea.Cmd {
	return nil
}

// Update handles tea.Msg and returns (cmd, ViewResult).
// Returns nil ViewResult for internal state changes.
// Returns non-nil ViewResult when navigation/cancel occurs.
//
// Expected:
//   - msg is a valid tea.Msg (may be any message type).
//
// Returns:
//   - A tea.Cmd and widgets.ViewResult value.
//
// Side effects:
//   - None.
func (v *List) Update(msg tea.Msg) (tea.Cmd, widgets.ViewResult) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.SetTerminalInfo(msg.Width, msg.Height)
		return nil, nil
	case tea.KeyMsg:
		return v.handleKeyMessage(msg)
	}
	return nil, nil
}

// handleKeyMessage processes keyboard input using key.Matches with keymaps.
func (v *List) handleKeyMessage(msg tea.KeyMsg) (tea.Cmd, widgets.ViewResult) {
	if result := v.handleGlobalKeys(msg); result != nil {
		return nil, result
	}

	if v.list.HandleKey(msg) {
		if key.Matches(msg, v.keys.Select) {
			selected := v.list.GetSelectedItem()
			if selected != nil {
				return nil, &widgets.NavigateViewResult{ResultData: Nav{Action: ActionView, Burst: *selected}}
			}
		}
		return nil, nil
	}

	return nil, v.handleActionKeys(msg)
}

// handleGlobalKeys handles escape and help keys.
func (v *List) handleGlobalKeys(msg tea.KeyMsg) widgets.ViewResult {
	switch {
	case key.Matches(msg, v.globalKeys.Help):
		return &widgets.NavigateViewResult{ResultData: Nav{Action: ActionHelp}}
	case key.Matches(msg, v.globalKeys.Back):
		return &widgets.CancelViewResult{}
	}
	return nil
}

// handleActionKeys handles domain-specific action keys.
func (v *List) handleActionKeys(msg tea.KeyMsg) widgets.ViewResult {
	switch {
	case key.Matches(msg, v.keys.Add):
		return &widgets.NavigateViewResult{ResultData: Nav{Action: ActionAdd}}
	case key.Matches(msg, v.keys.Edit):
		return v.navigateWithSelectedBurst(ActionEdit)
	case key.Matches(msg, v.keys.Delete):
		return v.navigateWithSelectedBurst(ActionDelete)
	}

	if msg.String() == "s" {
		return &widgets.NavigateViewResult{ResultData: Nav{Action: ActionSuggest}}
	}

	return nil
}

// navigateWithSelectedBurst returns a navigate result for the selected burst.
func (v *List) navigateWithSelectedBurst(action ActionKey) widgets.ViewResult {
	selected := v.list.GetSelectedItem()
	if selected != nil {
		return &widgets.NavigateViewResult{ResultData: Nav{Action: action, Burst: *selected}}
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
func (v *List) RenderContent() string {
	return v.list.RenderContent()
}

// HelpText returns the rendered key binding footer for this view.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (v *List) HelpText() string {
	var th theme.Theme
	if t, ok := v.GetTheme().(theme.Theme); ok && t != nil {
		th = t
	}
	v.list.SetTheme(th)
	return v.list.HelpText()
}

// GetBursts returns the bursts list (for testing).
//
// Returns:
//   - A []display.Burst value.
//
// Side effects:
//   - None.
func (v *List) GetBursts() []display.Burst {
	return v.bursts
}

// GetSelectedIndex returns the current table selection index (for testing).
//
// Returns:
//   - A int value.
//
// Side effects:
//   - None.
func (v *List) GetSelectedIndex() int {
	return v.list.GetSelectedIndex()
}
