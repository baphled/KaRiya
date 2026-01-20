package components

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/containers"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	tea "github.com/charmbracelet/bubbletea"
)

// ExportSuccessModal displays a success message after export completion.
// It shows export details and allows the user to acknowledge completion.
//
// Usage:
//
//	modal := components.NewExportSuccessModal("Career Events", "JSON", "File", "/path/to/file.json", 1024)
//	// In Update:
//	cmd, done := modal.Update(msg)
//	if done {
//	    // User acknowledged - continue workflow
//	}
type ExportSuccessModal struct {
	artifactType string // Type of artifact exported
	format       string // Export format
	destination  string // Export destination
	filePath     string // File path (for file exports) or "clipboard"
	size         int64  // Size in bytes
	visible      bool
	theme        themes.Theme
	width        int
	height       int
}

// NewExportSuccessModal creates a new export success modal.
//
// Parameters:
//   - artifactType: type of artifact exported (e.g., "Career Events")
//   - format: export format (e.g., "JSON", "CSV")
//   - destination: export destination (e.g., "File", "Clipboard")
//   - filePath: path to exported file or "clipboard"
//   - size: size of exported data in bytes
func NewExportSuccessModal(artifactType, format, destination, filePath string, size int64) *ExportSuccessModal {
	return &ExportSuccessModal{
		artifactType: artifactType,
		format:       format,
		destination:  destination,
		filePath:     filePath,
		size:         size,
		visible:      true,
		theme:        themes.NewDefaultTheme(),
		width:        100,
		height:       24,
	}
}

// Init initializes the modal (required by BubbleTea lifecycle).
func (m *ExportSuccessModal) Init() tea.Cmd {
	return nil
}

// Update handles keyboard input for the export success modal.
//
// Returns:
//   - tea.Cmd: command to execute (usually nil)
//   - bool: true if user acknowledged (pressed Enter or Esc)
func (m *ExportSuccessModal) Update(msg tea.Msg) (tea.Cmd, bool) {
	if !m.visible {
		return nil, false
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return nil, false

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyEsc:
			m.visible = false
			return nil, true
		}
	}

	return nil, false
}

// View renders the export success modal.
func (m *ExportSuccessModal) View() string {
	if !m.visible {
		return ""
	}

	// Build footer with primitives showing keyboard shortcuts
	footer := primitives.RenderHelpFooter(m.theme,
		primitives.HelpKeyBadge("Enter/Esc", "Done", m.theme),
	)

	// Build modal content
	var content strings.Builder

	// Success title with checkmark
	content.WriteString(primitives.SuccessText("✓ Export Complete", m.theme).Bold().MarginBottom(1).Render())
	content.WriteString("\n\n")

	// Export details
	content.WriteString(primitives.Body(fmt.Sprintf("Type:   %s", m.artifactType), m.theme).Width(50).Render())
	content.WriteString("\n")
	content.WriteString(primitives.Body(fmt.Sprintf("Format: %s", m.format), m.theme).Width(50).Render())
	content.WriteString("\n")

	// Location details based on destination
	if strings.ToLower(m.destination) == "clipboard" || m.filePath == "clipboard" {
		content.WriteString(primitives.Body("Data copied to clipboard", m.theme).Width(50).Render())
	} else {
		// Show filename (not full path for cleaner display)
		filename := filepath.Base(m.filePath)
		content.WriteString(primitives.Body(fmt.Sprintf("File:   %s", filename), m.theme).Width(50).Render())
	}
	content.WriteString("\n")

	// Size info
	sizeStr := formatBytes(m.size)
	content.WriteString(primitives.Muted(fmt.Sprintf("Size:   %s", sizeStr), m.theme).Width(50).Render())
	content.WriteString("\n\n")

	// Footer
	content.WriteString(footer)

	// Create modal box using UIKit
	modalWidth := 60
	if modalWidth > m.width-4 {
		modalWidth = m.width - 4
	}
	if modalWidth < 40 {
		modalWidth = 40
	}

	// Use UIKit Box with success variant
	return containers.NewBox(m.theme).
		Content(content.String()).
		Variant(containers.BoxDefault). // Default variant for success
		Width(modalWidth).
		Padding(2).
		Background(m.theme.BackgroundColor()).
		Render()
}

// formatBytes formats byte size to human-readable string.
func formatBytes(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%d bytes", bytes)
	}
}

// IsVisible returns whether the modal is currently visible.
func (m *ExportSuccessModal) IsVisible() bool {
	return m.visible
}

// Show makes the modal visible.
func (m *ExportSuccessModal) Show() {
	m.visible = true
}

// Hide hides the modal.
func (m *ExportSuccessModal) Hide() {
	m.visible = false
}

// SetTheme sets the theme for the modal.
func (m *ExportSuccessModal) SetTheme(theme themes.Theme) {
	m.theme = theme
}

// SetDimensions sets the terminal dimensions for responsive sizing.
func (m *ExportSuccessModal) SetDimensions(width, height int) {
	m.width = width
	m.height = height
}

// GetDimensions returns the current dimensions.
func (m *ExportSuccessModal) GetDimensions() (width, height int) {
	return m.width, m.height
}

// GetArtifactType returns the artifact type exported.
func (m *ExportSuccessModal) GetArtifactType() string {
	return m.artifactType
}

// GetFormat returns the export format.
func (m *ExportSuccessModal) GetFormat() string {
	return m.format
}

// GetDestination returns the export destination.
func (m *ExportSuccessModal) GetDestination() string {
	return m.destination
}

// GetFilePath returns the file path (or "clipboard").
func (m *ExportSuccessModal) GetFilePath() string {
	return m.filePath
}

// GetSize returns the export size in bytes.
func (m *ExportSuccessModal) GetSize() int64 {
	return m.size
}
