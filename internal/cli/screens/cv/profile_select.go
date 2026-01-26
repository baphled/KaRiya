// Package cv provides CV-related screens for the GenerateCV workflow.
package cv

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/cli/types"
)

// CVProfileSelectState represents the internal state constant for this screen.
const CVProfileSelectState = "profile_select"

// CVProfileSelectScreen allows users to select a CV profile.
type CVProfileSelectScreen struct {
	*base.BaseSelectScreen[*types.CVProfile]
}

// NewCVProfileSelectScreen creates a new CV profile selection screen.
func NewCVProfileSelectScreen(profiles []*types.CVProfile) *CVProfileSelectScreen {
	// Create item renderer for CV profiles
	renderer := func(item *types.CVProfile) string {
		// Build profile display
		lines := []string{
			item.Name,
			fmt.Sprintf("  Role: %s | Audience: %s", item.TargetRole, item.TargetAudience),
		}

		if item.Description != "" {
			lines = append(lines, fmt.Sprintf("  %s", item.Description))
		}

		return strings.Join(lines, "\n")
	}

	breadcrumbs := []string{"Main Menu", "Generate CV", "Select Profile"}
	title := "Select CV Profile"

	baseScreen := base.NewBaseSelectScreen(
		profiles,
		renderer,
		breadcrumbs,
		title,
	)

	return &CVProfileSelectScreen{
		BaseSelectScreen: baseScreen,
	}
}

// Init initializes the screen.
func (s *CVProfileSelectScreen) Init() tea.Cmd {
	return nil // BaseSelectScreen doesn't need initialization
}

// Update handles messages.
func (s *CVProfileSelectScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	return s.BaseSelectScreen.Update(msg)
}

// View renders the screen.
func (s *CVProfileSelectScreen) View() string {
	return s.BaseSelectScreen.View()
}
