package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/cli/terminal"
)

// StandardView provides a standardized view layout with logo, content, and help footer
//
// Example usage:
//
//	logo := NewASCIILogo(false, termInfo.Width)
//	view := NewStandardView(termInfo).
//	    WithLogo(logo, 2).
//	    WithBreadcrumbs("Main Menu", "Settings").
//	    WithContent("Your content here").
//	    WithHelp("↑/k Up  ↓/j Down  Enter Select  Esc Back").
//	    WithFooterSeparator(true)
//	output := view.Render()
type StandardView struct {
	ShowLogo            bool
	Logo                *ASCIILogo
	LogoSpacing         int
	ShowHeader          bool
	Breadcrumbs         []string
	Title               string
	Subtitle            string
	Content             string
	ContentStyle        lipgloss.Style
	ShowModal           bool
	Modal               *ModalContent
	HelpText            string
	ShowFooter          bool
	ShowFooterSeparator bool
	TerminalInfo        *terminal.Info
	UseFullWidth        bool
}

// NewStandardView creates a new StandardView with default settings
// If info is nil, uses sensible defaults (140x40)
func NewStandardView(info *terminal.Info) *StandardView {
	// Handle nil terminal info gracefully with defaults
	if info == nil {
		info = &terminal.Info{Width: 140, Height: 40}
	}

	return &StandardView{
		ShowLogo:            false, // Logo must be explicitly set
		LogoSpacing:         2,
		ShowHeader:          false,
		ShowFooter:          true,
		ShowFooterSeparator: false,
		TerminalInfo:        info,
		UseFullWidth:        true,
		ContentStyle:        lipgloss.NewStyle(),
	}
}

// WithLogo sets the logo to display at the top with optional spacing before it
func (sv *StandardView) WithLogo(logo *ASCIILogo, spacing int) *StandardView {
	sv.ShowLogo = true
	sv.Logo = logo
	sv.LogoSpacing = spacing
	return sv
}

// WithBreadcrumbs sets breadcrumbs for the header
func (sv *StandardView) WithBreadcrumbs(crumbs ...string) *StandardView {
	sv.ShowHeader = true
	sv.Breadcrumbs = crumbs
	return sv
}

// WithTitle sets the title and subtitle for the header
func (sv *StandardView) WithTitle(title, subtitle string) *StandardView {
	sv.ShowHeader = true
	sv.Title = title
	sv.Subtitle = subtitle
	return sv
}

// WithContent sets the main content to display
func (sv *StandardView) WithContent(content string) *StandardView {
	sv.Content = content
	return sv
}

// WithContentStyle sets a custom style for the content
func (sv *StandardView) WithContentStyle(style lipgloss.Style) *StandardView {
	sv.ContentStyle = style
	return sv
}

// WithHelp sets the help text to display in the footer
func (sv *StandardView) WithHelp(helpText string) *StandardView {
	sv.ShowFooter = true
	sv.HelpText = helpText
	return sv
}

// WithFooterSeparator enables/disables the footer separator line
func (sv *StandardView) WithFooterSeparator(show bool) *StandardView {
	sv.ShowFooterSeparator = show
	return sv
}

// ShowModalOverlay displays a modal overlay on top of the content
func (sv *StandardView) ShowModalOverlay(modal *ModalContent) *StandardView {
	sv.ShowModal = true
	sv.Modal = modal
	return sv
}

// SetUseFullWidth sets whether to use full terminal width for content
func (sv *StandardView) SetUseFullWidth(full bool) *StandardView {
	sv.UseFullWidth = full
	return sv
}

// Render renders the complete view with all components
func (sv *StandardView) Render() string {
	var parts []string

	// Add logo spacing (blank lines before logo)
	if sv.ShowLogo && sv.Logo != nil {
		for i := 0; i < sv.LogoSpacing; i++ {
			parts = append(parts, "")
		}

		// Render logo
		sv.Logo.SetExternalCentering(true)
		sv.Logo.SetWidth(sv.TerminalInfo.Width)
		logoOutput := sv.Logo.ViewStatic()
		parts = append(parts, logoOutput)

		// Add one blank line after logo
		parts = append(parts, "")
	}

	// Add header (breadcrumbs or title/subtitle)
	if sv.ShowHeader {
		if len(sv.Breadcrumbs) > 0 {
			// Convert string breadcrumbs to Breadcrumb structs with icons
			crumbs := make([]Breadcrumb, len(sv.Breadcrumbs))
			for i, label := range sv.Breadcrumbs {
				intent := strings.ToLower(strings.ReplaceAll(label, " ", "_"))
				crumbs[i] = Breadcrumb{
					Label:  label,
					Icon:   GetIconForIntent(intent),
					Intent: intent,
				}
			}
			bar := NewBreadcrumbBar(sv.TerminalInfo.Width, false)
			bar.SetCrumbs(crumbs)
			breadcrumbOutput := bar.View()
			centeredBreadcrumbs := styles.CenterHorizontal(breadcrumbOutput, sv.TerminalInfo.Width)
			parts = append(parts, centeredBreadcrumbs)
			parts = append(parts, "") // Blank line after breadcrumbs
		}

		if sv.Title != "" {
			titleStyle := lipgloss.NewStyle().
				Foreground(styles.ColorTextPrimary).
				Bold(true)
			styledTitle := titleStyle.Render(sv.Title)
			centeredTitle := styles.CenterHorizontal(styledTitle, sv.TerminalInfo.Width)
			parts = append(parts, centeredTitle)

			if sv.Subtitle != "" {
				subtitleStyle := lipgloss.NewStyle().
					Foreground(styles.ColorTextSecondary)
				styledSubtitle := subtitleStyle.Render(sv.Subtitle)
				centeredSubtitle := styles.CenterHorizontal(styledSubtitle, sv.TerminalInfo.Width)
				parts = append(parts, centeredSubtitle)
			}

			parts = append(parts, "") // Blank line after title
		}
	}

	// Add content
	if sv.Content != "" {
		contentToRender := sv.Content
		// Check if custom style has been applied
		hasStyle := sv.ContentStyle.GetBackground() != lipgloss.NoColor{} || sv.ContentStyle.GetForeground() != lipgloss.NoColor{}
		if hasStyle {
			contentToRender = sv.ContentStyle.Render(sv.Content)
		}
		parts = append(parts, contentToRender)
	}

	// Add footer
	if sv.ShowFooter && sv.HelpText != "" {
		// Add visual separator if enabled
		if sv.ShowFooterSeparator {
			separator := strings.Repeat("─", sv.TerminalInfo.Width)
			separatorStyle := lipgloss.NewStyle().
				Foreground(styles.ColorBorder)
			parts = append(parts, "", separatorStyle.Render(separator))
		} else {
			parts = append(parts, "") // Just blank line
		}

		// Render help text
		helpStyle := lipgloss.NewStyle().
			Foreground(styles.ColorTextMuted)
		styledHelp := helpStyle.Render(sv.HelpText)
		centeredHelp := styles.CenterHorizontal(styledHelp, sv.TerminalInfo.Width)
		parts = append(parts, centeredHelp)
	}

	// Join all parts
	output := strings.Join(parts, "\n")

	// Use SmartContainer for final centering
	container := NewSmartContainer(sv.TerminalInfo).
		SetContent(output).
		SetCenteringMode(CenterBoth)

	rendered := container.Render()

	// Add modal overlay if needed
	if sv.ShowModal && sv.Modal != nil {
		// Dim the background content
		dimmedContent := sv.dimContent(rendered)
		modalOutput := sv.Modal.Render(sv.TerminalInfo.Width, sv.TerminalInfo.Height)
		// Overlay modal on top of dimmed content
		rendered = sv.overlayModal(dimmedContent, modalOutput)
	}

	return rendered
}

// dimContent applies a dimming effect to the content
func (sv *StandardView) dimContent(content string) string {
	dimStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")) // Gray color for dimming
	return dimStyle.Render(content)
}

// overlayModal overlays modal content on top of background content
func (sv *StandardView) overlayModal(background, modal string) string {
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
			// Center modal line horizontally
			centeredModalLine := styles.CenterHorizontal(modalLine, sv.TerminalInfo.Width)
			result[lineIndex] = centeredModalLine
		}
	}

	return strings.Join(result, "\n")
}
