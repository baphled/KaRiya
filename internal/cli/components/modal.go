package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/cli/themes"
)

// ModalType defines the type of modal
type ModalType int

const (
	// ModalError displays an error message
	ModalError ModalType = iota
	// ModalLoading displays a loading message with spinner
	ModalLoading
	// ModalProgress displays progress with a progress bar
	ModalProgress
	// ModalSuccess displays a success message
	ModalSuccess
	// ModalWarning displays a warning message
	ModalWarning
)

// ModalContent represents the content and configuration of a modal
type ModalContent struct {
	Type           ModalType
	Title          string
	Message        string
	Progress       float64 // 0.0 to 1.0
	Actions        []string
	FadeInDuration time.Duration
	AutoDismiss    time.Duration
	Bell           bool
	Cancellable    bool
	fadeStartTime  time.Time
	spinner        *SimpleSpinner
	messageRotator *LoadingMessageRotator
	theme          themes.Theme
}

// NewErrorModal creates a new error modal
func NewErrorModal(title, message string) *ModalContent {
	return &ModalContent{
		Type:           ModalError,
		Title:          title,
		Message:        message,
		FadeInDuration: 150 * time.Millisecond,
		Bell:           true,
		Cancellable:    true,
		fadeStartTime:  time.Now(),
	}
}

// NewLoadingModal creates a new loading modal with optional spinner
func NewLoadingModal(message string, cancellable bool) *ModalContent {
	return &ModalContent{
		Type:           ModalLoading,
		Title:          "Loading",
		Message:        message,
		FadeInDuration: 150 * time.Millisecond,
		Cancellable:    cancellable,
		fadeStartTime:  time.Now(),
		spinner:        NewSimpleSpinner(),
	}
}

// NewProgressModal creates a new progress modal
func NewProgressModal(title, message string, progress float64) *ModalContent {
	return &ModalContent{
		Type:           ModalProgress,
		Title:          title,
		Message:        message,
		Progress:       progress,
		FadeInDuration: 150 * time.Millisecond,
		fadeStartTime:  time.Now(),
	}
}

// NewSuccessModal creates a new success modal with auto-dismiss
func NewSuccessModal(message string) *ModalContent {
	return &ModalContent{
		Type:           ModalSuccess,
		Title:          "Success",
		Message:        message,
		FadeInDuration: 150 * time.Millisecond,
		AutoDismiss:    3 * time.Second,
		fadeStartTime:  time.Now(),
	}
}

// NewWarningModal creates a new warning modal
func NewWarningModal(title, message string) *ModalContent {
	return &ModalContent{
		Type:           ModalWarning,
		Title:          title,
		Message:        message,
		FadeInDuration: 150 * time.Millisecond,
		Bell:           true,
		Cancellable:    true,
		fadeStartTime:  time.Now(),
	}
}

// SetMessageRotator sets a loading message rotator for dynamic messages
func (m *ModalContent) SetMessageRotator(rotator *LoadingMessageRotator) *ModalContent {
	m.messageRotator = rotator
	return m
}

// WithTheme sets the theme for the modal
func (m *ModalContent) WithTheme(theme themes.Theme) *ModalContent {
	m.theme = theme
	return m
}

// Theme helper methods for consistent themed styling.

// getErrorColor returns the error color from theme or fallback.
func (m *ModalContent) getErrorColor() lipgloss.Color {
	if m.theme != nil {
		return m.theme.ErrorColor()
	}
	return styles.ColorError
}

// getInfoColor returns the info color from theme or fallback.
func (m *ModalContent) getInfoColor() lipgloss.Color {
	if m.theme != nil {
		return m.theme.InfoColor()
	}
	return styles.ColorInfo
}

// getSuccessColor returns the success color from theme or fallback.
func (m *ModalContent) getSuccessColor() lipgloss.Color {
	if m.theme != nil {
		return m.theme.SuccessColor()
	}
	return styles.ColorSuccess
}

// getWarningColor returns the warning color from theme or fallback.
func (m *ModalContent) getWarningColor() lipgloss.Color {
	if m.theme != nil {
		return m.theme.WarningColor()
	}
	return styles.ColorWarning
}

// getBorderColor returns the border color from theme or fallback.
func (m *ModalContent) getBorderColor() lipgloss.Color {
	if m.theme != nil {
		return m.theme.BorderColor()
	}
	return styles.ColorBorder
}

// Render renders the modal centered in the given terminal dimensions
func (m *ModalContent) Render(terminalWidth, terminalHeight int) string {
	// Calculate fade-in opacity
	opacity := m.calculateOpacity()

	// Get modal styling based on type
	borderColor, icon := m.getStyleForType()

	// Build modal content
	var contentParts []string

	// Add icon and title
	if m.Title != "" {
		titleLine := fmt.Sprintf("%s %s", icon, m.Title)
		contentParts = append(contentParts, titleLine)
		contentParts = append(contentParts, "")
	}

	// Add message (with spinner if loading)
	message := m.Message
	if m.Type == ModalLoading && m.spinner != nil {
		spinnerFrame := m.spinner.GetFrame()
		message = fmt.Sprintf("%s %s", spinnerFrame, m.Message)
	}

	// Use rotator message if available
	if m.messageRotator != nil {
		rotatedMessage := m.messageRotator.GetCurrent()
		if m.Type == ModalLoading && m.spinner != nil {
			message = fmt.Sprintf("%s %s", m.spinner.GetFrame(), rotatedMessage)
		} else {
			message = rotatedMessage
		}
	}

	// Wrap message to fit modal width
	maxWidth := 76 // max width minus padding and borders
	wrappedMessage := wrapText(message, maxWidth)
	contentParts = append(contentParts, wrappedMessage)

	// Add progress bar if progress modal
	if m.Type == ModalProgress {
		contentParts = append(contentParts, "")
		progressBar := m.renderProgressBar(maxWidth)
		contentParts = append(contentParts, progressBar)
	}

	// Add actions if present
	if len(m.Actions) > 0 {
		contentParts = append(contentParts, "")
		actionsLine := strings.Join(m.Actions, "  ")
		contentParts = append(contentParts, actionsLine)
	}

	// Add dismissal hint
	if m.Cancellable {
		contentParts = append(contentParts, "")
		if m.Type == ModalError {
			contentParts = append(contentParts, "Press Esc to dismiss")
		} else if m.Type == ModalLoading {
			contentParts = append(contentParts, "Press Esc to cancel")
		} else {
			contentParts = append(contentParts, "Press Esc to close")
		}
	}

	if m.AutoDismiss > 0 {
		contentParts = append(contentParts, "")
		elapsed := time.Since(m.fadeStartTime)
		remaining := m.AutoDismiss - elapsed
		if remaining > 0 {
			contentParts = append(contentParts, fmt.Sprintf("Auto-dismiss in %.0fs", remaining.Seconds()))
		}
	}

	// Join content
	content := strings.Join(contentParts, "\n")

	// Calculate adaptive modal size
	contentLines := strings.Split(content, "\n")
	contentWidth := 0
	for _, line := range contentLines {
		lineWidth := lipgloss.Width(line)
		if lineWidth > contentWidth {
			contentWidth = lineWidth
		}
	}

	// Apply min/max constraints
	modalWidth := contentWidth + 4 // Add padding
	if modalWidth < 40 {
		modalWidth = 40
	}
	if modalWidth > 80 {
		modalWidth = 80
	}

	modalHeight := len(contentLines) + 4 // Add padding
	if modalHeight < 8 {
		modalHeight = 8
	}
	if modalHeight > 20 {
		modalHeight = 20
	}

	// Create modal box style
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1, 2).
		Width(modalWidth).
		MaxWidth(modalWidth)

	// Apply fade-in effect
	if opacity < 1.0 {
		boxStyle = boxStyle.Faint(true)
	}

	return boxStyle.Render(content)
}

// calculateOpacity calculates the current opacity based on fade-in duration
func (m *ModalContent) calculateOpacity() float64 {
	if m.FadeInDuration == 0 {
		return 1.0
	}

	elapsed := time.Since(m.fadeStartTime)
	if elapsed >= m.FadeInDuration {
		return 1.0
	}

	return float64(elapsed) / float64(m.FadeInDuration)
}

// getStyleForType returns border color and icon for the modal type
func (m *ModalContent) getStyleForType() (lipgloss.Color, string) {
	switch m.Type {
	case ModalError:
		return m.getErrorColor(), "⚠️"
	case ModalLoading:
		return m.getInfoColor(), "⏳"
	case ModalProgress:
		return m.getInfoColor(), "📊"
	case ModalSuccess:
		return m.getSuccessColor(), "✅"
	case ModalWarning:
		return m.getWarningColor(), "⚠️"
	default:
		return m.getBorderColor(), "ℹ️"
	}
}

// renderProgressBar renders a progress bar for the given width
func (m *ModalContent) renderProgressBar(width int) string {
	percentage := int(m.Progress * 100)
	barWidth := width - 8 // Leave space for percentage and brackets

	filled := int(float64(barWidth) * m.Progress)
	if filled > barWidth {
		filled = barWidth
	}

	filledBar := strings.Repeat("█", filled)
	emptyBar := strings.Repeat("░", barWidth-filled)

	return fmt.Sprintf("[%s%s] %d%%", filledBar, emptyBar, percentage)
}

// wrapText wraps text to fit within the specified width
func wrapText(text string, width int) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	var lines []string
	var currentLine string

	for _, word := range words {
		testLine := currentLine
		if testLine != "" {
			testLine += " "
		}
		testLine += word

		if lipgloss.Width(testLine) <= width {
			currentLine = testLine
		} else {
			if currentLine != "" {
				lines = append(lines, currentLine)
			}
			currentLine = word
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return strings.Join(lines, "\n")
}

// UpdateProgress updates the progress value (0.0 to 1.0)
func (m *ModalContent) UpdateProgress(progress float64) {
	if progress < 0.0 {
		progress = 0.0
	}
	if progress > 1.0 {
		progress = 1.0
	}
	m.Progress = progress
}

// AdvanceSpinner advances the spinner to the next frame
func (m *ModalContent) AdvanceSpinner() {
	if m.spinner != nil {
		m.spinner.Advance()
	}
}

// RotateMessage advances to the next message if a rotator is set
func (m *ModalContent) RotateMessage() string {
	if m.messageRotator != nil {
		return m.messageRotator.Rotate()
	}
	return m.Message
}

// ============================================================================
// Modal Overlay System
// ============================================================================

const (
	// MinOverlayWidth is the minimum width for overlay modals
	MinOverlayWidth = 40
	// MaxOverlayWidth is the maximum width for overlay modals
	MaxOverlayWidth = 120
	// DefaultOverlayWidth is the default width for overlay modals
	DefaultOverlayWidth = 60
)

// OverlayModal represents a modal that renders as an overlay on top of dimmed content.
// It supports configurable width, title, content, and footer.
type OverlayModal struct {
	Title   string
	Content string
	Footer  string
	Width   int
}

// NewOverlayModal creates a new overlay modal with the given title and content.
func NewOverlayModal(title, content string) *OverlayModal {
	return &OverlayModal{
		Title:   title,
		Content: content,
		Width:   DefaultOverlayWidth,
	}
}

// SetWidth sets the modal width, clamping to min/max bounds.
func (o *OverlayModal) SetWidth(width int) *OverlayModal {
	if width < MinOverlayWidth {
		width = MinOverlayWidth
	}
	if width > MaxOverlayWidth {
		width = MaxOverlayWidth
	}
	o.Width = width
	return o
}

// SetFooter sets the footer text for the modal.
func (o *OverlayModal) SetFooter(footer string) *OverlayModal {
	o.Footer = footer
	return o
}

// RenderCentered renders the modal centered over the dimmed background.
func (o *OverlayModal) RenderCentered(background string, termWidth, termHeight int) string {
	return RenderOverlay(background, o.buildContent(), termWidth, termHeight)
}

// buildContent assembles the modal content with title, body, and footer.
func (o *OverlayModal) buildContent() string {
	var parts []string

	// Add title with styling
	if o.Title != "" {
		titleStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(styles.ColorTextPrimary).
			MarginBottom(1)
		parts = append(parts, titleStyle.Render(o.Title))
	}

	// Add content
	if o.Content != "" {
		parts = append(parts, o.Content)
	}

	// Add footer with muted styling
	if o.Footer != "" {
		footerStyle := lipgloss.NewStyle().
			Foreground(styles.ColorTextMuted).
			MarginTop(1)
		parts = append(parts, footerStyle.Render(o.Footer))
	}

	return strings.Join(parts, "\n")
}

// DimContent applies a dimmed/faded style to the content.
// This is used to visually distinguish the background from the modal overlay.
func DimContent(content string) string {
	if content == "" {
		return ""
	}

	dimStyle := lipgloss.NewStyle().Faint(true)
	return dimStyle.Render(content)
}

// RenderOverlay renders modal content centered over a dimmed background.
// It handles the centering calculation and compositing.
func RenderOverlay(background, modalContent string, termWidth, termHeight int) string {
	// Dim the background
	dimmedBg := DimContent(background)

	// Create modal box with border
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Padding(1, 2).
		MaxWidth(MaxOverlayWidth)

	modalBox := modalStyle.Render(modalContent)

	// Calculate modal dimensions
	modalHeight := lipgloss.Height(modalBox)
	modalWidth := lipgloss.Width(modalBox)

	// Calculate center position
	centerX := (termWidth - modalWidth) / 2
	centerY := (termHeight - modalHeight) / 2

	// Ensure non-negative positions
	if centerX < 0 {
		centerX = 0
	}
	if centerY < 0 {
		centerY = 0
	}

	// Split background into lines
	bgLines := strings.Split(dimmedBg, "\n")

	// Ensure we have enough background lines
	for len(bgLines) < termHeight {
		bgLines = append(bgLines, "")
	}

	// Split modal into lines
	modalLines := strings.Split(modalBox, "\n")

	// Overlay modal onto background
	for i, modalLine := range modalLines {
		bgLineIdx := centerY + i
		if bgLineIdx < 0 || bgLineIdx >= len(bgLines) {
			continue
		}

		// Ensure background line is wide enough
		bgLine := bgLines[bgLineIdx]
		for lipgloss.Width(bgLine) < centerX {
			bgLine += " "
		}

		// Build the new line: prefix + modal line + suffix
		prefix := truncateToWidth(bgLine, centerX)
		suffix := ""
		afterModal := centerX + lipgloss.Width(modalLine)
		if lipgloss.Width(bgLine) > afterModal {
			suffix = substringFromWidth(bgLine, afterModal)
		}

		bgLines[bgLineIdx] = prefix + modalLine + suffix
	}

	return strings.Join(bgLines[:termHeight], "\n")
}

// truncateToWidth truncates a string to fit within the specified width.
func truncateToWidth(s string, width int) string {
	if width <= 0 {
		return ""
	}

	result := ""
	currentWidth := 0

	for _, r := range s {
		charWidth := lipgloss.Width(string(r))
		if currentWidth+charWidth > width {
			break
		}
		result += string(r)
		currentWidth += charWidth
	}

	// Pad with spaces if needed
	for currentWidth < width {
		result += " "
		currentWidth++
	}

	return result
}

// substringFromWidth returns the substring starting from the specified width position.
func substringFromWidth(s string, startWidth int) string {
	if startWidth <= 0 {
		return s
	}

	currentWidth := 0
	for i, r := range s {
		charWidth := lipgloss.Width(string(r))
		if currentWidth >= startWidth {
			return s[i:]
		}
		currentWidth += charWidth
	}

	return ""
}
