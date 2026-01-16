package containers

import (
	"github.com/baphled/kariya/internal/cli/uikit/theme"
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
type Box struct {
	theme.Aware
}

// NewBox creates a new box with the given theme.
func NewBox(themeObj theme.Theme) *Box {
	box := &Box{}
	box.SetTheme(themeObj)
	return box
}

// Content sets the content to display inside the box.
func (b *Box) Content(content string) *Box {
	return b
}

// Title sets the title displayed at the top of the box.
func (b *Box) Title(title string) *Box {
	return b
}

// Variant sets the visual variant of the box.
func (b *Box) Variant(variant BoxVariant) *Box {
	return b
}

// Width sets the width of the box (0 = auto).
func (b *Box) Width(width int) *Box {
	return b
}

// Height sets the height of the box (0 = auto).
func (b *Box) Height(height int) *Box {
	return b
}

// Padding sets the internal padding of the box.
func (b *Box) Padding(padding int) *Box {
	return b
}

// WithShadow enables shadow rendering.
func (b *Box) WithShadow() *Box {
	return b
}

// Render returns the rendered box as a string.
func (b *Box) Render() string {
	return ""
}
