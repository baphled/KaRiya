package primitives

import (
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/uikit/theme"
)

// ProgressBar is a theme-aware progress bar component for displaying progress or confidence levels.
// It supports customizable width, characters, labels, and percentage display.
//
// Example:
//
//	progressBar := primitives.NewProgressBar(0.75, theme).Width(20).ShowPercentage(true)
//	confidenceBar := primitives.ConfidenceBar(0.85, theme)
type ProgressBar struct {
	theme.Aware
	value          float64
	width          int
	showPercentage bool
	label          string
	filledChar     string
	emptyChar      string
}

// NewProgressBar creates a new progress bar with the given value (0.0 to 1.0) and theme.
// If theme is nil, the default theme is used.
// The value is clamped to the range [0.0, 1.0].
func NewProgressBar(value float64, th theme.Theme) *ProgressBar {
	if value < 0.0 {
		value = 0.0
	}
	if value > 1.0 {
		value = 1.0
	}

	pb := &ProgressBar{
		value:          value,
		width:          20,
		showPercentage: false,
		label:          "",
		filledChar:     "█",
		emptyChar:      "░",
	}
	if th != nil {
		pb.SetTheme(th)
	}
	return pb
}

// Width sets the width of the bar (number of characters).
// Returns the progress bar for method chaining.
func (pb *ProgressBar) Width(w int) *ProgressBar {
	if w > 0 {
		pb.width = w
	}
	return pb
}

// ShowPercentage controls whether to display the percentage value.
// Returns the progress bar for method chaining.
func (pb *ProgressBar) ShowPercentage(show bool) *ProgressBar {
	pb.showPercentage = show
	return pb
}

// Label sets a label prefix for the progress bar.
// Returns the progress bar for method chaining.
func (pb *ProgressBar) Label(label string) *ProgressBar {
	pb.label = label
	return pb
}

// FilledChar sets the character used for the filled portion.
// Returns the progress bar for method chaining.
func (pb *ProgressBar) FilledChar(char string) *ProgressBar {
	if char != "" {
		pb.filledChar = char
	}
	return pb
}

// EmptyChar sets the character used for the empty portion.
// Returns the progress bar for method chaining.
func (pb *ProgressBar) EmptyChar(char string) *ProgressBar {
	if char != "" {
		pb.emptyChar = char
	}
	return pb
}

// Render returns the styled progress bar as a string.
func (pb *ProgressBar) Render() string {
	var parts []string

	if pb.label != "" {
		parts = append(parts, pb.label)
	}

	filledWidth := int(pb.value * float64(pb.width))
	if filledWidth > pb.width {
		filledWidth = pb.width
	}
	emptyWidth := pb.width - filledWidth

	bar := "[" + strings.Repeat(pb.filledChar, filledWidth) + strings.Repeat(pb.emptyChar, emptyWidth) + "]"
	parts = append(parts, bar)

	if pb.showPercentage {
		percentage := int(pb.value * 100)
		parts = append(parts, fmt.Sprintf("%d%%", percentage))
	}

	return strings.Join(parts, " ")
}

// ConfidenceBar creates a progress bar configured for displaying confidence scores.
// It shows the percentage and uses the default width of 20.
func ConfidenceBar(value float64, th theme.Theme) *ProgressBar {
	return NewProgressBar(value, th).ShowPercentage(true)
}

// CompactBar creates a compact progress bar without label or percentage.
// Useful for inline display in tables or lists.
func CompactBar(value float64, width int, th theme.Theme) *ProgressBar {
	return NewProgressBar(value, th).Width(width).ShowPercentage(false)
}
