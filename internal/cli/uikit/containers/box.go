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
	// BoxSuccess uses success colors for positive feedback.
	BoxSuccess
	// BoxWarning uses warning colors for caution messages.
	BoxWarning
	// BoxInfo uses accent colors for informational content.
	BoxInfo
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
//	    Background(lipgloss.Color("#1e1e2e")).
//	    WithShadow()
//	rendered := box.Render()
type Box struct {
	theme.Aware

	// Configuration
	content     string
	title       string
	variant     BoxVariant
	width       int // 0 = auto
	height      int // 0 = auto
	maxHeight   int // 0 = no max height
	padding     int
	withShadow  bool
	background  *lipgloss.Color // Optional background color
	borderColor *lipgloss.Color // Optional border color override
}

// NewBox creates a new box with the given theme.
//
// Expected:
//   - th must be a valid theme instance (can be nil).
//
// Returns:
//   - A fully initialized Box ready for use.
//
// Side effects:
//   - None.
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
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized Box ready for use.
//
// Side effects:
//   - None.
func (b *Box) Content(content string) *Box {
	b.content = content
	return b
}

// Title sets the title displayed at the top of the box.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized Box ready for use.
//
// Side effects:
//   - None.
func (b *Box) Title(title string) *Box {
	b.title = title
	return b
}

// Variant sets the visual variant of the box.
//
// Expected:
//   - boxvariant must be valid.
//
// Returns:
//   - A fully initialized Box ready for use.
//
// Side effects:
//   - None.
func (b *Box) Variant(variant BoxVariant) *Box {
	b.variant = variant
	return b
}

// Width sets the width of the box (0 = auto).
//
// Expected:
//   - int must be valid.
//
// Returns:
//   - A fully initialized Box ready for use.
//
// Side effects:
//   - None.
func (b *Box) Width(width int) *Box {
	b.width = width
	return b
}

// Height sets the height of the box (0 = auto).
//
// Expected:
//   - int must be valid.
//
// Returns:
//   - A fully initialized Box ready for use.
//
// Side effects:
//   - None.
func (b *Box) Height(height int) *Box {
	b.height = height
	return b
}

// Padding sets the internal padding of the box.
//
// Expected:
//   - int must be valid.
//
// Returns:
//   - A fully initialized Box ready for use.
//
// Side effects:
//   - None.
func (b *Box) Padding(padding int) *Box {
	b.padding = padding
	return b
}

// WithShadow enables shadow rendering.
//
// Returns:
//   - A fully initialized Box ready for use.
//
// Side effects:
//   - None.
func (b *Box) WithShadow() *Box {
	b.withShadow = true
	return b
}

// Background sets a solid background color for the box.
//
// Expected:
//   - color must be valid.
//
// Returns:
//   - A fully initialized Box ready for use.
//
// Side effects:
//   - None.
func (b *Box) Background(color lipgloss.Color) *Box {
	b.background = &color
	return b
}

// MaxHeight sets the maximum height of the box.
//
// Expected:
//   - int must be valid.
//
// Returns:
//   - A fully initialized Box ready for use.
//
// Side effects:
//   - None.
func (b *Box) MaxHeight(maxHeight int) *Box {
	b.maxHeight = maxHeight
	return b
}

// BorderColor sets a custom border color, overriding the variant color.
//
// Expected:
//   - color must be valid.
//
// Returns:
//   - A fully initialized Box ready for use.
//
// Side effects:
//   - None.
func (b *Box) BorderColor(color lipgloss.Color) *Box {
	b.borderColor = &color
	return b
}

// Render returns the rendered box as a string.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (b *Box) Render() string {
	// Get border style based on variant
	borderStyle := b.getBorderStyle()

	// Get border color - use custom if set, otherwise variant-based
	var borderColor lipgloss.Color
	if b.borderColor != nil {
		borderColor = *b.borderColor
	} else {
		borderColor = b.getBorderColor()
	}

	// Build the base style
	style := lipgloss.NewStyle().
		Border(borderStyle).
		BorderForeground(borderColor).
		Padding(b.padding)

	// Apply background if set
	if b.background != nil {
		style = style.Background(*b.background)
	}

	// Apply width/height if set
	if b.width > 0 {
		style = style.Width(b.width - (b.padding * 2) - 2) // Subtract padding and borders
	}
	if b.height > 0 {
		style = style.Height(b.height - (b.padding * 2) - 2)
	}
	if b.maxHeight > 0 {
		style = style.MaxHeight(b.maxHeight)
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
				lines[i] += lipgloss.NewStyle().Foreground(shadowColor).Render(shadowChar)
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
	default:
		return lipgloss.RoundedBorder()
	}
}

// getBorderColor returns the border color for the current variant.
func (b *Box) getBorderColor() lipgloss.Color {
	th := b.Theme()

	// Use theme colors
	switch b.variant {
	case BoxDestructive:
		return th.ErrorColor()
	case BoxSuccess:
		return th.SuccessColor()
	case BoxWarning:
		return th.WarningColor()
	case BoxInfo:
		return th.InfoColor()
	case BoxSubtle:
		return th.MutedColor()
	case BoxEmphasized:
		return th.PrimaryColor()
	default:
		return th.BorderColor()
	}
}
