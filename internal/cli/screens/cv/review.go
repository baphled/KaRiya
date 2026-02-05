// Package cv provides CV-related screens for the GenerateCV workflow.
package cv

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/cli/themes"
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
	theme := s.getTheme()
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.AccentColor())

	sectionTitleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.AccentColor())

	labelStyle := lipgloss.NewStyle().
		Foreground(theme.SecondaryColor())

	valueStyle := lipgloss.NewStyle().
		Foreground(theme.ForegroundColor())

	footerStyle := lipgloss.NewStyle().
		Foreground(theme.SecondaryColor())

	b.WriteString(titleStyle.Render("📋 CV Review"))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("═", 60))
	b.WriteString("\n\n")

	if s.cv == nil {
		b.WriteString("No CV data available\n")
		b.WriteString("\n")
		b.WriteString(strings.Repeat("─", 60))
		b.WriteString("\n")
		b.WriteString(footerStyle.Render("esc: back"))
		return b.String()
	}

	b.WriteString(fmt.Sprintf("  %s %s\n", labelStyle.Render("Name:"), valueStyle.Render(s.cv.Name)))
	b.WriteString(fmt.Sprintf("  %s %s\n", labelStyle.Render("Role:"), valueStyle.Render(s.cv.TargetRole)))
	b.WriteString(fmt.Sprintf("  %s %s\n", labelStyle.Render("Audience:"), valueStyle.Render(s.cv.TargetAudience)))
	b.WriteString("\n")

	b.WriteString(strings.Repeat("─", 60))
	b.WriteString("\n")
	b.WriteString(sectionTitleStyle.Render("📊 Statistics"))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", 60))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("  %s %s\n", labelStyle.Render("Source Events:"), valueStyle.Render(strconv.Itoa(s.cv.SourceEventCount))))
	b.WriteString(fmt.Sprintf("  %s %s\n", labelStyle.Render("Source Facts:"), valueStyle.Render(strconv.Itoa(s.cv.SourceFactCount))))
	b.WriteString(fmt.Sprintf("  %s %s\n", labelStyle.Render("Sections:"), valueStyle.Render(strconv.Itoa(len(s.cv.Sections)))))

	totalBullets := 0
	for _, section := range s.cv.Sections {
		for _, group := range section.Content {
			totalBullets += len(group.Bullets)
		}
	}
	b.WriteString(fmt.Sprintf("  %s %s\n", labelStyle.Render("Total Bullets:"), valueStyle.Render(strconv.Itoa(totalBullets))))
	b.WriteString("\n")

	if len(s.cv.Sections) > 0 {
		b.WriteString(strings.Repeat("─", 60))
		b.WriteString("\n")
		b.WriteString(sectionTitleStyle.Render("📑 Sections"))
		b.WriteString("\n")
		b.WriteString(strings.Repeat("─", 60))
		b.WriteString("\n")

		for _, section := range s.cv.Sections {
			sectionBullets := 0
			for _, group := range section.Content {
				sectionBullets += len(group.Bullets)
			}

			bulletText := "bullets"
			if sectionBullets == 1 {
				bulletText = "bullet"
			}

			typeIndicator := "  •"
			if section.SectionType == "summary" {
				typeIndicator = "  ✎"
			}

			b.WriteString(fmt.Sprintf("%s %s %s\n",
				typeIndicator,
				valueStyle.Render(section.Title),
				labelStyle.Render(fmt.Sprintf("(%d %s)", sectionBullets, bulletText))))
		}
	} else {
		b.WriteString("\n  0 sections generated\n")
	}

	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", 60))
	b.WriteString("\n")
	b.WriteString(footerStyle.Render("enter/p: preview full CV  x: export  e: edit  esc: back"))

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

// getTheme returns the theme from Screen or a default theme.
func (s *ReviewScreen) getTheme() themes.Theme {
	if t := s.Theme(); t != nil {
		if theme, ok := t.(themes.Theme); ok {
			return theme
		}
	}
	return themes.NewDefaultTheme()
}
