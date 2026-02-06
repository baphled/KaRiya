package modals

import (
	"fmt"
	"time"

	"github.com/baphled/kariya/internal/cli/uikit/containers"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/cli/uikit/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ProgressModal displays a progress indicator for async operations.
// It shows a spinner animation with title/subtitle and optional cancellation.
//
// Usage:
//
//	modal := modals.NewProgressModal("Processing", "Please wait...", true, 120, 40)
//	cmd := modal.Init()
//	// In Update:
//	cmd := modal.Update(msg)
//	if modal.IsCompleted() {
//	    // Operation complete
//	}
type ProgressModal struct {
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

// NewProgressModal creates a new progress modal with the given title,
//
// Expected:
//   - Must be a valid string.
//   - bool must be valid.
//   - int must be valid.
//
// Returns:
//   - A fully initialized ProgressModal ready for use.
//
// Side effects:
//   - None.
func NewProgressModal(title, subtitle string, cancellable bool, width, height int) *ProgressModal {
	return &ProgressModal{
		title:       title,
		subtitle:    subtitle,
		cancellable: cancellable,
		visible:     true,
		width:       width,
		height:      height,
	}
}

// NewExtractingTechsProgress creates a progress modal for technology extraction.
//
// Expected:
//   - int must be valid.
//
// Returns:
//   - A fully initialized ProgressModal ready for use.
//
// Side effects:
//   - None.
func NewExtractingTechsProgress(width, height int) *ProgressModal {
	return NewProgressModal(
		"Extracting Technologies",
		"Analyzing your skills...",
		true,
		width,
		height,
	)
}

// NewGeneratingCVProgress creates a progress modal for CV generation.
//
// Expected:
//   - Must be a valid string.
//   - int must be valid.
//
// Returns:
//   - A fully initialized ProgressModal ready for use.
//
// Side effects:
//   - None.
func NewGeneratingCVProgress(profile, audience string, width, height int) *ProgressModal {
	subtitle := fmt.Sprintf("Profile: %s | Audience: %s", profile, audience)
	return NewProgressModal(
		"Generating CV",
		subtitle,
		true,
		width,
		height,
	)
}

// NewExportingProgress creates a progress modal for CV export.
//
// Expected:
//   - Must be a valid string.
//   - int must be valid.
//
// Returns:
//   - A fully initialized ProgressModal ready for use.
//
// Side effects:
//   - None.
func NewExportingProgress(format string, width, height int) *ProgressModal {
	subtitle := "Format: " + format
	return NewProgressModal(
		"Exporting CV",
		subtitle,
		false, // Export is not cancellable
		width,
		height,
	)
}

// Init initializes the progress modal and starts spinner animation.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (m *ProgressModal) Init() tea.Cmd {
	return m.tickSpinner()
}

// Update handles messages for the progress modal.
//
// Expected:
//   - msg must be valid.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (m *ProgressModal) Update(msg tea.Msg) tea.Cmd {
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
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *ProgressModal) View() string {
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

	// Title with spinner using UIKit Text
	spinnerChar := spinnerFrames[m.spinner]
	titleText := fmt.Sprintf("%s  %s", spinnerChar, m.title)
	title := primitives.Title(titleText, th).
		Width(modalWidth - 4).
		Center().
		Render()

	// Subtitle using UIKit Text
	subtitle := primitives.Subtitle(m.subtitle, th).
		Width(modalWidth - 4).
		Center().
		Render()

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

	// Wrap in styled container with solid background using UIKit Box
	return containers.NewBox(th).
		Content(content).
		Width(modalWidth).
		Height(modalHeight).
		Padding(1).
		Background(th.BackgroundColor()).
		Variant(containers.BoxInfo). // Use accent color border
		Render()
}

// buildFooter creates the keyboard shortcuts footer using UIKit primitives.
func (m *ProgressModal) buildFooter() string {
	th := theme.Default()

	if m.cancellable {
		return primitives.CancelBadge(th).Render()
	}

	// Non-cancellable operations show no footer
	return ""
}

// tickSpinner returns a command to tick the spinner animation.
func (m *ProgressModal) tickSpinner() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(_ time.Time) tea.Msg {
		return SpinnerTickMsg{}
	})
}

// Show makes the modal visible.
//
// Side effects:
//   - None.
func (m *ProgressModal) Show() {
	m.visible = true
}

// Hide makes the modal invisible.
//
// Side effects:
//   - None.
func (m *ProgressModal) Hide() {
	m.visible = false
}

// IsVisible returns whether the modal is currently visible.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *ProgressModal) IsVisible() bool {
	return m.visible
}

// IsCancellable returns whether the operation can be cancelled.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *ProgressModal) IsCancellable() bool {
	return m.cancellable
}

// IsCompleted returns whether the operation has completed.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *ProgressModal) IsCompleted() bool {
	return m.completed
}

// IsCancelled returns whether the operation was cancelled by user.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *ProgressModal) IsCancelled() bool {
	return m.cancelled
}

// GetTitle returns the modal title.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *ProgressModal) GetTitle() string {
	return m.title
}

// GetSubtitle returns the modal subtitle.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *ProgressModal) GetSubtitle() string {
	return m.subtitle
}

// GetSpinnerFrame returns the current spinner frame index.
//
// Returns:
//   - A int value.
//
// Side effects:
//   - None.
func (m *ProgressModal) GetSpinnerFrame() int {
	return m.spinner
}

// GetError returns any error that occurred during the operation.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (m *ProgressModal) GetError() error {
	return m.err
}

// GetDimensions returns the current modal dimensions.
//
// Returns:
//   - width: The modal width in characters.
//   - height: The modal height in characters.
//
// Side effects:
//   - None.
func (m *ProgressModal) GetDimensions() (width, height int) {
	return m.width, m.height
}

// Complete marks the operation as completed and hides the modal.
//
// Side effects:
//   - None.
func (m *ProgressModal) Complete() {
	m.completed = true
	m.visible = false
}

// SetError sets an error for the operation and marks it as completed.
//
// Expected:
//   - error must be valid.
//
// Side effects:
//   - None.
func (m *ProgressModal) SetError(err error) {
	m.err = err
	m.completed = true
	m.visible = false
}
