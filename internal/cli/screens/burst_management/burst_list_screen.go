package burst_management

import (
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

// BurstListState identifies the burst list screen in the state matrix.
// On this screen the user sees a paginated table of career bursts with
// columns for name, description preview, confirmation status, event count,
// and creation date. Arrow keys or j/k navigate the highlighted row.
// Pressing Enter opens the selected burst's detail view. Action keys
// allow adding a burst (a), editing the selected burst (e), deleting
// it (d), or triggering burst suggestions (s). Escape returns to the
// parent menu.
const BurstListState = "burst_list"

// burstRowFormatter formats a burst for table display.
func burstRowFormatter(burst *career.Burst, _ int) []string {
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

	eventCount := fmt.Sprintf("%d", len(burst.EventIDs))

	createdStr := burst.CreatedAt.Format("2006-01-02")

	return []string{nameStr, descStr, confirmedStr, eventCount, createdStr}
}

// BurstListScreen displays a list of career bursts.
//
// This screen provides:
// - List navigation with ↑/↓ or j/k (vim-style)
// - Burst selection with Enter key
// - Actions: Add (a), Edit (e), Delete (d), Suggest (s)
// - Cancel with Escape
// - Burst count display
// - Empty state handling
//
// Usage:
//
//	screen := burst_management.NewBurstListScreen(bursts)
//	cmd, result := screen.Update(msg)
//	if result != nil && result.Type() == screens.ResultNavigate {
//	    navResult := result.(*screens.NavigateResult)
//	    if action, ok := navResult.ResultData.(map[string]interface{}); ok {
//	        // Handle action (add, edit, delete, suggest)
//	    } else if burst, ok := navResult.ResultData.(*career.Burst); ok {
//	        // User selected a burst to view details
//	    }
//	}
//
// Related:
// - Screen provides the foundation
// - docs/TUI_DEVELOPER_GUIDE.md (Screen patterns)
// - docs/TUI_STANDARDS.md (Keyboard shortcuts).
type BurstListScreen struct {
	*base.Screen
	bursts        []*career.Burst
	tableBehavior *behaviors.TableBehavior[*career.Burst]
}

// NewBurstListScreen creates a new burst list screen.
//
// The screen:
// - Displays bursts in a scrollable list
// - Supports keyboard navigation (↑/↓, j/k)
// - Shows burst name, description, confirmed status, event count, created date
// - Provides actions (add, edit, delete, suggest, view)
//
// Parameters:
//   - bursts: List of career bursts to display (can be empty)
func NewBurstListScreen(bursts []*career.Burst) *BurstListScreen {
	columns := []behaviors.ColumnDef{
		{Title: "Name", Width: 30},
		{Title: "Description", Width: 35},
		{Title: "Confirmed", Width: 10},
		{Title: "Events", Width: 8},
		{Title: "Created", Width: 12},
	}

	tableBehavior := behaviors.NewTableBehavior[*career.Burst](nil, columns, burstRowFormatter).
		PageSize(15).
		PaginationPrefix("Bursts").
		EmptyMessage("No bursts found. Press 'a' to add or 's' to suggest.")

	tableBehavior.SetItems(bursts)

	screen := &BurstListScreen{
		Screen:        base.NewBaseScreen(),
		bursts:        bursts,
		tableBehavior: tableBehavior,
	}

	return screen
}

// Update handles messages and navigation.
func (s *BurstListScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.SetTerminalInfo(msg.Width, msg.Height)
		return nil, nil

	case tea.KeyMsg:
		return s.handleKeyMsg(msg)
	}

	return nil, nil
}

// handleKeyMsg processes all keyboard input for the list screen.
func (s *BurstListScreen) handleKeyMsg(msg tea.KeyMsg) (tea.Cmd, screens.ScreenResult) {
	if result := s.handleSpecialKey(msg); result != nil {
		return nil, result
	}

	if s.handleNavigationKey(msg) {
		return nil, nil
	}

	return s.handleActionKey(msg)
}

// handleSpecialKey handles escape and enter keys.
func (s *BurstListScreen) handleSpecialKey(msg tea.KeyMsg) screens.ScreenResult {
	switch msg.Type {
	case tea.KeyEsc:
		return &screens.CancelResult{}
	case tea.KeyEnter:
		if selected := s.tableBehavior.GetSelectedItem(); selected != nil {
			return &screens.NavigateResult{ResultData: *selected}
		}
	}
	return nil
}

// handleNavigationKey handles arrow keys and vim-style navigation.
// Returns true if a navigation key was handled.
func (s *BurstListScreen) handleNavigationKey(msg tea.KeyMsg) bool {
	switch msg.Type {
	case tea.KeyUp:
		s.tableBehavior.HandleNavigation("up")
		return true
	case tea.KeyDown:
		s.tableBehavior.HandleNavigation("down")
		return true
	case tea.KeyPgDown:
		s.tableBehavior.HandleNavigation("pgdn")
		return true
	case tea.KeyPgUp:
		s.tableBehavior.HandleNavigation("pgup")
		return true
	case tea.KeyHome:
		s.tableBehavior.HandleNavigation("home")
		return true
	case tea.KeyEnd:
		s.tableBehavior.HandleNavigation("end")
		return true
	case tea.KeyCtrlD:
		s.tableBehavior.HandleNavigation("ctrl+d")
		return true
	case tea.KeyCtrlU:
		s.tableBehavior.HandleNavigation("ctrl+u")
		return true
	}

	switch msg.String() {
	case "k":
		s.tableBehavior.HandleNavigation("up")
		return true
	case "j":
		s.tableBehavior.HandleNavigation("down")
		return true
	case "g":
		s.tableBehavior.HandleNavigation("home")
		return true
	case "G":
		s.tableBehavior.HandleNavigation("end")
		return true
	}

	return false
}

// handleActionKey handles action keys (a, e, d, s).
func (s *BurstListScreen) handleActionKey(msg tea.KeyMsg) (tea.Cmd, screens.ScreenResult) {
	switch msg.String() {
	case "a":
		return nil, &screens.NavigateResult{
			ResultData: map[string]interface{}{"action": "add"},
		}

	case "e":
		if selected := s.tableBehavior.GetSelectedItem(); selected != nil {
			return nil, &screens.NavigateResult{
				ResultData: map[string]interface{}{"action": "edit", "burst": *selected},
			}
		}

	case "d":
		if selected := s.tableBehavior.GetSelectedItem(); selected != nil {
			return nil, &screens.NavigateResult{
				ResultData: map[string]interface{}{"action": "delete", "burst": *selected},
			}
		}

	case "s":
		return nil, &screens.NavigateResult{
			ResultData: map[string]interface{}{"action": "suggest"},
		}
	}

	return nil, nil
}

// RenderContent returns just the content (table) without StandardView wrapper.
// This allows the intent to wrap it with proper breadcrumbs and themed footer.
func (s *BurstListScreen) RenderContent() string {
	return s.tableBehavior.Render()
}

// View renders the burst list screen using StandardView with table.
func (s *BurstListScreen) View() string {
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
		primitives.PageBadge(th),
		primitives.ViewBadge(th),
		primitives.AddBadge(th),
		primitives.EditBadge(th),
		primitives.DeleteBadge(th),
		primitives.SuggestBadge(th),
		primitives.BackBadge(th),
	)

	breadcrumbs := []string{"Main Menu", "Burst Management"}
	return s.CreateView(breadcrumbs, content, footer)
}

// SetTheme applies theme to the table (override Screen).
func (s *BurstListScreen) SetTheme(theme interface{}) {
	s.Screen.SetTheme(theme)
	if t, ok := theme.(themes.Theme); ok && t != nil {
		s.tableBehavior.SetTheme(t)
	}
}

// GetBursts returns the list of bursts.
func (s *BurstListScreen) GetBursts() []*career.Burst {
	return s.bursts
}

// GetSelectedIndex returns the currently selected index.
func (s *BurstListScreen) GetSelectedIndex() int {
	return s.tableBehavior.GetSelectedIndex()
}
