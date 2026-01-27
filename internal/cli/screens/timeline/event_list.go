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

// State constant for state matrix tracking (REQUIRED)
const TimelineEventListState = "timeline_event_list"

// eventRowFormatter formats a career event for table display
func eventRowFormatter(event *career.CareerEvent, _ int) []string {
	// Date formatting
	dateStr := event.Date.Format("2006-01-02")

	// Truncate text to 40 chars (reduced to make room for project)
	text := event.Text
	if len(text) > 40 {
		text = text[:40] + "..."
	}

	// Company (or dash if empty)
	company := event.Company
	if company == "" {
		company = "-"
	}

	// Project (or dash if empty)
	project := event.Project
	if project == "" {
		project = "-"
	}

	return []string{dateStr, text, company, project}
}

// TimelineEventListScreen displays a list of career events in chronological order.
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
//	screen := timeline.NewTimelineEventListScreen(events)
//	cmd, result := screen.Update(msg)
//	if result != nil && result.Type() == screens.ResultNavigate {
//	    navResult := result.(*screens.NavigateResult)
//	    if action, ok := navResult.ResultData.(map[string]interface{}); ok {
//	        // Handle action (add, edit, delete)
//	    } else if event, ok := navResult.ResultData.(*career.CareerEvent); ok {
//	        // User selected an event to view details
//	    }
//	}
//
// Related:
// - BaseScreen provides the foundation
// - docs/TUI_DEVELOPER_GUIDE.md (Screen patterns)
// - docs/TUI_STANDARDS.md (Keyboard shortcuts)
type TimelineEventListScreen struct {
	*base.BaseScreen
	events        []*career.CareerEvent
	tableBehavior *behaviors.TableBehavior[*career.CareerEvent]
}

// NewTimelineEventListScreen creates a new timeline event list screen.
//
// The screen:
// - Displays events in a scrollable list
// - Supports keyboard navigation (↑/↓, j/k)
// - Shows event date, text, and company
// - Provides actions (add, edit, delete, view)
//
// Parameters:
//   - events: List of career events to display (can be empty)
func NewTimelineEventListScreen(events []*career.CareerEvent) *TimelineEventListScreen {
	// Create table columns
	columns := []behaviors.ColumnDef{
		{Title: "Date", Width: 12},
		{Title: "Event", Width: 40},
		{Title: "Company", Width: 18},
		{Title: "Project", Width: 15},
	}

	// Create TableBehavior for events list (theme set later via SetTheme)
	tableBehavior := behaviors.NewTableBehavior[*career.CareerEvent](nil, columns, eventRowFormatter).
		PageSize(15).
		PaginationPrefix("Events").
		EmptyMessage("No events found.")

	// Set items after creation
	tableBehavior.SetItems(events)

	screen := &TimelineEventListScreen{
		BaseScreen:    base.NewBaseScreen(),
		events:        events,
		tableBehavior: tableBehavior,
	}

	return screen
}

// Update handles messages and navigation.
func (s *TimelineEventListScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.SetTerminalInfo(msg.Width, msg.Height)
		return nil, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// Cancel and return to main menu
			// Note: 'q' (quit) is handled by the intent before delegation
			return nil, &screens.CancelResult{}

		case "up", "k", "down", "j", "ctrl+d", "ctrl+u", "pgup", "pgdown", "home", "end", "g", "G":
			// Delegate navigation to TableBehavior
			s.tableBehavior.HandleNavigation(msg.String())
			return nil, nil

		case "enter":
			// View event details
			if selected := s.tableBehavior.GetSelectedItem(); selected != nil {
				return nil, &screens.NavigateResult{
					ResultData: *selected, // Dereference **T to get *T
				}
			}
			return nil, nil

		case "a":
			// Add new event
			return nil, &screens.NavigateResult{
				ResultData: map[string]interface{}{
					"action": "add",
				},
			}

		case "e":
			// Edit selected event
			if selected := s.tableBehavior.GetSelectedItem(); selected != nil {
				return nil, &screens.NavigateResult{
					ResultData: map[string]interface{}{
						"action": "edit",
						"event":  *selected, // Dereference **T to get *T
					},
				}
			}
			return nil, nil

		case "d":
			// Delete selected event
			if selected := s.tableBehavior.GetSelectedItem(); selected != nil {
				return nil, &screens.NavigateResult{
					ResultData: map[string]interface{}{
						"action": "delete",
						"event":  *selected, // Dereference **T to get *T
					},
				}
			}
			return nil, nil

		case "f":
			// Open filter screen
			return nil, &screens.NavigateResult{
				ResultData: map[string]interface{}{
					"action": "filter",
				},
			}
		}
	}

	return nil, nil
}

// RenderContent returns just the content (table) without StandardView wrapper.
// This allows the intent to wrap it with proper breadcrumbs and themed footer.
func (s *TimelineEventListScreen) RenderContent() string {
	// TableBehavior handles empty state and pagination internally
	return s.tableBehavior.Render()
}

// View renders the event list screen using StandardView with table.
func (s *TimelineEventListScreen) View() string {
	content := s.RenderContent()

	// Get theme for UIKit primitives (fall back to default if not set).
	var th themes.Theme
	if screenTheme := s.Theme(); screenTheme != nil {
		if t, ok := screenTheme.(themes.Theme); ok {
			th = t
		}
	}
	if th == nil {
		th = themes.NewDefaultTheme()
	}

	// Build footer using UIKit primitives for consistent styling.
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

	// Use BaseScreen's CreateView helper for StandardView integration.
	breadcrumbs := []string{"Main Menu", "Timeline"}
	return s.CreateView(breadcrumbs, content, footer)
}

// SetTheme applies theme to the table (override BaseScreen).
func (s *TimelineEventListScreen) SetTheme(theme interface{}) {
	s.BaseScreen.SetTheme(theme)
	// Apply themed table styles if theme is available
	if t, ok := theme.(themes.Theme); ok && t != nil {
		s.tableBehavior.SetTheme(t)
	}
}

// GetEvents returns the list of events.
func (s *TimelineEventListScreen) GetEvents() []*career.CareerEvent {
	return s.events
}

// GetSelectedIndex returns the currently selected index.
func (s *TimelineEventListScreen) GetSelectedIndex() int {
	return s.tableBehavior.GetSelectedIndex()
}
