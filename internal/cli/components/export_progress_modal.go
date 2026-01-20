package components

import (
	"fmt"
	"time"

	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/containers"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ExportProgressModal displays a progress indicator for export operations.
// It shows a spinner animation with export details and cancellation option.
//
// Usage:
//
//	modal := components.NewExportProgressModal("Career Events", "JSON", "File", 120, 40)
//	cmd := modal.Init()
//	// In Update:
//	cmd := modal.Update(msg)
//	if modal.IsCompleted() {
//	    // Operation complete
//	} else if modal.IsCancelled() {
//	    // User cancelled
//	}
type ExportProgressModal struct {
	artifactType string // Type of artifact being exported
	format       string // Export format
	destination  string // Export destination
	spinner      int    // Current spinner frame (0-9)
	visible      bool
	completed    bool
	cancelled    bool
	err          error
	theme        themes.Theme
	width        int
	height       int
}

// NewExportProgressModal creates a new export progress modal.
//
// Parameters:
//   - artifactType: type of artifact being exported (e.g., "Career Events")
//   - format: export format (e.g., "JSON", "CSV")
//   - destination: export destination (e.g., "File", "Clipboard")
//   - width, height: terminal dimensions
func NewExportProgressModal(artifactType, format, destination string, width, height int) *ExportProgressModal {
	return &ExportProgressModal{
		artifactType: artifactType,
		format:       format,
		destination:  destination,
		visible:      true,
		theme:        themes.NewDefaultTheme(),
		width:        width,
		height:       height,
	}
}

// Init initializes the progress modal and starts spinner animation.
func (m *ExportProgressModal) Init() tea.Cmd {
	return m.tickSpinner()
}

// Update handles messages for the progress modal.
func (m *ExportProgressModal) Update(msg tea.Msg) tea.Cmd {
	if !m.visible {
		return nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc && !m.completed {
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
func (m *ExportProgressModal) View() string {
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

	modalHeight := 12 // Fixed height for progress modal

	// Title with spinner
	spinnerChar := spinnerFrames[m.spinner]
	titleText := fmt.Sprintf("%s  Exporting %s", spinnerChar, m.artifactType)
	title := primitives.Title(titleText, m.theme).
		Width(modalWidth - 4).
		Center().
		Render()

	// Export details
	detailsText := fmt.Sprintf("Format: %s | Destination: %s", m.format, m.destination)
	details := primitives.Subtitle(detailsText, m.theme).
		Width(modalWidth - 4).
		Center().
		Render()

	// Please wait message
	waitMsg := primitives.Muted("Please wait...", m.theme).
		Width(modalWidth - 4).
		Center().
		Render()

	// Footer with cancel badge
	footer := primitives.RenderHelpFooter(m.theme,
		primitives.HelpKeyBadge("Esc", "Cancel", m.theme),
	)

	// Combine all parts
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		title,
		"",
		details,
		"",
		waitMsg,
		"",
		footer,
	)

	// Wrap in styled container with solid background using UIKit Box
	return containers.NewBox(m.theme).
		Content(content).
		Width(modalWidth).
		Height(modalHeight).
		Padding(1).
		Background(m.theme.BackgroundColor()).
		Variant(containers.BoxInfo). // Use accent color border
		Render()
}

// tickSpinner returns a command to tick the spinner animation.
func (m *ExportProgressModal) tickSpinner() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return SpinnerTickMsg{}
	})
}

// Show makes the modal visible.
func (m *ExportProgressModal) Show() {
	m.visible = true
}

// Hide makes the modal invisible.
func (m *ExportProgressModal) Hide() {
	m.visible = false
}

// IsVisible returns whether the modal is currently visible.
func (m *ExportProgressModal) IsVisible() bool {
	return m.visible
}

// IsCompleted returns whether the operation has completed.
func (m *ExportProgressModal) IsCompleted() bool {
	return m.completed
}

// IsCancelled returns whether the operation was cancelled by user.
func (m *ExportProgressModal) IsCancelled() bool {
	return m.cancelled
}

// Complete marks the operation as completed and hides the modal.
func (m *ExportProgressModal) Complete() {
	m.completed = true
	m.visible = false
}

// SetError sets an error for the operation and marks it as completed.
func (m *ExportProgressModal) SetError(err error) {
	m.err = err
	m.completed = true
	m.visible = false
}

// GetError returns any error that occurred during the operation.
func (m *ExportProgressModal) GetError() error {
	return m.err
}

// GetSpinnerFrame returns the current spinner frame index.
func (m *ExportProgressModal) GetSpinnerFrame() int {
	return m.spinner
}

// GetDimensions returns the current modal dimensions.
func (m *ExportProgressModal) GetDimensions() (width, height int) {
	return m.width, m.height
}

// GetArtifactType returns the artifact type being exported.
func (m *ExportProgressModal) GetArtifactType() string {
	return m.artifactType
}

// GetFormat returns the export format.
func (m *ExportProgressModal) GetFormat() string {
	return m.format
}

// GetDestination returns the export destination.
func (m *ExportProgressModal) GetDestination() string {
	return m.destination
}

// SetTheme sets the theme for the modal.
func (m *ExportProgressModal) SetTheme(theme themes.Theme) {
	m.theme = theme
}
