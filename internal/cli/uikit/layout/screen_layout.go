package layout

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/baphled/kariya/internal/cli/terminal"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/navigation"
)

// ModalRenderer is an interface for components that can render a modal overlay.
// This allows ScreenLayout to work with any modal implementation without
// creating import cycles. Both feedback.Modal and components.ModalContent satisfy this.
type ModalRenderer interface {
	Render(terminalWidth, terminalHeight int) string
}

// LogoRenderer is an interface for components that can render a logo.
// This allows ScreenLayout to work with any logo implementation without
// creating import cycles.
type LogoRenderer interface {
	ViewStatic() string
	SetWidth(width int)
}

// ScreenLayout provides a standardized screen layout with logo, content, and help footer.
//
// This is the UIKit version that uses theme-based styling exclusively.
//
// Example usage:
//
//	logo := display.NewLogo(false, termInfo.Width)
//	view := NewScreenLayout(termInfo).
//	    WithLogo(logo, 2).
//	    WithBreadcrumbs("Main Menu", "Settings").
//	    WithContent("Your content here").
//	    WithHelp("↑/k Up  ↓/j Down  Enter Select  Esc Back").
//	    WithFooterSeparator(true)
//	output := view.Render()
type ScreenLayout struct {
	ShowLogo            bool
	Logo                LogoRenderer
	LogoSpacing         int
	ShowHeader          bool
	Breadcrumbs         []string
	Title               string
	Subtitle            string
	Content             string
	ContentStyle        lipgloss.Style
	ShowModal           bool
	Modal               ModalRenderer
	HelpText            string
	ShowFooter          bool
	ShowFooterSeparator bool
	TerminalInfo        *terminal.Info
	UseFullWidth        bool
	theme               themes.Theme
}

// NewScreenLayout creates a new ScreenLayout with default settings.
// If info is nil, uses sensible defaults (140x40).
func NewScreenLayout(info *terminal.Info) *ScreenLayout {
	// Handle nil terminal info gracefully with defaults
	if info == nil {
		info = &terminal.Info{Width: 140, Height: 40}
	}

	return &ScreenLayout{
		ShowLogo:            false, // Logo must be explicitly set
		LogoSpacing:         2,
		ShowHeader:          false,
		ShowFooter:          true,
		ShowFooterSeparator: false,
		TerminalInfo:        info,
		UseFullWidth:        true,
		ContentStyle:        lipgloss.NewStyle(),
		theme:               themes.NewDefaultTheme(),
	}
}

// getTheme returns the theme or default if nil.
func (sl *ScreenLayout) getTheme() themes.Theme {
	if sl.theme != nil {
		return sl.theme
	}
	return themes.NewDefaultTheme()
}

// WithLogo sets the logo to display at the top with optional spacing before it.
func (sl *ScreenLayout) WithLogo(logo LogoRenderer, spacing int) *ScreenLayout {
	sl.ShowLogo = true
	sl.Logo = logo
	sl.LogoSpacing = spacing
	return sl
}

// WithBreadcrumbs sets breadcrumbs for the header.
func (sl *ScreenLayout) WithBreadcrumbs(crumbs ...string) *ScreenLayout {
	sl.ShowHeader = true
	sl.Breadcrumbs = crumbs
	return sl
}

// WithTitle sets the title and subtitle for the header.
func (sl *ScreenLayout) WithTitle(title, subtitle string) *ScreenLayout {
	sl.ShowHeader = true
	sl.Title = title
	sl.Subtitle = subtitle
	return sl
}

// WithContent sets the main content to display.
func (sl *ScreenLayout) WithContent(content string) *ScreenLayout {
	sl.Content = content
	return sl
}

// WithContentStyle sets a custom style for the content
func (sl *ScreenLayout) WithContentStyle(style lipgloss.Style) *ScreenLayout {
	sl.ContentStyle = style
	return sl
}

// WithHelp sets the help text to display in the footer
func (sl *ScreenLayout) WithHelp(helpText string) *ScreenLayout {
	sl.ShowFooter = true
	sl.HelpText = helpText
	return sl
}

// WithFooterSeparator enables/disables the footer separator line
func (sl *ScreenLayout) WithFooterSeparator(show bool) *ScreenLayout {
	sl.ShowFooterSeparator = show
	return sl
}

// ShowModalOverlay displays a modal overlay on top of the content.
// Accepts any ModalRenderer implementation (feedback.Modal, components.ModalContent, etc.)
func (sl *ScreenLayout) ShowModalOverlay(modal ModalRenderer) *ScreenLayout {
	sl.ShowModal = true
	sl.Modal = modal
	return sl
}

// SetUseFullWidth sets whether to use full terminal width for content
func (sl *ScreenLayout) SetUseFullWidth(full bool) *ScreenLayout {
	sl.UseFullWidth = full
	return sl
}

// WithTheme sets the theme for the view
func (sl *ScreenLayout) WithTheme(theme themes.Theme) *ScreenLayout {
	sl.theme = theme
	return sl
}

// Render renders the complete view with all components
func (sl *ScreenLayout) Render() string {
	theme := sl.getTheme()
	var parts []string

	// Add logo spacing (blank lines before logo)
	if sl.ShowLogo && sl.Logo != nil {
		for i := 0; i < sl.LogoSpacing; i++ {
			parts = append(parts, "")
		}

		// Render logo (logo will be centered by JoinVertical)
		sl.Logo.SetWidth(sl.TerminalInfo.Width)
		logoOutput := sl.Logo.ViewStatic()
		parts = append(parts, logoOutput)

		// Add one blank line after logo
		parts = append(parts, "")
	}

	// Add header (breadcrumbs or title/subtitle)
	if sl.ShowHeader {
		if len(sl.Breadcrumbs) > 0 {
			// Convert string breadcrumbs to Breadcrumb structs with icons
			crumbs := make([]navigation.Breadcrumb, len(sl.Breadcrumbs))
			for i, label := range sl.Breadcrumbs {
				intent := strings.ToLower(strings.ReplaceAll(label, " ", "_"))
				crumbs[i] = navigation.Breadcrumb{
					Label:  label,
					Icon:   navigation.GetIconForIntent(intent),
					Intent: intent,
				}
			}
			bar := navigation.NewBreadcrumbBar(sl.TerminalInfo.Width, false).
				WithTheme(theme)
			bar.SetCrumbs(crumbs)
			breadcrumbOutput := bar.View()
			parts = append(parts, breadcrumbOutput)
			parts = append(parts, "") // Blank line after breadcrumbs
		}

		if sl.Title != "" {
			titleStyle := lipgloss.NewStyle().
				Foreground(theme.ForegroundColor()).
				Bold(true)
			styledTitle := titleStyle.Render(sl.Title)
			parts = append(parts, styledTitle)

			if sl.Subtitle != "" {
				subtitleStyle := lipgloss.NewStyle().
					Foreground(theme.MutedColor())
				styledSubtitle := subtitleStyle.Render(sl.Subtitle)
				parts = append(parts, styledSubtitle)
			}

			parts = append(parts, "") // Blank line after title
		}
	}

	// Add content
	if sl.Content != "" {
		contentToRender := sl.Content
		// Check if custom style has been applied
		hasStyle := sl.ContentStyle.GetBackground() != lipgloss.NoColor{} || sl.ContentStyle.GetForeground() != lipgloss.NoColor{}
		if hasStyle {
			contentToRender = sl.ContentStyle.Render(sl.Content)
		}
		parts = append(parts, contentToRender)
	}

	// Add footer
	if sl.ShowFooter && sl.HelpText != "" {
		// Add visual separator if enabled
		if sl.ShowFooterSeparator {
			// Separator will match widest content line via JoinVertical
			separator := strings.Repeat("─", 100)
			separatorStyle := lipgloss.NewStyle().
				Foreground(theme.BorderColor())
			parts = append(parts, "", separatorStyle.Render(separator))
		} else {
			parts = append(parts, "") // Just blank line
		}

		// Render help text
		helpStyle := lipgloss.NewStyle().
			Foreground(theme.MutedColor())
		styledHelp := helpStyle.Render(sl.HelpText)
		parts = append(parts, styledHelp)
	}

	// Join all parts with center alignment - aligns to widest line
	combined := lipgloss.JoinVertical(lipgloss.Center, parts...)

	// Center within terminal (both horizontal and vertical)
	rendered := lipgloss.Place(sl.TerminalInfo.Width, sl.TerminalInfo.Height,
		lipgloss.Center, lipgloss.Center, combined)

	// Add modal overlay if needed
	if sl.ShowModal && sl.Modal != nil {
		// Dim the background content
		dimmedContent := sl.dimContent(rendered)
		modalOutput := sl.Modal.Render(sl.TerminalInfo.Width, sl.TerminalInfo.Height)
		// Overlay modal on top of dimmed content
		rendered = sl.overlayModal(dimmedContent, modalOutput)
	}

	return rendered
}

// dimContent applies a dimming effect to the content.
// Uses Faint(true) for consistency with feedback.DimContent().
func (sl *ScreenLayout) dimContent(content string) string {
	dimStyle := lipgloss.NewStyle().
		Faint(true)
	return dimStyle.Render(content)
}

// overlayModal overlays modal content on top of background content
func (sl *ScreenLayout) overlayModal(background, modal string) string {
	bgLines := strings.Split(background, "\n")
	modalLines := strings.Split(modal, "\n")

	// Calculate vertical position to center modal
	bgHeight := len(bgLines)
	modalHeight := len(modalLines)
	startLine := (bgHeight - modalHeight) / 2
	if startLine < 0 {
		startLine = 0
	}

	// Overlay modal lines onto background
	result := make([]string, len(bgLines))
	copy(result, bgLines)

	for i, modalLine := range modalLines {
		lineIndex := startLine + i
		if lineIndex >= 0 && lineIndex < len(result) {
			// Center modal line horizontally using lipgloss.PlaceHorizontal
			centeredModalLine := lipgloss.PlaceHorizontal(sl.TerminalInfo.Width, lipgloss.Center, modalLine)
			result[lineIndex] = centeredModalLine
		}
	}

	return strings.Join(result, "\n")
}
