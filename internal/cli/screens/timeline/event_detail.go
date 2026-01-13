package timeline

import (
	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

// State constant for state matrix tracking (REQUIRED)
const TimelineEventDetailState = "timeline_event_detail"

// TimelineEventDetailScreen displays detailed information about a career event.
//
// This screen provides:
// - Detailed view of event date, text, company, project
// - Display of tags, categories, and associated skills
// - Actions: Edit (e), Delete (d)
// - Back navigation with Escape, q, or backspace
//
// Usage:
//
//	screen := timeline.NewTimelineEventDetailScreen(event)
//	cmd, result := screen.Update(msg)
//	if result != nil && result.Type() == screens.ResultNavigate {
//	    navResult := result.(*screens.NavigateResult)
//	    if action, ok := navResult.ResultData.(map[string]interface{}); ok {
//	        // Handle action (edit, delete)
//	    }
//	}
//
// Related:
// - BaseScreen provides the foundation
// - docs/TUI_DEVELOPER_GUIDE.md (Screen patterns)
// - docs/TUI_STANDARDS.md (Keyboard shortcuts)
type TimelineEventDetailScreen struct {
	*base.BaseScreen
	event *career.CareerEvent
}

// NewTimelineEventDetailScreen creates a new event detail screen.
//
// The screen:
// - Displays all event fields in a formatted layout
// - Shows optional fields only when present
// - Provides edit and delete actions
// - Supports back navigation
//
// Parameters:
//   - event: The career event to display
func NewTimelineEventDetailScreen(event *career.CareerEvent) *TimelineEventDetailScreen {
	return &TimelineEventDetailScreen{
		BaseScreen: base.NewBaseScreen(),
		event:      event,
	}
}

// Update handles messages and actions.
func (s *TimelineEventDetailScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.SetTerminalInfo(msg.Width, msg.Height)
		return nil, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "backspace":
			// Back to event list
			// Note: 'q' (quit) is handled by the intent before delegation
			return nil, &screens.CancelResult{}

		case "e":
			// Edit event
			return nil, &screens.NavigateResult{
				ResultData: map[string]interface{}{
					"action": "edit",
					"event":  s.event,
				},
			}

		case "d":
			// Delete event
			return nil, &screens.NavigateResult{
				ResultData: map[string]interface{}{
					"action": "delete",
					"event":  s.event,
				},
			}
		}
	}

	return nil, nil
}

// RenderContent returns just the content (event detail card) without StandardView wrapper.
// This allows the intent to wrap it with proper breadcrumbs and themed footer.
func (s *TimelineEventDetailScreen) RenderContent() string {
	// Use the same component as legacy for consistency
	var themeObj themes.Theme
	if t := s.Theme(); t != nil {
		if th, ok := t.(themes.Theme); ok {
			themeObj = th
		}
	}
	return components.RenderEventDetailCard(s.event, themeObj)
}

// View renders the event detail screen using StandardView.
// This is kept for backward compatibility but RenderContent() is preferred
// when the intent manages the StandardView wrapper.
func (s *TimelineEventDetailScreen) View() string {
	content := s.RenderContent()

	// Footer with actions (matching legacy)
	// Note: 'q' is a global key handled by intent (quits app)
	footer := "e: Edit  d: Delete  Esc/Backspace: Back to timeline  q: Quit  ?: Help"

	// Use BaseScreen's CreateView helper for StandardView integration
	breadcrumbs := []string{"Main Menu", "Timeline", "Event Details"}
	return s.CreateView(breadcrumbs, content, footer)
}

// GetEvent returns the event being displayed.
func (s *TimelineEventDetailScreen) GetEvent() *career.CareerEvent {
	return s.event
}
