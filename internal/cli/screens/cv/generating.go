// Package cv provides CV-related screens for the GenerateCV workflow.
package cv

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
)

// CVGeneratingState represents the internal state constant for this screen
const CVGeneratingState = "generating"

// CVGeneratingScreen shows progress while generating a CV.
type CVGeneratingScreen struct {
	*base.BaseScreen

	profile  string
	audience string
	spinner  int
}

// NewCVGeneratingScreen creates a new CV generating screen.
func NewCVGeneratingScreen(profile, audience string) *CVGeneratingScreen {
	return &CVGeneratingScreen{
		BaseScreen: base.NewBaseScreen(),
		profile:    profile,
		audience:   audience,
		spinner:    0,
	}
}

// Init initializes the screen.
func (s *CVGeneratingScreen) Init() tea.Cmd {
	return nil
}

// Update handles messages.
func (s *CVGeneratingScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.BaseScreen.HandleWindowSizeMsg(msg)
		return nil, nil

	case tea.KeyMsg:
		// Allow escape to cancel (return to previous screen)
		if msg.Type == tea.KeyEsc {
			return nil, &screens.CancelResult{}
		}
	}

	return nil, nil
}

// View renders the screen.
func (s *CVGeneratingScreen) View() string {
	var b strings.Builder

	spinnerChars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	spinner := spinnerChars[s.spinner%len(spinnerChars)]

	b.WriteString("\n\n")
	b.WriteString(fmt.Sprintf("  %s Generating CV...\n\n", spinner))
	b.WriteString(fmt.Sprintf("  Profile: %s\n", s.profile))
	b.WriteString(fmt.Sprintf("  Audience: %s\n\n", s.audience))
	b.WriteString("  This may take a few moments...\n\n")
	b.WriteString("  esc: Cancel")

	return b.String()
}
