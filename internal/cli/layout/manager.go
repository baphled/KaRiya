package layout

import (
	"github.com/baphled/kariya/internal/cli/terminal"
)

// Manager handles responsive layout calculations based on terminal size
type Manager struct {
	terminalInfo *terminal.Info
	config       Config
}

// Config defines layout configuration and breakpoints
type Config struct {
	// Breakpoints for responsive behavior
	CompactBreakpoint int
	NormalBreakpoint  int
	LargeBreakpoint   int

	// Margin configurations per size category
	TinyMargins    terminal.Margins
	CompactMargins terminal.Margins
	NormalMargins  terminal.Margins
	LargeMargins   terminal.Margins
}

// DefaultConfig provides sensible layout defaults
var DefaultConfig = Config{
	CompactBreakpoint: 60,
	NormalBreakpoint:  80,
	LargeBreakpoint:   120,
	TinyMargins:       terminal.Margins{Top: 0, Right: 1, Bottom: 0, Left: 1},
	CompactMargins:    terminal.Margins{Top: 1, Right: 2, Bottom: 1, Left: 2},
	NormalMargins:     terminal.Margins{Top: 2, Right: 4, Bottom: 2, Left: 4},
	LargeMargins:      terminal.Margins{Top: 3, Right: 8, Bottom: 3, Left: 8},
}

// Rectangle represents a rectangular area
type Rectangle struct {
	X      int
	Y      int
	Width  int
	Height int
}

// NewManager creates a new layout manager with default configuration
func NewManager(info *terminal.Info) *Manager {
	return &Manager{
		terminalInfo: info,
		config:       DefaultConfig,
	}
}

// NewManagerWithConfig creates a layout manager with custom configuration
func NewManagerWithConfig(info *terminal.Info, config Config) *Manager {
	return &Manager{
		terminalInfo: info,
		config:       config,
	}
}

// GetContentArea calculates the available content area after margins
func (m *Manager) GetContentArea() Rectangle {
	margins := m.GetMargins()
	width, height := m.terminalInfo.GetSafeDimensions(terminal.DefaultConfig)

	contentWidth := width - margins.Left - margins.Right
	contentHeight := height - margins.Top - margins.Bottom

	// Ensure minimum content area
	contentWidth = max(contentWidth, 20)
	contentHeight = max(contentHeight, 5)

	return Rectangle{
		X:      margins.Left,
		Y:      margins.Top,
		Width:  contentWidth,
		Height: contentHeight,
	}
}

// GetMargins returns the appropriate margins based on terminal size category
func (m *Manager) GetMargins() terminal.Margins {
	category := m.terminalInfo.GetCategory()

	switch category {
	case terminal.SizeTiny:
		return m.config.TinyMargins
	case terminal.SizeCompact:
		return m.config.CompactMargins
	case terminal.SizeLarge, terminal.SizeXLarge:
		return m.config.LargeMargins
	default:
		return m.config.NormalMargins
	}
}

// CalculateColumns divides the content area into columns with gutters
// Returns an array of column widths
func (m *Manager) CalculateColumns(count int, gutterWidth int) []int {
	if count <= 0 {
		return []int{}
	}

	content := m.GetContentArea()
	totalGutter := gutterWidth * (count - 1)
	availableWidth := content.Width - totalGutter

	if availableWidth <= 0 {
		// Not enough space, return minimal widths
		columns := make([]int, count)
		for i := range columns {
			columns[i] = 1
		}
		return columns
	}

	columnWidth := availableWidth / count
	columns := make([]int, count)

	// Set base width for all columns
	for i := range columns {
		columns[i] = columnWidth
	}

	// Distribute remainder to first columns
	remainder := availableWidth - (columnWidth * count)
	for i := 0; i < remainder && i < count; i++ {
		columns[i]++
	}

	return columns
}

// ShouldUseCompactLayout returns true if terminal is too small for normal layout
func (m *Manager) ShouldUseCompactLayout() bool {
	category := m.terminalInfo.GetCategory()
	return category <= terminal.SizeCompact
}

// ShouldUseListLayout returns true if terminal is too narrow for multi-column layouts
func (m *Manager) ShouldUseListLayout() bool {
	if !m.terminalInfo.IsValid {
		return false
	}
	return m.terminalInfo.Width < 60
}

// UpdateTerminalInfo updates the terminal information
func (m *Manager) UpdateTerminalInfo(info *terminal.Info) {
	m.terminalInfo = info
}

// GetTerminalInfo returns the current terminal information
func (m *Manager) GetTerminalInfo() *terminal.Info {
	return m.terminalInfo
}

// max returns the larger of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
