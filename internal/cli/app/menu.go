package app

import (
	"strings"

	"github.com/baphled/kariya/internal/cli/terminal"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
)

// Menu column widths for different terminal sizes.
// These provide responsive layout for the main menu table.
const (
	// Tiny terminal (< 80 cols).
	menuActionWidthTiny = 18
	menuDescWidthTiny   = 28

	// Compact terminal (80-99 cols).
	menuActionWidthCompact = 20
	menuDescWidthCompact   = 35

	// Normal terminal (100-119 cols).
	menuActionWidthNormal = 22
	menuDescWidthNormal   = 40

	// Large terminal (120-159 cols).
	menuActionWidthLarge = 25
	menuDescWidthLarge   = 50

	// Extra large terminal (160+ cols).
	menuActionWidthXLarge = 28
	menuDescWidthXLarge   = 60

	// Menu separator character.
	menuSeparator = "━"

	// Menu indicator for selected item.
	menuSelectedIndicator   = "▶ "
	menuUnselectedIndicator = "  "
)

// viewMenu renders the main menu using UIKit primitives for consistent centering.
func (m *Model) viewMenu() string {
	// Ensure terminalInfo has current dimensions.
	if !m.terminalInfo.IsValid && m.width > 0 && m.height > 0 {
		m.terminalInfo.Width = m.width
		m.terminalInfo.Height = m.height
		m.terminalInfo.IsValid = true
	}

	// Build menu components.
	var parts []string

	// 1. Logo (static view, no animation during menu).
	logoView := m.logo.ViewStatic()
	parts = append(parts, logoView)

	// 2. Spacing between logo and menu.
	parts = append(parts, "")
	parts = append(parts, "")

	// 3. Menu table with responsive columns.
	tableView := m.renderResponsiveTable()
	parts = append(parts, tableView)

	// 4. Spacing between menu and help.
	parts = append(parts, "")
	parts = append(parts, "")

	// 5. Help text.
	helpText := "↑/k Up  ↓/j Down  Enter Select  ? Help  q Quit"
	parts = append(parts, helpText)

	// Join all parts with center alignment and center in terminal.
	combined := primitives.JoinVertical(primitives.AlignCenter, parts...)
	return primitives.CenterInTerminal(combined, m.width, m.height)
}

// renderResponsiveTable creates the menu as simple text lines that can be centered.
func (m *Model) renderResponsiveTable() string {
	var lines []string

	// Calculate responsive column widths based on terminal size.
	actionWidth, descWidth := m.getMenuColumnWidths()

	// Create header row.
	headerAction := padRight("Action", actionWidth)
	headerDesc := padRight("Description", descWidth)
	header := primitives.Title(headerAction+"  "+headerDesc, m.theme).Render()
	lines = append(lines, header)

	// Add separator.
	separatorWidth := actionWidth + descWidth + 2
	separator := primitives.Muted(strings.Repeat(menuSeparator, separatorWidth), m.theme).Render()
	lines = append(lines, separator)

	// Create menu rows.
	for i, item := range m.menuItems {
		indicator := menuUnselectedIndicator
		var rowText string

		actionText := padRight(indicator+item.Name, actionWidth)
		descText := padRight(item.Help, descWidth)
		content := actionText + "  " + descText

		if i == m.selectedMenuIndex {
			indicator = menuSelectedIndicator
			actionText = padRight(indicator+item.Name, actionWidth)
			content = actionText + "  " + descText
			rowText = primitives.Title(content, m.theme).Render()
		} else {
			rowText = primitives.Body(content, m.theme).Render()
		}

		lines = append(lines, rowText)
	}

	return primitives.JoinVertical(primitives.AlignLeft, lines...)
}

// getMenuColumnWidths returns the appropriate column widths based on terminal size.
func (m *Model) getMenuColumnWidths() (actionWidth, descWidth int) {
	category := m.terminalInfo.GetCategory()

	switch category {
	case terminal.SizeTiny:
		return menuActionWidthTiny, menuDescWidthTiny
	case terminal.SizeCompact:
		return menuActionWidthCompact, menuDescWidthCompact
	case terminal.SizeNormal:
		return menuActionWidthNormal, menuDescWidthNormal
	case terminal.SizeLarge:
		return menuActionWidthLarge, menuDescWidthLarge
	default: // SizeXLarge
		return menuActionWidthXLarge, menuDescWidthXLarge
	}
}

// GetMenuItems returns the menu items from the model.
func (m *Model) GetMenuItems() []MenuItem {
	return m.menuItems
}
