package components

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/cli/themes"
)

// ASCIILogo represents the KaRiya logo with optional animation
type ASCIILogo struct {
	animated     bool
	fadeProgress float64 // 0.0 to 1.0
	width        int
	tagline      string
	showTagline  bool
	version      string
	showVersion  bool
	theme        themes.Theme
}

const (
	// Bold ASCII art logo for KaRiya
	logoArt = `██╗  ██╗ █████╗ ██████╗ ██╗██╗   ██╗ █████╗ 
██║ ██╔╝██╔══██╗██╔══██╗██║╚██╗ ██╔╝██╔══██╗
█████╔╝ ███████║██████╔╝██║ ╚████╔╝ ███████║
██╔═██╗ ██╔══██║██╔══██╗██║  ╚██╔╝  ██╔══██║
██║  ██╗██║  ██║██║  ██║██║   ██║   ██║  ██║
╚═╝  ╚═╝╚═╝  ╚═╝╚═╝  ╚═╝╚═╝   ╚═╝   ╚═╝  ╚═╝`

	// Animation settings
	animationDuration = 300 * time.Millisecond
	animationFrames   = 10
	frameInterval     = 30 * time.Millisecond
)

// TickMsg is sent periodically to update the animation
type TickMsg time.Time

// NewASCIILogo creates a new ASCII logo component
func NewASCIILogo(animated bool, width int) *ASCIILogo {
	return &ASCIILogo{
		animated:     animated,
		fadeProgress: 0.0,
		width:        width,
		tagline:      "Career Event Management System",
		showTagline:  true,
		version:      "v1.0.0",
		showVersion:  true,
	}
}

// SetWidth sets the width for centering calculations
func (l *ASCIILogo) SetWidth(width int) {
	l.width = width
}

// SetTagline sets the tagline text
func (l *ASCIILogo) SetTagline(tagline string) {
	l.tagline = tagline
}

// ShowTagline controls tagline visibility
func (l *ASCIILogo) ShowTagline(show bool) {
	l.showTagline = show
}

// SetVersion sets the version text
func (l *ASCIILogo) SetVersion(version string) {
	l.version = version
}

// ShowVersion controls version visibility
func (l *ASCIILogo) ShowVersion(show bool) {
	l.showVersion = show
}

// WithTheme sets the theme for the logo
func (l *ASCIILogo) WithTheme(theme themes.Theme) *ASCIILogo {
	l.theme = theme
	return l
}

// Theme helper methods for consistent themed styling.

// getAccentColor returns the accent color from theme or fallback.
func (l *ASCIILogo) getAccentColor() lipgloss.Color {
	if l.theme != nil {
		return l.theme.PrimaryColor()
	}
	return styles.ColorAccentTeal
}

// getSecondaryColor returns the secondary text color from theme or fallback.
func (l *ASCIILogo) getSecondaryColor() lipgloss.Color {
	if l.theme != nil {
		return l.theme.MutedColor()
	}
	return styles.ColorTextSecondary
}

// getMutedColor returns the muted text color from theme or fallback.
func (l *ASCIILogo) getMutedColor() lipgloss.Color {
	if l.theme != nil {
		return l.theme.MutedColor()
	}
	return styles.ColorTextMuted
}

// Init initializes the logo component
func (l *ASCIILogo) Init() tea.Cmd {
	if l.animated {
		return l.tick()
	}
	l.fadeProgress = 1.0
	return nil
}

// Update handles animation updates
func (l *ASCIILogo) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case TickMsg:
		if l.animated && l.fadeProgress < 1.0 {
			l.fadeProgress += 0.1 // 10 frames to reach 1.0
			if l.fadeProgress < 1.0 {
				return l, l.tick()
			}
			l.fadeProgress = 1.0
		}
	}
	return l, nil
}

// tick returns a command that sends a TickMsg after the frame interval
func (l *ASCIILogo) tick() tea.Cmd {
	return tea.Tick(frameInterval, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// View renders the logo
func (l *ASCIILogo) View() string {
	return l.render()
}

// ViewStatic renders the logo without animation (instant display)
func (l *ASCIILogo) ViewStatic() string {
	savedProgress := l.fadeProgress
	l.fadeProgress = 1.0
	result := l.render()
	l.fadeProgress = savedProgress
	return result
}

// render generates the logo output with current fade progress
func (l *ASCIILogo) render() string {
	var parts []string

	// Render the ASCII art with fade effect
	logoLines := strings.Split(logoArt, "\n")
	styledLines := make([]string, len(logoLines))

	for i, line := range logoLines {
		styledLine := l.applyFadeStyle(line)
		styledLines[i] = styledLine
	}

	logoRendered := strings.Join(styledLines, "\n")
	parts = append(parts, logoRendered)

	// Add tagline if enabled
	if l.showTagline && l.tagline != "" {
		taglineStyle := lipgloss.NewStyle().
			Foreground(l.getSecondaryColor()).
			Faint(l.fadeProgress < 1.0)

		taglineText := taglineStyle.Render(l.tagline)
		parts = append(parts, "", taglineText)
	}

	// Add version if enabled
	if l.showVersion && l.version != "" {
		versionStyle := lipgloss.NewStyle().
			Foreground(l.getMutedColor()).
			Italic(true).
			Faint(l.fadeProgress < 1.0)

		versionText := versionStyle.Render(l.version)
		parts = append(parts, versionText)
	}

	return strings.Join(parts, "\n")
}

// applyFadeStyle applies the fade effect based on current progress
func (l *ASCIILogo) applyFadeStyle(text string) string {
	if l.fadeProgress >= 1.0 {
		// Full opacity - use primary accent color
		return lipgloss.NewStyle().
			Foreground(l.getAccentColor()).
			Bold(true).
			Render(text)
	}

	// Fading in - adjust opacity by making it faint
	return lipgloss.NewStyle().
		Foreground(l.getAccentColor()).
		Bold(true).
		Faint(true).
		Render(text)
}

// GetHeight returns the height of the logo in lines
func (l *ASCIILogo) GetHeight() int {
	height := 6 // Logo art is 6 lines

	if l.showTagline {
		height += 2 // Empty line + tagline
	}

	if l.showVersion {
		height += 1 // Version line
	}

	return height
}

// GetWidth returns the width of the logo
func (l *ASCIILogo) GetWidth() int {
	// The logo art is 51 characters wide
	return 51
}
