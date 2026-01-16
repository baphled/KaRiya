package containers

import (
	"strings"

	"github.com/baphled/kariya/internal/cli/uikit/theme"
	"github.com/charmbracelet/lipgloss"
)

// BoxVariant represents the visual style of a box.
type BoxVariant int

const (
	// BoxDefault is the standard box style.
	BoxDefault BoxVariant = iota
	// BoxEmphasized uses thick borders for important content.
	BoxEmphasized
	// BoxDestructive uses error colors for warnings/destructive actions.
	BoxDestructive
	// BoxSubtle uses muted colors for background content.
	BoxSubtle
)

// Box provides bordered container for modal frames, cards, panels.
//
// Usage:
//
//	box := containers.NewBox(theme).
//	    Title("Confirmation").
//	    Content("Are you sure?").
//	    Variant(containers.BoxDestructive).
//	    Width(50).
//	    Padding(2).
//	    WithShadow()
//	rendered := box.Render()
type Box struct {
	theme.Aware

	// Configuration
	content    string
	title      string
	variant    BoxVariant
	width      int // 0 = auto
	height     int // 0 = auto
	padding    int
	withShadow bool
}

// NewBox creates a new box with the given theme.
func NewBox(themeObj theme.Theme) *Box {
	box := &Box{
		variant: BoxDefault,
		padding: 1,
		width:   0,
		height:  0,
	}
	box.SetTheme(themeObj)
	return box
}

// Content sets the content to display inside the box.
func (b *Box) Content(content string) *Box {
	b.content = content
	return b
}

// Title sets the title displayed at the top of the box.
func (b *Box) Title(title string) *Box {
	b.title = title
	return b
}

// Variant sets the visual variant of the box.
func (b *Box) Variant(variant BoxVariant) *Box {
	b.variant = variant
	return b
}

// Width sets the width of the box (0 = auto).
func (b *Box) Width(width int) *Box {
	b.width = width
	return b
}

// Height sets the height of the box (0 = auto).
func (b *Box) Height(height int) *Box {
	b.height = height
	return b
}

// Padding sets the internal padding of the box.
func (b *Box) Padding(padding int) *Box {
	b.padding = padding
	return b
}

// WithShadow enables shadow rendering.
func (b *Box) WithShadow() *Box {
	b.withShadow = true
	return b
}

// Render returns the rendered box as a string.
func (b *Box) Render() string {
	// Get border style based on variant
	borderStyle := b.getBorderStyle()

	// Get border color based on variant
	borderColor := b.getBorderColor()

	// Build the base style
	style := lipgloss.NewStyle().
		Border(borderStyle).
		BorderForeground(borderColor).
		Padding(b.padding)

	// Apply width/height if set
	if b.width > 0 {
		style = style.Width(b.width - (b.padding * 2) - 2) // Subtract padding and borders
	}
	if b.height > 0 {
		style = style.Height(b.height - (b.padding * 2) - 2)
	}

	// Build content
	content := b.content
	if b.title != "" {
		titleStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(borderColor)
		content = titleStyle.Render(b.title) + "\n\n" + b.content
	}

	// Render the box
	rendered := style.Render(content)

	// Add shadow if enabled
	if b.withShadow {
		// Simple shadow: add gray background on bottom-right
		lines := strings.Split(rendered, "\n")
		shadowChar := "░"
		shadowColor := lipgloss.Color("#6C7086") // Catppuccin Overlay0

		for i := range lines {
			if i > 0 { // Skip first line
				lines[i] = lines[i] + lipgloss.NewStyle().Foreground(shadowColor).Render(shadowChar)
			}
		}
		// Add bottom shadow line
		if len(lines) > 0 {
			firstLineWidth := lipgloss.Width(lines[0])
			shadowLine := strings.Repeat(shadowChar, firstLineWidth+1)
			lines = append(lines, lipgloss.NewStyle().Foreground(shadowColor).Render(shadowLine))
		}
		rendered = strings.Join(lines, "\n")
	}

	return rendered
}

// getBorderStyle returns the lipgloss border style for the current variant.
func (b *Box) getBorderStyle() lipgloss.Border {
	switch b.variant {
	case BoxEmphasized:
		return lipgloss.ThickBorder()
	case BoxDefault, BoxDestructive, BoxSubtle:
		return lipgloss.RoundedBorder()
	default:
		return lipgloss.RoundedBorder()
	}
}

// getBorderColor returns the border color for the current variant.
func (b *Box) getBorderColor() lipgloss.Color {
	theme := b.Theme()
	if theme == nil {
		// Fallback colors if theme is not available
		switch b.variant {
		case BoxDestructive:
			return lipgloss.Color("#F38BA8") // Catppuccin Red
		case BoxSubtle:
			return lipgloss.Color("#6C7086") // Catppuccin Overlay0
		default:
			return lipgloss.Color("#CBA6F7") // Catppuccin Mauve
		}
	}

	// Use theme colors
	switch b.variant {
	case BoxDestructive:
		return lipgloss.Color(theme.ErrorColor())
	case BoxSubtle:
		return lipgloss.Color(theme.MutedColor())
	case BoxEmphasized:
		return lipgloss.Color(theme.PrimaryColor())
	default:
		return lipgloss.Color(theme.BorderColor())
	}
}
