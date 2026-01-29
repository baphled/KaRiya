package burst_management

import (
	"fmt"

	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/cli/uikit/theme"
	"github.com/baphled/kariya/internal/cli/uikit/widgets"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

// BurstDetailState identifies the burst detail screen in the state matrix.
// On this screen the user sees a single burst's name, description, event
// count, confirmation status, and created/updated timestamps rendered inside
// a UIKit DetailView card. Action keys allow viewing associated events (v),
// viewing extracted facts (f), editing the burst (e), deleting it (d), or
// marking it as confirmed (c). Pressing Escape or Backspace returns to the
// burst list.
const BurstDetailState = "burst_detail"

// BurstDetailScreen displays detailed information about a career burst.
//
// This screen provides:
// - Detailed view of burst name, description, event count, confirmation status
// - Display of created/updated timestamps
// - Actions: View Events (v), View Facts (f), Edit (e), Delete (d), Confirm (c)
// - Back navigation with Escape or backspace
//
// Usage:
//
//	screen := burst_management.NewBurstDetailScreen(burst)
//	cmd, result := screen.Update(msg)
//	if result != nil && result.Type() == screens.ResultNavigate {
//	    navResult := result.(*screens.NavigateResult)
//	    if action, ok := navResult.ResultData.(map[string]interface{}); ok {
//	        // Handle action (view_events, view_facts, edit, delete, confirm)
//	    }
//	}
//
// Related:
// - Screen provides the foundation
// - docs/TUI_DEVELOPER_GUIDE.md (Screen patterns)
// - docs/TUI_STANDARDS.md (Keyboard shortcuts).
type BurstDetailScreen struct {
	*base.Screen
	burst *career.Burst
}

// NewBurstDetailScreen creates a new burst detail screen.
//
// The screen:
// - Displays all burst fields in a formatted layout
// - Shows optional fields only when present
// - Provides multiple actions for burst management
// - Supports back navigation
//
// Parameters:
//   - burst: The career burst to display (can be nil)
func NewBurstDetailScreen(burst *career.Burst) *BurstDetailScreen {
	return &BurstDetailScreen{
		Screen: base.NewBaseScreen(),
		burst:  burst,
	}
}

// GetBurst returns the burst being displayed.
func (s *BurstDetailScreen) GetBurst() *career.Burst {
	return s.burst
}

// Update handles messages and actions.
func (s *BurstDetailScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.SetTerminalInfo(msg.Width, msg.Height)
		return nil, nil

	case tea.KeyMsg:
		return s.handleKeyMsg(msg)
	}

	return nil, nil
}

// handleKeyMsg processes keyboard input for detail screen actions.
func (s *BurstDetailScreen) handleKeyMsg(msg tea.KeyMsg) (tea.Cmd, screens.ScreenResult) {
	switch msg.Type {
	case tea.KeyEsc, tea.KeyBackspace:
		return nil, &screens.CancelResult{}
	}

	return s.handleActionKey(msg.String())
}

// handleActionKey processes action keys (v, f, e, d, c).
func (s *BurstDetailScreen) handleActionKey(key string) (tea.Cmd, screens.ScreenResult) {
	switch key {
	case "v":
		return nil, s.actionResult("view_events")

	case "f":
		return nil, s.actionResult("view_facts")

	case "e":
		return nil, s.actionResult("edit")

	case "d":
		return nil, s.actionResult("delete")

	case "c":
		return nil, s.actionResult("confirm")
	}

	return nil, nil
}

// actionResult creates a NavigateResult for the given action.
func (s *BurstDetailScreen) actionResult(action string) *screens.NavigateResult {
	return &screens.NavigateResult{
		ResultData: map[string]interface{}{
			"action": action,
			"burst":  s.burst,
		},
	}
}

// RenderContent returns just the content (burst detail card) without StandardView wrapper.
// This allows the intent to wrap it with proper breadcrumbs and themed footer.
func (s *BurstDetailScreen) RenderContent() string {
	if s.burst == nil {
		return "No burst selected."
	}

	th := theme.Default()
	dv := widgets.NewDetailView(th).
		Title(s.burst.Name).
		FieldIf("Description", s.burst.Description).
		Field("Events", fmt.Sprintf("%d", len(s.burst.EventIDs)))

	confirmedStatus := "No"
	if s.burst.Confirmed {
		confirmedStatus = "Yes"
	}
	dv.Field("Confirmed", confirmedStatus)

	dv.Field("Created", s.burst.CreatedAt.Format("2006-01-02 15:04:05")).
		Field("Updated", s.burst.UpdatedAt.Format("2006-01-02 15:04:05"))

	return dv.Render()
}

// View returns the full screen view (for standalone usage).
// Most callers should use RenderContent() instead.
func (s *BurstDetailScreen) View() string {
	return s.RenderContent()
}
