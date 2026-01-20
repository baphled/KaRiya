package configure

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/baphled/kariya/internal/cli/configtypes"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
)

// DomainSelectScreen allows users to select a configuration domain
type DomainSelectScreen struct {
	*base.BaseSelectScreen[configtypes.ConfigurationDomain]
}

// NewDomainSelectScreen creates a new domain selection screen
func NewDomainSelectScreen(domains []configtypes.ConfigurationDomain) *DomainSelectScreen {
	// Domain renderer - formats domain for display
	renderer := func(domain configtypes.ConfigurationDomain) string {
		return formatDomainLabel(domain)
	}

	// Breadcrumbs
	breadcrumbs := []string{"Main Menu", "Configure System"}

	// Title
	title := "Select Configuration Domain"

	baseScreen := base.NewBaseSelectScreen(
		domains,
		renderer,
		breadcrumbs,
		title,
	)

	return &DomainSelectScreen{
		BaseSelectScreen: baseScreen,
	}
}

// formatDomainLabel converts a domain to a display label
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

// Init initializes the screen
func (s *DomainSelectScreen) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (s *DomainSelectScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	// Delegate to base - it handles navigation and selection
	return s.BaseSelectScreen.Update(msg)
}

// View renders the screen
func (s *DomainSelectScreen) View() string {
	return s.BaseSelectScreen.View()
}
