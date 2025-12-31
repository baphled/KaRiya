package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/baphled/kariya/internal/cli/styles"
)

// MenuItem represents a single menu item with label, shortcut, and description
type MenuItem struct {
	Label       string
	Shortcut    string
	Description string
	Action      string // Identifier for the action to take when selected
}

// MenuLayout specifies how the menu should be displayed
type MenuLayout string

const (
	LayoutVertical   MenuLayout = "vertical"   // Items stacked vertically (default)
	LayoutHorizontal MenuLayout = "horizontal" // Items displayed horizontally
)

// NavigationMenuModel is a reusable BubbleTea component for displaying navigation menus
type NavigationMenuModel struct {
	items          []MenuItem
	selectedIndex  int
	width          int
	height         int
	layout         MenuLayout
	showBorder     bool
	showNumbers    bool
	maxVisible     int // Maximum visible items (for vertical layout)
	scrollOffset   int // Offset for scrolling
	selectedAction string
	isSelected     bool
}

// NewNavigationMenu creates a new navigation menu with the given items
func NewNavigationMenu(items []MenuItem, layout MenuLayout) NavigationMenuModel {
	return NavigationMenuModel{
		items:         items,
		selectedIndex: 0,
		width:         80,
		height:        24,
		layout:        layout,
		showBorder:    true,
		showNumbers:   false,
		maxVisible:    10,
		scrollOffset:  0,
		isSelected:    false,
	}
}

// SetWidth sets the width of the menu
func (m *NavigationMenuModel) SetWidth(width int) {
	m.width = width
}

// SetHeight sets the height of the menu
func (m *NavigationMenuModel) SetHeight(height int) {
	m.height = height
}

// SetShowBorder sets whether to show a border around the menu
func (m *NavigationMenuModel) SetShowBorder(show bool) {
	m.showBorder = show
}

// SetShowNumbers sets whether to show numbers beside menu items
func (m *NavigationMenuModel) SetShowNumbers(show bool) {
	m.showNumbers = show
}

// SetMaxVisible sets the maximum number of visible items (for vertical layout)
func (m *NavigationMenuModel) SetMaxVisible(max int) {
	m.maxVisible = max
}

// GetSelectedAction returns the action string of the selected item
func (m NavigationMenuModel) GetSelectedAction() string {
	return m.selectedAction
}

// IsSelected returns whether an item was selected
func (m NavigationMenuModel) IsSelected() bool {
	return m.isSelected
}

// Reset resets the menu to its initial state
func (m *NavigationMenuModel) Reset() {
	m.selectedIndex = 0
	m.scrollOffset = 0
	m.isSelected = false
	m.selectedAction = ""
}

// Init implements the BubbleTea Model interface
func (m NavigationMenuModel) Init() tea.Cmd {
	return nil
}

// Update implements the BubbleTea Model interface
func (m NavigationMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.layout == LayoutVertical {
				m.moveUp()
			}
		case "down", "j":
			if m.layout == LayoutVertical {
				m.moveDown()
			}
		case "left", "h":
			if m.layout == LayoutHorizontal {
				m.moveUp()
			}
		case "right", "l":
			if m.layout == LayoutHorizontal {
				m.moveDown()
			}
		case "enter":
			if m.selectedIndex >= 0 && m.selectedIndex < len(m.items) {
				m.selectedAction = m.items[m.selectedIndex].Action
				m.isSelected = true
			}
		case "esc":
			m.isSelected = false
			m.selectedAction = ""
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

// View implements the BubbleTea Model interface
func (m NavigationMenuModel) View() string {
	if len(m.items) == 0 {
		return "No menu items"
	}

	var content string

	if m.layout == LayoutVertical {
		content = m.renderVertical()
	} else {
		content = m.renderHorizontal()
	}

	if m.showBorder {
		return styles.WithBorder(lipgloss.NewStyle()).Render(content)
	}

	return content
}

// renderVertical renders the menu with items stacked vertically
func (m NavigationMenuModel) renderVertical() string {
	var lines []string

	// Calculate visible range
	start := m.scrollOffset
	end := m.scrollOffset + m.maxVisible
	if end > len(m.items) {
		end = len(m.items)
	}

	// Ensure selected item is visible
	if m.selectedIndex < start {
		m.scrollOffset = m.selectedIndex
		start = m.scrollOffset
		end = start + m.maxVisible
		if end > len(m.items) {
			end = len(m.items)
		}
	} else if m.selectedIndex >= end {
		m.scrollOffset = m.selectedIndex - m.maxVisible + 1
		if m.scrollOffset < 0 {
			m.scrollOffset = 0
		}
		start = m.scrollOffset
		end = start + m.maxVisible
		if end > len(m.items) {
			end = len(m.items)
		}
	}

	for i := start; i < end; i++ {
		item := m.items[i]
		line := m.renderMenuItem(i, item, i == m.selectedIndex)
		lines = append(lines, line)
	}

	// Add scroll indicator if needed
	if end < len(m.items) {
		lines = append(lines, styles.InfoHint.Render("... more items"))
	}

	return strings.Join(lines, "\n")
}

// renderHorizontal renders the menu with items displayed horizontally
func (m NavigationMenuModel) renderHorizontal() string {
	var items []string

	for i, item := range m.items {
		rendered := m.renderMenuItem(i, item, i == m.selectedIndex)
		items = append(items, rendered)
	}

	// Join with separator
	separator := " " + styles.InfoHint.Render("|") + " "
	return strings.Join(items, separator)
}

// renderMenuItem renders a single menu item with appropriate styling
func (m NavigationMenuModel) renderMenuItem(index int, item MenuItem, isSelected bool) string {
	var result string

	// Add number if enabled
	if m.showNumbers {
		numStr := fmt.Sprintf("%d. ", index+1)
		numStyle := lipgloss.NewStyle().Foreground(styles.ColorTextSecondary)
		result += numStyle.Render(numStr)
	}

	// Add selection indicator
	if isSelected {
		selectedStyle := lipgloss.NewStyle().
			Foreground(styles.ColorAccentTeal).
			Bold(true)
		result += selectedStyle.Render("► " + item.Label)

		shortcutStyle := lipgloss.NewStyle().Foreground(styles.ColorSuccess)
		result += " " + shortcutStyle.Render("["+item.Shortcut+"]")
	} else {
		primaryStyle := lipgloss.NewStyle().Foreground(styles.ColorTextPrimary)
		result += primaryStyle.Render(item.Label)
		if item.Shortcut != "" {
			shortcutStyle := lipgloss.NewStyle().Foreground(styles.ColorTextMuted)
			result += " " + shortcutStyle.Render("["+item.Shortcut+"]")
		}
	}

	// Add description if present
	if item.Description != "" && isSelected {
		result += "\n"
		result += styles.InfoHint.Render("  " + item.Description)
	}

	return result
}

// moveUp moves the selection up
func (m *NavigationMenuModel) moveUp() {
	if m.selectedIndex > 0 {
		m.selectedIndex--
	} else {
		m.selectedIndex = len(m.items) - 1
	}
}

// moveDown moves the selection down
func (m *NavigationMenuModel) moveDown() {
	if m.selectedIndex < len(m.items)-1 {
		m.selectedIndex++
	} else {
		m.selectedIndex = 0
	}
}

// GetSelectedItem returns the currently selected item
func (m NavigationMenuModel) GetSelectedItem() *MenuItem {
	if m.selectedIndex >= 0 && m.selectedIndex < len(m.items) {
		return &m.items[m.selectedIndex]
	}
	return nil
}

// GetItems returns the menu items
func (m NavigationMenuModel) GetItems() []MenuItem {
	return m.items
}

// SetItems replaces the menu items
func (m *NavigationMenuModel) SetItems(items []MenuItem) {
	m.items = items
	if m.selectedIndex >= len(items) {
		m.selectedIndex = 0
	}
	m.scrollOffset = 0
}

// GetSelectedIndex returns the index of the selected item
func (m NavigationMenuModel) GetSelectedIndex() int {
	return m.selectedIndex
}

// SetSelectedIndex sets the selected item by index
func (m *NavigationMenuModel) SetSelectedIndex(index int) {
	if index >= 0 && index < len(m.items) {
		m.selectedIndex = index
	}
}
