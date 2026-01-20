package components

import (
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/containers"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	tea "github.com/charmbracelet/bubbletea"
)

// ExportConfirmModal is a confirmation modal for export operations.
// It displays export details and asks for user confirmation.
//
// Usage:
//
//	modal := components.NewExportConfirmModal("Career Events", "JSON", "File")
//	cmd := modal.Init()
//	// In Update:
//	cmd, confirmed := modal.Update(msg)
//	if confirmed {
//	    // User confirmed - proceed with export
//	} else if !modal.IsVisible() {
//	    // User cancelled - abort export
//	}
type ExportConfirmModal struct {
	artifactType string // Type of artifact being exported
	format       string // Export format
	destination  string // Export destination
	visible      bool   // Whether modal is currently visible
	confirmed    bool   // Whether user confirmed export
	theme        themes.Theme
	width        int
	height       int
}

// NewExportConfirmModal creates a new export confirmation modal.
//
// Parameters:
//   - artifactType: type of artifact being exported (e.g., "Career Events")
//   - format: export format (e.g., "JSON", "CSV")
//   - destination: export destination (e.g., "File", "Clipboard")
//
// Returns a new ExportConfirmModal ready to use.
func NewExportConfirmModal(artifactType, format, destination string) *ExportConfirmModal {
	return &ExportConfirmModal{
		artifactType: artifactType,
		format:       format,
		destination:  destination,
		visible:      true,
		confirmed:    false,
		theme:        themes.NewDefaultTheme(),
		width:        100,
		height:       24,
	}
}

// Init initializes the modal (required by BubbleTea lifecycle).
func (m *ExportConfirmModal) Init() tea.Cmd {
	return nil
}

// Update handles keyboard input for the export confirmation modal.
//
// Returns:
//   - tea.Cmd: command to execute (usually nil)
//   - bool: true if user confirmed export, false otherwise
//
// The modal closes on:
//   - 'y' or Enter: confirms export (returns true)
//   - 'n' or Esc: cancels export (returns false)
func (m *ExportConfirmModal) Update(msg tea.Msg) (tea.Cmd, bool) {
	if !m.visible {
		return nil, false
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return nil, false

	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y", "enter":
			// User confirmed export
			m.confirmed = true
			m.visible = false
			return nil, true

		case "n", "N", "esc":
			// User cancelled export
			m.confirmed = false
			m.visible = false
			return nil, false
		}
	}

	return nil, false
}

// View renders the export confirmation modal as a centered overlay.
//
// The modal displays:
//   - Title "Confirm Export"
//   - Export details (artifact type, format, destination)
//   - Action instructions (y/Enter: Confirm, n/Esc: Cancel)
//
// Returns empty string if modal is not visible.
func (m *ExportConfirmModal) View() string {
	if !m.visible {
		return ""
	}

	// Build footer with primitives showing keyboard shortcuts
	footer := primitives.RenderHelpFooter(m.theme,
		primitives.HelpKeyBadge("y/Enter", "Confirm", m.theme),
		primitives.HelpKeyBadge("n/Esc", "Cancel", m.theme),
	)

	// Build modal content
	var content strings.Builder

	// Title using UIKit Text
	content.WriteString(primitives.Title("Confirm Export", m.theme).Bold().MarginBottom(1).Render())
	content.WriteString("\n\n")

	// Export details
	content.WriteString(primitives.Body(fmt.Sprintf("Type:        %s", m.artifactType), m.theme).Width(50).Render())
	content.WriteString("\n")
	content.WriteString(primitives.Body(fmt.Sprintf("Format:      %s", m.format), m.theme).Width(50).Render())
	content.WriteString("\n")
	content.WriteString(primitives.Body(fmt.Sprintf("Destination: %s", m.destination), m.theme).Width(50).Render())
	content.WriteString("\n\n")

	// Confirmation message
	content.WriteString(primitives.Muted("Are you sure you want to export?", m.theme).Render())
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

	// Use UIKit Box with emphasized variant (accent-colored border)
	return containers.NewBox(m.theme).
		Content(content.String()).
		Variant(containers.BoxEmphasized).
		Width(modalWidth).
		Padding(2).
		Background(m.theme.BackgroundColor()).
		Render()
}

// IsVisible returns whether the modal is currently visible.
func (m *ExportConfirmModal) IsVisible() bool {
	return m.visible
}

// WasConfirmed returns whether the user confirmed the export.
// Only meaningful after the modal is closed.
func (m *ExportConfirmModal) WasConfirmed() bool {
	return m.confirmed
}

// Show makes the modal visible and resets confirmation state.
func (m *ExportConfirmModal) Show() {
	m.visible = true
	m.confirmed = false
}

// Hide hides the modal without confirming.
func (m *ExportConfirmModal) Hide() {
	m.visible = false
	m.confirmed = false
}

// SetTheme sets the theme for the modal.
func (m *ExportConfirmModal) SetTheme(theme themes.Theme) {
	m.theme = theme
}

// SetDimensions sets the terminal dimensions for responsive sizing.
func (m *ExportConfirmModal) SetDimensions(width, height int) {
	m.width = width
	m.height = height
}

// GetDimensions returns the current dimensions.
func (m *ExportConfirmModal) GetDimensions() (width, height int) {
	return m.width, m.height
}

// GetArtifactType returns the artifact type being exported.
func (m *ExportConfirmModal) GetArtifactType() string {
	return m.artifactType
}

// GetFormat returns the export format.
func (m *ExportConfirmModal) GetFormat() string {
	return m.format
}

// GetDestination returns the export destination.
func (m *ExportConfirmModal) GetDestination() string {
	return m.destination
}
