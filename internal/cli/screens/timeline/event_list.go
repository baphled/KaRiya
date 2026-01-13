package timeline

import (
	"fmt"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

// State constant for state matrix tracking (REQUIRED)
const TimelineEventListState = "timeline_event_list"

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
	selectedIndex int
	table         table.Model
	listContainer *components.TableListContainer
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
	// Create table columns matching legacy format
	columns := []table.Column{
		{Title: "Date", Width: 12},
		{Title: "Event", Width: 50},
		{Title: "Company", Width: 20},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(15),
		table.WithWidth(100),
	)

	// Apply default styles - theme will be applied via SetTheme
	t.SetStyles(table.DefaultStyles())

	screen := &TimelineEventListScreen{
		BaseScreen:    base.NewBaseScreen(),
		events:        events,
		selectedIndex: 0,
		table:         t,
		listContainer: components.NewTableListContainer(t, "Career Timeline", 100),
	}

	// Update table rows with events
	screen.updateTableRows()

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
		case "esc", "q":
			// Cancel and return to main menu
			return nil, &screens.CancelResult{}

		case "up", "k":
			// Move selection up
			if s.selectedIndex > 0 {
				s.selectedIndex--
			}
			return nil, nil

		case "down", "j":
			// Move selection down
			if len(s.events) > 0 && s.selectedIndex < len(s.events)-1 {
				s.selectedIndex++
			}
			return nil, nil

		case "enter":
			// View event details
			if len(s.events) > 0 && s.selectedIndex < len(s.events) {
				return nil, &screens.NavigateResult{
					ResultData: s.events[s.selectedIndex],
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
			if len(s.events) > 0 && s.selectedIndex < len(s.events) {
				return nil, &screens.NavigateResult{
					ResultData: map[string]interface{}{
						"action": "edit",
						"event":  s.events[s.selectedIndex],
					},
				}
			}
			return nil, nil

		case "d":
			// Delete selected event
			if len(s.events) > 0 && s.selectedIndex < len(s.events) {
				return nil, &screens.NavigateResult{
					ResultData: map[string]interface{}{
						"action": "delete",
						"event":  s.events[s.selectedIndex],
					},
				}
			}
			return nil, nil
		}
	}

	return nil, nil
}

// updateTableRows updates the table rows based on events (with pagination).
func (s *TimelineEventListScreen) updateTableRows() {
	pageSize := 15
	total := len(s.events)

	// Determine which page current selection is on
	page := 0
	if pageSize > 0 && s.selectedIndex >= 0 {
		page = s.selectedIndex / pageSize
	}

	start := page * pageSize
	end := start + pageSize
	if end > total {
		end = total
	}

	pageEvents := s.events[start:end]

	rows := make([]table.Row, 0, len(pageEvents))
	for idx, event := range pageEvents {
		realIdx := start + idx

		// Date formatting
		dateStr := event.Date.Format("2006-01-02")

		// Add selection indicator for selected row
		if realIdx == s.selectedIndex {
			dateStr = "▶ " + dateStr
		} else {
			dateStr = "  " + dateStr
		}

		// Truncate text to 50 chars (matching legacy)
		text := event.Text
		if len(text) > 50 {
			text = text[:50] + "..."
		}

		// Company (or dash if empty)
		company := event.Company
		if company == "" {
			company = "-"
		}

		rows = append(rows, table.Row{dateStr, text, company})
	}

	s.table.SetRows(rows)

	// Calculate relative cursor position for this page
	relativeCursor := 0
	if s.selectedIndex >= start && s.selectedIndex < end {
		relativeCursor = s.selectedIndex - start
	}

	// Set table cursor to relative position
	s.table.SetCursor(relativeCursor)

	// Sync container's selected index
	s.listContainer.SetSelectedIdx(relativeCursor)
}

// View renders the event list screen using StandardView with table.
func (s *TimelineEventListScreen) View() string {
	// Handle empty state
	if len(s.events) == 0 {
		s.listContainer.SetEmptyStateMessage("No events found.")
		content := s.listContainer.Render()
		footer := "a: Add event  Esc/q: Back"
		return s.CreateView([]string{"Main Menu", "Timeline"}, content, footer)
	}

	// Ensure table rows are synchronized
	s.updateTableRows()

	// Build pagination info matching legacy format
	pageSize := 15
	totalItems := len(s.events)
	currentPage := (s.selectedIndex / pageSize) + 1
	totalPages := (totalItems + pageSize - 1) / pageSize
	paginationInfo := fmt.Sprintf("Events: %d | Page %d of %d", totalItems, currentPage, totalPages)
	s.listContainer.SetPaginationInfo(paginationInfo)

	// Render table via container
	content := s.listContainer.Render()

	// Footer with actions (matching legacy)
	footer := "↑/↓ or j/k: Navigate  Enter: View details  a: Add  e: Edit  d: Delete  Esc/q: Back"

	// Use BaseScreen's CreateView helper for StandardView integration
	breadcrumbs := []string{"Main Menu", "Timeline"}
	return s.CreateView(breadcrumbs, content, footer)
}

// SetTheme applies theme to the table (override BaseScreen).
func (s *TimelineEventListScreen) SetTheme(theme interface{}) {
	s.BaseScreen.SetTheme(theme)
	// Apply themed table styles if theme is available
	if t, ok := theme.(themes.Theme); ok && t != nil {
		s.table.SetStyles(themes.NewThemedTableStyles(t))
	}
}

// GetEvents returns the list of events.
func (s *TimelineEventListScreen) GetEvents() []*career.CareerEvent {
	return s.events
}

// GetSelectedIndex returns the currently selected index.
func (s *TimelineEventListScreen) GetSelectedIndex() int {
	return s.selectedIndex
}
