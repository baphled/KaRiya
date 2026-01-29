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

// GetAvailableContentHeight calculates the height available for content between header and footer.
// This is useful for screens that need to size their content (tables, viewports) to fill available space.
//
// Returns: terminalHeight - headerHeight - footerHeight
//
// Example:
//
//	contentHeight := screenLayout.GetAvailableContentHeight()
//	viewport := viewport.New(width, contentHeight)
func (sl *ScreenLayout) GetAvailableContentHeight() int {
	theme := sl.getTheme()
	var headerParts []string
	var footerParts []string

	// Calculate header height (same logic as Render)
	if sl.ShowLogo && sl.Logo != nil {
		headerParts = append(headerParts, "", "") // 2 blank lines before logo
		sl.Logo.SetWidth(sl.TerminalInfo.Width)
		logoOutput := sl.Logo.ViewStatic()
		headerParts = append(headerParts, logoOutput)
		headerParts = append(headerParts, "") // Blank line after logo
	}

	if sl.ShowHeader {
		if len(sl.Breadcrumbs) > 0 {
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
			headerParts = append(headerParts, breadcrumbOutput)
			headerParts = append(headerParts, "") // Blank line after breadcrumbs
		}

		if sl.Title != "" {
			titleStyle := lipgloss.NewStyle().
				Foreground(theme.ForegroundColor()).
				Bold(true)
			styledTitle := titleStyle.Render(sl.Title)
			headerParts = append(headerParts, styledTitle)

			if sl.Subtitle != "" {
				subtitleStyle := lipgloss.NewStyle().
					Foreground(theme.MutedColor())
				styledSubtitle := subtitleStyle.Render(sl.Subtitle)
				headerParts = append(headerParts, styledSubtitle)
			}

			headerParts = append(headerParts, "") // Blank line after title
		}
	}

	// Calculate footer height
	if sl.ShowFooter && sl.HelpText != "" {
		if sl.ShowFooterSeparator {
			separator := strings.Repeat("─", 100)
			separatorStyle := lipgloss.NewStyle().
				Foreground(theme.BorderColor())
			footerParts = append(footerParts, "", separatorStyle.Render(separator))
		} else {
			footerParts = append(footerParts, "") // Just blank line
		}

		helpStyle := lipgloss.NewStyle().
			Foreground(theme.MutedColor())
		styledHelp := helpStyle.Render(sl.HelpText)
		footerParts = append(footerParts, styledHelp)
	}

	// Get heights
	header := lipgloss.JoinVertical(lipgloss.Center, headerParts...)
	footer := lipgloss.JoinVertical(lipgloss.Center, footerParts...)

	headerHeight := lipgloss.Height(header)
	footerHeight := lipgloss.Height(footer)

	// Calculate available content height
	availableHeight := sl.TerminalInfo.Height - headerHeight - footerHeight

	// Ensure minimum height of 1
	if availableHeight < 1 {
		availableHeight = 1
	}

	return availableHeight
}

// Render renders the complete view with all components.
// Layout strategy: pin logo to top (line 0), footer to bottom, content flows after header.
func (sl *ScreenLayout) Render() string {
	theme := sl.getTheme()
	var headerParts []string
	var contentParts []string
	var footerParts []string

	// === HEADER SECTION (pinned to top) ===
	// Add 2 blank lines before logo for breathing room
	if sl.ShowLogo && sl.Logo != nil {
		headerParts = append(headerParts, "", "") // 2 blank lines before logo
		sl.Logo.SetWidth(sl.TerminalInfo.Width)
		logoOutput := sl.Logo.ViewStatic()
		headerParts = append(headerParts, logoOutput)
		// Add one blank line after logo
		headerParts = append(headerParts, "")
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
			headerParts = append(headerParts, breadcrumbOutput)
			headerParts = append(headerParts, "") // Blank line after breadcrumbs
		}

		if sl.Title != "" {
			titleStyle := lipgloss.NewStyle().
				Foreground(theme.ForegroundColor()).
				Bold(true)
			styledTitle := titleStyle.Render(sl.Title)
			headerParts = append(headerParts, styledTitle)

			if sl.Subtitle != "" {
				subtitleStyle := lipgloss.NewStyle().
					Foreground(theme.MutedColor())
				styledSubtitle := subtitleStyle.Render(sl.Subtitle)
				headerParts = append(headerParts, styledSubtitle)
			}

			headerParts = append(headerParts, "") // Blank line after title
		}
	}

	// === CONTENT SECTION (flows after header) ===
	if sl.Content != "" {
		contentToRender := sl.Content
		// Check if custom style has been applied
		hasStyle := sl.ContentStyle.GetBackground() != lipgloss.NoColor{} || sl.ContentStyle.GetForeground() != lipgloss.NoColor{}
		if hasStyle {
			contentToRender = sl.ContentStyle.Render(sl.Content)
		}
		contentParts = append(contentParts, contentToRender)
	}

	// === FOOTER SECTION (pinned to bottom) ===
	if sl.ShowFooter && sl.HelpText != "" {
		// Add visual separator if enabled
		if sl.ShowFooterSeparator {
			// Separator will match widest content line via JoinVertical
			separator := strings.Repeat("─", 100)
			separatorStyle := lipgloss.NewStyle().
				Foreground(theme.BorderColor())
			footerParts = append(footerParts, "", separatorStyle.Render(separator))
		} else {
			footerParts = append(footerParts, "") // Just blank line
		}

		// Render help text
		helpStyle := lipgloss.NewStyle().
			Foreground(theme.MutedColor())
		styledHelp := helpStyle.Render(sl.HelpText)
		footerParts = append(footerParts, styledHelp)
	}

	// === ASSEMBLE WITH SPACER ===
	// Build each section
	header := lipgloss.JoinVertical(lipgloss.Center, headerParts...)
	content := lipgloss.JoinVertical(lipgloss.Center, contentParts...)
	footer := lipgloss.JoinVertical(lipgloss.Center, footerParts...)

	// Calculate heights
	headerHeight := lipgloss.Height(header)
	contentHeight := lipgloss.Height(content)
	footerHeight := lipgloss.Height(footer)

	// Calculate spacer to fill the gap
	spacerHeight := sl.TerminalInfo.Height - headerHeight - contentHeight - footerHeight
	if spacerHeight < 0 {
		spacerHeight = 0 // Graceful: don't go negative
	}

	// Combine all sections with spacer lines added individually
	var allParts []string
	if header != "" {
		allParts = append(allParts, header)
	}
	if content != "" {
		allParts = append(allParts, content)
	}
	// Add spacer lines individually (not as a joined string)
	for i := 0; i < spacerHeight; i++ {
		allParts = append(allParts, "")
	}
	if footer != "" {
		allParts = append(allParts, footer)
	}

	combined := lipgloss.JoinVertical(lipgloss.Center, allParts...)

	// Place with Top vertical alignment (pins to top)
	rendered := lipgloss.Place(sl.TerminalInfo.Width, sl.TerminalInfo.Height,
		lipgloss.Center, lipgloss.Top, combined)

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
