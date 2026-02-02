package components

import (
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/containers"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
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
// This modal is purely for viewing. To edit or delete, use the 'e' or 'd' shortcuts
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
	event            *career.Event
	theme            themes.Theme
	visible          bool
	width            int
	height           int
	action           string
	viewport         viewport.Model
	ready            bool
	hasContent       bool
	showSkillsOption bool
}

// NewViewEventDetailModal creates a new event detail modal.
// By default, shows the "s: Skills" option in the footer.
// Use WithShowSkillsOption(false) to hide this option.
func NewViewEventDetailModal(event *career.Event, theme themes.Theme) *ViewEventDetailModal {
	// Ensure theme is not nil at initialization (UIKit pattern)
	if theme == nil {
		theme = themes.NewDefaultTheme()
	}
	return &ViewEventDetailModal{
		event:            event,
		theme:            theme,
		visible:          false,
		width:            80,
		height:           24,
		action:           "",
		showSkillsOption: true,
	}
}

// WithShowSkillsOption sets whether to show the "s: Skills" option in the footer.
// This should be set to false when viewing events from the ManageSkills workflow,
// since the user is already in a skills context.
func (m *ViewEventDetailModal) WithShowSkillsOption(show bool) *ViewEventDetailModal {
	m.showSkillsOption = show
	return m
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
		m.ready = false
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

	theme := m.theme

	// Calculate modal dimensions
	// Keep modal height reasonable: max 30 lines or 70% of terminal, whichever is smaller
	maxModalHeight := 30
	terminalMaxHeight := int(float64(m.height) * 0.7)
	if terminalMaxHeight < maxModalHeight {
		maxModalHeight = terminalMaxHeight
	}
	if maxModalHeight < 10 {
		maxModalHeight = 10
	}

	modalWidth := m.width - 12
	if modalWidth < 60 {
		modalWidth = 60
	}
	if modalWidth > 80 {
		modalWidth = 80
	}

	// Render event details
	content := RenderEventDetailCard(m.event, theme)
	contentLines := strings.Split(content, "\n")
	contentHeight := len(contentLines)

	// Calculate viewport height (modal height - borders - padding - footer)
	viewportHeight := maxModalHeight - 4
	if viewportHeight < 10 {
		viewportHeight = 10
	}

	// Initialize viewport if needed
	if !m.ready {
		vpWidth := modalWidth - 4
		if vpWidth < 10 {
			vpWidth = 10
		}
		if viewportHeight < 5 {
			viewportHeight = 5
		}
		m.viewport = viewport.New(vpWidth, viewportHeight)
		m.viewport.SetContent(content)
		m.hasContent = contentHeight > viewportHeight
		m.ready = true
	}

	// Build footer with UIKit primitives
	badges := []*primitives.Badge{
		primitives.HelpKeyBadge("↑↓/jk", "Scroll", theme),
	}
	if m.showSkillsOption {
		badges = append(badges, primitives.HelpKeyBadge("s", "Skills", theme))
	}
	badges = append(badges, primitives.HelpKeyBadge("Enter/Esc", "Close", theme))

	// Add scroll percentage if scrollable
	var scrollHint string
	if m.hasContent {
		percentScrolled := int(m.viewport.ScrollPercent() * 100)
		badges = append(badges, primitives.HelpKeyBadge(fmt.Sprintf("[%d%%]", percentScrolled), "", theme))
	}
	scrollHint = primitives.RenderHelpFooter(theme, badges...)

	// Build modal content
	modalContent := lipgloss.JoinVertical(lipgloss.Left, m.viewport.View(), "", scrollHint)

	// Wrap in styled box with solid background using UIKit
	return containers.NewBox(theme).
		Content(modalContent).
		Width(modalWidth).
		Height(maxModalHeight).
		Padding(2).
		Background(theme.BackgroundColor()).
		Render()
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
	m.ready = false
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
func (m *ViewEventDetailModal) SetEvent(event *career.Event) {
	m.event = event
	m.action = ""
}
