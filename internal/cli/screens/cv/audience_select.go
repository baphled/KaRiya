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

// AudienceSelectScreen allows users to select a target audience.
type AudienceSelectScreen struct {
	*base.SelectScreen[*AudienceOption]
}

// NewCVAudienceSelectScreen creates a new CV audience selection screen.
func NewCVAudienceSelectScreen(audiences []*AudienceOption) *AudienceSelectScreen {
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

	return &AudienceSelectScreen{
		SelectScreen: baseScreen,
	}
}

// Init initializes the screen.
func (s *AudienceSelectScreen) Init() tea.Cmd {
	return nil
}

// Update handles messages.
func (s *AudienceSelectScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	return s.SelectScreen.Update(msg)
}

// View renders the screen.
func (s *AudienceSelectScreen) View() string {
	return s.SelectScreen.View()
}
