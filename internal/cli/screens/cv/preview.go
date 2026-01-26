// Package cv provides CV-related screens for the GenerateCV workflow.
package cv

import (
	"fmt"
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

// CVPreviewState represents the internal state constant for this screen.
const CVPreviewState = "preview"

// CVPreviewScreen displays the full CV content in a scrollable viewport.
// This screen shows the actual CV content with bullet points and allows
// the user to scroll through the entire document.
type CVPreviewScreen struct {
	*base.BaseScreen

	cv            *career.CVView
	profileConfig *config.ProfileConfig
	viewport      viewport.Model
	ready         bool
	width         int
	height        int
}

// NewCVPreviewScreen creates a new CV preview screen with default profile.
func NewCVPreviewScreen(cv *career.CVView) *CVPreviewScreen {
	return NewCVPreviewScreenWithProfile(cv, nil)
}

// NewCVPreviewScreenWithProfile creates a new CV preview screen with custom profile config.
// If profileConfig is nil, falls back to default narrative profile.
func NewCVPreviewScreenWithProfile(cv *career.CVView, profileConfig *config.ProfileConfig) *CVPreviewScreen {
	return &CVPreviewScreen{
		BaseScreen:    base.NewBaseScreen(),
		cv:            cv,
		profileConfig: profileConfig,
		width:         80,
		height:        24,
		ready:         false,
	}
}

// Init initializes the screen.
func (s *CVPreviewScreen) Init() tea.Cmd {
	return nil
}

// Update handles messages.
func (s *CVPreviewScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.BaseScreen.HandleWindowSizeMsg(msg)
		s.width = msg.Width
		s.height = msg.Height
		s.ready = false // Force viewport recreation on resize
		return nil, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// Go back to review
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

		case "x":
			// Export CV
			return nil, &screens.NavigateResult{
				ResultData: "export",
			}

		// Vim-style navigation
		case "g":
			// Go to top
			s.viewport.GotoTop()
			return nil, nil

		case "G":
			// Go to bottom
			s.viewport.GotoBottom()
			return nil, nil

		// Scrolling keys - pass to viewport
		case "up", "k", "down", "j", "pgup", "pgdown", "ctrl+u", "ctrl+d":
			if s.ready {
				s.viewport, cmd = s.viewport.Update(msg)
				return cmd, nil
			}
		}
	}

	return nil, nil
}

// View renders the screen with full CV content in a scrollable viewport.
func (s *CVPreviewScreen) View() string {
	theme := s.getTheme()
	var b strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.AccentColor())

	b.WriteString(titleStyle.Render("📄 CV Preview"))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("═", min(60, s.width-4)))
	b.WriteString("\n\n")

	if s.cv == nil {
		b.WriteString("No CV data available\n")
		b.WriteString("\n")
		b.WriteString(s.renderFooter())
		return b.String()
	}

	// Build CV content for viewport
	content := s.renderCVContent()

	// Initialize viewport if needed
	if !s.ready {
		// Calculate viewport dimensions
		// Account for title (2 lines), separator (1 line), spacing (2 lines), footer (2 lines)
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

	// Render viewport
	b.WriteString(s.viewport.View())
	b.WriteString("\n")

	// Footer with scroll indicator
	b.WriteString(s.renderFooter())

	return b.String()
}

// renderCVContent renders the full CV content for the viewport.
func (s *CVPreviewScreen) renderCVContent() string {
	if s.cv == nil {
		return "No CV data available"
	}

	var b strings.Builder

	// Calculate content width for word wrapping (leave margin)
	contentWidth := s.width - 8
	if contentWidth < 40 {
		contentWidth = 40
	}

	// Section styles
	theme := s.getTheme()
	sectionTitleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.AccentColor())

	headerStyle := lipgloss.NewStyle().
		Bold(true)

	mutedStyle := lipgloss.NewStyle().
		Foreground(theme.SecondaryColor())

	bulletStyle := lipgloss.NewStyle().
		Foreground(theme.ForegroundColor())

	// Render personal details header
	b.WriteString(s.renderPersonalDetails(contentWidth))
	b.WriteString("\n")

	if len(s.cv.Sections) == 0 {
		b.WriteString("No sections generated yet\n")
		return b.String()
	}

	for idx, section := range s.cv.Sections {
		if idx > 0 {
			b.WriteString("\n")
		}

		// Section title
		b.WriteString(sectionTitleStyle.Render("## " + section.Title))
		b.WriteString("\n")
		b.WriteString(strings.Repeat("─", len(section.Title)+3))
		b.WriteString("\n")

		// Handle summary section (prose) with word wrap
		if section.SectionType == "summary" && section.Summary != "" {
			wrapped := wordWrap(section.Summary, contentWidth)
			b.WriteString(wrapped)
			b.WriteString("\n")
			continue
		}

		// Handle content groups (experience, projects, skills)
		for _, group := range section.Content {
			// Group header with dates
			if group.Header != "" {
				headerText := group.Header
				if group.StartDate != "" && group.EndDate != "" {
					if group.StartDate == group.EndDate {
						headerText = fmt.Sprintf("%s (%s)", group.Header, group.StartDate)
					} else {
						headerText = fmt.Sprintf("%s (%s - %s)", group.Header, group.StartDate, group.EndDate)
					}
				}
				b.WriteString(headerStyle.Render("  " + headerText))
				b.WriteString("\n")
			}

			// Bullets with word wrapping
			bulletWidth := contentWidth - 6 // Account for "    • " prefix
			for _, bullet := range group.Bullets {
				wrapped := wordWrap(bullet.Text, bulletWidth)
				lines := strings.Split(wrapped, "\n")
				for i, line := range lines {
					if i == 0 {
						b.WriteString(bulletStyle.Render("    • " + line))
					} else {
						b.WriteString(bulletStyle.Render("      " + line)) // Indent continuation
					}
					b.WriteString("\n")
				}
			}

			// Group separator
			if len(group.Bullets) > 0 {
				b.WriteString("\n")
			}
		}
	}

	// Add scroll position indicator at bottom if content is long
	contentLines := strings.Count(b.String(), "\n")
	if contentLines > s.height-7 {
		b.WriteString("\n")
		b.WriteString(mutedStyle.Render(fmt.Sprintf("(%d lines)", contentLines)))
	}

	return b.String()
}

// renderPersonalDetails renders the personal details header.
func (s *CVPreviewScreen) renderPersonalDetails(width int) string {
	// Use profile config if provided, otherwise fall back to defaults
	profile := cvservice.NarrativeProfileFromConfig(s.profileConfig)

	var b strings.Builder
	theme := s.getTheme()

	// Name (large, bold) - using warning color for warm appearance
	nameStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.WarningColor())

	roleStyle := lipgloss.NewStyle().
		Foreground(theme.ForegroundColor())

	contactStyle := lipgloss.NewStyle().
		Foreground(theme.SecondaryColor())

	// Name
	b.WriteString(nameStyle.Render(profile.Name))
	b.WriteString("\n")

	// Role
	b.WriteString(roleStyle.Render(profile.Role))
	b.WriteString("\n")

	// Contact line
	contactLine := fmt.Sprintf("%s  |  %s  |  %s",
		profile.Location,
		profile.Email,
		profile.GitHub,
	)
	b.WriteString(contactStyle.Render(contactLine))
	b.WriteString("\n")

	// Separator
	b.WriteString(strings.Repeat("═", min(width, 60)))
	b.WriteString("\n")

	return b.String()
}

// wordWrap wraps text to the specified width, breaking at word boundaries.
func wordWrap(text string, width int) string {
	if width <= 0 {
		return text
	}

	var result strings.Builder
	var currentLine strings.Builder
	currentLen := 0

	words := strings.Fields(text)
	for i, word := range words {
		wordLen := len(word)

		// If adding this word exceeds width, start a new line
		if currentLen > 0 && currentLen+1+wordLen > width {
			result.WriteString(currentLine.String())
			result.WriteString("\n")
			currentLine.Reset()
			currentLen = 0
		}

		// Add space before word (except at start of line)
		if currentLen > 0 {
			currentLine.WriteString(" ")
			currentLen++
		}

		currentLine.WriteString(word)
		currentLen += wordLen

		// Handle very long words that exceed width
		if wordLen > width && i < len(words)-1 {
			result.WriteString(currentLine.String())
			result.WriteString("\n")
			currentLine.Reset()
			currentLen = 0
		}
	}

	// Write remaining content
	if currentLine.Len() > 0 {
		result.WriteString(currentLine.String())
	}

	return result.String()
}

// renderFooter renders the footer with help text and scroll indicator.
func (s *CVPreviewScreen) renderFooter() string {
	theme := s.getTheme()
	footerStyle := lipgloss.NewStyle().
		Foreground(theme.SecondaryColor())

	var footer strings.Builder
	footer.WriteString(strings.Repeat("─", min(60, s.width-4)))
	footer.WriteString("\n")

	// Show scroll percentage if viewport is ready and has scrollable content
	if s.ready && s.viewport.TotalLineCount() > s.viewport.Height {
		pct := int(s.viewport.ScrollPercent() * 100)
		footer.WriteString(footerStyle.Render(
			fmt.Sprintf("↑↓/jk: scroll  g/G: top/bottom  [%d%%]  enter/y: confirm  x: export  esc: back", pct)))
	} else {
		footer.WriteString(footerStyle.Render("↑↓/jk: scroll  enter/y: confirm  x: export  e: edit  esc: back"))
	}

	return footer.String()
}

// GetCV returns the CV data.
func (s *CVPreviewScreen) GetCV() *career.CVView {
	return s.cv
}

// getTheme returns the theme from BaseScreen or a default theme.
func (s *CVPreviewScreen) getTheme() themes.Theme {
	if t := s.BaseScreen.Theme(); t != nil {
		if theme, ok := t.(themes.Theme); ok {
			return theme
		}
	}
	return themes.NewDefaultTheme()
}

// min returns the minimum of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
