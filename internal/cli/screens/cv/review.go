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

// CVReviewState represents the internal state constant for this screen.
const CVReviewState = "review"

// ReviewScreen displays a summary review of the generated CV.
// This screen shows metadata, statistics, and section overview before
// allowing the user to view the full CV preview with scrolling.
type ReviewScreen struct {
	*base.Screen

	cv *career.CVView
}

// NewCVReviewScreen creates a new CV review screen.
//
// Expected:
//   - cvview must be valid.
//
// Returns:
//   - A fully initialized ReviewScreen ready for use.
//
// Side effects:
//   - None.
func NewCVReviewScreen(cv *career.CVView) *ReviewScreen {
	return &ReviewScreen{
		Screen: base.NewBaseScreen(),
		cv:     cv,
	}
}

// Init initializes the screen.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (s *ReviewScreen) Init() tea.Cmd {
	return nil
}

// Update handles messages.
//
// Expected:
//   - msg must be a valid tea.Msg type.
//
// Returns:
//   - tea.Cmd: command to execute.
//   - screens.ScreenResult: result indicating user action.
//
// Side effects:
//   - May return CancelResult or NavigateResult.
func (s *ReviewScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.HandleWindowSizeMsg(msg)
		return nil, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// Go back to wizard
			return nil, &screens.CancelResult{}

		case "enter", "p":
			// Navigate to full preview
			return nil, &screens.NavigateResult{
				ResultData: "preview",
			}

		case "x":
			// Export CV directly
			return nil, &screens.NavigateResult{
				ResultData: "export",
			}

		case "e":
			// Edit CV
			return nil, &screens.NavigateResult{
				ResultData: "edit",
			}
		}
	}

	return nil, nil
}

// View renders the review screen with CV metadata and section summary.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (s *ReviewScreen) View() string {
	var b strings.Builder

	b.WriteString("📋 CV Review\n")
	b.WriteString(strings.Repeat("═", 60))
	b.WriteString("\n\n")

	if s.cv == nil {
		b.WriteString("No CV data available\n")
		b.WriteString("\n")
		b.WriteString(strings.Repeat("─", 60))
		b.WriteString("\n")
		b.WriteString("esc: back")
		return b.String()
	}

	// CV Metadata
	b.WriteString(fmt.Sprintf("  Name:     %s\n", s.cv.Name))
	b.WriteString(fmt.Sprintf("  Role:     %s\n", s.cv.TargetRole))
	b.WriteString(fmt.Sprintf("  Audience: %s\n", s.cv.TargetAudience))
	b.WriteString("\n")

	// Statistics
	b.WriteString(strings.Repeat("─", 60))
	b.WriteString("\n")
	b.WriteString("📊 Statistics\n")
	b.WriteString(strings.Repeat("─", 60))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("  Source Events: %d\n", s.cv.SourceEventCount))
	b.WriteString(fmt.Sprintf("  Source Facts:  %d\n", s.cv.SourceFactCount))
	b.WriteString(fmt.Sprintf("  Sections:      %d\n", len(s.cv.Sections)))

	// Calculate total bullets
	totalBullets := 0
	for _, section := range s.cv.Sections {
		for _, group := range section.Content {
			totalBullets += len(group.Bullets)
		}
	}
	b.WriteString(fmt.Sprintf("  Total Bullets: %d\n", totalBullets))
	b.WriteString("\n")

	// Section Summary
	if len(s.cv.Sections) > 0 {
		b.WriteString(strings.Repeat("─", 60))
		b.WriteString("\n")
		b.WriteString("📑 Sections\n")
		b.WriteString(strings.Repeat("─", 60))
		b.WriteString("\n")

		for _, section := range s.cv.Sections {
			// Count bullets in this section
			sectionBullets := 0
			for _, group := range section.Content {
				sectionBullets += len(group.Bullets)
			}

			// Format bullet count with proper pluralization
			bulletText := "bullets"
			if sectionBullets == 1 {
				bulletText = "bullet"
			}

			// Section type indicator
			typeIndicator := "  •"
			if section.SectionType == "summary" {
				typeIndicator = "  ✎"
			}

			b.WriteString(fmt.Sprintf("%s %s (%d %s)\n",
				typeIndicator, section.Title, sectionBullets, bulletText))
		}
	} else {
		b.WriteString("\n  0 sections generated\n")
	}

	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", 60))
	b.WriteString("\n")
	b.WriteString("enter/p: preview full CV  x: export  e: edit  esc: back")

	return b.String()
}

// GetCV returns the CV data.
//
// Returns:
//   - A fully initialized career.CVView ready for use.
//
// Side effects:
//   - None.
func (s *ReviewScreen) GetCV() *career.CVView {
	return s.cv
}
