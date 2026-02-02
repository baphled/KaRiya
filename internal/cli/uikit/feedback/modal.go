package feedback

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/display"
)

// ModalSpinnerTickMsg is sent periodically to advance the spinner animation.
type ModalSpinnerTickMsg struct{}

// ModalType defines the type of modal.
type ModalType int

const (
	// ModalError displays an error message.
	ModalError ModalType = iota
	// ModalLoading displays a loading message with spinner.
	ModalLoading
	// ModalProgress displays progress with a progress bar.
	ModalProgress
	// ModalSuccess displays a success message.
	ModalSuccess
	// ModalWarning displays a warning message.
	ModalWarning
)

// SimpleSpinner provides a simple text-based spinner animation.
type SimpleSpinner struct {
	frames []string
	index  int
}

// NewSimpleSpinner creates a new simple spinner.
//
// Returns:
//   - A fully initialized SimpleSpinner ready for use.
//
// Side effects:
//   - None.
func NewSimpleSpinner() *SimpleSpinner {
	return &SimpleSpinner{
		frames: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		index:  0,
	}
}

// GetFrame returns the current spinner frame.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (s *SimpleSpinner) GetFrame() string {
	return s.frames[s.index]
}

// Advance advances the spinner to the next frame.
//
// Side effects:
//   - None.
func (s *SimpleSpinner) Advance() {
	s.index = (s.index + 1) % len(s.frames)
}

// LoadingMessageRotator rotates through a list of loading messages.
type LoadingMessageRotator struct {
	messages []string
	index    int
}

// NewLoadingMessageRotator creates a new message rotator.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized LoadingMessageRotator ready for use.
//
// Side effects:
//   - None.
func NewLoadingMessageRotator(messages []string) *LoadingMessageRotator {
	if len(messages) == 0 {
		messages = []string{"Loading..."}
	}
	return &LoadingMessageRotator{
		messages: messages,
		index:    0,
	}
}

// GetCurrent returns the current message.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (r *LoadingMessageRotator) GetCurrent() string {
	return r.messages[r.index]
}

// Rotate advances to the next message and returns it.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (r *LoadingMessageRotator) Rotate() string {
	r.index = (r.index + 1) % len(r.messages)
	return r.messages[r.index]
}

// Modal represents the content and configuration of a modal
// This is the UIKit version that uses theme-based styling exclusively.
type Modal struct {
	Type           ModalType
	Title          string
	Message        string
	Progress       float64
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

// getTheme returns the theme or default if nil.
func (m *Modal) getTheme() themes.Theme {
	if m.theme != nil {
		return m.theme
	}
	return themes.NewDefaultTheme()
}

// NewErrorModal creates a new error modal.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized Modal ready for use.
//
// Side effects:
//   - None.
func NewErrorModal(title, message string) *Modal {
	return &Modal{
		Type:           ModalError,
		Title:          title,
		Message:        message,
		FadeInDuration: 150 * time.Millisecond,
		Bell:           true,
		Cancellable:    true,
		fadeStartTime:  time.Now(),
		theme:          themes.NewDefaultTheme(),
	}
}

// NewLoadingModal creates a new loading modal with optional spinner.
//
// Expected:
//   - Must be a valid string.
//   - bool must be valid.
//
// Returns:
//   - A fully initialized Modal ready for use.
//
// Side effects:
//   - None.
func NewLoadingModal(message string, cancellable bool) *Modal {
	return &Modal{
		Type:           ModalLoading,
		Title:          "Loading",
		Message:        message,
		FadeInDuration: 150 * time.Millisecond,
		Cancellable:    cancellable,
		fadeStartTime:  time.Now(),
		spinner:        NewSimpleSpinner(),
		theme:          themes.NewDefaultTheme(),
	}
}

// NewProgressModal creates a new progress modal.
//
// Expected:
//   - Must be a valid string.
//   - float64 must be valid.
//
// Returns:
//   - A fully initialized Modal ready for use.
//
// Side effects:
//   - None.
func NewProgressModal(title, message string, progress float64) *Modal {
	return &Modal{
		Type:           ModalProgress,
		Title:          title,
		Message:        message,
		Progress:       progress,
		FadeInDuration: 150 * time.Millisecond,
		fadeStartTime:  time.Now(),
		theme:          themes.NewDefaultTheme(),
	}
}

// NewSuccessModal creates a new success modal with auto-dismiss.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized Modal ready for use.
//
// Side effects:
//   - None.
func NewSuccessModal(message string) *Modal {
	return &Modal{
		Type:           ModalSuccess,
		Title:          "Success",
		Message:        message,
		FadeInDuration: 150 * time.Millisecond,
		AutoDismiss:    3 * time.Second,
		fadeStartTime:  time.Now(),
		theme:          themes.NewDefaultTheme(),
	}
}

// NewWarningModal creates a new warning modal.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized Modal ready for use.
//
// Side effects:
//   - None.
func NewWarningModal(title, message string) *Modal {
	return &Modal{
		Type:           ModalWarning,
		Title:          title,
		Message:        message,
		FadeInDuration: 150 * time.Millisecond,
		Bell:           true,
		Cancellable:    true,
		fadeStartTime:  time.Now(),
		theme:          themes.NewDefaultTheme(),
	}
}

// SetMessageRotator sets a loading message rotator for dynamic messages.
//
// Expected:
//   - loadingmessagerotator must be valid.
//
// Returns:
//   - A fully initialized Modal ready for use.
//
// Side effects:
//   - None.
func (m *Modal) SetMessageRotator(rotator *LoadingMessageRotator) *Modal {
	m.messageRotator = rotator
	return m
}

// WithTheme sets the theme for the modal.
//
// Expected:
//   - th must be a valid theme instance (can be nil).
//
// Returns:
//   - A fully initialized Modal ready for use.
//
// Side effects:
//   - None.
func (m *Modal) WithTheme(theme themes.Theme) *Modal {
	m.theme = theme
	return m
}

// Render renders the modal centered in the given terminal dimensions.
//
// Expected:
//   - int must be valid.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *Modal) Render(terminalWidth, terminalHeight int) string {
	theme := m.getTheme()

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

	// Calculate maximum modal width based on terminal size.
	// Reserve 6 chars for border (2) + margin (4) to prevent cutoff.
	maxModalWidth := 100
	if terminalWidth > 0 && terminalWidth-6 < maxModalWidth {
		maxModalWidth = terminalWidth - 6
	}
	if maxModalWidth < 40 {
		maxModalWidth = 40
	}

	// Calculate text wrap width: modal width minus padding (4 chars).
	maxTextWidth := maxModalWidth - 4
	if maxTextWidth < 20 {
		maxTextWidth = 20
	}

	wrappedMessage := wrapText(message, maxTextWidth)
	contentParts = append(contentParts, wrappedMessage)

	// Add progress bar if progress modal.
	if m.Type == ModalProgress {
		contentParts = append(contentParts, "")
		progressBar := m.renderProgressBar(maxTextWidth, theme)
		contentParts = append(contentParts, progressBar)
	}

	// Add actions if present.
	if len(m.Actions) > 0 {
		contentParts = append(contentParts, "")
		actionsLine := strings.Join(m.Actions, "  ")
		contentParts = append(contentParts, actionsLine)
	}

	// Add dismissal hint.
	if m.Cancellable {
		contentParts = append(contentParts, "")
		switch m.Type {
		case ModalError:
			contentParts = append(contentParts, "Press Esc to dismiss")
		case ModalLoading:
			contentParts = append(contentParts, "Press Esc to cancel")
		default:
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

	// Join content.
	content := strings.Join(contentParts, "\n")

	// Calculate adaptive modal size based on actual content.
	contentLines := strings.Split(content, "\n")
	contentWidth := 0
	for _, line := range contentLines {
		lineWidth := lipgloss.Width(line)
		if lineWidth > contentWidth {
			contentWidth = lineWidth
		}
	}

	// Apply min/max constraints.
	// modalWidth is the lipgloss Width which includes padding but not border.
	modalWidth := contentWidth + 4
	if modalWidth < 40 {
		modalWidth = 40
	}
	if modalWidth > maxModalWidth {
		modalWidth = maxModalWidth
	}

	// Create modal box style.
	//
	// Width sets the content+padding width. Border is added outside.
	// We don't use MaxWidth as it would truncate the border characters.
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1, 2).
		Width(modalWidth)

	// Apply fade-in effect
	if opacity < 1.0 {
		boxStyle = boxStyle.Faint(true)
	}

	return boxStyle.Render(content)
}

// calculateOpacity calculates the current opacity based on fade-in duration.
func (m *Modal) calculateOpacity() float64 {
	if m.FadeInDuration == 0 {
		return 1.0
	}

	elapsed := time.Since(m.fadeStartTime)
	if elapsed >= m.FadeInDuration {
		return 1.0
	}

	return float64(elapsed) / float64(m.FadeInDuration)
}

// getStyleForType returns border color and icon for the modal type.
func (m *Modal) getStyleForType() (lipgloss.Color, string) {
	theme := m.getTheme()

	switch m.Type {
	case ModalError:
		return theme.ErrorColor(), "⚠️"
	case ModalLoading:
		return theme.InfoColor(), "⏳"
	case ModalProgress:
		return theme.InfoColor(), "📊"
	case ModalSuccess:
		return theme.SuccessColor(), "✅"
	case ModalWarning:
		return theme.WarningColor(), "⚠️"
	default:
		return theme.BorderColor(), "ℹ️"
	}
}

// renderProgressBar renders a progress bar for the given width.
func (m *Modal) renderProgressBar(width int, _ themes.Theme) string {
	percentage := int(m.Progress * 100)
	barWidth := width - 8

	filled := int(float64(barWidth) * m.Progress)
	if filled > barWidth {
		filled = barWidth
	}

	filledBar := strings.Repeat("█", filled)
	emptyBar := strings.Repeat("░", barWidth-filled)

	return fmt.Sprintf("[%s%s] %d%%", filledBar, emptyBar, percentage)
}

// wrapText wraps text to fit within the specified width.
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

// UpdateProgress updates the progress value (0.0 to 1.0).
//
// Expected:
//   - float64 must be valid.
//
// Side effects:
//   - None.
func (m *Modal) UpdateProgress(progress float64) {
	if progress < 0.0 {
		progress = 0.0
	}
	if progress > 1.0 {
		progress = 1.0
	}
	m.Progress = progress
}

// AdvanceSpinner advances the spinner to the next frame.
//
// Side effects:
//   - None.
func (m *Modal) AdvanceSpinner() {
	if m.spinner != nil {
		m.spinner.Advance()
	}
}

// RotateMessage advances to the next message if a rotator is set.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *Modal) RotateMessage() string {
	if m.messageRotator != nil {
		return m.messageRotator.Rotate()
	}
	return m.Message
}

// Init initializes the modal and starts spinner animation for loading modals.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (m *Modal) Init() tea.Cmd {
	if m.Type == ModalLoading && m.spinner != nil {
		return m.tickSpinner()
	}
	return nil
}

// Update handles messages for the modal, advancing the spinner on tick.
//
// Expected:
//   - msg must be valid.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (m *Modal) Update(msg tea.Msg) tea.Cmd {
	switch msg.(type) {
	case ModalSpinnerTickMsg:
		if m.Type == ModalLoading && m.spinner != nil {
			m.spinner.Advance()
			if m.messageRotator != nil {
				m.messageRotator.Rotate()
			}
			return m.tickSpinner()
		}
	}
	return nil
}

// tickSpinner returns a command to tick the spinner animation.
func (m *Modal) tickSpinner() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return ModalSpinnerTickMsg{}
	})
}

// ============================================================================
// Modal Overlay System
// ============================================================================

const (
	// MinOverlayWidth is the minimum width for overlay modals.
	MinOverlayWidth = 40
	// MaxOverlayWidth is the maximum width for overlay modals.
	MaxOverlayWidth = 120
	// DefaultOverlayWidth is the default width for overlay modals.
	DefaultOverlayWidth = 60
)

// OverlayModal represents a modal that renders as an overlay on top of dimmed content.
// It supports configurable width, title, content, and footer.
type OverlayModal struct {
	Title   string
	Content string
	Footer  string
	Width   int
	theme   themes.Theme
}

// NewOverlayModal creates a new overlay modal with the given title and content.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized OverlayModal ready for use.
//
// Side effects:
//   - None.
func NewOverlayModal(title, content string) *OverlayModal {
	return &OverlayModal{
		Title:   title,
		Content: content,
		Width:   DefaultOverlayWidth,
		theme:   themes.NewDefaultTheme(),
	}
}

// getTheme returns the theme or default if nil.
func (o *OverlayModal) getTheme() themes.Theme {
	if o.theme != nil {
		return o.theme
	}
	return themes.NewDefaultTheme()
}

// WithTheme sets the theme for the overlay modal.
//
// Expected:
//   - th must be a valid theme instance (can be nil).
//
// Returns:
//   - A fully initialized OverlayModal ready for use.
//
// Side effects:
//   - None.
func (o *OverlayModal) WithTheme(theme themes.Theme) *OverlayModal {
	o.theme = theme
	return o
}

// SetWidth sets the modal width, clamping to min/max bounds.
//
// Expected:
//   - int must be valid.
//
// Returns:
//   - A fully initialized OverlayModal ready for use.
//
// Side effects:
//   - None.
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
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized OverlayModal ready for use.
//
// Side effects:
//   - None.
func (o *OverlayModal) SetFooter(footer string) *OverlayModal {
	o.Footer = footer
	return o
}

// RenderCentered renders the modal centered over the dimmed background.
//
// Expected:
//   - Must be a valid string.
//   - int must be valid.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (o *OverlayModal) RenderCentered(background string, termWidth, termHeight int) string {
	return RenderOverlay(background, o.buildContent(), termWidth, termHeight, o.getTheme())
}

// buildContent assembles the modal content with title, body, and footer.
func (o *OverlayModal) buildContent() string {
	theme := o.getTheme()
	var parts []string

	// Add title with styling
	if o.Title != "" {
		titleStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.PrimaryColor()).
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
			Foreground(theme.MutedColor()).
			MarginTop(1)
		parts = append(parts, footerStyle.Render(o.Footer))
	}

	return strings.Join(parts, "\n")
}

// DimContent applies a dimmed/faded style to the content.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func DimContent(content string) string {
	if content == "" {
		return ""
	}

	dimStyle := lipgloss.NewStyle().Faint(true)
	return dimStyle.Render(content)
}

// RenderOverlay renders modal content centered over a dimmed background.
//
// Expected:
//   - Must be a valid string.
//   - int must be valid.
//   - th must be a valid theme instance (can be nil).
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func RenderOverlay(background, modalContent string, termWidth, termHeight int, theme themes.Theme) string {
	if theme == nil {
		theme = themes.NewDefaultTheme()
	}

	// Dim the entire background first
	dimmedBg := DimContent(background)

	// Create modal box with border
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.BorderColor()).
		Padding(1, 2).
		MaxWidth(MaxOverlayWidth)

	modalBox := modalStyle.Render(modalContent)

	// Split both into lines
	bgLines := strings.Split(dimmedBg, "\n")
	modalLines := strings.Split(modalBox, "\n")

	// Calculate modal position (below logo, not vertically centered)
	// Use display.DefaultLogoHeight as the single source of truth for logo dimensions.
	// Position modal to start just below the logo.
	modalHeight := len(modalLines)
	logoHeight := display.DefaultLogoHeight
	startY := logoHeight + 1

	// Handle very small terminals gracefully
	if termHeight < 15 {
		// For very small terminals, center the modal (no room for logo positioning)
		startY = (termHeight - modalHeight) / 2
		if startY < 0 {
			startY = 0
		}
	}

	// If modal is too tall to fit below logo, constrain it
	availableHeight := termHeight - startY - 2
	if availableHeight <= 0 {
		// Terminal too small - use all available space
		availableHeight = termHeight - 2
		if availableHeight < 1 {
			availableHeight = termHeight
		}
		startY = 0
	}

	if modalHeight > availableHeight && availableHeight > 0 {
		// Modal is too tall - truncate it and add scroll indicator
		if availableHeight <= len(modalLines) {
			modalLines = modalLines[:availableHeight]
			modalHeight = availableHeight
		}
		// Add scroll indicator at bottom
		if modalHeight > 0 {
			lastLine := modalLines[modalHeight-1]
			modalLines[modalHeight-1] = lastLine + " ↓"
		}
	}

	// Ensure we have exactly termHeight background lines
	for len(bgLines) < termHeight {
		bgLines = append(bgLines, strings.Repeat(" ", termWidth))
	}
	if len(bgLines) > termHeight {
		bgLines = bgLines[:termHeight]
	}

	// Create result by copying background
	result := make([]string, termHeight)
	copy(result, bgLines)

	// Overlay modal lines (centered horizontally) onto the background
	for i, modalLine := range modalLines {
		lineIndex := startY + i
		if lineIndex >= 0 && lineIndex < termHeight {
			// Use lipgloss.PlaceHorizontal to center the modal line
			// This creates a new line of exactly termWidth with the modal centered
			centeredModalLine := lipgloss.PlaceHorizontal(termWidth, lipgloss.Center, modalLine)
			result[lineIndex] = centeredModalLine
		}
	}

	return strings.Join(result, "\n")
}

// RenderOverlayWithDefaultTheme renders modal content using the default theme.
//
// Expected:
//   - Must be a valid string.
//   - int must be valid.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func RenderOverlayWithDefaultTheme(background, modalContent string, termWidth, termHeight int) string {
	return RenderOverlay(background, modalContent, termWidth, termHeight, themes.NewDefaultTheme())
}
