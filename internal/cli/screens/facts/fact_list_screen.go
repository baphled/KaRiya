package facts

import (
	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

// FactListState identifies the fact list screen in the state matrix.
// On this screen the user sees a read-only, paginated table of career
// facts with two columns: the fact text (truncated to 80 characters) and
// its creation date. Arrow keys or j/k scroll through rows; Page Up/Down
// and Ctrl+D/U jump by page. Pressing Escape or Backspace returns to the
// parent screen. No editing or selection actions are available.
const FactListState = "fact_list"

// factRowFormatter formats a fact for table display.
func factRowFormatter(fact *career.Fact, _ int) []string {
	text := fact.Text
	if len(text) > 80 {
		text = text[:77] + "..."
	}

	createdStr := fact.CreatedAt.Format("2006-01-02")

	return []string{text, createdStr}
}

// FactListScreen displays a list of facts in a table.
//
// This screen provides:
// - List navigation with ↑/↓ or j/k (vim-style)
// - Scrollable list of facts
// - Back navigation with Escape or backspace
// - Empty state handling
//
// Usage:
//
//	screen := facts.NewFactListScreen(factList)
//	cmd, result := screen.Update(msg)
//	if result != nil && result.Type() == screens.ResultCancel {
//	    // User pressed escape/backspace to go back
//	}
//
// Related:
// - Screen provides the foundation
// - docs/TUI_DEVELOPER_GUIDE.md (Screen patterns)
// - docs/TUI_STANDARDS.md (Keyboard shortcuts).
type FactListScreen struct {
	*base.Screen
	facts         []*career.Fact
	tableBehavior *behaviors.TableBehavior[*career.Fact]
}

// NewFactListScreen creates a new fact list screen.
//
// The screen:
// - Displays facts in a scrollable list
// - Supports keyboard navigation (↑/↓, j/k)
// - Shows fact text and created date
// - Provides back navigation
//
// Parameters:
//   - facts: List of facts to display (can be empty or nil)
func NewFactListScreen(facts []*career.Fact) *FactListScreen {
	if facts == nil {
		facts = []*career.Fact{}
	}

	columns := []behaviors.ColumnDef{
		{Title: "Fact", Width: 80},
		{Title: "Created", Width: 12},
	}

	tableBehavior := behaviors.NewTableBehavior[*career.Fact](nil, columns, factRowFormatter).
		PageSize(15).
		PaginationPrefix("Facts").
		EmptyMessage("No facts found.")

	tableBehavior.SetItems(facts)

	screen := &FactListScreen{
		Screen:        base.NewBaseScreen(),
		facts:         facts,
		tableBehavior: tableBehavior,
	}

	return screen
}

// Update handles messages and navigation.
func (s *FactListScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.SetTerminalInfo(msg.Width, msg.Height)
		return nil, nil

	case tea.KeyMsg:
		return s.handleKeyMsg(msg)
	}

	return nil, nil
}

// handleKeyMsg processes keyboard input for navigation and actions.
func (s *FactListScreen) handleKeyMsg(msg tea.KeyMsg) (tea.Cmd, screens.ScreenResult) {
	if nav := s.handleKeyType(msg.Type); nav != "" {
		if nav == "cancel" {
			return nil, &screens.CancelResult{}
		}
		s.tableBehavior.HandleNavigation(nav)
		return nil, nil
	}

	if nav := s.handleVimKey(msg.String()); nav != "" {
		s.tableBehavior.HandleNavigation(nav)
		return nil, nil
	}

	return nil, nil
}

// handleKeyType maps key types to navigation commands.
// Returns "cancel" for exit keys, navigation string for nav keys, or empty for unhandled.
func (s *FactListScreen) handleKeyType(keyType tea.KeyType) string {
	switch keyType {
	case tea.KeyEsc, tea.KeyBackspace:
		return "cancel"
	case tea.KeyUp:
		return "up"
	case tea.KeyDown:
		return "down"
	case tea.KeyPgDown:
		return "pgdn"
	case tea.KeyPgUp:
		return "pgup"
	case tea.KeyHome:
		return "home"
	case tea.KeyEnd:
		return "end"
	case tea.KeyCtrlD:
		return "ctrl+d"
	case tea.KeyCtrlU:
		return "ctrl+u"
	default:
		return ""
	}
}

// handleVimKey maps vim-style key strings to navigation commands.
func (s *FactListScreen) handleVimKey(key string) string {
	switch key {
	case "k":
		return "up"
	case "j":
		return "down"
	case "g":
		return "home"
	case "G":
		return "end"
	default:
		return ""
	}
}

// RenderContent returns just the table content without StandardView wrapper.
// This allows the intent to wrap it with proper breadcrumbs and themed footer.
func (s *FactListScreen) RenderContent() string {
	return s.tableBehavior.Render()
}

// View returns the full screen view (for standalone usage).
// Most callers should use RenderContent() instead.
func (s *FactListScreen) View() string {
	return s.RenderContent()
}

// SetTheme sets the theme for the screen and its table behavior.
func (s *FactListScreen) SetTheme(theme interface{}) {
	s.Screen.SetTheme(theme)
	if t, ok := theme.(themes.Theme); ok && t != nil {
		s.tableBehavior.SetTheme(t)
	}
}
