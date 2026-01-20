package export

import (
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/types"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Preview displays a preview of the export content with a scrollable viewport.
type Preview struct {
	*base.BaseScreen

	content      string
	artifactType types.ExportArtifactType
	format       types.ExportFormat
	breadcrumbs  []string
	viewport     viewport.Model
	ready        bool
	width        int
	height       int
}

// NewPreview creates a new preview screen with viewport support.
func NewPreview(content string, artifactType types.ExportArtifactType, format types.ExportFormat, breadcrumbs []string) *Preview {
	return &Preview{
		BaseScreen:   base.NewBaseScreen(),
		content:      content,
		artifactType: artifactType,
		format:       format,
		breadcrumbs:  breadcrumbs,
		width:        80,
		height:       24,
		ready:        false,
	}
}

// Init initializes the screen.
func (s *Preview) Init() tea.Cmd {
	return nil
}

// Update handles messages.
func (s *Preview) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
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
			// Go back
			return nil, &screens.CancelResult{}

		case "enter":
			// Continue to confirmation
			return nil, &screens.NavigateResult{
				ResultData: "confirm",
			}

		// Vim-style navigation
		case "g":
			// Go to top
			if s.ready {
				s.viewport.GotoTop()
			}
			return nil, nil

		case "G":
			// Go to bottom
			if s.ready {
				s.viewport.GotoBottom()
			}
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

// View renders the screen.
func (s *Preview) View() string {
	theme := s.getTheme()
	var b strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.AccentColor())

	b.WriteString(titleStyle.Render(fmt.Sprintf("📄 Preview Export (%s as %s)", s.artifactType, s.format)))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("═", min(60, s.width-4)))
	b.WriteString("\n\n")

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
		s.viewport.SetContent(s.content)
		s.ready = true
	}

	// Render viewport
	b.WriteString(s.viewport.View())
	b.WriteString("\n")

	// Footer with scroll indicator
	b.WriteString(s.renderFooter())

	return b.String()
}

// RenderContent returns just the content without StandardView wrapper.
func (s *Preview) RenderContent() string {
	return s.View()
}

// renderFooter renders the footer with help text and scroll indicator.
func (s *Preview) renderFooter() string {
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
			fmt.Sprintf("↑↓/jk: scroll  g/G: top/bottom  [%d%%]  Enter: continue  Esc: back", pct)))
	} else {
		footer.WriteString(footerStyle.Render("↑↓/jk: scroll  Enter: continue  Esc: back"))
	}

	return footer.String()
}

// getTheme returns the theme from BaseScreen or a default theme.
func (s *Preview) getTheme() themes.Theme {
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
