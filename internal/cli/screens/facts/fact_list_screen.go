package facts

import (
	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

// State constant for state matrix tracking (REQUIRED)
const FactListState = "fact_list"

// factRowFormatter formats a fact for table display
func factRowFormatter(fact *career.Fact, index int) []string {
	// Truncate text to 80 chars for display
	text := fact.Text
	if len(text) > 80 {
		text = text[:77] + "..."
	}

	// Created date
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
// - BaseScreen provides the foundation
// - docs/TUI_DEVELOPER_GUIDE.md (Screen patterns)
// - docs/TUI_STANDARDS.md (Keyboard shortcuts)
type FactListScreen struct {
	*base.BaseScreen
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

	// Create table columns
	columns := []behaviors.ColumnDef{
		{Title: "Fact", Width: 80},
		{Title: "Created", Width: 12},
	}

	// Create TableBehavior for facts list (theme set later via SetTheme)
	tableBehavior := behaviors.NewTableBehavior[*career.Fact](nil, columns, factRowFormatter).
		PageSize(15).
		PaginationPrefix("Facts").
		EmptyMessage("No facts found.")

	// Set items after creation
	tableBehavior.SetItems(facts)

	screen := &FactListScreen{
		BaseScreen:    base.NewBaseScreen(),
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
		switch msg.String() {
		case "esc", "backspace":
			// Back to previous screen
			// Note: 'q' (quit) is handled by the intent before delegation
			return nil, &screens.CancelResult{}

		case "up", "k", "down", "j", "ctrl+d", "ctrl+u", "pgup", "pgdown", "home", "end", "g", "G":
			// Delegate navigation to TableBehavior
			s.tableBehavior.HandleNavigation(msg.String())
			return nil, nil
		}
	}

	return nil, nil
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
	s.BaseScreen.SetTheme(theme)
	// TableBehavior theme is set via its own SetTheme method when rendering
	// The theme is passed during rendering, not stored
}
