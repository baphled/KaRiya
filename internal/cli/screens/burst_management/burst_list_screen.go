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

// State constant for state matrix tracking (REQUIRED).
const BurstListState = "burst_list"

// burstRowFormatter formats a burst for table display.
func burstRowFormatter(burst *career.Burst, index int) []string {
	// Column 1: Name (truncate to 27 chars).
	nameStr := burst.Name
	if len(nameStr) > 27 {
		nameStr = nameStr[:27] + "..."
	}

	// Column 2: Description (truncated preview, max 32 chars).
	descStr := strings.TrimSpace(burst.Description)
	descStr = strings.ReplaceAll(descStr, "\n", " ")
	descStr = strings.ReplaceAll(descStr, "\r", " ")
	if descStr == "" {
		descStr = "-"
	} else if len(descStr) > 32 {
		descStr = descStr[:32] + "..."
	}

	// Column 3: Confirmed Status.
	confirmedStr := "✗ No"
	if burst.Confirmed {
		confirmedStr = "✓ Yes"
	}

	// Column 4: Event Count.
	eventCount := fmt.Sprintf("%d", len(burst.EventIDs))

	// Column 5: Created Date (YYYY-MM-DD).
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
// - BaseScreen provides the foundation
// - docs/TUI_DEVELOPER_GUIDE.md (Screen patterns)
// - docs/TUI_STANDARDS.md (Keyboard shortcuts)
type BurstListScreen struct {
	*base.BaseScreen
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
	// Create table columns.
	columns := []behaviors.ColumnDef{
		{Title: "Name", Width: 30},
		{Title: "Description", Width: 35},
		{Title: "Confirmed", Width: 10},
		{Title: "Events", Width: 8},
		{Title: "Created", Width: 12},
	}

	// Create TableBehavior for bursts list (theme set later via SetTheme).
	tableBehavior := behaviors.NewTableBehavior[*career.Burst](nil, columns, burstRowFormatter).
		PageSize(15).
		PaginationPrefix("Bursts").
		EmptyMessage("No bursts found. Press 'a' to add or 's' to suggest.")

	// Set items after creation.
	tableBehavior.SetItems(bursts)

	screen := &BurstListScreen{
		BaseScreen:    base.NewBaseScreen(),
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
		switch msg.String() {
		case "esc":
			// Cancel and return to main menu.
			return nil, &screens.CancelResult{}

		case "up", "k", "down", "j", "ctrl+d", "ctrl+u", "pgup", "pgdown", "home", "end", "g", "G":
			// Delegate navigation to TableBehavior.
			s.tableBehavior.HandleNavigation(msg.String())
			return nil, nil

		case "enter":
			// View burst details.
			if selected := s.tableBehavior.GetSelectedItem(); selected != nil {
				return nil, &screens.NavigateResult{
					ResultData: *selected, // Dereference **T to get *T
				}
			}
			return nil, nil

		case "a":
			// Add new burst.
			return nil, &screens.NavigateResult{
				ResultData: map[string]interface{}{
					"action": "add",
				},
			}

		case "e":
			// Edit selected burst.
			if selected := s.tableBehavior.GetSelectedItem(); selected != nil {
				return nil, &screens.NavigateResult{
					ResultData: map[string]interface{}{
						"action": "edit",
						"burst":  *selected, // Dereference **T to get *T
					},
				}
			}
			return nil, nil

		case "d":
			// Delete selected burst.
			if selected := s.tableBehavior.GetSelectedItem(); selected != nil {
				return nil, &screens.NavigateResult{
					ResultData: map[string]interface{}{
						"action": "delete",
						"burst":  *selected, // Dereference **T to get *T
					},
				}
			}
			return nil, nil

		case "s":
			// Trigger burst suggestion (AI detection).
			return nil, &screens.NavigateResult{
				ResultData: map[string]interface{}{
					"action": "suggest",
				},
			}
		}
	}

	return nil, nil
}

// RenderContent returns just the content (table) without StandardView wrapper.
// This allows the intent to wrap it with proper breadcrumbs and themed footer.
func (s *BurstListScreen) RenderContent() string {
	// TableBehavior handles empty state and pagination internally.
	return s.tableBehavior.Render()
}

// View renders the burst list screen using StandardView with table.
func (s *BurstListScreen) View() string {
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
		primitives.HelpKeyBadge("s", "Suggest", th),
		primitives.BackBadge(th),
		primitives.QuitBadge(th),
		primitives.HelpBadge(th),
	)

	// Use BaseScreen's CreateView helper for StandardView integration.
	breadcrumbs := []string{"Main Menu", "Burst Management"}
	return s.CreateView(breadcrumbs, content, footer)
}

// SetTheme applies theme to the table (override BaseScreen).
func (s *BurstListScreen) SetTheme(theme interface{}) {
	s.BaseScreen.SetTheme(theme)
	// Apply themed table styles if theme is available.
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
