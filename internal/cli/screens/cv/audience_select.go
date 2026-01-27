// Package cv provides CV-related screens for the GenerateCV workflow.
package cv

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
)

// CVAudienceSelectState represents the internal state constant for this screen.
const CVAudienceSelectState = "audience_select"

// AudienceOption represents a target audience option.
type AudienceOption struct {
	ID          string
	Name        string
	Description string
}

// CVAudienceSelectScreen allows users to select a target audience.
type CVAudienceSelectScreen struct {
	*base.BaseSelectScreen[*AudienceOption]
}

// NewCVAudienceSelectScreen creates a new CV audience selection screen.
func NewCVAudienceSelectScreen(audiences []*AudienceOption) *CVAudienceSelectScreen {
	// Create item renderer for audiences
	renderer := func(item *AudienceOption) string {
		return fmt.Sprintf("%s\n  %s", item.Name, item.Description)
	}

	breadcrumbs := []string{"Main Menu", "Generate CV", "Select Audience"}
	title := "Select Target Audience"

	baseScreen := base.NewBaseSelectScreen(
		audiences,
		renderer,
		breadcrumbs,
		title,
	)

	return &CVAudienceSelectScreen{
		BaseSelectScreen: baseScreen,
	}
}

// Init initializes the screen.
func (s *CVAudienceSelectScreen) Init() tea.Cmd {
	return nil
}

// Update handles messages.
func (s *CVAudienceSelectScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	return s.BaseSelectScreen.Update(msg)
}

// View renders the screen.
func (s *CVAudienceSelectScreen) View() string {
	return s.BaseSelectScreen.View()
}
