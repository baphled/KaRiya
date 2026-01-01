package components

import (
	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/charmbracelet/lipgloss"
)

// PaddingMode defines the padding strategy for a screen container
type PaddingMode int

const (
	// PaddingNormal provides standard padding (1 vertical, 2 horizontal)
	PaddingNormal PaddingMode = iota
	// PaddingCompact provides minimal padding (0 vertical, 1 horizontal)
	PaddingCompact
	// PaddingSpacious provides generous padding (2 vertical, 4 horizontal)
	PaddingSpacious
)

// ScreenContainer wraps screen content with consistent margins and width constraints.
// It is a stateless rendering component that applies padding and width constraints
// to ensure consistent screen layout across the application.
type ScreenContainer struct {
	content          string
	paddingMode      PaddingMode
	customVertical   int
	customHorizontal int
	maxWidth         int
	useCustomPadding bool
}

// NewScreenContainer creates a new ScreenContainer with default padding (Normal mode).
func NewScreenContainer(content string) *ScreenContainer {
	return &ScreenContainer{
		content:     content,
		paddingMode: PaddingNormal,
		maxWidth:    styles.MaxWidth(120), // Use default max width
	}
}

// WithPaddingMode sets the padding mode for the container.
// This method uses the builder pattern to allow method chaining.
func (sc *ScreenContainer) WithPaddingMode(mode PaddingMode) *ScreenContainer {
	sc.paddingMode = mode
	sc.useCustomPadding = false
	return sc
}

// WithCustomPadding sets custom padding values (vertical, horizontal).
// This overrides any padding mode that was previously set.
// This method uses the builder pattern to allow method chaining.
func (sc *ScreenContainer) WithCustomPadding(vertical, horizontal int) *ScreenContainer {
	sc.customVertical = vertical
	sc.customHorizontal = horizontal
	sc.useCustomPadding = true
	return sc
}

// WithMaxWidth sets a custom maximum width for the container.
// This method uses the builder pattern to allow method chaining.
func (sc *ScreenContainer) WithMaxWidth(width int) *ScreenContainer {
	sc.maxWidth = width
	return sc
}

// Render returns the styled screen container as a string.
// It applies padding, width constraints, and text color based on the container's configuration.
func (sc *ScreenContainer) Render() string {
	// Determine padding based on mode or custom values
	vertical, horizontal := sc.getPadding()

	// Create base style with padding and width constraints
	baseStyle := lipgloss.NewStyle().
		Padding(vertical, horizontal).
		MaxWidth(sc.maxWidth).
		Foreground(styles.ColorTextPrimary)

	return baseStyle.Render(sc.content)
}

// getPadding returns the vertical and horizontal padding based on the current mode or custom values.
func (sc *ScreenContainer) getPadding() (int, int) {
	if sc.useCustomPadding {
		return sc.customVertical, sc.customHorizontal
	}

	switch sc.paddingMode {
	case PaddingCompact:
		return 0, 1
	case PaddingSpacious:
		return 2, 4
	case PaddingNormal:
		fallthrough
	default:
		return 1, 2
	}
}
