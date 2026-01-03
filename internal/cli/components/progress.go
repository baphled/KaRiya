package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// ProgressBar represents a progress bar for operations with known duration.
type ProgressBar struct {
	label       string
	current     int
	total       int
	width       int
	showLabel   bool
	showPercent bool
}

// NewProgressBar creates a new progress bar with the given label and total.
func NewProgressBar(label string, total int) *ProgressBar {
	return &ProgressBar{
		label:       label,
		current:     0,
		total:       total,
		width:       40,
		showLabel:   true,
		showPercent: true,
	}
}

// SetCurrent updates the current progress value.
func (pb *ProgressBar) SetCurrent(current int) {
	if current < 0 {
		pb.current = 0
	} else if current > pb.total {
		pb.current = pb.total
	} else {
		pb.current = current
	}
}

// Increment increases the current progress by 1.
func (pb *ProgressBar) Increment() {
	pb.SetCurrent(pb.current + 1)
}

// IncrementBy increases the current progress by the given amount.
func (pb *ProgressBar) IncrementBy(amount int) {
	pb.SetCurrent(pb.current + amount)
}

// GetProgress returns the current progress value and total.
func (pb *ProgressBar) GetProgress() (int, int) {
	return pb.current, pb.total
}

// IsComplete returns true if progress has reached the total.
func (pb *ProgressBar) IsComplete() bool {
	return pb.current >= pb.total
}

// GetPercentage returns the progress as a percentage (0-100).
func (pb *ProgressBar) GetPercentage() int {
	if pb.total == 0 {
		return 0
	}
	return (pb.current * 100) / pb.total
}

// SetWidth sets the width of the progress bar in characters.
func (pb *ProgressBar) SetWidth(width int) {
	if width < 10 {
		pb.width = 10
	} else {
		pb.width = width
	}
}

// SetShowLabel sets whether to show the label.
func (pb *ProgressBar) SetShowLabel(show bool) {
	pb.showLabel = show
}

// SetShowPercent sets whether to show the percentage.
func (pb *ProgressBar) SetShowPercent(show bool) {
	pb.showPercent = show
}

// Render returns the rendered progress bar string.
func (pb *ProgressBar) Render() string {
	percentage := pb.GetPercentage()
	filledWidth := (pb.width * pb.current) / pb.total

	// Build the bar
	bar := ""
	for i := 0; i < pb.width; i++ {
		if i < filledWidth {
			bar += "█"
		} else {
			bar += "░"
		}
	}

	bar = lipgloss.NewStyle().
		Foreground(lipgloss.Color("42")).
		Render(bar[:filledWidth]) +
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			Render(bar[filledWidth:])

	// Build the result string
	result := ""

	if pb.showLabel {
		result += pb.label + " "
	}

	result += fmt.Sprintf("[%s]", bar)

	if pb.showPercent {
		result += fmt.Sprintf(" %3d%%", percentage)
	}

	// Add current/total if showing
	result += fmt.Sprintf(" (%d/%d)", pb.current, pb.total)

	return result
}

// ProgressIndicator represents a simple progress indicator for indeterminate operations.
type ProgressIndicator struct {
	frames    []string
	index     int
	label     string
	showLabel bool
}

// NewProgressIndicator creates a new progress indicator with the given label.
func NewProgressIndicator(label string) *ProgressIndicator {
	return &ProgressIndicator{
		frames: []string{
			"⠋",
			"⠙",
			"⠹",
			"⠸",
			"⠼",
			"⠴",
			"⠦",
			"⠧",
			"⠇",
			"⠏",
		},
		index:     0,
		label:     label,
		showLabel: true,
	}
}

// Next advances the indicator to the next frame.
func (pi *ProgressIndicator) Next() {
	pi.index = (pi.index + 1) % len(pi.frames)
}

// SetLabel sets the indicator label.
func (pi *ProgressIndicator) SetLabel(label string) {
	pi.label = label
}

// SetShowLabel sets whether to show the label.
func (pi *ProgressIndicator) SetShowLabel(show bool) {
	pi.showLabel = show
}

// GetFrame returns the current frame.
func (pi *ProgressIndicator) GetFrame() string {
	return pi.frames[pi.index]
}

// Render returns the rendered indicator string.
func (pi *ProgressIndicator) Render() string {
	frame := lipgloss.NewStyle().
		Foreground(lipgloss.Color("45")).
		Render(pi.GetFrame())

	result := frame

	if pi.showLabel {
		result += " " + pi.label
	}

	return result
}

// SimpleProgressBar renders a simple progress bar without animation.
// This is useful for static progress display.
func SimpleProgressBar(current, total int, width int) string {
	if width < 10 {
		width = 10
	}

	percentage := 0
	if total > 0 {
		percentage = (current * 100) / total
	}

	filledWidth := 0
	if total > 0 {
		filledWidth = (width * current) / total
	}

	bar := ""
	for i := 0; i < width; i++ {
		if i < filledWidth {
			bar += "█"
		} else {
			bar += "░"
		}
	}

	bar = lipgloss.NewStyle().
		Foreground(lipgloss.Color("42")).
		Render(bar[:filledWidth]) +
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			Render(bar[filledWidth:])

	return fmt.Sprintf("[%s] %3d%% (%d/%d)", bar, percentage, current, total)
}

// SimpleProgressIndicator renders a simple spinning indicator.
// This is useful for inline progress indication.
func SimpleProgressIndicator(index int) string {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	frame := frames[index%len(frames)]
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("45")).
		Render(frame)
}
