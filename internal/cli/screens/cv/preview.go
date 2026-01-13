// Package cv provides CV-related screens for the GenerateCV workflow.
package cv

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/domain/career"
)

// CVPreviewState represents the internal state constant for this screen
const CVPreviewState = "preview"

// CVPreviewScreen displays a preview of the generated CV.
type CVPreviewScreen struct {
	*base.BaseScreen

	cv           *career.CVView
	scrollOffset int
}

// NewCVPreviewScreen creates a new CV preview screen.
func NewCVPreviewScreen(cv *career.CVView) *CVPreviewScreen {
	return &CVPreviewScreen{
		BaseScreen:   base.NewBaseScreen(),
		cv:           cv,
		scrollOffset: 0,
	}
}

// Init initializes the screen.
func (s *CVPreviewScreen) Init() tea.Cmd {
	return nil
}

// Update handles messages.
func (s *CVPreviewScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.BaseScreen.HandleWindowSizeMsg(msg)
		return nil, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// Go back
			return nil, &screens.CancelResult{}

		case "enter", "y":
			// Confirm CV
			return nil, &screens.NavigateResult{
				ResultData: s.cv,
			}

		case "e":
			// Edit CV
			return nil, &screens.NavigateResult{
				ResultData: "edit",
			}

		case "down", "j":
			s.scrollOffset++
			return nil, nil

		case "up", "k":
			if s.scrollOffset > 0 {
				s.scrollOffset--
			}
			return nil, nil
		}
	}

	return nil, nil
}

// View renders the screen.
func (s *CVPreviewScreen) View() string {
	var b strings.Builder

	b.WriteString("CV Preview\n")
	b.WriteString(strings.Repeat("═", 50))
	b.WriteString("\n\n")

	if s.cv != nil {
		// Display CV metadata
		b.WriteString(fmt.Sprintf("Name: %s\n", s.cv.Name))
		b.WriteString(fmt.Sprintf("Role: %s\n", s.cv.TargetRole))
		b.WriteString(fmt.Sprintf("Audience: %s\n", s.cv.TargetAudience))
		b.WriteString(fmt.Sprintf("Events: %d | Facts: %d\n",
			s.cv.SourceEventCount, s.cv.SourceFactCount))
		b.WriteString("\n")

		// Display CV sections (simplified)
		if len(s.cv.Sections) > 0 {
			for _, section := range s.cv.Sections {
				b.WriteString(fmt.Sprintf("## %s\n", section.Title))
				// Content is a slice of groups, just show count for now
				b.WriteString(fmt.Sprintf("  (%d content groups)\n", len(section.Content)))
				b.WriteString("\n")
			}
		} else {
			b.WriteString("No sections generated yet\n")
		}
	} else {
		b.WriteString("No CV data available\n")
	}

	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", 50))
	b.WriteString("\n")
	b.WriteString("↑/k: scroll up  ↓/j: scroll down  enter/y: confirm  e: edit  esc: back")

	return b.String()
}
