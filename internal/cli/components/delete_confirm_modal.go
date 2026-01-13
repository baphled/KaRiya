package components

import (
	"strings"

	"github.com/baphled/kariya/internal/cli/themes"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// DeleteConfirmModal is a reusable confirmation modal for delete operations.
// It provides a generic delete confirmation dialog that can be used for any entity type.
//
// Usage:
//
//	modal := components.NewDeleteConfirmModal("Go Programming", "Delete Skill", "Are you sure you want to delete 'Go Programming'?")
//	cmd := modal.Init()
//	// In Update:
//	cmd, confirmed := modal.Update(msg)
//	if confirmed {
//	    // User confirmed - proceed with deletion
//	} else if !modal.IsVisible() {
//	    // User cancelled - abort deletion
//	}
type DeleteConfirmModal struct {
	entityName string // Name of entity being deleted (for display)
	title      string // Modal title (e.g., "Delete Event", "Delete Skill")
	message    string // Confirmation message (e.g., "Are you sure you want to delete 'Go Programming'?")
	visible    bool   // Whether modal is currently visible
	confirmed  bool   // Whether user confirmed deletion
	theme      themes.Theme
	width      int
	height     int
}

// NewDeleteConfirmModal creates a new delete confirmation modal.
//
// Parameters:
//   - entityName: name of entity being deleted (shown in message)
//   - title: modal title (e.g., "Delete Event", "Delete Skill")
//   - message: confirmation message to display
//
// Returns a new DeleteConfirmModal ready to use.
func NewDeleteConfirmModal(entityName, title, message string) *DeleteConfirmModal {
	return &DeleteConfirmModal{
		entityName: entityName,
		title:      title,
		message:    message,
		visible:    true,
		confirmed:  false,
		theme:      themes.NewDefaultTheme(),
		width:      100,
		height:     24,
	}
}

// Init initializes the modal (required by BubbleTea lifecycle).
func (m *DeleteConfirmModal) Init() tea.Cmd {
	return nil
}

// Update handles keyboard input for the delete confirmation modal.
//
// Returns:
//   - tea.Cmd: command to execute (usually nil)
//   - bool: true if user confirmed deletion, false otherwise
//
// The modal closes on:
//   - 'y' or Enter: confirms deletion (returns true)
//   - 'n' or Esc: cancels deletion (returns false)
func (m *DeleteConfirmModal) Update(msg tea.Msg) (tea.Cmd, bool) {
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
			// User confirmed deletion
			m.confirmed = true
			m.visible = false
			return nil, true

		case "n", "N", "esc":
			// User cancelled deletion
			m.confirmed = false
			m.visible = false
			return nil, false
		}
	}

	return nil, false
}

// View renders the delete confirmation modal as a centered overlay.
//
// The modal displays:
//   - Title (e.g., "Delete Event")
//   - Confirmation message
//   - Action instructions (y: Confirm, n/Esc: Cancel)
//
// Returns empty string if modal is not visible.
func (m *DeleteConfirmModal) View() string {
	if !m.visible {
		return ""
	}

	// Build footer with KeyBadge components (Pattern #2 - Themed Footer Building)
	footer := RenderHelpFooter(m.theme,
		NewKeyBadge("y/Enter", "Confirm"),
		NewKeyBadge("n/Esc", "Cancel"),
	)

	// Build modal content
	var content strings.Builder

	// Title (styled as heading)
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(m.theme.ErrorColor()). // Red for delete warning
		MarginBottom(1)
	content.WriteString(titleStyle.Render(m.title))
	content.WriteString("\n\n")

	// Message (wrapped to modal width)
	messageStyle := lipgloss.NewStyle().
		Width(50).
		MarginBottom(1)
	content.WriteString(messageStyle.Render(m.message))
	content.WriteString("\n\n")

	// Footer
	content.WriteString(footer)

	// Create modal box (centered, with border)
	modalWidth := 60
	if modalWidth > m.width-4 {
		modalWidth = m.width - 4
	}
	if modalWidth < 40 {
		modalWidth = 40
	}

	modalStyle := lipgloss.NewStyle().
		Width(modalWidth).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.ErrorColor()).
		Background(m.theme.BackgroundColor()). // Add solid background for overlay
		Padding(1, 2).
		Align(lipgloss.Center)

	modalBox := modalStyle.Render(content.String())

	// Return just the modal box - bubbletea-overlay will handle positioning
	// No need to center it ourselves anymore
	return modalBox
}

// IsVisible returns whether the modal is currently visible.
func (m *DeleteConfirmModal) IsVisible() bool {
	return m.visible
}

// WasConfirmed returns whether the user confirmed the deletion.
// Only meaningful after the modal is closed.
func (m *DeleteConfirmModal) WasConfirmed() bool {
	return m.confirmed
}

// Show makes the modal visible.
func (m *DeleteConfirmModal) Show() {
	m.visible = true
	m.confirmed = false
}

// Hide hides the modal without confirming.
func (m *DeleteConfirmModal) Hide() {
	m.visible = false
	m.confirmed = false
}

// SetTheme sets the theme for the modal (useful for testing or custom themes).
func (m *DeleteConfirmModal) SetTheme(theme themes.Theme) {
	m.theme = theme
}

// SetDimensions sets the terminal dimensions for responsive sizing.
func (m *DeleteConfirmModal) SetDimensions(width, height int) {
	m.width = width
	m.height = height
}
