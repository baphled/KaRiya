package facts

import (
	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/cli/themes"
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
		// Handle special keys by type (use tea.Key* constants).
		switch msg.Type {
		case tea.KeyEsc, tea.KeyBackspace:
			// Back to previous screen.
			// Note: 'q' (quit) is handled by the intent before delegation.
			return nil, &screens.CancelResult{}
		case tea.KeyUp:
			s.tableBehavior.HandleNavigation("up")
			return nil, nil
		case tea.KeyDown:
			s.tableBehavior.HandleNavigation("down")
			return nil, nil
		case tea.KeyPgDown:
			s.tableBehavior.HandleNavigation("pgdn")
			return nil, nil
		case tea.KeyPgUp:
			s.tableBehavior.HandleNavigation("pgup")
			return nil, nil
		case tea.KeyHome:
			s.tableBehavior.HandleNavigation("home")
			return nil, nil
		case tea.KeyEnd:
			s.tableBehavior.HandleNavigation("end")
			return nil, nil
		case tea.KeyCtrlD:
			s.tableBehavior.HandleNavigation("ctrl+d")
			return nil, nil
		case tea.KeyCtrlU:
			s.tableBehavior.HandleNavigation("ctrl+u")
			return nil, nil
		}

		// Handle vim-style keys (rune-based).
		switch msg.String() {
		case "k":
			s.tableBehavior.HandleNavigation("up")
			return nil, nil
		case "j":
			s.tableBehavior.HandleNavigation("down")
			return nil, nil
		case "g":
			s.tableBehavior.HandleNavigation("home")
			return nil, nil
		case "G":
			s.tableBehavior.HandleNavigation("end")
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
	// Propagate theme to TableBehavior for consistent styling.
	if t, ok := theme.(themes.Theme); ok && t != nil {
		s.tableBehavior.SetTheme(t)
	}
}
