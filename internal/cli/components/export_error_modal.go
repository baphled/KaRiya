package components

import (
	"strings"

	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/containers"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	tea "github.com/charmbracelet/bubbletea"
)

// ErrorResult represents the user's choice in the error modal.
type ErrorResult int

const (
	// ErrorResultNone means no choice was made yet.
	ErrorResultNone ErrorResult = iota
	// ErrorResultRetry means the user wants to retry the operation.
	ErrorResultRetry
	// ErrorResultCancel means the user wants to cancel.
	ErrorResultCancel
)

// ExportErrorModal displays an error message with retry/cancel options.
//
// Usage:
//
//	modal := components.NewExportErrorModal("export_failed", "Failed to write file")
//	// In Update:
//	cmd, result := modal.Update(msg)
//	if result == components.ErrorResultRetry {
//	    // User wants to retry
//	} else if result == components.ErrorResultCancel {
//	    // User wants to cancel
//	}
type ExportErrorModal struct {
	errorCode    string // Error code
	errorMessage string // Error message
	visible      bool
	wantsRetry   bool // Whether user chose retry
	theme        themes.Theme
	width        int
	height       int
}

// NewExportErrorModal creates a new export error modal.
//
// Parameters:
//   - errorCode: machine-readable error code (e.g., "export_failed")
//   - errorMessage: human-readable error message
func NewExportErrorModal(errorCode, errorMessage string) *ExportErrorModal {
	return &ExportErrorModal{
		errorCode:    errorCode,
		errorMessage: errorMessage,
		visible:      true,
		wantsRetry:   false,
		theme:        themes.NewDefaultTheme(),
		width:        100,
		height:       24,
	}
}

// Init initializes the modal (required by BubbleTea lifecycle).
func (m *ExportErrorModal) Init() tea.Cmd {
	return nil
}

// Update handles keyboard input for the export error modal.
//
// Returns:
//   - tea.Cmd: command to execute (usually nil)
//   - ErrorResult: the user's choice (None, Retry, or Cancel)
func (m *ExportErrorModal) Update(msg tea.Msg) (tea.Cmd, ErrorResult) {
	if !m.visible {
		return nil, ErrorResultNone
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return nil, ErrorResultNone

	case tea.KeyMsg:
		switch msg.String() {
		case "r", "R":
			// User wants to retry
			m.wantsRetry = true
			m.visible = false
			return nil, ErrorResultRetry

		case "esc", "q", "Q":
			// User wants to cancel
			m.wantsRetry = false
			m.visible = false
			return nil, ErrorResultCancel
		}
	}

	return nil, ErrorResultNone
}

// View renders the export error modal.
func (m *ExportErrorModal) View() string {
	if !m.visible {
		return ""
	}

	// Build footer with primitives showing keyboard shortcuts
	footer := primitives.RenderHelpFooter(m.theme,
		primitives.HelpKeyBadge("r", "Retry", m.theme),
		primitives.HelpKeyBadge("Esc/q", "Cancel", m.theme),
	)

	// Build modal content
	var content strings.Builder

	// Error title
	content.WriteString(primitives.ErrorText("Export Failed", m.theme).Bold().MarginBottom(1).Render())
	content.WriteString("\n\n")

	// Error message (wrapped to modal width)
	content.WriteString(primitives.Body(m.errorMessage, m.theme).Width(50).Render())
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

	// Use UIKit Box with destructive variant (error-colored border)
	return containers.NewBox(m.theme).
		Content(content.String()).
		Variant(containers.BoxDestructive).
		Width(modalWidth).
		Padding(2).
		Background(m.theme.BackgroundColor()).
		Render()
}

// IsVisible returns whether the modal is currently visible.
func (m *ExportErrorModal) IsVisible() bool {
	return m.visible
}

// WantsRetry returns whether the user chose to retry.
// Only meaningful after the modal is closed.
func (m *ExportErrorModal) WantsRetry() bool {
	return m.wantsRetry
}

// Show makes the modal visible and resets state.
func (m *ExportErrorModal) Show() {
	m.visible = true
	m.wantsRetry = false
}

// Hide hides the modal.
func (m *ExportErrorModal) Hide() {
	m.visible = false
}

// SetTheme sets the theme for the modal.
func (m *ExportErrorModal) SetTheme(theme themes.Theme) {
	m.theme = theme
}

// SetDimensions sets the terminal dimensions for responsive sizing.
func (m *ExportErrorModal) SetDimensions(width, height int) {
	m.width = width
	m.height = height
}

// GetDimensions returns the current dimensions.
func (m *ExportErrorModal) GetDimensions() (width, height int) {
	return m.width, m.height
}

// GetErrorCode returns the error code.
func (m *ExportErrorModal) GetErrorCode() string {
	return m.errorCode
}

// GetErrorMessage returns the error message.
func (m *ExportErrorModal) GetErrorMessage() string {
	return m.errorMessage
}
