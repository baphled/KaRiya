package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/cli/terminal"
)

// CenteringMode defines how content should be centered
type CenteringMode int

const (
	// CenterNone disables all centering
	CenterNone CenteringMode = iota
	// CenterHorizontal centers content horizontally only
	CenterHorizontal
	// CenterVertical centers content vertically only
	CenterVertical
	// CenterBoth centers content both horizontally and vertically
	CenterBoth
)

// OverflowMode defines how to handle content that exceeds available space
type OverflowMode int

const (
	// OverflowWrap wraps long lines to fit width
	OverflowWrap OverflowMode = iota
	// OverflowTruncate truncates content that's too long
	OverflowTruncate
	// OverflowScroll allows content to scroll (future enhancement)
	OverflowScroll
)

// SmartContainer provides intelligent content centering and layout
// based on terminal dimensions and content size
type SmartContainer struct {
	content        string
	terminalInfo   *terminal.Info
	centeringMode  CenteringMode
	overflowMode   OverflowMode
	customMargins  *Margins
	useMargins     bool
	minContentSize Size
	maxContentSize Size
	dimContent     bool
}

// Size represents width and height dimensions
type Size struct {
	Width  int
	Height int
}

// NewSmartContainer creates a new SmartContainer with the given terminal info
func NewSmartContainer(info *terminal.Info) *SmartContainer {
	return &SmartContainer{
		terminalInfo:  info,
		centeringMode: CenterBoth,
		overflowMode:  OverflowWrap,
		useMargins:    false, // Use responsive margins by default
	}
}

// SetContent sets the content to be rendered
func (s *SmartContainer) SetContent(content string) *SmartContainer {
	s.content = content
	return s
}

// SetCenteringMode sets how content should be centered
func (s *SmartContainer) SetCenteringMode(mode CenteringMode) *SmartContainer {
	s.centeringMode = mode
	return s
}

// SetOverflowMode sets how overflow content should be handled
func (s *SmartContainer) SetOverflowMode(mode OverflowMode) *SmartContainer {
	s.overflowMode = mode
	return s
}

// SetCustomMargins sets custom margins (overrides responsive margins)
func (s *SmartContainer) SetCustomMargins(margins Margins) *SmartContainer {
	s.customMargins = &margins
	s.useMargins = true
	return s
}

// SetMinContentSize sets minimum content dimensions
func (s *SmartContainer) SetMinContentSize(width, height int) *SmartContainer {
	s.minContentSize = Size{Width: width, Height: height}
	return s
}

// SetMaxContentSize sets maximum content dimensions
func (s *SmartContainer) SetMaxContentSize(width, height int) *SmartContainer {
	s.maxContentSize = Size{Width: width, Height: height}
	return s
}

// SetDimContent enables or disables content dimming (for modal backgrounds)
func (s *SmartContainer) SetDimContent(dim bool) *SmartContainer {
	s.dimContent = dim
	return s
}

// Render renders the container with appropriate layout strategy
func (s *SmartContainer) Render() string {
	// Handle terminal too small
	if !s.terminalInfo.CanRender(terminal.DefaultConfig) {
		return s.renderMinimalMode()
	}

	// Get margins (either custom or responsive)
	margins := s.getMargins()

	// Get available content area
	contentWidth, contentHeight := s.terminalInfo.ContentArea(margins)

	// Apply size constraints
	contentWidth = s.applyWidthConstraints(contentWidth)
	contentHeight = s.applyHeightConstraints(contentHeight)

	// Render content with appropriate strategy based on terminal size
	switch s.terminalInfo.GetCategory() {
	case terminal.SizeTiny:
		return s.renderTinyMode(contentWidth, contentHeight)
	case terminal.SizeCompact:
		return s.renderCompactMode(contentWidth, contentHeight)
	default:
		return s.renderNormalMode(contentWidth, contentHeight, margins)
	}
}

// renderNormalMode renders with full centering and margins
func (s *SmartContainer) renderNormalMode(width, height int, margins terminal.Margins) string {
	if s.content == "" {
		return ""
	}

	// Split content into lines
	lines := strings.Split(s.content, "\n")

	// Apply horizontal centering if needed
	if s.centeringMode == CenterHorizontal || s.centeringMode == CenterBoth {
		centeredLines := make([]string, len(lines))
		for i, line := range lines {
			centeredLines[i] = s.centerLine(line, width)
		}
		lines = centeredLines
	}

	// Join lines
	content := strings.Join(lines, "\n")

	// Apply vertical centering if needed
	if s.centeringMode == CenterVertical || s.centeringMode == CenterBoth {
		content = s.centerVertically(content, height)
	}

	// Apply margins
	style := lipgloss.NewStyle().
		Padding(margins.Top, margins.Right, margins.Bottom, margins.Left).
		Foreground(styles.ColorTextPrimary)

	// Apply dimming if enabled
	if s.dimContent {
		style = style.Foreground(lipgloss.Color("240")) // Gray color for dimming
	}

	return style.Render(content)
}

// renderCompactMode renders with reduced margins and simplified layout
func (s *SmartContainer) renderCompactMode(width, height int) string {
	if s.content == "" {
		return ""
	}

	// Simplified margins for compact mode
	compactMargins := terminal.Margins{Top: 1, Right: 2, Bottom: 1, Left: 2}

	// Apply basic centering
	lines := strings.Split(s.content, "\n")
	if s.centeringMode == CenterHorizontal || s.centeringMode == CenterBoth {
		centeredLines := make([]string, len(lines))
		for i, line := range lines {
			centeredLines[i] = s.centerLine(line, width)
		}
		lines = centeredLines
	}

	content := strings.Join(lines, "\n")

	style := lipgloss.NewStyle().
		Padding(compactMargins.Top, compactMargins.Right, compactMargins.Bottom, compactMargins.Left)

	return style.Render(content)
}

// renderTinyMode renders with minimal margins, no decorations
func (s *SmartContainer) renderTinyMode(width, height int) string {
	if s.content == "" {
		return ""
	}

	// Minimal margins for tiny terminals
	style := lipgloss.NewStyle().
		Width(width).
		MaxWidth(width).
		Padding(0, 1)

	return style.Render(s.content)
}

// renderMinimalMode renders fallback for terminals too small
func (s *SmartContainer) renderMinimalMode() string {
	msg := "Terminal too small\nMinimum: 40x15"
	return lipgloss.NewStyle().
		Foreground(styles.ColorError).
		Render(msg)
}

// centerLine centers a single line within the given width
func (s *SmartContainer) centerLine(line string, width int) string {
	visualWidth := lipgloss.Width(line)
	if visualWidth >= width {
		return line
	}

	padding := (width - visualWidth) / 2
	return strings.Repeat(" ", padding) + line
}

// centerVertically centers content vertically within the given height
func (s *SmartContainer) centerVertically(content string, height int) string {
	lines := strings.Split(content, "\n")
	contentHeight := len(lines)

	if contentHeight >= height {
		return content
	}

	topPadding := (height - contentHeight) / 2
	paddingLines := make([]string, topPadding)
	for i := range paddingLines {
		paddingLines[i] = ""
	}

	return strings.Join(paddingLines, "\n") + "\n" + content
}

// getMargins returns either custom or responsive margins
func (s *SmartContainer) getMargins() terminal.Margins {
	if s.useMargins && s.customMargins != nil {
		return terminal.Margins{
			Top:    s.customMargins.Top,
			Right:  s.customMargins.Right,
			Bottom: s.customMargins.Bottom,
			Left:   s.customMargins.Left,
		}
	}

	// Return responsive margins based on terminal size category
	switch s.terminalInfo.GetCategory() {
	case terminal.SizeTiny:
		return terminal.Margins{Top: 0, Right: 1, Bottom: 0, Left: 1}
	case terminal.SizeCompact:
		return terminal.Margins{Top: 1, Right: 2, Bottom: 1, Left: 2}
	case terminal.SizeLarge, terminal.SizeXLarge:
		return terminal.Margins{Top: 3, Right: 8, Bottom: 3, Left: 8}
	default:
		return terminal.Margins{Top: 2, Right: 4, Bottom: 2, Left: 4}
	}
}

// applyWidthConstraints applies min/max width constraints
func (s *SmartContainer) applyWidthConstraints(width int) int {
	if s.maxContentSize.Width > 0 && width > s.maxContentSize.Width {
		width = s.maxContentSize.Width
	}
	if s.minContentSize.Width > 0 && width < s.minContentSize.Width {
		width = s.minContentSize.Width
	}
	return width
}

// applyHeightConstraints applies min/max height constraints
func (s *SmartContainer) applyHeightConstraints(height int) int {
	if s.maxContentSize.Height > 0 && height > s.maxContentSize.Height {
		height = s.maxContentSize.Height
	}
	if s.minContentSize.Height > 0 && height < s.minContentSize.Height {
		height = s.minContentSize.Height
	}
	return height
}
