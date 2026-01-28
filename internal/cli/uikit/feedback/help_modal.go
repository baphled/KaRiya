package feedback

import (
	"strings"

	"github.com/baphled/kariya/internal/cli/uikit/containers"
	"github.com/baphled/kariya/internal/cli/uikit/theme"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// HelpModal displays context-sensitive keyboard shortcuts in a modal overlay.
// It integrates with bubbles/help to auto-generate help text from KeyMaps.
//
// This is the UIKit version that uses theme-based styling exclusively.
type HelpModal struct {
	help       help.Model
	keyMap     help.KeyMap
	visible    bool
	showingAll bool
	width      int
	height     int
	theme      theme.Theme
}

// HelpModalKeyMap defines the keys for controlling the help modal itself.
type HelpModalKeyMap struct {
	Toggle   key.Binding
	Close    key.Binding
	FullHelp key.Binding
}

// DefaultHelpModalKeyMap returns the default key bindings for the help modal.
func DefaultHelpModalKeyMap() HelpModalKeyMap {
	return HelpModalKeyMap{
		Toggle: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		Close: key.NewBinding(
			key.WithKeys("esc", "?"),
			key.WithHelp("esc/?", "close help"),
		),
		FullHelp: key.NewBinding(
			key.WithKeys("f"),
			key.WithHelp("f", "full help"),
		),
	}
}

// NewHelpModal creates a new help modal with the given keymap.
func NewHelpModal(keyMap help.KeyMap) *HelpModal {
	h := help.New()
	h.ShowAll = false
	h.ShortSeparator = " • "
	h.FullSeparator = "   "

	return &HelpModal{
		help:       h,
		keyMap:     keyMap,
		visible:    false,
		showingAll: false,
		width:      80,
		height:     24,
		theme:      theme.Default(),
	}
}

// getTheme returns the theme or default if nil.
func (m *HelpModal) getTheme() theme.Theme {
	if m.theme != nil {
		return m.theme
	}
	return theme.Default()
}

// WithTheme sets the theme for the help modal.
// This method uses the builder pattern to allow method chaining.
func (m *HelpModal) WithTheme(t theme.Theme) *HelpModal {
	m.theme = t
	return m
}

// SetKeyMap updates the keymap displayed in the help modal
func (m *HelpModal) SetKeyMap(keyMap help.KeyMap) {
	m.keyMap = keyMap
}

// SetSize sets the available size for the help modal
func (m *HelpModal) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.help.Width = width - 8 // Account for padding/borders
}

// Show makes the help modal visible
func (m *HelpModal) Show() {
	m.visible = true
}

// Hide makes the help modal invisible
func (m *HelpModal) Hide() {
	m.visible = false
	m.showingAll = false
}

// Toggle toggles the help modal visibility
func (m *HelpModal) Toggle() {
	if m.visible {
		m.Hide()
	} else {
		m.Show()
	}
}

// IsVisible returns whether the help modal is currently visible
func (m *HelpModal) IsVisible() bool {
	return m.visible
}

// ToggleFullHelp toggles between short and full help display
func (m *HelpModal) ToggleFullHelp() {
	m.showingAll = !m.showingAll
	m.help.ShowAll = m.showingAll
}

// Update handles key events for the help modal
// Returns true if the event was consumed by the modal
func (m *HelpModal) Update(msg tea.Msg) (consumed bool, cmd tea.Cmd) {
	if !m.visible {
		// Check if help key was pressed to open
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			if key.Matches(keyMsg, DefaultHelpModalKeyMap().Toggle) {
				m.Show()
				return true, nil
			}
		}
		return false, nil
	}

	// Modal is visible - handle keys
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(keyMsg, DefaultHelpModalKeyMap().Close):
			m.Hide()
			return true, nil
		case key.Matches(keyMsg, DefaultHelpModalKeyMap().FullHelp):
			m.ToggleFullHelp()
			return true, nil
		}
	}

	return true, nil // Consume all events when modal is visible
}

// View renders the help modal
func (m *HelpModal) View() string {
	if !m.visible || m.keyMap == nil {
		return ""
	}

	t := m.getTheme()

	var content string
	if m.showingAll {
		content = m.help.FullHelpView(m.keyMap.FullHelp())
	} else {
		content = m.help.ShortHelpView(m.keyMap.ShortHelp())
	}

	// Build the modal content
	var sb strings.Builder

	// Title using theme colors
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(t.AccentColor()).
		MarginBottom(1)

	sb.WriteString(titleStyle.Render("Keyboard Shortcuts"))
	sb.WriteString("\n\n")

	// Help content
	sb.WriteString(content)
	sb.WriteString("\n\n")

	// Footer with toggle hint using theme colors
	footerStyle := lipgloss.NewStyle().
		Foreground(t.SecondaryColor()).
		Italic(true)

	if m.showingAll {
		sb.WriteString(footerStyle.Render("Press ? or esc to close • f for short help"))
	} else {
		sb.WriteString(footerStyle.Render("Press ? or esc to close • f for full help"))
	}

	// Modal container using UIKit Box with solid background
	modalWidth := m.width - 10
	if modalWidth < 40 {
		modalWidth = 40
	}

	return containers.NewBox(t).
		Content(sb.String()).
		Width(modalWidth).
		Padding(1).
		Background(t.BackgroundColor()).
		Render()
}

// ShortHelp returns a short help string for display in footer
func (m *HelpModal) ShortHelp() string {
	if m.keyMap == nil {
		return ""
	}
	return m.help.ShortHelpView(m.keyMap.ShortHelp())
}

// RenderOverlay renders the help modal as an overlay on top of existing content
func (m *HelpModal) RenderOverlay(baseContent string) string {
	if !m.visible {
		return baseContent
	}

	t := m.getTheme()
	modalContent := m.View()

	// Center the modal in the available space
	lines := strings.Split(baseContent, "\n")
	modalLines := strings.Split(modalContent, "\n")

	// Calculate vertical position (center)
	startY := (len(lines) - len(modalLines)) / 2
	if startY < 0 {
		startY = 0
	}

	// Build overlay
	var result strings.Builder
	for i, line := range lines {
		if i >= startY && i < startY+len(modalLines) {
			modalIdx := i - startY
			if modalIdx < len(modalLines) {
				// Center horizontally
				padding := (m.width - lipgloss.Width(modalLines[modalIdx])) / 2
				if padding < 0 {
					padding = 0
				}
				result.WriteString(strings.Repeat(" ", padding))
				result.WriteString(modalLines[modalIdx])
			}
		} else {
			// Dim the background content using theme
			dimStyle := lipgloss.NewStyle().Foreground(t.MutedColor())
			result.WriteString(dimStyle.Render(line))
		}
		if i < len(lines)-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}
