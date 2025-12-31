package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/baphled/kariya/internal/cli/styles"
)

// FooterModel renders a consistent screen footer with status, mode, and help text
type FooterModel struct {
	statusMessage string
	modeContext   string // Current mode/context (e.g., "Capture Mode: Timeline")
	width         int
	height        int
	showStatus    bool
	showMode      bool
	showHelp      bool
	helpFooter    *HelpFooterModel // Optional help footer
}

// NewFooter creates a new footer with width
func NewFooter(width int) FooterModel {
	return FooterModel{
		statusMessage: "",
		modeContext:   "",
		width:         width,
		height:        1,
		showStatus:    true,
		showMode:      true,
		showHelp:      false,
		helpFooter:    nil,
	}
}

// SetStatusMessage sets the status message
func (f *FooterModel) SetStatusMessage(message string) {
	f.statusMessage = message
}

// SetModeContext sets the mode/context string
func (f *FooterModel) SetModeContext(context string) {
	f.modeContext = context
}

// SetHelpFooter sets the help footer component
func (f *FooterModel) SetHelpFooter(helpFooter *HelpFooterModel) {
	f.helpFooter = helpFooter
	f.showHelp = true
}

// SetWidth sets the footer width
func (f *FooterModel) SetWidth(width int) {
	f.width = width
	if f.helpFooter != nil {
		f.helpFooter.width = width
	}
}

// SetHeight sets the footer height
func (f *FooterModel) SetHeight(height int) {
	f.height = height
}

// SetShowStatus sets whether to display the status message
func (f *FooterModel) SetShowStatus(show bool) {
	f.showStatus = show
}

// SetShowMode sets whether to display the mode context
func (f *FooterModel) SetShowMode(show bool) {
	f.showMode = show
}

// SetShowHelp sets whether to display the help footer
func (f *FooterModel) SetShowHelp(show bool) {
	f.showHelp = show && f.helpFooter != nil
}

// GetStatusMessage returns the current status message
func (f FooterModel) GetStatusMessage() string {
	return f.statusMessage
}

// GetModeContext returns the current mode context
func (f FooterModel) GetModeContext() string {
	return f.modeContext
}

// View renders the footer
func (f FooterModel) View() string {
	if f.width <= 0 {
		return ""
	}

	var parts []string

	// Render status and mode on top line if both present
	if f.showStatus && f.statusMessage != "" && f.showMode && f.modeContext != "" {
		statusModeLine := f.renderStatusAndMode()
		parts = append(parts, statusModeLine)
	} else {
		// Render them separately if only one is shown
		if f.showStatus && f.statusMessage != "" {
			parts = append(parts, f.renderStatus())
		}
		if f.showMode && f.modeContext != "" {
			parts = append(parts, f.renderMode())
		}
	}

	// Render help footer if present
	if f.showHelp && f.helpFooter != nil {
		helpView := f.helpFooter.View()
		if helpView != "" {
			parts = append(parts, helpView)
		}
	}

	return strings.Join(parts, "\n")
}

// renderStatusAndMode renders status and mode on the same line
func (f FooterModel) renderStatusAndMode() string {
	statusStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary)

	modeStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextMuted).
		Italic(true)

	// Format: "3/10 events  |  Capture Mode: Timeline"
	status := statusStyle.Render(f.statusMessage)
	mode := modeStyle.Render(f.modeContext)
	separator := styles.InfoHint.Render("  |  ")

	line := status + separator + mode

	// Truncate if too long
	if len(line) > f.width {
		// Just show status if combined is too long
		return status
	}

	return line
}

// renderStatus renders just the status message
func (f FooterModel) renderStatus() string {
	statusStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary)

	status := f.statusMessage
	if len(status) > f.width {
		status = status[:f.width-3] + "..."
	}

	return statusStyle.Render(status)
}

// renderMode renders just the mode context
func (f FooterModel) renderMode() string {
	modeStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextMuted).
		Italic(true)

	mode := f.modeContext
	if len(mode) > f.width {
		mode = mode[:f.width-3] + "..."
	}

	return modeStyle.Render(mode)
}
