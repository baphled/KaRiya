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

// ProfileSelectScreen allows users to select a CV profile.
type ProfileSelectScreen struct {
	*base.SelectScreen[*types.CVProfile]
}

// NewCVProfileSelectScreen creates a new CV profile selection screen.
//
// Expected:
//   - cvprofile must be valid.
//
// Returns:
//   - A fully initialized ProfileSelectScreen ready for use.
//
// Side effects:
//   - None.
func NewCVProfileSelectScreen(profiles []*types.CVProfile) *ProfileSelectScreen {
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

	return &ProfileSelectScreen{
		SelectScreen: baseScreen,
	}
}

// Init initializes the screen.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (s *ProfileSelectScreen) Init() tea.Cmd {
	return nil
}

// Update handles messages.
//
// Expected:
//   - msg must be a valid tea.Msg type.
//
// Returns:
//   - tea.Cmd: command to execute.
//   - screens.ScreenResult: result indicating selection.
//
// Side effects:
//   - Delegates to underlying SelectScreen.
func (s *ProfileSelectScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	return s.SelectScreen.Update(msg)
}

// View renders the screen.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (s *ProfileSelectScreen) View() string {
	return s.SelectScreen.View()
}
