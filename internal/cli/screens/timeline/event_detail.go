package timeline

import (
	"fmt"

	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/cli/uikit/theme"
	"github.com/baphled/kariya/internal/cli/uikit/widgets"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

// TimelineEventDetailState identifies the timeline event detail view in the
// state matrix. On this screen the user sees a DetailView card displaying
// the event date, company, project, and full text body, followed by
// optional lists of tags, categories, and an associated skills count.
// Pressing e opens the event editor, d starts the delete confirmation, and
// Escape or Backspace returns to the event list.
const TimelineEventDetailState = "timeline_event_detail"

// EventDetailScreen displays detailed information about a career event.
//
// This screen provides:
// - Detailed view of event date, text, company, project
// - Display of tags, categories, and associated skills
// - Actions: Edit (e), Delete (d)
// - Back navigation with Escape, q, or backspace
//
// Usage:
//
//	screen := browsetimeline.NewTimelineEventDetailScreen(event)
//	cmd, result := screen.Update(msg)
//	if result != nil && result.Type() == screens.ResultNavigate {
//	    navResult := result.(*screens.NavigateResult)
//	    if action, ok := navResult.ResultData.(map[string]interface{}); ok {
//	        // Handle action (edit, delete)
//	    }
//	}
//
// Related:
// - Screen provides the foundation
// - docs/TUI_DEVELOPER_GUIDE.md (Screen patterns)
// - docs/TUI_STANDARDS.md (Keyboard shortcuts).
type EventDetailScreen struct {
	*base.Screen
	event *career.Event
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
func NewTimelineEventDetailScreen(event *career.Event) *EventDetailScreen {
	return &EventDetailScreen{
		Screen: base.NewBaseScreen(),
		event:  event,
	}
}

// Update handles messages and actions.
func (s *EventDetailScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.SetTerminalInfo(msg.Width, msg.Height)
		return nil, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "backspace":
			return nil, &screens.CancelResult{}

		case "e":
			return nil, &screens.NavigateResult{
				ResultData: map[string]interface{}{
					"action": "edit",
					"event":  s.event,
				},
			}

		case "d":
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
func (s *EventDetailScreen) RenderContent() string {
	if s.event == nil {
		return "No event selected."
	}

	var th theme.Theme
	if screenTheme := s.Theme(); screenTheme != nil {
		if t, ok := screenTheme.(theme.Theme); ok {
			th = t
		}
	}
	if th == nil {
		th = theme.Default()
	}

	dv := widgets.NewDetailView(th).
		Title("Event Details").
		Field("Date", s.event.Date.Format("2006-01-02")).
		FieldIf("Company", s.event.Company).
		FieldIf("Project", s.event.Project).
		Field("Text", s.event.Text)

	if len(s.event.Tags) > 0 {
		dv.List("Tags", s.event.Tags)
	}
	if len(s.event.Categories) > 0 {
		dv.List("Categories", s.event.Categories)
	}
	if len(s.event.Skills) > 0 {
		dv.Field("Skills", fmt.Sprintf("%d associated", len(s.event.Skills)))
	}

	return dv.Render()
}

// View renders the event detail screen using StandardView.
func (s *EventDetailScreen) View() string {
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
		primitives.EditBadge(th),
		primitives.DeleteBadge(th),
		primitives.HelpKeyBadge("Esc/Backspace", "Back", th),
		primitives.QuitBadge(th),
		primitives.HelpBadge(th),
	)

	breadcrumbs := []string{"Main Menu", "Timeline", "Event Details"}
	return s.CreateView(breadcrumbs, content, footer)
}

// GetEvent returns the event being displayed.
func (s *EventDetailScreen) GetEvent() *career.Event {
	return s.event
}
