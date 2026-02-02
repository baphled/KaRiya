package timeline

import (
	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

// TimelineEventListState identifies the timeline event list view in the
// state matrix. On this screen the user sees a paginated table of career
// events with columns for date, truncated event text, company, and
// project. Arrow keys or j/k navigate the highlighted row, and Ctrl+D/U
// or PgDn/PgUp page through larger lists. Pressing Enter opens the
// selected event detail view. Action keys allow adding an event (a),
// editing the selected event (e), deleting it (d), or opening the filter
// screen (f). Escape returns to the main menu.
const TimelineEventListState = "timeline_event_list"

// eventRowFormatter formats a career event for table display.
func eventRowFormatter(event *career.Event, _ int) []string {
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

// EventListScreen displays a list of career events in chronological order.
//
// This screen provides:
// - List navigation with ↑/↓ or j/k (vim-style)
// - Event selection with Enter key
// - Actions: Add (a), Edit (e), Delete (d)
// - Cancel with Escape or q
// - Event count display
// - Empty state handling
//
// Usage:
//
//	screen := browsetimeline.NewTimelineEventListScreen(events)
//	cmd, result := screen.Update(msg)
//	if result != nil && result.Type() == screens.ResultNavigate {
//	    navResult := result.(*screens.NavigateResult)
//	    if action, ok := navResult.ResultData.(map[string]interface{}); ok {
//	        // Handle action (add, edit, delete)
//	    } else if event, ok := navResult.ResultData.(*career.Event); ok {
//	        // User selected an event to view details
//	    }
//	}
//
// Related:
// - Screen provides the foundation
// - docs/TUI_DEVELOPER_GUIDE.md (Screen patterns)
// - docs/TUI_STANDARDS.md (Keyboard shortcuts).
type EventListScreen struct {
	*base.Screen
	events        []*career.Event
	tableBehavior *behaviors.TableBehavior[*career.Event]
}

// NewTimelineEventListScreen creates a new timeline event list screen.
//
// Expected:
//   - event must be valid.
//
// Returns:
//   - A fully initialized EventListScreen ready for use.
//
// Side effects:
//   - None.
func NewTimelineEventListScreen(events []*career.Event) *EventListScreen {
	columns := []behaviors.ColumnDef{
		{Title: "Date", Width: 12},
		{Title: "Event", Width: 40},
		{Title: "Company", Width: 18},
		{Title: "Project", Width: 15},
	}

	tableBehavior := behaviors.NewTableBehavior[*career.Event](nil, columns, eventRowFormatter).
		PageSize(15).
		PaginationPrefix("Events").
		EmptyMessage("No events found.")

	tableBehavior.SetItems(events)

	screen := &EventListScreen{
		Screen:        base.NewBaseScreen(),
		events:        events,
		tableBehavior: tableBehavior,
	}

	return screen
}

// Update handles messages and navigation.
//
// Expected:
//   - msg must be a valid tea.Msg type.
//
// Returns:
//   - tea.Cmd: command to execute.
//   - screens.ScreenResult: result indicating navigation or action.
//
// Side effects:
//   - May update table selection.
//   - May return CancelResult or NavigateResult.
func (s *EventListScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.SetTerminalInfo(msg.Width, msg.Height)
		return nil, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return nil, &screens.CancelResult{}

		case "up", "k", "down", "j", "ctrl+d", "ctrl+u", "pgup", "pgdown", "home", "end", "g", "G":
			s.tableBehavior.HandleNavigation(msg.String())
			return nil, nil

		case "enter":
			if selected := s.tableBehavior.GetSelectedItem(); selected != nil {
				return nil, &screens.NavigateResult{
					ResultData: *selected,
				}
			}
			return nil, nil

		case "a":
			return nil, &screens.NavigateResult{
				ResultData: map[string]interface{}{
					"action": "add",
				},
			}

		case "e":
			if selected := s.tableBehavior.GetSelectedItem(); selected != nil {
				return nil, &screens.NavigateResult{
					ResultData: map[string]interface{}{
						"action": "edit",
						"event":  *selected,
					},
				}
			}
			return nil, nil

		case "d":
			if selected := s.tableBehavior.GetSelectedItem(); selected != nil {
				return nil, &screens.NavigateResult{
					ResultData: map[string]interface{}{
						"action": "delete",
						"event":  *selected,
					},
				}
			}
			return nil, nil

		case "f":
			return nil, &screens.NavigateResult{
				ResultData: map[string]interface{}{
					"action": "filter",
				},
			}
		}
	}

	return nil, nil
}

// SetContentHeight configures the table to use the specified height with viewport scrolling.
//
// Expected:
//   - int must be valid.
//
// Side effects:
//   - None.
func (s *EventListScreen) SetContentHeight(height int) {
	if height < 10 {
		return
	}

	s.tableBehavior.SetHeight(height)
}

// RenderContent returns just the content (table) without StandardView wrapper.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (s *EventListScreen) RenderContent() string {
	return s.tableBehavior.Render()
}

// View renders the event list screen using StandardView with table.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (s *EventListScreen) View() string {
	content := s.RenderContent()

	var th themes.Theme
	if screenTheme := s.Theme(); screenTheme != nil {
		if t, ok := screenTheme.(themes.Theme); ok {
			th = t
		}
	}
	if th == nil {
		th = themes.NewDefaultTheme()
	}

	footer := primitives.RenderHelpFooter(th,
		primitives.NavigateBadge(th),
		primitives.HelpKeyBadge("Ctrl+D/U", "Page", th),
		primitives.HelpKeyBadge("Enter", "View", th),
		primitives.AddBadge(th),
		primitives.EditBadge(th),
		primitives.DeleteBadge(th),
		primitives.BackBadge(th),
		primitives.QuitBadge(th),
		primitives.HelpBadge(th),
	)

	breadcrumbs := []string{"Main Menu", "Timeline"}
	return s.CreateView(breadcrumbs, content, footer)
}

// SetTheme applies theme to the table (override Screen).
//
// Expected:
//   - interface{} must be valid.
//
// Side effects:
//   - None.
func (s *EventListScreen) SetTheme(theme interface{}) {
	s.Screen.SetTheme(theme)
	if t, ok := theme.(themes.Theme); ok && t != nil {
		s.tableBehavior.SetTheme(t)
	}
}

// GetEvents returns the list of events.
//
// Returns:
//   - A []*career.Event value.
//
// Side effects:
//   - None.
func (s *EventListScreen) GetEvents() []*career.Event {
	return s.events
}

// GetSelectedIndex returns the currently selected index.
//
// Returns:
//   - A int value.
//
// Side effects:
//   - None.
func (s *EventListScreen) GetSelectedIndex() int {
	return s.tableBehavior.GetSelectedIndex()
}
