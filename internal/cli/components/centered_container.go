package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/baphled/kariya/internal/cli/styles"
)

// Margins represents spacing around content
type Margins struct {
	Top    int
	Right  int
	Bottom int
	Left   int
}

// CenteredContainer provides automatic centering of content
type CenteredContainer struct {
	content        string
	width          int
	height         int
	verticalCenter bool
	minMargins     Margins
	borderStyle    lipgloss.Style
	showBorder     bool
}

// NewCenteredContainer creates a new centered container
func NewCenteredContainer(width, height int) *CenteredContainer {
	return &CenteredContainer{
		content:        "",
		width:          width,
		height:         height,
		verticalCenter: false,
		minMargins: Margins{
			Top:    2,
			Right:  4,
			Bottom: 2,
			Left:   4,
		},
		borderStyle: lipgloss.NewStyle(),
		showBorder:  false,
	}
}

// SetContent sets the content to be centered
func (c *CenteredContainer) SetContent(content string) *CenteredContainer {
	c.content = content
	return c
}

// SetWidth sets the container width
func (c *CenteredContainer) SetWidth(width int) *CenteredContainer {
	c.width = width
	return c
}

// SetHeight sets the container height
func (c *CenteredContainer) SetHeight(height int) *CenteredContainer {
	c.height = height
	return c
}

// SetVerticalCenter enables vertical centering
func (c *CenteredContainer) SetVerticalCenter(center bool) *CenteredContainer {
	c.verticalCenter = center
	return c
}

// SetMinMargins sets minimum margins
func (c *CenteredContainer) SetMinMargins(margins Margins) *CenteredContainer {
	c.minMargins = margins
	return c
}

// SetBorderStyle sets the border style
func (c *CenteredContainer) SetBorderStyle(style lipgloss.Style) *CenteredContainer {
	c.borderStyle = style
	return c
}

// ShowBorder enables border display
func (c *CenteredContainer) ShowBorder(show bool) *CenteredContainer {
	c.showBorder = show
	return c
}

// Render renders the centered container
func (c *CenteredContainer) Render() string {
	if c.content == "" {
		return ""
	}

	// Calculate available width after margins
	availableWidth := c.width - c.minMargins.Left - c.minMargins.Right
	if availableWidth < 20 {
		availableWidth = c.width // Ignore margins if too small
	}

	// Split content into lines
	contentLines := strings.Split(c.content, "\n")

	// Center each line horizontally
	centeredLines := make([]string, len(contentLines))
	for i, line := range contentLines {
		// Remove any existing ANSI codes' width for proper centering
		visualWidth := lipgloss.Width(line)

		if visualWidth < availableWidth {
			padding := (availableWidth - visualWidth) / 2
			centeredLines[i] = strings.Repeat(" ", c.minMargins.Left+padding) + line
		} else {
			centeredLines[i] = strings.Repeat(" ", c.minMargins.Left) + line
		}
	}

	result := strings.Join(centeredLines, "\n")

	// Apply vertical centering if enabled
	if c.verticalCenter && c.height > 0 {
		contentHeight := len(centeredLines)
		if contentHeight < c.height {
			topPadding := (c.height - contentHeight) / 2
			if topPadding < c.minMargins.Top {
				topPadding = c.minMargins.Top
			}

			paddingLines := make([]string, topPadding)
			for i := range paddingLines {
				paddingLines[i] = ""
			}

			result = strings.Join(paddingLines, "\n") + "\n" + result
		}
	} else if c.minMargins.Top > 0 {
		// Just add top margin
		topPadding := make([]string, c.minMargins.Top)
		for i := range topPadding {
			topPadding[i] = ""
		}
		result = strings.Join(topPadding, "\n") + "\n" + result
	}

	// Apply border if enabled
	if c.showBorder {
		// Note: lipgloss styles are immutable, assignment is sufficient for copy
		borderStyle := c.borderStyle.
			Width(c.width - 4).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(styles.ColorBorder)

		result = borderStyle.Render(result)
	}

	return result
}

// RenderWithParts renders multiple content parts centered together
func (c *CenteredContainer) RenderWithParts(parts ...string) string {
	combined := strings.Join(parts, "\n")
	return c.SetContent(combined).Render()
}

// CenterText is a convenience function to center a single line of text
func CenterText(text string, width int) string {
	visualWidth := lipgloss.Width(text)
	if visualWidth >= width {
		return text
	}

	padding := (width - visualWidth) / 2
	return strings.Repeat(" ", padding) + text
}

// CenterBlock is a convenience function to center a block of text
func CenterBlock(content string, width, height int, verticalCenter bool) string {
	container := NewCenteredContainer(width, height)
	container.SetContent(content)
	container.SetVerticalCenter(verticalCenter)
	return container.Render()
}
