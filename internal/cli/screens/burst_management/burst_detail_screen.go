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

// State constant for state matrix tracking (REQUIRED)
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
// - BaseScreen provides the foundation
// - docs/TUI_DEVELOPER_GUIDE.md (Screen patterns)
// - docs/TUI_STANDARDS.md (Keyboard shortcuts)
type BurstDetailScreen struct {
	*base.BaseScreen
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
		BaseScreen: base.NewBaseScreen(),
		burst:      burst,
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
	// Handle special keys by type.
	switch msg.Type {
	case tea.KeyEsc, tea.KeyBackspace:
		// Back to burst list.
		// Note: 'q' (quit) is handled by the intent before delegation.
		return nil, &screens.CancelResult{}
	}

	// Handle rune-based action keys.
	return s.handleActionKey(msg.String())
}

// handleActionKey processes action keys (v, f, e, d, c).
func (s *BurstDetailScreen) handleActionKey(key string) (tea.Cmd, screens.ScreenResult) {
	switch key {
	case "v":
		// View events in burst.
		return nil, s.actionResult("view_events")

	case "f":
		// View facts extracted from burst.
		return nil, s.actionResult("view_facts")

	case "e":
		// Edit burst.
		return nil, s.actionResult("edit")

	case "d":
		// Delete burst.
		return nil, s.actionResult("delete")

	case "c":
		// Confirm burst (mark as confirmed).
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

	// Use UIKit DetailView for consistent styling.
	th := theme.Default()
	dv := widgets.NewDetailView(th).
		Title(s.burst.Name).
		FieldIf("Description", s.burst.Description).
		Field("Events", fmt.Sprintf("%d", len(s.burst.EventIDs)))

	// Confirmation status.
	confirmedStatus := "No"
	if s.burst.Confirmed {
		confirmedStatus = "Yes"
	}
	dv.Field("Confirmed", confirmedStatus)

	// Created and updated dates.
	dv.Field("Created", s.burst.CreatedAt.Format("2006-01-02 15:04:05")).
		Field("Updated", s.burst.UpdatedAt.Format("2006-01-02 15:04:05"))

	return dv.Render()
}

// View returns the full screen view (for standalone usage).
// Most callers should use RenderContent() instead.
func (s *BurstDetailScreen) View() string {
	return s.RenderContent()
}
