package app

import (
	"strings"

	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/terminal"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/charmbracelet/lipgloss"
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

// viewMenu renders the main menu using pinned layout (logo at top, help at bottom).
func (m *Model) viewMenu() string {
	// Ensure terminalInfo has current dimensions.
	if !m.terminalInfo.IsValid && m.width > 0 && m.height > 0 {
		m.terminalInfo.Width = m.width
		m.terminalInfo.Height = m.height
		m.terminalInfo.IsValid = true
	}

	// Build menu sections for pinned layout.
	// Header: Logo
	logoView := m.logo.ViewStatic()

	// Content: Menu table with spacing
	tableView := m.renderResponsiveTable()
	content := primitives.JoinVertical(primitives.AlignCenter,
		"", "", // Spacing after logo
		tableView,
	)

	helpText := primitives.RenderHelpFooter(m.theme,
		primitives.NavigateBadge(m.theme),
		primitives.SelectBadge(m.theme),
		primitives.HelpKeyBadge(",", "Settings", m.theme),
		primitives.HelpBadge(m.theme),
		primitives.QuitBadge(m.theme),
	)

	// Calculate heights for spacer
	logoHeight := lipgloss.Height(logoView)
	contentHeight := lipgloss.Height(content)
	helpHeight := lipgloss.Height(helpText)

	logoMargin := 2
	helpMargin := 1
	spacerHeight := m.height - logoMargin - helpMargin - logoHeight - contentHeight - helpHeight
	if spacerHeight < 0 {
		spacerHeight = 0
	}

	allParts := []string{"", "", logoView, content}
	for range spacerHeight {
		allParts = append(allParts, "")
	}
	allParts = append(allParts, "")
	allParts = append(allParts, helpText)

	combined := primitives.JoinVertical(primitives.AlignCenter, allParts...)
	background := primitives.PlaceInTerminal(combined, m.width, m.height)

	if m.configModal != nil {
		m.configModal.SetDimensions(m.width, m.height)
		return behaviors.RenderModalOverlay(m.configModal, background)
	}

	return background
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
//
// Returns:
//   - A []MenuItem value.
//
// Side effects:
//   - None.
func (m *Model) GetMenuItems() []MenuItem {
	return m.menuItems
}
