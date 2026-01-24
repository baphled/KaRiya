// Package fact provides screens for fact management.
package fact

import (
	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
	domain "github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

// ListScreen displays a list of facts.
type ListScreen struct {
	*base.BaseScreen
	tableBehavior *behaviors.TableBehavior[*domain.Fact]
}

// NewListScreen creates a new facts list screen.
func NewListScreen(tableBehavior *behaviors.TableBehavior[*domain.Fact]) *ListScreen {
	return &ListScreen{
		BaseScreen:    base.NewBaseScreen(),
		tableBehavior: tableBehavior,
	}
}

// Init initializes the screen.
func (s *ListScreen) Init() tea.Cmd {
	return nil
}

// Update processes messages.
func (s *ListScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	return nil, nil
}

// View renders the screen.
func (s *ListScreen) View() string {
	if s.tableBehavior == nil {
		return "No table behavior"
	}
	return s.tableBehavior.Render()
}

// RenderContent returns the table content for embedding.
func (s *ListScreen) RenderContent() string {
	return s.tableBehavior.Render()
}
