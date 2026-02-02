package configure

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/baphled/kariya/internal/cli/configtypes"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
)

// DomainSelectScreen allows users to select a configuration domain.
type DomainSelectScreen struct {
	*base.SelectScreen[configtypes.ConfigurationDomain]
}

// NewDomainSelectScreen creates a new domain selection screen.
//
// Expected:
//   - config must be a valid configuration object.
//
// Returns:
//   - A fully initialized DomainSelectScreen ready for use.
//
// Side effects:
//   - None.
func NewDomainSelectScreen(domains []configtypes.ConfigurationDomain) *DomainSelectScreen {
	// Breadcrumbs
	breadcrumbs := []string{"Main Menu", "Configure System"}

	// Title
	title := "Select Configuration Domain"

	baseScreen := base.NewBaseSelectScreen(
		domains,
		formatDomainLabel,
		breadcrumbs,
		title,
	)

	return &DomainSelectScreen{
		SelectScreen: baseScreen,
	}
}

// formatDomainLabel converts a domain to a display label.
func formatDomainLabel(domain configtypes.ConfigurationDomain) string {
	switch domain {
	case configtypes.DomainSystem:
		return "System"
	case configtypes.DomainProfile:
		return "Profile"
	case configtypes.DomainExport:
		return "Export"
	case configtypes.DomainUI:
		return "UI"
	default:
		return string(domain)
	}
}

// Init initializes the screen.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (s *DomainSelectScreen) Init() tea.Cmd {
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
func (s *DomainSelectScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	// Delegate to base - it handles navigation and selection
	return s.SelectScreen.Update(msg)
}

// View renders the screen.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (s *DomainSelectScreen) View() string {
	return s.SelectScreen.View()
}
