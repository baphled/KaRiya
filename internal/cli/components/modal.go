package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/baphled/kariya/internal/cli/styles"
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
		return styles.ColorError, "⚠️"
	case ModalLoading:
		return styles.ColorInfo, "⏳"
	case ModalProgress:
		return styles.ColorInfo, "📊"
	case ModalSuccess:
		return styles.ColorSuccess, "✅"
	case ModalWarning:
		return styles.ColorWarning, "⚠️"
	default:
		return styles.ColorBorder, "ℹ️"
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
