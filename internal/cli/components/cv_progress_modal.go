package components

import (
	"fmt"
	"time"

	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/cli/uikit/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// CVProgressModal displays a progress indicator for async operations.
// It shows a spinner animation with title/subtitle and optional cancellation.
//
// Usage:
//
//	modal := components.NewCVProgressModal("Processing", "Please wait...", true, 120, 40)
//	cmd := modal.Init()
//	// In Update:
//	cmd := modal.Update(msg)
//	if modal.IsCompleted() {
//	    // Operation complete
//	}
type CVProgressModal struct {
	title       string
	subtitle    string
	spinner     int // Current spinner frame (0-9)
	visible     bool
	cancellable bool
	completed   bool
	cancelled   bool
	err         error
	width       int
	height      int
}

// SpinnerTickMsg is sent periodically to advance the spinner animation.
type SpinnerTickMsg struct{}

// spinnerFrames defines the 10-frame spinner animation.
var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// NewCVProgressModal creates a new progress modal.
// title: Main title text
// subtitle: Descriptive subtitle text
// cancellable: Whether user can press Esc to cancel
// width, height: Terminal dimensions
func NewCVProgressModal(title, subtitle string, cancellable bool, width, height int) *CVProgressModal {
	return &CVProgressModal{
		title:       title,
		subtitle:    subtitle,
		cancellable: cancellable,
		visible:     true,
		width:       width,
		height:      height,
	}
}

// NewExtractingTechsProgress creates a progress modal for technology extraction.
func NewExtractingTechsProgress(width, height int) *CVProgressModal {
	return NewCVProgressModal(
		"Extracting Technologies",
		"Analyzing your skills...",
		true,
		width,
		height,
	)
}

// NewGeneratingCVProgress creates a progress modal for CV generation.
func NewGeneratingCVProgress(profile, audience string, width, height int) *CVProgressModal {
	subtitle := fmt.Sprintf("Profile: %s | Audience: %s", profile, audience)
	return NewCVProgressModal(
		"Generating CV",
		subtitle,
		true,
		width,
		height,
	)
}

// NewExportingProgress creates a progress modal for CV export.
func NewExportingProgress(format string, width, height int) *CVProgressModal {
	subtitle := fmt.Sprintf("Format: %s", format)
	return NewCVProgressModal(
		"Exporting CV",
		subtitle,
		false, // Export is not cancellable
		width,
		height,
	)
}

// Init initializes the progress modal and starts spinner animation.
func (m *CVProgressModal) Init() tea.Cmd {
	return m.tickSpinner()
}

// Update handles messages for the progress modal.
func (m *CVProgressModal) Update(msg tea.Msg) tea.Cmd {
	if !m.visible {
		return nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc && m.cancellable && !m.completed {
			m.cancelled = true
			m.visible = false
			return nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case SpinnerTickMsg:
		if !m.completed {
			m.spinner = (m.spinner + 1) % len(spinnerFrames)
			return m.tickSpinner()
		}
	}

	return nil
}

// View renders the progress modal.
func (m *CVProgressModal) View() string {
	if !m.visible {
		return ""
	}

	// Calculate modal dimensions
	modalWidth := m.width - 20
	if modalWidth > 60 {
		modalWidth = 60
	}
	if modalWidth < 40 {
		modalWidth = 40
	}

	modalHeight := 10 // Fixed height for progress modal

	// Use default theme for colors
	th := theme.Default()

	// Title with spinner
	titleStyle := lipgloss.NewStyle().
		Foreground(th.AccentColor()).
		Bold(true).
		Align(lipgloss.Center).
		Width(modalWidth - 4)

	spinnerChar := spinnerFrames[m.spinner]
	titleText := fmt.Sprintf("%s  %s", spinnerChar, m.title)
	title := titleStyle.Render(titleText)

	// Subtitle
	subtitleStyle := lipgloss.NewStyle().
		Foreground(th.SecondaryColor()).
		Align(lipgloss.Center).
		Width(modalWidth - 4)

	subtitle := subtitleStyle.Render(m.subtitle)

	// Footer with KeyBadge
	footer := m.buildFooter()

	// Combine all parts
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		title,
		"",
		subtitle,
		"",
		"",
		footer,
	)

	// Wrap in styled container with solid background
	styledContent := lipgloss.NewStyle().
		Width(modalWidth).
		Height(modalHeight).
		Background(th.BackgroundColor()). // CRITICAL: Solid background
		Border(lipgloss.RoundedBorder()).
		BorderForeground(th.AccentColor()).
		Padding(1).
		Render(content)

	return styledContent
}

// buildFooter creates the keyboard shortcuts footer using UIKit primitives.
func (m *CVProgressModal) buildFooter() string {
	th := theme.Default()

	if m.cancellable {
		return primitives.CancelBadge(th).Render()
	}

	// Non-cancellable operations show no footer
	return ""
}

// tickSpinner returns a command to tick the spinner animation.
func (m *CVProgressModal) tickSpinner() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return SpinnerTickMsg{}
	})
}

// Show makes the modal visible.
func (m *CVProgressModal) Show() {
	m.visible = true
}

// Hide makes the modal invisible.
func (m *CVProgressModal) Hide() {
	m.visible = false
}

// IsVisible returns whether the modal is currently visible.
func (m *CVProgressModal) IsVisible() bool {
	return m.visible
}

// IsCancellable returns whether the operation can be cancelled.
func (m *CVProgressModal) IsCancellable() bool {
	return m.cancellable
}

// IsCompleted returns whether the operation has completed.
func (m *CVProgressModal) IsCompleted() bool {
	return m.completed
}

// IsCancelled returns whether the operation was cancelled by user.
func (m *CVProgressModal) IsCancelled() bool {
	return m.cancelled
}

// GetTitle returns the modal title.
func (m *CVProgressModal) GetTitle() string {
	return m.title
}

// GetSubtitle returns the modal subtitle.
func (m *CVProgressModal) GetSubtitle() string {
	return m.subtitle
}

// GetSpinnerFrame returns the current spinner frame index.
func (m *CVProgressModal) GetSpinnerFrame() int {
	return m.spinner
}

// GetError returns any error that occurred during the operation.
func (m *CVProgressModal) GetError() error {
	return m.err
}

// GetDimensions returns the current modal dimensions.
func (m *CVProgressModal) GetDimensions() (width, height int) {
	return m.width, m.height
}

// Complete marks the operation as completed and hides the modal.
func (m *CVProgressModal) Complete() {
	m.completed = true
	m.visible = false
}

// SetError sets an error for the operation and marks it as completed.
func (m *CVProgressModal) SetError(err error) {
	m.err = err
	m.completed = true
	m.visible = false
}
