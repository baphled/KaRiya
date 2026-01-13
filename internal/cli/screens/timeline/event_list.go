package timeline

import (
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	return &TimelineEventListScreen{
		BaseScreen:    base.NewBaseScreen(),
		events:        events,
		selectedIndex: 0,
	}
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

// View renders the event list screen using StandardView.
func (s *TimelineEventListScreen) View() string {
	var b strings.Builder

	// Header with event count
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	if len(s.events) == 0 {
		b.WriteString(headerStyle.Render("Career Timeline"))
		b.WriteString("\n\n")
		b.WriteString("No events found. Press 'a' to add your first event.")
	} else {
		countText := fmt.Sprintf("Career Timeline (%d events)", len(s.events))
		b.WriteString(headerStyle.Render(countText))
		b.WriteString("\n\n")

		// Render event list
		for i, event := range s.events {
			s.renderEventItem(&b, i, event)
		}
	}

	content := b.String()

	// Footer with actions
	footer := "↑/↓ or j/k: Navigate  Enter: View details  a: Add  e: Edit  d: Delete  Esc/q: Back"

	// Use BaseScreen's CreateView helper for StandardView integration
	breadcrumbs := []string{"Main Menu", "Timeline"}
	return s.CreateView(breadcrumbs, content, footer)
}

// renderEventItem renders a single event in the list.
func (s *TimelineEventListScreen) renderEventItem(b *strings.Builder, index int, event *career.CareerEvent) {
	// Date formatting
	dateStr := event.Date.Format("2006-01-02")

	// Selected indicator
	indicator := "  "
	if index == s.selectedIndex {
		indicator = "▶ "
	}

	// Style for selected item
	itemStyle := lipgloss.NewStyle()
	if index == s.selectedIndex {
		itemStyle = itemStyle.Bold(true).Foreground(lipgloss.Color("12"))
	}

	// Company tag if present
	companyTag := ""
	if event.Company != "" {
		companyTag = fmt.Sprintf(" [%s]", event.Company)
	}

	// Truncate text if too long
	text := event.Text
	maxTextLength := s.Width() - 20 - len(dateStr) - len(companyTag)
	if maxTextLength < 20 {
		maxTextLength = 20
	}
	if len(text) > maxTextLength {
		text = text[:maxTextLength-3] + "..."
	}

	// Format: "▶ 2024-01-01  Event text [Company]"
	line := fmt.Sprintf("%s%s  %s%s", indicator, dateStr, text, companyTag)
	b.WriteString(itemStyle.Render(line))
	b.WriteString("\n")
}

// GetEvents returns the list of events.
func (s *TimelineEventListScreen) GetEvents() []*career.CareerEvent {
	return s.events
}

// GetSelectedIndex returns the currently selected index.
func (s *TimelineEventListScreen) GetSelectedIndex() int {
	return s.selectedIndex
}
