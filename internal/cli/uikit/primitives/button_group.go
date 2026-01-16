package primitives

import (
	"github.com/baphled/kariya/internal/cli/uikit/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ButtonGroup is a theme-aware container for multiple buttons with keyboard navigation.
// It supports both horizontal and vertical layouts and handles focus management.
//
// Example:
//
//	group := primitives.NewButtonGroup(theme)
//	group.AddPrimary("Save").AddSecondary("Cancel").AddDanger("Delete")
//	group.Update(tea.KeyMsg{Type: tea.KeyTab}) // Navigate with keyboard
//	rendered := group.Render()
type ButtonGroup struct {
	theme.Aware
	buttons      []*Button
	focusedIndex int
	horizontal   bool
}

// NewButtonGroup creates a new button group with the given theme.
// If theme is nil, the default theme is used.
// The group defaults to horizontal layout.
func NewButtonGroup(th theme.Theme) *ButtonGroup {
	g := &ButtonGroup{
		buttons:      []*Button{},
		focusedIndex: 0,
		horizontal:   true,
	}
	if th != nil {
		g.SetTheme(th)
	}
	return g
}

// Add adds a secondary button with the given label to the group.
// Returns the group for method chaining.
func (g *ButtonGroup) Add(label string) *ButtonGroup {
	btn := NewButton(label, g.Theme())
	g.buttons = append(g.buttons, btn)
	return g
}

// AddPrimary adds a primary button with the given label to the group.
// Returns the group for method chaining.
func (g *ButtonGroup) AddPrimary(label string) *ButtonGroup {
	btn := PrimaryButton(label, g.Theme())
	g.buttons = append(g.buttons, btn)
	return g
}

// AddSecondary adds a secondary button with the given label to the group.
// Returns the group for method chaining.
func (g *ButtonGroup) AddSecondary(label string) *ButtonGroup {
	btn := SecondaryButton(label, g.Theme())
	g.buttons = append(g.buttons, btn)
	return g
}

// AddDanger adds a danger button with the given label to the group.
// Returns the group for method chaining.
func (g *ButtonGroup) AddDanger(label string) *ButtonGroup {
	btn := DangerButton(label, g.Theme())
	g.buttons = append(g.buttons, btn)
	return g
}

// Horizontal sets whether buttons are laid out horizontally or vertically.
// Returns the group for method chaining.
func (g *ButtonGroup) Horizontal(h bool) *ButtonGroup {
	g.horizontal = h
	return g
}

// FocusIndex returns the index of the currently focused button.
func (g *ButtonGroup) FocusIndex() int {
	return g.focusedIndex
}

// FocusedLabel returns the label of the currently focused button.
// Returns empty string if there are no buttons.
func (g *ButtonGroup) FocusedLabel() string {
	if len(g.buttons) == 0 {
		return ""
	}
	return g.buttons[g.focusedIndex].label
}

// FocusNext moves focus to the next button, wrapping at the end.
func (g *ButtonGroup) FocusNext() {
	if len(g.buttons) == 0 {
		return
	}
	g.focusedIndex = (g.focusedIndex + 1) % len(g.buttons)
}

// FocusPrev moves focus to the previous button, wrapping at the beginning.
func (g *ButtonGroup) FocusPrev() {
	if len(g.buttons) == 0 {
		return
	}
	g.focusedIndex--
	if g.focusedIndex < 0 {
		g.focusedIndex = len(g.buttons) - 1
	}
}

// FocusFirst moves focus to the first button.
func (g *ButtonGroup) FocusFirst() {
	if len(g.buttons) > 0 {
		g.focusedIndex = 0
	}
}

// FocusLast moves focus to the last button.
func (g *ButtonGroup) FocusLast() {
	if len(g.buttons) > 0 {
		g.focusedIndex = len(g.buttons) - 1
	}
}

// Update handles keyboard input for navigation between buttons.
// Implements the tea.Model interface for Bubble Tea integration.
//
// Supported keys:
// - Tab / Right / 'l': Focus next button
// - Shift+Tab / Left / 'h': Focus previous button
// - Home: Focus first button
// - End: Focus last button
func (g *ButtonGroup) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab:
			g.FocusNext()
		case tea.KeyShiftTab:
			g.FocusPrev()
		case tea.KeyRight:
			g.FocusNext()
		case tea.KeyLeft:
			g.FocusPrev()
		case tea.KeyHome:
			g.FocusFirst()
		case tea.KeyEnd:
			g.FocusLast()
		case tea.KeyRunes:
			// Handle vim-style navigation
			if len(msg.Runes) > 0 {
				switch msg.Runes[0] {
				case 'l':
					g.FocusNext()
				case 'h':
					g.FocusPrev()
				}
			}
		}
	}
	return g, nil
}

// View returns the rendered button group as a string.
// This is an alias for Render() to support Bubble Tea's tea.Model interface.
func (g *ButtonGroup) View() string {
	return g.Render()
}

// Init implements tea.Model interface. Returns nil as no initialization is needed.
func (g *ButtonGroup) Init() tea.Cmd {
	return nil
}

// Render returns the styled button group as a string.
func (g *ButtonGroup) Render() string {
	if len(g.buttons) == 0 {
		return ""
	}

	// Render each button with appropriate focus state
	var renderedButtons []string
	for i, btn := range g.buttons {
		focused := i == g.focusedIndex
		rendered := btn.Focused(focused).Render()
		renderedButtons = append(renderedButtons, rendered)
	}

	// Join buttons based on layout
	if g.horizontal {
		return lipgloss.JoinHorizontal(lipgloss.Top, renderedButtons...)
	}
	return lipgloss.JoinVertical(lipgloss.Left, renderedButtons...)
}
