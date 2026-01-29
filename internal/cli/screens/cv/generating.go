// Package cv provides CV-related screens for the GenerateCV workflow.
package cv

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
)

// CVGeneratingState represents the internal state constant for this screen.
const CVGeneratingState = "generating"

// GeneratingScreen shows progress while generating a CV.
type GeneratingScreen struct {
	*base.Screen

	profile  string
	audience string
	spinner  int
}

// NewCVGeneratingScreen creates a new CV generating screen.
func NewCVGeneratingScreen(profile, audience string) *GeneratingScreen {
	return &GeneratingScreen{
		Screen:   base.NewBaseScreen(),
		profile:  profile,
		audience: audience,
		spinner:  0,
	}
}

// Init initializes the screen.
func (s *GeneratingScreen) Init() tea.Cmd {
	return nil
}

// Update handles messages.
func (s *GeneratingScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.Screen.HandleWindowSizeMsg(msg)
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
func (s *GeneratingScreen) View() string {
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
