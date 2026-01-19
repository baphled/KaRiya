package components

import (
	"strings"

	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/charmbracelet/bubbles/viewport"
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
	event      *career.CareerEvent
	theme      themes.Theme
	visible    bool
	width      int
	height     int
	action     string // Always empty for read-only modal (kept for API compatibility)
	viewport   viewport.Model
	ready      bool
	hasContent bool
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

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = false // Force viewport recreation on resize
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "backspace", "enter", "q":
			// Close modal
			m.action = ""
			m.Hide()
			return m, nil

		// Scrolling keys - pass to viewport if content is scrollable
		case "up", "k", "down", "j", "pgup", "pgdown", "ctrl+u", "ctrl+d":
			if m.ready && m.hasContent {
				m.viewport, cmd = m.viewport.Update(msg)
				return m, cmd
			}
		}
	}

	return m, nil
}

// View renders the modal content with solid background and scrolling.
func (m *ViewEventDetailModal) View() string {
	if !m.visible {
		return ""
	}

	// Calculate modal dimensions
	// Keep modal height reasonable: max 30 lines or 70% of terminal, whichever is smaller
	maxModalHeight := 30
	terminalMaxHeight := int(float64(m.height) * 0.7)
	if terminalMaxHeight < maxModalHeight {
		maxModalHeight = terminalMaxHeight
	}
	if maxModalHeight < 10 {
		maxModalHeight = 10 // Minimum usable height
	}

	modalWidth := m.width - 12 // Leave margins
	if modalWidth < 60 {
		modalWidth = 60 // Ensure minimum readable width
	}
	if modalWidth > 80 {
		modalWidth = 80 // Max width for readability
	}

	// Render event details
	content := RenderEventDetailCard(m.event, m.theme)
	contentLines := strings.Split(content, "\n")
	contentHeight := len(contentLines)

	// Calculate viewport height (modal height - borders - padding - footer)
	viewportHeight := maxModalHeight - 4 // Account for border (2), footer (2)
	if viewportHeight < 10 {
		viewportHeight = 10 // Ensure minimum usable height
	}

	// Initialize viewport if needed
	if !m.ready {
		vpWidth := modalWidth - 4
		if vpWidth < 10 {
			vpWidth = 10 // Minimum viewport width
		}
		if viewportHeight < 5 {
			viewportHeight = 5 // Minimum viewport height
		}
		m.viewport = viewport.New(vpWidth, viewportHeight)
		m.viewport.SetContent(content)
		m.hasContent = contentHeight > viewportHeight
		m.ready = true
	}

	// Build footer with scroll indicator
	scrollHint := "Enter/Esc: Close"
	if m.hasContent {
		percentScrolled := int(m.viewport.ScrollPercent() * 100)
		scrollHint = lipgloss.NewStyle().
			Foreground(styles.ColorTextSecondary).
			Render("↑↓/j/k: Scroll | " + "Enter/Esc: Close " + lipgloss.NewStyle().Faint(true).Render("["+string(rune(percentScrolled/10+'0'))+string(rune(percentScrolled%10+'0'))+"%]"))
	} else {
		scrollHint = lipgloss.NewStyle().
			Foreground(styles.ColorTextSecondary).
			Render(scrollHint)
	}

	// Build modal content
	modalContent := lipgloss.JoinVertical(lipgloss.Left, m.viewport.View(), "", scrollHint)

	// Wrap in styled box with solid background
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Background(styles.ColorBackground).
		Padding(1, 2).
		Width(modalWidth).
		MaxHeight(maxModalHeight).
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
	m.ready = false // Reset viewport when showing
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
