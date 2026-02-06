// Package cv provides CV-related screens for the GenerateCV workflow.
package cv

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/domain/career"
	cvservice "github.com/baphled/kariya/internal/service/career/cv"
)

// CVReviewState represents the internal state constant for this screen.
const CVReviewState = "review"

// ReviewScreen displays a summary review of the generated CV.
// This screen shows metadata, statistics, and section overview before
// allowing the user to view the full CV preview with scrolling.
type ReviewScreen struct {
	*base.Screen

	cv            *career.CVView
	profileConfig *config.ProfileConfig
	viewport      viewport.Model
	ready         bool
	width         int
	height        int
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
	return NewCVReviewScreenWithProfile(cv, nil)
}

// NewCVReviewScreenWithProfile creates a new CV review screen with custom profile config.
//
// Expected:
//   - cvview must be valid.
//   - profileConfig may be nil.
//
// Returns:
//   - A fully initialized ReviewScreen ready for use.
//
// Side effects:
//   - None.
func NewCVReviewScreenWithProfile(cv *career.CVView, profileConfig *config.ProfileConfig) *ReviewScreen {
	return &ReviewScreen{
		Screen:        base.NewBaseScreen(),
		cv:            cv,
		profileConfig: profileConfig,
		width:         80,
		height:        24,
		ready:         false,
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
//   - May update viewport dimensions.
//   - May return CancelResult or NavigateResult.
func (s *ReviewScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.HandleWindowSizeMsg(msg)
		s.width = msg.Width
		s.height = msg.Height
		s.ready = false
		return nil, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return nil, &screens.CancelResult{}

		case "enter", "p":
			return nil, &screens.NavigateResult{
				ResultData: "preview",
			}

		case "x":
			return nil, &screens.NavigateResult{
				ResultData: "export",
			}

		case "e":
			return nil, &screens.NavigateResult{
				ResultData: "edit",
			}

		case "g":
			s.viewport.GotoTop()
			return nil, nil

		case "G":
			s.viewport.GotoBottom()
			return nil, nil

		case "up", "k", "down", "j", "pgup", "pgdown", "ctrl+u", "ctrl+d":
			if s.ready {
				s.viewport, cmd = s.viewport.Update(msg)
				return cmd, nil
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

	footerStyle := lipgloss.NewStyle().
		Foreground(theme.SecondaryColor())

	b.WriteString(titleStyle.Render("📋 CV Review"))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("═", minInt(60, s.width-4)))
	b.WriteString("\n\n")

	if s.cv == nil {
		b.WriteString("No CV data available\n")
		b.WriteString("\n")
		b.WriteString(strings.Repeat("─", 60))
		b.WriteString("\n")
		b.WriteString(footerStyle.Render("esc: back"))
		return b.String()
	}

	content := s.renderContent()

	if !s.ready {
		viewportHeight := s.height - 7
		if viewportHeight < 5 {
			viewportHeight = 5
		}
		viewportWidth := s.width - 4
		if viewportWidth < 40 {
			viewportWidth = 40
		}

		s.viewport = viewport.New(viewportWidth, viewportHeight)
		s.viewport.SetContent(content)
		s.ready = true
	}

	b.WriteString(s.viewport.View())
	b.WriteString("\n")
	b.WriteString(s.renderFooter())

	return b.String()
}

// renderContent renders the scrollable content.
func (s *ReviewScreen) renderContent() string {
	var b strings.Builder

	b.WriteString(s.renderPersonalDetails())
	b.WriteString(s.renderCVDetails())
	b.WriteString(s.renderStatistics())
	b.WriteString(s.renderSectionsList())

	return b.String()
}

func (s *ReviewScreen) renderPersonalDetails() string {
	theme := s.getTheme()
	sectionTitleStyle := lipgloss.NewStyle().Bold(true).Foreground(theme.AccentColor())
	labelStyle := lipgloss.NewStyle().Foreground(theme.SecondaryColor())
	valueStyle := lipgloss.NewStyle().Foreground(theme.ForegroundColor())

	profile := cvservice.NarrativeProfileFromConfig(s.profileConfig)

	var b strings.Builder
	b.WriteString(sectionTitleStyle.Render("👤 Personal Details"))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", 60))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("  %s %s\n", labelStyle.Render("Name:"), valueStyle.Render(profile.Name)))
	b.WriteString(fmt.Sprintf("  %s %s\n", labelStyle.Render("Email:"), valueStyle.Render(profile.Email)))
	b.WriteString(fmt.Sprintf("  %s %s\n", labelStyle.Render("Location:"), valueStyle.Render(profile.Location)))
	b.WriteString("\n")

	return b.String()
}

func (s *ReviewScreen) renderCVDetails() string {
	theme := s.getTheme()
	sectionTitleStyle := lipgloss.NewStyle().Bold(true).Foreground(theme.AccentColor())
	labelStyle := lipgloss.NewStyle().Foreground(theme.SecondaryColor())
	valueStyle := lipgloss.NewStyle().Foreground(theme.ForegroundColor())

	var b strings.Builder
	b.WriteString(sectionTitleStyle.Render("📄 CV Details"))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", 60))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("  %s %s\n", labelStyle.Render("CV Name:"), valueStyle.Render(s.cv.Name)))
	b.WriteString(fmt.Sprintf("  %s %s\n", labelStyle.Render("Role:"), valueStyle.Render(s.cv.TargetRole)))
	b.WriteString(fmt.Sprintf("  %s %s\n", labelStyle.Render("Audience:"), valueStyle.Render(s.cv.TargetAudience)))
	b.WriteString("\n")

	return b.String()
}

func (s *ReviewScreen) renderStatistics() string {
	theme := s.getTheme()
	sectionTitleStyle := lipgloss.NewStyle().Bold(true).Foreground(theme.AccentColor())
	labelStyle := lipgloss.NewStyle().Foreground(theme.SecondaryColor())
	valueStyle := lipgloss.NewStyle().Foreground(theme.ForegroundColor())

	var b strings.Builder
	b.WriteString(sectionTitleStyle.Render("📊 Statistics"))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", 60))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("  %s %s\n", labelStyle.Render("Source Events:"), valueStyle.Render(strconv.Itoa(s.cv.SourceEventCount))))
	b.WriteString(fmt.Sprintf("  %s %s\n", labelStyle.Render("Source Facts:"), valueStyle.Render(strconv.Itoa(s.cv.SourceFactCount))))
	b.WriteString(fmt.Sprintf("  %s %s\n", labelStyle.Render("Sections:"), valueStyle.Render(strconv.Itoa(len(s.cv.Sections)))))

	totalBullets := s.countTotalBullets()
	b.WriteString(fmt.Sprintf("  %s %s\n", labelStyle.Render("Total Bullets:"), valueStyle.Render(strconv.Itoa(totalBullets))))
	b.WriteString("\n")

	return b.String()
}

func (s *ReviewScreen) countTotalBullets() int {
	total := 0
	for _, section := range s.cv.Sections {
		for _, group := range section.Content {
			total += len(group.Bullets)
		}
	}
	return total
}

func (s *ReviewScreen) renderSectionsList() string {
	theme := s.getTheme()
	sectionTitleStyle := lipgloss.NewStyle().Bold(true).Foreground(theme.AccentColor())
	labelStyle := lipgloss.NewStyle().Foreground(theme.SecondaryColor())
	valueStyle := lipgloss.NewStyle().Foreground(theme.ForegroundColor())

	if len(s.cv.Sections) == 0 {
		return "\n  0 sections generated\n"
	}

	var b strings.Builder
	b.WriteString(sectionTitleStyle.Render("📑 Sections"))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", 60))
	b.WriteString("\n")

	for _, section := range s.cv.Sections {
		sectionBullets := s.countSectionBullets(section)
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

	return b.String()
}

func (s *ReviewScreen) countSectionBullets(section *career.CVSection) int {
	count := 0
	for _, group := range section.Content {
		count += len(group.Bullets)
	}
	return count
}

// renderFooter renders the help footer with scroll indicator.
func (s *ReviewScreen) renderFooter() string {
	theme := s.getTheme()
	footerStyle := lipgloss.NewStyle().
		Foreground(theme.SecondaryColor())

	scrollInfo := ""
	if s.ready {
		scrollInfo = fmt.Sprintf(" (%d%%)", int(s.viewport.ScrollPercent()*100))
	}

	return footerStyle.Render(fmt.Sprintf("↑↓/jk: scroll%s  enter/p: preview  x: export  e: edit  esc: back", scrollInfo))
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
