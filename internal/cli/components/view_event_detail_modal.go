package components

import (
	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ViewEventDetailModal displays event details in a modal overlay.
// Unlike the full-screen detail view, this shows the event info over the timeline list.
//
// Features:
// - Read-only display of event details
// - Close with Escape, backspace, Enter, or 'q'
// - Solid background to prevent transparency issues
// - Uses bubbletea-overlay for compositing
//
// Note: This modal is purely for viewing. To edit or delete, use the 'e' or 'd' shortcuts
// directly from the timeline list, or access edit/delete from the timeline context menu.
//
// Usage:
//
//	modal := NewViewEventDetailModal(event, theme)
//	modal.SetDimensions(width, height)
//	modal.Show()
//
//	// In Update:
//	cmd := modal.Update(msg)
//
//	// In View:
//	if modal.IsVisible() {
//	    return renderModalOverlay(modal, background)
//	}
type ViewEventDetailModal struct {
	event   *career.CareerEvent
	theme   themes.Theme
	visible bool
	width   int
	height  int
	action  string // Always empty for read-only modal (kept for API compatibility)
}

// NewViewEventDetailModal creates a new event detail modal.
func NewViewEventDetailModal(event *career.CareerEvent, theme themes.Theme) *ViewEventDetailModal {
	return &ViewEventDetailModal{
		event:   event,
		theme:   theme,
		visible: false,
		width:   80,
		height:  24,
		action:  "",
	}
}

// Init initializes the modal (implements tea.Model for bubbletea-overlay).
func (m *ViewEventDetailModal) Init() tea.Cmd {
	return nil
}

// Update handles keyboard input and window sizing.
func (m *ViewEventDetailModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.visible {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "backspace", "enter", "q":
			// Close modal
			m.action = ""
			m.Hide()
			return m, nil
		}
	}

	return m, nil
}

// View renders the modal content with solid background.
func (m *ViewEventDetailModal) View() string {
	if !m.visible {
		return ""
	}

	// Render event details using the existing component
	content := RenderEventDetailCard(m.event, m.theme)

	// Add footer with close hint
	footer := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		Render("Enter/Esc: Close")

	modalContent := lipgloss.JoinVertical(lipgloss.Left, content, "", footer)

	// Wrap in styled box with solid background to prevent transparency
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Background(styles.ColorBackground). // Solid background
		Padding(1, 2).
		MaxWidth(m.width - 8).   // Leave margins
		MaxHeight(m.height - 8). // Leave margins
		Render(modalContent)
}

// SetDimensions updates the modal's available dimensions.
func (m *ViewEventDetailModal) SetDimensions(width, height int) {
	m.width = width
	m.height = height
}

// Show makes the modal visible.
func (m *ViewEventDetailModal) Show() {
	m.visible = true
	m.action = ""
}

// Hide hides the modal.
func (m *ViewEventDetailModal) Hide() {
	m.visible = false
}

// IsVisible returns whether the modal is currently visible.
func (m *ViewEventDetailModal) IsVisible() bool {
	return m.visible
}

// GetAction returns the action selected by the user ("", "edit", or "delete").
// Should be called after the modal is hidden to determine what action to take.
func (m *ViewEventDetailModal) GetAction() string {
	return m.action
}

// SetEvent updates the event being displayed (useful for reusing the modal).
func (m *ViewEventDetailModal) SetEvent(event *career.CareerEvent) {
	m.event = event
	m.action = ""
}
